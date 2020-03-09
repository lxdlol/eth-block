#ethereum-block


pulling
从区块链节点里拉取数据,保存到数据库。

1.监听到最新区块
2.读取最新区块，保存到block
3.读取最新区块里的交易数据,保存到transaction
3.1 解析transaction里的交易数据，如果是代币，新增或更新代币相关信息，并保存交易记录。
4.计算node.metric