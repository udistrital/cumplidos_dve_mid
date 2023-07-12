package xray

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecs"
	ecs2 "github.com/aws/aws-xray-sdk-go/awsplugins/ecs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/sirupsen/logrus"
)

type customLogger struct {
	logger *logrus.Logger
}

var globalContext context.Context

func InitXRay() error {
	//cfg := aws.NewConfig().WithRegion("us-west-2") // Reemplaza con tu región deseada
	os.Setenv("AWS_XRAY_NOOP_ID", "false")
	XraySess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})
	if err != nil {
		return err
	}
	//xray.SetLogger(xraylog.NewDefaultLogger(os.Stderr, xraylog.LogLevelError))
	xray.Configure(xray.Config{
		//DaemonAddr: "127.0.0.1:2000", // Dirección y puerto del demonio de X-Ray local
		DaemonAddr: "",
		LogLevel:   "info", // Nivel de log deseado
		LogFormat:  "json", // Formato de log deseado (text o json)
	})

	// S3 and ECS Clients
	ecrClient := ecr.New(XraySess)
	ecsClient := ecs.New(XraySess)

	ecs2.Init()

	// XRay Setup
	xray.AWS(ecrClient.Client)
	xray.AWS(ecsClient.Client)

	fmt.Println("Listed buckets successfully")

	return nil
}

func BeginSegmentWithContext(ctx context.Context, segmentName string, method string, url string, code int, origin string) (context.Context, *xray.Segment) {

	ctx, seg := xray.BeginSegment(ctx, segmentName)
	seg.Origin = origin
	seg.HTTP = &xray.HTTPData{
		Request: &xray.RequestData{
			Method: method,
			URL:    url,
		},
		Response: &xray.ResponseData{
			Status: code,
		},
	}
	return ctx, seg

}

func BeginSubSegmentWithContext(subseg *xray.Segment, method string, url string, code int) {
	subseg.HTTP = &xray.HTTPData{
		Request: &xray.RequestData{
			Method: method, // Método de solicitud
			URL:    url,    // URL de solicitud
		},
		Response: &xray.ResponseData{
			Status: code,
		},
	}
}

func BeginSegmentWithContextTP(ctx context.Context, segmentName string, method string, url string, code int, origin string, traceID []string) (context.Context, *xray.Segment) {

	if traceID != nil {
		traceID := strings.Trim(traceID[0], "[]")
		id, parent := GetTraceIDAndParentID(traceID)
		ctx, seg := xray.BeginSegment(ctx, segmentName)
		seg.Origin = origin
		seg.HTTP = &xray.HTTPData{
			Request: &xray.RequestData{
				Method: method,
				URL:    url,
			},
			Response: &xray.ResponseData{
				Status: code,
			},
		}
		seg.TraceID = id
		seg.ParentID = parent

		return ctx, seg
	} else {
		ctx, seg := xray.BeginSegment(ctx, segmentName)
		seg.Origin = origin
		seg.HTTP = &xray.HTTPData{
			Request: &xray.RequestData{
				Method: method,
				URL:    url,
			},
			Response: &xray.ResponseData{
				Status: code,
			},
		}
		return ctx, seg
	}

}

func GetTraceIDAndParentID(traceID string) (trace string, parent string) {
	if traceID != "" {
		// Separar los pares clave-valor utilizando el delimitador ";"
		traceIDParts := strings.Split(traceID, ";")
		Id := ""
		IdParent := ""
		// Recorrer los pares clave-valor para encontrar el ID de Root
		for _, part := range traceIDParts {
			if strings.HasPrefix(part, "Root=") {
				// Extraer el valor del ID de Root
				Id = strings.TrimPrefix(part, "Root=")
			}
			if strings.HasPrefix(part, "Parent=") {
				// Extraer el valor del ID de Root
				IdParent = strings.TrimPrefix(part, "Parent=")
			}
		}
		return Id, IdParent
	} else {
		return
	}
}

func GetContext() context.Context {
	return globalContext
}

func SetContext(ctx context.Context) {
	globalContext = ctx
}
