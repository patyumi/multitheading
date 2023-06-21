package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	for _, cep := range os.Args[1:] {
		url := fmt.Sprintf("http://viacep.com.br/ws/" + cep + "/json/")
		url2 := fmt.Sprintf("https://cdn.apicep.com/file/apicep/" + cep + ".json")

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(1)

		go BuscaCep(url, &waitGroup)
		go BuscaCep(url2, &waitGroup)

		waitGroup.Wait()
	}
}

func BuscaCep(url string, wg *sync.WaitGroup) {
	fmt.Println("Buscando CEP na URL " + url)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	select {
	case <-ctxTimeout.Done():
		log.Printf("Requisição para [%s] cancelada :: Timeout", url)
	default:
	}

	req, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição. %v", err)
		return
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao parsear resposta. %v", err)
		return
	}

	fmt.Println(string(res))
	wg.Done()
}
