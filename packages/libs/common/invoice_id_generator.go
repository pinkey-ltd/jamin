package common

import (
	"fmt"
	"math/rand"
)

func GetInvoiceId() string {
	return fmt.Sprint(int64(rand.Int()))
}
