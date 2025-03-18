package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type Request struct {
	Text string `json:"text"`
}

type Response struct {
	Result string `json:"result"`
}

func ollamaHandler(w http.ResponseWriter, r *http.Request) {
	// Limita apenas a métodos POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodifica o corpo da requisição JSON
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao processar JSON", http.StatusBadRequest)
		return
	}

	// Executa o comando do Ollama via CLI
	cmd := exec.Command("ollama", "run", "llama3.2", "", req.Text)

	// Captura a saída do comando
	output, err := cmd.CombinedOutput()
    fmt.Println(string(output))
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao executar comando: %v", err), http.StatusInternalServerError)
		return
	}

	// Define o cabeçalho de resposta como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

    encoded := base64.StdEncoding.EncodeToString(output)
    decoded, err := base64.StdEncoding.DecodeString(encoded)
	// Estrutura de resposta
	resp := map[string]string{
		"result": string(decoded),
	}
    
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Erro ao gerar resposta JSON", http.StatusInternalServerError)
	}

}

func main() {
	// Define o endpoint de API para o Ollama
	http.HandleFunc("/ollama", ollamaHandler)

	// Inicia o servidor HTTP
	port := ":8080"
	fmt.Printf("Servidor ouvindo na porta %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
