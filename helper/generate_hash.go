//go:build generate
// +build generate

package main

import (
	"fmt"
	"PROJECT_UAS/helper"
)

func main() {
	password := "admin123" // ubah sesuai kebutuhan
	hash, err := helper.HashPassword(password)
	if err != nil {
		panic(err)
	}

	fmt.Println("Password:", password)
	fmt.Println("Hash:", hash)
}