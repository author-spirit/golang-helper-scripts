package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const CHUNK_SIZE_IN_BYTES = 10 * 1024 * 1024 // 10 MB in bytes
const DOWNLOAD_PATH = "downloads/"

//export download_file_chunks
func downloadFileInChunks(sourceFile string, destinationFile string, accessToken string) bool {
	if strings.TrimSpace(sourceFile) == "" || strings.TrimSpace(destinationFile) == "" {
		fmt.Println("Failed to download")
		return false
	}

	req, err := http.NewRequest("GET", sourceFile, nil)
	if err != nil {
		fmt.Println("Failed to Create Request Wrap")
		return false
	}

	// If Access Token is provided then apply Bearer token
	// Access Token is mainly for Securely Stored files
	// FUTURE Readers: extend additional options if required
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	client := &http.Client{}
	fileRequest, err := client.Do(req)

	if err != nil {
		fmt.Print("Error Performing Request", err)
		return false
	}

	// close the stream after return
	defer fileRequest.Body.Close()

	destFile, err := os.Create(destinationFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return false
	}

	// close the new created file
	defer destFile.Close()

	// Create buffer for chunks, 1
	bufReader := bufio.NewReader(fileRequest.Body)
	chunk := make([]byte, CHUNK_SIZE_IN_BYTES)

	for {
		n, err := bufReader.Read(chunk)
		if err != nil {
			if err != io.EOF {
				fmt.Println("failed to read chunk: %w", err)
				return false
			}
			break // EOF reached, exit loop
		}

		// Write the read bytes only
		_, err = destFile.Write(chunk[:n])
		if err != nil {
			fmt.Println("error writing chunk: %w", err)
			return false
		}
	}

	fmt.Println("\nDownload Complete", destinationFile)
	return true
}

func getFileName(fullPath string) (string, error) {
	fileParts := strings.Split(fullPath, "/")
	if len(fileParts) == 0 {
		// get the last part
		return "", nil
	}
	return DOWNLOAD_PATH + fileParts[len(fileParts)-1], nil
}

func main() {

	_, dirErr := os.ReadDir(DOWNLOAD_PATH)
	if dirErr != nil {
		log.Fatal(dirErr)
		os.Exit(0)
	}

	// fmt.Println("(" + strings.TrimSpace(" . ") + ")")

	// Enabled Profiling if you want to know performance of download

	// f, _ := os.Create("profiling.prof")
	// mem, err := os.Create("mem.prof")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return

	// }
	// pprof.StartCPUProfile(f)

	// defer f.Close()
	// defer mem.Close()

	// defer pprof.StopCPUProfile()

	start := time.Now()

	url := "<PROVIDE REMOTE FILE PATH>"
	listOfFiles := [3]string{"FILE_1", "FILE_2", "FILE_3"}

	for _, filename := range listOfFiles {

		if filename == "" {
			continue
		}

		filename = url + filename
		destFile, err := getFileName(filename)
		if err != nil {
			fmt.Println("Failed to get filename", err)
			return
		}
		downloadFileInChunks(filename, destFile, "")
	}

	elapsed := time.Since(start)

	// fmt.Println("Memory allocated:", endMem-startMem, "bytes")
	fmt.Println("Elapsed Time : ", elapsed)

	// Only enable when pprof is enabled
	// pprof.WriteHeapProfile(mem)
}
