package services

import "fmt"

func GetDocentesProyecto(proyectoId int) ([]map[string]interface{}, error) {
	if proyectoId <= 0 {
		return nil, fmt.Errorf("proyectoId inválido")
	}

	resOk, err := GetResolucionesEnEstadoActivo(proyectoId)
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
		})
	}

	return limpios
}
