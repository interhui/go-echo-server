package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	SERVER_VERSION = "Go Echo Server - 1.0"
)

var port string
var server string

func init() {
	flag.StringVar(&port,"port","4001","Echo Server Port")
	flag.StringVar(&server,"server",SERVER_VERSION,"Echo Server Name")
}

func main() {
	flag.Parse()

	http.HandleFunc("/", Echo)
	log.Println("Echo Server Start at", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))

}

func Echo(w http.ResponseWriter, r *http.Request) {

	request := make(map[string]interface{})
	request["client"] = r.RemoteAddr
	request["protocol"] = r.Proto
	request["method"] = r.Method
	request["url"] = r.URL.Path

	requestHeader := make(map[string]string)
	for headerKey, headerValues := range(r.Header) {
		headerValue := ArrayToString(headerValues)
		requestHeader[headerKey] = headerValue
	}
	request["header"] = requestHeader

	requestParam := make(map[string]string)
	for paramKey, paramValues := range(r.URL.Query()) {
		paramValue := ArrayToString(paramValues)
		requestParam[paramKey] = paramValue
	}
	request["parameter"] = requestParam

	if r.Method == "POST" || r.Method == "PUT" {
		body, _ := ioutil.ReadAll(r.Body)
		request["body"] = string(body)
	}

	current := time.Now()
	now := current.Format("2006-01-02T15:04:05")

	var status int
	var contentType string
	var contentLength int

	var content []byte
	var err error

	accept := r.Header.Get("Accept")
	if accept == "application/xml" {
		content, err = xml.Marshal(request)
	} else {
		content, err = json.Marshal(request)
	}

	if err == nil {
		contentType = "application/json"
		status = http.StatusOK

	} else {
		contentType = "text/plain"
		content = []byte(err.Error())
		status = http.StatusInternalServerError
	}

	contentLength = len(content)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Server", server)
	w.WriteHeader(status)
	w.Write(content)

	log.Printf("%s - [%s] \"%s %s %s\" %d %d %s\n", r.RemoteAddr, now, r.Method, r.URL, r.Proto, status, contentLength, r.Header.Get("User-Agent"))

}

func ArrayToString(array[] string) string {
	var str string
	for _, item_val := range(array) {
		str = str + "," + item_val
	}
	return str[1:]
}