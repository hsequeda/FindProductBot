package products

type QuintaY42Product struct {
	Name      string `css:".product-name" extract:"attr" attr:"title"`
	Price     string `css:".product-price"`
	Link      string `css:".product-name" extract:"attr" attr:"href"`
	Available string `css:".ajax_add_to_cart_button" extract:"attr" attr:"href"`
	Store     string
}

func (q QuintaY42Product) GetName() string {
	return q.Name
}

func (q QuintaY42Product) GetPrice() string {
	return q.Price
}

func (q QuintaY42Product) GetLink() string {
	return q.Link
}

func (q QuintaY42Product) GetStore() string {
	return "5taY42"
}

func (q QuintaY42Product) IsAvailable() bool {
	return q.Available != ""
}
