package free

import (
	"encoding/json"
	"fmt"
	types "lamoda-test/internal/handler/types"
	goods "lamoda-test/internal/storage"
	"net/http"
	"sync"
)

func FreeGoods(w http.ResponseWriter, r *http.Request, db goods.Storage) {
	switch r.Method {
	case "POST":
		var req types.GoodsReq
		var err error
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body, %v", err), http.StatusBadRequest)
			return
		}
		if len(req.Goods) == 0 {
			http.Error(w, "error empty data", http.StatusBadRequest)
			return
		}
		wg := sync.WaitGroup{}
		wg.Add(len(req.Goods))
		for _, good := range req.Goods {
			go func(good goods.Goods) error {
				err := db.ReserveGoods(good)
				wg.Done()
				return err
			}(good)
		}
		wg.Wait()
		if err != nil {
			http.Error(w, fmt.Sprintf("error creating, %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, fmt.Sprintf("method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}
