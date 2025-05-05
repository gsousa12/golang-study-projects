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

type transactionServer struct {
	proto.UnimplementedTransactionServiceServer
}


func (s *transactionServer) ValidateBalance(ctx context.Context, req *proto.TransferRequest) (*proto.TransferResponse, error) {

	/* Ler banco de dados in memory */
	accountsData, err := os.ReadFile("../data/AccountDB.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read AccountDB.json: %v", err)
	}

	/* Deserializar Json */
	type AccountDB struct {
		accounts []struct {
			ID      string `json:"id"`
			Balance int64  `json:"balance"`
			Country string `json:"country"`
		} `json:"accounts"`
	}

	var db AccountDB
	if err := json.Unmarshal(accountsData, &db); err != nil {
		return nil, fmt.Errorf("failed to parse AccountDB.json: %v", err)
	}

	/* Buscar conta do remetente */
	var senderBalance int64
	found := false

	for _, account := range db.accounts {
		if account.ID == req.SenderAccount.Id && account.Country == req.SenderAccount.Country {
			senderBalance = account.Balance
			found = true
			break
		}
	}

	if !found {
		return &proto.TransferResponse{
			Status: "Failure",
			Reason: "sender account not found",
		}, nil
	}

	/* Validar Saldo */
	if senderBalance < req.TransferAmount {
		return &proto.TransferResponse{
			Status: "FAILURE",
			Reason: "insufficient balance",
		}, nil 
	}

	return &proto.TransferResponse{
		Status: "SUCCESS",
		Reason: "",
	}, nil
}

func main() {
	/* Configuração do servidor gRPC */
	lis,err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterTransactionServiceServer(grpcServer, &transactionServer{})

	log.Println("Transaction Service running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}