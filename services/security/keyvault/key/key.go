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

func getKey(sec *wssdcloudsecurity.Key, vaultName string, getCustomValue func(*wssdcloudsecurity.Key) (string, error)) (keyvault.Key, error) {
	keysize, err := getKeySize(sec.Size)
	if err != nil {
		return keyvault.Key{}, errors.Wrapf(err, "INVESTIAGE")
	}

	key := keyvault.Key{
		ID:      &sec.Id,
		Name:    &sec.Name,
		Version: &sec.Status.Version.Number,
		KeyProperties: &keyvault.KeyProperties{
			Statuses:                      status.GetStatuses(sec.GetStatus()),
			KeyType:                       getKeyType(sec.Type),
			KeySize:                       keysize,
			KeyRotationFrequencyInSeconds: &sec.KeyRotationFrequencyInSeconds,
		},
	}
	value := ""
	switch sec.Type {
	case wssdcloudcommon.JsonWebKeyType_RSA:
		fallthrough
	case wssdcloudcommon.JsonWebKeyType_RSA_HSM:
		// Allow callers to this function to choose how to construct the value from the returned key
		if getCustomValue != nil {
			value, err = getCustomValue(sec)
			if err != nil {
				return keyvault.Key{}, errors.Wrapf(err, "Failed to create custom value from returned key")
			}
		} else { // Default value is pem public key
			pubBlk := pem.Block{
				Type:    "RSA PUBLIC KEY",
				Headers: nil,
				Bytes:   sec.PublicKey,
			}
			value = string(pem.EncodeToMemory(&pubBlk))
		}
	case wssdcloudcommon.JsonWebKeyType_AES:
		// The only time we should be getting a customvalue is in the export case.
		// If we do its not actually public or private data its the wrapped key.
		if getCustomValue != nil {
			value, err = getCustomValue(sec)
			if err != nil {
				return keyvault.Key{}, errors.Wrapf(err, "Failed to create custom value from returned key")
			}
		}
	}
	key.Value = &value
	return key, nil
}

// keyID optional in getWssdKeyByVaultName function
func getWssdKeyByVaultName(name string, groupName,
	vaultName, keyID string, opType wssdcloudcommon.Operation) (*wssdcloudsecurity.Key, error) {
	key := &wssdcloudsecurity.Key{
		Name:       name,
		VaultName:  vaultName,
		GroupName:  groupName,
		Type:       wssdcloudcommon.JsonWebKeyType_EC,
		Size:       wssdcloudcommon.KeySize_K_UNKNOWN,
		KeyOps:     []wssdcloudcommon.KeyOperation{},
		Status:     status.InitStatus(),
		KeyVersion: keyID,
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

	var keyRotationValue int64
	keyRotationValue = 0
	if sec.KeyRotationFrequencyInSeconds != nil {
		keyRotationValue = *sec.KeyRotationFrequencyInSeconds
	}

	key := &wssdcloudsecurity.Key{
		Name:                          name,
		VaultName:                     vaultName,
		GroupName:                     groupName,
		Type:                          getMOCKeyType(sec.KeyType),
		Size:                          keysize,
		KeyOps:                        []wssdcloudcommon.KeyOperation{},
		Status:                        status.InitStatus(),
		KeyRotationFrequencyInSeconds: keyRotationValue,
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
	case keyvault.A256CBC:
		return wssdcloudcommon.Algorithm_A256CBC, nil
	}
	return wssdcloudcommon.Algorithm_A_UNKNOWN, errors.Wrapf(errors.InvalidInput, "Invalid Algorithm [%s]", algo)
}

func getMOCSigningAlgorithm(algo keyvault.JSONWebKeySignatureAlgorithm) (wssdcloudcommon.JSONWebKeySignatureAlgorithm, error) {
	switch algo {
	case keyvault.ES256:
		return wssdcloudcommon.JSONWebKeySignatureAlgorithm_ES256, nil
	case keyvault.ES256K:
		return wssdcloudcommon.JSONWebKeySignatureAlgorithm_ES256K, nil
	case keyvault.ES384:
		return wssdcloudcommon.JSONWebKeySignatureAlgorithm_ES384, nil
	case keyvault.ES512:
		return wssdcloudcommon.JSONWebKeySignatureAlgorithm_ES512, nil
	case keyvault.PS256:
		return wssdcloudcommon.JSONWebKeySignatureAlgorithm_PS256, nil
	case keyvault.PS384:
		return wssdcloudcommon.JSONWebKeySignatureAlgorithm_PS384, nil
	case keyvault.PS512:
		return wssdcloudcommon.JSONWebKeySignatureAlgorithm_PS512, nil
	case keyvault.RS256:
		return wssdcloudcommon.JSONWebKeySignatureAlgorithm_RS256, nil
	case keyvault.RS384:
		return wssdcloudcommon.JSONWebKeySignatureAlgorithm_RS384, nil
	case keyvault.RS512:
		return wssdcloudcommon.JSONWebKeySignatureAlgorithm_RS512, nil
	}
	return wssdcloudcommon.JSONWebKeySignatureAlgorithm_RSNULL, errors.Wrapf(errors.InvalidInput, "Invalid Algorithm [%s]", algo)
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
	case "A-256-CBC":
		return keyvault.A256CBC, nil
	}
	return keyvault.RSA15, errors.Wrapf(errors.InvalidInput, "Invalid Algorithm [%s]", algo)
}

func GetMOCKeyWrappingAlgorithm(algo keyvault.KeyWrappingAlgorithm) (wrappingAlgo wssdcloudcommon.KeyWrappingAlgorithm, err error) {
	switch algo {
	case keyvault.CKM_RSA_AES_KEY_WRAP:
		wrappingAlgo = wssdcloudcommon.KeyWrappingAlgorithm_CKM_RSA_AES_KEY_WRAP
	case keyvault.RSA_AES_KEY_WRAP_256:
		wrappingAlgo = wssdcloudcommon.KeyWrappingAlgorithm_RSA_AES_KEY_WRAP_256
	case keyvault.RSA_AES_KEY_WRAP_384:
		wrappingAlgo = wssdcloudcommon.KeyWrappingAlgorithm_RSA_AES_KEY_WRAP_384
	default:
		err = errors.Wrapf(errors.InvalidInput, "Invalid Algorithm [%s]", algo)
	}
	return
}

func GetKeyWrappingAlgorithm(algo wssdcloudcommon.KeyWrappingAlgorithm) (wrappingAlgo keyvault.KeyWrappingAlgorithm, err error) {
	switch algo {
	case wssdcloudcommon.KeyWrappingAlgorithm_CKM_RSA_AES_KEY_WRAP:
		wrappingAlgo = keyvault.CKM_RSA_AES_KEY_WRAP
	case wssdcloudcommon.KeyWrappingAlgorithm_RSA_AES_KEY_WRAP_256:
		wrappingAlgo = keyvault.RSA_AES_KEY_WRAP_256
	case wssdcloudcommon.KeyWrappingAlgorithm_RSA_AES_KEY_WRAP_384:
		wrappingAlgo = keyvault.RSA_AES_KEY_WRAP_384
	default:
		err = errors.Wrapf(errors.Failed, "Invalid Algorithm [%s]", algo)
	}
	return
}
