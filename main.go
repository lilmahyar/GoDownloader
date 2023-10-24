package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {

	fmt.Println("Url : ")
	var fileURL string

	fileURL = "https://dl2.soft98.ir/soft/m/Mozilla.Firefox.119.0.EN.x64.zip?1698160384"

	// _, err := fmt.Scan(&fileURL)

	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	urlParsed, err := url.Parse(fileURL)
	if err != nil {
		fmt.Println("Error:", err)
	}

	path := urlParsed.Path
	splitedPath := strings.Split(path, "/")
	fileName := splitedPath[len(splitedPath)-1]

	tempFile, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Error:", err)
	}

	resp, err := http.Head(fileURL)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	_, err = io.Copy(tempFile, resp.Body)

	defer tempFile.Close()

	fmt.Println(fileName)
}
