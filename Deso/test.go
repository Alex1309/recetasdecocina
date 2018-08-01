package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
   "os"
	"database/sql"
 _ "github.com/lib/pq"
    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
) 

func getConnection() *sql.DB{

    dsn:=os.Getenv("DATABASE_URL")
	db,err := sql.Open("postgres", dsn)
	
	if err!= nil{
		log.Fatal(err)
	}
	return db
}

type Receta struct {
	Titulo string
	Ingredientes string 
	Descripcion string 
}
func Index(w http.ResponseWriter, r *http.Request) {
    db := getConnection()
    selDB, err := db.Query("SELECT nombre,ingredientes,descripcion FROM receta ORDER BY nombre")
    if err != nil {
        panic(err.Error())
    }
    rec := Receta{}
    res := []Receta{}
    for selDB.Next() {
        var titulo, ingredientes,descripcion string
        err = selDB.Scan(&titulo, &ingredientes, &descripcion)
        if err != nil {
            panic(err.Error())
        }
        rec.Titulo = titulo
        rec.Ingredientes = ingredientes
        rec.Descripcion = descripcion
        res = append(res, rec)
    }
    js, err := json.Marshal(res)  	
   		 if err != nil {   		
    		log.Fatalln(err)   	
    	}   	
    log.Print(js) 	
    w.Header().Set("Content-Type", "application/json") 	
    w.Write(js)
    defer db.Close()
}
func Insert(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close() 	
    var msj Receta 	
    if err := json.NewDecoder(r.Body).Decode(&msj); err != nil { 	   
    	http.Error(w, err.Error(), http.StatusInternalServerError) 	   
    	return 	
    }
    db := getConnection()
    if r.Method == "POST" {
        titulo := msj.Titulo
        ingredientes := msj.Ingredientes
        descripcion:= msj.Descripcion
        insForm, err := db.Prepare("INSERT INTO receta(nombre,ingredientes, descripcion) VALUES($1,$2,$3)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(titulo,ingredientes,descripcion)
    }
    defer db.Close()
    
}
func Update(w http.ResponseWriter, r *http.Request) {
    db := getConnection()
    defer r.Body.Close() 	
    var msj Receta 	
    if err := json.NewDecoder(r.Body).Decode(&msj); err != nil { 	   
    	http.Error(w, err.Error(), http.StatusInternalServerError) 	   
    	return 	
    }
    if r.Method == "PUT" {
        titulo := msj.Titulo
        ingredientes := msj.Ingredientes
        descripcion := msj.Descripcion
        insForm, err := db.Prepare("UPDATE receta SET ingredientes=$1, descripcion=$2 WHERE nombre=$3")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(ingredientes, descripcion,titulo)
    }
    defer db.Close()
}
func Delete(w http.ResponseWriter, r *http.Request) {

    db := getConnection()
    defer r.Body.Close() 	
    var msj Receta 	
    if err := json.NewDecoder(r.Body).Decode(&msj); err != nil { 	   
    	http.Error(w, err.Error(), http.StatusInternalServerError) 	   
    	return 	
    }
    if r.Method =="DELETE" {
    	titulo := msj.Titulo
    	delForm, err := db.Prepare("DELETE FROM receta WHERE nombre=$1")
    	if err != nil {
        	panic(err.Error())
    	}
    	delForm.Exec(titulo)
    }
        defer db.Close()
}

func handleRequests(){

    allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With","Origin","Content-Type","Accept","Authorization"})
    allowedOrigins := handlers.AllowedOrigins([]string{"*"})
    allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

    r := mux.NewRouter()
	r.HandleFunc("/recetas", Index).Methods("GET")
	r.HandleFunc("/insert", Insert).Methods("POST")
    r.HandleFunc("/update", Update).Methods("PUT")
	r.HandleFunc("/delete", Delete).Methods("DELETE")
   log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(r))); 
}
func  main() {
	fmt.Println("Correcto")
	handleRequests()

}