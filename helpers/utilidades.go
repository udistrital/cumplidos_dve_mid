package helpers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

const (
	Calibri             string = "Calibri"
	CalibriBold         string = "Calibri-Bold"
	MinionProBoldCn     string = "MinionPro-BoldCn"
	MinionProMediumCn   string = "MinionPro-MediumCn"
	MinionProBoldItalic string = "MinionProBoldItalic"
)

func GetRequestWSO2(service string, route string, target interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") +
		beego.AppConfig.String("CumplidosDveUrlWso2") +
		beego.AppConfig.String(service) + "/" + route
	if response, err := request.GetJsonWSO2Test(url, &target); response == 200 && err == nil {
		return nil
	} else {
		return err
	}
}
