package models

type ResolucionRef struct {
	Id int `json:"Id"`
}

type ResolucionEstadoCrud struct {
	Id                 int         `json:"Id"`
	EstadoResolucionId int         `json:"EstadoResolucionId"`
	Activo             bool        `json:"Activo"`
	ResolucionId       interface{} `json:"ResolucionId"`
}
