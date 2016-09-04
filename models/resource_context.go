package models

type ResourceContext interface {
	Context
	IdParameter() *string
	SetIdParameter(*string)
	IsEmpty() bool
	SetIsEmpty(bool)
	JoinFields() []JoinField
	SetJoinFields([]JoinField)
}

type BaseResourceContext struct {
	Context
	idParameter *string
	isEmpty     bool
	joinFields  []JoinField
}

func (c *BaseResourceContext) IdParameter() *string {
	return c.idParameter
}

func (c *BaseResourceContext) SetIdParameter(value *string) {
	c.idParameter = value
}

func (c *BaseResourceContext) IsEmpty() bool {
	return c.isEmpty
}

func (c *BaseResourceContext) SetIsEmpty(value bool) {
	c.isEmpty = value
}

func (c *BaseResourceContext) JoinFields() []JoinField {
	return c.joinFields
}

func (c *BaseResourceContext) SetJoinFields(value []JoinField) {
	c.joinFields = value
}

func NewResourceContext() ResourceContext {
	return &BaseResourceContext{NewContext(), nil, false, []JoinField{}}
}
