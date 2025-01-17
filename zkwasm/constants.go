package zkwasm

import (
	"fmt"
	"math/big"
)

var Constants = map[string]*Field{
	"c":  NewField(toBigInt("0")),
	"a":  NewField(toBigInt("21888242871839275222246405745257275088548364400416034343698204186575808495616")),
	"d":  NewField(toBigInt("12181644023421730124874158521699555681764249180949974110617291017600649128846")),
	"gX": NewField(toBigInt("21237458262955047976410108958495203094252581401952870797780751629344472264183")),
	"gY": NewField(toBigInt("2544379904535866821506503524998632645451772693132171985463128613946158519479")),
}

func toBigInt(numStr string) *big.Int {
	// 创建一个新的 big.Int
	num := new(big.Int)

	// 将字符串转换为 big.Int
	num, success := num.SetString(numStr, 10)
	if !success {
		fmt.Println("无法转换字符串为大整数")
		return nil
	}
	return num
}
