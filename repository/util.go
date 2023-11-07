package repository

import "strings"

func GoSliceToPostgresArray(arr []string) string {
	return "{" + strings.Join(arr, ",") + "}"
}
