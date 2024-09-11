// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 License.

//go:build windows
// +build windows

//
// This file contains wrapper function calls that c++ component
// can leverage to call into MocStack
//

package main

// the below blob is NOT a comment, it is required to compile the c wrapper

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"context"
	"time"
	"unsafe"

	"github.com/microsoft/moc-sdk-for-go/services/security/authentication"
	"github.com/microsoft/moc-sdk-for-go/services/security/keyvault"
	"github.com/microsoft/moc-sdk-for-go/services/security/keyvault/key"
	"github.com/microsoft/moc-sdk-for-go/wrapper"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/config"
)

// possible Win32 return values
const (
	Win32Succeed int = 0               // ERROR_SUCCESS
	Win32ErrorInsufficientBuffer int = 122    // ERROR_INSUFFICIENT_BUFFER
	Win32ErrorBadArg int = 160         // ERROR_BAD_ARGUMENTS
	Win32ErrorFunctionFail int = 1627  // ERROR_FUNCTION_FAILED
)

// helper function to best effort copy over the error message
func copyErrorMessage(errMessage *C.char, errMessageBuffer *C.char, errMessageSize C.ulonglong) {
	if (errMessage != nil) {
		if (errMessageBuffer != nil && errMessageSize > 0) {
			msglength := C.strlen(errMessage)
			if (errMessageSize < msglength) {
				msglength = errMessageSize
			}
			C.strncpy(errMessageBuffer, errMessage, msglength)
		}
		C.free(unsafe.Pointer(errMessage))
	}
}

// This function exists to maintain backwards compatability. Please use SecurityLoginV2.
//
//export SecurityLogin
func SecurityLogin(serverName *C.char, groupName *C.char, loginFilePath *C.char, timeoutInSeconds C.int) *C.char {
	return SecurityLoginCV(serverName, groupName, loginFilePath, C.CString(""), timeoutInSeconds)
}

