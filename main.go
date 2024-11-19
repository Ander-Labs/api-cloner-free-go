package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

// Renderiza la página principal
func Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Error cargando la plantilla", http.StatusInternalServerError)
		fmt.Println("Error cargando la plantilla:", err)
		return
	}
	tmpl.Execute(w, nil)
}

// Maneja la clonación de APIs
func CloneAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	apiURL := r.FormValue("api_url")
	if apiURL == "" {
		http.Error(w, "La URL de la API es requerida", http.StatusBadRequest)
		return
	}

	// Realiza la solicitud a la API
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Error al conectar con la API", http.StatusInternalServerError)
		fmt.Println("Error en la solicitud:", err)
		return
	}
	defer resp.Body.Close()

	// Crea un archivo local para almacenar los datos
	file, err := os.Create("data/response.json")
	if err != nil {
		http.Error(w, "Error al guardar los datos", http.StatusInternalServerError)
		fmt.Println("Error creando archivo:", err)
		return
	}
	defer file.Close()

	// Guarda la respuesta en el archivo
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		http.Error(w, "Error al guardar los datos", http.StatusInternalServerError)
		fmt.Println("Error guardando datos:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Datos clonados exitosamente."))
	fmt.Println("Datos clonados exitosamente.")
}

func main() {
	// Crea la carpeta de datos si no existe
	os.MkdirAll("data", os.ModePerm)

	http.HandleFunc("/", Index)
	http.HandleFunc("/clone", CloneAPI)

	fmt.Println("Servidor corriendo en el puerto http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
