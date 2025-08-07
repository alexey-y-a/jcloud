package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jcloud/bins"
	"jcloud/config"
	"jcloud/files"
	"jcloud/storage"
	"net/http"
)

const (
	baseURL = "https://api.jsonbin.io/v3/b"
)

type API struct {
	cfg     *config.Config
	storage storage.Storage
}

func New(cfg *config.Config, storage storage.Storage) *API {
	return &API{
		cfg:     cfg,
		storage: storage,
	}
}

func (a *API) makeRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, baseURL+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", a.cfg.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s", string(respBody))
	}

	return respBody, nil
}

func (a *API) CreateBin(filePath, name string) (*bins.Bin, error) {
	fs := files.FileSystem{}
	data, err := fs.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var content interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return nil, fmt.Errorf("file must contain valid JSON: %v", err)
	}

	resp, err := a.makeRequest("POST", "", map[string]interface{}{
		"name":    name,
		"private": true,
		"data":    content,
	})
	if err != nil {
		return nil, err
	}

	var result struct {
		Metadata bins.Bin `json:"metadata"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	binList, _ := a.storage.LoadBins("bins.json")
	binList.Bins = append(binList.Bins, result.Metadata)
	if err := a.storage.SaveBins(binList, "bins.json"); err != nil {
		return nil, fmt.Errorf("error saving bin: %v", err)
	}

	return &result.Metadata, nil
}
