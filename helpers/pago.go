package helpers

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/phpdave11/gofpdf"
	"github.com/udistrital/cumplidos_dve_mid/models"
)

func CargarCertificacionDocumentosAprobados(dependencia string, mes string, anio string) (personas []models.Persona, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "CargarCertificacionDocumentosAprobados", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var contrato_ordenador_dependencia models.ContratoOrdenadorDependencia
	var pagos_mensuales []models.PagoMensual
	var persona models.Persona
	var vinculaciones_docente []models.VinculacionDocente
	var mes_cer, _ = strconv.Atoi(mes)
	var parametro []models.Parametro

	if mes_cer < 10 {
		mes = "0" + mes
	}

	contrato_ordenador_dependencia = GetContratosOrdenadorDependencia(dependencia, anio+"-"+mes, anio+"-"+mes)
	if len(contrato_ordenador_dependencia.ContratosOrdenadorDependencia.InformacionContratos) != 0 {
		for _, contrato := range contrato_ordenador_dependencia.ContratosOrdenadorDependencia.InformacionContratos {
			if err := GetRequestNew("CumplidosDveUrlCrudResoluciones", "vinculacion_docente/?limit=-1&query=NumeroContrato:"+contrato.NumeroContrato+",Vigencia:"+contrato.Vigencia, &vinculaciones_docente); err == nil {
				for _, vinculacion_docente := range vinculaciones_docente {
					if vinculacion_docente.NumeroContrato != "" {
						if err := GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion:AP_DVE", &parametro); err == nil {
							if err := GetRequestNew("CumplidosDveUrlCrud", "pago_mensual/?query=EstadoPagoMensualId:"+strconv.Itoa(parametro[0].Id)+",NumeroContrato:"+contrato.NumeroContrato+",VigenciaContrato:"+contrato.Vigencia+",Mes:"+strconv.Itoa(mes_cer)+",Ano:"+anio, &pagos_mensuales); err == nil {
								if len((pagos_mensuales)) == 0 {
									persona.NumDocumento = contrato.Documento
									persona.Nombre = contrato.NombreContratista
									persona.NumeroContrato = contrato.NumeroContrato
									persona.Vigencia, _ = strconv.Atoi(contrato.Vigencia)
									personas = append(personas, persona)
								}
							} else {
								panic(err.Error())
							}
						}
					}
				}
			} else {
				panic(err.Error())
			}
		}
	} else {
		panic("EL ORDENADOR NO TIENE CONTRATOS ASOCIADOS A LA DEPENDENCIA")
	}

	return personas, outputError
}

func GetContratosOrdenadorDependencia(dependencia string, fechaInicio string, fechaFin string) (contratos_ordenador_dependencia models.ContratoOrdenadorDependencia) {
	if err := GetJsonWSO2(beego.AppConfig.String("CumplidosDveUrlWso2")+beego.AppConfig.String("CumplidosDveAdministrativa")+"/contratos_ordenador_dependencia/"+dependencia+"/"+fechaInicio+"/"+fechaFin, &contratos_ordenador_dependencia); err == nil {

	}
	return contratos_ordenador_dependencia
}

func ConsultarPagoAprobado(numero_contrato string, vigencia string, mes string, anio string) (resultado bool, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "ConsultarPagoAprobado", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var pagos_mensuales []models.PagoMensual
	var parametro []models.Parametro

	if err := GetRequestNew("CumplidosDveUrlCrud", "pago_mensual/?query=NumeroContrato:"+numero_contrato+",VigenciaContrato:"+vigencia+",Mes:"+mes+",Ano:"+anio, &pagos_mensuales); err == nil {
		if pagos_mensuales != nil {
			if err := GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion:AP_DVE", &parametro); err == nil {
				for _, pago_mensual := range pagos_mensuales {
					if pago_mensual.EstadoPagoMensualId == parametro[0].Id {
						resultado = true
					} else {
						resultado = false
					}
				}
			}
		} else {
			resultado = false
		}
	} else {
		panic(err.Error())
	}
	return resultado, outputError
}

