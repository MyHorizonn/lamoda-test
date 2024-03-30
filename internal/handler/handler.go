package handler

import (
	check "lamoda-test/internal/handler/check_handler"
	free "lamoda-test/internal/handler/free_handler"
	reserve "lamoda-test/internal/handler/reserve_handler"
	goods "lamoda-test/internal/storage"
	"log"
	"net/http"
)

func StartServer(db goods.Storage) {
	http.HandleFunc("/reserve_goods", func(w http.ResponseWriter, r *http.Request) {
		reserve.ReserveGoods(w, r, db)
	})

	http.HandleFunc("/free_goods", func(w http.ResponseWriter, r *http.Request) {
		free.FreeGoods(w, r, db)
	})

	http.HandleFunc("/check_goods", func(w http.ResponseWriter, r *http.Request) {
		check.CheckGoods(w, r, db)
	})

	log.Fatalln(http.ListenAndServe(":8000", nil))
}
