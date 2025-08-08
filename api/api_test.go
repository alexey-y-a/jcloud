package api_test

import (
	"encoding/json"
	"fmt"
	"jcloud/api"
	"jcloud/bins"
	"jcloud/config"
	"jcloud/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type mockFileSystem struct {
	files map[string][]byte
}

func (m *mockFileSystem) ReadFile(filename string) ([]byte, error) {
	data, exists := m.files[filename]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", filename)
	}
	return data, nil
}

func (m *mockFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	m.files[filename] = data
	return nil
}

func (m *mockFileSystem) IsJson(filename string) bool {
	return filepath.Ext(filename) == ".json"
}

func setupTest(t *testing.T) (*api.API, *mockFileSystem, *httptest.Server, func()) {
	t.Helper()

	cfg := &config.Config{APIKey: "test-key"}
	mockFS := &mockFileSystem{files: make(map[string][]byte)}
	storage := storage.NewJsonStorage(mockFS)
	apiClient := api.New(cfg, storage)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	origClient := http.DefaultClient
	http.DefaultClient = server.Client()
	cleanup := func() {
		http.DefaultClient = origClient
		server.Close()
		for k := range mockFS.files {
			delete(mockFS.files, k)
		}
	}

	return apiClient, mockFS, server, cleanup
}

func TestCreateBin(t *testing.T) {
	apiClient, mockFS, server, cleanup := setupTest(t)
	defer cleanup()

	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/" {
			t.Fatalf("Expected POST to /, got %s to %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-Master-Key") != "test-key" {
			t.Fatalf("Expected X-Master-Key: test-key, got: %s", r.Header.Get("X-Master-Key"))
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody["name"] != "test-bin" || !reqBody["private"].(bool) {
			t.Fatalf("Unexpected request body: %+v", reqBody)
		}

		resp := struct {
			Metadata bins.Bin `json:"metadata"`
		}{
			Metadata: bins.Bin{
				ID:        "bin-123",
				Name:      "test-bin",
				Private:   true,
				CreatedAt: bins.TimeNow(),
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	testContent := `{"key": "value"}`
	mockFS.files["test.json"] = []byte(testContent)

	bin, err := apiClient.CreateBin("test.json", "test-bin")
	if err != nil {
		t.Fatalf("CreateBin failed: %v", err)
	}

	if bin.ID != "bin-123" || bin.Name != "test-bin" || !bin.Private {
		t.Fatalf("Unexpected bin: %+v", bin)
	}

	binList, err := apiClient.ListBins()
	if err != nil {
		t.Fatalf("Failed to load bins: %v", err)
	}
	if len(binList.Bins) != 1 || binList.Bins[0].ID != "bin-123" {
		t.Fatalf("Expected one bin in storage with ID bin-123, got: %+v", binList)
	}
}

func TestUpdateBin(t *testing.T) {
	apiClient, mockFS, server, cleanup := setupTest(t)
	defer cleanup()

	binList := bins.BinList{
		Bins: []bins.Bin{
			{ID: "bin-123", Name: "test-bin", Private: true, CreatedAt: bins.TimeNow()},
		},
	}
	data, _ := json.Marshal(binList)
	mockFS.files["bins.json"] = data

	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" || r.URL.Path != "/bin-123" {
			t.Fatalf("Expected PUT to /bin-123, got %s to %s", r.Method, r.URL.Path)
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if _, ok := reqBody["data"]; !ok {
			t.Fatalf("Expected data in request body, got: %+v", reqBody)
		}

		resp := struct {
			Metadata bins.Bin `json:"metadata"`
		}{
			Metadata: bins.Bin{
				ID:        "bin-123",
				Name:      "updated-bin",
				Private:   true,
				CreatedAt: bins.TimeNow(),
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	testContent := `{"new-key": "new-value"}`
	mockFS.files["test.json"] = []byte(testContent)

	bin, err := apiClient.UpdateBin("test.json", "bin-123")
	if err != nil {
		t.Fatalf("UpdateBin failed: %v", err)
	}

	if bin.ID != "bin-123" || bin.Name != "updated-bin" {
		t.Fatalf("Unexpected bin: %+v", bin)
	}

	binList, err = apiClient.ListBins()
	if err != nil {
		t.Fatalf("Failed to load bins: %v", err)
	}
	if len(binList.Bins) != 1 || binList.Bins[0].Name != "updated-bin" {
		t.Fatalf("Expected updated bin in storage, got: %+v", binList)
	}
}

func TestGetBin(t *testing.T) {
	apiClient, _, server, cleanup := setupTest(t)
	defer cleanup()

	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/bin-123" {
			t.Fatalf("Expected GET to /bin-123, got %s to %s", r.Method, r.URL.Path)
		}

		resp := struct {
			Record   interface{} `json:"record"`
			Metadata bins.Bin    `json:"metadata"`
		}{
			Metadata: bins.Bin{
				ID:        "bin-123",
				Name:      "test-bin",
				Private:   true,
				CreatedAt: bins.TimeNow(),
			},
			Record: map[string]string{"key": "value"},
		}
		json.NewEncoder(w).Encode(resp)
	})

	bin, err := apiClient.GetBin("bin-123")
	if err != nil {
		t.Fatalf("GetBin failed: %v", err)
	}

	if bin.ID != "bin-123" || bin.Name != "test-bin" || !bin.Private {
		t.Fatalf("Unexpected bin: %+v", bin)
	}
}

func TestDeleteBin(t *testing.T) {
	apiClient, mockFS, server, cleanup := setupTest(t)
	defer cleanup()

	binList := bins.BinList{
		Bins: []bins.Bin{
			{ID: "bin-123", Name: "test-bin", Private: true, CreatedAt: bins.TimeNow()},
		},
	}
	data, _ := json.Marshal(binList)
	mockFS.files["bins.json"] = data

	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" || r.URL.Path != "/bin-123" {
			t.Fatalf("Expected DELETE to /bin-123, got %s to %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := apiClient.DeleteBin("bin-123")
	if err != nil {
		t.Fatalf("DeleteBin failed: %v", err)
	}

	binList, err = apiClient.ListBins()
	if err != nil {
		t.Fatalf("Failed to load bins: %v", err)
	}
	if len(binList.Bins) != 0 {
		t.Fatalf("Expected empty bin list, got: %+v", binList)
	}
}

func TestListBins(t *testing.T) {
	apiClient, mockFS, _, cleanup := setupTest(t)
	defer cleanup()

	binList := bins.BinList{
		Bins: []bins.Bin{
			{ID: "bin-123", Name: "test-bin-1", Private: true, CreatedAt: bins.TimeNow()},
			{ID: "bin-456", Name: "test-bin-2", Private: false, CreatedAt: bins.TimeNow()},
		},
	}
	data, _ := json.Marshal(binList)
	mockFS.files["bins.json"] = data

	result, err := apiClient.ListBins()
	if err != nil {
		t.Fatalf("ListBins failed: %v", err)
	}

	if len(result.Bins) != 2 {
		t.Fatalf("Expected 2 bins, got %d: %+v", len(result.Bins), result)
	}
	if result.Bins[0].ID != "bin-123" || result.Bins[0].Name != "test-bin-1" {
		t.Fatalf("Unexpected first bin: %+v", result.Bins[0])
	}
	if result.Bins[1].ID != "bin-456" || result.Bins[1].Name != "test-bin-2" {
		t.Fatalf("Unexpected second bin: %+v", result.Bins[1])
	}
}
