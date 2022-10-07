package main

import (
	"database/sql"
	"log"
	"text/template"

	"net/http"
	"text/template"
	_ "github.com/mattn/go-sqlite3"
)

type empleado struct{
	Id int
	Nombre string
	Cargo string
	Salario int
}

func dbConn()(db *sql.DB)  {
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
func Index(w http.ResponseWriter, r *http.Request)  {
	db:=dbConn()
	selDB,err:=db.Query("SELECT * FROM empleado ORDER BY id DESC")

	if err != nil {
		panic(err.Error())
	}
	emp:=empleado{}
	res:=[]empleado{}

	for selDB.Next(){
		var id, salario int
		var nombre, cargo string

		err=selDB.Scan(&id, &nombre, &cargo, &salario)
		if err != nil {
			panic(err.Error())
		}

		emp.Id = id
		emp.Nombre=nombre
		emp.Cargo=cargo
		emp.Salario=salario

		res=append(res, emp)
	}

	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func main()  {
	
}