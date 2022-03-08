// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package key

import (
	"encoding/pem"
	"github.com/microsoft/moc-sdk-for-go/services/security/keyvault"

	"github.com/microsoft/moc/pkg/convert"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	wssdcloudsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdcloudcommon "github.com/microsoft/moc/rpc/common"
)

func getKey(sec *wssdcloudsecurity.Key, vaultName string) (keyvault.Key, error) {
	keysize, err := getKeySize(sec.Size)
	if err != nil {
		return keyvault.Key{}, errors.Wrapf(err, "INVESTIAGE")
	}
	value := ""
	switch sec.Type {
	case wssdcloudcommon.JsonWebKeyType_RSA:
		fallthrough
	case wssdcloudcommon.JsonWebKeyType_RSA_HSM:
		pubBlk := pem.Block{
			Type:    "RSA PUBLIC KEY",
			Headers: nil,
			Bytes:   sec.PublicKey,
		}
		value = string(pem.EncodeToMemory(&pubBlk))
	}
	return keyvault.Key{
		ID:      &sec.Id,
		Name:    &sec.Name,
		Version: &sec.Status.Version.Number,
		Value:   &value,
		KeyProperties: &keyvault.KeyProperties{
			Statuses:                      status.GetStatuses(sec.GetStatus()),
			KeyType:                       getKeyType(sec.Type),
			KeySize:                       keysize,
			KeyRotationFrequencyInSeconds: &sec.KeyRotationFrequencyInSeconds,
		},
	}, nil
}

func getWssdKeyByVaultName(name string, groupName,
	vaultName string, opType wssdcloudcommon.Operation) (*wssdcloudsecurity.Key, error) {
	key := &wssdcloudsecurity.Key{
		Name:      name,
		VaultName: vaultName,
		GroupName: groupName,
		Type:      wssdcloudcommon.JsonWebKeyType_EC,
		Size:      wssdcloudcommon.KeySize_K_UNKNOWN,
		KeyOps:    []wssdcloudcommon.KeyOperation{},
		Status:    status.InitStatus(),
	}

	// No Update support
	return key, nil
}

func getWssdKey(name string, sec *keyvault.Key,
	groupName, vaultName string, opType wssdcloudcommon.Operation) (*wssdcloudsecurity.Key, error) {
	keysize, err := getMOCKeySize(*sec.KeySize)
	if err != nil {
		return nil, err
	}
	key := &wssdcloudsecurity.Key{
		Name:                          name,
		VaultName:                     vaultName,
		GroupName:                     groupName,
		Type:                          getMOCKeyType(sec.KeyType),
		Size:                          keysize,
		KeyOps:                        []wssdcloudcommon.KeyOperation{},
		Status:                        status.InitStatus(),
		KeyRotationFrequencyInSeconds: *sec.KeyRotationFrequencyInSeconds,
	}

	// No Update support
	return key, nil
}

func getMOCKeyType(ktype keyvault.JSONWebKeyType) wssdcloudcommon.JsonWebKeyType {
	switch ktype {
	case keyvault.EC:
		return wssdcloudcommon.JsonWebKeyType_EC
	case keyvault.ECHSM:
		return wssdcloudcommon.JsonWebKeyType_EC_HSM
	case keyvault.Oct:
		return wssdcloudcommon.JsonWebKeyType_OCT
	case keyvault.RSA:
		return wssdcloudcommon.JsonWebKeyType_RSA
	case keyvault.RSAHSM:
		return wssdcloudcommon.JsonWebKeyType_RSA_HSM
	case keyvault.AES:
		return wssdcloudcommon.JsonWebKeyType_AES
	default:
		return wssdcloudcommon.JsonWebKeyType_EC
	}
}

func getKeyType(ktype wssdcloudcommon.JsonWebKeyType) keyvault.JSONWebKeyType {
	switch ktype {
	case wssdcloudcommon.JsonWebKeyType_EC:
		return keyvault.EC
	case wssdcloudcommon.JsonWebKeyType_EC_HSM:
		return keyvault.ECHSM
	case wssdcloudcommon.JsonWebKeyType_OCT:
		return keyvault.Oct
	case wssdcloudcommon.JsonWebKeyType_RSA:
		return keyvault.RSA
	case wssdcloudcommon.JsonWebKeyType_RSA_HSM:
		return keyvault.RSAHSM
	case wssdcloudcommon.JsonWebKeyType_AES:
		return keyvault.AES
	default:
		return keyvault.EC
	}
}

