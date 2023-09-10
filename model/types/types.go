package types

import (
	"encoding/hex"
	"errors"
	"github.com/bitxx/ethutil/util"
	"github.com/ethereum/go-ethereum"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type UrlParam struct {
	RpcUrl string
	WsUrl  string
}

type CallMethodOpts struct {
	Nonce                int64
	Value                string
	GasPrice             string // MaxFeePerGas
	GasLimit             string
	IsPredictError       bool
	MaxPriorityFeePerGas string
}

type CallMethodOptsBigInt struct {
	Nonce                uint64
	Value                *big.Int
	GasPrice             *big.Int // MaxFeePerGas
	GasLimit             uint64
	IsPredictError       bool
	MaxPriorityFeePerGas *big.Int
}

type BuildTxResult struct {
	SignedTx *types.Transaction
	TxHex    string
}

type TransactionByHashResult struct {
	SignedTx    *types.Transaction
	From        common.Address
	IsPending   bool   // 交易是否处于Pending状态
	Status      string // 0: 交易失败, 1: 交易成功
	GasUsed     string // 实际花费gas
	BlockNumber string // 区块高度
}

type Erc20TxParams struct {
	ToAddress string `json:"toAddress"`
	Amount    string `json:"amount"`
	Method    string `json:"method"`
}

// CallMsg contains parameters for contract calls.
type CallMsg struct {
	Msg ethereum.CallMsg
}

// NewCallMsg creates an empty contract call parameter list.
func NewCallMsg() *CallMsg {
	return new(CallMsg)
}

func (msg *CallMsg) GetFrom() string     { return msg.Msg.From.String() }
func (msg *CallMsg) GetGasLimit() string { return strconv.FormatUint(msg.Msg.Gas, 10) }
func (msg *CallMsg) GetGasPrice() string { return msg.Msg.GasPrice.String() }
func (msg *CallMsg) GetValue() string    { return msg.Msg.Value.String() }
func (msg *CallMsg) GetData() []byte     { return msg.Msg.Data }
func (msg *CallMsg) GetDataHex() string  { return util.HexEncodeToString(msg.Msg.Data) }
func (msg *CallMsg) GetTo() string       { return msg.Msg.To.String() }

func (msg *CallMsg) SetFrom(address string) { msg.Msg.From = common.HexToAddress(address) }
func (msg *CallMsg) SetGasLimit(gas string) {
	i, _ := strconv.ParseUint(gas, 10, 64)
	msg.Msg.Gas = i
}
func (msg *CallMsg) SetGasPrice(price string) {
	i, _ := new(big.Int).SetString(price, 10)
	msg.Msg.GasPrice = i
}

// Set amount with decimal number
func (msg *CallMsg) SetValue(value string) {
	i, _ := new(big.Int).SetString(value, 10)
	msg.Msg.Value = i
}

// Set amount with hexadecimal number
func (msg *CallMsg) SetValueHex(hex string) {
	hex = strings.TrimPrefix(hex, "0x") // must trim 0x !!
	i, _ := new(big.Int).SetString(hex, 16)
	msg.Msg.Value = i
}
func (msg *CallMsg) SetData(data []byte) { msg.Msg.Data = common.CopyBytes(data) }
func (msg *CallMsg) SetDataHex(hex string) {
	data, err := util.HexDecodeString(hex)
	if err != nil {
		return
	}
	msg.Msg.Data = data
}
func (msg *CallMsg) SetTo(address string) {
	if address == "" {
		msg.Msg.To = nil
	} else {
		a := common.HexToAddress(address)
		msg.Msg.To = &a
	}
}

func (msg *CallMsg) TransferToTransaction() *Transaction {
	return &Transaction{
		GasPrice: msg.GetGasPrice(),
		GasLimit: msg.GetGasLimit(),
		To:       msg.GetTo(),
		Value:    msg.GetValue(),
		//Data:     msg.GetDataHex(),
		Data: msg.GetData(),
	}
}

type Transaction struct {
	Nonce    uint64 // nonce of sender account
	GasPrice string // wei per gas
	GasLimit string // gas limit
	To       string // receiver
	Value    string // wei amount
	Data     []byte // contract invocation input data

	// EIP1559, Default is ""
	MaxPriorityFeePerGas string
}

func NewTransaction(nonce uint64, gasPrice, gasLimit, to, value string, data []byte) *Transaction {
	return &Transaction{nonce, gasPrice, gasLimit, to, value, data, ""}
}

func NewTransactionFromHex(hexData string) (*Transaction, error) {
	rawBytes, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, err
	}
	decodeTx := types.NewTx(&types.DynamicFeeTx{})
	err = decodeTx.UnmarshalBinary(rawBytes)
	if err != nil {
		return nil, err
	}
	tx := NewTransaction(
		decodeTx.Nonce(),
		decodeTx.GasFeeCap().String(),
		strconv.Itoa(int(decodeTx.Gas())),
		decodeTx.To().String(),
		decodeTx.Value().String(),
		decodeTx.Data())
	//hex.EncodeToString(decodeTx.Data()))
	// not equal, is eip1559; legacy feecap equal tipcap
	if decodeTx.GasTipCap().Cmp(decodeTx.GasFeeCap()) != 0 {
		tx.MaxPriorityFeePerGas = decodeTx.GasTipCap().String()
	}
	return tx, nil
}

