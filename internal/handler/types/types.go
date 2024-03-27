package types

import goods "lamoda-test/internal/storage"

type GoodsResp struct {
	Goods []goods.Goods `json:"goods"`
}

type GoodsReq struct {
	Goods []goods.Goods `json:"goods"`
}