// This function exists to maintain backwards compatability. Please use SecurityLoginV2.
//
//export SecurityLoginCV
func SecurityLoginCV(serverName *C.char, groupName *C.char, loginFilePath *C.char, cv *C.char, timeoutInSeconds C.int) *C.char {
	loginconfig := auth.LoginConfig{}
	err := config.LoadYAMLFile(C.GoString(loginFilePath), &loginconfig)
	if err != nil {
		telemetry.EmitWrapperTelemetry("SecurityLoginCV", C.GoString(cv), err.Error(), "config.LoadYAMLFile", C.GoString(serverName))
		return C.CString(telemetry.FilterSensitiveData(err.Error()))
	}

	authenticationClient, err := authentication.NewAuthenticationClientAuthMode(C.GoString(serverName), loginconfig)
	if err != nil {
		telemetry.EmitWrapperTelemetry("SecurityLoginCV", C.GoString(cv), err.Error(), "authentication.NewAuthenticationClientAuthMode", C.GoString(serverName))
		return C.CString(telemetry.FilterSensitiveData(err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	// Login with config stores the access file in the WSSD_CONFIG environment variable
	// set true to auto renew
	_, err = authenticationClient.LoginWithConfig(ctx, C.GoString(groupName), loginconfig, true)
	if err != nil {
		telemetry.EmitWrapperTelemetry("SecurityLoginCV", C.GoString(cv), err.Error(), "authenticationClient.LoginWithConfig", C.GoString(serverName))
		return C.CString(telemetry.FilterSensitiveData(err.Error()))
	}

	//Provide moc version information after login
	telemetry.EmitWrapperTelemetry("SecurityLoginCV", C.GoString(cv), "", "", C.GoString(serverName))
	return nil
}

// This function performs security login
// parameters:
//   - serverName
//   - groupName
//   - loginFilePath
//   - cv: for logging purpose
//   - timeoutInSeconds
//   - errMessageBuffer: the pre allocated buffer to accept error message if there is any. If the buffer is not pre allocated, the error messaged won't be copied
//   - errMessageSize: the size of pre allocated errMessageBuffer. If the actual error message size is bigger, the error message will be truncated to errMessageSize
// return:
//   - int: the win32 code
//
//export SecurityLoginV2
func SecurityLoginV2(serverName *C.char, groupName *C.char, loginFilePath *C.char, cv *C.char, timeoutInSeconds C.int, errMessageBuffer *C.char, errMessageSize C.ulonglong) C.int {
	if (serverName == nil || groupName == nil || loginFilePath == nil || cv == nil) {
		telemetry.EmitWrapperTelemetry("SecurityLoginV2", C.GoString(cv), "", "InvalidArgument", C.GoString(serverName))
		copyErrorMessage(C.CString("invalid argument"), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorBadArg)
	}

	loginconfig := auth.LoginConfig{}
	err := config.LoadYAMLFile(C.GoString(loginFilePath), &loginconfig)
	if err != nil {
		telemetry.EmitWrapperTelemetry("SecurityLoginV2", C.GoString(cv), err.Error(), "config.LoadYAMLFile", C.GoString(serverName))
		copyErrorMessage(C.CString(telemetry.FilterSensitiveData(err.Error())), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorFunctionFail)
	}

	authenticationClient, err := authentication.NewAuthenticationClientAuthMode(C.GoString(serverName), loginconfig)
	if err != nil {
		telemetry.EmitWrapperTelemetry("SecurityLoginV2", C.GoString(cv), err.Error(), "authentication.NewAuthenticationClientAuthMode", C.GoString(serverName))
		copyErrorMessage(C.CString(telemetry.FilterSensitiveData(err.Error())), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorFunctionFail)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	// Login with config stores the access file in the WSSD_CONFIG environment variable
	// set true to auto renew
	_, err = authenticationClient.LoginWithConfig(ctx, C.GoString(groupName), loginconfig, true)
	if err != nil {
		telemetry.EmitWrapperTelemetry("SecurityLoginV2", C.GoString(cv), err.Error(), "authenticationClient.LoginWithConfig", C.GoString(serverName))
		copyErrorMessage(C.CString(telemetry.FilterSensitiveData(err.Error())), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorFunctionFail)
	}

	//Provide moc version information after login
	telemetry.EmitWrapperTelemetry("SecurityLoginV2", C.GoString(cv), "", "", C.GoString(serverName))
	return C.int(Win32Succeed)
}

// This function exists to maintain backwards compatability. Please use KeyvaultKeyEncryptDataV2.
//
//export KeyvaultKeyEncryptData
func KeyvaultKeyEncryptData(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, input *C.char, timeoutInSeconds C.int) *C.char {
	return KeyvaultKeyEncryptDataCV(serverName, groupName, keyvaultName, keyName, input, C.CString(""), timeoutInSeconds)
}

// This function exists to maintain backwards compatability. Please use KeyvaultKeyEncryptDataV2.
//
//export KeyvaultKeyEncryptDataCV
func KeyvaultKeyEncryptDataCV(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, input *C.char, cv *C.char, timeoutInSeconds C.int) *C.char {
	keyClient, err := getKeyvaultKeyClient(C.GoString(serverName), C.GoString(cv))
	// if errror occurs, return an empty string so that caller can tell between error and encrypted blob
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyEncryptDataCV", C.GoString(cv), err.Error(), "getKeyvaultKeyClient", C.GoString(serverName))
		return C.CString("")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	// input is base64 encoded
	var value string
	value = C.GoString(input)

	parameters := &keyvault.KeyOperationsParameters{
		Value:     &value,
		Algorithm: keyvault.A256CBC,
	}

	response, err := keyClient.Encrypt(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName), parameters)
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyEncryptDataCV", C.GoString(cv), err.Error(), "keyClient.Encrypt", C.GoString(serverName))
		return C.CString("")
	}

	// retrun base64 encoded string
	return C.CString(*response.Result)
}

// This function performs Encrypt
// parameters:
//   - serverName
//   - groupName
//   - keyvaultName
//   - keyName
//   - input: the data to be encrypted
//   - algorithm: the encryption algorithm
//   - cv: for logging purpose
//   - timeoutInSeconds
//   - outputBuffer: the buffer to accept encrypted data
//   - outputBufferSize: the size of the outputBuffer. If this is not allocated or size is not big enough the expected size will be set to outputBufferSize.
//   - errMessageBuffer: the pre allocated buffer to accept error message if there is any. If the buffer is not pre allocated, the error messaged won't be copied
//   - errMessageSize: the size of pre allocated errMessageBuffer. If the actual error message size is bigger, the error message will be truncated to errMessageSize
// return:
//   - int: the win32 code. If output buffer is not big enough this function will return Win32ErrorInsufficientBuffer which is win32 error ERROR_INSUFFICIENT_BUFFER.
//     Caller when seeing this error should allocate outputBuffer to size outputBufferSize
// Calling pattern:
//   call KeyvaultKeyEncryptDataV2 for the first time. If returned error Win32ErrorInsufficientBuffer, caller allocate outputBuffer to size outputBufferSize and call KeyvaultKeyEncryptDataV2 again.
//
//export KeyvaultKeyEncryptDataV2
func KeyvaultKeyEncryptDataV2(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, input *C.char, algorithm *C.char, cv *C.char, timeoutInSeconds C.int, outputBuffer *C.char, outputBufferSize *C.ulonglong, errMessageBuffer *C.char, errMessageSize C.ulonglong) C.int {
	if (serverName == nil || groupName == nil || keyvaultName == nil || keyName == nil || input == nil || algorithm == nil || cv == nil || outputBufferSize == nil) {
		copyErrorMessage(C.CString("Invalid Argument"), errMessageBuffer, errMessageSize)
		telemetry.EmitWrapperTelemetry("KeyvaultKeyEncryptDataV2", C.GoString(cv), "", "InvalidArgument", C.GoString(serverName))
		return C.int(Win32ErrorBadArg)
	}

	keyClient, err := getKeyvaultKeyClient(C.GoString(serverName), C.GoString(cv))
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyEncryptDataV2", C.GoString(cv), err.Error(), "getKeyvaultKeyClient", C.GoString(serverName))
		copyErrorMessage(C.CString(telemetry.FilterSensitiveData(err.Error())), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorFunctionFail)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	// input is base64 encoded
	var value string
	value = C.GoString(input)

	alg :=keyvault.JSONWebKeyEncryptionAlgorithm(C.GoString(algorithm))
	parameters := &keyvault.KeyOperationsParameters{
		Value:     &value,
		Algorithm: alg,
	}

	response, err := keyClient.Encrypt(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName), parameters)
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyEncryptDataV2", C.GoString(cv), err.Error(), "keyClient.Encrypt", C.GoString(serverName))
		copyErrorMessage(C.CString(telemetry.FilterSensitiveData(err.Error())), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorFunctionFail)
	}

	return handleOutputBuffer(serverName, *response.Result, cv, outputBuffer, outputBufferSize, "KeyvaultKeyDecryptDataV2")
}

