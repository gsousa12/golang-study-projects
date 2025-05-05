package main

import (
	"encoding/json"
	"fmt"
	"log"
	"money-transfer-system/golang-simple_microservice-grpc/proto"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Conectar ao Transaction Service
	transactionConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Transaction Service: %v", err)
	}

	defer transactionConn.Close()
	transactionClient := proto.NewTransactionServiceClient(transactionConn)

	// Conectar ao Conversion Service
	conversionConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Conversion Service: %v", err)
	}
	defer conversionConn.Close()
	conversionClient := proto.NewConversionServiceClient(conversionConn)

	// Configurar o handler HTTP
	http.HandleFunc("/transfer", func(w http.ResponseWriter, r *http.Request) {
		handleTransfer(w, r, transactionClient, conversionClient)
	})

	// Iniciar o servidor HTTP
	log.Println("API Gateway running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}

func handleTransfer(w http.ResponseWriter, r *http.Request, transactionClient proto.TransactionServiceClient, conversionClient proto.ConversionServiceClient) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar a requisição JSON do cliente
	var req proto.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validar saldo
	transactionResp, err := transactionClient.ValidateBalance(r.Context(), &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Transaction Service error: %v", err), http.StatusInternalServerError)
		return
	}

	if transactionResp.Status == "FAILURE" {
		json.NewEncoder(w).Encode(transactionResp)
		return
	}

	// Converter
	conversionResp, err := conversionClient.ConvertAmount(r.Context(), &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Conversion Service error: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornar resposta ao cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversionResp)
}