func CargarSolicitudesOrdenador(doc_ordenador string, limit int, offset int, err0 error) (pagos_personas_proyecto []models.PagoPersonaProyecto, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "CargarSolicitudesOrdenador", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	_ = err0
	var pagos_mensuales []models.PagoMensual
	var contratistas []models.InformacionProveedor
	var pago_personas_proyecto models.PagoPersonaProyecto
	var vinculaciones_docente []models.VinculacionDocente
	var parametro []models.Parametro

	if err := GetRequestNew("CumplidosDveUrlParametros", "parametro/?query=CodigoAbreviacion:PAD_DVE", &parametro); err == nil {
		if err := GetRequestNew("CumplidosDveUrlCrud", "pago_mensual/?limit=-1&query=EstadoPagoMensualId:"+strconv.Itoa(parametro[0].Id)+",Responsable:"+doc_ordenador, &pagos_mensuales); err == nil {
			for x, pago_mensual := range pagos_mensuales {
				if err := GetRequestLegacy("CumplidosDveUrlCrudAgora", "informacion_proveedor/?query=NumDocumento:"+pago_mensual.Persona, &contratistas); err == nil {
					for _, contratista := range contratistas {
						if err := GetRequestNew("CumplidosDveUrlCrudResoluciones", "vinculacion_docente/?limit=-1&query=NumeroContrato:"+pago_mensual.NumeroContrato+",Vigencia:"+strconv.FormatFloat(pago_mensual.VigenciaContrato, 'f', 0, 64), &vinculaciones_docente); err == nil {
							for _, vinculacion := range vinculaciones_docente {
								var dep models.Dependencia

								if err := GetRequestLegacy("CumplidosDveUrlCrudOikos", "dependencia/"+strconv.Itoa(vinculacion.ProyectoCurricularId), &dep); err == nil {
									pago_personas_proyecto.PagoMensual = &pagos_mensuales[x]
									pago_personas_proyecto.NombrePersona = contratista.NomProveedor
									pago_personas_proyecto.Dependencia = &dep
									pagos_personas_proyecto = append(pagos_personas_proyecto, pago_personas_proyecto)
								} else {
									panic(err.Error())
								}
							}
						} else {
							panic(err.Error())
						}
					}
				} else {
					panic(err.Error())
				}
			}
		} else {
			panic(err.Error())
		}
	} else {
		panic(err.Error())
	}
	return pagos_personas_proyecto, outputError
}

func ObtenerDependenciaOrdenador(doc_ordenador string) (resultado int, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "ObtenerDependenciaOrdenador", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var ordenadores_gasto []models.OrdenadorGasto
	var jefes_dependencia []models.JefeDependencia

	if err := GetRequestLegacy("CumplidosDveUrlCore", "jefe_dependencia/?query=TerceroId:"+doc_ordenador+"&sortby=FechaFin&order=desc&limit=1", &jefes_dependencia); err == nil {
		for _, jefe := range jefes_dependencia {
			if err := GetRequestLegacy("CumplidosDveUrlCore", "ordenador_gasto/?query=DependenciaId:"+strconv.Itoa(jefe.DependenciaId), &ordenadores_gasto); err == nil {
				for _, ordenador := range ordenadores_gasto {
					resultado = ordenador.DependenciaId
				}
			} else {
				fmt.Println("Error 1")
				panic(err.Error())
			}
		}
	} else {
		fmt.Println("Error 2")
		panic(err.Error())
	}
	return resultado, outputError
}

