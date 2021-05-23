/**
2 * @Author: Nico
3 * @Date: 2021/5/21 22:01
4 */
package http

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

var httpClient = &http.Client{}

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(30) * time.Second,
	}
}

func Get(url string) (string, error){
	resp, err := httpClient.Get(url)
	if err != nil{
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}