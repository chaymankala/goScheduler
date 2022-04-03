package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gojektech/heimdall/httpclient"
	"github.com/pkg/errors"
)

func postHTTP(URI string, params string) (string, error) {
	timeout := 20000 * time.Millisecond
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	headers := http.Header{}
	// fmt.Println(vmURI + "?" + params)
	res, err := client.Get(URI, headers)
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 200 {
		return string(body), nil
	} else {
		return "", errors.Errorf("Some error occurred")
	}
}
