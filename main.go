package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type empleado struct {
	Id      int
	Nombre  string
	Cargo   string
	Salario int
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "./base.db")
	db.Exec("create table if not exists empleado(id INTEGER PRIMARY KEY AutoIncrement,nombre varchar(30) NOT NULL,cargo varchar(30) NOT NULL,salario INT NOT NULL)")

	if err != nil {
		panic(err.Error())
	}
	log.Println("Base de datos conectada")
	return db
}

//Plantillas
var tmpl = template.Must(template.ParseGlob("vista/*"))

//CRUD
//Página para mostrar lista de registros
func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM empleado ORDER BY id DESC")

	if err != nil {
		panic(err.Error())
	}
	emp := empleado{}
	res := []empleado{}

	for selDB.Next() {
		var id, salario int
		var nombre, cargo string

		err = selDB.Scan(&id, &nombre, &cargo, &salario)
		if err != nil {
			panic(err.Error())
		}

		emp.Id = id
		emp.Nombre = nombre
		emp.Cargo = cargo
		emp.Salario = salario

		res = append(res, emp)
	}

	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

//Página para mostrar registros de forma individual
func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM empleado WHERE id = ?", nId)

	if err != nil {
		panic(err.Error())
	}

	emp := empleado{}

	for selDB.Next() {
		var id, salario int
		var nombre, cargo string

		err = selDB.Scan(&id, &nombre, &cargo, &salario)
		if err != nil {
			panic(err.Error())
		}

		emp.Id = id
		emp.Nombre = nombre
		emp.Cargo = cargo
		emp.Salario = salario
	}

	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

//Función para crear nuevos registros
func New(w http.ResponseWriter, r *http.Request)  {
	tmpl.ExecuteTemplate(w, "New", nil)
}

//Función para editar los registros
func Edit(w http.ResponseWriter, r *http.Request)  {
	db := dbConn()
	nId:= r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM empleado WHERE id = ?", nId)

	if err != nil {
		panic(err.Error())
	}
	emp:=empleado{}
	for selDB.Next() {
		var id, salario int
		var nombre, cargo string

		err = selDB.Scan(&id, &nombre, &cargo, &salario)
		if err != nil {
			panic(err.Error())
		}

		emp.Id = id
		emp.Nombre = nombre
		emp.Cargo = cargo
		emp.Salario = salario
	}

	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

// Funciones para insertar registros nuevos creados
func Insert(w http.ResponseWriter, r *http.Request)  {
	db := dbConn()

	if r.Method=="POST" {
		nombre:=r.FormValue("nombre")
		cargo:=r.FormValue("cargo")
		salario:=r.FormValue("salario")

		insForm, err := db.Prepare("INSERT INTO empleado(nombre, cargo, salario) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(nombre, cargo, salario)
		log.Println("Nuevo Registro: "+ nombre+", "+cargo+", "+salario)
	}
	
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

//Función para actualizar registros
func Update(w http.ResponseWriter, r *http.Request)  {
	db := dbConn()

	if r.Method == "POST"{
		nombre := r.FormValue("nombre")
		cargo := r.FormValue("cargo")
		salario := r.FormValue("salario")
		id := r.FormValue("uid")

		insForm, err := db.Prepare("UPDATE empleado SET nombre=?, cargo=?, salario=? WHERE id = ?")

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(nombre, cargo, salario, id)
		log.Println("Registro actualizado: "+ nombre+", "+ cargo+ ", "+salario)
	}

	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// Funcion eliminar registros
func Delete(w http.ResponseWriter, r *http.Request)  {
	db := dbConn()

	id := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM empleado WHERE id = ?")

	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(id)
	log.Println("Registro eliminado")

	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Servidor corriendo en http://localhost:1981")

	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)

	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)



	http.ListenAndServe(":1981", nil)
}
