package models

type CollectionContext interface {
	Context
	PerPage() uint64
	SetPerPage(uint64)
	Page() uint64
	SetPage(uint64)
	OrderBy() string
	SetOrderBy(string)
	FilteringValues() []FilteringValue
	SetFilteringValues([]FilteringValue)
	IgnoreCase() bool
	SetIgnoreCase(bool)
}

type BaseCollectionContext struct {
	Context
	PerPageValue         uint64
	PageValue            uint64
	OrderByValue         string
	FilteringValuesValue []FilteringValue
	IgnoreCaseValue      bool
}

func (c *BaseCollectionContext) PerPage() uint64 {
	return c.PerPageValue
}

func (c *BaseCollectionContext) SetPerPage(value uint64) {
	c.PerPageValue = value
}

func (c *BaseCollectionContext) Page() uint64 {
	return c.PageValue
}

func (c *BaseCollectionContext) SetPage(value uint64) {
	c.PageValue = value
}

func (c *BaseCollectionContext) OrderBy() string {
	return c.OrderByValue
}

func (c *BaseCollectionContext) SetOrderBy(value string) {
	c.OrderByValue = value
}

func (c *BaseCollectionContext) FilteringValues() []FilteringValue {
	return c.FilteringValuesValue
}

func (c *BaseCollectionContext) SetFilteringValues(value []FilteringValue) {
	c.FilteringValuesValue = value
}

func (c *BaseCollectionContext) IgnoreCase() bool {
	return c.IgnoreCaseValue
}

func (c *BaseCollectionContext) SetIgnoreCase(value bool) {
	c.IgnoreCaseValue = value
}

func NewCollectionContext() CollectionContext {
	return &BaseCollectionContext{
		NewContext(), 10, 1, "id", []FilteringValue{}, false,
	}
}
