package cheerapp

import (
	"context"
	"fmt"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/chwjbn/go4sky"
	"github.com/chwjbn/go4sky/reporter"
	"os"
	"time"

	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

var(
	gSkyapmOapGrpcAddr ="127.0.0.1:11800"
	gSkyapmTracer *go4sky.Tracer=nil
	gSkyapmLogger *go4sky.SkyLogger
)


func StartSkyapm()  {

	cheerlib.StdInfo("cheerapp.initSkyapm Begin")
	defer func() {
		go4sky.SetGlobalTracer(gSkyapmTracer)
		cheerlib.StdInfo("cheerapp.initSkyapm End")
	}()

	var xSkyapmError error

	xSkyapmAppName:=os.Getenv("SKYAPM_APP_NAME")
	if len(xSkyapmAppName)<1{
		xSkyapmAppName=fmt.Sprintf("arch::%s",cheerlib.ApplicationFileName())
	}

	cheerlib.StdInfo(fmt.Sprintf("cheerapp.initSkyapm SkyapmAppName=[%s]",xSkyapmAppName))

	gSkyapmTracer, xSkyapmError =go4sky.NewTracer(xSkyapmAppName)
	if xSkyapmError!=nil{
		cheerlib.StdError(fmt.Sprintf("cheerapp.initSkyapm go4sky.NewTracer Error=[%s]",xSkyapmError.Error()))
		return
	}

	xSkyapmOapGrpcAddr:=os.Getenv("SKYAPM_OAP_GRPC_ADDR")

	cheerlib.StdInfo(fmt.Sprintf("cheerapp.initSkyapm SkyapmOapGrpcAddr=[%s]",xSkyapmOapGrpcAddr))

	if len(xSkyapmOapGrpcAddr)>0{
		gSkyapmOapGrpcAddr=xSkyapmOapGrpcAddr

		var xSkyapmReporter go4sky.Reporter
		xSkyapmReporter,xSkyapmError=reporter.NewGRPCReporter(xSkyapmOapGrpcAddr, reporter.WithCheckInterval(time.Second))
		if xSkyapmError!=nil{
			cheerlib.StdError(fmt.Sprintf("cheerapp.initSkyapm reporter.NewGRPCReporter Error=[%s]",xSkyapmError.Error()))
			return
		}

		var xSkyapmTracer *go4sky.Tracer
		xSkyapmTracer, xSkyapmError =go4sky.NewTracer(xSkyapmAppName,go4sky.WithReporter(xSkyapmReporter))
		if xSkyapmError!=nil{
			cheerlib.StdError(fmt.Sprintf("cheerapp.initSkyapm go4sky.NewTracer Error=[%s]",xSkyapmError.Error()))
			return
		}

		gSkyapmTracer=xSkyapmTracer

		var xSkyapmLogger *go4sky.SkyLogger
		xSkyapmError,xSkyapmLogger=go4sky.NewSkyLogger(xSkyapmReporter)
		if xSkyapmError!=nil{
			cheerlib.StdError(fmt.Sprintf("cheerapp.initSkyapm go4sky.NewSkyLogger Error=[%s]",xSkyapmError.Error()))
			return
		}

		gSkyapmLogger=xSkyapmLogger
	}


}

func SpanEnd(span go4sky.Span)  {
	if span==nil{
		return
	}
	span.End()
}

func SpanTag(span go4sky.Span,tagKey go4sky.Tag,tagVal string)  {
	if span==nil{
		return
	}

	span.Tag(tagKey,tagVal)
}

func SpanError(span go4sky.Span,ll ...string)  {
	if span==nil{
		return
	}

	span.Error(time.Now(),ll...)
}

func SpanLog(span go4sky.Span,ll ...string)  {
	if span==nil{
		return
	}

	span.Log(time.Now(),ll...)
}

func SpanBeginDbService(ctx context.Context,operName string) go4sky.Span  {

	var xSpan go4sky.Span

	xSkyapmTracer:=go4sky.GetGlobalTracer()
	if xSkyapmTracer==nil{
		return xSpan
	}

	xSpan, _, _ = xSkyapmTracer.CreateLocalSpan(ctx)
	if xSpan==nil{
		return xSpan
	}

	xSpan.SetComponent(42)
	xSpan.SetOperationName(operName)
	xSpan.SetPeer(fmt.Sprintf("%s@%s",cheerlib.OsIPV4(),cheerlib.OsHostName()))
	xSpan.SetSpanLayer(agentv3.SpanLayer_Database)

	return xSpan

}

func SpanBeginBizFunction(ctx context.Context,operName string) go4sky.Span  {

	var xSpan go4sky.Span

	xSkyapmTracer:=go4sky.GetGlobalTracer()
	if xSkyapmTracer==nil{
		return xSpan
	}

	xSpan, _, _ = xSkyapmTracer.CreateLocalSpan(ctx)
	if xSpan==nil{
		return xSpan
	}

	xSpan.SetComponent(0)
	xSpan.SetOperationName(operName)
	xSpan.SetPeer(fmt.Sprintf("%s@%s",cheerlib.OsIPV4(),cheerlib.OsHostName()))
	xSpan.SetSpanLayer(agentv3.SpanLayer_RPCFramework)

	return xSpan

}