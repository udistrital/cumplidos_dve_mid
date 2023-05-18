package models

import (
	"time"
)

type Contrato struct {
	Id					int
	NumeroContrato		string
	Vigencia			int
	NombreCompleto		string
	Documento			string
	PersonaId			int
	TipoNominaId		int
	FechaInicio			time.Time
	FechaFin			time.Time
	ValorContrato		int
	Vacaciones 			int
	DependenciaId		int
	ProyectoId 			int
	Cdp 				int
	Rp 					int
	Unico				bool
	Completo 			bool
	Activo				bool
	FechaCreacion  		string
	FechaModificacion	string
	Desagregado			string
}