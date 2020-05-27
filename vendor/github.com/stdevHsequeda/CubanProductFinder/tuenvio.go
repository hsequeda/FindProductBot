package storeClient

import "time"

type TuEnvioSection struct {
	Name      string `css:"div ul li a"`
	Url       string `css:"div ul li a" extract:"attr" attr:"href"`
	Parent    string
	Store     *Store
	Priority  int
	ReadyTime time.Time
}

func (t *TuEnvioSection) SetReadyTime(readyTime time.Time) {
	t.ReadyTime = readyTime
}

func (t *TuEnvioSection) GetName() string {
	return t.Name
}

func (t *TuEnvioSection) GetUrl() string {
	return t.Url
}

func (t *TuEnvioSection) GetParent() string {
	return t.Parent
}

func (t *TuEnvioSection) GetStore() *Store {
	return t.Store
}

func (t *TuEnvioSection) GetPriority(string) int {
	return t.Priority
}

func (t *TuEnvioSection) GetReadyTime() time.Time {
	return t.ReadyTime
}

type TuEnvioProduct struct {
	Name    string `css:".thumbTitle",redis:"name"`
	Price   string `css:".thumbPrice",redis:"price"`
	Link    string `css:".thumbnail a" extract:"attr" attr:"href",redis:"link"`
	Section Section
}

func (t *TuEnvioProduct) GetName() string {
	return t.Name
}

func (t *TuEnvioProduct) GetPrice() string {
	return t.Price
}

func (t *TuEnvioProduct) GetLink() string {
	return t.Link
}

func (t *TuEnvioProduct) GetSection() Section {
	return t.Section
}
