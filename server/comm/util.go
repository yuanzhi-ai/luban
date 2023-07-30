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

// 判断手机号是否合法
func IsPhoneLegal(phone string) bool {
	mobileExp := `^1\d{10}$`
	mobileReg := regexp.MustCompile(mobileExp)
	return mobileReg.MatchString(phone)
}

const (
	minPswdLne   = 8  // 最小密码长度
	maxPswdLen   = 16 //最大密码长度
	minPswdLevel = 2  //最小密码强度
)

// 判断密码是否合法
func IsPswdLegal(pswd string) bool {
	// 密码长度
	if len(pswd) < minPswdLne || len(pswd) > maxPswdLen {
		return false
	}
	// 过滤掉这四类字符以外的密码串,直接判断不合法
	re, err := regexp.Compile(`^[a-zA-Z0-9.@$!%*#_~?&^]{8,16}$`)
	if err != nil {
		return false
	}
	match := re.MatchString(pswd)
	if !match {
		return false
	}
	// 密码强度
	var level = 0
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[.@$!%*#_~?&^]+`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, pswd)
		if match {
			level++
		}
	}
	if level < minPswdLevel {
		return false
	}
	return true
}