func ObtenerInfoCoordinador(DependenciaOikosId int) (info_coordinador models.InformacionCoordinador, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "CargarInformacionCoordinador", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var err error
	var temp map[string]interface{}
	var temp_snies map[string]interface{}

	if err = GetJsonWSO2(beego.AppConfig.String("CumplidosDveUrlWso2")+beego.AppConfig.String("CumplidosDveHomologacion")+"/"+"proyecto_curricular_oikos/"+strconv.Itoa(DependenciaOikosId), &temp); err == nil && temp != nil {
		json_proyecto_curricular, error_json := json.Marshal(temp)
		if error_json == nil {
			var temp_homologacion models.ObjetoProyectoCurricular
			if err = json.Unmarshal(json_proyecto_curricular, &temp_homologacion); err == nil {
				id_proyecto_snies := temp_homologacion.Homologacion.IDSnies
				if err = GetJsonWSO2(beego.AppConfig.String("CumplidosDveUrlWso2")+beego.AppConfig.String("CumplidosDveAcademica")+"/"+"carrera_snies/"+id_proyecto_snies, &temp_snies); err == nil && temp_snies != nil {
					json_info_coordinador, error_json := json.Marshal(temp_snies)
					if error_json == nil {
						var temp_info_coordinador models.InformacionCoordinador
						if err = json.Unmarshal(json_info_coordinador, &temp_info_coordinador); err == nil {
							info_coordinador = temp_info_coordinador
						} else {
							panic(err.Error())
						}
					} else {
						panic(error_json.Error())
					}
				} else {
					panic(err.Error())
				}
			} else {
				panic(err.Error())
			}
		} else {
			panic(error_json.Error())
		}
	} else {
		panic(err.Error())
	}
	return info_coordinador, outputError
}

func ObtenerInfoOrdenador(numero_contrato string, vigencia string) (informacion_ordenador models.InformacionOrdenador, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "ObtenerInfoOrdenador", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var temp map[string]interface{}
	var contrato_elaborado models.ContratoElaborado
	var ordenadores_gasto []models.OrdenadorGasto
	var jefes_dependencia []models.JefeDependencia
	var informacion_proveedores []models.InformacionProveedor
	//var ordenadores []models.Ordenador

	if err := GetJsonWSO2(beego.AppConfig.String("CumplidosDveUrlWso2")+beego.AppConfig.String("CumplidosDveAdministrativa")+"/"+"contrato_elaborado/"+numero_contrato+"/"+vigencia, &temp); err == nil && temp != nil {
		json_contrato_elaborado, error_json := json.Marshal(temp)
		if error_json == nil {
			if err := json.Unmarshal(json_contrato_elaborado, &contrato_elaborado); err == nil {
				if contrato_elaborado.Contrato.TipoContrato == "2" || contrato_elaborado.Contrato.TipoContrato == "3" || contrato_elaborado.Contrato.TipoContrato == "18" {
					if err := GetRequestLegacy("CumplidosDveUrlCore", "ordenador_gasto/?query=Id:"+contrato_elaborado.Contrato.OrdenadorGasto, &ordenadores_gasto); err == nil {
						for _, ordenador_gasto := range ordenadores_gasto {
							if err := GetRequestLegacy("CumplidosDveUrlCore", "jefe_dependencia/?query=DependenciaId:"+strconv.Itoa(ordenador_gasto.DependenciaId)+"&sortby=FechaInicio&order=desc&limit=1", &jefes_dependencia); err == nil {
								for _, jefe_dependencia := range jefes_dependencia {
									if err := GetRequestLegacy("CumplidosDveUrlCrudAgora", "informacion_proveedor/?query=NumDocumento:"+strconv.Itoa(jefe_dependencia.TerceroId), &informacion_proveedores); err == nil {
										for _, informacion_proveedor := range informacion_proveedores {
											informacion_ordenador.NumeroDocumento = jefe_dependencia.TerceroId
											informacion_ordenador.Cargo = ordenador_gasto.Cargo
											informacion_ordenador.Nombre = informacion_proveedor.NomProveedor
											informacion_ordenador.IdDependencia = jefe_dependencia.DependenciaId
										}
									} else {
										panic(err.Error())
									}
								}
							} else {
								panic(err.Error())
							}
						}
					} else {
						panic(err.Error())
					}
				} else {
					panic(err.Error())
				}
			} else {
				panic(err.Error())
			}
		} else {
			panic(error_json.Error())
		}
	} else {
		panic(err.Error())
	}
	return informacion_ordenador, outputError
}

