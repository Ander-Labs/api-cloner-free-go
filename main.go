package main

import (
	// "encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Estructura de solicitud
type Request struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    string
}

// Estructura de respuesta
type Response struct {
	StatusCode int
	Headers    map[string][]string
	Body       string
}

var history []Request

// Función para renderizar la página principal
func Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Error cargando la plantilla", http.StatusInternalServerError)
		fmt.Println("Error cargando la plantilla:", err)
		return
	}

	// Pasar el historial de solicitudes a la plantilla
	tmpl.Execute(w, history)
}

// Función para manejar la solicitud y obtener la respuesta
func SendRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener los datos de la solicitud
	url := r.FormValue("url")
	method := r.FormValue("method")
	headers := make(map[string]string)
	for key, values := range r.Form {
		if key != "url" && key != "method" && key != "body" {
			headers[key] = values[0]
		}
	}
	body := r.FormValue("body")

	// Crear un objeto de solicitud
	request := Request{
		URL:     url,
		Method:  method,
		Headers: headers,
		Body:    body,
	}

	// Guardar la solicitud en el historial
	history = append(history, request)

	// Enviar la solicitud HTTP
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		http.Error(w, "Error creando la solicitud", http.StatusInternalServerError)
		fmt.Println("Error creando solicitud:", err)
		return
	}

	// Agregar encabezados
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Si hay cuerpo, agregarlo
	if body != "" {
		req.Body = ioutil.NopCloser(strings.NewReader(body))
	}

	// Obtener la respuesta
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error al hacer la solicitud", http.StatusInternalServerError)
		fmt.Println("Error en la solicitud:", err)
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error al leer la respuesta", http.StatusInternalServerError)
		fmt.Println("Error leyendo la respuesta:", err)
		return
	}

	// Mostrar la respuesta en la plantilla
	response := Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       string(respBody),
	}

	// Pasar la respuesta a la plantilla
	tmpl, err := template.ParseFiles("template/response.html")
	if err != nil {
		http.Error(w, "Error cargando la plantilla de respuesta", http.StatusInternalServerError)
		fmt.Println("Error cargando la plantilla de respuesta:", err)
		return
	}

	tmpl.Execute(w, response)
}

func main() {
	// Crear la carpeta de datos si no existe
	err := os.MkdirAll("data", os.ModePerm)
	if err != nil {
		fmt.Println("Error creando la carpeta de datos:", err)
		return
	}

	// Manejar las rutas
	http.HandleFunc("/", Index)
	http.HandleFunc("/send", SendRequest)

	// Iniciar el servidor
	fmt.Println("Servidor corriendo en el puerto http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