// This function exists to maintain backwards compatability. Please use KeyvaultKeyDecryptDataV2.
//
//export KeyvaultKeyDecryptData
func KeyvaultKeyDecryptData(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, input *C.char, timeoutInSeconds C.int) *C.char {
	return KeyvaultKeyDecryptDataCV(serverName, groupName, keyvaultName, keyName, input, C.CString(""), timeoutInSeconds)
}

// This function exists to maintain backwards compatability. Please use KeyvaultKeyDecryptDataV2.
//
//export KeyvaultKeyDecryptDataCV
func KeyvaultKeyDecryptDataCV(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, input *C.char, cv *C.char, timeoutInSeconds C.int) *C.char {
	keyClient, err := getKeyvaultKeyClient(C.GoString(serverName), C.GoString(cv))
	// if errror occurs, return an empty string so that caller can tell between error and decrypted blob
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyDecryptDataCV", C.GoString(cv), err.Error(), "getKeyvaultKeyClient", C.GoString(serverName))
		return C.CString("")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	var value string
	value = C.GoString(input)

	parameters := &keyvault.KeyOperationsParameters{
		Value:     &value,
		Algorithm: keyvault.A256CBC,
	}

	response, err := keyClient.Decrypt(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName), parameters)
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyDecryptDataCV", C.GoString(cv), err.Error(), "keyClient.Decrypt", C.GoString(serverName))
		return C.CString("")
	}

	return C.CString(*response.Result)
}

