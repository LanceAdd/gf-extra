package ops

import (
	"fmt"
	"strings"
)

type ConcatWsOpsImpl struct {
	items []string
	split string
}

func ConcatWsOps() *ConcatWsOpsImpl {
	return &ConcatWsOpsImpl{}
}
func (s *ConcatWsOpsImpl) ConcatWs(split string, str ...string) *ConcatWsOpsImpl {
	if split == "" || len(str) == 0 {
		return s
	}
	s.split = split
	s.items = append(s.items, str...)
	return s
}

func (s *ConcatWsOpsImpl) ConcatWsSplit(split string) *ConcatWsOpsImpl {
	if split == "" {
		return s
	}
	s.split = split
	return s
}

func (s *ConcatWsOpsImpl) ConcatColumn(column string) *ConcatWsOpsImpl {
	s.items = append(s.items, column)
	return s

}

func (s *ConcatWsOpsImpl) ConcatString(plain string) *ConcatWsOpsImpl {
	s.items = append(s.items, fmt.Sprintf("'%s'", plain))
	return s
}

func (s *ConcatWsOpsImpl) Items() string {
	return strings.Join(s.items, ", ")
}

func (s *ConcatWsOpsImpl) GetSplitString() string {
	return s.split
}

func (s *ConcatWsOpsImpl) Omitempty() *ConcatWsOpsImpl {
	var data []string
	for _, value := range s.items {
		if value != "" {
			data = append(data, value)
		}
	}
	s.items = data
	return s
}
func (s *ConcatWsOpsImpl) As(alias ...string) string {
	if len(alias) > 0 {
		return fmt.Sprintf("CONCAT_WS('%s', %s) AS `%s`", s.split, s.Omitempty().Items(), alias[0])
	}
	return fmt.Sprintf("CONCAT_WS('%s', %s)", s.split, s.Omitempty().Items())
}
