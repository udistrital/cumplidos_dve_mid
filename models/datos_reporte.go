package models

type DatosReporte struct {
	Id                 int
	ProyectoCurricular string
	Documento          string
	NombrePersona      string
	NumeroContrato     string
	Mes                float64
	Ano                float64
	Estado             string
}
