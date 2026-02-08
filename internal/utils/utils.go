package utils

import "strings"

func NormalizeOrder(order string) string {
	switch strings.ToUpper(order) {
	case "ASC":
		return "ASC"
	default:
		return "DESC"
	}
}
