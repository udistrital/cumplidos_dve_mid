package models

type Resolucion struct {
	Id                      int
	NumeroResolucion        string
	FechaExpedicion         string
	Vigencia                int
	DependenciaId           int
	TipoResolucionId        int
	PreambuloResolucion     string
	ConsideracionResolucion string
	NumeroSemanas           int
	Periodo                 int
	Titulo                  string
	DependenciaFirmaId      int
	VigenciaCarga           int
	PeriodoCarga            int
	CuadroResponsabilidades string
	NuxeoUid				string
	Activo                  bool
	FechaCreacion           string
	FechaModificacion		string
}
