package helpers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/cumplidos_dve_mid/models"
)

func CargarInformacionCoordinador(DependenciaOikosId int) (info_coordinador models.InformacionCoordinador, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "CargarInformacionCoordinador", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var err error
	var temp map[string]interface{}
	var temp_snies map[string]interface{}

	if err = GetJsonWSO2(beego.AppConfig.String("CumplidosDveUrlWso2")+beego.AppConfig.String("CumplidosDveHomologacion")+"/"+"proyecto_curricular_oikos/"+strconv.Itoa(DependenciaOikosId), &temp); err == nil && temp != nil {
		json_proyecto_curricular, error_json := json.Marshal(temp)
		if error_json == nil {
			var temp_homologacion models.ObjetoProyectoCurricular
			if err = json.Unmarshal(json_proyecto_curricular, &temp_homologacion); err == nil {
				id_proyecto_snies := temp_homologacion.Homologacion.IDSnies
				if err = GetJsonWSO2(beego.AppConfig.String("CumplidosDveUrlWso2")+"servicios_academicos"+"/"+"carrera_snies/"+id_proyecto_snies, &temp_snies); err == nil && temp_snies != nil {
					json_info_coordinador, error_json := json.Marshal(temp_snies)
					if error_json == nil {
						var temp_info_coordinador models.InformacionCoordinador
						if err = json.Unmarshal(json_info_coordinador, &temp_info_coordinador); err == nil {
							info_coordinador = temp_info_coordinador
						} else {
							panic(err.Error())
						}
					} else {
						panic(error_json.Error())
					}
				} else {
					panic(err.Error())
				}
			} else {
				panic(err.Error())
			}
		} else {
			panic(error_json.Error())
		}
	} else {
		panic(err.Error())
	}
	return info_coordinador, outputError
}

func CargarContratosDocente(numDocumento int) (contratosDocentes []models.ContratosDocente, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "CargarContratosDocente", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var err error
	var cd models.ContratosDocente
	var proveedor []models.InformacionProveedor
	var vinculaciones []models.VinculacionDocente
	var actasInicio []models.ActaInicio
	var res []models.Resolucion
	var dep models.Dependencia
	var parametro []models.Parametro

	//CONSULTA LA INFORMACION DE PROVEEDOR
	if err = GetRequestLegacy("CumplidosDveUrlCrudAgora", "informacion_proveedor/?query=num_documento:"+strconv.Itoa(numDocumento), &proveedor); err == nil {
		//CONSULTA LA VINCULACION DEL DOCENTE
		if err = GetRequestNew("CumplidosDveUrlCrudResoluciones", "vinculacion_docente/?query=PersonaId:"+strconv.Itoa(numDocumento)+"&limit=-1", &vinculaciones); err == nil {
			for _, vinculacion := range vinculaciones {
				if err = GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=Id:"+strconv.FormatInt(vinculacion.DedicacionId, 10), &parametro); err == nil {
					//CONSULTA LA DEPENDENCIA DE CADA VINCULACION
					if err = GetRequestLegacy("CumplidosDveUrlCrudOikos", "dependencia/"+strconv.Itoa(vinculacion.ProyectoCurricularId), &dep); err == nil {
						//CONSULTA A RESOLUCION
						if err = GetRequestNew("CumplidosDveUrlCrudResoluciones", "resolucion/?query=Id:"+strconv.Itoa(vinculacion.ResolucionVinculacionDocenteId.Id), &res); err == nil {
							//CONSULTA LAS ACTAS DE INICIO
							if err = GetRequestLegacy("CumplidosDveUrlCrudAgora", "acta_inicio/?query=NumeroContrato:"+vinculacion.NumeroContrato+",Vigencia:"+strconv.FormatInt(vinculacion.Vigencia, 10), &actasInicio); err == nil {
								for _, actaInicio := range actasInicio {
									actaInicio.FechaInicio = actaInicio.FechaInicio.UTC()
									actaInicio.FechaFin = actaInicio.FechaFin.UTC()

									//fechaInicio<fechaActual<(fechaFin+2 meses) Se da holgura de 2 meses luego de fechaFin para subir cumplidos
									//if time.Now().After(actaInicio.FechaInicio) && time.Now().Before(actaInicio.FechaFin.AddDate(0, 2, 0)) {

									cd.NumeroVinculacion = vinculacion.NumeroContrato
									cd.Vigencia = vinculacion.Vigencia
									cd.Resolucion = res[0].NumeroResolucion
									cd.Dependencia = dep.Nombre
									cd.IdDependencia = dep.Id
									cd.NombreDocente = proveedor[0].NomProveedor
									cd.Dedicacion = parametro[0].Nombre
									contratosDocentes = append(contratosDocentes, cd)
									//}
								}
							} else {
								panic(err.Error())
							}
						} else {
							panic(err.Error())
						}
					}
				} else {
					panic(err.Error())
				}
			}
		} else {
			panic(err.Error())
		}
	} else {
		panic(err.Error())
	}
	return contratosDocentes, outputError
}
