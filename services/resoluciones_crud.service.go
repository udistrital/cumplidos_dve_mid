package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udistrital/cumplidos_dve_mid/helpers"
)

func GetResolucionesEnEstadoActivo(proyectoId int) (res []map[string]interface{}, err error) {
	var outputError map[string]interface{}

	defer func() {
		if e := recover(); e != nil {
			outputError = map[string]interface{}{
				"funcion": "GetResolucionesEnEstadoActivo",
				"err":     e,
				"status":  "500",
			}
			fmt.Println("[PANIC]", outputError)
			panic(outputError)
		}
	}()

	var vinculaciones []map[string]interface{}
	if e := helpers.GetRequestNew(
		"CumplidosDveUrlCrudResoluciones",
		"vinculacion_docente/?query=ProyectoCurricularId:"+strconv.Itoa(proyectoId)+",NumeroContrato__isnull:false",
		&vinculaciones,
	); e != nil {
		panic(e.Error())
	}

	res = make([]map[string]interface{}, 0, len(vinculaciones))

	for _, v := range vinculaciones {

		num, ok := v["NumeroContrato"]
		if !ok || num == nil {
			continue
		}
		if s, ok := num.(string); ok && strings.TrimSpace(s) == "" {
			continue
		}

		rawRvd, ok := v["ResolucionVinculacionDocenteId"]
		if !ok || rawRvd == nil {
			continue
		}

		rvdID := 0
		if m, ok := rawRvd.(map[string]interface{}); ok {
			if idVal, ok := m["Id"]; ok && idVal != nil {
				switch t := idVal.(type) {
				case float64:
					rvdID = int(t)
				case int:
					rvdID = t
				case string:
					if ii, e := strconv.Atoi(strings.TrimSpace(t)); e == nil {
						rvdID = ii
					}
				}
			}
		}

		if rvdID <= 0 {
			continue
		}

		var estados []map[string]interface{}
		if e := helpers.GetRequestNew(
			"CumplidosDveUrlCrudResoluciones",
			"resolucion_estado/?limit=1&query=ResolucionId__id:"+strconv.Itoa(rvdID)+",EstadoResolucionId:598", //671
			&estados,
		); e != nil {
			continue
		}

		if len(estados) == 0 {
			continue
		}

		res = append(res, v)
	}

	return res, nil
}

func keysOf(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
