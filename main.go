package main

import "ethereum-block/pulling"

func main() {
	go pulling.DealBlockInfo()
	go pulling.UpdateCountBalance()
	go pulling.UpdateTps()
	pulling.GetBlock()
}
