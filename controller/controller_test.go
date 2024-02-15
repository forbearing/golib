package controller

import (
	"fmt"
	"testing"
)

func TestEncryptPasswd(t *testing.T) {
	fmt.Println(encryptPasswd("admin"))
}
