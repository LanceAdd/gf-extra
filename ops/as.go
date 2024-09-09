package ops

import "fmt"

func As(plain string, alias string) string {
	return fmt.Sprintf("%s AS `%s`", plain, alias)
}
