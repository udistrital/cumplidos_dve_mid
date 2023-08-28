package controllers

import (
	"net/http"
	"testing"
	"bytes"
)

func TestSolicitudesSupervisor(t *testing.T){
	if response, err := http.Get("http://localhost:9001/v1/aprobacion_documentos/solicitudes_supervisor/52310001"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error SolicitudesSupervisor Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("SolicitudesSupervisor Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error SolicitudesSupervisor:", err.Error())
		t.Fail()
	}
}

func TestSolicitudesCoordinador(t *testing.T){
	if response, err := http.Get("http://localhost:9001/v1/aprobacion_documentos/solicitudes_coordinador/19257731"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestSolicitudesCoordinador Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestSolicitudesCoordinador Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestSolicitudesCoordinador:", err.Error())
		t.Fail()
	}
}

func TestCertificacionVistoBueno(t *testing.T){
	if response, err := http.Get("http://localhost:9001/v1/aprobacion_documentos/certificacion_visto_bueno/72/2/2023"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestCertificacionVistoBueno Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestCertificacionVistoBueno Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestCertificacionVistoBueno:", err.Error())
		t.Fail()
	}
}

func TestGenerarCertificado(t *testing.T){
	if response, err := http.Get("http://localhost:9001/v1/aprobacion_documentos/generar_certificado/HERRERA%20CUBIDES%20JHON%20FRANCINED/INGENIER%C3%8DA%20DE%20SISTEMAS/72/FACULTAD%20DE%20INGENIERIA/FEBRERO/2023/2023-1"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGenerarCertificado Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestGenerarCertificado Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestGenerarCertificado:", err.Error())
		t.Fail()
	}
}

func TestAprobarSolicitudes(t *testing.T){
	body := []byte(`[{
		"Id": 91,
		"NumeroContrato": "DVE2001",
		"VigenciaContrato": 2023,
		"Mes": 5,
		"Persona": "51653275",
		"EstadoPagoMensualId": 4493,
		"Responsable": "79777053",
		"FechaCreacion": "2023-06-22 19:22:20",
		"FechaModificacion": "2023-06-22 19:22:20",
		"CargoResponsable": "SUPERVISOR",
		"Ano": 2023
	  }]`)
	
	if response, err := http.Post("http://localhost:9001/v1/aprobacion_documentos/aprobar_documentos", "application/json", bytes.NewBuffer(body)); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestAprobarSolicitudes Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		}else{
			t.Log("TestAprobarSolicitudes Finalizado Correctamente (OK)")
		}
	}else{
		t.Error("Error TestAprobarSolicitudes:", err.Error())
		t.Fail()
	}
}