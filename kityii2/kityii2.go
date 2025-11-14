package kityii2

import "golang.org/x/crypto/bcrypt"

// GeneratePasswordHash 生成密码哈希
func GeneratePasswordHash(password string) (_hash string) {
    fromPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(fromPassword)
}

// ValidatePasswordHash 验证密码与哈希是否匹配
func ValidatePasswordHash(password, passwordHash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
}
