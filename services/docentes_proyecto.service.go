package services

import (
	"fmt"
	"strconv"

	"github.com/udistrital/cumplidos_dve_mid/helpers"
	"github.com/udistrital/cumplidos_dve_mid/models"
)

func EnviarYAprobarSolicitudesCoordinador(req models.EnviarAprobarSolicitudesCoordinadorRequest) (res models.EnviarAprobarSolicitudesCoordinadorResult, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "EnviarYAprobarSolicitudesCoordinador", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var paramPRC []models.Parametro
	if err := helpers.GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion:PRC_DVE", &paramPRC); err != nil || len(paramPRC) == 0 {
		panic("No se pudo obtener el parámetro PRC_DVE")
	}

	var paramPAD []models.Parametro
	if err := helpers.GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion:PAD_DVE", &paramPAD); err != nil || len(paramPAD) == 0 {
		panic("No se pudo obtener el parámetro PAD_DVE")
	}

	for _, d := range req.Docentes {

		var pagos []models.PagoMensual
		query := "pago_mensual/?query=NumeroContrato:" + d.NumeroContrato +
			",VigenciaContrato:" + strconv.Itoa(d.VigenciaContrato) +
			",Mes:" + strconv.Itoa(d.Mes) +
			",Ano:" + strconv.Itoa(d.Anio) +
			",Persona:" + d.Persona

		if err := helpers.GetRequestNew("CumplidosDveUrlCrud", query, &pagos); err != nil {
			res.Fallidos = append(res.Fallidos, models.EnviarAprobarSolicitudesCoordinadorItemResult{
				Docente: d,
				Error:   "Error consultando pago_mensual: " + err.Error(),
			})
			continue
		}

		var pm models.PagoMensual
		if len(pagos) == 0 {
			pm = models.PagoMensual{
				NumeroContrato:      d.NumeroContrato,
				VigenciaContrato:    float64(d.VigenciaContrato),
				Mes:                 float64(d.Mes),
				Ano:                 float64(d.Anio),
				Persona:             d.Persona,
				EstadoPagoMensualId: paramPRC[0].Id,
				Responsable:         req.Coordinador,
				CargoResponsable:    "COORDINADOR",
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
				var pagos2 []models.PagoMensual
				if err := helpers.GetRequestNew("CumplidosDveUrlCrud", query, &pagos2); err == nil && len(pagos2) > 0 {
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
		}

		if pm.Responsable != req.Coordinador {
			res.Fallidos = append(res.Fallidos, models.EnviarAprobarSolicitudesCoordinadorItemResult{
				Docente: d,
				PagoId:  pm.Id,
				Error:   "El pago no está asignado al coordinador enviado (Responsable distinto)",
			})
			continue
		}

		if pm.EstadoPagoMensualId != paramPRC[0].Id {
			res.Fallidos = append(res.Fallidos, models.EnviarAprobarSolicitudesCoordinadorItemResult{
				Docente: d,
				PagoId:  pm.Id,
				Error:   "El pago no está en estado PRC_DVE (bandeja coordinador)",
			})
			continue
		}

		infoOrdenador, errOrd := helpers.ObtenerInfoOrdenador(d.NumeroContrato, strconv.Itoa(d.VigenciaContrato))
		if errOrd != nil {
			res.Fallidos = append(res.Fallidos, models.EnviarAprobarSolicitudesCoordinadorItemResult{
				Docente: d,
				PagoId:  pm.Id,
				Error:   "Error obteniendo ordenador: " + fmt.Sprint(errOrd),
			})
			continue
		}

		pm.EstadoPagoMensualId = paramPAD[0].Id
		pm.Responsable = strconv.Itoa(infoOrdenador.NumeroDocumento)
		pm.CargoResponsable = infoOrdenador.Cargo

		var response interface{}
		if err := helpers.SendRequestNew("CumplidosDveUrlCrud", "pago_mensual/"+strconv.Itoa(pm.Id), "PUT", &response, &pm); err != nil {
			res.Fallidos = append(res.Fallidos, models.EnviarAprobarSolicitudesCoordinadorItemResult{
				Docente: d,
				PagoId:  pm.Id,
				Error:   "Error actualizando pago_mensual (PUT): " + err.Error(),
			})
			continue
		}

		res.Actualizados = append(res.Actualizados, models.EnviarAprobarSolicitudesCoordinadorItemResult{
			Docente: d,
			PagoId:  pm.Id,
			Estado:  "OK",
		})
	}

	return res, outputError
}
