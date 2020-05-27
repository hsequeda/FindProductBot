package storeClient

type Store struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	// Address        string `json:"address"`
	Province string `json:"province"`
	Online   bool   `json:"online"`
	// PickUpOnStore  bool   `json:"pickUpOnStore"`
	// HomeDelivery   bool   `json:"homeDelivery"`
	// FreezeDelivery string `json:"freezeDelivery"`
	// DeliveryTime   string `json:"deliveryTime"`
	// Cost           string `json:"cost"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Url   string `json:"url"`
	// Cadena         string `json:"cadena"`
}
