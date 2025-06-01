package _utility

import cowsay "github.com/Code-Hex/Neo-cowsay/v2"

func IsValidCowName(name string) bool {
	if name == "default" {
		return true
	}
	availableCows := cowsay.CowsInBinary()
	for _, validName := range availableCows {
		if name == validName {
			return true
		}
	}
	return false
}
