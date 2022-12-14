package main

import (
	"database/sql"
	"net/http"
	"text/template"

	//"log"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// conexion base d dtos
func dbconexion() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "dbtienda"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	fmt.Println("servidor corriendo...")
	http.ListenAndServe(":8080", nil)
}

type Camiseta struct{
	Id int
	Nombre string
	Marca string
	Precio int
	Foto string
}
func Inicio(w http.ResponseWriter, r *http.Request) {
	 ConexionEstablecida := dbconexion()
	 consultarRegistros, err := ConexionEstablecida.Query("SELECT * FROM camisetas")

	 if err != nil {
	 	panic(err.Error())
	 }
	 camiseta:=Camiseta{}
	 arregloCamiseta:=[]Camiseta{}

	 for consultarRegistros.Next(){
		var id,precio int
		var nombre,marca,foto string
		err=consultarRegistros.Scan(&id,&nombre, &marca, &precio, &foto)
		if err != nil {
			panic(err.Error())
		}
		camiseta.Id=id
		camiseta.Nombre=nombre
		camiseta.Marca=marca
		camiseta.Precio=precio
		camiseta.Foto=foto
		arregloCamiseta=append(arregloCamiseta, camiseta)
	 }
	 fmt.Println(arregloCamiseta)
	plantillas.ExecuteTemplate(w, "inicio", arregloCamiseta)
}
func Crear(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "HOlaaaaa")
	plantillas.ExecuteTemplate(w, "crear", nil)
}
func Insertar(w http.ResponseWriter,r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		marca := r.FormValue("marca")
		precio := r.FormValue("precio")
		foto := r.FormValue("foto")

		ConexionEstablecida := dbconexion()
		insertarRegistros, err := ConexionEstablecida.Prepare("INSERT INTO camisetas(nombre,marca,precio,foto) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, marca, precio, foto)

		http.Redirect(w, r, "/",http.StatusMovedPermanently)

	}
}
