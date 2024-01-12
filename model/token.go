package model

import (
	"context"
	"errors"
	"github.com/bitxx/evm-utils/model/types"
	"github.com/bitxx/evm-utils/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"time"
)

type Token struct {
	chain *Chain
}

func NewToken(chain *Chain) *Token {
	return &Token{
		chain: chain,
	}
}

func (t *Token) BalanceOf(address string) (balance string, err error) {
	if t.chain == nil {
		return "", errors.New("the chain node is empty")
	}
	if !util.IsValidAddress(address) {
		return "", errors.New("invalid hex address")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.chain.Timeout)*time.Second)
	defer cancel()
	balanceResult, err := t.chain.RemoteRpcClient.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return "", err
	}
	return balanceResult.String(), nil
}

func (t *Token) Transfer(privateKey, nonce, gasPrice, gasLimit, maxPriorityFeePerGas, value, to, data string) (hash string, err error) {
	if gasPrice == "" || gasLimit == "" || to == "" || value == "" {
		return "", errors.New("param is error")
	}
	tx := types.NewTransaction(nonce, gasPrice, gasLimit, maxPriorityFeePerGas, to, value, data)

	priData, err := util.HexDecodeString(privateKey)
	if err != nil {
		return "", err
	}
	privateKeyECDSA, err := crypto.ToECDSA(priData)
	address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex()

	//get no sign tx
	txUnSign, err := t.chain.BuildTxUnSign(address, tx)
	if err != nil {
		return "", err
	}

	//tx sign
	txSign, err := t.chain.BuildTxSign(privateKeyECDSA, txUnSign)
	if err != nil {
		return "", err
	}

	//send tx
	return txSign.TxHex, t.chain.SendTx(txSign.SignedTx)
}

func (t *Token) EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string, data []byte) (string, error) {
	msg := types.NewCallMsg()
	msg.SetFrom(fromAddress)
	msg.SetTo(receiverAddress)
	msg.SetGasPrice(gasPrice)
	msg.SetValue(amount)
	if data != nil {
		msg.SetData(data)
	}
	return t.chain.EstimateGasLimit(msg)
}
