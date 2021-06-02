package utils

import (
	"fmt"
	"testing"
)

func TestPassword(*testing.T) {
	pwd := "test"
	hash, _ := HashPassword(pwd)

	fmt.Println(pwd, hash)
	fmt.Println(VerifyPassword(hash, pwd))
}
