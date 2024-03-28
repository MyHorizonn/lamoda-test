package postgres

import (
	"database/sql"
	"fmt"
	goods "lamoda-test/internal/storage"
	"log"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Client *sql.DB
}

func (p *Postgres) ReserveGood(good goods.Goods) error {
	log.Printf("start reserving %s good\n", good.Uuid)
	txn, err := p.Client.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = txn.Rollback()
	}()
	res, err := txn.Exec(`SELECT * FROM goods WHERE uuid = $1 FOR UPDATE;`, good.Uuid)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return fmt.Errorf("no goods with this uuid")
	}
	res, err = txn.Exec(`
		UPDATE 
			goods 
		SET
			amount = amount - $2
		WHERE 
			uuid = $1 AND amount >= $2`,
		good.Uuid, good.Amount)
	if err != nil {
		return nil
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return fmt.Errorf("can't reserve that much")
	}
	return txn.Commit()
}

func (p *Postgres) FreeGood(good goods.Goods) error {
	log.Printf("start freeing %s good\n", good.Uuid)
	txn, err := p.Client.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = txn.Rollback()
	}()
	res, err := txn.Exec(`SELECT * FROM goods WHERE uuid = $1 FOR UPDATE;`, good.Uuid)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return fmt.Errorf("no goods with this uuid")
	}
	_, err = txn.Exec(`
		UPDATE 
			goods 
		SET
			amount = amount + $2
		WHERE 
			uuid = $1`,
		good.Uuid, good.Amount)
	if err != nil {
		log.Printf("err: %v\n", err)
		return err
	}
	return txn.Commit()
}

func (p *Postgres) CheckGoods() ([]goods.Goods, error) {
	goodsResult := make([]goods.Goods, 0)
	good := goods.Goods{}
	rows, err := p.Client.Query("select * from goods")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&good.Uuid, &good.Name, &good.Size, &good.Amount)
		if err != nil {
			return goodsResult, err
		}
		goodsResult = append(goodsResult, good)
	}
	return goodsResult, err
}
