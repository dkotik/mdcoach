package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Exec runs an external command with a timeout and tracks errors.
func Exec(cmd string, args ...string) (err error) {
	// TODO: set time out in program arguments? or at least raise a proper time out error? right now just quits without doing anything
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Minute)
	defer cancel() // interrupts the process running
	p := exec.CommandContext(ctx, cmd, args...)
	// log.Printf(`Running %s %v`, cmd, args) // @TODO: manage this under verbose flag?

	var errorMessage bytes.Buffer
	p.Stdout = os.Stdout
	p.Stderr = &errorMessage
	err = p.Run()

	if err == nil && errorMessage.Len() > 0 {
		// log.Fatal(errorMessage.String())
		err = fmt.Errorf(`could not run %s command, reason: %s`, cmd, strings.TrimSpace(errorMessage.String()))
	}

	if err != nil {
		// log.Printf(`external command %s error: %s`, cmd, err.Error())
		// https://www.amazon.com/gp/feature.html?ie=UTF8&docId=1000765211 - kindlegen
		if _, err2 := exec.LookPath(`weasyprint`); err2 != nil {
			log.Printf(`Weasyprint PDF generator has not been found. Run "sudo python3 -m pip install weasyprint".`)
		}
		if _, err2 := exec.LookPath(`youtube-dl`); err2 != nil {
			// sudo pip install -U youtube-dl
			log.Printf(`Youtude-dl video downloader has not been found. Run "sudo python3 -m pip install youtube-dl".`)
		}
	}
	return err
}
