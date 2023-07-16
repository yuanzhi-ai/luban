package comm

import (
	"crypto/rand"
	"math"
	"math/big"
	"regexp"
	"strconv"
)

func RangeRand(min int64, max int64) int64 {
	left := min
	right := max
	if min > max {
		left = max
		right = min
	}
	if left < 0 {
		f64Min := math.Abs(float64(left))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(right+1+i64Min))
		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(right-left+1))

		return left + result.Int64()
	}
}

func GetRandDigitStr(length int) string {
	chars := "0123456789"
	randStr := ""
	for i := 0; i < length; i++ {
		randIndex := RangeRand(0, int64(len(chars)-1))
		randStr += chars[randIndex : randIndex+1]
	}
	return randStr
}

// 判读手机验证码是否合法
func IsPhoneCodeLegal(code string, needLen int) bool {
	if len(code) != needLen {
		return false
	}
	_, err := strconv.Atoi(code)
	return err == nil
}

func IsPhoneLegal(phone string) bool {
	mobileExp := `^(1[3-9]d{9})$`
	mobileReg := regexp.MustCompile(mobileExp)
	return mobileReg.MatchString(phone)
}
