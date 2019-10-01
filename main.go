// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	// Guest passes URL
	GPURL = "https://www.crunchyroll.com/forumtopic-803801/the-official-guest-pass-thread-read-opening-post-first?pg=last"

	// Guest passes redeem URL
	GPRURL = "https://www.crunchyroll.com/coupon_redeem?code="
)

var (
	codesRE = regexp.MustCompile("[A-Z0-9]{11}")

	commentsPrefix = []byte("showforumtopic-message-contents-text")
	commentsSuffix = []byte("</div>")
)

func main() {
	req, err := http.NewRequest(http.MethodGet, GPURL, nil)
	if err != nil {
		log.Printf("can't create the request -> %+v", err)
		return
	}

	codes := make(chan string)
	go ListenForCodes(codes)

	for {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("can't do the request -> %+v", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("bad response from Crunchyroll -> %d", resp.StatusCode)
			continue
		}

		go GetCodes(resp.Body, codes)
		time.Sleep(time.Second * 30)
	}
}

func GetCodes(rc io.ReadCloser, codes chan<- string) {
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Printf("can't read response from Crunchyroll -> %+v", err)
		return
	}

	for _, comment := range bytes.Split(data, commentsPrefix)[1:] {
		comment = bytes.SplitN(comment, commentsSuffix, 2)[0]
		comment = bytes.Trim(comment, " ")

		for _, code := range codesRE.FindAll(comment, -1) {
			codes <- string(code)
		}
	}
}

func ListenForCodes(codes <-chan string) {
	db := make(map[string]struct{}, 1000)

	for code := range codes {
		if _, ok := db[code]; ok {
			continue
		}

		if err := RedeemCode(code); err != nil {
			log.Printf("can't open the web browser -> %+v", err)
			continue
		}

		db[code] = struct{}{}
	}
}

func RedeemCode(code string) error {
	return redeemCode(code)
}
