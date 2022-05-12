// 4scrape
// Simple utility to download all files from a 4chan thread
// Usage: ./4scrape <url>
// Copyright (C) Fred Sanchez, 1999

package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	CLASS_NAME = "fileThumb"
)

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func HandleUser() (string, bool) {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./4scrape <url>")
	}
	path := os.Args[1]

	verbose := len(os.Args) == 3 && (os.Args[2] == "-v" || os.Args[2] == "--verbose")

	return path, verbose
}

func ParsePage(path string) *goquery.Selection {
	res, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Invalid status code: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc.Find("a")
}

func CreateAndGetNameOfDirectory(path string) string {
	parts := strings.Split(path, "/")
	dirname := parts[len(parts)-1]
	err := os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return dirname
}

func DownloadEntry(s *goquery.Selection, verbose bool, dirname string) {
	name, exists := s.Attr("class")
	if !exists || name != CLASS_NAME {
		return
	}
	href, _ := s.Attr("href")
	parts := strings.Split(href, "/")
	filename := parts[len(parts)-1]
	err := DownloadFile("./"+dirname+"/"+filename, "http:"+href)
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		log.Println("Downloaded " + "http:" + href + " as " + "./" + dirname + "/" + filename)
	}
}

func main() {
	path, verbose := HandleUser()

	entries := ParsePage(path)
	dirname := CreateAndGetNameOfDirectory(path)

	if verbose {
		log.Println("Downloading files to directory ./" + dirname)
	}

	entries.Each(func(_ int, s *goquery.Selection) {
		DownloadEntry(s, verbose, dirname)
	})
}
