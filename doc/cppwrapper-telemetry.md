# Emit Telemetry from Cpp Wrapper

## Problem
1. Teams leveraging the cpp wrapper in the moc-sdk-for-go (sdk) have no method to investigate failures which occur in the moc stack because no error information is returned from wrapper calls.
2. There is currently no means to correlate events in the moc stack to calls made by outside programs at a per instance level; that is, the ability to look at the telemetry events emited as the result from a specifc call made to the cpp wrapper.

## Goals
1. To provide telemetry about errors which occur after making a call into the cpp wrapper. 
2. To provide the groundwork to support the usage of a correlation vector in the cpp wrapper.
   
## Proposal
1. Instrument the cpp wrapper to emit telemetry about errors received from calling into moc.
2. Refactor the cpp wrapper to include a correlation vector that can be provided by calling processes and emited in telemetry from the wrapper.


Although it would be possible to add error information as part of the returned value in the cpp wrapper, this would require calling programs to parse error information and increase their complexity.

## Design
Initial design will include creating a new function in the wrapper to emit ETW events. This function will leverage the github.com/Microsoft/go-winio/pkg/etw package and follow the structure of other emitters in the moc stack. However, this function will not include any private configuration information. 
The following methods in the cpp wrapper will be instrumented to emit error telemetry:
- SecurityLogin
- KeyvaultKeyEncryptData
- KeyvaultKeyDecryptData
- KeyvaultKeyExist
- KeyvaultKeyCreateOrUpdate
- KeyvaultKeySignData
- KeyvaultKeyVerifyData
- KeyVaultGetPublicKey
- KeyVaultKeyClient

The telemetry will provide a timestamp, the function in the cpp where the error was received, the error string returned from the called subroutine (which method the wrapper called that returned an error), and a correlation vector. For example a telemetry event may look like: 
```ERROR:2023/07/18 10:38:11 main.go:229: KeyvaultKeyDecryptData Failed. Call to keyclient.decrypt failed.  "rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp [::1]:55000: connectex: No connection could be made because the target machine actively refused it.\"" Correlation Vector: ABCDEF12345
```




To provide support for a correlation vector without breaking existing clients that leverage the cpp wrapper, the existing functions will be renamed with CV appended and the parameters will be updated to include a CV.
  E.g.
  
  ```func KeyvaultKeyDecryptData(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, input *C.char, timeoutInSeconds C.int) *C.char``` 
  
  becomes 
  
  ```func KeyvaultKeyDecryptDataCV(serverName *C.char, groupName *C.char, keyvaultName *C.char, keyName *C.char, input *C.char, , CV *C.char,  timeoutInSeconds C.int) *C.char```
  
Then, a wrapper function will be created for each of these functions that has the same name and parameters as the original function prototype. The wrapper can provide a hard coded CV when calling the CV vernsion of the function. This allows existing clients to continue to leverage the wrapper as usual, while new clients can leverage the new functions' support for correlation vectors.


## Changes to other repositories in MOC Stack
None.
