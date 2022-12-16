package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

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
