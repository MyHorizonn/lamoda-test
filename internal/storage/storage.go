package goods

type Storage interface {
	ReserveGoods(good Goods) error
	FreeGoods(good Goods) error
	CheckGood() ([]Goods, error)
}

type Goods struct {
	Uuid   string `json:"uuid" db:"uuid"`
	Name   string `json:"name" db:"name"`
	Size   string `json:"size" db:"size"`
	Amount int    `json:"amount" db:"amount"`
}

type Store struct {
	Name          string `json:"name" db:"name"`
	Accessibility bool   `json:"bool" db:"bool"`
}
