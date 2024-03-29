package check

import (
	"encoding/json"
	"fmt"
	"lamoda-test/internal/handler/types"
	goods "lamoda-test/internal/storage"
	"log"
	"net/http"
)

func CheckGoods(w http.ResponseWriter, r *http.Request, db goods.Storage) {
	switch r.Method {
	case "POST":
		var req types.CheckGoodsReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("error reading request body, %v\n", err)
			http.Error(w, fmt.Sprintf("error reading request body, %v", err), http.StatusBadRequest)
			return
		}
		goods, err := db.CheckGoods(req.Store)
		if err != nil {
			log.Printf("error checking goods: %v\n", err)
			http.Error(w, fmt.Sprintf("error checking goods: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		response := types.GoodsResp{Goods: goods}
		jsonErr := json.NewEncoder(w).Encode(response)
		if jsonErr != nil {
			log.Printf("error while encoding response: %v\n", jsonErr)
			http.Error(w, fmt.Sprintf("error while encoding response: %v", jsonErr), http.StatusInternalServerError)
		}
	default:
		http.Error(w, fmt.Sprintf("method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}