func GenerarPDFOrdenador(nombre string, facultad string, dependencia string, docentes_incumplidos []models.Persona, mes string, anio string, periodo string) (encodedPdf string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "GenerarPDF", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var pdf *gofpdf.Fpdf
	var err map[string]interface{}

	if pdf, err = ConstruirDocumentoOrdenador(nombre, facultad, dependencia, docentes_incumplidos, mes, anio, periodo); err != nil {
		panic(err)
	}
	if pdf.Err() {
		logs.Error(pdf.Error())
		panic(pdf.Error())
	}
	if pdf.Ok() {
		encodedPdf = encodePDFOrdenador(pdf)
	}
	return
}

func ConstruirDocumentoOrdenador(nombre string, facultad string, dependencia string, docentes_incumplidos []models.Persona, mes string, anio string, periodo string) (doc *gofpdf.Fpdf, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "ConstruirDocumento", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	fontPath := filepath.Join(beego.AppConfig.String("StaticPath"), "fonts")
	imgPath := filepath.Join(beego.AppConfig.String("StaticPath"), "img")
	fontSize := 11.0
	lineHeight := 4.0

	//DESCIFRAR DEPENDENCIA
	dependencia_nombre, err := url.QueryUnescape(dependencia)
	if err != nil {
		fmt.Println("Error al decodificar:", err)
	}

	//DESCIFRAR NOMBRE
	decano, err := url.QueryUnescape(nombre)
	if err != nil {
		fmt.Println("Error al decodificar:", err)
	}

	//GENERAR FECHA DEL DÍA DE HOY
	now := time.Now()

	meses := map[time.Month]string{
		time.January:   "Enero",
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

	pdf := gofpdf.New("P", "mm", "A4", fontPath)
	pdf.AddUTF8Font(Calibri, "", "calibri.ttf")
	pdf.AddUTF8Font(CalibriBold, "B", "calibrib.ttf")
	pdf.AddUTF8Font(MinionProBoldCn, "B", "MinionPro-BoldCn.ttf")
	pdf.AddUTF8Font(MinionProMediumCn, "", "MinionPro-MediumCn.ttf")
	pdf.AddUTF8Font(MinionProBoldItalic, "BI", "MinionProBoldItalic.ttf")

	pdf.SetTopMargin(85)

	pdf.SetHeaderFuncMode(func() {

		pdf.SetLeftMargin(10)
		pdf.SetRightMargin(10)

		pdf.ImageOptions(filepath.Join(imgPath, "escudo.png"), 82, 8, 45, 45, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
		pdf.SetY(65)
		pdf.SetFont(MinionProBoldCn, "B", fontSize)
		pdf.MultiCell(0, lineHeight+1, "EL SUSCRITO DECANO/A DE LA "+dependencia_nombre+" DE LA UNIVERSIDAD DISTRITAL FRANCISCO JOSÉ DE CALDAS", "", "C", false)
		pdf.Ln(lineHeight + 2)
	}, true)

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
	pdf.MultiCell(0, lineHeight+1, "De acuerdo a la información suministrada por los proyectos curriculares de la "+dependencia_nombre+", los profesores de Vinculación Especial contratados para el periodo académico "+periodo+", cumplieron a cabalidad con las funciones docentes en el mes de "+mes+" del presente año.(De acuerdo a calendario académico)", "", "J", false)
	pdf.Ln(lineHeight * 3)

	if docentes_incumplidos != nil {
		pdf.WriteAligned(0, lineHeight+1, "A excepción de las siguientes novedades: ", "")
		pdf.Ln(lineHeight * 2)
		for _, docente := range docentes_incumplidos {
			pdf.WriteAligned(0, lineHeight+1, docente.NumDocumento+" "+docente.Nombre+" "+docente.NumeroContrato+", no se le aprueba cumplido.", "")
			pdf.Ln(lineHeight * 2)
			_, h := pdf.GetPageSize()
			_, _, _, b := pdf.GetMargins()
			if pdf.GetY() > h-b-(lineHeight*10) {
				pdf.AddPage()
			}
		}
	}

	pdf.Ln(lineHeight * 3)
	pdf.WriteAligned(0, lineHeight+1, "La presente certificación se expide con destino a la División de Recursos Humanos el día "+strconv.Itoa(now.Day())+" del mes de "+meses[now.Month()]+" de "+strconv.Itoa(now.Year())+".", "")
	pdf.Ln(lineHeight * 12)

	pdf.SetFont(MinionProBoldCn, "B", fontSize)
	pdf.WriteAligned(0, lineHeight+1, decano, "C")
	pdf.Ln(lineHeight)
	pdf.WriteAligned(0, lineHeight+1, "DECANO/A "+dependencia_nombre, "C")

	return pdf, outputError
}

func encodePDFOrdenador(pdf *gofpdf.Fpdf) string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	//pdf.OutputFileAndClose("Certificado.pdf") // para guardar el archivo localmente
	pdf.Output(writer)
	writer.Flush()
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return encodedFile
}

func EnviarTitan(m models.PagoMensual) (resultado models.CumplidoRp, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "EnviarTitan", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var informacion_rp []models.VinculacionDocente
	var cumplido_rp models.CumplidoRp

	if err := GetRequestNew("CumplidosDveUrlCrudResoluciones", "vinculacion_docente/?query=NumeroContrato:"+m.NumeroContrato+",Vigencia:"+strconv.FormatFloat(m.VigenciaContrato, 'f', 0, 64), &informacion_rp); err == nil {
		if err2 := GetRequestNew("CumplidosDveUrlTitan", "contrato_preliquidacion/cumplido_rp/"+strconv.FormatFloat(m.Ano, 'f', 0, 64)+"/"+strconv.FormatFloat(m.Mes, 'f', 0, 64)+"/"+m.NumeroContrato+"/"+strconv.FormatFloat(informacion_rp[0].VigenciaRp, 'f', 0, 64)+"/"+strconv.FormatFloat(informacion_rp[0].NumeroRp, 'f', 0, 64), &cumplido_rp); err2 == nil {
			resultado = cumplido_rp
		} else {
			panic(err2.Error())
		}
	} else {
		panic(err.Error())
	}

	return resultado, outputError
}

func AprobarMultiplesPagos(m []models.PagoMensual) (resultado string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "AprobarMultiplesPagos", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var response interface{}
	var informacion_rp []models.VinculacionDocente
	var cumplido_rp models.CumplidoRp
	//var pagos_mensuales []*models.PagoMensual
	//var pago_mensual *models.PagoMensual

	for _, pm := range m {
		//pago_mensual = pm.PagoMensual
		//pagos_mensuales = append(pagos_mensuales, pago_mensual)
		if err := GetRequestNew("CumplidosDveUrlCrudResoluciones", "vinculacion_docente/?query=NumeroContrato:"+pm.NumeroContrato+",Vigencia:"+strconv.FormatFloat(pm.VigenciaContrato, 'f', 0, 64), &informacion_rp); err == nil {
			if err2 := GetRequestNew("CumplidosDveUrlTitan", "contrato_preliquidacion/cumplido_rp/"+strconv.FormatFloat(pm.Ano, 'f', 0, 64)+"/"+strconv.FormatFloat(pm.Mes, 'f', 0, 64)+"/"+pm.NumeroContrato+"/"+strconv.FormatFloat(informacion_rp[0].VigenciaRp, 'f', 0, 64)+"/"+strconv.FormatFloat(informacion_rp[0].NumeroRp, 'f', 0, 64), &cumplido_rp); err2 == nil {
				if cumplido_rp.Id != 0 {
					if err := SendRequestNew("CumplidosDveUrlCrud", "pago_mensual/"+strconv.Itoa(pm.Id), "PUT", &response, &pm); err == nil {
						resultado = "OK"
					} else {
						panic(err.Error())
					}
				} else {
					panic("Error al envíar a titan, no existe el cumplido")
				}
			} else {
				panic(err2.Error())
			}
		} else {
			panic(err.Error())
		}
	}
	return resultado, outputError
}
