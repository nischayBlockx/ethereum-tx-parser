package util

import (
	"fmt"
	"regexp"
)

// Validate Ethereum contract address format
func ValidateAddress(address string) error {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(address) {
		return fmt.Errorf("input address [%s] is invalid", address)
	}

	return nil
}
