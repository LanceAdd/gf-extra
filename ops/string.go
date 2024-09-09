package ops

import (
	"fmt"
)

type StringOpsImpl struct {
	item string
}

func StringOps() *StringOpsImpl {
	return &StringOpsImpl{}
}
func (s *StringOpsImpl) FromStr(str string) *StringOpsImpl {
	s.item = fmt.Sprintf("'%s'", str)
	return s
}

func (s *StringOpsImpl) FromColumn(column string) *StringOpsImpl {
	s.item = fmt.Sprintf("%s", column)
	return s
}

func (s *StringOpsImpl) SubString(start int, length int) *StringOpsImpl {
	s.item = fmt.Sprintf("SUBSTRING(%s, %d, %d)", s.item, start, length)
	return s
}

func (s *StringOpsImpl) Trim() *StringOpsImpl {
	s.item = fmt.Sprintf("TRIM(%s)", s.item)
	return s
}

func (s *StringOpsImpl) Replace(find string, replace string) *StringOpsImpl {
	s.item = fmt.Sprintf("REPLACE(%s, '%s', '%s')", s.item, find, replace)
	return s
}

func (s *StringOpsImpl) Upper() *StringOpsImpl {
	s.item = fmt.Sprintf("UPPER(%s)", s.item)
	return s
}

func (s *StringOpsImpl) Lower() *StringOpsImpl {
	s.item = fmt.Sprintf("LOWER(%s)", s.item)
	return s
}

func (s *StringOpsImpl) Length() *StringOpsImpl {
	s.item = fmt.Sprintf("LENGTH(%s)", s.item)
	return s
}

func (s *StringOpsImpl) Locate(find string, start int) *StringOpsImpl {
	s.item = fmt.Sprintf("LOCATE(%s, %s, %d)", find, s.item, start)
	return s
}

func (s *StringOpsImpl) Instr(find string) *StringOpsImpl {
	s.item = fmt.Sprintf("INSTR(%s, %s)", s.item, find)
	return s
}

func (s *StringOpsImpl) Reverse() *StringOpsImpl {
	s.item = fmt.Sprintf("REVERSE(%s)", s.item)
	return s
}

func (s *StringOpsImpl) Left(length int) *StringOpsImpl {
	s.item = fmt.Sprintf("LEFT(%s, %d)", s.item, length)
	return s
}

func (s *StringOpsImpl) Right(length int) *StringOpsImpl {
	s.item = fmt.Sprintf("RIGHT(%s, %d)", s.item, length)
	return s
}

func (s *StringOpsImpl) As(alias ...string) string {
	if IsBlank(s.item) {
		return ""
	}
	if len(alias) > 0 {
		return fmt.Sprintf("%s AS %s", s.item, alias[0])
	}
	return s.item
}
