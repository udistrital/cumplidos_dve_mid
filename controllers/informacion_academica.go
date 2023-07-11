package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/udistrital/cumplidos_dve_mid/helpers"
	xray2 "github.com/udistrital/cumplidos_dve_mid/xray"
)

// InformacionAcademicaController operations for InformacionAcademica
type InformacionAcademicaController struct {
	beego.Controller
}

// URLMapping ...
func (c *InformacionAcademicaController) URLMapping() {
	c.Mapping("ObtenerInfoCoordinador", c.ObtenerInfoCoordinador)
	c.Mapping("GetContratosDocente", c.GetContratosDocente)
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

	///Segmento
	ctx := c.Ctx.Request.Context()
	ctx, seg := xray2.BeginSegmentWithContextTP(ctx, "Cumplidos_DVE_MID", c.Ctx.Request.Method, c.Ctx.Request.URL.String(), 200, c.Ctx.Request.URL.String(), c.Ctx.Request.Header.Values("X-Amzn-Trace-Id"))
	defer seg.Close(nil)

	//subsegmento
	_, subseg := xray.BeginSubsegment(ctx, "ObtenerInfoCoordinador")
	defer subseg.Close(nil)
	xray2.SetContext(ctx)

	id_oikos := c.Ctx.Input.Param(":id_dependencia_oikos")
	id, err2 := strconv.Atoi(id_oikos)

	if err2 != nil || id <= 0 {
		xray2.BeginSubSegmentWithContext(subseg, c.Ctx.Request.Method, c.Ctx.Request.URL.String(), 400)
		xray.AddError(ctx, fmt.Errorf("%v", helpers.ErrorParametros))
		panic(map[string]interface{}{"funcion": "ObtenerInfoCoordinador", "err": helpers.ErrorParametros, "status": "400", "message": "Error de registro"})
	}

	if data, err3 := helpers.CargarInformacionCoordinador(id); err3 == nil && data.CarreraSniesCollection.CarreraSnies != nil {
		c.Ctx.Output.SetStatus(200)
		xray2.BeginSubSegmentWithContext(subseg, c.Ctx.Request.Method, c.Ctx.Request.URL.String(), 200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Información del coordinador cargada con exito", "Data": data}
	} else {
		// Manejo de error y rastreo en caso de fallo
		xray2.BeginSubSegmentWithContext(subseg, c.Ctx.Request.Method, c.Ctx.Request.URL.String(), 400)
		xray.AddError(ctx, fmt.Errorf("%v", helpers.ErrorParametros))
		panic(map[string]interface{}{"funcion": "ObtenerInfoCoordinador", "err": err3, "status": "400"})
	}

	c.Ctx.Request = c.Ctx.Request.WithContext(ctx)
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

	///Segmento
	ctx := c.Ctx.Request.Context()
	ctx, seg := xray2.BeginSegmentWithContextTP(ctx, "Cumplidos_DVE_MID", c.Ctx.Request.Method, c.Ctx.Request.URL.String(), 200, c.Ctx.Request.URL.String(), c.Ctx.Request.Header.Values("X-Amzn-Trace-Id"))
	defer seg.Close(nil)

	//subsegmento
	_, subseg := xray.BeginSubsegment(ctx, "GetContratosDocente")
	defer subseg.Close(nil)
	xray2.SetContext(ctx)

	numDocumento := c.Ctx.Input.Param(":numDocumento")
	doc, err := strconv.Atoi(numDocumento)

	if err != nil || doc <= 0 {
		xray2.BeginSubSegmentWithContext(subseg, c.Ctx.Request.Method, c.Ctx.Request.URL.String(), 400)
		xray.AddError(ctx, fmt.Errorf("%v", helpers.ErrorParametros))
		panic(map[string]interface{}{"funcion": "GetContratosDocente", "err": helpers.ErrorParametros, "status": "400"})
	}

	if data, err2 := helpers.CargarContratosDocente(doc); err2 == nil {
		c.Ctx.Output.SetStatus(200)
		xray2.BeginSubSegmentWithContext(subseg, c.Ctx.Request.Method, c.Ctx.Request.URL.String(), 200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Contratos del Docente cargados con exito", "Data": data}
	} else {
		xray2.BeginSubSegmentWithContext(subseg, c.Ctx.Request.Method, c.Ctx.Request.URL.String(), 400)
		xray.AddError(ctx, fmt.Errorf("%v", helpers.ErrorParametros))
		panic(map[string]interface{}{"funcion": "GetContratosDocente", "err": err2, "status": "400"})
	}

	c.Ctx.Request = c.Ctx.Request.WithContext(ctx)
	c.ServeJSON()
}
