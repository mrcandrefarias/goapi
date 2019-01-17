package main

import (
	"testing"
	"time"
	"net/http"
	"net/http/httptest"
)

func TestGetPool(test *testing.T) {
	if getPool() == nil {
	    test.Error("Erro ao obter pool de conexao")
	}	
}

func TestIsClienteImpactado(test *testing.T) {
	currentTime := time.Now().UTC()
	dt := currentTime.Format("2006-01-02T15:04:05.999999+00:00") 
	
	if !IsClienteImpactado(dt) {
		test.Error("Erro ao obter id cliente")
	}
	
	if IsClienteImpactado(dt) {
		test.Error("Erro cliente impactado anteriormente")
	}
}

func TestDelete(test *testing.T) {
    currentTime := time.Now().UTC()
	dt := currentTime.Format("2006-01-02T15:04:05.999999+00:00")	
	
	if !DelValue(dt) {
		test.Error("Erro ao deletar id cliente")
	} 
}

func TestIsClienteImpactadoEndpoint(t *testing.T) {
	req, _ := http.NewRequest("GET", "/123", nil)
	res := httptest.NewRecorder()

	IsClienteImpactadoEndpoint(res, req)
    isClienteImpactado := res.Body.String()
	if len(isClienteImpactado) < 1 {
		t.Error("Falha no servico IsClienteImpactadoEndpoint")
	}
}

func TestDeleteClienteImpactadoEndpoint(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/123", nil)
	res := httptest.NewRecorder()

	DeleteClienteImpactadoEndpoint(res, req)
    delete := res.Body.String()
	if len(delete) < 1 {
		t.Error("Falha no servico DeleteClienteImpactadoEndpoint")
	}
}
