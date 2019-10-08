// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	// Guest passes URL
	GPURL = "https://www.crunchyroll.com/forumtopic-803801/the-official-guest-pass-thread-read-opening-post-first?pg=last" // nolint: lll

	// Guest passes redeem URL
	GPRURL = "https://www.crunchyroll.com/coupon_redeem?code="

	// Waiting time between requests
	WaitTime = time.Second * 15

	// Preallocated size for the database
	DBPreallocatedSize = 100
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
	go GetComments(req, codes)
	ListenForCodes(codes)
}

// GetCodes reads r and sends every code it finds to codes.
func GetCodes(r io.Reader, codes chan<- string) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("can't read response from Crunchyroll -> %w", err)
	}

	for _, comment := range bytes.Split(data, commentsPrefix)[1:] {
		comment = bytes.SplitN(comment, commentsSuffix, 2)[0]

		for _, code := range codesRE.FindAll(comment, -1) {
			codes <- string(code)
		}
	}

	return nil
}

// GetComments gets comments from Crunchyroll and looks for codes in them.
func GetComments(req *http.Request, codes chan<- string) {
	for {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("can't do the request -> %+v", err)
			goto wait
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("bad response from Crunchyroll -> %d", resp.StatusCode)
			goto closeRespBody
		}

		if err := GetCodes(resp.Body, codes); err != nil {
			log.Printf("%+v", err)
		}

	closeRespBody:
		resp.Body.Close()
	wait:
		time.Sleep(WaitTime)
	}
}

// ListenForCodes listens for new codes and sends them to RedeemCode.
func ListenForCodes(codes <-chan string) {
	db := make(map[string]struct{}, DBPreallocatedSize)

	for code := range codes {
		if _, ok := db[code]; ok {
			continue
		}

		if err := RedeemCode(code); err != nil {
			log.Printf("can't create notification -> %+v", err)
			continue
		}

		db[code] = struct{}{}
	}
}

// RedeemCode notifies when new codes are ready for redeem. If the current
// platform doesn't support notifications, it just opens the web browser for
// redeeming the code.
func RedeemCode(code string) error {
	return redeemCode(code)
}
