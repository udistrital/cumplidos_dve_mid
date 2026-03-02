package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/cumplidos_dve_mid/helpers"
	"github.com/udistrital/cumplidos_dve_mid/models"
)

func GetDocentesProyecto(proyectoId int, vigencia int, mes int, anio int) ([]map[string]interface{}, error) {
	if proyectoId <= 0 {
		return nil, fmt.Errorf("proyectoId inválido")
	}
	if vigencia <= 0 {
		return nil, fmt.Errorf("vigencia inválida")
	}

	idOikos, err := GetProyectoOikosIdDesdeSnies(proyectoId)
	if err != nil {
		return nil, err
	}

	resOk, err := GetResolucionesEnEstadoActivo(idOikos, vigencia, mes, anio)
	if err != nil {
		return nil, err
	}

	limpios := limpiarDocentesCamposSlice(resOk)
	return limpios, nil
}

func GetProyectoOikosIdDesdeSnies(proyectoId int) (int, error) {
	if proyectoId <= 0 {
		return 0, fmt.Errorf("proyectoId inválido")
	}

	var err error
	var temp map[string]interface{}

	url := beego.AppConfig.String("CumplidosDveUrlWso2") +
		beego.AppConfig.String("CumplidosDveHomologacion") + "/" +
		"proyecto_curricular_snies/" + strconv.Itoa(proyectoId)

	if err = helpers.GetJsonWSO2(url, &temp); err == nil && temp != nil {

		jsonHomologacion, errorJSON := json.Marshal(temp)
		if errorJSON == nil {

			var proyecto models.ObjetoProyectoCurricular
			if err = json.Unmarshal(jsonHomologacion, &proyecto); err == nil {

				idStr := strings.TrimSpace(proyecto.Homologacion.IDOikos)
				if idStr == "" {
					return 0, fmt.Errorf("no se encontró id_oikos para proyectoId=%d", proyectoId)
				}

				idOikos, convErr := strconv.Atoi(idStr)
				if convErr != nil {
					return 0, fmt.Errorf("id_oikos inválido (%s)", idStr)
				}

				return idOikos, nil

			} else {
				panic(err.Error())
			}

		} else {
			panic(errorJSON.Error())
		}

	} else {
		if err != nil {
			panic(err.Error())
		}
		return 0, fmt.Errorf("respuesta vacía de homologación")
	}
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