// This is an alias property for GasPrice in order to support EIP1559
func (tx *Transaction) MaxFee() string {
	return tx.GasPrice
}

// This is an alias property for GasPrice in order to support EIP1559
func (tx *Transaction) SetMaxFee(maxFee string) {
	tx.GasPrice = maxFee
}

func (tx *Transaction) GetRawTx() (*types.Transaction, error) {
	var (
		gasPrice, value, maxFeePerGas *big.Int // default nil

		nonce     uint64 = 0
		gasLimit  uint64 = 90000 // reference https://eth.wiki/json-rpc/API method eth_sendTransaction
		toAddress common.Address
		data      []byte
		valid     bool
		err       error
	)
	if tx.GasPrice != "" {
		if gasPrice, valid = big.NewInt(0).SetString(tx.GasPrice, 10); !valid {
			return nil, errors.New("invalid gasPrice")
		}
	}
	if tx.Value != "" {
		if value, valid = big.NewInt(0).SetString(tx.Value, 10); !valid {
			return nil, errors.New("invalid value")
		}
	}
	if tx.MaxPriorityFeePerGas != "" {
		if maxFeePerGas, valid = big.NewInt(0).SetString(tx.MaxPriorityFeePerGas, 10); !valid {
			return nil, errors.New("invalid max priority fee per gas")
		}
	}
	if tx.GasLimit != "" {
		if gasLimit, err = strconv.ParseUint(tx.GasLimit, 10, 64); err != nil {
			return nil, errors.New("invalid gas limit")
		}
	}
	if tx.To != "" && !common.IsHexAddress(tx.To) {
		return nil, errors.New("invalid toAddress")
	}
	toAddress = common.HexToAddress(tx.To)
	/*if len(tx.Data) > 0 {
		if data, err = util.HexDecodeString(tx.Data); err != nil {
			return nil, errors.New("invalid data string")
		}
	}*/

	if maxFeePerGas == nil || maxFeePerGas.Int64() == 0 {
		// is legacy tx
		return types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &toAddress,
			Value:    value,
			Gas:      gasLimit,
			GasPrice: gasPrice,
			Data:     data,
		}), nil
	} else {
		// is dynamic fee tx
		return types.NewTx(&types.DynamicFeeTx{
			Nonce:     nonce,
			To:        &toAddress,
			Value:     value,
			Gas:       gasLimit,
			GasFeeCap: gasPrice,
			GasTipCap: maxFeePerGas,
			Data:      data,
		}), nil
	}
}

// @return gasPrice * gasLimit + value
func (tx *Transaction) TotalAmount() string {
	priceInt, ok := big.NewInt(0).SetString(tx.GasPrice, 10)
	if !ok {
		return "0"
	}
	limitInt, ok := big.NewInt(0).SetString(tx.GasLimit, 10)
	if !ok {
		return "0"
	}
	amount, ok := big.NewInt(0).SetString(tx.Value, 10)
	if !ok {
		return "0"
	}
	return amount.Add(amount, priceInt.Mul(priceInt, limitInt)).String()
}
