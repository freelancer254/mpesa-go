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
    "github.com/freelancer254/mpesa-go/client"
    "log"
)

func main() {
    mpesa := client.NewMpesa()
    token, err := mpesa.GetAccessToken("consumerKey", "consumerSecret")
    if err != nil {
        log.Fatal(err)
    }
    log.Println(token)
}
```
## Prerequisites
- M-Pesa API credentials (Consumer Key, Consumer Secret, ShortCode, Passkey).
- Go 1.18 or higher.