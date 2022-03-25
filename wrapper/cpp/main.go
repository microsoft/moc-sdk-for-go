// Copyright (c) Microsoft Corporation.
// Licensed under the Apache v2.0 License.

//
// This file contains wrapper function calls that c++ component
// can leverage to call into MocStack
//

package main

import "C"

import (
    "context"
    "time"
    "github.com/microsoft/moc/pkg/auth"
    "github.com/microsoft/moc-sdk-for-go/services/security/keyvault"
    "github.com/microsoft/moc-sdk-for-go/services/security/keyvault/key"
)

//export KeyvaultKeyEncryptData
func KeyvaultKeyEncryptData(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, input *C.char, timeoutInSeconds C.int) *C.char {
    keyClient, err := getKeyvaultKeyClient(C.GoString(serverName))
    // if errror occurs, return an empty string so that caller can tell between error and encrypted blob
    if err != nil {
        return C.CString("")
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
    defer cancel()

    // input is base64 encoded
    var value string
    value = C.GoString(input)

    parameters := &keyvault.KeyOperationsParameters{
        Value:     &value,
        Algorithm: keyvault.A256KW,
    }

    response, err := keyClient.Encrypt(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName), parameters)
    if err != nil {
        return C.CString("")
    }

    // retrun base64 encoded string
    return  C.CString(*response.Result)
}

//export KeyvaultKeyDecryptData
func KeyvaultKeyDecryptData(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, input *C.char, timeoutInSeconds C.int) *C.char {
    keyClient, err := getKeyvaultKeyClient(C.GoString(serverName))
    // if errror occurs, return an empty string so that caller can tell between error and decrypted blob
    if err != nil {
        return C.CString("")
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
    defer cancel()

    var value string
    value = C.GoString(input)

    parameters := &keyvault.KeyOperationsParameters{
        Value:     &value,
        Algorithm: keyvault.A256KW,
    }

    response, err := keyClient.Decrypt(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName), parameters)
    if err != nil {
        return C.CString("")
    }

    return  C.CString(*response.Result)
}

//export KeyvaultKeyExist
func KeyvaultKeyExist(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, timeoutInSeconds C.int) C.int {
    keyClient, err := getKeyvaultKeyClient(C.GoString(serverName))
    if err != nil {
        return 0
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
    defer cancel()

    keys, err := keyClient.Get(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName))
    if err != nil {
        return 0
    }

    // check the length and return 1 (means key exists) if there is more than one key
    if keys != nil && len(*keys) > 0 {
        return  1
    }

    return 0
}

//export KeyvaultKeyCreateOrUpdate
func KeyvaultKeyCreateOrUpdate(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, keyTypeName *C.char, timeoutInSeconds C.int) *C.char {
    keyClient, err := getKeyvaultKeyClient(C.GoString(serverName))
    if err != nil {
        return C.CString(err.Error())
    }

    var kvConfig *keyvault.Key
    kvConfig = &keyvault.Key{}

    var keyNameString string
    keyNameString = C.GoString(keyName)
    kvConfig.Name = &keyNameString
    kvConfig.KeyProperties = &keyvault.KeyProperties{}

    kvConfig.KeyType = keyvault.JSONWebKeyType(C.GoString(keyTypeName))
    var keySize int32
    keySize =256 // hardcode for AES key
    kvConfig.KeySize = &keySize

    var keyRotation int64
    keyRotation = -1
    kvConfig.KeyRotationFrequencyInSeconds = &keyRotation // -1 means disable key rotation

    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutInSeconds)*time.Second)
    defer cancel()

    _, err = keyClient.CreateOrUpdate(ctx, C.GoString(groupName), C.GoString(keyvaultName), C.GoString(keyName), kvConfig)
    if err != nil {
        return C.CString(err.Error())
    }

    return nil
}

func getKeyvaultKeyClient(serverName string) (*key.KeyClient, error) {
    authorizer, err := auth.NewAuthorizerFromEnvironment(serverName)
    if err != nil {
        return nil, err
    }

    return key.NewKeyClient(serverName, authorizer)
}

func main() {}
