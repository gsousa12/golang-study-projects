package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"money-transfer-system/golang-simple_microservice-grpc/proto"
	"net"
	"os"

	"google.golang.org/grpc"
)

type conversionServer struct {
	proto.UnimplementedConversionServiceServer
}


func (s *conversionServer) ConvertAmount(ctx context.Context, req *proto.TransferRequest) (*proto.TransferResponse, error) {
	// Ler o banco de dados de conversões (CoinDB.json)
	conversionsData, err := os.ReadFile("../data/CoinDB.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read CoinDB.json: %v", err)
	}

	// Estrutura para deserializar o JSON
	type ConversionDB struct {
		Conversions map[string]struct {
			Code    string `json:"code"`
			Codein  string `json:"codein"`
			Name    string `json:"name"`
			Value   int64  `json:"value"`
		} `json:"conversions"`
	}

	var db ConversionDB
	if err := json.Unmarshal(conversionsData, &db); err != nil {
		return nil, fmt.Errorf("failed to parse CoinDB.json: %v", err)
	}

	// Determinar a chave de conversão com base nos países do remetente e destinatário
	senderCountry := req.SenderAccount.Country
	receiverCountry := req.ReceiverAccount.Country
	var conversionKey string

	if senderCountry == "BR" && receiverCountry == "USA" {
		conversionKey = "BRLUSD"
	} else if senderCountry == "USA" && receiverCountry == "BR" {
		conversionKey = "USDBRL"
	} else {
		return &proto.TransferResponse{
			Status: "FAILURE",
			Reason: "unsupported currency conversion",
		}, nil
	}

	// Obter a taxa de conversão
	conversion, exists := db.Conversions[conversionKey]
	if !exists {
		return &proto.TransferResponse{
			Status: "FAILURE",
			Reason: "conversion rate not found",
		}, nil
	}

	// Calcular o valor convertido
	convertedAmount := (req.TransferAmount * conversion.Value) / 100

	// Retornar resposta de sucesso com o valor convertido
	return &proto.TransferResponse{
		Status:          "SUCCESS",
		Reason:          "",
		ConvertedAmount: convertedAmount,
	}, nil
}

func main() {
	// Configurar o servidor gRPC
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterConversionServiceServer(grpcServer, &conversionServer{})

	log.Println("Conversion Service running on port 50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}