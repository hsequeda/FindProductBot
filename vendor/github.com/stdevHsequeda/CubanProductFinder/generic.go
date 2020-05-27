package storeClient

import "time"

type GenericSection struct {
	Name      string
	Url       string
	Parent    string
	Store     *Store
	Priority  int
	ReadyTime time.Time
}

func (g *GenericSection) SetReadyTime(readyTime time.Time) {
	g.ReadyTime = readyTime
}

func (g *GenericSection) GetName() string {
	return g.Name
}

func (g *GenericSection) GetUrl() string {
	return g.Url
}

func (g *GenericSection) GetParent() string {
	return g.Parent
}

func (g *GenericSection) GetStore() *Store {
	return g.Store
}

func (g *GenericSection) GetPriority(string) int {
	return g.Priority
}

func (g *GenericSection) GetReadyTime() time.Time {
	return g.ReadyTime
}

type GenericProduct struct {
	Name    string
	Price   string
	Link    string
	Section Section
}

func (g *GenericProduct) GetName() string {
	return g.Name
}

func (g *GenericProduct) GetPrice() string {
	return g.Price
}

func (g *GenericProduct) GetLink() string {
	return g.Link
}

func (g *GenericProduct) GetSection() Section {
	return g.Section
}
