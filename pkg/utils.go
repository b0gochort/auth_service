package pkg

import (
	"fmt"
	"math/rand"
)

//random 8 num
func GenerateCode() string {
	return fmt.Sprint(10000000 + rand.Intn(89999999))
}
