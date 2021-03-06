package main

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

func indexHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var tmpl = template.Must(
		template.ParseFiles("views/layout.html", "views/index.html", "views/head.html"),
	)

	err := tmpl.ExecuteTemplate(w, "layout", struct {
		Title string
	}{
		Title: "TODO List",
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func getTODOHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rows, err := db.Table("todos").Select("*").Rows()
	if err != nil {
		ErrorMessage(w, http.StatusBadGateway, "Unable to retrieve database")
		return
	}

	var todos []TODO

	todos = make([]TODO, 0)

	for rows.Next() {
		var todo struct {
			ID   uint
			Todo string `gorm:"todo"`
		}

		db.ScanRows(rows, &todo)

		todos = append(todos, TODO{
			Index: todo.ID,
			Item:  todo.Todo,
		})
	}

	data := TODOs{
		todos,
	}

	json, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func addTODOHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var t TODO
	err := decoder.Decode(&t)
	if err != nil {
		ErrorMessage(w, http.StatusUnprocessableEntity, "Failed to parse input")
		return
	}

	db.Table("todos").Create(struct {
		Todo string `gorm:"todo"`
	}{
		Todo: t.Item,
	})

	var rec struct {
		ID   uint
		Todo string `gorm:"todo"`
	}

	db.Table("todos").Last(&rec)

	data := TODO{
		Index: rec.ID,
		Item:  rec.Todo,
	}

	json, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func updateTODOHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var t TODO
	err := decoder.Decode(&t)
	if err != nil {
		ErrorMessage(w, http.StatusUnprocessableEntity, "Failed to parse input")
		return
	}

	db.Table("todos").Where("id == ?", t.Index).Update(struct {
		Todo string `gorm:"todo"`
	}{
		Todo: t.Item,
	})

	data := TODO{
		Index: t.Index,
		Item:  t.Item,
	}

	json, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func deleteTODOHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var t TODO
	err := decoder.Decode(&t)
	if err != nil {
		ErrorMessage(w, http.StatusUnprocessableEntity, "Failed to parse input")
		return
	}

	db.Table("todos").Where("id == ?", t.Index).Delete(struct {
		ID   uint
		Todo string
	}{})

	data := struct {
		Message string `json:"message"`
	}{
		Message: "TODO item deleted",
	}

	json, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
