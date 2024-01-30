package evmutils

import (
	"errors"
	"github.com/bitxx/evm-utils/model"
	"github.com/bitxx/evm-utils/util/signutil"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type EthClient struct {
	RpcUrl  string
	timeout int64
}

// NewEthClient
//
//	@Description: if rpcUrl and timeout is empty,you can‘t connect the node,but you can use the function about wallet
//	@param rpcUrl
//	@param timeout
//	@return *EthClient
func NewEthClient(rpcUrl string, timeout int64) *EthClient {
	return &EthClient{
		RpcUrl:  rpcUrl,
		timeout: timeout,
	}
}

// NewSimpleEthClient
//
//	@Description: not support connect to the node
//	@return *EthClient
func NewSimpleEthClient() *EthClient {
	return &EthClient{}
}

func (o *EthClient) AccountByMnemonic() (account *model.Account, err error) {
	return model.NewAccount().AccountByMnemonic()
}

// AccountInfoByMnemonic
//
//	@Description:
//	@receiver o
//	@param mnemonic
//	@return account
//	@return err
func (o *EthClient) AccountInfoByMnemonic(mnemonic string) (account *model.Account, err error) {
	return model.NewAccount().AccountInfoByMnemonic(mnemonic)
}

func (o *EthClient) AccountWithPrivateKey(privateKey string) (account *model.Account, err error) {
	return model.NewAccount().AccountWithPrivateKey(privateKey)
}

func (o *EthClient) TokenBalanceOf(address string) (balance string, err error) {
	chain, err := o.Chain()
	if err != nil {
		return "", err
	}
	token := model.NewToken(chain)
	return token.BalanceOf(address)
}

// TokenEstimateGasLimit
//
//	@Description: 估算gas ，如果是合约地址，data肯定不得为空
//	@receiver o
//	@param fromAddress
//	@param receiverAddress
//	@param gasPrice
//	@param amount
//	@param data
//	@return balance
//	@return err
func (o *EthClient) TokenEstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string, data []byte) (balance string, err error) {
	chain, err := o.Chain()
	if err != nil {
		return "", err
	}
	token := model.NewToken(chain)
	return token.EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount, data)
}

func (o *EthClient) Chain() (*model.Chain, error) {
	return model.GetChain(o.RpcUrl, o.timeout)
}

func (o *EthClient) Nonce(address string) (nonce uint64, err error) {
	chain, err := o.Chain()
	if err != nil {
		return 0, err
	}
	return chain.Nonce(address)
}

func (o *EthClient) TokenTransfer(privateKey, nonce, gasPrice, gasLimit, maxPriorityFeePerGas, value, to, data string) (hash string, err error) {
	chain, err := o.Chain()
	if err != nil {
		return "", err
	}
	token := model.NewToken(chain)
	return token.Transfer(privateKey, nonce, gasPrice, gasLimit, maxPriorityFeePerGas, value, to, data)
}

// TransactionsByBlockNum
//
//	@Description: 获取一个块的所有交易
//	@receiver o
//	@param number
//	@return []model.Transaction
//	@return error
func (o *EthClient) TransactionsByBlockNum(number uint64) ([]model.Transaction, error) {
	chain, err := o.Chain()
	if err != nil {
		return nil, err
	}
	transaction := model.NewTransaction(chain)
	return transaction.TransactionsByBlockNumber(number)
}

// MetamaskSignLogin
//
//	@Description: metamask sign login
//	@receiver o
//	@param message
//	@param privateKey
//	@return string
//	@return error
func (o *EthClient) MetamaskSignLogin(message, privateKey string) (string, error) {
	if message == "" || privateKey == "" {
		return "", errors.New("param is empty")
	}
	return signutil.MetamaskSignLogin(message, privateKey)
}

// SignEip721
//
//	@Description: eip721 sign
//	@receiver o
//	@param privateKey
//	@param typedData
//	@return string
//	@return error
func (o *EthClient) SignEip721(privateKey string, typedData *apitypes.TypedData) (string, error) {
	if typedData == nil || privateKey == "" {
		return "", errors.New("param is empty")
	}
	return signutil.SignEip721(privateKey, typedData)
}