// This function performs Decrypt
// parameters:
//   - serverName
//   - groupName
//   - keyvaultName
//   - keyName
//   - input: the data to be decrypted
//   - algorithm: the decryption algorithm
//   - cv: for logging purpose
//   - timeoutInSeconds
//   - outputBuffer: the buffer to accept decrypted data
//   - outputBufferSize: the size of the outputBuffer. If this is not allocated or size is not big enough the expected size will be set to outputBufferSize.
//   - errMessageBuffer: the pre allocated buffer to accept error message if there is any. If the buffer is not pre allocated, the error messaged won't be copied
//   - errMessageSize: the size of pre allocated errMessageBuffer. If the actual error message size is bigger, the error message will be truncated to errMessageSize
// return:
//   - int: the win32 code. If output buffer is not big enough this function will return Win32ErrorInsufficientBuffer which is win32 error ERROR_INSUFFICIENT_BUFFER.
//     Caller when seeing this error should allocate outputBuffer to size outputBufferSize
// Calling pattern:
//   call KeyvaultKeyDecryptDataV2 for the first time. If returned error Win32ErrorInsufficientBuffer, caller allocate outputBuffer to size outputBufferSize and call KeyvaultKeyDecryptDataV2 again.
//
//export KeyvaultKeyDecryptDataV2
func KeyvaultKeyDecryptDataV2(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, input *C.char, algorithm *C.char, cv *C.char, timeoutInSeconds C.int, outputBuffer *C.char, outputBufferSize *C.ulonglong, errMessageBuffer *C.char, errMessageSize C.ulonglong) C.int {
	if (serverName == nil || groupName == nil || keyvaultName == nil || keyName == nil || input == nil || algorithm == nil || cv == nil || outputBufferSize == nil) {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyDecryptDataV2", C.GoString(cv), "", "InvalidArgument", C.GoString(serverName))
		copyErrorMessage(C.CString("Invalid Argument"), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorBadArg)
	}

	keyClient, err := getKeyvaultKeyClient(C.GoString(serverName), C.GoString(cv))
	if err != nil {
		copyErrorMessage(C.CString(telemetry.FilterSensitiveData(err.Error())), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorFunctionFail)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	var value string
	value = C.GoString(input)

	alg :=keyvault.JSONWebKeyEncryptionAlgorithm(C.GoString(algorithm))
	parameters := &keyvault.KeyOperationsParameters{
		Value:     &value,
		Algorithm: alg,
	}

	response, err := keyClient.Decrypt(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName), parameters)
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyDecryptDataV2", C.GoString(cv), err.Error(), "keyClient.Decrypt", C.GoString(serverName))
		copyErrorMessage(C.CString(telemetry.FilterSensitiveData(err.Error())), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorFunctionFail)
	}

	return handleOutputBuffer(serverName, *response.Result, cv, outputBuffer, outputBufferSize, "KeyvaultKeyDecryptDataV2")
}

// This function exists to maintain backwards compatability. Please use KeyvaultKeyExistCV.
//
//export KeyvaultKeyExist
func KeyvaultKeyExist(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, timeoutInSeconds C.int) C.int {
	return KeyvaultKeyExistCV(serverName, groupName, keyvaultName, keyName, C.CString(""), timeoutInSeconds)
}

//export KeyvaultKeyExistCV
func KeyvaultKeyExistCV(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, cv *C.char, timeoutInSeconds C.int) C.int {
	keyClient, err := getKeyvaultKeyClient(C.GoString(serverName), C.GoString(cv))
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyExistCV", C.GoString(cv), err.Error(), "getKeyvaultKeyClient", C.GoString(serverName))
		return -1
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	keys, err := keyClient.Get(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName))
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyExistCV", C.GoString(cv), err.Error(), "keyClient.Get", C.GoString(serverName))
		return -1
	}

	// check the length and return 1 (means key exists) if there is at least one key
	if keys != nil && len(*keys) > 0 {
		return 1
	}

	return 0
}

// This function exists to maintain backwards compatability. Please use KeyvaultKeyCreateOrUpdateV2.
//
//export KeyvaultKeyCreateOrUpdate
func KeyvaultKeyCreateOrUpdate(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, keyTypeName *C.char, keySize C.int, timeoutInSeconds C.int) *C.char {
	return KeyvaultKeyCreateOrUpdateCV(serverName, groupName, keyvaultName, keyName, keyTypeName, keySize, C.CString(""), timeoutInSeconds)
}

