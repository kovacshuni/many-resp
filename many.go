package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sync/atomic"
	"time"
)

// App is the app
type App struct {
	count  *uint64
	router *mux.Router
	ticker *time.Ticker
}

func main() {
	app := newApp()
	go app.report()
	app.router.Path("/").HandlerFunc(app.serveMany)
	err := http.ListenAndServe(":8080", app.router)
	if err != nil {
		panic(err)
	}
}

func newApp() App {
	var c uint64
	return App{
		count:  &c,
		router: mux.NewRouter(),
		ticker: time.NewTicker(time.Second),
	}
}

func (app App) serveMany(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	atomic.AddUint64(app.count, 1)
}

func (app App) report() {
	for range app.ticker.C {
		finalN := atomic.LoadUint64(app.count)
		atomic.StoreUint64(app.count, 0)
		fmt.Printf("nServedRequests=%d\n", finalN)
	}
}
