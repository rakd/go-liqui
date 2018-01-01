package liqui

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBalance(t *testing.T) {
	key := os.Getenv("LIQUI_KEY")
	secret := os.Getenv("LIQUI_SECRET")

	api := New(key, secret)
	ret, body, err := api.GetBalance()

	log.Printf("err:%v", err)
	log.Printf("body:%s", string(body))
	log.Printf("ret:%v", ret)
	require.NoError(t, err, nil)

	return
}

func TestGetTicker(t *testing.T) {
	key := os.Getenv("LIQUI_KEY")
	secret := os.Getenv("LIQUI_SECRET")

	api := New(key, secret)
	ret, body, err := api.GetTicker()
	log.Printf("err:%v", err)
	log.Printf("body:%s", string(body))
	log.Printf("ret:%v", ret)
	require.NoError(t, err, nil)

	return
}
