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

func (t *Token) Transfer(privateKey, gasPrice, gasLimit, maxPriorityFeePerGas, value, to, data string) (hash string, err error) {
	if gasPrice == "" || gasLimit == "" || to == "" || value == "" {
		return "", errors.New("param is error")
	}
	tx := types.NewTransaction("", gasPrice, gasLimit, maxPriorityFeePerGas, to, value, data)

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

// 现实：0x02f901b814860ae9f7bcc0008302a7fc94bd927011759b2c4f2602c3008f8ef3407db5347388b469471f80140000b90145 73c45c98000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000540000000000000000000000000000000000000000000000000de0b6b3a764000000000000000000000000000000000000000000000000000000000000000000147a547a149a79a03f4dd441b6806ffcbb1b63f3830000000000000000000000000000000000000000000000000000000000000000000000000000000000000008a688906bd8b0000000000000000000000000000000000000000000000000000083018ff7a09b057a4f9ebae5c92483b37acddbe9df12acbf195a939008d38b61b0e47fd106a04e1397d8e06d136209b64890e3f2841cbf9a726e5809f3355ab1344b1edd780a
// 标准：0x02f901c082c7ea82040c851695a68a0086012dfb396a00830179b5 94bd927011759b2c4f2602c3008f8ef3407db5347388b469471f80140000b90144 73c45c98000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000540000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000000000000000014e5e69b292170459a4e4cc77f94491681ff1f16360000000000000000000000000000000000000000000000000000000000000000000000000000000000000008a688906bd8b00000000000000000000000000000000000000000000000000000c080a0136e287d6d49355d79d366cd3e5fadd640a19b983a199eee548a4e2cade729d7a01cfcb8d6222d54debfd95e71849583f884e08099055814f1b7abf5fe212d93a3
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
