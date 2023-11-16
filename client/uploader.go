package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func UploadFile(filePath string, targetURL string) error {
	// Open the file we want to send.
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Prepare a form that you will submit to the target URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Create a form field writer for the file part.
	fw, err := w.CreateFormFile("uploadFile", file.Name())
	if err != nil {
		return err
	}

	// Copy the file into the form field writer.
	if _, err = io.Copy(fw, file); err != nil {
		return err
	}

	// Close the multipart writer to set the terminating boundary.
	if err = w.Close(); err != nil {
		return err
	}

	// Create a request to send our multipart/ form-data.
	req, err := http.NewRequest("POST", targetURL, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request.
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check the response.
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return err
	}

	// Read the response body.
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	fmt.Println("Response from the server:", string(responseBody))

	return nil
}

func DownloadFile(fileNum string, targetURL string) error {
	// Create the file
	out, err := os.Create("test/data/" + fileNum + ".txt")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(targetURL + "?filename=" + fileNum)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
