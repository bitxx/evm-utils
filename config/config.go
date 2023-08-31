package config

// 默认gas limit估算失败后，21000 * 3 = 63000
const (
	DefaultContractGasLimit = "63000"
	DefaultEthGasLimit      = "21000"
	DefaultEthGasPrice      = "12000000000000" // 当前网络 standard gas price

	GasFactor = 1.8
)