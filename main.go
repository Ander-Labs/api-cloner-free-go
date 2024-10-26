package main

import (
	"fmt"
	"html/template"
	"net/http"
	// "encoding/json"
	// "log"
)



func Index(rw http.ResponseWriter, r *http.Request){
	template, err:= template.ParseFiles("template/index.html")
	if err != nil{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}else{
	template.Execute(rw, nil)
	fmt.Println(rw, "Hola desde el server")
}
}

func main(){
	http.HandleFunc("/", Index)
	fmt.Println("Servidor corriendo en el puerto http://localhost:8080")
	http.ListenAndServe(":8080", nil)
	
}