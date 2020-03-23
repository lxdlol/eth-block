package pulling

import (
	"context"
	"encoding/json"
	"ethereum-block/db"
	token "ethereum-block/erc20"
	"ethereum-block/ethconnect"
	"ethereum-block/log"
	"ethereum-block/models"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"gopkg.in/mgo.v2/bson"
	"math"
	"math/big"
	"strconv"
)

var (
	Block       = make(chan *types.Block, 10000)
	Account     = make(chan models.Transfer, 100000)
	err         error
	TransferNum int
	Tps         = make(chan models.Metric, 10000)
)

const tranhash = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

//获取区块信息
func GetBlock() {
	headers := make(chan *types.Header)
	sub, err := ethconnect.Client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Log.Fatal(err)
	}
	fmt.Println("start")
	for {
		select {
		case err := <-sub.Err():
			log.Log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			block, err := ethconnect.Client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Log.Fatal(err)
			}
			rightBlock(block)
		}
	}
}

//判断落盘区块是否连续
func rightBlock(block *types.Block) {
	maxBlock := models.MaxBlock()
	max := block.Number().Int64()
	for i := int64(maxBlock); i < max; i++ {
		block, err := ethconnect.Client.BlockByNumber(context.Background(), big.NewInt(int64(i+1)))
		if err != nil {
			log.Log.Fatal(err)
		}
		fmt.Println(block.Number())
		Block <- block
	}
}

func DealBlockInfo() {
	for {
		var n int
		fmt.Println("开始解析数据")
		select {
		case block := <-Block:
			n++
			InsertBlock(block)
			InsertBlockTransfer(block)
			InsertTransfer(block)
			fmt.Println("一个区块完成", n)
		}
	}
}

//插入区块信息
func InsertBlock(block *types.Block) {
	session, collection := db.Connect(db.DB, "block")
	defer session.Close()
	var block_tmp models.Block
	number, _ := bson.ParseDecimal128(block.Number().String())
	diff, _ := bson.ParseDecimal128(block.Difficulty().String())
	totaldiff, _ := bson.ParseDecimal128("0")
	block_tmp = models.Block{
		number,
		diff,
		string(block.Extra()),
		strconv.FormatUint(block.Header().GasLimit, 32),
		strconv.FormatUint(block.GasUsed(), 32),
		block.Hash().Hex(),
		string(block.Bloom().Bytes()),
		block.Header().Coinbase.String(),
		"",
		block.MixDigest().Hex(),
		strconv.FormatUint(block.Nonce(), 32),
		block.ParentHash().Hex(),
		block.ReceiptHash().Hex(),
		"",
		block.Size().String(),
		block.Root().String(),
		strconv.FormatUint(block.Time(), 32),
		totaldiff,
		block.Transactions().Len(),
		block.Header().TxHash.Hex(),
		[]models.Block{},
	}
	fmt.Println("区块信息:", block_tmp.Number)
	fmt.Println("....")
	e := collection.Insert(&block_tmp)
	fmt.Println(e)
	if e != nil {
		log.Log.Error(e.Error())
		panic(e)
	}
}

