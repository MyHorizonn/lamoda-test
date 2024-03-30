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
	_, err = txn.Exec(`
	DECLARE goods_cursor CURSOR FOR
	SELECT a.id, (amount - reserved) as left
	FROM goods_in_store a
	JOIN store b ON a.store_id = b.id
	WHERE goods_uuid = $1 AND b.accessibility AND (amount - reserved) > 0
	ORDER BY a.store_id
	FOR UPDATE;`,
		good.Uuid)
	if err != nil {
		return err
	}
	leftGoodsAmount := 0
	err = txn.QueryRow(`
		SELECT SUM(amount - reserved) as left
		FROM goods_in_store a
		JOIN store b ON a.store_id = b.id
		WHERE goods_uuid = $1 AND b.accessibility
	`, good.Uuid).Scan(&leftGoodsAmount)
	if err != nil {
		return err
	}
	if leftGoodsAmount < good.Amount {
		return fmt.Errorf("not enough free goods in stores")
	}
	amountToReserve := good.Amount
	for amountToReserve > 0 {
		var rowId int
		var canReserve int
		if err := txn.QueryRow(
			`FETCH NEXT FROM goods_cursor`,
		).Scan(&rowId, &canReserve); err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return err
		}
		_, err := txn.Exec(`
			UPDATE
			goods_in_store
			SET
			reserved = CASE WHEN (amount - reserved) > $2 THEN reserved + $2 ELSE amount END
			WHERE
			id = $1`,
			rowId, amountToReserve)
		if err != nil {
			return err
		}
		amountToReserve -= canReserve
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
	_, err = txn.Exec(`
	DECLARE goods_cursor CURSOR FOR
	SELECT a.id, reserved
	FROM goods_in_store a
	JOIN store b ON a.store_id = b.id
	WHERE goods_uuid = $1 AND b.accessibility AND reserved > 0
	ORDER BY a.store_id
	FOR UPDATE;`,
		good.Uuid)
	if err != nil {
		return err
	}
	reservedAmount := 0
	err = txn.QueryRow(`
		SELECT SUM(reserved)
		FROM goods_in_store a
		JOIN store b ON a.store_id = b.id
		WHERE goods_uuid = $1 AND b.accessibility
	`, good.Uuid).Scan(&reservedAmount)
	if err != nil {
		return err
	}
	if reservedAmount < good.Amount {
		return fmt.Errorf("not enough reserved goods in stores")
	}
	amountToFree := good.Amount
	for amountToFree > 0 {
		var rowId int
		var canFree int
		if err := txn.QueryRow(
			`FETCH NEXT FROM goods_cursor`,
		).Scan(&rowId, &canFree); err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return err
		}
		_, err := txn.Exec(`
			UPDATE
			goods_in_store
			SET
			reserved = CASE WHEN reserved > $2 THEN reserved - $2 ELSE 0 END
			WHERE
			id = $1`,
			rowId, amountToFree)
		if err != nil {
			return err
		}
		amountToFree -= canFree
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
