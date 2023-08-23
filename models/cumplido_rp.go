package models

type CumplidoRp struct {
	Id                   int
	ContratoId           *Contrato
	PreliquidacionId     *Preliquidacion
	Cumplido             bool
	Preliquidado         bool
	ResponsableIva       bool
	Dependientes         bool
	Pensionado           bool
	InteresesVivienda    int
	MedicinaPrepagadaUvt float64
	PensionVoluntaria    int
	Afc                  int
	Activo               bool
	FechaCreacion        string
	FechaModificacion    string
}
