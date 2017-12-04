package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func setup() {
	fmt.Println("Initializing the database")
	os.MkdirAll("assets/config/", 0755)
}

func Unzip(file string, dest string) {

}
