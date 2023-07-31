package ethutil

import (
	"fmt"

	"github.com/bitxx/ethutil/config"

	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

const (
	rpcUrl  = "以太坊生态网络地址"
	timeout = 60 //second

	testAccountFromAddress           = "0x7a547A149A79A03F4dd441B6806ffCBb1b63F383"
	testAccountFromAddressPrivateKey = ""
	testAccountToAddress             = "0x8B63293748e058F47a31c0D2Af0B1b3FeDdc4D4C"
	accountFile                      = "./account.txt"
	addressFile                      = "./address.txt"
	privateKeyFile                   = "./privateKey.txt"
)

func MyClient() *EthClient {
	return NewEthClient(rpcUrl, timeout)
}

func TestAccountByMnemonic(t *testing.T) {
	num := 20
	result := "account：\n"
	addresses := "\n\n\nSummary of the above account addresses:"
	for i := 0; i < num; i++ {
		account, _ := MyClient().AccountByMnemonic()
		addresses = addresses + "\n" + account.Address
		result = result + fmt.Sprintf("NO. %d group account：\nmnemonic：%s\naddress：%s\nprivateKey：%s\npublicKey：%s\n\n", i+1, account.Mnemonic, account.Address, account.PrivateKey, account.PublicKey)
	}
	result = result + addresses
	_ = os.WriteFile(accountFile, []byte(result), 0666)
}

func TestAccountWithPrivateKey(t *testing.T) {
	account1, _ := MyClient().AccountByMnemonic()
	account2, _ := MyClient().AccountWithPrivateKey(account1.PrivateKey)
	require.Equal(t, account1.PublicKey, account2.PublicKey)
}

func TestTokenBalance(t *testing.T) {
	balance, err := MyClient().TokenBalanceOf(testAccountToAddress)
	require.Nil(t, err)
	t.Log("balance：", balance)
}

func TestBatchTokenBalance(t *testing.T) {
	bytes, err := os.ReadFile(addressFile)
	require.Nil(t, err)
	addresses := strings.Split(string(bytes), "\n")
	client := MyClient()
	for i, address := range addresses {
		address = strings.TrimSpace(address)
		if len(address) <= 0 {
			continue
		}
		balance, err := client.TokenBalanceOf(address)
		if err != nil {
			t.Error(fmt.Sprintf("index: %d ,address:%s, request err: %s", i, address, err.Error()))
			break
		}
		/*value, err := decimal.NewFromString(balance)
		if err != nil {
			t.Error(fmt.Sprintf("index: %d ,address:%s, value err: %s", i, address, err.Error()))
			break
		}
		if value.Cmp(decimal.Zero) > 0 {
			println("the", i+1, "group,address: ", address)
		}*/

		t.Log(fmt.Sprintf("address: %s,balance: %s", address, balance))
	}
}

func TestNonce(t *testing.T) {
	nonce, err := MyClient().Nonce(testAccountToAddress)
	require.Nil(t, err)
	t.Log("nonce: ", nonce)
}

func TestTokenEstimateGasLimit(t *testing.T) {
	value := "1000000000000000000"
	gasLimit, err := MyClient().TokenEstimateGasLimit(testAccountFromAddress, testAccountToAddress, config.DefaultEthGasPrice, value)
	require.Nil(t, err)
	t.Log("estimate gas limit: ", gasLimit)
}

func TestTokenTransfer(t *testing.T) {
	value := "1000000000000000000"
	gasLimit, err := MyClient().TokenEstimateGasLimit(testAccountFromAddress, testAccountToAddress, config.DefaultEthGasPrice, value)
	require.Nil(t, err)
	hash, err := MyClient().TokenTransfer(testAccountFromAddressPrivateKey, config.DefaultEthGasPrice, gasLimit, value, testAccountToAddress, "")
	require.Nil(t, err)
	t.Log("hash:", hash)
}

func TestBatchTokenTransfer(t *testing.T) {
	toAddress := ""
	value := "25000000000000000000000"
	gasLimit := "21000"
	bytes, err := os.ReadFile(privateKeyFile)
	require.Nil(t, err)
	privateKeys := strings.Split(string(bytes), "\n")
	client := MyClient()
	for i, privateKey := range privateKeys {
		privateKey = strings.TrimSpace(privateKey)
		if len(privateKey) <= 0 {
			continue
		}
		hash, err := client.TokenTransfer(privateKey, config.DefaultEthGasPrice, gasLimit, value, toAddress, "")
		if err != nil {
			t.Error(fmt.Sprintf("index:%d,privateKey: %s,error: %s", i, privateKey, err.Error()))
			continue
		}
		t.Log(fmt.Sprintf("index: %d,privateKey: %s,hash: %s", i, privateKey, hash))
	}
}
