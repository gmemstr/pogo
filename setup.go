package main

import (
	"archive/zip"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
)

func Setup() {
	defer LockFile()
	// Create users SQLite3 file
	fmt.Println("Initializing the database")

	os.MkdirAll("assets/config/", 0755)
	os.Create("assets/config/users.db")

	db, err := sql.Open("sqlite3", "assets/config/users.db")
	if err != nil {
		fmt.Sprintf("Problem opening database file! %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `users` ( `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, `username` TEXT UNIQUE, `hash` TEXT, `realname` TEXT, `email` TEXT, `permissions` INTEGER )")
	if err != nil {
		fmt.Sprintf("Problem creating database! %v", err)
	}

	db.Close()

	// Download web assets
	fmt.Println("Downloading web assets")
	os.MkdirAll("assets/web/", 0755)

	client := github.NewClient(nil).Repositories

	ctx := context.Background()
	res, _, err := client.GetLatestRelease(ctx, "gmemstr", "pogo")
	if err != nil {
		fmt.Sprintf("Problem creating database! %v", err)
	}
	for i := 0; i < len(res.Assets); i++ {
		if res.Assets[i].GetName() == "webassets.zip" {
			download := res.Assets[i]
			fmt.Sprintf("Release found: %v", download.GetBrowserDownloadURL())
			tmpfile, err := os.Create(download.GetName())
			if err != nil {
				fmt.Sprintf("Problem creating webassets file! %v", err)
			}
			var j io.Reader = (*os.File)(tmpfile)
			defer tmpfile.Close()

			j, s, err := client.DownloadReleaseAsset(ctx, "gmemstr", "pogo", download.GetID())
			if err != nil {
				fmt.Sprintf("Problem downloading webassets! %v", err)
			}
			if j == nil {
				resp, err := http.Get(s)
				defer resp.Body.Close()
				_, err = io.Copy(tmpfile, resp.Body)
				if err != nil {
					fmt.Sprintf("Problem creating webassets file! %v", err)
				}
				fmt.Println("Download complete\nUnzipping")
				err = Unzip(download.GetName(), "assets/web")
				defer os.Remove(download.GetName()) // Remove zip
			} else {
				fmt.Println("Unexpected error, please open an issue!")
			}
		}
	}
}

func LockFile() {
	lock, err := os.Create("run.lockfile")
	if err != nil {
		fmt.Println("Error: %v", err)
	}
	lock.Write([]byte("# You can leave this file emtpy, it simply tells Pogo you're set up!"))
	defer lock.Close()
}

// From https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