// This function exists to maintain backwards compatability. Please use KeyvaultKeyCreateOrUpdateV2.
//
//export KeyvaultKeyCreateOrUpdateCV
func KeyvaultKeyCreateOrUpdateCV(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, keyTypeName *C.char, keySize C.int, cv *C.char, timeoutInSeconds C.int) *C.char {
	keyClient, err := getKeyvaultKeyClient(C.GoString(serverName), C.GoString(cv))
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyCreateOrUpdateCV", C.GoString(cv), err.Error(), "getKeyvaultKeyClient", C.GoString(serverName))
		return C.CString(telemetry.FilterSensitiveData(err.Error()))
	}

	var kvConfig *keyvault.Key
	kvConfig = &keyvault.Key{}

	var keyNameString string
	keyNameString = C.GoString(keyName)
	kvConfig.Name = &keyNameString
	kvConfig.KeyProperties = &keyvault.KeyProperties{}

	kvConfig.KeyType = keyvault.JSONWebKeyType(C.GoString(keyTypeName))
	var tKeySize int32
	tKeySize = int32(keySize)
	kvConfig.KeySize = &tKeySize

	var keyRotation int64
	keyRotation = -1
	kvConfig.KeyRotationFrequencyInSeconds = &keyRotation // -1 means disable key rotation

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	_, err = keyClient.CreateOrUpdate(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName), kvConfig)
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyCreateOrUpdateCV", C.GoString(cv), err.Error(), "keyClient.CreateOrUpdate", C.GoString(serverName))
		//This return cannot be empty!
		return C.CString(telemetry.FilterSensitiveData(err.Error()))
	}

	return nil
}

// This function performs key creation or update
// parameters:
//   - serverName
//   - groupName
//   - keyvaultName
//   - keyName
//   - keyTypeName
//   - keySize
//   - cv: for logging purpose
//   - timeoutInSeconds
//   - errMessageBuffer: the pre allocated buffer to accept error message if there is any. If the buffer is not pre allocated, the error messaged won't be copied
//   - errMessageSize: the size of pre allocated errMessageBuffer. If the actual error message size is bigger, the error message will be truncated to errMessageSize
// return:
//   - int: the win32 code.
//
//export KeyvaultKeyCreateOrUpdateV2
func KeyvaultKeyCreateOrUpdateV2(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, keyTypeName *C.char, keySize C.int, cv *C.char, timeoutInSeconds C.int, errMessageBuffer *C.char, errMessageSize C.ulonglong) C.int {
	if (serverName == nil || groupName == nil || keyvaultName == nil || keyName == nil || keyTypeName == nil || cv == nil) {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyCreateOrUpdateV2", C.GoString(cv), "", "InvalidArgument", C.GoString(serverName))
		copyErrorMessage(C.CString("Invalid Argument"), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorBadArg)
	}

	keyClient, err := getKeyvaultKeyClient(C.GoString(serverName), C.GoString(cv))
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyCreateOrUpdateV2", C.GoString(cv), err.Error(), "getKeyvaultKeyClient", C.GoString(serverName))
		copyErrorMessage(C.CString(telemetry.FilterSensitiveData(err.Error())), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorFunctionFail)
	}

	var kvConfig *keyvault.Key
	kvConfig = &keyvault.Key{}

	var keyNameString string
	keyNameString = C.GoString(keyName)
	kvConfig.Name = &keyNameString
	kvConfig.KeyProperties = &keyvault.KeyProperties{}

	kvConfig.KeyType = keyvault.JSONWebKeyType(C.GoString(keyTypeName))
	var tKeySize int32
	tKeySize = int32(keySize)
	kvConfig.KeySize = &tKeySize

	var keyRotation int64
	keyRotation = -1
	kvConfig.KeyRotationFrequencyInSeconds = &keyRotation // -1 means disable key rotation

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	_, err = keyClient.CreateOrUpdate(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName), kvConfig)
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyCreateOrUpdateV2", C.GoString(cv), err.Error(), "keyClient.CreateOrUpdate", C.GoString(serverName))
		//This return cannot be empty!
		copyErrorMessage(C.CString(telemetry.FilterSensitiveData(err.Error())), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorFunctionFail)
	}

	return C.int(Win32Succeed)
}

// This function exists to maintain backwards compatability. Please use KeyvaultKeySignDataCV.
//
//export KeyvaultKeySignData
func KeyvaultKeySignData(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, input *C.char, algorithm *C.char, timeoutInSeconds C.int) *C.char {
	return KeyvaultKeySignDataCV(serverName, groupName, keyvaultName, keyName, input, algorithm, C.CString(""), timeoutInSeconds)
}

