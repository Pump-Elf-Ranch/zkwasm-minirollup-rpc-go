package zkwasm

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func bytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func bytesToDecimal(bytes []byte) string {
	var sb strings.Builder
	for _, b := range bytes {
		sb.WriteString(fmt.Sprintf("%02d", b))
	}
	return sb.String()
}

func composeWithdrawParams(address string, amount *big.Int) ([]*big.Int, error) {
	addressBytes, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}
	firstLimb := new(big.Int).SetBytes(reverseBytes(addressBytes[:4]))
	sndLimb := new(big.Int).SetBytes(reverseBytes(addressBytes[4:12]))
	thirdLimb := new(big.Int).SetBytes(reverseBytes(addressBytes[12:20]))
	one := new(big.Int).Add(new(big.Int).Lsh(firstLimb, 32), amount)
	return []*big.Int{one, sndLimb, thirdLimb}, nil
}

func decodeWithdraw(txdata []byte) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	if len(txdata) > 1 {
		for i := 0; i < len(txdata); i += 32 {
			extra := txdata[i : i+4]
			address := txdata[i+4 : i+24]
			amount := txdata[i+24 : i+32]
			amountInWei, err := strconv.ParseInt(bytesToDecimal(amount), 10, 64)
			if err != nil {
				return nil, err
			}
			result = append(result, map[string]interface{}{
				"op":      extra[0],
				"index":   extra[1],
				"address": "0x" + bytesToHex(address),
				"amount":  amountInWei,
			})
		}
	}
	return result, nil
}

func reverseBytes(bytes []byte) []byte {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return bytes
}

type PlayerConvention struct {
	processingKey   string
	rpc             *ZKWasmAppRpc
	commandDeposit  *big.Int
	commandWithdraw *big.Int
}

func NewPlayerConvention(key string, rpc *ZKWasmAppRpc, commandDeposit, commandWithdraw *big.Int) *PlayerConvention {
	return &PlayerConvention{
		processingKey:   key,
		rpc:             rpc,
		commandDeposit:  commandDeposit,
		commandWithdraw: commandWithdraw,
	}
}

func (pc *PlayerConvention) createCommand(nonce, command, objindex *big.Int) *big.Int {
	bigNonce0 := new(big.Int).Lsh(nonce, 16) // cmd[1] << 16
	bigCmd0 := new(big.Int).Lsh(command, 8)  // cmd[2] << 8
	cmd := new(big.Int).Add(bigNonce0, bigCmd0)
	cmd = cmd.Add(cmd, objindex)
	return cmd
}

func (pc *PlayerConvention) getConfig() (map[string]interface{}, error) {
	return pc.rpc.QueryConfig()
}

func (pc *PlayerConvention) getState() (map[string]interface{}, error) {
	state, err := pc.rpc.QueryState(pc.processingKey)
	if err != nil {
		return nil, err
	}
	var parsedState map[string]interface{}
	if err := json.Unmarshal([]byte(state["data"].(string)), &parsedState); err != nil {
		return nil, err
	}
	return parsedState, nil
}

func (pc *PlayerConvention) getNonce() (*big.Int, error) {
	data, err := pc.getState()
	if err != nil {
		return nil, err
	}
	player := data["player"].(map[string]interface{})
	nonce := big.NewInt(int64(player["nonce"].(float64)))
	return nonce, nil
}

func (pc *PlayerConvention) Deposit(pid1, pid2, amount *big.Int) (string, error) {
	nonce, err := pc.getNonce()
	if err != nil {
		return "", err
	}
	return pc.rpc.SendTransaction([4]*big.Int{
		pc.createCommand(nonce, pc.commandDeposit, big.NewInt(0)),
		pid1,
		pid2,
		amount,
	}, pc.processingKey)
}

func (pc *PlayerConvention) WithdrawRewards(address string, amount *big.Int) (string, error) {
	nonce, err := pc.getNonce()
	if err != nil {
		return "", err
	}
	params, err := composeWithdrawParams(address, amount)
	if err != nil {
		return "", err
	}
	return pc.rpc.SendTransaction([4]*big.Int{
		pc.createCommand(nonce, pc.commandWithdraw, big.NewInt(0)),
		params[0],
		params[1],
		params[2],
	}, pc.processingKey)
}
