package ops

import (
	"fmt"
	"strings"
)

type ConcatOpsImpl struct {
	items []string
}

func ConcatOps() *ConcatOpsImpl {
	return &ConcatOpsImpl{}
}
func (s *ConcatOpsImpl) Concat(str ...string) *ConcatOpsImpl {
	s.items = append(s.items, str...)
	return s
}

func (s *ConcatOpsImpl) ConcatColumn(column string) *ConcatOpsImpl {
	s.items = append(s.items, column)
	return s

}

func (s *ConcatOpsImpl) ConcatString(plain string) *ConcatOpsImpl {
	s.items = append(s.items, fmt.Sprintf("'%s'", plain))
	return s
}

func (s *ConcatOpsImpl) Items() string {
	return strings.Join(s.items, ", ")
}

func (s *ConcatOpsImpl) Omitempty() *ConcatOpsImpl {
	var data []string
	for _, value := range s.items {
		if value != "" {
			data = append(data, value)
		}
	}
	s.items = data
	return s
}
func (s *ConcatOpsImpl) As(alias ...string) string {
	if len(alias) > 0 {
		return fmt.Sprintf("CONCAT(%s) AS `%s`", s.Omitempty().Items(), alias[0])
	}
	return fmt.Sprintf("CONCAT(%s)", s.Omitempty().Items())
}
