package pulling

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

var (
	Block  = make(chan *types.Block, 100)
	Client *ethclient.Client
	err    error
)

func GetBlock() {
	Client, err = ethclient.Dial("ws://192.168.8.126:8561")
	if err != nil {
		log.Log.Fatal(err)
	}

	headers := make(chan *types.Header)
	sub, err := Client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			block, err := Client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Log.Fatal(err)
			}
			Block <- block
		}
	}
}
func DealBlockInfo(ch chan *types.Block) {
	for block := range ch {
		//InsertBlock(block)
		InsertBlockTransfer(block)
	}
}

//插入区块信息
func InsertBlock(block *types.Block) {
	session, collection := db.Connect(DB, "transaction")
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
	session, collection := db.Connect(DB, "block")
	defer session.Close()
	for k, v := range tx {
		//获取交易信息
		var transaction models.Transaction
		v2, r, s := v.RawSignatureValues()
		transaction = models.Transaction{
			BlockHash:          block.Hash().String(),
			BlockNumber:        block.Number().Int64(),
			From:               GetSendAddr(v),
			Gas:                v.Gas(),
			GasPrice:           v.GasPrice().Int64(),
			Hash:               v.Hash().String(),
			Input:              string(v.Data()),
			Nonce:              v.Nonce(),
			To:                 v.To().String(),
			TransactionIndex:   int64(k + 1),
			Value:              v.Value(),
			V:                  v2,
			R:                  r.String(),
			S:                  s.String(),
			TransactionReceipt: models.TransactionReceipt{},
		}
		if e := collection.Insert(&transaction); e != nil {
			log.Log.Error(e.Error())
			panic(e)
		}
	}
}

//解析区块交易信息，分解
func InsertTransfer(block *types.Block) {

}

//计算节点tps

//获取发送地址
func GetSendAddr(tx *types.Transaction) string {
	chainID, err := Client.NetworkID(context.Background())
	if err != nil {
		log.Log.Fatal(err)
	}
	if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err != nil {
		return msg.From().String()
	} else {
		return ""
	}
}

func GetReceipt(tx *types.Transaction) {
	receipt, err := Client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Log.Fatal(err)
	}
	for _, k := range receipt.Logs {
		k.Topics
	}
}
