package app

import (
	"sc-app/model"
	"log"
	"github.com/alexandrevicenzi/go-sse"
	"net/http"

	"sc-app/mdware"
	"sc-app/handler"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
)

//App struct
type App struct {
	Router *chi.Mux
	Db     *gorm.DB
	Sse     *sse.Server
}

func (a *App) setRouter() {
	a.Router.Use(mdware.JWTAuth)

	a.Router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		handler.Login(w, r, a.Db)
	})
	a.Router.Post("/register",func(w http.ResponseWriter, r *http.Request) {
		handler.Register(w, r, a.Db)
	})

	a.Router.Route("/post", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			handler.CreatePost(w, r, a.Db, a.Sse)
		})
		r.Get("/{PostID}", func(w http.ResponseWriter, r *http.Request) {
			handler.GetPost(w, r, a.Db, a.Sse)
		})
		r.Put("/{PostID}", func(w http.ResponseWriter, r *http.Request) {
			handler.UpdatePost(w, r, a.Db, a.Sse)
		})
		r.Delete("/{PostID}", func(w http.ResponseWriter, r *http.Request) {
			handler.DeletePost(w, r, a.Db, a.Sse)
		})
	})
	
	a.Router.Mount("/events/", a.Sse)

	//example listen event 
	a.Router.Handle("/", http.FileServer(http.Dir("./static")))
	mdware.AddNoAuthRoute("/")

}

//Initalize router and SSE
func (a *App) Initalize() {
	//migrate model 
	a.Db = model.DbMigrate(a.Db)
	//set Router
	a.setRouter()
}

//Run application
func (a *App) Run() {
	a.Initalize()
	defer a.Sse.Shutdown() 
	log.Println("Server listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", a.Router))
}
