package models

type Context interface {
	CollectionName() string
	SetCollectionName(string)
}

type BaseContext struct {
	collectionName *string
}

func (c *BaseContext) CollectionName() string {
	return *c.collectionName
}

func (c *BaseContext) SetCollectionName(value string) {
	c.collectionName = &value
}

func NewContext() Context {
	return &BaseContext{nil}
}

func NewContextWithCollectionName(name string) Context {
	return &BaseContext{&name}
}
