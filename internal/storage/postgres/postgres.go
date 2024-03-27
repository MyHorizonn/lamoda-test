package postgres

import (
	"database/sql"
	goods "lamoda-test/internal/storage"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Client *sql.DB
}

func (p *Postgres) ReserveGoods(good goods.Goods) error {
	_, err := p.Client.Exec("")
	return err
}

func (p *Postgres) FreeGoods(good goods.Goods) error {
	_, err := p.Client.Exec("")
	return err
}

func (p *Postgres) CheckGood() ([]goods.Goods, error) {
	_, err := p.Client.Exec("")
	return []goods.Goods{}, err
}
