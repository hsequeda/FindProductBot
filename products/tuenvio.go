package products

type TuEnvioProduct struct {
	Name  string `css:".thumbTitle"`
	Price string `css:".thumbPrice"`
	Link  string `css:".thumbnail a" extract:"attr" attr:"href"`
	Store string
}

func (t TuEnvioProduct) GetName() string {
	return t.Name
}

func (t TuEnvioProduct) GetPrice() string {
	return t.Price
}

func (t TuEnvioProduct) GetLink() string {
	return t.Link
}

func (t TuEnvioProduct) GetStore() string {
	return t.Store
}

func (t TuEnvioProduct) IsAvailable() bool {
	return true
}
