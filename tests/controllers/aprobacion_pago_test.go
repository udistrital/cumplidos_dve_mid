package controllers

import (
	"net/http"
	"testing"
	"bytes"
)

func TestCertificacionDocumentosAprobados(t *testing.T){
	if response, err := http.Get("http://localhost:9001/v1/aprobacion_pago/certificacion_documentos_aprobados/14/2/2023"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestCertificacionDocumentosAprobados Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestCertificacionDocumentosAprobados Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestCertificacionDocumentosAprobados:", err.Error())
		t.Fail()
	}
}

func TestPagoAprobado(t *testing.T){
	if response, err := http.Get("http://localhost:9001/v1/aprobacion_pago/pago_aprobado/DVE341/2018/5/2018"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestPagoAprobado Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestPagoAprobado Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestPagoAprobado:", err.Error())
		t.Fail()
	}
}

func TestSolicitudesOrdenador(t *testing.T){
	if response, err := http.Get("http://localhost:9001/v1/aprobacion_pago/solicitudes_ordenador/52310001"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestSolicitudesOrdenador Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestSolicitudesOrdenador Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestSolicitudesOrdenador:", err.Error())
		t.Fail()
	}
}

func TestDependenciaOrdenador(t *testing.T){
	if response, err := http.Get("http://localhost:9001/v1/aprobacion_pago/dependencia_ordenador/52310001"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestDependenciaOrdenador Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestDependenciaOrdenador Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestDependenciaOrdenador:", err.Error())
		t.Fail()
	}
}

func TestInfoOrdenador(t *testing.T){
	if response, err := http.Get("http://localhost:9001/v1/aprobacion_pago/informacion_ordenador/DVE2891/2018"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestInfoOrdenador Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestInfoOrdenador Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestInfoOrdenador:", err.Error())
		t.Fail()
	}
}

func TestGenerarCertificadoOrdenador(t *testing.T){
	if response, err := http.Get("http://localhost:9001/v1/aprobacion_pago/generar_certificado/LUZ%20ESPERANZA%20BOHORQUEZ%20AREVALO/14/FACULTAD%20DE%20INGENIERIA/FEBRERO/2023/2023-1"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGenerarCertificadoOrdenador Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestGenerarCertificadoOrdenador Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestGenerarCertificadoOrdenador:", err.Error())
		t.Fail()
	}
}

func TestEnviarTitan(t *testing.T){
	body := []byte(`{
		"Id": 66,
		"NumeroContrato": "DVE2109",
		"VigenciaContrato": 2023,
		"Mes": 4,
		"Persona": "79362987",
		"EstadoPagoMensualId": 4497,
		"Responsable": "52310001",
		"FechaCreacion": "2023-06-22 19:22:20",
		"FechaModificacion": "2023-06-22 19:22:20",
		"CargoResponsable": "ORDENADOR DEL GASTO",
		"Ano": 2023
	  }`)
	
	if response, err := http.Post("http://localhost:9001/v1/aprobacion_pago/enviar_titan", "application/json", bytes.NewBuffer(body)); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestEnviarTitan Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestEnviarTitan Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestEnviarTitan:", err.Error())
		t.Fail()
	}
}

func TestAprobarPagos(t *testing.T){
	body := []byte(`[{
		"Id": 66,
		"NumeroContrato": "DVE2109",
		"VigenciaContrato": 2023,
		"Mes": 4,
		"Persona": "79362987",
		"EstadoPagoMensualId": 4497,
		"Responsable": "52310001",
		"FechaCreacion": "2023-06-22 19:22:20",
		"FechaModificacion": "2023-06-22 19:22:20",
		"CargoResponsable": "ORDENADOR DEL GASTO",
		"Ano": 2023
	  }]`)
	
	if response, err := http.Post("http://localhost:9001/v1/aprobacion_pago/aprobar_pagos", "application/json", bytes.NewBuffer(body)); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestAprobarPagos Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestAprobarPagos Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestAprobarPagos:", err.Error())
		t.Fail()
	}
}