package helpers

import (
	"strconv"
	//"fmt"

	//"github.com/astaxie/beego"
	"github.com/udistrital/cumplidos_dve_mid/models"
)

func GetSolicitudesSupervisor(doc_supervisor string) (pagos_personas_proyecto []models.PagoPersonaProyecto, outputError map[string]interface{}){
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "GetSolicitudesSupervisor", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var parametro []models.Parametro
	var pagos_mensuales []models.PagoMensual
	var contratistas []models.InformacionProveedor
	var pago_personas_proyecto models.PagoPersonaProyecto
	var vinculaciones_docente []models.VinculacionDocente
	
	if err:=GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion:PRS_DVE", &parametro); err == nil{
		if err := GetRequestNew("CumplidosDveUrlCrud", "pago_mensual/?limit=-1&query=Responsable:" + doc_supervisor + ",EstadoPagoMensualId:" + strconv.Itoa(parametro[0].Id), &pagos_mensuales); err == nil{
			for x, pago_mensual := range pagos_mensuales{
				if err := GetRequestLegacy("CumplidosDveUrlCrudAgora", "informacion_proveedor/?query=NumDocumento:" + pago_mensual.Persona, &contratistas); err == nil{
					for _, contratista := range contratistas{
						if err := GetRequestNew("CumplidosDveUrlCrudResoluciones", "vinculacion_docente/?limit=-1&query=NumeroContrato:" + pago_mensual.NumeroContrato + ",Vigencia:" + strconv.FormatFloat(pago_mensual.VigenciaContrato, 'f', 0, 64), &vinculaciones_docente); err == nil{
							for _, vinculacion := range vinculaciones_docente{
								var dep models.Dependencia
								if err := GetRequestLegacy("CumplidosDveUrlCrudOikos", "dependencia/" + strconv.Itoa(vinculacion.ProyectoCurricularId), &dep); err == nil{
									pago_personas_proyecto.PagoMensual = &pagos_mensuales[x]
									pago_personas_proyecto.NombrePersona = contratista.NomProveedor
									pago_personas_proyecto.Dependencia = &dep
									pagos_personas_proyecto = append(pagos_personas_proyecto, pago_personas_proyecto)
								}else{
									panic(err.Error())
								}
							}
						}else{
							panic(err.Error())
						}
					}
				}else{
					panic(err.Error())
				}
			}
		}else{
			panic(err.Error())
		}
	}else{
		panic(err.Error())
	}

	return pagos_personas_proyecto, outputError
}

func GetSolicitudesCoordinador(doc_coordinador string) (pagos_personas_proyecto []models.PagoPersonaProyecto, outputError map[string]interface{}){
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "GetSolicitudesCoordinador", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var parametro []models.Parametro
	var pagos_mensuales []models.PagoMensual
	var contratistas []models.InformacionProveedor
	var pago_personas_proyecto models.PagoPersonaProyecto
	var vinculaciones_docente []models.VinculacionDocente
	
	if err:= GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion:PRC_DVE", &parametro); err == nil{
		if err := GetRequestNew("CumplidosDveUrlCrud", "pago_mensual/?limit=-1&query=Responsable:" + doc_coordinador +",EstadoPagoMensualId:" + strconv.Itoa(parametro[0].Id), &pagos_mensuales); err == nil{
			for x, _ := range pagos_mensuales{
				if err := GetRequestLegacy("CumplidosDveUrlCrudAgora", "informacion_proveedor/?query=NumDocumento:" + pagos_mensuales[x].Persona, &contratistas); err == nil{
					for _, contratista := range contratistas{
						if err := GetRequestNew("CumplidosDveUrlCrudResoluciones", "vinculacion_docente/?limit=-1&query=NumeroContrato:" + pagos_mensuales[x].NumeroContrato + ",Vigencia:" + strconv.FormatFloat(pagos_mensuales[x].VigenciaContrato, 'f', 0, 64), &vinculaciones_docente);err == nil{
							for y, _ := range vinculaciones_docente{
								var dep []models.Dependencia
								if err := GetRequestLegacy("CumplidosDveUrlCrudOikos", "dependencia/?query=Id:" + strconv.Itoa(vinculaciones_docente[y].ProyectoCurricularId), &dep); err == nil{
									for z, _ := range dep{
										pago_personas_proyecto.PagoMensual = &pagos_mensuales[x]
										pago_personas_proyecto.NombrePersona = contratista.NomProveedor
										pago_personas_proyecto.Dependencia = &dep[z]
										pagos_personas_proyecto = append(pagos_personas_proyecto, pago_personas_proyecto)
									}
								}else{
									panic(err.Error())
								}
							}
						}else{
							panic(err.Error())
						}
					}
				}else{
					panic(err.Error())
				}
			}
		}else{
			panic(err.Error())
		}
	}else{
		panic(err.Error())
	}
	return pagos_personas_proyecto, outputError
}

