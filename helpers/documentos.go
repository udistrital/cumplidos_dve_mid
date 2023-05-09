package helpers

import (
	"fmt"
	"time"
	"bufio"
	"bytes"
	"strconv"
	"net/url"
	"path/filepath"
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/phpdave11/gofpdf"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/cumplidos_dve_mid/models"
)

func GetSolicitudesSupervisor(doc_supervisor string) (pagos_personas_proyecto []models.PagoPersonaProyecto, outputError map[string]interface{}){
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "GetSolicitudesSupervisor", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	var parametro []models.Parametro
	var pagos_mensuales []models.PagoMensual
	var contratistas []models.InformacionProveedor
	var pago_personas_proyecto models.PagoPersonaProyecto
	var vinculaciones_docente []models.VinculacionDocente
	
	if err:=GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion:PRS_DVE", &parametro); err == nil{
		if err := GetRequestNew("CumplidosDveUrlCrud", "pago_mensual/?limit=-1&query=Responsable:" + doc_supervisor + ",EstadoPagoMensualId:" + strconv.Itoa(parametro[0].Id), &pagos_mensuales); err == nil{
			for x, pago_mensual := range pagos_mensuales{
				if err := GetRequestLegacy("CumplidosDveUrlCrudAgora", "informacion_proveedor/?query=NumDocumento:" + pago_mensual.Persona, &contratistas); err == nil{
					for _, contratista := range contratistas{
						if err := GetRequestNew("CumplidosDveUrlCrudResoluciones", "vinculacion_docente/?limit=-1&query=NumeroContrato:" + pago_mensual.NumeroContrato + ",Vigencia:" + strconv.FormatFloat(pago_mensual.VigenciaContrato, 'f', 0, 64), &vinculaciones_docente); err == nil{
							for _, vinculacion := range vinculaciones_docente{
								var dep models.Dependencia
								if err := GetRequestLegacy("CumplidosDveUrlCrudOikos", "dependencia/" + strconv.Itoa(vinculacion.ProyectoCurricularId), &dep); err == nil{
									pago_personas_proyecto.PagoMensual = &pagos_mensuales[x]
									pago_personas_proyecto.NombrePersona = contratista.NomProveedor
									pago_personas_proyecto.Dependencia = &dep
									pagos_personas_proyecto = append(pagos_personas_proyecto, pago_personas_proyecto)
								}else{
									panic(err.Error())
								}
							}
						}else{
							panic(err.Error())
						}
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

	return pagos_personas_proyecto, outputError
}

func GetSolicitudesCoordinador(doc_coordinador string) (pagos_personas_proyecto []models.PagoPersonaProyecto, outputError map[string]interface{}){
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "GetSolicitudesCoordinador", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var parametro []models.Parametro
	var pagos_mensuales []models.PagoMensual
	var contratistas []models.InformacionProveedor
	var pago_personas_proyecto models.PagoPersonaProyecto
	var vinculaciones_docente []models.VinculacionDocente
	
	if err:= GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion:PRC_DVE", &parametro); err == nil{
		if err := GetRequestNew("CumplidosDveUrlCrud", "pago_mensual/?limit=-1&query=Responsable:" + doc_coordinador +",EstadoPagoMensualId:" + strconv.Itoa(parametro[0].Id), &pagos_mensuales); err == nil{
			for x, _ := range pagos_mensuales{
				if err := GetRequestLegacy("CumplidosDveUrlCrudAgora", "informacion_proveedor/?query=NumDocumento:" + pagos_mensuales[x].Persona, &contratistas); err == nil{
					for _, contratista := range contratistas{
						if err := GetRequestNew("CumplidosDveUrlCrudResoluciones", "vinculacion_docente/?limit=-1&query=NumeroContrato:" + pagos_mensuales[x].NumeroContrato + ",Vigencia:" + strconv.FormatFloat(pagos_mensuales[x].VigenciaContrato, 'f', 0, 64), &vinculaciones_docente);err == nil{
							for y, _ := range vinculaciones_docente{
								var dep []models.Dependencia
								if err := GetRequestLegacy("CumplidosDveUrlCrudOikos", "dependencia/?query=Id:" + strconv.Itoa(vinculaciones_docente[y].ProyectoCurricularId), &dep); err == nil{
									for z, _ := range dep{
										pago_personas_proyecto.PagoMensual = &pagos_mensuales[x]
										pago_personas_proyecto.NombrePersona = contratista.NomProveedor
										pago_personas_proyecto.Dependencia = &dep[z]
										pagos_personas_proyecto = append(pagos_personas_proyecto, pago_personas_proyecto)
									}
								}else{
									panic(err.Error())
								}
							}
						}else{
							panic(err.Error())
						}
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
	return pagos_personas_proyecto, outputError
}

func CertificacionVistoBueno(dependencia string, mes string, anio string) (personas []models.Persona, outputError map[string]interface{}){
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "CertificacionVistoBueno", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var parametros []models.Parametro
	var vinculaciones_docente []models.VinculacionDocente
	var pagos_mensuales []models.PagoMensual
	var contratistas []models.InformacionProveedor
	var persona models.Persona
	var actasInicio []models.ActaInicio
	var mes_cer, _ = strconv.Atoi(mes)
	var anio_cer, _ = strconv.Atoi(anio)
	var resultado = false

	if err := GetRequestNew("CumplidosDveUrlCrudResoluciones", "vinculacion_docente/?limit=-1&query=ProyectoCurricularId:" + dependencia, &vinculaciones_docente); err == nil{
		for _, vinculacion_docente := range vinculaciones_docente {
			if vinculacion_docente.Activo == true {
				if err := GetRequestLegacy("CumplidosDveUrlCrudAgora", "acta_inicio/?query=NumeroContrato:" + vinculacion_docente.NumeroContrato + ",Vigencia:" + strconv.FormatInt(vinculacion_docente.Vigencia, 10), &actasInicio); err == nil{
					for _, actaInicio := range actasInicio{
						if (actaInicio.FechaInicio.Year() == actaInicio.FechaFin.Year() && int(actaInicio.FechaInicio.Month()) <= mes_cer && actaInicio.FechaInicio.Year() <= anio_cer && int(actaInicio.FechaFin.Month()) >= mes_cer && actaInicio.FechaFin.Year() >= anio_cer) ||
							(actaInicio.FechaInicio.Year() < actaInicio.FechaFin.Year() && int(actaInicio.FechaInicio.Month()) <= mes_cer && actaInicio.FechaInicio.Year() <= anio_cer && int(actaInicio.FechaFin.Month()) <= mes_cer && actaInicio.FechaFin.Year() > anio_cer) ||
							(actaInicio.FechaInicio.Year() < actaInicio.FechaFin.Year() && int(actaInicio.FechaInicio.Month()) >= mes_cer && actaInicio.FechaInicio.Year() < anio_cer && int(actaInicio.FechaFin.Month()) >= mes_cer && actaInicio.FechaFin.Year() >= anio_cer){
								if err := GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion.in:PAD_DVE|AD_DVE|AP_DVE", &parametros); err == nil{
									for _, parametro := range parametros{
										if resultado == false{
											if err := GetRequestNew("CumplidosDveUrlCrud", "pago_mensual/?query=EstadoPagoMensualId:" + strconv.Itoa(parametro.Id) + ",NumeroContrato:" + vinculacion_docente.NumeroContrato + ",VigenciaContrato:" + strconv.FormatInt(vinculacion_docente.Vigencia, 10) + ",Mes:" + mes + ",Ano:" + anio,  &pagos_mensuales); err == nil{
												if pagos_mensuales == nil {
													if err := GetRequestLegacy("CumplidosDveUrlCrudAgora", "informacion_proveedor/?query=NumDocumento:" + strconv.Itoa(vinculacion_docente.PersonaId), &contratistas); err == nil{
														for  _, contratista := range contratistas{
															persona.NumDocumento = contratista.NumDocumento
															persona.Nombre = contratista.NomProveedor
															persona.NumeroContrato = actaInicio.NumeroContrato
															persona.Vigencia = actaInicio.Vigencia
															personas = append(personas, persona)
														}
													}else{
														panic(err.Error())
													}
													resultado = true
												}
											}else{
												panic(err.Error())
											}
										}	
									}
									resultado = false
								}
						}
					}
				}else{
					panic(err.Error())
				}
			}
		}
	}else{
		panic(err.Error())
	}
	return personas, outputError
}

func GenerarPDF(nombre string, proyecto_curricular string, facultad string, mes string, anio string, periodo string) (encodedPdf string, outputError map[string]interface{}){
	defer func(){
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "GenerarPDF", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var pdf *gofpdf.Fpdf
	var err map[string]interface{}

	if pdf, err = ConstruirDocumento(nombre, proyecto_curricular, facultad, mes, anio, periodo); err != nil {
		panic(err)
	}
	if pdf.Err() {
		logs.Error(pdf.Error())
		panic(pdf.Error())
	}
	if pdf.Ok() {
		encodedPdf = encodePDF(pdf)
	}
	return
}

func ConstruirDocumento(nombre string, proyecto_curricular string, facultad string, mes string, anio string, periodo string) (doc *gofpdf.Fpdf, outputError map[string]interface{}){
	defer func(){
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "ConstruirDocumento", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	fontPath := filepath.Join(beego.AppConfig.String("StaticPath"), "fonts")
	imgPath := filepath.Join(beego.AppConfig.String("StaticPath"), "img")
	fontSize := 11.0
	lineHeight := 4.0

	fmt.Println("BANDERA NUMERO UNO")
	
	//DESCIFRAR PROYECTO CURRICULAR 
	proyecto, err := url.QueryUnescape(proyecto_curricular)
	if err != nil {
		fmt.Println("Error al decodificar:", err)
	}

	//DESCIFRAR FACULTAD
	Facultad, err := url.QueryUnescape(facultad)
	if err != nil {
		fmt.Println("Error al decodificar:", err)
	}

	//DESCIFRAR NOMBRE
	Coordinador, err := url.QueryUnescape(nombre)
	if err != nil {
		fmt.Println("Error al decodificar:", err)
	}

	fmt.Println("BANDERA NUMERO DOS")
	//GENERAR FECHA DEL DÍA DE HOY
	now:=time.Now()

	meses := map[time.Month]string{
		time.January:	"Enero",
		time.February:  "Febrero",
        time.March:     "Marzo",
        time.April:     "Abril",
        time.May:       "Mayo",
        time.June:      "Junio",
        time.July:      "Julio",
        time.August:    "Agosto",
        time.September: "Septiembre",
        time.October:   "Octubre",
        time.November:  "Noviembre",
        time.December:  "Diciembre",
	}

	fmt.Println("BANDERA NUMERO TRES")
	pdf := gofpdf.New("P", "mm", "A4", fontPath)
	pdf.AddUTF8Font(Calibri, "", "calibri.ttf")
	pdf.AddUTF8Font(CalibriBold, "B", "calibrib.ttf")
	pdf.AddUTF8Font(MinionProBoldCn, "B", "MinionPro-BoldCn.ttf")
	pdf.AddUTF8Font(MinionProMediumCn, "", "MinionPro-MediumCn.ttf")
	pdf.AddUTF8Font(MinionProBoldItalic, "BI", "MinionProBoldItalic.ttf")

	pdf.SetTopMargin(85)
	fmt.Println("BANDERA NUMERO CUATRO")

	pdf.SetHeaderFuncMode(func() {

		pdf.SetLeftMargin(10)
		pdf.SetRightMargin(10)

		pdf.ImageOptions(filepath.Join(imgPath, "escudo.png"), 82, 8, 45, 45, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
		pdf.SetY(65)
		pdf.SetFont(MinionProBoldCn, "B", fontSize)
		pdf.WriteAligned(0, lineHeight+1, "EL SUSCRITO COORDINADOR DEL PROYECTO CURRICULAR DE " + proyecto + " DE LA " + Facultad + " DE LA UNIVERSIDAD DISTRITAL FRANCISCO JOSÉ DE CALDAS", "C")
		pdf.Ln(lineHeight + 2)
	}, true)

	fmt.Println("BANDERA NUMERO CINCO")
	
	pdf.AliasNbPages("")
	pdf.AddPage()

	pdf.SetAutoPageBreak(false, 25)

	pdf.SetLeftMargin(20)
	pdf.SetRightMargin(20)

	pdf.Ln(lineHeight + 10)

	pdf.SetFont(MinionProBoldCn, "B", fontSize)
	pdf.WriteAligned(0, lineHeight+1, "CERTIFICA QUE:", "C")
	pdf.Ln(lineHeight + 18)

	pdf.SetFont(Calibri, "", fontSize)
	pdf.MultiCell(0, lineHeight+1, "Los Docentes de Vinculación Especial contratados para el periodo Académico " + periodo + ", del Proyecto Curricular de " + proyecto + " cumplieron a cabalidad con las funciones docentes durante el mes de " + mes + " de " + anio + " (según calendario académico).", "", "J", false)
	pdf.Ln(lineHeight * 3)
	
	/*if docentes_incumplidos != nil{
		pdf.WriteAligned(0, lineHeight+1, "A excepción de las siguientes novedades: ", "")
		pdf.Ln(lineHeight * 2)
		for _, docente := range docentes_incumplidos{
			pdf.WriteAligned(0, lineHeight+1, docente.NumDocumento + " " + docente.Nombre +" " + docente.NumeroContrato + ", no se le aprueba cumplido.", "")
			pdf.Ln(lineHeight * 2)
			_, h := pdf.GetPageSize()
			_, _, _, b := pdf.GetMargins()
			if pdf.GetY() > h-b-(lineHeight*10) {
				pdf.AddPage()
			}
		}
	}*/
	
	fmt.Println("BANDERA NUMERO SEIS")
	pdf.Ln(lineHeight * 3)
	pdf.WriteAligned(0, lineHeight+1, "La presente certificación se expide el día " + strconv.Itoa(now.Day()) + " del mes de " + meses[now.Month()] + " de " + strconv.Itoa(now.Year()) + ".", "")
	pdf.Ln(lineHeight * 12)

	pdf.SetFont(MinionProBoldCn, "B", fontSize)
	pdf.WriteAligned(0, lineHeight+1, Coordinador, "C")
	pdf.Ln(lineHeight)
	pdf.WriteAligned(0, lineHeight+1, "Coordinador", "C")
	pdf.Ln(lineHeight)
	pdf.WriteAligned(0, lineHeight+1, "Proyecto Curricular " + proyecto, "C")
	
	fmt.Println("BANDERA FINAL")
	return pdf, outputError
}

func encodePDF(pdf *gofpdf.Fpdf) string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	//pdf.OutputFileAndClose("Certificado.pdf") // para guardar el archivo localmente
	pdf.Output(writer)
	writer.Flush()
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return encodedFile
}

func AprobarMultiplesSolicitudes(v []models.PagoMensual) (resultado string, outputError map[string]interface{}){
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "AprobarMultiplesSolicitudes", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	
	var response interface{}
	for _, pm := range v{
		if err := SendRequestNew("CumplidosDveUrlCrud", "pago_mensual/" + strconv.Itoa(pm.Id), "PUT",&response, &pm); err == nil{
			resultado = "OK"
		}else{
			panic(err.Error())
		}
	}
	return resultado, outputError
}