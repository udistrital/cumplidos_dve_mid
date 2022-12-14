package controllers

import (
	"encoding/json"
	"fmt"
	
	"github.com/astaxie/beego"
	"github.com/udistrital/cumplidos_dve_mid/helpers"
	"github.com/udistrital/cumplidos_dve_mid/models"
)

// InformacionAcademicaController operations for InformacionAcademica
type InformacionAcademicaController struct{
	beego.Controller
}

// URLMapping ...
func (c *InformacionAcademicaController) URLMapping(){
	c.Mapping("ObtenerInfoCoordinador", c.ObtenerInfoCoordinador)
}

// InformacionAcademicaController ...
// @Title ObtenerInfoCoordinador
// @Description create ObtenerInfoCoordinador
// @Param id_dependencia_oikos query int true "Proyecto a obtener información coordinador"
// @Success 201 {int} models.InformacionCoordinador
// @Failure 403 :id_dependencia_oikos is empty
// @router /informacion_coordinador/:id_dependencia_oikos [get]
func (c *InformacionAcademicaController) ObtenerInfoCoordinador(){
	defer helpers.ErrorController(c.Controller, "InformacionAcademicaController")

	id_oikos := c.GetString(":id_dependencia_oikos")
	var temp map[string]interface{}
	var temp_snies map[string]interface{}
	var info_coordinador models.InformacionCoordinador

	if err := helpers.GetJsonWSO2(beego.AppConfig.String("CumplidosDveUrlWso2") + beego.AppConfig.String("CumplidosDveHomologacion")+"/"+"proyecto_curricular_oikos/"+id_oikos, &temp); err == nil && temp != nil {
		json_proyecto_curricular, error_json := json.Marshal(temp)

		if error_json == nil {
			var temp_homologacion models.ObjetoProyectoCurricular
			if err := json.Unmarshal(json_proyecto_curricular, &temp_homologacion); err == nil {
				id_proyecto_snies := temp_homologacion.Homologacion.IDSnies

				if err := helpers.GetJsonWSO2(beego.AppConfig.String("CumplidosDveUrlWso2") + beego.AppConfig.String("CumplidosDveAcademica")+"/"+"carrera_snies/"+id_proyecto_snies, &temp_snies); err == nil && temp_snies != nil {
					json_info_coordinador, error_json := json.Marshal(temp_snies)

					if error_json == nil {
						var temp_info_coordinador models.InformacionCoordinador
						if err := json.Unmarshal(json_info_coordinador, &temp_info_coordinador); err == nil {

							fmt.Println(temp_info_coordinador)
							info_coordinador = temp_info_coordinador
						} else {
							fmt.Println(err)
						}
					} else {
						fmt.Println(error_json.Error())
					}
				}

			} else {
				fmt.Println(err)
			}

		} else {
			fmt.Println(error_json.Error())
		}
	} else {
		fmt.Println(err)

	}

	c.Data["json"] = info_coordinador
	c.ServeJSON()
}