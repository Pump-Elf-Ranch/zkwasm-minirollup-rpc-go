package main

import (
	"fmt"
	"zkwasm-minirollup-rpc-go/zkwasm"
)

func main() {
	prikey := "1234"
	pid1, pid2 := zkwasm.GetPid(prikey)

	fmt.Println("pid1:", pid1)
	fmt.Println("pid2:", pid2)

	data := zkwasm.Query(prikey)
	fmt.Println("data:", data)

	//zkwamRpc := zkwasm.NewZKWasmAppRpc("https://zk-server.pumpelf.ai")

	zkwamRpc := zkwasm.NewZKWasmAppRpc("https://zk-server.pumpelf.ai")
	state, _ := zkwamRpc.QueryState(prikey)
	fmt.Println("state:", state)
	// 收集金币
	nonce, _ := zkwamRpc.GetNonce(prikey)
	cmd := zkwamRpc.CreateCommand(nonce, 11, 0)
	fmt.Println("cmd:", cmd)
	zkwamRpc.SendTransaction([4]uint64{cmd, 1, 1, 0}, prikey)

}
