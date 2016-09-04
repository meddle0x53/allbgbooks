package models

type Context interface {
	CollectionName() string
	SetCollectionName(string)
}
