package ethconnect

import (
	"ethereum-block/log"
	"github.com/ethereum/go-ethereum/ethclient"
)

var Client *ethclient.Client
var err error

func init() {
	Client, err = ethclient.Dial("ws://192.168.8.126:8561")
	if err != nil {
		log.Log.Fatal(err)
	}
}
