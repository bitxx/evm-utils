package signutil

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/storyicon/sigverify"
)

// VerifyEip721Signature 可以在opside_web_test.go中查看使用方式
func VerifyEip721Signature(address, signature string, typedData apitypes.TypedData) (bool, error) {
	//验证案例：
	// ok, err = VerifyEip721Signature(address, hexutil.Encode(signature), typedData)
	//	if err != nil {
	//		log.Fatal("VerifyEip721Signature fail.", err)
	//	}
	valid, err := sigverify.VerifyTypedDataHexSignatureEx(
		common.HexToAddress(address),
		typedData,
		signature,
	)
	return valid, err
}

func SignEip721(privateKey string, typedData *apitypes.TypedData) (string, error) {
	if privateKey == "" || typedData == nil {
		return "", errors.New("invalid parameter")
	}

	ecdsaPrivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}

	// 1、获取需要签名的数据的 Keccak-256 的哈希
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return "", err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return "", err
	}
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	sigHash := crypto.Keccak256(rawData)

	// 2、使用私钥签名哈希，得到签名
	signature, err := crypto.Sign(sigHash, ecdsaPrivateKey)
	if err != nil {
		return "", err
	}
	if signature[64] < 27 {
		signature[64] += 27
	}
	return hexutil.Encode(signature), nil
}

func MetamaskSignLogin(message string, privateKey string) (string, error) {
	ecdsaPrivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}
	fullMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(fullMessage))
	signatureBytes, err := crypto.Sign(hash.Bytes(), ecdsaPrivateKey)
	if err != nil {
		return "", err
	}
	signatureBytes[64] += 27
	return hexutil.Encode(signatureBytes), nil
}
