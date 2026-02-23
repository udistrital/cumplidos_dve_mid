package models

type DocenteContratoOut struct {
	NumeroContrato string `json:"numero_contrato"`
	Vigencia       int    `json:"vigencia"`
	ProyectoId     int    `json:"proyecto_id"`
}

type DocenteCargoOut struct {
	Documento int                  `json:"documento"`
	Contratos []DocenteContratoOut `json:"contratos"`
}
