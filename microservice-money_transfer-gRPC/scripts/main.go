package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

const (
	url         = "http://localhost:8080/transfer"
	numRequests = 500
)

type Account struct {
	ID      string `json:"id"`
	Amount  int64  `json:"amount"`
	Country string `json:"country"`
}

type Transfer struct {
	Amount int64 `json:"amount"`
}

type RequestBody struct {
	SenderAccount   Account  `json:"senderAccount"`
	ReceiverAccount Account  `json:"receiverAccount"`
	Transfer        Transfer `json:"transfer"`
	Status          string   `json:"status"`
	Reason          string   `json:"reason"`
}

func sendRequest(i int, wg *sync.WaitGroup) {
	defer wg.Done()

	body := RequestBody{
		SenderAccount: Account{
			ID:      fmt.Sprintf("sender-%d", i),
			Amount:  1200000,
			Country: "BR",
		},
		ReceiverAccount: Account{
			ID:      fmt.Sprintf("receiver-%d", i),
			Amount:  500000,
			Country: "USA",
		},
		Transfer: Transfer{
			Amount: 100000,
		},
		Status: "PENDING",
		Reason: "",
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Erro ao serializar JSON da request %d: %v\n", i, err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Erro na request %d: %v\n", i, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Request %d retornou status: %s\n", i, resp.Status)
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go sendRequest(i, &wg)
	}

	wg.Wait()
	fmt.Println("Todas as requisições foram enviadas.")
}
