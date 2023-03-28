package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	// magic
	A = "YzYwY2UwNDVkMjV"
	B = "iYzkwY2I1NjAyNm"
	C = "E4ZGQ2MjFlZWJlZ"
	D = "WY5OTVjYmVjYzUx"
	E = "OTUxMTkyZGE3NTM"
	F = "0OGM5NzdjZA=="

	// api
	ApiKey  = "ebe3d27ec52a35fc8d1835c6531c37bd72b7a54337666d5bd759379b72ae16f0"
	SiteURL = "https://api.letterboxd.com"
)

// I DON'T HANDLE AUTHENTICATION, if you must handle authentication, set the header yourself!
func signRequest(request *http.Request) {
	method := request.Method
	url := request.URL

	// get the query parameters of the url
	query := url.Query()

	// set api key parameter
	query.Set("apikey", ApiKey)
	// set nonce parameter
	query.Set("nonce", uuid.NewString())
	// set timestamp parameter
	query.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))

	//apply the parameters
	url.RawQuery = query.Encode()

	var body []byte
	if request.Body != nil {
		var err error
		body, err = io.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
	}

	request.Body = io.NopCloser(bytes.NewReader(body))

	// generate and set signature parameter
	query.Set("signature", signature(method, url.String(), body))

	//apply the signature parameter
	url.RawQuery = query.Encode()

	request.URL = url
}

func signature(method string, url string, body []byte) string {
	// get secret key using the silly little magic constants found in the Letterboxd APK
	key, _ := b64.StdEncoding.DecodeString(A + B + C + D + E + F)

	var data []byte

	// append all bytes in string "method" to the slice
	data = append(data, []byte(method)...)

	// append a '0' to separate data
	data = append(data, 0)

	// append all bytes in string "url" to the slice
	data = append(data, []byte(url)...)

	// append a '0' to separate data
	data = append(data, 0)

	// lastly, append the body byte array to the slice
	data = append(data, body...)

	// great! now let's sign this data with our key
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	result := mac.Sum(nil)

	return hex.EncodeToString(result)
}
