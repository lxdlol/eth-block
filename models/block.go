package models

import (
	"gopkg.in/mgo.v2/bson"
	"math/big"
)

type Block struct {
	/**
	AttributeDict(
	{'difficulty': 572954,
	'extraData': HexBytes('0xd883010908846765746888676f312e31332e34856c696e7578'),
	'gasLimit': 8000000,
	'gasUsed': 0,
	'hash': HexBytes('0x10f5102c4d5907fa387b042fddb47230854604786603cd30db39a43d78c79402'),
	'logsBloom': HexBytes('0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000'),
	'miner': '0x0d8c6aBa421723b3bCE849C70C06592f696E4399',
	'mixHash': HexBytes('0x5cd73e66b5483277d3231b74ed42a145b8790d6f79928f3ed8acea2d7ed11478'),
	'nonce': HexBytes('0x2a1415c2f834debd'),
	'number': 11111,
	'parentHash': HexBytes('0x4db4693dec24eafc2b018ec1a6752fb1ca979593ec7a8afab799a3fba77cbd0c'),
	'receiptsRoot': HexBytes('0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421'),
	'sha3Uncles': HexBytes('0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347'),
	'size': 538,
	'stateRoot': HexBytes('0x3be1b45f7a6d62226dd02244f7b76b54433a580787b5329c3cf726d4d9305505'),
	'timestamp': 1583421959,
	'totalDifficulty': 4909026460,
	'transactions': [],
	'transactionsRoot': HexBytes('0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421'),
	'uncles': []})

	*/
	Number            int64   `json:"number" bson:"number"` //unique
	Difficulty        int64   `json:"difficulty"`
	ExtraData         string  `json:"extra_data"` //附加数据
	GasLimit          int64   `json:"gas_limit"`
	GasUsed           int64   `json:"gas_used"`
	Hash              string  `json:"hash"`
	LogsBloom         string  `json:"logs_bloom"`
	Miner             string  `json:"miner"`
	MinerName         string  `json:"miner_name"` //播报方
	MixHash           string  `json:"mix_hash"`
	Nonce             string  `json:"nonce"`
	ParentHash        string  `json:"parent_hash"`
	ReceiptsRoot      string  `json:"receipts_root"`
	Sha3Uncles        string  `json:"sha_3_uncles"`
	Size              int64   `json:"size"`
	StateRoot         string  `json:"state_root"`
	Timestamp         string  `json:"timestamp"`
	TotalDifficulty   int64   `json:"total_difficulty"`
	TransactionsCount int64   `json:"transactions_counts"` //交易
	TransactionsRoot  string  `json:"transactions_root"`
	Uncles            []Block `json:"uncles"`
}

type Transaction struct {
	/**

	交易

	AttributeDict({
	'blockHash': HexBytes('0x34b38b733ef9661068e70f5e26a9315731e351206c6931881002bf7fb6da3d2d'),
	'blockNumber': 16407,
	'from': '0x0d8c6aBa421723b3bCE849C70C06592f696E4399',
	'gas': 178261,
	'gasPrice': 1000000000,
	'hash': HexBytes('0xad2297deee3b8809fde6bdaff6e63cbb912d863a3ef05165c9ff9149ccff9b02'),
	'input': '0x608060405234801561001057600080fd5b50610248806100206000396000f3fe608060405260043610610057576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630dbe671f1461005c578063289f8d3a146100875780634df7e3d0146100c4575b600080fd5b34801561006857600080fd5b506100716100ef565b60405161007e9190610191565b60405180910390f35b34801561009357600080fd5b506100ae60048036036100a9919081019061011c565b6100f5565b6040516100bb9190610176565b60405180910390f35b3480156100d057600080fd5b506100d9610102565b6040516100e69190610191565b60405180910390f35b60005481565b6000818316905092915050565b60015481565b600061011482356101e2565b905092915050565b6000806040838503121561012f57600080fd5b600061013d85828601610108565b925050602061014e85828601610108565b9150509250929050565b610161816101ac565b82525050565b610170816101d8565b82525050565b600060208201905061018b6000830184610158565b92915050565b60006020820190506101a66000830184610167565b92915050565b60007fff0000000000000000000000000000000000000000000000000000000000000082169050919050565b6000819050919050565b60007fff000000000000000000000000000000000000000000000000000000000000008216905091905056fea265627a7a72305820bf34a3233172b51540dfe4e4c3456faf448a8557c6156ff9b0ec1ce42981bec96c6578706572696d656e74616cf50037',
	'nonce': 50,
	'to': None,
	'transactionIndex': 0,
	'value': 0,
	'v': 74078,
	'r': HexBytes('0xd95987fb9b61b5cf5110e2f9924d0a47867b7deb9d2e644fe478afba03f14ce6'),
	's': HexBytes('0x600a9fe0c68724ece3b2f4a5238470716856cd7deb76501cec3b514e7eb4b144')
	})
	*/

	BlockHash          string             `json:"block_hash" bson:"block_hash"`
	BlockNumber        int64              `json:"block_number" bson:"block_number"`
	From               string             `json:"from" bson:"from"`
	Gas                uint64             `json:"gas" bson:"gas"`
	GasPrice           int64              `json:"gas_price" bson:"gas_price"`
	Hash               string             `json:"hash"`
	Input              string             `json:"input"` //输入数据
	Nonce              uint64             `json:"nonce"`
	To                 string             `json:"to"`
	TransactionIndex   int64              `json:"transaction_index"` //position
	Value              *big.Int           `json:"value"`
	V                  *big.Int           `json:"v"`
	R                  string             `json:"r"`
	S                  string             `json:"s"`
	TransactionReceipt TransactionReceipt `json:"transaction_receipt" bson:"transaction_receipt"` //交易结果
}

