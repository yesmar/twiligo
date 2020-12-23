# Twiligo sample

The Twiligo sample illustrates how to use the Twiligo package. The sample requires three environment variables, from which it loads part of its runtime configuration:

```zsh
export TWILIO_ACCOUNT_SID=your_sid_here
export TWILIO_AUTH_TOKEN=your_auth_token_here
export TWILIO_PHONE_NUMBER=your_twilio_phone_number_here
```

The only other configuration required is the message and the phone number to send it to. These two items are entered via the command line:

```zsh
cd cmd/twiligo
go build
./twiligo -to +12155551212 -msg hello
```

Assuming all goes well, the program will emit `ok`, otherwise it will print out the error or HTTP status of the Twilio API call.

By default, the sample driver will wait 5 seconds for the Twilio API call to complete. The optional `-timeout` flag can be specified with a different timeout value, e.g., `-timeout 10`. The timeout value is always interpreted as seconds and is `unsigned`.
