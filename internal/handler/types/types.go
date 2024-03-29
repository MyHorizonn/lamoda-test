package types

import goods "lamoda-test/internal/storage"

type GoodsResp struct {
	Goods []goods.Goods `json:"goods"`
}

type GoodsReq struct {
	Store int           `json:"store"`
	Goods []goods.Goods `json:"goods"`
}

type CheckGoodsReq struct {
	Store int `json:"store"`
}

type GoodsResult struct {
	Uuid   string `json:"uuid"`
	Status string `json:"status"`
}

type WorkOnGoodsResult struct {
	Result []GoodsResult `json:"result"`
}
