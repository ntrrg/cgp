// +build !android
// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package main

import (
	"os/exec"
)

func redeemCode(code string) error {
	args := []string{
		"call", "--session", "--dest", "org.freedesktop.Notifications",
		"--object-path", "/org/freedesktop/Notifications",
		"--method", "org.freedesktop.Notifications.Notify",
		"CGP", "0", "web-browser",
		"New Crunchyroll Guest Pass",
		"<a href='" + GPRURL + code + "'>" + code + "</a>",
		"['default', 'Redeem']", "{}", "int32 5000",
	}

	c := exec.Command("gdbus", args...)
	return c.Run()
}
