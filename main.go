package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	/*driver para conectar a la base de datos*/
	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contraseña := ""
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contraseña+"@tcp(127.0.0.1)/"+Nombre)
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
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("Servidor Corriendo ...")

	http.ListenAndServe(":8080", nil)

}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	// fmt.Println(idEmpleado)

	conexionEstablecida := conexionBD()

	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?;")

	if err != nil {
		panic(err.Error())
	}
	borrarRegistro.Exec(idEmpleado)

	http.Redirect(w, r, "/", 301)

}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()

	Registro, err := conexionEstablecida.Query("SELECT * FROM empleados;")

	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for Registro.Next() {
		var id int
		var nombre string
		var correo string

		err = Registro.Scan(&id, &nombre, &correo)

		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	//fmt.Println(arregloEmpleado)

	//fmt.Fprintf(w, "Hola mundo")

	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	// fmt.Println(idEmpleado)

	conexionEstablecida := conexionBD()

	Registros, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?;", idEmpleado)

	empleado := Empleado{}
	for Registros.Next() {
		var id int
		var nombre string
		var correo string

		err = Registros.Scan(&id, &nombre, &correo)

		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

	}
	//fmt.Println(empleado)
	plantillas.ExecuteTemplate(w, "editar", empleado)

}

func Crear(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hola mundo")
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()

		insertarRegistro, err := conexionEstablecida.Prepare("INSERT INTO `sistema`.`empleados` (`Nombre`, `Correo`) VALUES (?,?);")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistro.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301)
	}

}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()

		modificarRegistro, err := conexionEstablecida.Prepare("UPDATE empleados SET Nombre=?, Correo=? WHERE id=?;")
		if err != nil {
			panic(err.Error())
		}
		modificarRegistro.Exec(nombre, correo, id)

		http.Redirect(w, r, "/", 301)
	}

}
