package main

import (
	"context"
	"ethereum-block/db"
	"ethereum-block/log"
	"ethereum-block/models"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strconv"
)

const (
	DB = "ethereum-block"
)

var Block = make(chan *types.Block, 100)

func GetBlock() {
	client, err := ethclient.Dial("ws://192.168.8.126:8561")
	if err != nil {
		log.Log.Fatal(err)
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Log.Fatal(err)
			}
			Block <- block
		}
	}
}
func DealBlockInfo(ch chan *types.Block) {
	for block := range ch {
		InsertBlock(block)

	}

}

//插入区块信息
func InsertBlock(block *types.Block) {
	session, collection := db.Connect(DB, "block")
	defer session.Close()
	block_tmp := models.Block{
		block.Header().Number,
		block.Header().Difficulty,
		string(block.Extra()),
		block.Header().GasLimit,
		block.GasUsed(),
		block.Header().Hash().String(),
		string(block.Bloom().Bytes()),
		block.Header().Coinbase.String(),
		"",
		block.MixDigest().String(),
		block.Nonce(),
		block.ParentHash().String(),
		block.ReceiptHash().String(),
		"",
		block.Size().String(),
		block.Root().String(),
		strconv.FormatUint(block.Time(), 64),
		big.NewInt(0),
		block.Transactions().Len(),
		block.Header().TxHash.String(),
		[]models.Block{},
	}
	e := collection.Insert(&block_tmp)
	log.Log.Error(e.Error())
	panic(e)

}

//插入区块交易信息
func InsertBlockTransfer(block *types.Block) {
	tx := block.Transactions()

}

//解析区块交易信息，分解
func InsertTransfer(block *types.Block) {

}

//计算节点tps
