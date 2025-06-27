# M-Pesa Go
A high-performance Go module for interacting with Safaricom's M-Pesa Daraja API.

## Installation
```bash
go get github.com/freelancer254/mpesa-go
```
## Usage

```go
package main

import (
	"context"
	"github.com/freelancer254/mpesa-go/client"
	"github.com/freelancer254/mpesa-go/types"
	"github.com/freelancer254/mpesa-go/utils"
	"log"
)

func main() {
	mpesa := client.NewMpesa()
	ctx := context.Background()

	// Get access token
	token, err := mpesa.GetAccessToken(ctx, "consumer_key", "consumer_secret")
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	// Encode password
	password := utils.EncodePassword("123456", "passkey")

	// Initiate STK Push
	payload := types.STKPushRequest{
		AccessToken:       token.AccessToken,
		BusinessShortCode: "123456",
		Password:          password,
		Amount:            "100",
		PartyA:            "254700000000",
		PartyB:            "123456",
		PhoneNumber:       "254700000000",
		CallBackURL:       "https://callback.example.com",
		AccountReference:  "Test123",
		TransactionDesc:   "Payment",
	}

	response, err := mpesa.STKPush(ctx, payload)
	if err != nil {
		log.Printf("STK Push failed: %v", err) // e.g., "STK Push failed: Request canceled by user. (code: 1032)"
		return
	}
	log.Printf("STK Push response: %+v", response)
}
```
## Prerequisites
- M-Pesa API credentials (Consumer Key, Consumer Secret, ShortCode, Passkey).
- Go 1.18 or higher.