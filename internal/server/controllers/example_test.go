package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/GrebenschikovDI/metalsys.git/internal/server/config"
	"github.com/GrebenschikovDI/metalsys.git/internal/server/storages"
)

func ExampleMetricsRouter() {
	storage := storages.NewMemStorage()
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error with config")
	}
	ctx := NewControllerContext(storage, *cfg)
	router := MetricsRouter(ctx)
	server := httptest.NewServer(router)
	defer server.Close()

	req, err := http.NewRequest("POST", server.URL+"/update/gauge/temp/23.5", nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println("Status Code:", resp.StatusCode)

	getReq, err := http.NewRequest("GET", server.URL+"/value/gauge/temp", nil)
	if err != nil {
		panic(err)
	}
	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		panic(err)
	}
	fmt.Println("Get Status Code:", getResp.StatusCode)

	// Output:
	// Status Code: 200
	// Get Status Code: 200

}
