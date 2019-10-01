// +build termux
// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package main

import (
	"os/exec"
)

func redeemCode(code string) error {
	args := []string{
		"-c", "New Crunchyroll Guest Pass ["+code+"]",
		"--action", "termux-open " + GPRURL + code,
	}

	c := exec.Command("termux-notification", args...)
	return c.Run()
}
