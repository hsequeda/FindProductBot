package storeClient

import "time"

type QuintaY42Section struct {
	Name      string `css:"a"`
	Url       string `css:"a" extract:"attr" attr:"href"`
	Parent    string
	Store     *Store
	Priority  int
	ReadyTime time.Time
}

func (q *QuintaY42Section) SetReadyTime(readyTime time.Time) {
	q.ReadyTime = readyTime
}

func (q *QuintaY42Section) GetName() string {
	return q.Name
}

func (q *QuintaY42Section) GetUrl() string {
	return q.Url
}

func (q *QuintaY42Section) GetParent() string {
	return q.Parent
}

func (q *QuintaY42Section) GetStore() *Store {
	return q.Store
}

func (q *QuintaY42Section) GetPriority(string) int {
	return q.Priority
}

func (q *QuintaY42Section) GetReadyTime() time.Time {
	return q.ReadyTime
}

type QuintaY42Product struct {
	Name      string `css:".product-name" extract:"attr" attr:"title"`
	Price     string `css:".product-price"`
	Link      string `css:".product-name" extract:"attr" attr:"href"`
	Available string `css:".ajax_add_to_cart_button" extract:"attr" attr:"href"`
	Section   Section
}

func (q *QuintaY42Product) GetName() string {
	return q.Name
}

func (q *QuintaY42Product) GetPrice() string {
	return q.Price
}

func (q *QuintaY42Product) GetLink() string {
	return q.Link
}

func (q *QuintaY42Product) GetSection() Section {
	return q.Section
}
