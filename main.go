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

type DataChunk struct {
	index int64
	chunk []byte
}

func main() {

	_, err := http.Get("https://www.google.com/")

	if err != nil {
		log.Fatal(err)
	}

	processorCount := 6

	fmt.Println("Url : ")
	var fileURL string

	fileURL = "https://dl.yasdl.com/2023/Software/Mozilla.Firefox.119.0.EN.x64_YasDL.com.rar?a"

	resp, err := http.Head(fileURL)

	if err != nil {
		log.Fatal(err)
	}

	urlParsed, err := url.Parse(fileURL)

	path := urlParsed.Path
	splitedPath := strings.Split(path, "/")
	fileName := splitedPath[len(splitedPath)-1]

	tempFile, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Error:", err)
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	fileSize := resp.ContentLength

	eachChunk := fileSize / int64(processorCount)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	dataChan := make(chan DataChunk)
	for i := 0; i < processorCount; i++ {
		go downloadChunk(fileURL, int64(i), eachChunk, dataChan)
	}

	counter := 0

	allData := make([][]byte, processorCount)

	for chunk := range dataChan {
		counter++

		allData[chunk.index] = chunk.chunk

		if counter == processorCount {
			close(dataChan)
		}
	}

	finalContent := []byte{}
	for _, data := range allData {
		finalContent = append(finalContent, data...)
	}

	_, err = tempFile.Write(finalContent)

	defer tempFile.Close()

	fmt.Println(fileName)
}

func downloadChunk(url string, index, size int64, downloadedChunk chan DataChunk) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	start := index * size

	contentRange := fmt.Sprintf("bytes=%d-%d", start, start+size-1)
	req.Header.Add("Range", contentRange)

	client := &http.Client{}

	resp, err := client.Do(req)

	body, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	downloadedChunk <- DataChunk{index: index, chunk: body}
}