func getMOCKeySize(size int32) (ksize wssdcloudcommon.KeySize, err error) {
	switch size {
	case 256:
		ksize = wssdcloudcommon.KeySize__256
	case 2048:
		ksize = wssdcloudcommon.KeySize__2048
	case 3072:
		ksize = wssdcloudcommon.KeySize__3072
	case 4096:
		ksize = wssdcloudcommon.KeySize__4096
	default:
		err = errors.Wrapf(errors.InvalidInput, "Invalid key size")
	}

	return
}

func getKeySize(ksize wssdcloudcommon.KeySize) (size *int32, err error) {
	switch ksize {
	case wssdcloudcommon.KeySize__256:
		size = convert.ToInt32Ptr(256)
	case wssdcloudcommon.KeySize__2048:
		size = convert.ToInt32Ptr(2048)
	case wssdcloudcommon.KeySize__3072:
		size = convert.ToInt32Ptr(3072)
	case wssdcloudcommon.KeySize__4096:
		size = convert.ToInt32Ptr(4096)
	default:
		err = errors.Wrapf(errors.InvalidInput, "Invalid key size -- INVESTIGATE BUG -- ")
	}
	return
}

func getMOCKeyOperations(ops *[]keyvault.JSONWebKeyOperation) (kops []wssdcloudcommon.KeyOperation, err error) {
	if ops == nil {
		err = errors.Wrapf(errors.InvalidInput, "")
		return
	}

	kops = []wssdcloudcommon.KeyOperation{}

	for _, k := range *ops {
		tmp, err1 := getMOCKeyOperation(k)
		if err1 != nil {
			err = err1
			return
		}
		kops = append(kops, tmp)
	}
	return
}

func getMOCKeyOperation(op keyvault.JSONWebKeyOperation) (ko wssdcloudcommon.KeyOperation, err error) {
	switch op {
	case keyvault.Encrypt:
		ko = wssdcloudcommon.KeyOperation_ENCRYPT
	case keyvault.Decrypt:
		ko = wssdcloudcommon.KeyOperation_DECRYPT
	case keyvault.WrapKey:
		ko = wssdcloudcommon.KeyOperation_WRAPKEY
	case keyvault.UnwrapKey:
		ko = wssdcloudcommon.KeyOperation_UNWRAPKEY
	default:
		err = errors.Wrapf(errors.InvalidInput, "Invalid KeyOperation")
	}
	return
}

func getKeyOperations(kops []wssdcloudcommon.KeyOperation) (ops *[]keyvault.JSONWebKeyOperation, err error) {
	tmp := []keyvault.JSONWebKeyOperation{}
	for _, k := range kops {
		tmpkey, err1 := getKeyOperation(k)
		if err1 != nil {
			err = err1
			return
			// Something is wrong. Investigate
		}
		tmp = append(tmp, tmpkey)
	}

	ops = &tmp
	return
}

func getKeyOperation(ko wssdcloudcommon.KeyOperation) (op keyvault.JSONWebKeyOperation, err error) {
	switch ko {
	case wssdcloudcommon.KeyOperation_ENCRYPT:
		op = keyvault.Encrypt
	case wssdcloudcommon.KeyOperation_DECRYPT:
		op = keyvault.Decrypt
	case wssdcloudcommon.KeyOperation_WRAPKEY:
		op = keyvault.WrapKey
	case wssdcloudcommon.KeyOperation_UNWRAPKEY:
		op = keyvault.UnwrapKey
	default:
		err = errors.Wrapf(errors.InvalidInput, "--- INVESTIGATE BUG --- ")
	}
	return
}

func getMOCAlgorithm(algo keyvault.JSONWebKeyEncryptionAlgorithm) (wssdcloudcommon.Algorithm, error) {
	switch algo {
	case keyvault.RSA15:
		return wssdcloudcommon.Algorithm_RSA15, nil
	case keyvault.RSAOAEP:
		return wssdcloudcommon.Algorithm_RSAOAEP, nil
	case keyvault.RSAOAEP256:
		return wssdcloudcommon.Algorithm_RSAOAEP256, nil
	case keyvault.A256KW:
		return wssdcloudcommon.Algorithm_A256KW, nil
	}
	return wssdcloudcommon.Algorithm_A_UNKNOWN, errors.Wrapf(errors.InvalidInput, "Invalid Algorithm [%s]", algo)
}

func GetMOCAlgorithmType(algo string) (keyvault.JSONWebKeyEncryptionAlgorithm, error) {
	switch algo {
	case "RSA1_5":
		return keyvault.RSA15, nil
	case "RSA-OAEP":
		return keyvault.RSAOAEP, nil
	case "RSA-OAEP-256":
		return keyvault.RSAOAEP256, nil
	case "A-256-KW":
		return keyvault.A256KW, nil
	}
	return keyvault.RSA15, errors.Wrapf(errors.InvalidInput, "Invalid Algorithm [%s]", algo)
}