//插入区块交易信息
func InsertBlockTransfer(block *types.Block) {
	fmt.Println("处理区块交易")
	tx := block.Transactions()
	session, collection := db.Connect(db.DB, "transaction")
	defer session.Close()
	for k, v := range tx {
		//获取交易信息
		fmt.Println("处理区块的每一笔交易", k)
		var transaction models.Transaction
		v2, r, s := v.RawSignatureValues()
		fromString, _ := decimal.NewFromString(v.Value().String())
		i := fromString.Div(decimal.NewFromInt(int64(math.Pow10(18)))).String()
		value, _ := bson.ParseDecimal128(i)
		v22, _ := bson.ParseDecimal128(v2.String())
		var to string
		if v.To() == nil {
			to = ""
		} else {
			to = v.To().String()
		}
		transaction = models.Transaction{
			BlockHash:          block.Hash().Hex(),
			BlockNumber:        block.Number().Int64(),
			From:               GetSendAddr(v),
			Gas:                strconv.FormatUint(v.Gas(), 32),
			GasPrice:           v.GasPrice().Int64(),
			Hash:               v.Hash().Hex(),
			Input:              "",
			Nonce:              strconv.FormatUint(v.Nonce(), 32),
			To:                 to,
			TransactionIndex:   int64(k + 1),
			Value:              value,
			V:                  v22,
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
	fmt.Println("处理每一笔交易")
	tx := block.Transactions()
	var num int64
	for _, v := range tx {
		n := DoReceipt(v, block)
		num += n
	}
	var node models.Metric
	node.BlockNumber = block.Number().Int64()
	node.Timestamp = int64(block.Time())
	node.TransactionCount = num
	fmt.Println("计算tps")
	Tps <- node
	fmt.Println("单笔交易end")
}

//获取发送地址
func GetSendAddr(tx *types.Transaction) string {
	chainID, err := ethconnect.Client.NetworkID(context.Background())
	if err != nil {
		log.Log.Fatal(err)
	}
	if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err != nil {
		log.Log.Fatal(err)
		return ""
	} else {
		return msg.From().String()
	}
}

//处理每条交易的log生成可索引的交易记录,
func DoReceipt(tx *types.Transaction, block *types.Block) int64 {
	fmt.Println("every log")
	receipt, err := ethconnect.Client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Log.Fatal(err)
	}
	s, collection := db.Connect(db.DB, "transfer")
	defer s.Close()
	//区块交易次数
	var num int64 = 0
	fmt.Println(len(receipt.Logs))
	if len(receipt.Logs) > 0 {
		//代币交易
		var Transfer models.Transfer
		for index, k := range receipt.Logs {
			fmt.Println(index)
			//查询token表是否存在这种代币地址
			//大于0有代笔交易，否则没有代笔交易，eth交易
			if len(k.Topics) > 0 && k.Address.Hex() != "" && k.Topics[0].Hex() == tranhash {
				num++
				var result models.Token
				var symbol string
				var dec int8
				var addr string
				session, c := db.Connect(db.DB, "token")
				defer session.Close()
				//是否已知代币
				if err := c.Find(bson.M{"contract_address": k.Address.Hex()}).One(&result); err != nil {
					//未知代币交易，用rpc调用查询到这种代币合约，存入token表
					tokenAddress := common.HexToAddress(k.Address.Hex())
					newToken, err := token.NewToken(tokenAddress, ethconnect.Client)
					if err != nil {
						log.Log.Error(err.Error())
					}
					s, err := newToken.Symbol(&bind.CallOpts{})
					if err != nil {
						log.Log.Error(k.Address.Hex() + "," + err.Error())
						goto Here
					}
					supply, err := newToken.TotalSupply(&bind.CallOpts{})
					if err != nil {
						log.Log.Error(err.Error())
					}
					decimals, err := newToken.Decimals(&bind.CallOpts{})
					if err != nil {
						log.Log.Error(err.Error())
					}
					name, err := newToken.Name(&bind.CallOpts{})
					if err != nil {
						log.Log.Error(err.Error())
					}
					decimal128, _ := bson.ParseDecimal128(supply.String())
					var token models.Token = models.Token{
						ContractAddress: k.Address.Hex(),
						Name:            name,
						Symbol:          s,
						Decimals:        int8(decimals),
						TotalSupply:     decimal128,
					}
					fmt.Println("插入新token")
					if err := models.InsertToken(token); err != nil {
						log.Log.Error(err)
						delErr(err)
					}
					symbol = s
					dec = int8(decimals)
					addr = k.Address.Hex()
				} else {
					//已知代币交易
					symbol = result.Symbol
					dec = result.Decimals
					addr = result.ContractAddress
				}
				var num string
				json.Unmarshal(k.Data, &num)
				fromString, _ := decimal.NewFromString(num)
				fmt.Println("交易额度：", fromString)
				div := fromString.Div(decimal.NewFromInt(int64(math.Pow10(int(dec)))))
				decimal128, _ := bson.ParseDecimal128(div.String())
				Transfer = models.Transfer{
					ContractAddress: addr,
					Symbol:          symbol,
					From:            k.Topics[1].Hex(),
					To:              k.Topics[2].Hex(),
					Value:           decimal128,
					TransactionHash: tx.Hash().Hex(),
					TransferIndex:   int8(index + 1),
					BlockNumber:     block.Number().Int64(),
					Timestamp:       int64(block.Time()),
					Hash:            tx.Hash().Hex() + strconv.Itoa(index+1),
				}
				Account <- Transfer
				if err := collection.Insert(&Transfer); err != nil {
					delErr(err)
				}
			}
		}
	}
Here:
	//是否有以太坊交易
	if tx.Value().Int64() > 0 {
		num++
		//太坊交易
		var Transfer models.Transfer
		decimal128, _ := bson.ParseDecimal128(strconv.Itoa(int(tx.Value().Int64() / int64(math.Pow10(18)))))
		var to string
		if tx.To() == nil {
			to = ""
		} else {
			to = tx.To().Hex()
		}
		Transfer = models.Transfer{
			ContractAddress: "BASE",
			Symbol:          "ETH",
			From:            GetSendAddr(tx),
			To:              to,
			Value:           decimal128,
			TransactionHash: tx.Hash().Hex(),
			TransferIndex:   1,
			BlockNumber:     block.Number().Int64(),
			Timestamp:       int64(block.Time()),
			Hash:            tx.Hash().Hex() + "1",
		}
		if err := collection.Insert(&Transfer); err != nil {
			log.Log.Error(err.Error())
			panic(err)
		}
		Account <- Transfer
	}
	fmt.Println("every log end")
	//传输区块交易交易数量信息
	return num
}

//错误处理
func delErr(err error) {
	log.Log.Error(err.Error())
	panic(err)
}

//账户更新
func UpdateCountBalance() {
	//轮寻account管道，获取账户数据
	var n int
	for {
		select {
		case v := <-Account:
			n++
			fmt.Println("账户越更新：", v.ContractAddress)
			//查询账户是否存在，存在更新账户余额，不存在创建
			models.AccountIsExist(v.To, v.ContractAddress, v.BlockNumber, v.Symbol)
			models.AccountIsExist(v.From, v.ContractAddress, v.BlockNumber, v.Symbol)
			fmt.Println("完成一次", n)
		}
	}
}

func UpdateTps() {
	for {
		select {
		case v := <-Tps:
			//查找上一个区块时间，计算duration
			fmt.Println("上一个区块高度", v.BlockNumber-1)
			metric, e := models.FindMetric(v.BlockNumber - 1)
			fmt.Println(e)
			if e != nil {
				metric = models.Metric{
					BlockNumber:           v.BlockNumber,
					TransactionCount:      v.TransactionCount,
					Timestamp:             v.Timestamp,
					Duration:              0,
					Tps:                   0,
					TotalTransactionCount: int64(v.TransactionCount),
				}
				e := models.InsertMetric(metric)
				if e != nil {
					delErr(e)
				}
			} else {
				duration := (v.Timestamp - metric.Timestamp)
				fmt.Println(duration)
				tps := metric.TransactionCount / duration
				fmt.Println(tps)
				total := metric.TotalTransactionCount + v.TransactionCount
				v.TotalTransactionCount = total
				models.InsertMetric(v)
				models.UpdateMetric(v.BlockNumber-1, tps, duration)
			}
			fmt.Println("tps计算完成")
		}
	}
}
