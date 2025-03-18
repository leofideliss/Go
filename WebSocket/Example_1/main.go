package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

// Definindo um upgrade de HTTP para WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Função para lidar com as conexões WebSocket
func handleConnection(w http.ResponseWriter, r *http.Request) {
	// Fazendo o upgrade de HTTP para WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao fazer o upgrade: ", err)
		return
	}
	defer conn.Close()

	// Loop de leitura e escrita
	for {
		// Lê a mensagem do cliente
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Erro ao ler a mensagem:", err)
			break
		}

		// Exibe a mensagem recebida
		fmt.Printf("Mensagem recebida: %s\n", p)

		// Envia a mesma mensagem de volta para o cliente
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			log.Println("Erro ao enviar a mensagem:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)
	
	// Inicia o servidor na porta 8080
	fmt.Println("Servidor WebSocket rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
