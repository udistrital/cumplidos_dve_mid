package helpers

import (
	"encoding/json"
	"strconv"
	//"time"

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

	if err = GetJsonWSO2(beego.AppConfig.String("CumplidosDveUrlWso2") + beego.AppConfig.String("CumplidosDveHomologacion") + "/" + "proyecto_curricular_oikos/" + strconv.Itoa(DependenciaOikosId), &temp); err == nil && temp != nil {
		json_proyecto_curricular, error_json := json.Marshal(temp)
		if error_json == nil {
			var temp_homologacion models.ObjetoProyectoCurricular
			if err = json.Unmarshal(json_proyecto_curricular, &temp_homologacion); err == nil {
				id_proyecto_snies := temp_homologacion.Homologacion.IDSnies
				if err = GetJsonWSO2(beego.AppConfig.String("CumplidosDveUrlWso2") + beego.AppConfig.String("CumplidosDveAcademica") + "/" + "carrera_snies/" + id_proyecto_snies, &temp_snies); err == nil && temp_snies != nil {
					json_info_coordinador, error_json := json.Marshal(temp_snies)
					if error_json == nil {
						var temp_info_coordinador models.InformacionCoordinador
						if err = json.Unmarshal(json_info_coordinador, &temp_info_coordinador); err == nil {
							info_coordinador = temp_info_coordinador
						} else{
							panic(err.Error())
						}
					}else{
						panic(error_json.Error())
					}
				}else{
					panic(err.Error())
				}
			}else{
				panic(err.Error())
			}
		}else{
			panic(error_json.Error())
		}
	}else{
		panic(err.Error())
	}
	return info_coordinador, outputError
}

func CargarContratosDocente(numDocumento int) (contratosDocentes []models.ContratosDocente, outputError map[string]interface{}){
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
	var res models.Resolucion
	var dep models.Dependencia

	//If informacion_proveedor get
	if err = GetJson(beego.AppConfig.String("CumplidosDveUrlCrudAgora") + "informacion_proveedor/?query=num_documento:" + strconv.Itoa(numDocumento), &proveedor); err == nil {
		//If vinculacion_docente get
		if err = GetJson(beego.AppConfig.String("CumplidosDveUrlCrudAdmin") + "vinculacion_docente/?query=IdPersona:" + strconv.Itoa(numDocumento) + "&limit=-1", &vinculaciones); err == nil {
			//for vinculaciones
			for _, vinculacion := range vinculaciones {
				//If dependencia get
				if err = GetJson(beego.AppConfig.String("CumplidosDveUrlCrudOikos") + "dependencia/" + strconv.Itoa(vinculacion.IdProyectoCurricular), &dep); err == nil{
					//If resolucion get
					if err = GetJson(beego.AppConfig.String("CumplidosDveUrlCrudAdmin") + "resolucion/" + strconv.Itoa(vinculacion.IdResolucion.Id), &res); err == nil{
						//If nulo
						//if vinculacion.NumeroContrato.Valid == true {
							if err = GetJson(beego.AppConfig.String("CumplidosDveUrlCrudAgora") + "acta_inicio/?query=NumeroContrato:" + vinculacion.NumeroContrato.String + ",Vigencia:" + strconv.FormatInt(vinculacion.Vigencia.Int64, 10), &actasInicio); err == nil{

								//If Estado = 4
								for _, actaInicio := range actasInicio {
									actaInicio.FechaInicio = actaInicio.FechaInicio.UTC()
									actaInicio.FechaFin = actaInicio.FechaFin.UTC()

									//fechaInicio<fechaActual<(fechaFin+2 meses) Se da holgura de 2 meses luego de fechaFin para subir cumplidos 
									//if time.Now().After(actaInicio.FechaInicio) && time.Now().Before(actaInicio.FechaFin.AddDate(0, 2, 0)) {
									
										cd.NumeroVinculacion = vinculacion.NumeroContrato.String
										cd.Vigencia = vinculacion.Vigencia.Int64
										cd.Resolucion = res.NumeroResolucion
										cd.Dependencia = dep.Nombre
										cd.IdDependencia = dep.Id
										cd.NombreDocente = proveedor[0].NomProveedor
										cd.Dedicacion = vinculacion.IdDedicacion.NombreDedicacion
										contratosDocentes = append(contratosDocentes, cd)
									//}
								}
							}else{
								panic(err.Error())
							}
						//}
					}else{
						panic(err.Error())
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
	return contratosDocentes, outputError
}