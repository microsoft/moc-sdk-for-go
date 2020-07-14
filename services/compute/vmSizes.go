// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package compute

import (
	cloudcompute "github.com/microsoft/moc-proto/rpc/common"
)

func GetCloudSdkVirtualMachineSizeFromCloudVirtualMachineSize(size cloudcompute.VirtualMachineSizeType) VirtualMachineSizeTypes {
	sizeInt := int32(size)
	value, found := cloudcompute.VirtualMachineSizeType_name[sizeInt]
	if !found {
		return VirtualMachineSizeTypesDefault // Not found, return default
	}
	return VirtualMachineSizeTypes(value)
}

func GetCloudVirtualMachineSizeFromCloudSdkVirtualMachineSize(size VirtualMachineSizeTypes) cloudcompute.VirtualMachineSizeType {
	// Convert sdk enum to string representation
	sizeString := string(size)

	// Find the corresponding string in size map
	value, found := cloudcompute.VirtualMachineSizeType_value[sizeString]
	if !found {
		// Not found, user supplied unsupported size
		return cloudcompute.VirtualMachineSizeType_Unsupported
	}
	return cloudcompute.VirtualMachineSizeType(value)
}

// VirtualMachineSizeTypes enumerates the values for virtual machine size types.
type VirtualMachineSizeTypes string

// For more information about virtual machine sizes, see 'Sizes for virtual machines':
//  https://docs.microsoft.com/en-us/azure/virtual-machines/windows/sizes
//  NOTE: Kubernetes requires 2 CPU cores. [ERROR NumCPU]: the number of available CPUs 1 is less than the required 2.
/*
The following Size Types are supported:

				CPU  GBRAM
Default            4    4
Standard_A2_v2     2    4
Standard_A4_v2     4    8
-
Standard_D2s_v3    2    8
Standard_D4s_v3    4   16
Standard_D8s_v3    8   32
Standard_D16s_v3  16   64
Standard_D32s_v3  32  128
-
Standard_DS2_v2    2    7
Standard_DS3_v2    2   14
Standard_DS4_v2    8   28
Standard_DS5_v2   16   56
-
Standard_DS13_v2   8   56
-
Standard_K8S_v1    4    2 (custom for IoT)
Standard_K8S2_v1   2    2 (custom for IoT)
Standard_K8S3_v1   4    6 (custom for WAC)
Standard_K8S4_v1   4    4 (WSSD Default size)
*/

