# Cloud SDK

moc-sdk-for-go contains Go packages for creating and managing cloud resources on the AzureStackHCI cloud.

## Debug Mode

If you would like to test changes without cert generation or identity setup, please set environment variable `WSSD_DEBUG_MODE` to `on`

### Windows

To do this on windows open an admin cmd. Then run

`setx WSSD_DEBUG_MODE "on"`

Close this window and open a new cmd. Ensure the proper value is set with 

`set`

### Linux

Run

`export WSSD_DEBUG_MODE=on`


## Building wssdcloudctl

### Linux

First build the vendor.

`$ make vendor`

Then build `wssdcloudctl.exe`

`$ make`

You should see `wssdcloudctl.exe` in the `./bin` folder

```
$ tree bin

bin
└── wssdcloudctl.exe

```

## Getting Started and Running Pester Tests

### Prerequisites

* Physical machine or VM with nested virtualization turned on.
* Containers Role and Hyper-V Enabled
* Copy `wssdcloudagent.exe`, `wssdagent.exe`, `wssdcloudctl.exe` and `.\test\pester` to the test machine 

### Starting the Node Agent and the Cloud Agent

To run the node agent, copy it to a test machine and execute it from the command line.

`PS C:\>.\wssdagent.exe`

From another shell run the cloud agent

`PS C:\>.\wssdcloudagent.exe`

Add `wssdcloudctl.exe` to the `PATH`

`PS C:\>cp .\wssdcloudctl.exe C:\Windows\System32`

### Running the Tests

In `.\test\pester` run `testwssd.ps1` to execute the pester tests

This should also generate sample `.yaml` needed to run commands manually.

For example the following is the workflow for creating a VM:

```
wssdcloudctl.exe cloud node create --config c:\wssd\samplenode.yaml
wssdcloudctl.exe cloud group create --config c:\wssd\samplegroup.yaml
wssdcloudctl.exe network virtualnetwork create --config c:\wssd\samplevirtualnetwork.yaml --group testgroup
wssdcloudctl.exe network networkinterface create --config c:\wssd\samplenetworkinterface.yaml --group testgroup
wssdcloudctl.exe storage vhd create --config c:\wssd\samplevhd.yaml --group testgroup
wssdcloudctl.exe compute virtualmachine create --config c:\wssd\samplevirtualmachine.yaml --group testgroup
```

Checked into the `.\test\pester` directory is `test.vhdx`. This file is just a shell used for testing. 
If you wanted to create a usable VM you will need to use a real vhdx.

The easiest way to get a Linux vhdx is to copy down the one from netapplinux.

`PS C:\>Start-BitsTransfer http://netapplinux/AzureEdge/Ubuntu2.vhdx`

Then edit `samplevhd.yaml` to point to the the Ubuntu vhdx.

# Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
