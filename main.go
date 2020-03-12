package main

import "ethereum-block/pulling"

func main() {
	go pulling.GetBlock()
	pulling.DealBlockInfo(pulling.Block)
}
