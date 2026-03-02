package models

type DocenteSolicitudPago struct {
	Persona          string `json:"persona"`
	NumeroContrato   string `json:"numero_contrato"`
	VigenciaContrato int    `json:"vigencia_contrato"`
	Mes              int    `json:"mes"`
	Anio             int    `json:"anio"`
}

type EnviarAprobarSolicitudesCoordinadorRequest struct {
	Coordinador string                 `json:"coordinador"`
	Docentes    []DocenteSolicitudPago `json:"docentes"`
}

type EnviarAprobarSolicitudesCoordinadorItemResult struct {
	Docente DocenteSolicitudPago `json:"docente"`
	PagoId  int                  `json:"pago_id,omitempty"`
	Estado  string               `json:"estado,omitempty"`
	Error   string               `json:"error,omitempty"`
}

type EnviarAprobarSolicitudesCoordinadorResult struct {
	Actualizados []EnviarAprobarSolicitudesCoordinadorItemResult `json:"actualizados"`
	Fallidos     []EnviarAprobarSolicitudesCoordinadorItemResult `json:"fallidos"`
}
