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
		"vinculacion_docente/?limit=-1&query=ProyectoCurricularId:"+strconv.Itoa(proyectoId)+",Activo:true,NumeroContrato__isnull:false,Vigencia:"+strconv.Itoa(vigencia),
		&vinculaciones,
	); e != nil {
		panic(e.Error())
	}

	// Contadores
	contadorContratoValido := 0
	contadorConResolucionVinculacion := 0
	contadorConResolucionActiva := 0
	contadorConPersonaValida := 0
	contadorConPagoMensual := 0
	contadorSinPagoMensual := 0
	contadorFinal := 0

	var paramRSOL []models.Parametro
	if err := helpers.GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion:RSOL", &paramRSOL); err != nil || len(paramRSOL) == 0 {
		panic("No se pudo obtener el parámetro RSOL")
	}
	var RSOL = strconv.Itoa(paramRSOL[0].Id)

	res = make([]map[string]interface{}, 0, len(vinculaciones))

	for _, v := range vinculaciones {

		num, ok := v["NumeroContrato"]
		if !ok || num == nil {
			continue
		}
		if s, ok := num.(string); ok && strings.TrimSpace(s) == "" {
			continue
		}
		contadorContratoValido++

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
		contadorConResolucionVinculacion++

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
		contadorConResolucionActiva++

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
		case string:
			personaStr = strings.TrimSpace(t)
		default:
			continue
		}

		if personaStr == "" {
			continue
		}
		contadorConPersonaValida++

		var pagos []map[string]interface{}
		qPago := "pago_mensual/?limit=1&query=NumeroContrato:" + fmt.Sprint(num) +
			",VigenciaContrato:" + strconv.Itoa(vigencia) +
			",Mes:" + strconv.Itoa(mes) +
			",Ano:" + strconv.Itoa(anio) +
			",Persona:" + personaStr

		if e := helpers.GetRequestNew("CumplidosDveUrlCrud", qPago, &pagos); e != nil {
			v["TienePagoMensual"] = nil
			contadorSinPagoMensual++
			res = append(res, v)
			contadorFinal++
			continue
		}

		tienePago := len(pagos) > 0
		v["TienePagoMensual"] = tienePago

		if tienePago {
			contadorConPagoMensual++
		} else {
			contadorSinPagoMensual++
		}

		res = append(res, v)
		contadorFinal++
	}

	return res, nil
}
