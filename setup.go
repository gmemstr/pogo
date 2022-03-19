package main

import (
	"archive/zip"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
)

type Configuration struct {
	Name        string `json:"Name"`
	Host        string `json:"Host"`
	Email       string `json:"Email"`
	Description string `json:"Description"`
	Image       string `json:"Image"`
	PodcastUrl  string `json:"PodcastUrl"`
}

func Setup() {
	// Create directories
	os.MkdirAll("assets/config/", 0755)
	os.Mkdir("podcasts", 0755)

	// Write basic configuration file
	WriteSkeletonConfig()

	// Generate neccesary feed files
	GenerateRss()

	// Create "first run" lockfile when function exits
	// defer LockFile()

	// Create users SQLite3 file
	CreateDatabase()

	// Download web assets
	GetWebAssets()
}

func GetWebAssets() {
	fmt.Println("Downloading web assets")
	os.MkdirAll("assets/web/", 0755)

	client := github.NewClient(nil).Repositories

	ctx := context.Background()
	res, _, err := client.GetLatestRelease(ctx, "gmemstr", "pogo-vue")
	if err != nil {
		fmt.Println("Problem getting latest pogo-vue release! %v", err)
	}
	for i := 0; i < len(res.Assets); i++ {
		if res.Assets[i].GetName() == "webassets.zip" {
			download := res.Assets[i]
			fmt.Printf("Release found: %v\n", download.GetBrowserDownloadURL())
			tmpfile, err := os.Create(download.GetName())
			if err != nil {
				fmt.Printf("Problem creating webassets file! %v\n", err)
			}
			var j io.Reader = (*os.File)(tmpfile)
			defer tmpfile.Close()

			j, s, err := client.DownloadReleaseAsset(ctx, "gmemstr", "pogo-vue", download.GetID())
			if err != nil {
				fmt.Printf("Problem downloading webassets! %v\n", err)
			}
			if j == nil {
				resp, err := http.Get(s)
				defer resp.Body.Close()
				_, err = io.Copy(tmpfile, resp.Body)
				if err != nil {
					fmt.Printf("Problem creating webassets file! %v\n", err)
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

func CreateDatabase() {
	fmt.Println("Initializing the database")
	os.Create("assets/config/users.db")

	db, err := sql.Open("sqlite3", "assets/config/users.db")
	if err != nil {
		fmt.Println("Problem opening database file! %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `users` ( `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, `username` TEXT UNIQUE, `hash` TEXT, `realname` TEXT, `email` TEXT, `permissions` INTEGER )")
	if err != nil {
		fmt.Println("Problem creating database! %v", err)
	}

	text, err := GenerateRandomString(12)
	if err != nil {
		fmt.Println("Error randomly generating password", err)
	}
	fmt.Println("Admin password: ", text)
	hash, err := bcrypt.GenerateFromPassword([]byte(text), 4)
	if err != nil {
		fmt.Println("Error generating hash", err)
	}
	if bcrypt.CompareHashAndPassword(hash, []byte(text)) == nil {
		fmt.Println("Password hashed")
	}
	_, err = db.Exec("INSERT INTO users(id,username,hash,realname,email,permissions) VALUES (0,'admin','" + string(hash) + "','Administrator','admin@localhost',2)")
	if err != nil {
		fmt.Println("Problem creating database! %v", err)
	}
	defer db.Close()
}

func LockFile() {
	lock, err := os.Create(".lock")
	if err != nil {
		fmt.Println("Error: %v", err)
	}
	lock.Write([]byte("This file left intentionally empty"))
	defer lock.Close()
}

func WriteSkeletonConfig() {
	fmt.Println("Writing basic config file to disk")

	os.Create("assets/config/config.json")

	config := Configuration{
		"Pogo Podcast",
		"John Doe",
		"johndoe@localhost",
		"A Podcast About Stuff",
		"localhost:3000/assets/podcastimage.png",
		"http://localhost:3000",
	}
	c, err := json.Marshal(config)

	filename := "config.json"

	err = ioutil.WriteFile("./assets/config/"+filename, c, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
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
