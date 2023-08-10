// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 License.

//
// This file provides telemetry support for wrapper functions
//

package telemetry

import (
	"context"
	"regexp"
	"runtime"
	"strings"

	"github.com/Microsoft/go-winio/pkg/etw"
	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/microsoft/moc-sdk-for-go/services/admin/version"
	"github.com/microsoft/moc/pkg/auth"
)

// GUID: {37009198-5B9B-4B40-B1A9-6C561FB69D26}
var providerGUID guid.GUID = guid.FromArray([16]byte{
	0x37, 0x00,
	0x91, 0x98,
	0x5B, 0x9B,
	0x4B, 0x40,
	0xB1, 0xA9,
	0x6C, 0x56,
	0x1F, 0xB6,
	0x9D, 0x26,
})

const ProviderName string = "Microsoft.AKSHCI.MocSdkForGo.Wrapper"
const WrapperTelemetryVersion string = "v1.0.1"

type WrapperEventFields struct {
	EventName         string
	CorrelationVector string
	FunctionName      string
	Subroutine        string
	ErrorString       string
	MocVersion        string
	MocAgentVersion   string
}

// serverName is used when requesting MocVersion info and is NOT to be emitted in telemetry.
func EmitWrapperTelemetry(eventName string, correlationVector string, errorString string, subroutine string, serverName string) {

	// calling for version information is costly, only perform at login
	mocVersion, mocAgentVersion := "", ""
	if eventName == "SecurityLoginCV" {
		mocVersion, mocAgentVersion = getMocVersion(serverName)
	}
	functionName := ""
	var details *runtime.Func
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		details = runtime.FuncForPC(pc)
	}
	if ok && details != nil {
		functionName = details.Name()
	}

	fields := WrapperEventFields{
		eventName,
		correlationVector,
		functionName,
		subroutine,
		FilterSensitiveData(errorString),
		mocVersion,
		mocAgentVersion,
	}
	fields.emitWrapperTelemetry()
}

func (eventTelemetry WrapperEventFields) emitWrapperTelemetry() {

	provider, err := etw.NewProviderWithOptions(ProviderName, etw.WithID(providerGUID))
	if err != nil {
		return
	}
	defer provider.Close()

	err = provider.WriteEvent(
		eventTelemetry.EventName,
		etw.WithEventOpts(),
		etw.WithFields(
			etw.StringField("CorrelationVector", eventTelemetry.CorrelationVector),
			etw.StringField("Function", eventTelemetry.FunctionName),
			etw.StringField("Subroutine", eventTelemetry.Subroutine),
			etw.StringField("Error", eventTelemetry.ErrorString),
			etw.StringField("WrapperTelemetryVersion", WrapperTelemetryVersion),
			etw.StringField("MocVersion", eventTelemetry.MocVersion),
			etw.StringField("MocAgentVersion", eventTelemetry.MocAgentVersion),
		))
	if err != nil {
		return
	}
}

func getMocVersion(serverName string) (string, string) {
	authorizer, err := auth.NewAuthorizerFromEnvironment(serverName)
	if err != nil {
		return "", ""
	}
	client, err := version.NewVersionClient(serverName, authorizer)
	if err != nil {
		return "", ""
	}
	mocVersion, agentVersion, err := client.GetVersion(context.Background())
	if err != nil {
		return "", ""
	}
	return mocVersion, agentVersion
}

// Based on moc-pkg
// This function is used by KeyvaultKeyCreateOrUpdateCV - if a nonempty string is provided, a nonempty string must be returned.
func CheckElementExists(regexStr *regexp.Regexp, inputStr string) string {
	submatchall := regexStr.FindAllString(inputStr, -1)
	for _, toreplace := range submatchall {
		r := strings.NewReplacer(toreplace, "**Redacted**")
		inputStr = r.Replace(inputStr)
	}
	return inputStr
}

// Based on moc-pkg
func FilterSensitiveData(inputStr string) string {
	regexIp := regexp.MustCompile(`\b((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\.)){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\b`)
	regexEmail := regexp.MustCompile("(?i)([A-Z0-9._%+-]+@[A-Z0-9.-]+\\.[A-Z]{2,24})")
	regexFilePaths := regexp.MustCompile(`(:[A-Z]:|\\|(:\.{1,2}[\/\\])+)[\w+\\\s_\(\)\/]+(:\.\w+)*`)

	validatedForIp := CheckElementExists(regexIp, inputStr)
	validatedForEmail := CheckElementExists(regexEmail, validatedForIp)
	validatedAllRegex := CheckElementExists(regexFilePaths, validatedForEmail)

	return validatedAllRegex
}
