// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/cumplidos_dve_mid/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/informacion_academica",
			beego.NSInclude(
				&controllers.InformacionAcademicaController{},
			),
		),
		beego.NSNamespace("/aprobacion_pago",
			beego.NSInclude(
				&controllers.AprobacionPagoController{},
			),
		),
		beego.NSNamespace("/aprobacion_documentos",
			beego.NSInclude(
				&controllers.AprobacionDocumentosController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
