package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/udistrital/cumplidos_dve_mid/helpers"
	"github.com/udistrital/cumplidos_dve_mid/models"
)

// AprobacionDocumentosController operations for AprobacionDocumentos
type AprobacionDocumentosController struct{
	beego.Controller
}

// URL Mapping ...
func (c *AprobacionDocumentosController) URLMapping() {
	c.Mapping("SolicitudesSupervisor", c.SolicitudesSupervisor)
	c.Mapping("SolicitudesCoordinador", c.SolicitudesCoordinador)
	c.Mapping("CertificacionVistoBueno", c.CertificacionVistoBueno)
	c.Mapping("AprobarSolicitudes", c.AprobarSolicitudes)
}

// AprobacionDocumentosController ...
// @Title SolicitudesSupervisor
// @Description create SolicitudesSupervisor
// @Param docsupervisor query string true "Número del documento del supervisor"
// @Success 201
// @Failure 403 :docsupervisor is empty
// @router /solicitudes_supervisor/:docsupervisor [get]
func (c *AprobacionDocumentosController) SolicitudesSupervisor(){
	defer helpers.ErrorController(c.Controller, "AprobacionDocumentosController")

	doc_supervisor := c.GetString(":docsupervisor")

	if doc_supervisor == "" {
		panic(map[string]interface{}{"funcion": "SolicitudesSupervisor", "err": helpers.ErrorParametros, "status": "400"})
	}

	if data, err2:= helpers.GetSolicitudesSupervisor(doc_supervisor); err2 == nil && data != nil{
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Solicitudes del supervisor cargadas con exito", "Data": data}
	}else{
		panic(map[string]interface{}{"funcion": "SolicitudesSupervisor", "err": err2, "status": "400"})
	}
	c.ServeJSON()
}

// AprobacionDocumentosController ...
// @Title SolicitudesCoordinador
// @Description create SolicitudesCoordinador
// @Param doccoordinador query string true "Número del documento del coordinador"
// @Success 201
// @Failure 403 :doccoordinador is empty
// @router /solicitudes_coordinador/:doccoordinador [get]
func (c *AprobacionDocumentosController) SolicitudesCoordinador() {
	defer helpers.ErrorController(c.Controller, "AprobacionDocumentosController")

	doc_coordinador := c.GetString(":doccoordinador")

	if doc_coordinador == "" {
		panic(map[string]interface{}{"funcion": "SolicitudesCoordinador", "err": helpers.ErrorParametros, "status": "400"})
	}

	if data, err2:= helpers.GetSolicitudesCoordinador(doc_coordinador); err2 == nil && data != nil{
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Solicitudes del coordinador cargadas con exito", "Data": data}
	}else{
		panic(map[string]interface{}{"funcion": "SolicitudesCoordinador", "err": err2, "status": "400"})
	}
	c.ServeJSON()
}

// AprobacionDocumentosController ...
// @Title CertificacionVistoBueno
// @Description create CertificacionVistoBueno
// @Param dependencia query int true "Dependencia del contrato en la tabla vinculacion"
// @Param mes query int true "Mes del pago mensual"
// @Param anio query int true "Año del pago mensual"
// @Success 201
// @Failure 403 :dependencia is empty
// @Failure 403 :mes is empty
// @Failure 403 :anio is empty
// @router /certificacion_visto_bueno/:dependencia/:mes/:anio [get]
func (c *AprobacionDocumentosController) CertificacionVistoBueno(){
	defer helpers.ErrorController(c.Controller, "AprobacionDocumentosController")
	
	dependencia := c.GetString(":dependencia")
	mes := c.GetString(":mes")
	anio := c.GetString(":anio")

	if dependencia == "" && mes == "" && anio == ""{
		panic(map[string]interface{}{"funcion": "CertificacionVistoBueno", "err": helpers.ErrorParametros, "status": "400"})
	}

	if data, err2:= helpers.CertificacionVistoBueno(dependencia, mes, anio); err2 == nil && data != nil{
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Certificaciones de visto bueno cargadas con exito", "Data": data}
	}else{
		panic(map[string]interface{}{"funcion": "CertificacionVistoBueno", "err": err2, "status": "400"})
	}
	c.ServeJSON()
}

// AprobacionDocumentosController ...
// @Title AprobarSolicitudes
// @Description create AprobarSolicitudes
// @Success 201
// @Failure 403
// @router /aprobar_documentos [post]
func (c *AprobacionDocumentosController) AprobarSolicitudes(){
	defer helpers.ErrorController(c.Controller, "AprobacionDocumentosController")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e !=nil{
		panic(map[string]interface{}{"funcion": "AprobarSolicitudes", "err": helpers.ErrorBody, "status": "400"})
	}

	var m []models.PagoPersonaProyecto

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil{
		if res, err := helpers.AprobarMultiplesSolicitudes(m); err == nil{
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": res, "Data": m}
		}else{
			panic(err)
		}
	}else{
		panic(map[string]interface{}{"funcion": "AprobarSolicitudes", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}