package controllers

import (
	"net/http"
	"testing"
)

func TestObtenerInfoCoordinador(t *testing.T){
	if response, err := http.Get("http://localhost:9001/v1/informacion_academica/informacion_coordinador/28"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestObtenerInfoCoordinador Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestObtenerInfoCoordinador Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestObtenerInfoCoordinador:", err.Error())
		t.Fail()
	}
}

func TestGetContratosDocentes(t *testing.T){
	if response, err := http.Get("http://localhost:9001/v1/informacion_academica/contratos_docente/51653275"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetContratosDocentes Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestGetContratosDocentes Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestGetContratosDocentes:", err.Error())
		t.Fail()
	}
}