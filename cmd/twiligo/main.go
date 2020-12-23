// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yesmar/twiligo"
)

// loadConfiguration initializes Twiligo from specific environment variables.
func loadConfiguration(timeout uint) (*twiligo.Config, error) {
	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	if sid == "" {
		return nil, errors.New("please export TWILIO_ACCOUNT_SID")
	}

	token := os.Getenv("TWILIO_AUTH_TOKEN")
	if token == "" {
		return nil, errors.New("please export TWILIO_AUTH_TOKEN")
	}

	from := os.Getenv("TWILIO_PHONE_NUMBER")
	if from == "" {
		return nil, errors.New("please export TWILIO_PHONE_NUMBER")
	}

	return twiligo.New(sid, token, from, timeout, true)
}

func main() {
	to := flag.String("to", "", "+phone number")
	msg := flag.String("msg", "", "message")
	timeout := flag.Uint("timeout", 5, "timeout")

	flag.Parse()

	twilio, err := loadConfiguration(*timeout)
	if err != nil {
		log.Fatal(err)
	}

	if err = twilio.SendMessage(*to, *msg); err != nil {
		log.Fatal(err)
	}
	fmt.Println("ok")
}
