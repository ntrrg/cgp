// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package main

import (
	"os/exec"
)

// gdbus monitor --session --dest org.freedesktop.Notifications --object-path /org/freedesktop/Notifications

func redeemCode(code string) error {
	args := []string{
		"call", "--session", "--dest", "org.freedesktop.Notifications",
		"--object-path", "/org/freedesktop/Notifications",
		"--method", "org.freedesktop.Notifications.Notify",
		"CGP", "0", "web-browser",
		"New Crunchyroll Guest Pass", "<b>"+code+"</b>",
		"['default', 'Redeem']", "{}", "int32 5000",
	}

	c := exec.Command("gdbus", args...)
	return c.Run()
}
