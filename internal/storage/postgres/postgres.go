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

func (p *Postgres) ReserveGood(good goods.Goods, store int) error {
	log.Printf("start reserving %s good\n", good.Uuid)
	txn, err := p.Client.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = txn.Rollback()
	}()
	res, err := txn.Exec(`
	SELECT store_id, goods_uuid, amount, reserved 
	FROM goods_in_store a
	JOIN store b ON a.store_id = b.id
	WHERE goods_uuid = $1 AND store_id = $2 AND b.accessibility
	FOR UPDATE;`,
		good.Uuid, store)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return fmt.Errorf("no goods with this uuid or can't get store with this id or store is not avaliable in that moment")
	}
	res, err = txn.Exec(`
	UPDATE
	goods_in_store
	SET
	reserved = reserved + $2
	WHERE
	goods_uuid = $1 AND store_id = $3 AND amount - reserved >= $2`,
		good.Uuid, good.Amount, store)
	if err != nil {
		return err
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

func (p *Postgres) FreeGood(good goods.Goods, store int) error {
	log.Printf("start freeing %s good\n", good.Uuid)
	txn, err := p.Client.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = txn.Rollback()
	}()
	res, err := txn.Exec(`
	SELECT * 
	FROM goods_in_store a
	JOIN store b ON a.store_id = b.id
	WHERE goods_uuid = $1 AND store_id = $2 AND b.accessibility
	FOR UPDATE;`,
		good.Uuid, store)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return fmt.Errorf("no goods with this uuid or can't get store with this id or store is not avaliable in that moment")
	}
	res, err = txn.Exec(`
	UPDATE 
	goods_in_store
	SET
	reserved = reserved - $2
	WHERE 
	goods_uuid = $1 AND store_id = $3 AND reserved >= $2`,
		good.Uuid, good.Amount, store)
	if err != nil {
		return err
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return fmt.Errorf("can't free that much")
	}
	return txn.Commit()
}

func (p *Postgres) CheckGoods(store int) ([]goods.Goods, error) {
	txn, err := p.Client.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = txn.Rollback()
	}()
	goodsResult := make([]goods.Goods, 0)
	good := goods.Goods{}
	rows, err := p.Client.Query(`
		SELECT a.uuid, a.name, a.size, (b.amount - b.reserved) as amount
		FROM goods a
		JOIN goods_in_store b ON a.uuid = b.goods_uuid
		JOIN store c ON b.store_id = c.id
		WHERE store_id = $1 AND c.accessibility
	`, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	zeroRows := 0
	for rows.Next() {
		zeroRows++
		err := rows.Scan(&good.Uuid, &good.Name, &good.Size, &good.Amount)
		if err != nil {
			return goodsResult, err
		}
		goodsResult = append(goodsResult, good)
	}
	if zeroRows < 1 {
		return nil, fmt.Errorf("can't get store with this id or store is not avaliable in that moment")
	}
	return goodsResult, err
}
