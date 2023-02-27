package models

type PagoMensual struct {
	Id                int             
	NumeroContrato    string             
	VigenciaContrato  float64            
	Mes               float64            
	Persona           string             
	EstadoPagoMensualId int  
	Responsable       string  
	FechaCreacion	  string
	FechaModificacion string        
	CargoResponsable  string             
	Ano               float64            
}