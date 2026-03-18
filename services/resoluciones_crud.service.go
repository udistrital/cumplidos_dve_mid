package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udistrital/cumplidos_dve_mid/helpers"
	"github.com/udistrital/cumplidos_dve_mid/models"
)

func GetResolucionesEnEstadoActivo(proyectoId int, vigencia int, mes int, anio int) (res []map[string]interface{}, err error) {
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
		"vinculacion_docente/?limit=-1&query=ProyectoCurricularId:"+strconv.Itoa(proyectoId)+",NumeroContrato__isnull:false,Vigencia:"+strconv.Itoa(vigencia),
		&vinculaciones,
	); e != nil {
		panic(e.Error())
	}

	var paramRSOL []models.Parametro
	if e := helpers.GetRequestNew(
		"CumplidosDveUrlParametros",
		"parametro/?query=CodigoAbreviacion:RSOL",
		&paramRSOL,
	); e != nil || len(paramRSOL) == 0 {
		panic("No se pudo obtener el parámetro RSOL")
	}
	RSOL := strconv.Itoa(paramRSOL[0].Id)

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
				case int64:
					rvdID = int(t)
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

		var resolucion map[string]interface{}
		if e := helpers.GetRequestNew(
			"CumplidosDveUrlCrudResoluciones",
			"resolucion/"+strconv.Itoa(rvdID),
			&resolucion,
		); e != nil {
			continue
		}

		activoResolucion := false
		numeroResolucion := ""

		dataResolucion := map[string]interface{}{}
		if data, ok := resolucion["Data"]; ok {
			if data == nil {
				continue
			}
			if dataMap, ok := data.(map[string]interface{}); ok {
				dataResolucion = dataMap
			} else {
				continue
			}
		} else {
			dataResolucion = resolucion
		}

		if activoVal, ok := dataResolucion["Activo"]; ok && activoVal != nil {
			switch t := activoVal.(type) {
			case bool:
				activoResolucion = t
			case string:
				activoResolucion = strings.EqualFold(strings.TrimSpace(t), "true")
			}
		}

		if !activoResolucion {
			continue
		}

		if nroVal, ok := dataResolucion["NumeroResolucion"]; ok && nroVal != nil {
			switch t := nroVal.(type) {
			case string:
				numeroResolucion = strings.TrimSpace(t)
			case float64:
				numeroResolucion = strconv.Itoa(int(t))
			case int:
				numeroResolucion = strconv.Itoa(t)
			case int64:
				numeroResolucion = strconv.Itoa(int(t))
			default:
				numeroResolucion = fmt.Sprint(t)
			}
		}

		var estados []map[string]interface{}
		if e := helpers.GetRequestNew(
			"CumplidosDveUrlCrudResoluciones",
			"resolucion_estado/?limit=1&query=ResolucionId__id:"+strconv.Itoa(rvdID)+",EstadoResolucionId:"+RSOL,
			&estados,
		); e != nil {
			continue
		}

		if len(estados) == 0 {
			continue
		}

		personaVal, ok := v["PersonaId"]
		if !ok || personaVal == nil {
			continue
		}

		personaStr := ""
		switch t := personaVal.(type) {
		case float64:
			personaStr = strconv.Itoa(int(t))
		case int:
			personaStr = strconv.Itoa(t)
		case int64:
			personaStr = strconv.Itoa(int(t))
		case string:
			personaStr = strings.TrimSpace(t)
		default:
			continue
		}

		if personaStr == "" {
			continue
		}

		var pagos []map[string]interface{}
		qPago := "pago_mensual/?limit=1&query=NumeroContrato:" + fmt.Sprint(num) +
			",VigenciaContrato:" + strconv.Itoa(vigencia) +
			",Mes:" + strconv.Itoa(mes) +
			",Ano:" + strconv.Itoa(anio) +
			",Persona:" + personaStr

		nuevo := make(map[string]interface{})
		for k, val := range v {
			nuevo[k] = val
		}

		if e := helpers.GetRequestNew("CumplidosDveUrlCrud", qPago, &pagos); e != nil {
			nuevo["TienePagoMensual"] = nil
			nuevo["NumeroResolucion"] = numeroResolucion
			res = append(res, nuevo)
			continue
		}

		nuevo["TienePagoMensual"] = len(pagos) > 0
		nuevo["NumeroResolucion"] = numeroResolucion
		res = append(res, nuevo)
	}

	return res, nil
}
