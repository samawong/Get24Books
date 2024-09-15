package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Book struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	Year      int    `json:"year"`
	IpfsCid   string `json:"ipfs_cid"`
	Extension string `json:"extension"`
	FileSize  int    `json:"file_size"`
}

func main() {
	// Replace with your API URL
	//apiURL := "https://d.24hbook.store/search?query=+24hb&limit=20&offset=0"
	apiURL := ""
	downURL := "https://b.24hbook.store/ipfs"

	for i := 10000; i < 100000; i += 15 {
		apiURL = fmt.Sprintf("https://d.24hbook.store/search?query=+24hb&limit=15&offset=%d", i)
		resp, err := http.Get(apiURL)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			continue
		}

		var jsonResponse map[string]interface{}
		err = json.Unmarshal(body, &jsonResponse)
		if err != nil {
			fmt.Println("Error unmarshalling JSON response:", err)
			continue
		}

		for _, book := range jsonResponse["books"].([]interface{}) {
			bookMap := book.(map[string]interface{})
			bookStruct := Book{
				Title:     bookMap["title"].(string),
				Author:    bookMap["author"].(string),
				Publisher: bookMap["publisher"].(string),
				Year:      int(bookMap["year"].(float64)),
				IpfsCid:   bookMap["ipfs_cid"].(string),
				Extension: bookMap["extension"].(string),
			}
			fmt.Println("Book:", bookStruct)
			downloadFile(downURL+"/"+bookStruct.IpfsCid+"?filename="+bookStruct.Title+"."+bookStruct.Extension, bookStruct.Title+"."+bookStruct.Extension)
			time.Sleep(time.Second * 5)
		}

	}
}

func downloadFile(url string, destFile string) error {
	// Create the file.
	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data from the url.
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the data to file.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("The file", destFile, "has been downloaded")
	return nil
}
