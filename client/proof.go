package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/6l20/zama-test/common/merkle"
)

func Proof(index int, targetURL string) (*merkle.Proof, error) {
	var b bytes.Buffer
	
	// Create a request to send our multipart/ form-data.
	req, err := http.NewRequest("GET", targetURL + "/" + strconv.Itoa(index), &b)
	if err != nil {
		return nil, err
	}
	

	// Submit the request.
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check the response.
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return nil, err
	}

	// Read the response body.
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var proof merkle.Proof

	json.Unmarshal(responseBody, &proof)

	fmt.Println("Response from the server:", string(responseBody))

	return &proof, nil
}