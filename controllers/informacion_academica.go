package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/cumplidos_dve_mid/helpers"
	"github.com/udistrital/cumplidos_dve_mid/services"
)

// InformacionAcademicaController operations for InformacionAcademica
type InformacionAcademicaController struct {
	beego.Controller
}

// URLMapping ...
func (c *InformacionAcademicaController) URLMapping() {
	c.Mapping("ObtenerInfoCoordinador", c.ObtenerInfoCoordinador)
	c.Mapping("GetContratosDocente", c.GetContratosDocente)
	c.Mapping("GetDocentesCoordinador", c.GetDocentesCoordinador)
}

// InformacionAcademicaController ...
// @Title ObtenerInfoCoordinador
// @Description create ObtenerInfoCoordinador
// @Param id_dependencia_oikos query int true "Proyecto a obtener información coordinador"
// @Success 201 {int} models.InformacionCoordinador
// @Failure 400 Bad Request
// @Failure 403 :id_dependencia_oikos is empty
// @router /informacion_coordinador/:id_dependencia_oikos [get]
func (c *InformacionAcademicaController) ObtenerInfoCoordinador() {
	defer helpers.ErrorController(c.Controller, "InformacionAcademicaController")

	id_oikos := c.Ctx.Input.Param(":id_dependencia_oikos")
	id, err2 := strconv.Atoi(id_oikos)

	if err2 != nil || id <= 0 {
		panic(map[string]interface{}{"funcion": "ObtenerInfoCoordinador", "err": helpers.ErrorParametros, "status": "400", "message": "Error de registro"})
	}

	if data, err3 := helpers.CargarInformacionCoordinador(id); err3 == nil && data.CarreraSniesCollection.CarreraSnies != nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Información del coordinador cargada con exito", "Data": data}
	} else {
		panic(map[string]interface{}{"funcion": "ObtenerInfoCoordinador", "err": err3, "status": "400"})
	}
	c.ServeJSON()
}

// InformacionAcademicaController ...
// @Title GetContratosDocente
// @Description create GetContratosDocente
// @Param numDocumento query string true "Docente a consultar"
// @Success 201 {object} []models.ContratosDocentes
// @Failure 403 body is empty
// @router /contratos_docente/:numDocumento [get]
func (c *InformacionAcademicaController) GetContratosDocente() {
	defer helpers.ErrorController(c.Controller, "InformacionAcademicaController")

	numDocumento := c.Ctx.Input.Param(":numDocumento")
	doc, err := strconv.Atoi(numDocumento)

	if err != nil || doc <= 0 {
		panic(map[string]interface{}{"funcion": "GetContratosDocente", "err": helpers.ErrorParametros, "status": "400"})
	}

	if data, err2 := helpers.CargarContratosDocente(doc); err2 == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Contratos del Docente cargados con exito", "Data": data}
	} else {
		panic(map[string]interface{}{"funcion": "GetContratosDocente", "err": err2, "status": "400"})
	}

	c.ServeJSON()
}

// @Title GetDocentesCoordinador
// @Description Retorna docentes a cargo del proyecto curricular (por proyectoId, vigencia, mes y año)
// @Param   proyectoId   path   int  true  "ID del proyecto curricular (OIKOS)"
// @Param   vigencia     path   int  true  "Vigencia del contrato"
// @Param   mes          path   int  true  "Mes (1-12)"
// @Param   anio         path   int  true  "Año (ej: 2026)"
// @Success 200 {object} models.DocentesCoordinadorResponse
// @Failure 400 Bad Request
// @router /docentes_coordinador/:proyectoId/:vigencia/:mes/:anio [get]
func (c *InformacionAcademicaController) GetDocentesCoordinador() {
	defer helpers.ErrorController(c.Controller, "InformacionAcademicaController")

	proyectoStr := c.Ctx.Input.Param(":proyectoId")
	proyectoId, err := strconv.Atoi(proyectoStr)
	if err != nil || proyectoId <= 0 {
		panic(map[string]interface{}{"funcion": "GetDocentesCoordinador", "err": helpers.ErrorParametros, "status": "400"})
	}

	vigenciaStr := c.Ctx.Input.Param(":vigencia")
	vigencia, err := strconv.Atoi(vigenciaStr)
	if err != nil || vigencia <= 0 {
		panic(map[string]interface{}{"funcion": "GetDocentesCoordinador", "err": helpers.ErrorParametros, "status": "400"})
	}

	mesStr := c.Ctx.Input.Param(":mes")
	mes, err := strconv.Atoi(mesStr)
	if err != nil || mes < 1 || mes > 12 {
		panic(map[string]interface{}{"funcion": "GetDocentesCoordinador", "err": helpers.ErrorParametros, "status": "400"})
	}

	anioStr := c.Ctx.Input.Param(":anio")
	anio, err := strconv.Atoi(anioStr)
	if err != nil || anio <= 0 {
		panic(map[string]interface{}{"funcion": "GetDocentesCoordinador", "err": helpers.ErrorParametros, "status": "400"})
	}

	data, svcErr := services.GetDocentesProyecto(proyectoId, vigencia, mes, anio)
	if svcErr != nil {
		panic(map[string]interface{}{"funcion": "GetDocentesCoordinador", "err": fmt.Errorf("%v", svcErr), "status": "400"})
	}

	c.Data["json"] = map[string]interface{}{
		"Success": true,
		"Status":  "200",
		"Message": "Docentes cargados",
		"Data":    data,
	}
	c.ServeJSON()
}
