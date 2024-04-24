package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/", app.home)

	// auth routes
	mux.HandleFunc("GET /user/signup", app.signupForm)
	mux.HandleFunc("POST /user/signup", app.signup)
	mux.HandleFunc("GET /user/login", app.loginForm)
	mux.HandleFunc("POST /user/login", app.login)

	// protected routes
	protected := alice.New(app.protectRoute)
	mux.Handle("POST /user/logout", protected.ThenFunc(app.logout))
	// -- "board/create" handler should only respond to POST requests
	mux.Handle("POST /board/create", protected.ThenFunc(app.boardCreate))
	mux.Handle("/board/view/{id}", protected.ThenFunc(app.boardView))
	mux.Handle("PUT /board/update/title/{id}", protected.ThenFunc(app.boardTitleUpdate))
	// -- task handlers
	mux.Handle("POST /task/create/{boardId}/{groupId}", protected.ThenFunc(app.taskCreate))
	mux.Handle("DELETE /task/delete/{boardId}/{groupId}/{taskId}", protected.ThenFunc(app.taskDelete))
	mux.Handle("PUT /task/swap/{boardId}/{groupId}/{swappedId}/{swappedOrder}/{targetId}/{targetOrder}", protected.ThenFunc(app.taskSwap))
	mux.Handle("PUT /task/update/colVal/{boardId}/{groupId}/{taskId}/{columnId}", protected.ThenFunc(app.taskColValUpdate))
	mux.Handle("PUT /task/update/title/{boardId}/{groupId}/{taskId}", protected.ThenFunc(app.taskContentUpdate))
	// -- group handlers
	mux.Handle("POST /group/create/{boardId}", protected.ThenFunc(app.groupCreate))
	mux.Handle("PUT /group/update/title/{boardId}/{groupId}", protected.ThenFunc(app.groupNameUpdate))
	// -- column handlers
	mux.Handle("POST /column/create/{boardId}", protected.ThenFunc(app.columnCreate))
	mux.Handle("PUT /column/drag/{boardId}", protected.ThenFunc(app.reorderColumns))
	// -- status handlers
	mux.Handle("POST /status/create/{boardId}/{groupId}/{taskId}", protected.ThenFunc(app.statusCreate))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders, app.requireAuth)

	return standard.Then(mux)
}
