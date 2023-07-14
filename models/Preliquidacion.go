package models

type Preliquidacion struct {
	Id						int
	Descripcion				string
	Mes						int
	Ano						int
	EstadoPreliquidacionId	int
	NominaId 				int
	Activo					bool
	FechaCreacion			string
	FechaModificacion		string
}