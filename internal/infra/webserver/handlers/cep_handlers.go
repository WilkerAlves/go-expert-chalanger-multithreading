package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/WilkerAlves/go-expert-chalanger-client-server-api/internal/dto"
	"github.com/go-chi/chi/v5"
)

type CepHandler struct{}

func (c *CepHandler) GetCep(res http.ResponseWriter, req *http.Request) {
	cep := chi.URLParam(req, "cep")
	if cep == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	chApiCep := make(chan string)
	go getApiCep(req.Context(), cep, chApiCep)

	chViaCep := make(chan string)
	go getViaCep(req.Context(), cep, chViaCep)

	select {
	case msg := <-chViaCep:
		fmt.Println("GetViaCep")
		fmt.Println(msg)
		res.WriteHeader(http.StatusOK)
		return
	case msg := <-chApiCep:
		fmt.Println("GetApiCep")
		fmt.Println(msg)
		res.WriteHeader(http.StatusOK)
		return
	case <-req.Context().Done():
		fmt.Println("timeout")
		res.WriteHeader(http.StatusRequestTimeout)
		return
	}
}

func getApiCep(ctx context.Context, cep string, ch chan<- string) {
	url := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	var c dto.ApiCep
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Erro ao converter para bytes. %v", err)
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Fatalf("Erro ao fazer o Unmarshal. %v", err)
	}

	marshal, err := json.Marshal(&c)
	if err != nil {
		return
	}
	ch <- fmt.Sprint(string(marshal))
	close(ch)
}

func getViaCep(ctx context.Context, cep string, ch chan<- string) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	var c dto.ViaCep
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Erro ao converter para bytes. %v", err)
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Fatalf("Erro ao fazer o Unmarshal. %v", err)
	}

	marshal, err := json.Marshal(&c)
	if err != nil {
		return
	}

	ch <- fmt.Sprint(string(marshal))
	close(ch)
}

func NewCepHandler() *CepHandler {
	return &CepHandler{}
}
