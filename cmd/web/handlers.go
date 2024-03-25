package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"todo.zaidalghurabi.net/internal/validator"
)

type Validator struct {
	validator.Validator
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// GET handlers
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	data := app.newTemplateData(r)
	userId := app.GetUserId(r)
	if userId != "" {
		boards, err := app.boards.GetAll(userId)
		if err != nil {
			app.serverError(w, err)
			return
		}
		data.Boards = boards
	}
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}
func (app *application) boardView(w http.ResponseWriter, r *http.Request) {
	userId := app.GetUserId(r)
	if userId == "" {
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	boardId := r.PathValue("id")
	if boardId == "" {
		app.notFound(w)
		return
	}
	board, err := app.boards.Get(userId, boardId)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Board = board

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}
func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}
func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

// POST handlers
func (app *application) boardCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := r.PostFormValue("title")
	description := "board description TBD"
	userId := app.GetUserId(r)
	if userId == "" {
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	id, err := app.boards.Insert(title, description, userId)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/board/view/%s", id), http.StatusSeeOther)
}
func (app *application) taskCreate(w http.ResponseWriter, r *http.Request) {
	userId := app.GetUserId(r)
	if userId == "" {
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	boardId := r.PathValue("boardId")
	if boardId == "" {
		app.notFound(w)
		return
	}
	groupId := r.PathValue("groupId")
	if groupId == "" {
		app.notFound(w)
		return
	}
	content := r.PostFormValue("task")
	// - validate the content
	v := Validator{}
	v.Check(v.NotBlank(content), "task", "task content should not be empty")
	if !v.Valid() {
		app.clientError(w, http.StatusBadRequest)
		app.infoLog.Println(v.Errors)
		return
	}
	task, err := app.tasks.Insert(userId, boardId, groupId, content)
	if err != nil {
		app.serverError(w, err)
		return
	}
	task.GroupId = groupId
	task.BoardId = boardId
	app.renderPart(w, http.StatusOK, "couple", "task", task)
}
func (app *application) groupCreate(w http.ResponseWriter, r *http.Request) {
	userId := app.GetUserId(r)
	if userId == "" {
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	boardId := r.PathValue("boardId")
	if boardId == "" {
		app.notFound(w)
		return
	}
	name := r.PostFormValue("name")
	group, err := app.groups.Insert(userId, boardId, name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	group.BoardId = boardId
	app.renderPart(w, http.StatusOK, "couple", "group", group)
}
func (app *application) columnCreate(w http.ResponseWriter, r *http.Request) {
	userId := app.GetUserId(r)
	if userId == "" {
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	boardId := r.PathValue("boardId")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	name := r.PostFormValue("name")
	colType := r.PostFormValue("type")
	columnId, err := app.columns.Insert(userId, boardId, name, colType)
	if err != nil || columnId == "" {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// send the columnId as json as a response
	json.NewEncoder(w).Encode(map[string]string{"columnId": columnId})
}

// DELETE handlers
func (app *application) taskDelete(w http.ResponseWriter, r *http.Request) {
	userId := app.GetUserId(r)
	if userId == "" {
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	taskId := r.PathValue("taskId")
	if taskId == "" {
		app.notFound(w)
		return
	}
	groupId := r.PathValue("groupId")
	if groupId == "" {
		app.notFound(w)
		return
	}
	boardId := r.PathValue("boardId")
	if boardId == "" {
		app.notFound(w)
		return
	}
	err := app.tasks.Delete(userId, boardId, groupId, taskId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// auth handlers

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email string `json:"email"`
	}
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Println(err)
		return
	}
	email := req.Email
	if email == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// verify firebase token in "Authorization" header
	token := r.Header.Get("Authorization")
	if token == "" {
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	ctx := context.Background()
	client, err := app.users.Auth.Auth(ctx)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// verify the token
	t, err := client.VerifyIDToken(ctx, token)
	if err != nil {
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	// create firestore user
	err = app.users.CreateUser(email, t.UID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	// redirect user to homepage
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) login(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		app.clientError(w, http.StatusUnauthorized)
		app.errorLog.Println("no token provided")
		return
	}
	ctx := context.Background()
	client, err := app.users.Auth.Auth(ctx)
	if err != nil {
		app.serverError(w, err)
		return
	}
	_, err = client.VerifyIDToken(ctx, token)
	if err != nil {
		app.clientError(w, http.StatusUnauthorized)
		app.errorLog.Printf("invalid token: %s", err)
		return
	}
	// set a cookie (long-lived) to authenticate the user
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    token,
		Expires:  time.Now().Add(72 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