const (
	// VirtualMachineSizeTypesDefault ...
	VirtualMachineSizeTypesDefault VirtualMachineSizeTypes = "Default"
	// VirtualMachineSizeTypesStandardK8SV1 ...
	VirtualMachineSizeTypesStandardK8SV1 VirtualMachineSizeTypes = "Standard_K8S_v1"
	// VirtualMachineSizeTypesStandardK8S2V1 ...
	VirtualMachineSizeTypesStandardK8S2V1 VirtualMachineSizeTypes = "Standard_K8S2_v1"
	// VirtualMachineSizeTypesStandardK8S3V1 ...
	VirtualMachineSizeTypesStandardK8S3V1 VirtualMachineSizeTypes = "Standard_K8S3_v1"
	// VirtualMachineSizeTypesStandardK8S4V1 ...
	VirtualMachineSizeTypesStandardK8S4V1 VirtualMachineSizeTypes = "Standard_K8S4_v1"
	// VirtualMachineSizeTypesBasicA0 ...
	VirtualMachineSizeTypesBasicA0 VirtualMachineSizeTypes = "Basic_A0"
	// VirtualMachineSizeTypesBasicA1 ...
	VirtualMachineSizeTypesBasicA1 VirtualMachineSizeTypes = "Basic_A1"
	// VirtualMachineSizeTypesBasicA2 ...
	VirtualMachineSizeTypesBasicA2 VirtualMachineSizeTypes = "Basic_A2"
	// VirtualMachineSizeTypesBasicA3 ...
	VirtualMachineSizeTypesBasicA3 VirtualMachineSizeTypes = "Basic_A3"
	// VirtualMachineSizeTypesBasicA4 ...
	VirtualMachineSizeTypesBasicA4 VirtualMachineSizeTypes = "Basic_A4"
	// VirtualMachineSizeTypesStandardA0 ...
	VirtualMachineSizeTypesStandardA0 VirtualMachineSizeTypes = "Standard_A0"
	// VirtualMachineSizeTypesStandardA1 ...
	VirtualMachineSizeTypesStandardA1 VirtualMachineSizeTypes = "Standard_A1"
	// VirtualMachineSizeTypesStandardA10 ...
	VirtualMachineSizeTypesStandardA10 VirtualMachineSizeTypes = "Standard_A10"
	// VirtualMachineSizeTypesStandardA11 ...
	VirtualMachineSizeTypesStandardA11 VirtualMachineSizeTypes = "Standard_A11"
	// VirtualMachineSizeTypesStandardA1V2 ...
	VirtualMachineSizeTypesStandardA1V2 VirtualMachineSizeTypes = "Standard_A1_v2"
	// VirtualMachineSizeTypesStandardA2 ...
	VirtualMachineSizeTypesStandardA2 VirtualMachineSizeTypes = "Standard_A2"
	// VirtualMachineSizeTypesStandardA2mV2 ...
	VirtualMachineSizeTypesStandardA2mV2 VirtualMachineSizeTypes = "Standard_A2m_v2"
	// VirtualMachineSizeTypesStandardA2V2 ...
	VirtualMachineSizeTypesStandardA2V2 VirtualMachineSizeTypes = "Standard_A2_v2"
	// VirtualMachineSizeTypesStandardA3 ...
	VirtualMachineSizeTypesStandardA3 VirtualMachineSizeTypes = "Standard_A3"
	// VirtualMachineSizeTypesStandardA4 ...
	VirtualMachineSizeTypesStandardA4 VirtualMachineSizeTypes = "Standard_A4"
	// VirtualMachineSizeTypesStandardA4mV2 ...
	VirtualMachineSizeTypesStandardA4mV2 VirtualMachineSizeTypes = "Standard_A4m_v2"
	// VirtualMachineSizeTypesStandardA4V2 ...
	VirtualMachineSizeTypesStandardA4V2 VirtualMachineSizeTypes = "Standard_A4_v2"
	// VirtualMachineSizeTypesStandardA5 ...
	VirtualMachineSizeTypesStandardA5 VirtualMachineSizeTypes = "Standard_A5"
	// VirtualMachineSizeTypesStandardA6 ...
	VirtualMachineSizeTypesStandardA6 VirtualMachineSizeTypes = "Standard_A6"
	// VirtualMachineSizeTypesStandardA7 ...
	VirtualMachineSizeTypesStandardA7 VirtualMachineSizeTypes = "Standard_A7"
	// VirtualMachineSizeTypesStandardA8 ...
	VirtualMachineSizeTypesStandardA8 VirtualMachineSizeTypes = "Standard_A8"
	// VirtualMachineSizeTypesStandardA8mV2 ...
	VirtualMachineSizeTypesStandardA8mV2 VirtualMachineSizeTypes = "Standard_A8m_v2"
	// VirtualMachineSizeTypesStandardA8V2 ...
	VirtualMachineSizeTypesStandardA8V2 VirtualMachineSizeTypes = "Standard_A8_v2"
	// VirtualMachineSizeTypesStandardA9 ...
	VirtualMachineSizeTypesStandardA9 VirtualMachineSizeTypes = "Standard_A9"
	// VirtualMachineSizeTypesStandardB1ms ...
	VirtualMachineSizeTypesStandardB1ms VirtualMachineSizeTypes = "Standard_B1ms"
	// VirtualMachineSizeTypesStandardB1s ...
	VirtualMachineSizeTypesStandardB1s VirtualMachineSizeTypes = "Standard_B1s"
	// VirtualMachineSizeTypesStandardB2ms ...
	VirtualMachineSizeTypesStandardB2ms VirtualMachineSizeTypes = "Standard_B2ms"
	// VirtualMachineSizeTypesStandardB2s ...
	VirtualMachineSizeTypesStandardB2s VirtualMachineSizeTypes = "Standard_B2s"
	// VirtualMachineSizeTypesStandardB4ms ...
	VirtualMachineSizeTypesStandardB4ms VirtualMachineSizeTypes = "Standard_B4ms"
	// VirtualMachineSizeTypesStandardB8ms ...
	VirtualMachineSizeTypesStandardB8ms VirtualMachineSizeTypes = "Standard_B8ms"
	// VirtualMachineSizeTypesStandardDV2 ...
	VirtualMachineSizeTypesStandardDV2 VirtualMachineSizeTypes = "Standard_D_v2"
	// VirtualMachineSizeTypesStandardDV3 ...
	VirtualMachineSizeTypesStandardDV3 VirtualMachineSizeTypes = "Standard_D_v3"
	// VirtualMachineSizeTypesStandardD1 ...
	VirtualMachineSizeTypesStandardD1 VirtualMachineSizeTypes = "Standard_D1"
	// VirtualMachineSizeTypesStandardD11 ...
	VirtualMachineSizeTypesStandardD11 VirtualMachineSizeTypes = "Standard_D11"
	// VirtualMachineSizeTypesStandardD11V2 ...
	VirtualMachineSizeTypesStandardD11V2 VirtualMachineSizeTypes = "Standard_D11_v2"
	// VirtualMachineSizeTypesStandardD12 ...
	VirtualMachineSizeTypesStandardD12 VirtualMachineSizeTypes = "Standard_D12"
	// VirtualMachineSizeTypesStandardD12V2 ...
	VirtualMachineSizeTypesStandardD12V2 VirtualMachineSizeTypes = "Standard_D12_v2"
	// VirtualMachineSizeTypesStandardD13 ...
	VirtualMachineSizeTypesStandardD13 VirtualMachineSizeTypes = "Standard_D13"
	// VirtualMachineSizeTypesStandardD13V2 ...
	VirtualMachineSizeTypesStandardD13V2 VirtualMachineSizeTypes = "Standard_D13_v2"
	// VirtualMachineSizeTypesStandardD14 ...
	VirtualMachineSizeTypesStandardD14 VirtualMachineSizeTypes = "Standard_D14"
	// VirtualMachineSizeTypesStandardD14V2 ...
	VirtualMachineSizeTypesStandardD14V2 VirtualMachineSizeTypes = "Standard_D14_v2"
	// VirtualMachineSizeTypesStandardD15V2 ...
	VirtualMachineSizeTypesStandardD15V2 VirtualMachineSizeTypes = "Standard_D15_v2"
	// VirtualMachineSizeTypesStandardD16sV3 ...
	VirtualMachineSizeTypesStandardD16sV3 VirtualMachineSizeTypes = "Standard_D16s_v3"
	// VirtualMachineSizeTypesStandardD16V3 ...
	VirtualMachineSizeTypesStandardD16V3 VirtualMachineSizeTypes = "Standard_D16_v3"
	// VirtualMachineSizeTypesStandardD1V2 ...
	VirtualMachineSizeTypesStandardD1V2 VirtualMachineSizeTypes = "Standard_D1_v2"
	// VirtualMachineSizeTypesStandardD2 ...
	VirtualMachineSizeTypesStandardD2 VirtualMachineSizeTypes = "Standard_D2"
	// VirtualMachineSizeTypesStandardD2sV3 ...
	VirtualMachineSizeTypesStandardD2sV3 VirtualMachineSizeTypes = "Standard_D2s_v3"
	// VirtualMachineSizeTypesStandardD2V2 ...
	VirtualMachineSizeTypesStandardD2V2 VirtualMachineSizeTypes = "Standard_D2_v2"
	// VirtualMachineSizeTypesStandardD2V3 ...
	VirtualMachineSizeTypesStandardD2V3 VirtualMachineSizeTypes = "Standard_D2_v3"
	// VirtualMachineSizeTypesStandardD3 ...
	VirtualMachineSizeTypesStandardD3 VirtualMachineSizeTypes = "Standard_D3"
	// VirtualMachineSizeTypesStandardD32sV3 ...
	VirtualMachineSizeTypesStandardD32sV3 VirtualMachineSizeTypes = "Standard_D32s_v3"
	// VirtualMachineSizeTypesStandardD32V3 ...
	VirtualMachineSizeTypesStandardD32V3 VirtualMachineSizeTypes = "Standard_D32_v3"
	// VirtualMachineSizeTypesStandardD3V2 ...
	VirtualMachineSizeTypesStandardD3V2 VirtualMachineSizeTypes = "Standard_D3_v2"
	// VirtualMachineSizeTypesStandardD4 ...
	VirtualMachineSizeTypesStandardD4 VirtualMachineSizeTypes = "Standard_D4"
	// VirtualMachineSizeTypesStandardD4sV3 ...
	VirtualMachineSizeTypesStandardD4sV3 VirtualMachineSizeTypes = "Standard_D4s_v3"
	// VirtualMachineSizeTypesStandardD4V2 ...
	VirtualMachineSizeTypesStandardD4V2 VirtualMachineSizeTypes = "Standard_D4_v2"
	// VirtualMachineSizeTypesStandardD4V3 ...
	VirtualMachineSizeTypesStandardD4V3 VirtualMachineSizeTypes = "Standard_D4_v3"
	// VirtualMachineSizeTypesStandardD5V2 ...
	VirtualMachineSizeTypesStandardD5V2 VirtualMachineSizeTypes = "Standard_D5_v2"
	// VirtualMachineSizeTypesStandardD64sV3 ...
	VirtualMachineSizeTypesStandardD64sV3 VirtualMachineSizeTypes = "Standard_D64s_v3"
	// VirtualMachineSizeTypesStandardD64V3 ...
	VirtualMachineSizeTypesStandardD64V3 VirtualMachineSizeTypes = "Standard_D64_v3"
	// VirtualMachineSizeTypesStandardD8sV3 ...
	VirtualMachineSizeTypesStandardD8sV3 VirtualMachineSizeTypes = "Standard_D8s_v3"
	// VirtualMachineSizeTypesStandardD8V3 ...
	VirtualMachineSizeTypesStandardD8V3 VirtualMachineSizeTypes = "Standard_D8_v3"
	// VirtualMachineSizeTypesStandardDS1 ...
	VirtualMachineSizeTypesStandardDS1 VirtualMachineSizeTypes = "Standard_DS1"
	// VirtualMachineSizeTypesStandardDS11 ...
	VirtualMachineSizeTypesStandardDS11 VirtualMachineSizeTypes = "Standard_DS11"
	// VirtualMachineSizeTypesStandardDS11V2 ...
	VirtualMachineSizeTypesStandardDS11V2 VirtualMachineSizeTypes = "Standard_DS11_v2"
	// VirtualMachineSizeTypesStandardDS12 ...
	VirtualMachineSizeTypesStandardDS12 VirtualMachineSizeTypes = "Standard_DS12"
	// VirtualMachineSizeTypesStandardDS12V2 ...
	VirtualMachineSizeTypesStandardDS12V2 VirtualMachineSizeTypes = "Standard_DS12_v2"
	// VirtualMachineSizeTypesStandardDS13 ...
	VirtualMachineSizeTypesStandardDS13 VirtualMachineSizeTypes = "Standard_DS13"
	// VirtualMachineSizeTypesStandardDS132V2 ...
	VirtualMachineSizeTypesStandardDS132V2 VirtualMachineSizeTypes = "Standard_DS13-2_v2"
	// VirtualMachineSizeTypesStandardDS134V2 ...
	VirtualMachineSizeTypesStandardDS134V2 VirtualMachineSizeTypes = "Standard_DS13-4_v2"
	// VirtualMachineSizeTypesStandardDS13V2 ...
	VirtualMachineSizeTypesStandardDS13V2 VirtualMachineSizeTypes = "Standard_DS13_v2"
	// VirtualMachineSizeTypesStandardDS14 ...
	VirtualMachineSizeTypesStandardDS14 VirtualMachineSizeTypes = "Standard_DS14"
	// VirtualMachineSizeTypesStandardDS144V2 ...
	VirtualMachineSizeTypesStandardDS144V2 VirtualMachineSizeTypes = "Standard_DS14-4_v2"
	// VirtualMachineSizeTypesStandardDS148V2 ...
	VirtualMachineSizeTypesStandardDS148V2 VirtualMachineSizeTypes = "Standard_DS14-8_v2"
	// VirtualMachineSizeTypesStandardDS14V2 ...
	VirtualMachineSizeTypesStandardDS14V2 VirtualMachineSizeTypes = "Standard_DS14_v2"
	// VirtualMachineSizeTypesStandardDS15V2 ...
	VirtualMachineSizeTypesStandardDS15V2 VirtualMachineSizeTypes = "Standard_DS15_v2"
	// VirtualMachineSizeTypesStandardDS1V2 ...
	VirtualMachineSizeTypesStandardDS1V2 VirtualMachineSizeTypes = "Standard_DS1_v2"
	// VirtualMachineSizeTypesStandardDS2 ...
	VirtualMachineSizeTypesStandardDS2 VirtualMachineSizeTypes = "Standard_DS2"
	// VirtualMachineSizeTypesStandardDS2V2 ...
	VirtualMachineSizeTypesStandardDS2V2 VirtualMachineSizeTypes = "Standard_DS2_v2"
	// VirtualMachineSizeTypesStandardDS3 ...
	VirtualMachineSizeTypesStandardDS3 VirtualMachineSizeTypes = "Standard_DS3"
	// VirtualMachineSizeTypesStandardDS3V2 ...
	VirtualMachineSizeTypesStandardDS3V2 VirtualMachineSizeTypes = "Standard_DS3_v2"
	// VirtualMachineSizeTypesStandardDS4 ...
	VirtualMachineSizeTypesStandardDS4 VirtualMachineSizeTypes = "Standard_DS4"
	// VirtualMachineSizeTypesStandardDS4V2 ...
	VirtualMachineSizeTypesStandardDS4V2 VirtualMachineSizeTypes = "Standard_DS4_v2"
	// VirtualMachineSizeTypesStandardDS5V2 ...
	VirtualMachineSizeTypesStandardDS5V2 VirtualMachineSizeTypes = "Standard_DS5_v2"
	// VirtualMachineSizeTypesStandardE16sV3 ...
	VirtualMachineSizeTypesStandardE16sV3 VirtualMachineSizeTypes = "Standard_E16s_v3"
	// VirtualMachineSizeTypesStandardE16V3 ...
	VirtualMachineSizeTypesStandardE16V3 VirtualMachineSizeTypes = "Standard_E16_v3"
	// VirtualMachineSizeTypesStandardE2sV3 ...
	VirtualMachineSizeTypesStandardE2sV3 VirtualMachineSizeTypes = "Standard_E2s_v3"
	// VirtualMachineSizeTypesStandardE2V3 ...
	VirtualMachineSizeTypesStandardE2V3 VirtualMachineSizeTypes = "Standard_E2_v3"
	// VirtualMachineSizeTypesStandardE3216V3 ...
	VirtualMachineSizeTypesStandardE3216V3 VirtualMachineSizeTypes = "Standard_E32-16_v3"
	// VirtualMachineSizeTypesStandardE328sV3 ...
	VirtualMachineSizeTypesStandardE328sV3 VirtualMachineSizeTypes = "Standard_E32-8s_v3"
	// VirtualMachineSizeTypesStandardE32sV3 ...
	VirtualMachineSizeTypesStandardE32sV3 VirtualMachineSizeTypes = "Standard_E32s_v3"
	// VirtualMachineSizeTypesStandardE32V3 ...
	VirtualMachineSizeTypesStandardE32V3 VirtualMachineSizeTypes = "Standard_E32_v3"
	// VirtualMachineSizeTypesStandardE4sV3 ...
	VirtualMachineSizeTypesStandardE4sV3 VirtualMachineSizeTypes = "Standard_E4s_v3"
	// VirtualMachineSizeTypesStandardE4V3 ...
	VirtualMachineSizeTypesStandardE4V3 VirtualMachineSizeTypes = "Standard_E4_v3"
	// VirtualMachineSizeTypesStandardE6416sV3 ...
	VirtualMachineSizeTypesStandardE6416sV3 VirtualMachineSizeTypes = "Standard_E64-16s_v3"
	// VirtualMachineSizeTypesStandardE6432sV3 ...
	VirtualMachineSizeTypesStandardE6432sV3 VirtualMachineSizeTypes = "Standard_E64-32s_v3"
	// VirtualMachineSizeTypesStandardE64sV3 ...
	VirtualMachineSizeTypesStandardE64sV3 VirtualMachineSizeTypes = "Standard_E64s_v3"
	// VirtualMachineSizeTypesStandardE64V3 ...
	VirtualMachineSizeTypesStandardE64V3 VirtualMachineSizeTypes = "Standard_E64_v3"
	// VirtualMachineSizeTypesStandardE8sV3 ...
	VirtualMachineSizeTypesStandardE8sV3 VirtualMachineSizeTypes = "Standard_E8s_v3"
	// VirtualMachineSizeTypesStandardE8V3 ...
	VirtualMachineSizeTypesStandardE8V3 VirtualMachineSizeTypes = "Standard_E8_v3"
	// VirtualMachineSizeTypesStandardF1 ...
	VirtualMachineSizeTypesStandardF1 VirtualMachineSizeTypes = "Standard_F1"
	// VirtualMachineSizeTypesStandardF16 ...
	VirtualMachineSizeTypesStandardF16 VirtualMachineSizeTypes = "Standard_F16"
	// VirtualMachineSizeTypesStandardF16s ...
	VirtualMachineSizeTypesStandardF16s VirtualMachineSizeTypes = "Standard_F16s"
	// VirtualMachineSizeTypesStandardF16sV2 ...
	VirtualMachineSizeTypesStandardF16sV2 VirtualMachineSizeTypes = "Standard_F16s_v2"
	// VirtualMachineSizeTypesStandardF1s ...
	VirtualMachineSizeTypesStandardF1s VirtualMachineSizeTypes = "Standard_F1s"
	// VirtualMachineSizeTypesStandardF2 ...
	VirtualMachineSizeTypesStandardF2 VirtualMachineSizeTypes = "Standard_F2"
	// VirtualMachineSizeTypesStandardF2s ...
	VirtualMachineSizeTypesStandardF2s VirtualMachineSizeTypes = "Standard_F2s"
	// VirtualMachineSizeTypesStandardF2sV2 ...
	VirtualMachineSizeTypesStandardF2sV2 VirtualMachineSizeTypes = "Standard_F2s_v2"
	// VirtualMachineSizeTypesStandardF32sV2 ...
	VirtualMachineSizeTypesStandardF32sV2 VirtualMachineSizeTypes = "Standard_F32s_v2"
	// VirtualMachineSizeTypesStandardF4 ...
	VirtualMachineSizeTypesStandardF4 VirtualMachineSizeTypes = "Standard_F4"
	// VirtualMachineSizeTypesStandardF4s ...
	VirtualMachineSizeTypesStandardF4s VirtualMachineSizeTypes = "Standard_F4s"
	// VirtualMachineSizeTypesStandardF4sV2 ...
	VirtualMachineSizeTypesStandardF4sV2 VirtualMachineSizeTypes = "Standard_F4s_v2"
	// VirtualMachineSizeTypesStandardF64sV2 ...
	VirtualMachineSizeTypesStandardF64sV2 VirtualMachineSizeTypes = "Standard_F64s_v2"
	// VirtualMachineSizeTypesStandardF72sV2 ...
	VirtualMachineSizeTypesStandardF72sV2 VirtualMachineSizeTypes = "Standard_F72s_v2"
	// VirtualMachineSizeTypesStandardF8 ...
	VirtualMachineSizeTypesStandardF8 VirtualMachineSizeTypes = "Standard_F8"
	// VirtualMachineSizeTypesStandardF8s ...
	VirtualMachineSizeTypesStandardF8s VirtualMachineSizeTypes = "Standard_F8s"
	// VirtualMachineSizeTypesStandardF8sV2 ...
	VirtualMachineSizeTypesStandardF8sV2 VirtualMachineSizeTypes = "Standard_F8s_v2"
	// VirtualMachineSizeTypesStandardG1 ...
	VirtualMachineSizeTypesStandardG1 VirtualMachineSizeTypes = "Standard_G1"
	// VirtualMachineSizeTypesStandardG2 ...
	VirtualMachineSizeTypesStandardG2 VirtualMachineSizeTypes = "Standard_G2"
	// VirtualMachineSizeTypesStandardG3 ...
	VirtualMachineSizeTypesStandardG3 VirtualMachineSizeTypes = "Standard_G3"
	// VirtualMachineSizeTypesStandardG4 ...
	VirtualMachineSizeTypesStandardG4 VirtualMachineSizeTypes = "Standard_G4"
	// VirtualMachineSizeTypesStandardG5 ...
	VirtualMachineSizeTypesStandardG5 VirtualMachineSizeTypes = "Standard_G5"
	// VirtualMachineSizeTypesStandardGS1 ...
	VirtualMachineSizeTypesStandardGS1 VirtualMachineSizeTypes = "Standard_GS1"
	// VirtualMachineSizeTypesStandardGS2 ...
	VirtualMachineSizeTypesStandardGS2 VirtualMachineSizeTypes = "Standard_GS2"
	// VirtualMachineSizeTypesStandardGS3 ...
	VirtualMachineSizeTypesStandardGS3 VirtualMachineSizeTypes = "Standard_GS3"
	// VirtualMachineSizeTypesStandardGS4 ...
	VirtualMachineSizeTypesStandardGS4 VirtualMachineSizeTypes = "Standard_GS4"
	// VirtualMachineSizeTypesStandardGS44 ...
	VirtualMachineSizeTypesStandardGS44 VirtualMachineSizeTypes = "Standard_GS4-4"
	// VirtualMachineSizeTypesStandardGS48 ...
	VirtualMachineSizeTypesStandardGS48 VirtualMachineSizeTypes = "Standard_GS4-8"
	// VirtualMachineSizeTypesStandardGS5 ...
	VirtualMachineSizeTypesStandardGS5 VirtualMachineSizeTypes = "Standard_GS5"
	// VirtualMachineSizeTypesStandardGS516 ...
	VirtualMachineSizeTypesStandardGS516 VirtualMachineSizeTypes = "Standard_GS5-16"
	// VirtualMachineSizeTypesStandardGS58 ...
	VirtualMachineSizeTypesStandardGS58 VirtualMachineSizeTypes = "Standard_GS5-8"
	// VirtualMachineSizeTypesStandardH16 ...
	VirtualMachineSizeTypesStandardH16 VirtualMachineSizeTypes = "Standard_H16"
	// VirtualMachineSizeTypesStandardH16m ...
	VirtualMachineSizeTypesStandardH16m VirtualMachineSizeTypes = "Standard_H16m"
	// VirtualMachineSizeTypesStandardH16mr ...
	VirtualMachineSizeTypesStandardH16mr VirtualMachineSizeTypes = "Standard_H16mr"
	// VirtualMachineSizeTypesStandardH16r ...
	VirtualMachineSizeTypesStandardH16r VirtualMachineSizeTypes = "Standard_H16r"
	// VirtualMachineSizeTypesStandardH8 ...
	VirtualMachineSizeTypesStandardH8 VirtualMachineSizeTypes = "Standard_H8"
	// VirtualMachineSizeTypesStandardH8m ...
	VirtualMachineSizeTypesStandardH8m VirtualMachineSizeTypes = "Standard_H8m"
	// VirtualMachineSizeTypesStandardL16s ...
	VirtualMachineSizeTypesStandardL16s VirtualMachineSizeTypes = "Standard_L16s"
	// VirtualMachineSizeTypesStandardL32s ...
	VirtualMachineSizeTypesStandardL32s VirtualMachineSizeTypes = "Standard_L32s"
	// VirtualMachineSizeTypesStandardL4s ...
	VirtualMachineSizeTypesStandardL4s VirtualMachineSizeTypes = "Standard_L4s"
	// VirtualMachineSizeTypesStandardL8s ...
	VirtualMachineSizeTypesStandardL8s VirtualMachineSizeTypes = "Standard_L8s"
	// VirtualMachineSizeTypesStandardM12832ms ...
	VirtualMachineSizeTypesStandardM12832ms VirtualMachineSizeTypes = "Standard_M128-32ms"
	// VirtualMachineSizeTypesStandardM12864ms ...
	VirtualMachineSizeTypesStandardM12864ms VirtualMachineSizeTypes = "Standard_M128-64ms"
	// VirtualMachineSizeTypesStandardM128ms ...
	VirtualMachineSizeTypesStandardM128ms VirtualMachineSizeTypes = "Standard_M128ms"
	// VirtualMachineSizeTypesStandardM128s ...
	VirtualMachineSizeTypesStandardM128s VirtualMachineSizeTypes = "Standard_M128s"
	// VirtualMachineSizeTypesStandardM6416ms ...
	VirtualMachineSizeTypesStandardM6416ms VirtualMachineSizeTypes = "Standard_M64-16ms"
	// VirtualMachineSizeTypesStandardM6432ms ...
	VirtualMachineSizeTypesStandardM6432ms VirtualMachineSizeTypes = "Standard_M64-32ms"
	// VirtualMachineSizeTypesStandardM64ms ...
	VirtualMachineSizeTypesStandardM64ms VirtualMachineSizeTypes = "Standard_M64ms"
	// VirtualMachineSizeTypesStandardM64s ...
	VirtualMachineSizeTypesStandardM64s VirtualMachineSizeTypes = "Standard_M64s"
	// VirtualMachineSizeTypesStandardNC12 ...
	VirtualMachineSizeTypesStandardNC12 VirtualMachineSizeTypes = "Standard_NC12"
	// VirtualMachineSizeTypesStandardNC12sV2 ...
	VirtualMachineSizeTypesStandardNC12sV2 VirtualMachineSizeTypes = "Standard_NC12s_v2"
	// VirtualMachineSizeTypesStandardNC12sV3 ...
	VirtualMachineSizeTypesStandardNC12sV3 VirtualMachineSizeTypes = "Standard_NC12s_v3"
	// VirtualMachineSizeTypesStandardNC24 ...
	VirtualMachineSizeTypesStandardNC24 VirtualMachineSizeTypes = "Standard_NC24"
	// VirtualMachineSizeTypesStandardNC24r ...
	VirtualMachineSizeTypesStandardNC24r VirtualMachineSizeTypes = "Standard_NC24r"
	// VirtualMachineSizeTypesStandardNC24rsV2 ...
	VirtualMachineSizeTypesStandardNC24rsV2 VirtualMachineSizeTypes = "Standard_NC24rs_v2"
	// VirtualMachineSizeTypesStandardNC24rsV3 ...
	VirtualMachineSizeTypesStandardNC24rsV3 VirtualMachineSizeTypes = "Standard_NC24rs_v3"
	// VirtualMachineSizeTypesStandardNC24sV2 ...
	VirtualMachineSizeTypesStandardNC24sV2 VirtualMachineSizeTypes = "Standard_NC24s_v2"
	// VirtualMachineSizeTypesStandardNC24sV3 ...
	VirtualMachineSizeTypesStandardNC24sV3 VirtualMachineSizeTypes = "Standard_NC24s_v3"
	// VirtualMachineSizeTypesStandardNC6 ...
	VirtualMachineSizeTypesStandardNC6 VirtualMachineSizeTypes = "Standard_NC6"
	// VirtualMachineSizeTypesStandardNC6sV2 ...
	VirtualMachineSizeTypesStandardNC6sV2 VirtualMachineSizeTypes = "Standard_NC6s_v2"
	// VirtualMachineSizeTypesStandardNC6sV3 ...
	VirtualMachineSizeTypesStandardNC6sV3 VirtualMachineSizeTypes = "Standard_NC6s_v3"
	// VirtualMachineSizeTypesStandardND12s ...
	VirtualMachineSizeTypesStandardND12s VirtualMachineSizeTypes = "Standard_ND12s"
	// VirtualMachineSizeTypesStandardND24rs ...
	VirtualMachineSizeTypesStandardND24rs VirtualMachineSizeTypes = "Standard_ND24rs"
	// VirtualMachineSizeTypesStandardND24s ...
	VirtualMachineSizeTypesStandardND24s VirtualMachineSizeTypes = "Standard_ND24s"
	// VirtualMachineSizeTypesStandardND6s ...
	VirtualMachineSizeTypesStandardND6s VirtualMachineSizeTypes = "Standard_ND6s"
	// VirtualMachineSizeTypesStandardNV12 ...
	VirtualMachineSizeTypesStandardNV12 VirtualMachineSizeTypes = "Standard_NV12"
	// VirtualMachineSizeTypesStandardNV24 ...
	VirtualMachineSizeTypesStandardNV24 VirtualMachineSizeTypes = "Standard_NV24"
	// VirtualMachineSizeTypesStandardNV6 ...
	VirtualMachineSizeTypesStandardNV6 VirtualMachineSizeTypes = "Standard_NV6"
	// VirtualMachineSizeTypesStandardNK6 ...
	VirtualMachineSizeTypesStandardNK6 VirtualMachineSizeTypes = "Standard_NK6"
	// VirtualMachineSizeTypesStandardNK12 ...
	VirtualMachineSizeTypesStandardNK12 VirtualMachineSizeTypes = "Standard_NK12"
)

func GetVirtualMachineSizes() (vmsizes *[]VirtualMachineSizeTypes) {
	tmp := []VirtualMachineSizeTypes{}
	for key := range cloudcompute.VirtualMachineSizeType_name {
		tmp = append(tmp, GetCloudSdkVirtualMachineSizeFromCloudVirtualMachineSize(cloudcompute.VirtualMachineSizeType(key)))
	}

	vmsizes = &tmp
	return
}
