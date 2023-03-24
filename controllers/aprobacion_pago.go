package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/cumplidos_dve_mid/helpers"
	"github.com/udistrital/cumplidos_dve_mid/models"
)

//AprobacionPagoController operations for AprobacionPago
type AprobacionPagoController struct{
	beego.Controller
}

// URLMapping ...
func (c *AprobacionPagoController) URLMapping(){
	c.Mapping("CertificacionDocumentosAprobados", c.CertificacionDocumentosAprobados)
	c.Mapping("PagoAprobado", c.PagoAprobado)
	c.Mapping("SolicitudesOrdenador", c.SolicitudesOrdenador)
	c.Mapping("DependenciaOrdenador", c.DependenciaOrdenador)
	c.Mapping("AprobarPagos", c.AprobarPagos)
	c.Mapping("InfoOrdenador", c.InfoOrdenador)
}

// AprobacionPagoController ...
// @Title CertificacionDocumentosAprobados
// @Description create CertificacionDocumentosAprobados
// @Param dependencia query int true "Dependencia del contrato en la tabla ordenador_gasto"
// @Param mes query int true "Mes del pago mensual"
// @Param anio query int true "Año del pago mensual"
// @Success 201
// @Failure 403 :dependencia is empty
// @Failure 403 :mes is empty
// @Failure 403 :anio is empty
// @router /certificacion_documentos_aprobados/:dependencia/:mes/:anio [get]
func (c *AprobacionPagoController) CertificacionDocumentosAprobados() {
	defer helpers.ErrorController(c.Controller, "AprobacionPagoController")

	dependencia := c.GetString(":dependencia")
	mes := c.GetString(":mes")
	anio := c.GetString(":anio")
	
	if dependencia == "" && mes == "" && anio == "" {
		panic(map[string]interface{}{"funcion": "CertificacionDocumentosAprobados", "err": helpers.ErrorParametros, "status": "400"})
	}

	if data, err2:= helpers.CargarCertificacionDocumentosAprobados(dependencia, mes, anio); err2 == nil && data != nil{
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Certificación documentos aprobados cargada con exito", "Data": data}
	}else{
		panic(map[string]interface{}{"funcion": "CertificacionDocumentosAprobados", "err": err2, "status": "400"})
	}
	c.ServeJSON()
}

// AprobacionPagoController ...
// @Title PagoAprobado
// @Description create PagoAprobado
// @Param numero_contrato query int true "Numero de contrato en la tabla contrato general"
// @Param vigencia query int true "Vigencia del contrato en la tabla contrato general"
// @Param mes query int true "Mes del pago mensual"
// @Param anio query int true "Año del pago mensual"
// @Success 201
// @Failure 403 :numero_contrato is empty
// @Failure 403 :vigencia is empty
// @Failure 403 :mes is empty
// @Failure 403 :anio is empty
// @router /pago_aprobado/:numero_contrato/:vigencia/:mes/:anio [get]
func (c *AprobacionPagoController) PagoAprobado(){
	defer helpers.ErrorController(c.Controller, "AprobacionPagoController")

	numero_contrato := c.GetString(":numero_contrato")
	vigencia := c.GetString(":vigencia")
	mes := c.GetString(":mes")
	anio := c.GetString(":anio")

	if numero_contrato == "" && vigencia == "" && mes == "" && anio == "" {
		panic(map[string]interface{}{"funcion": "PagoAprobado", "err": helpers.ErrorParametros, "status": "400"})
	}

	if data, err2:= helpers.ConsultarPagoAprobado(numero_contrato, vigencia, mes, anio); err2 == nil{
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Estado del pago cargado con exito", "Data": data}
	}else{
		panic(map[string]interface{}{"funcion": "PagoAprobado", "err": err2, "status": "400"})
	}
	c.ServeJSON()
}

