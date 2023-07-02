package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomId() string {
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	day := time.Now().Format("02")

	return fmt.Sprintf("%v%v%v%v", year, month, day, rand.Int())
}
