package models

type VinculacionDocente struct {
	Id								int
	NumeroContrato					string
	Vigencia						int64
	PersonaId						int
	NumeroHorasSemanales			float64
	NumeroSemanas 					float64
	PuntoSalarialId					int
	SalarioMinimoId					int
	ResolucionVinculacionDocenteId	*ResolucionVinculacionDocente
	DedicacionId					int64
	ProyectoCurricularId			int
	ValorContrato					float64
	Categoria						string
	DependenciaAcademica			float64
	NumeroRp						float64
	VigenciaRp						float64
	FechaInicio						string
	Activo							bool
	FechaCreacion					string
	FechaModificacion				string
	
}