//export KeyvaultKeySignDataCV
func KeyvaultKeySignDataCV(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, input *C.char, algorithm *C.char, cv *C.char, timeoutInSeconds C.int) *C.char {
	keyClient, err := getKeyvaultKeyClient(C.GoString(serverName), C.GoString(cv))
	// if errror occurs, return an empty string so that caller can tell between error and decrypted blob
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeySignDataCV", C.GoString(cv), err.Error(), "getKeyvaultKeyClient", C.GoString(serverName))
		return C.CString("")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	// input is base64 encoded
	var value string
	value = C.GoString(input)
	algo := keyvault.JSONWebKeySignatureAlgorithm(C.GoString(algorithm))
	parameters := &keyvault.KeySignParameters{
		Value:     &value,
		Algorithm: algo,
	}

	response, err := keyClient.Sign(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName), parameters)
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeySignDataCV", C.GoString(cv), err.Error(), "keyClient.Sign", C.GoString(serverName))
		return C.CString("")
	}

	// retrun base64 encoded string
	return C.CString(*response.Result)
}

// This function exists to maintain backwards compatability. Please use KeyvaultKeyVerifyDataCV.
//
//export KeyvaultKeyVerifyData
func KeyvaultKeyVerifyData(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, digest *C.char, signature *C.char, algorithm *C.char, timeoutInSeconds C.int) (ret C.int) {
	return KeyvaultKeyVerifyDataCV(serverName, groupName, keyvaultName, keyName, digest, signature, algorithm, C.CString(""), timeoutInSeconds)
}

//export KeyvaultKeyVerifyDataCV
func KeyvaultKeyVerifyDataCV(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, digest *C.char, signature *C.char, algorithm *C.char, cv *C.char, timeoutInSeconds C.int) (ret C.int) {
	keyClient, err := getKeyvaultKeyClient(C.GoString(serverName), C.GoString(cv))
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyVerifyDataCV", C.GoString(cv), err.Error(), "getKeyvaultKeyClient", C.GoString(serverName))
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	// input is base64 encoded
	var value string
	value = C.GoString(digest)

	var sig string
	sig = C.GoString(signature)
	algo := keyvault.JSONWebKeySignatureAlgorithm(C.GoString(algorithm))
	parameters := &keyvault.KeyVerifyParameters{
		Signature: &sig,
		Digest:    &value,
		Algorithm: algo,
	}

	response, err := keyClient.Verify(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName), parameters)
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultKeyVerifyDataCV", C.GoString(cv), err.Error(), "keyClient.Verify", C.GoString(serverName))
		return C.int(0)
	}

	//Convert to Int, (is there away to use c.bool?)
	if *response.Value {
		return C.int(1)
	} else {
		return C.int(0)
	}
}

// This function exists to maintain backwards compatability. Please use KeyvaultGetPublicKeyV2.
//
//export KeyvaultGetPublicKey
func KeyvaultGetPublicKey(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, timeoutInSeconds C.int) *C.char {
	return KeyvaultGetPublicKeyCV(serverName, groupName, keyvaultName, keyName, C.CString(""), timeoutInSeconds)
}

// This function exists to maintain backwards compatability. Please use KeyvaultGetPublicKeyV2.
//
//export KeyvaultGetPublicKeyCV
func KeyvaultGetPublicKeyCV(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, cv *C.char, timeoutInSeconds C.int) *C.char {
	keyClient, err := getKeyvaultKeyClient(C.GoString(serverName), C.GoString(cv))
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultGetPublicKeyCV", C.GoString(cv), err.Error(), "getKeyvaultKeyClient", C.GoString(serverName))
		return C.CString("")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	keys, err := keyClient.Get(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName))
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultGetPublicKeyCV", C.GoString(cv), err.Error(), "keyClient.Get", C.GoString(serverName))
		return C.CString("")
	}

	// check the length of the key, if it is zero or we don't have any keys, there is no public key to return
	if keys == nil || len(*keys) <= 0 {
		return C.CString("")
	}

	pemPkcs1KeyPub := (*keys)[0].Value

	return C.CString(*pemPkcs1KeyPub)
}

