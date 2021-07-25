package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var note Creature

// Punto de ejecución del ejecutable.
func main() {
	// Instancia de http.DefaultServeMux
	mux := http.NewServeMux()

	// Rutas a manejar
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/creatures", CreaturesHandler)

	// Log informativo
	log.Println("Corriendo en http://localhost:8085")

	// Servidor escuchando en el puerto 8080
	http.ListenAndServe(":8085", mux)
}

// IndexHandler nos permite manejar la petición a la ruta '/' y retornar "hola mundo"
// como respuesta al cliente.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hola mundo")
}

// GetNotesHandler nos permite manejar las peticiones a la ruta '/notes' con el método GET.
func GetCreatures(w http.ResponseWriter, r *http.Request) {
	// Puntero a una estructura de tipo Note
	n := NewCreatureApi()
	// Solicitando todas las notas en la base de datos
	creatures, err := n.GetAll()
	if err != nil {
		println("error en handler")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Conviertiendo el slice de Note a formato JSON, retorna un []byte
	j, err := json.Marshal(creatures)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Escribiendo el código de respuesta.
	w.WriteHeader(http.StatusOK)
	// Estableciendo el tipo de contenido del cuerpo de la respuesta.
	w.Header().Set("Content-Type", "application/json")
	// Escribiendo la respuesta, es decir nuestro slice de notas en formato JSON.
	w.Write(j)
}

/*
// CreateNotesHandler nos permite manejar las peticiones a la ruta '/notes' con el método POST.
func CreateNotesHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	// Tomando el cuerpo de la petición, en formato JSON, y decodificandola en
	// la variable note que acabamos de declarar.
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Creamos la nueva nota gracias al método Create.
	err = note.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateNotesHandler nos permite manejar las peticiones a la ruta '/notes' con el método UPDATE.
func UpdateNotesHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Actualizamos la nota correspondiente.
	err = note.Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteNotesHandler nos permite manejar las peticiones a la ruta '/notes' con el método DELETE.
func DeleteNotesHandler(w http.ResponseWriter, r *http.Request) {
	// obtendemosel valor pasado en la url como query correspondiente a id, del tipo ?id=3.
	idStr := r.URL.Query().Get("id")
	// Verificamos que no esté vacío.
	if idStr == "" {
		http.Error(w, "Query id es requerido", http.StatusBadRequest)
		return
	}
	// Convertimos el valor obtenido del query a un int, de ser posible.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Query id debe ser un número", http.StatusBadRequest)
		return
	}

	var note Note
	// Borramos la nota con el id correspondiente.
	err = note.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
*/

func CreaturesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetCreatures(w, r)
	default:
		http.Error(w, "Metodo no permitido", http.StatusBadRequest)
		return
	}
}