// AprobacionPagoController ...
// @Title SolicitudesOrdenador
// @Description create SolicitudesOrdenador
// @Param docordenador query string true "Número del documento del ordenador"
// @Success 201
// @Failure 403 :docordenador is empty
// @router /solicitudes_ordenador/:docordenador [get]
func (c *AprobacionPagoController) SolicitudesOrdenador(){
	defer helpers.ErrorController(c.Controller, "AprobacionPagoController")

	doc_ordenador := c.GetString(":docordenador")
	limit, err0 := c.GetInt("limit")
	offset, err0 := c.GetInt("offset")

	if doc_ordenador == "" && limit <= 0 && offset <= 0{
		panic(map[string]interface{}{"funcion": "SolicitudesOrdenador", "err": helpers.ErrorParametros, "status": "400"})
	}

	if data, err2:= helpers.CargarSolicitudesOrdenador(doc_ordenador, limit, offset, err0); err2 == nil && data != nil{
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Solicitudes del ordenador cargadas con exito", "Data": data}
	}else{
		panic(map[string]interface{}{"funcion": "SolicitudesOrdenador", "err": err2, "status": "400"})
	}
	c.ServeJSON()
}

// AprobacionPagoController ...
// @Title DependenciaOrdenador
// @Description create DependenciaOrdenador
// @Param docordenador query string true "Número del documento del ordenador"
// @Success 201
// @Failure 403 :docordenador is empty
// @router /dependencia_ordenador/:docordenador [get]
func (c *AprobacionPagoController) DependenciaOrdenador(){
	defer helpers.ErrorController(c.Controller, "AprobacionPagoController")

	doc_ordenador := c.GetString(":docordenador")

	if doc_ordenador == ""{
		panic(map[string]interface{}{"funcion": "DependenciaOrdenador", "err": helpers.ErrorParametros, "status": "400"})
	}

	if data, err2:= helpers.ObtenerDependenciaOrdenador(doc_ordenador); err2 == nil && data != 0{
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Dependencia del ordenador cargadas con exito", "Data": data}
	}else{
		panic(map[string]interface{}{"funcion": "DependenciaOrdenador", "err": err2, "status": "400"})
	}
	c.ServeJSON()
}

// AprobacionPagoController ...
// @Title ObtenerInfoOrdenador
// @Description create ObtenerInfoOrdenador
// @Param numero_contrato query int true "Numero de contrato en la tabla contrato general"
// @Param vigencia query int true "Vigencia del contrato en la tabla contrato general"
// @Success 201 {int} models.InformacionOrdenador
// @Failure 403 :numero_contrato is empty
// @Failure 403 :vigencia is empty
// @router /informacion_ordenador/:numero_contrato/:vigencia [get]
func (c *AprobacionPagoController) InfoOrdenador() {
	defer helpers.ErrorController(c.Controller, "AprobacionPagoController")

	numero_contrato := c.GetString(":numero_contrato")
	vigencia := c.GetString(":vigencia")

	if numero_contrato == "" && vigencia == ""{
		panic(map[string]interface{}{"funcion": "InfoOrdenador", "err": helpers.ErrorParametros, "status": "400"})
	}

	if data, err2:= helpers.ObtenerInfoOrdenador(numero_contrato, vigencia); err2 == nil && data.NumeroDocumento != 0{
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": 200, "Message": "Informacion ordenador cargada con exito", "Data": data}
	}else{
		panic(map[string]interface{}{"funcion": "InfoOrdenador", "err": err2, "status": "400"})
	}
	c.ServeJSON()
}

// AprobacionPagoController ...
// @Title AprobarPagos
// @Description create AprobarPagos
// @Success 201
// @Failure 403
// @router /aprobar_pagos [post]
func (c *AprobacionPagoController) AprobarPagos(){
	defer helpers.ErrorController(c.Controller, "AprobacionPagoController")

	if v, e := helpers.ValidarBody(c.Ctx.Input.RequestBody); !v || e !=nil{
		panic(map[string]interface{}{"funcion": "AprobarPagos", "err": helpers.ErrorBody, "status": "400"})
	}

	var m []models.PagoMensual

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &m); err == nil{
		if res, err := helpers.AprobarMultiplesPagos(m); err == nil{
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": res, "Data": m}
		}else{
			panic(err)
		}
	}else{
		panic(map[string]interface{}{"funcion": "AprobarPagos", "err": err.Error(), "status": "400"})
	}
	c.ServeJSON()
}