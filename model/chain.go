package model

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/bitxx/ethutil/config"
	"github.com/bitxx/ethutil/model/types"
	"github.com/bitxx/ethutil/util"
	"github.com/ethereum/go-ethereum/common"
	eTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"strconv"
	"sync"
	"time"
)

var chainConnections = make(map[string]*Chain)
var lock sync.RWMutex

type Chain struct {
	RemoteRpcClient *ethclient.Client
	Timeout         int64
	rpcClient       *rpc.Client
	chainId         *big.Int
	rpcUrl          string
}

// GetChain
//
//	@Description: get connect from cache
//	@param rpcUrl
//	@param timeout
//	@return *EthChain
//	@return error
func GetChain(rpcUrl string, timeout int64) (*Chain, error) {
	if rpcUrl == "" {
		return nil, errors.New("rpc url can't empty")
	}

	chain, ok := chainConnections[rpcUrl]
	if ok {
		return chain, nil
	}

	// 通过加锁范围
	lock.Lock()
	defer lock.Unlock()

	// 再判断一次
	chain, ok = chainConnections[rpcUrl]
	if ok {
		return chain, nil
	}

	// 创建并存储
	chain, err := newChain(rpcUrl, timeout)
	if err != nil {
		return nil, err
	}

	chainConnections[rpcUrl] = chain
	return chain, nil
}

// newChain
//
//	@Description:
//	@param timeout the net connect time, second,default is 60
//	@return *Chain
func newChain(rpcUrl string, timeout int64) (chain *Chain, err error) {
	if timeout <= 0 {
		timeout = 60
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	rpcClient, err := rpc.DialContext(ctx, rpcUrl)
	if err != nil {
		return
	}

	remoteRpcClient := ethclient.NewClient(rpcClient)
	chainId, err := remoteRpcClient.ChainID(ctx)
	if err != nil {
		return
	}

	chain = &Chain{
		chainId:         chainId,
		rpcClient:       rpcClient,
		RemoteRpcClient: remoteRpcClient,
		rpcUrl:          rpcUrl,
		Timeout:         timeout,
	}
	return
}

func (c *Chain) Close() {
	if c.RemoteRpcClient != nil {
		c.RemoteRpcClient.Close()
	}
	if c.rpcClient != nil {
		c.rpcClient.Close()
	}
}

func (c *Chain) EstimateGasLimit(msg *types.CallMsg) (gas string, err error) {

	if len(msg.Msg.Data) > 0 {
		// any contract transaction
		gas = config.DefaultContractGasLimit
	} else {
		// nomal transfer
		gas = config.DefaultEthGasLimit
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()
	gasLimit, err := c.RemoteRpcClient.EstimateGas(ctx, msg.Msg)
	if err != nil {
		return
	}
	gasString := ""
	if len(msg.Msg.Data) > 0 {
		gasFloat := big.NewFloat(0).SetUint64(gasLimit)
		gasFloat = gasFloat.Mul(gasFloat, big.NewFloat(config.GasFactor))
		gasInt, _ := gasFloat.Int(nil)
		gasString = gasInt.String()
	} else {
		gasString = strconv.FormatUint(gasLimit, 10)
	}

	return gasString, nil
}

func (c *Chain) Nonce(spenderAddressHex string) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()
	nonce, err := c.RemoteRpcClient.PendingNonceAt(ctx, common.HexToAddress(spenderAddressHex))
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

// BuildTxUnSign
//
//	@Description: the transaction no sign
//	@receiver c
//	@param privateKey
//	@param transaction
//	@return *eTypes.Transaction
//	@return error
func (c *Chain) BuildTxUnSign(address string, transaction *types.Transaction) (*eTypes.Transaction, error) {
	if transaction.Nonce <= 0 {
		if !util.IsValidAddress(address) {
			return nil, errors.New("address format is error")
		}
		nonce, err := c.Nonce(address)
		if err != nil {
			nonce = 0
			err = nil
		}
		transaction.Nonce = nonce
	}
	return transaction.GetRawTx()
}

func (c *Chain) BuildTxSign(privateKey *ecdsa.PrivateKey, txNoSign *eTypes.Transaction) (*types.BuildTxResult, error) {
	if privateKey == nil || txNoSign == nil {
		return nil, errors.New("param is empty")
	}
	signedTx, err := eTypes.SignTx(txNoSign, eTypes.LatestSignerForChainID(c.chainId), privateKey)
	if err != nil {
		return nil, err
	}
	return &types.BuildTxResult{
		SignedTx: signedTx,
		TxHex:    signedTx.Hash().String(),
	}, nil
}

func (c *Chain) SendTx(signedTx *eTypes.Transaction) error {
	if signedTx == nil {
		return errors.New("signed transaction can't be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()
	err := c.RemoteRpcClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return err
	}
	return nil
}
