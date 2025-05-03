package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	urlStore    = make(map[string]string)
	mu          sync.Mutex
	lettersRune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func validateKey(key string) []byte {
	key = strings.TrimSpace(key)
	if len(key) != 32 {
		log.Fatalf("A chave deve ter exatamente 32 caracteres (tamanho atual: %d)", len(key))
	}
	return []byte(key)
}

func encryptOriginalUrl(originalUrl string, secret_key string) string {
	key := validateKey(secret_key)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal("Erro ao criar cipher:", err)
	}

	plainText := []byte(originalUrl)
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		log.Fatal("Erro ao gerar IV:", err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return hex.EncodeToString(cipherText)
}

func generateId() string {
	b := make([]rune, 8)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(lettersRune))))
		if err != nil {
			log.Fatal("Erro ao gerar ID aleatório:", err)
		}
		b[i] = lettersRune[num.Int64()]
	}
	return string(b)
}

func shortenUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	originalUrl := r.URL.Query().Get("url")
	if originalUrl == "" {
		http.Error(w, "Parâmetro 'url' é obrigatório", http.StatusBadRequest)
		return
	}

	secret_key := "minhachave32bytes1234567890123456"
	base_url := "http://localhost:8080/"

	encryptedUrl := encryptOriginalUrl(originalUrl, secret_key)
	id := generateId()

	mu.Lock()
	urlStore[id] = encryptedUrl
	mu.Unlock()

	shortUrl := base_url + id

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Sua URL encurtada: %s", shortUrl)
}

func main() {
	port := "8080"

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	http.HandleFunc("/short", shortenUrl)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, "Servidor de URL encurtadas está rodando!")
	})

	log.Printf("Servidor rodando na porta %s\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}