// SPDX-License-Identifier: MIT

package twiligo

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiVersion = "2010-04-01"
	maxBackoff = 4 // Seconds.
)

// Config stores the Twiligo configuration.
type Config struct {
	AccountSid string        // Account SID.
	authToken  string        // Authorization token.
	From       string        // Twilio phone number.
	MaxMsgLen  int           // Maximum SMS message length.
	Timeout    time.Duration // HTTP timeout in seconds for communicating with the Twilio API.
}

// New returns a Twiligo instance initialized using the specified parameters. It returns an error if initialization failed.
func New(sid, token, phoneNumber string, timeout uint, enableMessageConcatenation bool) (*Config, error) {
	if sid == "" {
		return nil, errors.New("no account sid specified")
	}
	if token == "" {
		return nil, errors.New("no authentication token specified")
	}
	if phoneNumber == "" {
		return nil, errors.New("no Twilio phone number specified")
	}

	// https://support.twilio.com/hc/en-us/articles/223181508-Does-Twilio-support-concatenated-SMS-messages-or-messages-over-160-characters-
	maxLen := 160
	if enableMessageConcatenation {
		maxLen = 1600
	}

	return &Config{
		AccountSid: sid,
		authToken:  token,
		From:       phoneNumber,
		MaxMsgLen:  maxLen,
		Timeout:    time.Second * time.Duration(timeout),
	}, nil
}

// String returns the Account SID of the Twiligo instance.
func (tc *Config) String() string {
	return tc.AccountSid
}

// SendMessage sends a message to the specified mobile phone number. It returns nil on success and an error otherwise.
func (tc *Config) SendMessage(to, message string) error {
	if to == "" {
		return errors.New("no phone number specified")
	}
	if message == "" {
		return errors.New("no message specified")
	}
	if len(message) > tc.MaxMsgLen {
		return fmt.Errorf("message exceeeds %d bytes", tc.MaxMsgLen)
	}

	ep := fmt.Sprintf("https://api.twilio.com/%s/Accounts/%s/Messages.json", apiVersion, tc.AccountSid)

	v := url.Values{}
	v.Set("To", to)
	v.Set("From", tc.From)
	v.Set("Body", message)
	vReader := *strings.NewReader(v.Encode())

	client := &http.Client{
		Timeout: tc.Timeout,
	}
	req, err := http.NewRequest("POST", ep, &vReader)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(tc.AccountSid, tc.authToken)

	var resp *http.Response

	for backOff := time.Second * 0; backOff <= maxBackoff; backOff *= 2 {
		time.Sleep(backOff)
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 429 { // Too many requests.
			break
		}
	}

	// Although the response is documented here https://www.twilio.com/docs/usage/twilios-response
	// we don't actually care about any of the fields. We return nil in the case of a successful
	// response, and we return the HTTP status message (e.g., "400 BAD REQUEST") otherwise.
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil // Success.
	}
	return errors.New(resp.Status) // Failure.
}
