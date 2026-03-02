package services

import "fmt"

func GetDocentesProyecto(proyectoId int, vigencia int, mes int, anio int) ([]map[string]interface{}, error) {
	if proyectoId <= 0 {
		return nil, fmt.Errorf("proyectoId inválido")
	}
	if vigencia <= 0 {
		return nil, fmt.Errorf("vigencia inválida")
	}

	resOk, err := GetResolucionesEnEstadoActivo(proyectoId, vigencia, mes, anio)
	if err != nil {
		return nil, err
	}

	limpios := limpiarDocentesCamposSlice(resOk)
	return limpios, nil
}

func limpiarDocentesCamposSlice(resOk []map[string]interface{}) []map[string]interface{} {
	limpios := make([]map[string]interface{}, 0, len(resOk))

	for _, m := range resOk {
		limpios = append(limpios, map[string]interface{}{
			"PersonaId":                      m["PersonaId"],
			"ProyectoCurricularId":           m["ProyectoCurricularId"],
			"Vigencia":                       m["Vigencia"],
			"NumeroContrato":                 m["NumeroContrato"],
			"ResolucionVinculacionDocenteId": m["ResolucionVinculacionDocenteId"],
			"Id":                             m["Id"],
			"PagoMensual":                    m["TienePagoMensual"],
		})
	}

	return limpios
}
