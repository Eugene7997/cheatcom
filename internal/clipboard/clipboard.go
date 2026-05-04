package clipboard

import "github.com/atotto/clipboard"

func Copy(s string) error {
	return clipboard.WriteAll(s)
}