// This function performs public key retrieval
// parameters:
//   - serverName
//   - groupName
//   - keyvaultName
//   - keyName
//   - cv: for logging purpose
//   - timeoutInSeconds
//   - outputBuffer: the buffer to accept public key data
//   - outputBufferSize: the size of the outputBuffer. If this is not allocated or size is not big enough the expected size will be set to outputBufferSize.
//   - errMessageBuffer: the pre allocated buffer to accept error message if there is any. If the buffer is not pre allocated, the error messaged won't be copied
//   - errMessageSize: the size of pre allocated errMessageBuffer. If the actual error message size is bigger, the error message will be truncated to errMessageSize
// return:
//   - int: the win32 code. If output buffer is not big enough this function will return Win32ErrorInsufficientBuffer which is win32 error ERROR_INSUFFICIENT_BUFFER.
//     Caller when seeing this error should allocate outputBuffer to size outputBufferSize
// Calling pattern:
//   call KeyvaultGetPublicKeyV2 for the first time. If returned error Win32ErrorInsufficientBuffer, caller allocate outputBuffer to size outputBufferSize and call KeyvaultGetPublicKeyV2 again.
//
//export KeyvaultGetPublicKeyV2
func KeyvaultGetPublicKeyV2(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, cv *C.char, timeoutInSeconds C.int, outputBuffer *C.char, outputBufferSize *C.ulonglong, errMessageBuffer *C.char, errMessageSize C.ulonglong) C.int {
	if (serverName == nil || groupName == nil || keyvaultName == nil || keyName == nil || outputBufferSize == nil || cv == nil) {
		telemetry.EmitWrapperTelemetry("KeyvaultGetPublicKeyV2", C.GoString(cv), "", "InvalidArgument", C.GoString(serverName))
		copyErrorMessage(C.CString("Invalid Argument"), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorBadArg)
	}

	keyClient, err := getKeyvaultKeyClient(C.GoString(serverName), C.GoString(cv))
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultGetPublicKeyV2", C.GoString(cv), err.Error(), "getKeyvaultKeyClient", C.GoString(serverName))
		copyErrorMessage(C.CString(telemetry.FilterSensitiveData(err.Error())), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorFunctionFail)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	keys, err := keyClient.Get(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName))
	if err != nil {
		telemetry.EmitWrapperTelemetry("KeyvaultGetPublicKeyV2", C.GoString(cv), err.Error(), "keyClient.Get", C.GoString(serverName))
		copyErrorMessage(C.CString(telemetry.FilterSensitiveData(err.Error())), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorFunctionFail)
	}

	if keys == nil || len(*keys) <= 0 {
		telemetry.EmitWrapperTelemetry("KeyvaultGetPublicKeyV2", C.GoString(cv), "", "EmptyKey", C.GoString(serverName))
		copyErrorMessage(C.CString("Returned key is empty"), errMessageBuffer, errMessageSize)
		return C.int(Win32ErrorFunctionFail)
	}

	pemPkcs1KeyPub := (*keys)[0].Value
	return handleOutputBuffer(serverName, *pemPkcs1KeyPub, cv, outputBuffer, outputBufferSize, "KeyvaultGetPublicKeyV2")
}

func handleOutputBuffer(serverName *C.char, resultGoString string, cv *C.char, outputBuffer *C.char, outputBufferSize *C.ulonglong, funcName string) C.int {
	if (outputBufferSize == nil) {
		return C.int(Win32ErrorBadArg)
	}

	resultCString := C.CString(resultGoString)
	var resultCStringLength C.ulonglong = C.strlen(resultCString)

	if (outputBuffer == nil || *outputBufferSize < resultCStringLength) {
		telemetry.EmitWrapperTelemetry(funcName, C.GoString(cv), "", "InsufficientBuffer", C.GoString(serverName))
		*outputBufferSize = resultCStringLength;
		C.free(unsafe.Pointer(resultCString))
		return C.int(Win32ErrorInsufficientBuffer)
	}

	C.strncpy(outputBuffer, resultCString, resultCStringLength)
	*outputBufferSize = resultCStringLength
	C.free(unsafe.Pointer(resultCString))
	return C.int(Win32Succeed)
}

func getKeyvaultKeyClient(serverName string, cv string) (*key.KeyClient, error) {
	authorizer, err := auth.NewAuthorizerFromEnvironment(serverName)
	if err != nil {
		telemetry.EmitWrapperTelemetry("getKeyvaultKeyClient", cv, err.Error(), "auth.NewAuthorizerFromEnvironment", serverName)
		return nil, err
	}

	return key.NewKeyClient(serverName, authorizer)
}

func main() {}