type TransactionReceipt struct {
	/**
	AttributeDict({
	'blockHash': HexBytes('0x34b38b733ef9661068e70f5e26a9315731e351206c6931881002bf7fb6da3d2d'),
	'blockNumber': 16407,
	'contractAddress': '0xDc4AB61c43DAAe1f25240b2a09C8D442415AaBD5',
	'cumulativeGasUsed': 178261,
	'from': '0x0d8c6aba421723b3bce849c70c06592f696e4399',
	'gasUsed': 178261,
	'logs': [],
	'logsBloom': HexBytes('0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000'),
	'status': 1,
	'to': None,
	'transactionHash': HexBytes('0xad2297deee3b8809fde6bdaff6e63cbb912d863a3ef05165c9ff9149ccff9b02'),
	'transactionIndex': 0}
	)
	*/

	ContractAddress   string `json:"contract_address" bson:"contract_address"`
	CumulativeGasUsed string `json:"cumulative_gas_used" bson:"cumulative_gas_used"`
	GasUsed           int64  `json:"gas_used" bson:"gas_used"`
	Logs              []Log  `json:"logs" bson:"logs"` //事件日记
	LogsBloom         string `json:"logs_bloom" bson:"logs_bloom"`
	Status            int8   `json:"status" bson:"status"` //交易结果 1:成功,0:失败。

}

type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address string `json:"address" bson:"address"`
	// list of topics provided by the contract.
	Topics []string `json:"topics" bson:"topics"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data" bson:"data"`
	// index of the log in the block
	Index uint `json:"log_index" bson:"log_index"`
	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed" bson:"removed"`
}

type Token struct {
	ContractAddress string   `json:"contract_address" bson:"contract_address"`
	Name            string   `json:"name" bson:"name"`
	Symbol          string   `json:"symbol" bson:"symbol"`
	Decimals        int8     `json:"decimals" bson:"decimals"`
	TotalSupply     *big.Int `json:"total_supply" bson:"total_supply"`
}

type Account struct {
	/**
	账户
	index 1:address
	index 2:ContractAddress Balance
	index 3:address ContractAddress
	*/
	ContractAddress string          `json:"contract_address"` //contract address,如果是系统币, contract_address=BASE
	Symbol          string          `json:"symbol"`           //代币的符号
	Address         string          `json:"address"`          //账户地址
	Balance         bson.Decimal128 `json:"balance"`          //账户余额
	BlockNumber     int64           `json:"block_number"`     //更新余额时的区块链高度
}

type Transfer struct {
	/**
	转账

	index 1: address,timestamp
	index 2: timestamp
	index 3: contract_address timestamp
	index 4: From To
	index 5: TransactionHash
	unique:Hash
	*/
	ContractAddress string          `json:"contract_address"` //合约地址, 不空
	Symbol          string          `json:"symbol"`           //代币
	From            string          `json:"from"`             //发送方
	To              string          `json:"to"`               //接收方
	Value           bson.Decimal128 `json:"value"`            //数量
	TransactionHash string          `json:"transaction_hash"` //交易哈希
	TransferIndex   int8            `json:"transfer_index"`   //输入顺序
	BlockNumber     int64           `json:"block_number"`     //区块高度
	Timestamp       int64           `json:"timestamp"`        //时间戳
	Hash            string          `json:"hash"`             //unique:transactionHash+TransferIndex
}
