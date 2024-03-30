package goods

type Storage interface {
	ReserveGood(good Goods) error
	FreeGood(good Goods) error
	CheckGoods(store int) ([]Goods, error)
}

type Goods struct {
	Uuid   string `json:"uuid" db:"uuid"`
	Name   string `json:"name" db:"name"`
	Size   string `json:"size" db:"size"`
	Amount int    `json:"amount" db:"amount"`
}

type Store struct {
	Name          string `json:"name" db:"name"`
	Accessibility bool   `json:"accessibility" db:"accessibility"`
}
