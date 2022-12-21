package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:AprobacionPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:AprobacionPagoController"],
        beego.ControllerComments{
            Method: "AprobarPagos",
            Router: "/aprobar_pagos",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:AprobacionPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:AprobacionPagoController"],
        beego.ControllerComments{
            Method: "CertificacionDocumentosAprobados",
            Router: "/certificacion_documentos_aprobados/:dependencia/:mes/:anio",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:AprobacionPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:AprobacionPagoController"],
        beego.ControllerComments{
            Method: "DependenciaOrdenador",
            Router: "/dependencia_ordenador/:docordenador",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:AprobacionPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:AprobacionPagoController"],
        beego.ControllerComments{
            Method: "InfoCoordinador",
            Router: "/informacion_coordinador/:id_dependencia_oikos",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:AprobacionPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:AprobacionPagoController"],
        beego.ControllerComments{
            Method: "PagoAprobado",
            Router: "/pago_aprobado/:numero_contrato/:vigencia/:mes/:anio",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:AprobacionPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:AprobacionPagoController"],
        beego.ControllerComments{
            Method: "SolicitudesOrdenador",
            Router: "/solicitudes_ordenador/:docordenador",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:InformacionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:InformacionAcademicaController"],
        beego.ControllerComments{
            Method: "GetContratosDocente",
            Router: "/contratos_docente/:numDocumento",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:InformacionAcademicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/cumplidos_dve_mid/controllers:InformacionAcademicaController"],
        beego.ControllerComments{
            Method: "ObtenerInfoCoordinador",
            Router: "/informacion_coordinador/:id_dependencia_oikos",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
