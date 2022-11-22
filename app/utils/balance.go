package utils

import "math/big"

func ConvertBalance(wei *big.Int, decimals *big.Int) *big.Float {
	ten := big.NewInt(10)
	return new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(ten.Exp(ten, decimals, nil)))
}
