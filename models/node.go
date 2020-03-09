package models

type Metric struct {
	/**
	区块链性能度量数据
	*/
	BlockNumber           int8  `json:"block_number"`      //unique,index
	TransactionCount      int8  `json:"transaction_count"` //该区块链交易数
	Timestamp             int64 `json:"timestamp"`
	Duration              int8  `json:"duration"`                //区块处理时间，秒。
	Tps                   int8  `json:"tps"`                     //index,TransactionCount/Duration
	TotalTransactionCount int64 `json:"total_transaction_count"` //交易总数,TotalTransactionCount=this.TransactionCount+pre.TotalTransactionCount
}

type Node struct {
	NodeType string `json:"node_type"` //授权节点,同步节点
	Coinbase string `json:"coinbase"`
	Name     string `json:"name"`
	Ip       string `json:"ip"`
	Geo      string `json:"geo"`
}

type CandidateNode struct {
	/*
		候选节点
	*/
	Coinbase    string   `json:"coinbase"`
	Name        string   `json:"name"`
	Ip          string   `json:"ip"`
	Geo         string   `json:"geo"`
	BlockNumber int8     `json:"block_number"`
	Timestamp   int64    `json:"timestamp"`
	Votes       []string `json:"votes"` //授权节点投票，得票超过50%，候选节点变为授权节点。
}
