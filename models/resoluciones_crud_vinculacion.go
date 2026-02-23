package models

type VinculacionDocenteCrud struct {
	Id                   int     `json:"Id"`
	Activo               bool    `json:"Activo"`
	PersonaId            float64 `json:"PersonaId"`
	NumeroContrato       string  `json:"NumeroContrato"`
	Vigencia             int     `json:"Vigencia"`
	ProyectoCurricularId int     `json:"ProyectoCurricularId"`

	ResolucionId *struct {
		Id int `json:"Id"`
	} `json:"ResolucionId,omitempty"`
}
