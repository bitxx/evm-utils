package model

import (
	"context"
	"errors"
	"github.com/bitxx/ethutil/model/types"
	"github.com/bitxx/ethutil/util"
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

func (t *Token) Transfer(privateKey, gasPrice, gasLimit, value, to, data string) (hash string, err error) {
	if gasPrice == "" || gasLimit == "" || to == "" || value == "" {
		return "", errors.New("param is error")
	}
	tx := types.NewTransaction(0, gasPrice, gasLimit, to, value, data)

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

func (t *Token) EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string) (string, error) {
	msg := types.NewCallMsg()
	msg.SetFrom(fromAddress)
	msg.SetTo(receiverAddress)
	msg.SetGasPrice(gasPrice)
	msg.SetValue(amount)
	return t.chain.EstimateGasLimit(msg)
}
