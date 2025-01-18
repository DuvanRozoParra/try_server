
// Define el tipo de enumeración
type ActionType int
type EventServer int

const (
	FirstHoverEntered ActionType = iota
	LastHoverExited
	HoverEntered
	HoverExited
	FirstSelectEntered
	LastSelectExited
	SelectEntered
	SelectExited
	FirstFocusEntered
	LastFocusExited
	FocusEntered
	FocusExited
	Activated
	Desactivated
)

const (
	RayInteraction EventServer = iota
)

// Define las estructuras
type InteractionEvent struct {
	PlayerID  string     `json:"playerId"`
	EventName string     `json:"eventName"`
	Action    ActionType `json:"action"`
}

type ServerMessage struct {
	EventName EventServer      `json:"eventName"`
	EventData InteractionEvent `json:"eventData"`
}

// teoricamente las variables de la clase
type WebSocketServerUnimeta struct {
	upgrader       websocket.Upgrader         // Para manejar el "handshake" de HTTP a WebSocket
	connections    map[string]*websocket.Conn // Conexiones activas mapeadas por una clave única
	reconnectDelay time.Duration              // Retardo para reconectar en caso de desconexión
	address        string                     // Dirección en la que el servidor escucha
}

// esta es la instancia para llamarla es literalmente el constructor
func NewWebSocketServerUnimeta(address string, reconnectDelay time.Duration) *WebSocketServerUnimeta {
	return &WebSocketServerUnimeta{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // Permitir cualquier origen
		},
		connections:    make(map[string]*websocket.Conn),
		reconnectDelay: reconnectDelay,
		address:        address,
	}
}

// creacion de los metodos
func (s *WebSocketServerUnimeta) handleConnection(conn *websocket.Conn, clientKey string) {
	// Guardar la conexión en el mapa de conexiones activas
	s.connections[clientKey] = conn
	fmt.Println("USUARIO CONECTADO: ", clientKey)
	//fmt.Println("CONEXION: ", s.connections[clientKey])

	// Realizar tareas de lectura y escritura
	go s.listenForMessages(conn, clientKey)

	// Enviar pings a los clientes a intervalos regulares
	// go s.sendPings(clientKey)

}

// Método para escuchar mensajes de un cliente
func (s *WebSocketServerUnimeta) listenForMessages(conn *websocket.Conn, clientKey string) {
	for {
		// Leer mensajes del cliente
		_, message, err := conn.ReadMessage()
		if err != nil {
			// En caso de error (desconexión), eliminar la conexión
			delete(s.connections, clientKey)
			conn.Close()
			return
		}

		// variable for almace the resutl
		var eventMessage ServerMessage

		// Deserializa el JSON
		err = json.Unmarshal([]byte(message), &eventMessage)
		if err != nil {
			fmt.Println("Error al deserializar:", err)
			return
		}

		// Imprimir el mensaje recibido como JSON (string)
		fmt.Printf("EVENT RECEIVED: %d\n", eventMessage.EventName)
		switch eventMessage.EventName {
		case RayInteraction:
			fmt.Println("Ray Interaction Triggered")
			// Serializar el mensaje de nuevo a JSON
			eventMessage.EventData.PlayerID = clientKey
			response, err := json.Marshal(eventMessage)

			if err != nil {
				fmt.Println("Error al serializar el mensaje:", err)
				return
			}

			// Enviar el mismo objeto como respuesta
			err = conn.WriteMessage(websocket.TextMessage, response)
			if err != nil {
				fmt.Println("Error al enviar mensaje:", err)
				return
			}
		default:
			fmt.Println("Unknown Event")
		}
		// Aquí puedes procesar el mensaje recibido
		// Por ejemplo, enviarlo de vuelta al cliente o hacer algo más
		// conn.WriteMessage(websocket.TextMessage, message)
	}
}

// Metodo para escuchar los movimiento de un cliente
func (s *WebSocketServerUnimeta) listenForMove(conn *websocket.Conn, clienteKey string) {

}

// Método para iniciar el servidor WebSocket
func (s *WebSocketServerUnimeta) start() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Realizar el "handshake" y obtener la conexión WebSocket
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		// Asignar una clave única al cliente (puede ser un UUID o algún identificador)
		clientKey := generateClientKey()

		// Manejar la conexión y escuchar mensajes
		s.handleConnection(conn, clientKey)
	})

	// Iniciar el servidor en la dirección especificada
	log.Fatal(http.ListenAndServe(s.address, nil))
}

// Método para generar una clave única para cada cliente (puedes usar un UUID aquí)
func generateClientKey() string {
	// Esto es un ejemplo, puedes usar algo más robusto como un UUID
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
