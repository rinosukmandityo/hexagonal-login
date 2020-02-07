package redis

import (
	"fmt"
)

func generateKey(code string) string {
	return fmt.Sprintf("login<>%s", code)
}

func generateUsernameKey(code string) string {
	return fmt.Sprintf("login<>username<>%s", code)
}
