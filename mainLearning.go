package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Definir el upgrader, que maneja la solicitud de HTTP a WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permite conexiones desde cualquier origen
	},
}

// Handler para manejar la conexión WebSocket
func handler(w http.ResponseWriter, r *http.Request) {
	// Establecer la conexión WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Leer y escribir mensajes
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// Enviar el mensaje de vuelta al cliente
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func mainLearning() {
	http.HandleFunc("/ws", handler)              // Ruta WebSocket
	log.Fatal(http.ListenAndServe(":8080", nil)) // Iniciar servidor
}
