#include <stdlib.h>
#include <stdio.h>
#include <string>
#include <algorithm>

#include <google/protobuf/io/coded_stream.h>
#include <google/protobuf/extension_set.h>
#include <google/protobuf/wire_format_lite.h>
#include <google/protobuf/descriptor.h>
#include <google/protobuf/generated_message_reflection.h>
#include <google/protobuf/reflection_ops.h>
#include <google/protobuf/wire_format.h>
// @@protoc_insertion_point(includes)
#include <google/protobuf/port_def.inc>
// for RpcController
#include <google/protobuf/service.h>

#include "moc_cloudagent_key.pb.h"
#include "moc_common_security.pb.h"

#define E_FAIL ((int)0x80004005)
#define E_INVALID_ARG ((int)0x80070057)
#define SUCCESS 0

using namespace namespace google::protobuf

EXTERN_C
int testEncrypt(const char* serverName, const char* groupName, const char* keyvaultName, const char* keyName, const char* plaintext, size_t plaintextSize, const char* cv, int timeout, char** ciphertext, int* ciphertextSize)
{
	if (serverName == nullptr
		|| groupName == nullptr
		|| keyvaultName == nullptr
		|| keyName == nullptr
		|| plaintext == nullptr
		|| plaintextSize == 0
		|| cv == nullptr
		|| ciphertext == nullptr
		|| ciphertextSize == nullptr)
	{
		return E_INVALID_ARG;
	}

	//
	// Step 1.
	// func getKeyvaultKeyClient(serverName string, cv string) (*key.KeyClient, error) {
	//     authorizer, err : = auth.NewAuthorizerFromEnvironment(serverName)
	//		return key.NewKeyClient(serverName, authorizer)}
	//
	// NewAuthorizerFromEnvironment is implemented in https://github.com/microsoft/moc/blob/31c12f09d373898fbdfc466a0630c365c473d0cb/pkg/auth/auth.go#L256 not available in protobuf
	//
	// newKeyClient creates client session with the backend cloudagent, implemented in sdk layer
	// https://github.com/microsoft/moc-sdk-for-go/blob/main/services/security/keyvault/key/wssd.go#L26
	// https://github.com/microsoft/moc-sdk-for-go/blob/9a2dab9e9aae584d042f9df09f825e5a7a19c63e/pkg/client/security.go#L34
	// calls in getClientConnection() which calls into more functions implemented in sdk layer, and with some cache mechanism to maintain connections:
	// https://github.com/microsoft/moc-sdk-for-go/blob/9a2dab9e9aae584d042f9df09f825e5a7a19c63e/pkg/client/client.go#L100
	//
	
	// 
	// Step 2. implement https://github.com/microsoft/moc-sdk-for-go/blob/9a2dab9e9aae584d042f9df09f825e5a7a19c63e/services/security/keyvault/key/wssd.go#L304
	// Note I believe we need to reimplement this whole file in cpp
	//
	KeyAgent agent{};

	//
	// construct Key
	//
	::moc::cloudagent::security::Key keyinfo{};
	keyinfo.set_name(keyName);
	keyinfo.set_vaultname(keyvaultName);
	keyinfo.set_groupname(groupName);

	//
    //::PROTOBUF_NAMESPACE_ID::RpcController, how to construct this? I think it should contains server name and point to the right cloudagent
    // https://protobuf.dev/reference/cpp/api-docs/google.protobuf.service/#RpcController
    // this need to point to step 1
    //
	::PROTOBUF_NAMESPACE_ID::RpcController controller{};

	//
	// Construct Key request
	// 
	::moc::cloudagent::security::KeyRequest keyRequest{};
	keyRequest.add_keys(); // TODO
	keyRequest.set_operationtype(::moc::cloudagent::common::Key_Encrypt);

	::moc::cloudagent::security::KeyResponse keyResponse{};

	//
	// call invoke to get keys
	agent.Invoke(controller,
		&keyRequest,
		&keyResponse,
		/*::google::protobuf::Closure* done*/);

	// error handling
	const std::string& errstring = keyResponse.error();
	if (errstring.size() > 0)
	{
		// error handling, logging the error
		return E_FAIL;
	}

	//
	// construct KeyOperationRequest
	//
	const ::moc::cloudagent::security::KeyOperationRequest request{};
	request.set_algorithm(::moc::Algorithm::A256CBC);
	request.set_data(input, plaintextSize);
	request.Key = keyResponse.keys; // double check on this

	::moc::cloudagent::security::KeyOperationResponse response{};
	
	agent.Operate(controller,
		&request,
		&response,
		/*::google::protobuf::Closure* done*/);

	errstring = response.error();
	if (errstring.size() > 0)
	{
		// error handling, logging the error
		return E_FAIL;
	}
	
	//
	// check whether the buffer has enough space
	//
	size_t responseSize = response.data().size();
	if (*ciphertextSize >= responseSize)
	{
		*ciphertextSize = responseSize;
		if (0 != std::memcpy_s(*ciphertext, *ciphertextSize, response.data(), responseSize))
		{
			// error handling, logging the error
			return E_FAIL;
		}

		return SUCCESS;
	}
	else
	{
		*ciphertextSize = responseSize;
		return /* find the error code for buffer too small */
	}
}

EXTERN_C
int testDecrypt(const char* serverName, const char* groupName, const char* keyvaultName, const char* keyName, const char* ciphertext, size_t ciphertextSize, const char* cv, int timeout, char** plaintext, int* plaintextSize)
{
}

int main(int argc, char** argv, char** env) {
	return 0;
}