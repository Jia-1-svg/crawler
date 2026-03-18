package pkg

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func CreateOrderSn() string {
	format := time.Now().Format("20060102150405")
	code := uuid.New().String()[:4]
	return fmt.Sprintf("%s%s", format, code)

}
