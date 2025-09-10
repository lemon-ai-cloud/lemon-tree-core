package utils

import (
	"math/big"
	"strings"

	"github.com/google/uuid"
)

func ShortUUID(u string) (string, error) {
	// base62字符集
	const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// 解析UUID（会自动处理带-的36位UUID）
	parsed, err := uuid.Parse(u)
	if err != nil {
		return "", err
	}

	// 转换成大整数
	num := new(big.Int).SetBytes(parsed[:])

	// 转成base62
	if num.Cmp(big.NewInt(0)) == 0 {
		return "0", nil
	}

	base := big.NewInt(62)
	var result strings.Builder
	for num.Cmp(big.NewInt(0)) > 0 {
		mod := new(big.Int)
		num.DivMod(num, base, mod)
		result.WriteByte(base62Chars[mod.Int64()])
	}

	// 由于是反向拼接的，需要翻转
	out := []rune(result.String())
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}

	return string(out), nil
}
