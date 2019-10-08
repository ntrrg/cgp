// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package main

import (
	"bytes"
	"fmt"
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

func main() {
	codes := make(chan string)

	go func(url string, codes chan<- string) {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Printf("can't create the request -> %+v", err)
			return
		}

		for {
			if err := GetCodes(req, codes); err != nil {
				log.Printf("can't get codes from Crunchyroll -> %+v", err)
			}

			time.Sleep(WaitTime)
		}
	}(GPURL, codes)

	db := make(map[string]struct{}, DBPreallocatedSize)

	for code := range codes {
		if _, ok := db[code]; ok {
			continue
		}

		if err := RedeemCode(code); err != nil {
			log.Printf("can't redeem the code [%s] -> %+v", code, err)
			continue
		}

		db[code] = struct{}{}
	}
}

// GetCodes gets comments from Crunchyroll and looks for codes in them.
func GetCodes(req *http.Request, codes chan<- string) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("can't do the request -> %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from Crunchyroll -> %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("can't read response from Crunchyroll -> %w", err)
	}

	codesRE := regexp.MustCompile("[A-Z0-9]{11}")
	commentsPrefix := []byte("showforumtopic-message-contents-text")
	commentsSuffix := []byte("</div>")

	for _, comment := range bytes.Split(data, commentsPrefix)[1:] {
		comment = bytes.SplitN(comment, commentsSuffix, 2)[0]

		for _, code := range codesRE.FindAll(comment, -1) {
			codes <- string(code)
		}
	}

	return nil
}

// RedeemCode notifies when new codes are ready for redeem. If the current
// platform doesn't support notifications, it just opens the web browser for
// redeeming the code.
func RedeemCode(code string) error {
	return redeemCode(code)
}
