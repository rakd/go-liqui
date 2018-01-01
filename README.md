go-liqui
==========

go-liqui is an implementation of the liqui.io API in Golang.

It's just for me to handle liqui for now, so now not support all of the APIs.

Based off of https://github.com/toorop/go-bittrex/

## Import
```
import "github.com/rakd/go-liqui"
```

## Usage
~~~ go
package main

import (
	"fmt"
	"github.com/rakd/go-liqui"
)

const (
	API_KEY    = "YOUR_API_KEY"
	API_SECRET = "YOUR_API_SECRET"
)

func main() {
	// liqui client
	client := liqui.New(API_KEY, API_SECRET)

	// Get balance
	balance, resp, err := client.GetBalance()
	log.Printf("err:%v", err)
	log.Printf("resp:%s", string(resp))
	log.Printf("balance:%v", balance)
}
~~~


## Donate

- BTC: 1Ah8sarQ4w9FnsCs8LoG6JuYiFHmrAAy6F
