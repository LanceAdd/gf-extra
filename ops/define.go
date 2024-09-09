package ops

type FieldsInterface interface {
	GetQueryFields() []string
	GetGroupFields() []string
}
