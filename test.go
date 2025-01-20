package main

import (
	"fmt"
	"math/big"
	"zkwasm-minirollup-rpc-go/zkwasm"
)

var (
	InitPlayerCmd  = big.NewInt(1)
	BuyElfCmd      = big.NewInt(2)
	CollectCoinCmd = big.NewInt(11)
	CleanRanchCmd  = big.NewInt(4)
	DepositCmd     = big.NewInt(8)
)

func main() {
	prikey := "1234"
	pid1, pid2 := zkwasm.GetPid(prikey)

	fmt.Println("pid1:", pid1.Int64())
	fmt.Println("pid2:", pid2.Int64())

	data := zkwasm.Query(prikey)
	fmt.Println("data:", data)

	zkwamRpc := zkwasm.NewZKWasmAppRpc("http://localhost:3000")

	//zkwamRpc := zkwasm.NewZKWasmAppRpc("https://zk-server.pumpelf.ai")
	state, _ := zkwamRpc.QueryState(prikey)
	fmt.Println("state:", state)
	// 收集金币
	//nonce, _ := zkwamRpc.GetNonce(prikey)
	//cmd := zkwamRpc.CreateCommand(nonce, big.NewInt(11), big.NewInt(0))
	//fmt.Println("cmd:", cmd)
	//transaction, _ := zkwamRpc.SendTransaction([4]*big.Int{cmd, big.NewInt(1), big.NewInt(1), big.NewInt(0)}, prikey)
	//fmt.Println("transaction:", transaction)
	// 初始化玩家
	initPlayer(zkwamRpc, prikey)
	// 购买宠物
	buyElf(zkwamRpc, prikey)
	// 收集金币
	//collectCoin(zkwamRpc, prikey)

	// 清理牧场
	//clearRanch(zkwamRpc, prikey)

	// 充值
	deposit(zkwamRpc, prikey)

	// 查询状态
	state, _ = zkwamRpc.QueryState(prikey)
	fmt.Println("state:", state)
}

// 初始化玩家
func initPlayer(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	// 获取配置
	cmd := zkwamRpc.CreateCommand(big.NewInt(0), InitPlayerCmd, big.NewInt(0))
	transaction, _ := zkwamRpc.SendTransaction([4]*big.Int{cmd, big.NewInt(0), big.NewInt(0), big.NewInt(0)}, prikey)
	fmt.Println("transaction:", transaction)
}

// 购买宠物
func buyElf(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranchId := big.NewInt(1)
	elfType := big.NewInt(1)
	cmd := zkwamRpc.CreateCommand(nonce, BuyElfCmd, big.NewInt(0))
	transaction, _ := zkwamRpc.SendTransaction([4]*big.Int{cmd, ranchId, elfType, big.NewInt(0)}, prikey)
	fmt.Println("transaction:", transaction)
}

// 收集金币
func collectCoin(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranchId := big.NewInt(1)
	elfId := big.NewInt(1)
	cmd := zkwamRpc.CreateCommand(nonce, CollectCoinCmd, big.NewInt(0))
	transaction, _ := zkwamRpc.SendTransaction([4]*big.Int{cmd, ranchId, elfId, big.NewInt(0)}, prikey)
	fmt.Println("transaction:", transaction)
}

// 清理牧场
func clearRanch(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranchId := big.NewInt(1)
	cmd := zkwamRpc.CreateCommand(nonce, CleanRanchCmd, big.NewInt(0))
	transaction, _ := zkwamRpc.SendTransaction([4]*big.Int{cmd, ranchId, big.NewInt(0), big.NewInt(0)}, prikey)
	fmt.Println("transaction:", transaction)
}

func deposit(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	pid1, pid2 := zkwasm.GetPid(prikey)
	ranchId := big.NewInt(590)
	propType := big.NewInt(89)
	depositP := new(big.Int).Lsh(ranchId, 32)
	depositP.Add(depositP, propType)
	cmd := zkwamRpc.CreateCommand(big.NewInt(0), DepositCmd, big.NewInt(0))
	transaction, _ := zkwamRpc.SendTransaction([4]*big.Int{cmd, pid1, pid2, depositP}, prikey)
	fmt.Println("transaction:", transaction)
}
