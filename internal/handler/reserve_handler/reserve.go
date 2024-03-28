package reserve

import (
	"encoding/json"
	"fmt"
	types "lamoda-test/internal/handler/types"
	goods "lamoda-test/internal/storage"
	"log"
	"net/http"
	"sync"
)

func ReserveGoods(w http.ResponseWriter, r *http.Request, db goods.Storage) {
	switch r.Method {
	case "POST":
		var req types.GoodsReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("error reading request body, %v\n", err)
			http.Error(w, fmt.Sprintf("error reading request body, %v", err), http.StatusBadRequest)
			return
		}
		if len(req.Goods) == 0 {
			log.Println("error empty data")
			http.Error(w, "error empty data", http.StatusBadRequest)
			return
		}
		errsArr := make([]types.GoodsResult, len(req.Goods))
		wg := sync.WaitGroup{}
		wg.Add(len(req.Goods))
		for i, good := range req.Goods {
			go func(good goods.Goods, idx int) {
				err := db.ReserveGood(good)
				reserveErr := types.GoodsResult{Uuid: good.Uuid, Status: fmt.Sprintf("reserved %d items", good.Amount)}
				if err != nil {
					reserveErr.Status = err.Error()
				}
				errsArr[idx] = reserveErr
				wg.Done()
			}(good, i)
		}
		wg.Wait()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response := types.WorkOnGoodsResult{Result: errsArr}
		jsonErr := json.NewEncoder(w).Encode(response)
		if jsonErr != nil {
			log.Printf("error while encoding response: %v\n", jsonErr)
			http.Error(w, fmt.Sprintf("error while encoding response: %v", jsonErr), http.StatusInternalServerError)
		}
	default:
		http.Error(w, fmt.Sprintf("method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}
