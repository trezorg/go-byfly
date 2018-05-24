package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
)

type formValue struct {
	Name  string
	Value string
}

func readBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(data), nil
}

func prepareClient() (*http.Client, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			if len(via) == 0 {
				return nil
			}
			for attr, val := range via[0].Header {
				if _, ok := req.Header[attr]; !ok {
					req.Header[attr] = val
				}
			}
			return nil
		},
	}
	return client, nil
}

func postFormRequest(
	urlStr string,
	formValues []formValue,
	headers http.Header) (*http.Response, error) {

	client, err := prepareClient()
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	for _, formValue := range formValues {
		data.Set(formValue.Name, formValue.Value)
	}
	encodedData := data.Encode()
	body := bytes.NewBufferString(encodedData)
	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(encodedData)))
	for key, values := range headers {
		for _, header := range values {
			req.Header.Add(key, header)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, _ := readBody(resp)
		return nil, fmt.Errorf("Response status code: %d\n%s", resp.StatusCode, body)
	}
	req, err = http.NewRequest("GET", urlStr, nil)
	for key, values := range headers {
		for _, header := range values {
			req.Header.Add(key, header)
		}
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, _ := readBody(resp)
		return nil, fmt.Errorf("Response status code: %d\n%s", resp.StatusCode, body)
	}
	return resp, err
}
