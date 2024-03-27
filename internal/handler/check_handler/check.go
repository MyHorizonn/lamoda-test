package check

import (
	"encoding/json"
	"fmt"
	"lamoda-test/internal/handler/types"
	goods "lamoda-test/internal/storage"
	"net/http"
)

func CheckGoods(w http.ResponseWriter, r *http.Request, db goods.Storage) {
	switch r.Method {
	case "POST":
		goods, err := db.CheckGood()
		if err != nil {
			http.Error(w, fmt.Sprintf("error creating, %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		response := types.GoodsResp{Goods: goods}
		jsonErr := json.NewEncoder(w).Encode(response)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error %v", jsonErr), http.StatusInternalServerError)
		}
	default:
		http.Error(w, fmt.Sprintf("method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}