func CertificacionVistoBueno(dependencia string, mes string, anio string) (personas []models.Persona, outputError map[string]interface{}){
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "CertificacionVistoBueno", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var parametros []models.Parametro
	var vinculaciones_docente []models.VinculacionDocente
	var pagos_mensuales []models.PagoMensual
	var contratistas []models.InformacionProveedor
	var persona models.Persona
	var actasInicio []models.ActaInicio
	var mes_cer, _ = strconv.Atoi(mes)
	var anio_cer, _ = strconv.Atoi(anio)

	if err := GetRequestNew("CumplidosDveUrlCrudResoluciones", "vinculacion_docente/?limit=-1&query=ProyectoCurricularId:" + dependencia, &vinculaciones_docente); err == nil{
		for _, vinculacion_docente := range vinculaciones_docente {
			if vinculacion_docente.Activo == true {
				if err := GetRequestLegacy("CumplidosDveUrlCrudAgora", "acta_inicio/?query=NumeroContrato:" + vinculacion_docente.NumeroContrato + ",Vigencia:" + strconv.FormatInt(vinculacion_docente.Vigencia, 10), &actasInicio); err == nil{
					for _, actaInicio := range actasInicio{
						if (actaInicio.FechaInicio.Year() == actaInicio.FechaFin.Year() && int(actaInicio.FechaInicio.Month()) <= mes_cer && actaInicio.FechaInicio.Year() <= anio_cer && int(actaInicio.FechaFin.Month()) >= mes_cer && actaInicio.FechaFin.Year() >= anio_cer) ||
							(actaInicio.FechaInicio.Year() < actaInicio.FechaFin.Year() && int(actaInicio.FechaInicio.Month()) <= mes_cer && actaInicio.FechaInicio.Year() <= anio_cer && int(actaInicio.FechaFin.Month()) <= mes_cer && actaInicio.FechaFin.Year() > anio_cer) ||
							(actaInicio.FechaInicio.Year() < actaInicio.FechaFin.Year() && int(actaInicio.FechaInicio.Month()) >= mes_cer && actaInicio.FechaInicio.Year() < anio_cer && int(actaInicio.FechaFin.Month()) >= mes_cer && actaInicio.FechaFin.Year() >= anio_cer){
								if err := GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion.in:PAD_DVE|AD_DVE|AP_DVE", &parametros); err == nil{
									for _, parametro := range parametros{
										if err := GetRequestNew("CumplidosDveUrlCrud", "pago_mensual/?query=EstadoPagoMensualId.in:" + strconv.Itoa(parametro.Id) + ",NumeroContrato:" + vinculacion_docente.NumeroContrato + ",VigenciaContrato:" + strconv.FormatInt(vinculacion_docente.Vigencia, 10) + ",Mes:" + mes + ",Ano:" + anio,  &pagos_mensuales); err == nil{
											if pagos_mensuales == nil {
												if err := GetRequestLegacy("CumplidosDveUrlCrudAgora", "informacion_proveedor/?query=NumDocumento:" + strconv.Itoa(vinculacion_docente.PersonaId), &contratistas); err == nil{
													for  _, contratista := range contratistas{
														persona.NumDocumento = contratista.NumDocumento
														persona.Nombre = contratista.NomProveedor
														persona.NumeroContrato = actaInicio.NumeroContrato
														persona.Vigencia = actaInicio.Vigencia
														personas = append(personas, persona)
													}
												}else{
													panic(err.Error())
												}
											}
										}else{
											panic(err.Error())
										}
									}
								}
						}
					}
				}else{
					panic(err.Error())
				}
			}
		}
	}else{
		panic(err.Error())
	}
	return personas, outputError
}

func AprobarMultiplesSolicitudes(v []models.PagoMensual) (resultado string, outputError map[string]interface{}){
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "AprobarMultiplesSolicitudes", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	
	var response interface{}
	for _, pm := range v{
		if err := SendRequestNew("CumplidosDveUrlCrud", "pago_mensual/" + strconv.Itoa(pm.Id), "PUT",&response, &pm); err == nil{
			resultado = "OK"
		}else{
			panic(err.Error())
		}
	}
	return resultado, outputError
}