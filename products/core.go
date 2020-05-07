package products

type Product interface {
	GetName() string
	GetPrice() string
	GetLink() string
	GetStore() string
	IsAvailable() bool
}
