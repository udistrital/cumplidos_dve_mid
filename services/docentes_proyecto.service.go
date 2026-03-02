package services

import (
	"fmt"
	"strconv"

	"github.com/udistrital/cumplidos_dve_mid/helpers"
	"github.com/udistrital/cumplidos_dve_mid/models"
)

func obtenerInfoOrdenadorSafe(numeroContrato, vigencia string) (info models.InformacionOrdenador, ok bool, errMsg string) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
			errMsg = fmt.Sprintf("%v", r)
		}
	}()

	info, outputErr := helpers.ObtenerInfoOrdenador(numeroContrato, vigencia)
	if outputErr != nil {
		ok = false
		errMsg = fmt.Sprintf("%v", outputErr)
		return
	}

	ok = true
	return
}

func EnviarYAprobarSolicitudesCoordinador(req models.EnviarAprobarSolicitudesCoordinadorRequest) (res models.EnviarAprobarSolicitudesCoordinadorResult, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "EnviarYAprobarSolicitudesCoordinador", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	const cargoOrdenadorFijo = "	// Param PAD_DVEORDENADOR DEL GASTO"

	var paramPAD []models.Parametro
	if err := helpers.GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion:PAD_DVE", &paramPAD); err != nil || len(paramPAD) == 0 {
		panic("No se pudo obtener el parámetro PAD_DVE")
	}

	getPago := func(query string) ([]models.PagoMensual, error) {
		var pagos []models.PagoMensual
		if err := helpers.GetRequestNew("CumplidosDveUrlCrud", query, &pagos); err != nil {
			return nil, err
		}
		return pagos, nil
	}

	for _, d := range req.Docentes {

		query := "pago_mensual/?query=NumeroContrato:" + d.NumeroContrato +
			",VigenciaContrato:" + strconv.Itoa(d.VigenciaContrato) +
			",Mes:" + strconv.Itoa(d.Mes) +
			",Ano:" + strconv.Itoa(d.Anio) +
			",Persona:" + d.Persona
		pagos, err := getPago(query)
		if err != nil {
			res.Fallidos = append(res.Fallidos, models.EnviarAprobarSolicitudesCoordinadorItemResult{
				Docente: d,
				Error:   "Error consultando pago_mensual: " + err.Error(),
			})
			continue
		}
		infoOrdenador, okOrd, errOrdMsg := obtenerInfoOrdenadorSafe(d.NumeroContrato, strconv.Itoa(d.VigenciaContrato))
		if !okOrd {
			res.Fallidos = append(res.Fallidos, models.EnviarAprobarSolicitudesCoordinadorItemResult{
				Docente: d,
				Error:   "Error obteniendo ordenador (WSO2/busservicios): " + errOrdMsg,
			})
			continue
		}
		responsableOrdenador := strconv.Itoa(infoOrdenador.NumeroDocumento)
		var pm models.PagoMensual
		if len(pagos) == 0 {
			pm = models.PagoMensual{
				NumeroContrato:      d.NumeroContrato,
				VigenciaContrato:    float64(d.VigenciaContrato),
				Mes:                 float64(d.Mes),
				Ano:                 float64(d.Anio),
				Persona:             d.Persona,
				EstadoPagoMensualId: paramPAD[0].Id,
				Responsable:         responsableOrdenador,
				CargoResponsable:    cargoOrdenadorFijo,
			}

			var created interface{}
			if err := helpers.SendRequestNew("CumplidosDveUrlCrud", "pago_mensual", "POST", &created, &pm); err != nil {
				res.Fallidos = append(res.Fallidos, models.EnviarAprobarSolicitudesCoordinadorItemResult{
					Docente: d,
					Error:   "No existe pago_mensual y falló la creación (POST): " + err.Error(),
				})
				continue
			}
			if id := helpers.HelpersGetID(created); id > 0 {
				pm.Id = id
			} else {
				pagos2, err2 := getPago(query)
				if err2 == nil && len(pagos2) > 0 {
					pm = pagos2[0]
				} else {
					res.Fallidos = append(res.Fallidos, models.EnviarAprobarSolicitudesCoordinadorItemResult{
						Docente: d,
						Error:   "Se creó el pago_mensual pero no pude obtener el Id",
					})
					continue
				}
			}

		} else {
			pm = pagos[0]
			pm.EstadoPagoMensualId = paramPAD[0].Id
			pm.Responsable = responsableOrdenador
			pm.CargoResponsable = cargoOrdenadorFijo

			var response interface{}
			if err := helpers.SendRequestNew("CumplidosDveUrlCrud", "pago_mensual/"+strconv.Itoa(pm.Id), "PUT", &response, &pm); err != nil {
				res.Fallidos = append(res.Fallidos, models.EnviarAprobarSolicitudesCoordinadorItemResult{
					Docente: d,
					PagoId:  pm.Id,
					Error:   "Error actualizando pago_mensual (PUT): " + err.Error(),
				})
				continue
			}
		}

		res.Actualizados = append(res.Actualizados, models.EnviarAprobarSolicitudesCoordinadorItemResult{
			Docente: d,
			PagoId:  pm.Id,
			Estado:  "OK",
		})
	}

	return res, outputError
}
