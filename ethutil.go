package ethutil

import (
	"github.com/bitxx/ethutil/model"
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

func (o *EthClient) TokenTransfer(privateKey, gasPrice, gasLimit, value, to string, data []byte) (hash string, err error) {
	chain, err := o.Chain()
	if err != nil {
		return "", err
	}
	token := model.NewToken(chain)
	return token.Transfer(privateKey, gasPrice, gasLimit, value, to, data)
}
