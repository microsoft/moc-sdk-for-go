// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package internal

import (
	"context"
	"errors"
	"github.com/microsoft/moc-proto/pkg/auth"
	loggingHelpers "github.com/microsoft/moc-proto/pkg/logging"
	wssdadmin "github.com/microsoft/moc-proto/rpc/cloudagent/admin"
	wssdclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"io"
	"os"
	"strconv"
	"time"
)

type client struct {
	wssdadmin.LogAgentClient
}

// NewLoggingClient - creates a client session with the backend wssd agent
func NewLoggingClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetLogClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) GetLogFile(ctx context.Context, location, filename string) error {
	request := getLoggingRequest(location)
	fileStreamClient, err := c.LogAgentClient.Get(ctx, request)
	if err != nil {
		return err
	}

	doneErr := errors.New("done")

	for err == nil {
		filename := "bad.log"
		recFunc := func() ([]byte, error) {
			getLogFileResponse, innerErr := fileStreamClient.Recv()
			if innerErr != nil {
				return []byte{}, innerErr
			}
			if getLogFileResponse.Error == io.EOF.Error() {
				filename = getLogFileResponse.Filename
				return getLogFileResponse.File, io.EOF
			}
			if getLogFileResponse.Done {
				return []byte{}, doneErr
			}
			return getLogFileResponse.File, nil

		}
		tempFilename := strconv.FormatInt(time.Now().Unix(), 16) + ".log"
		err = loggingHelpers.ReceiveFile(ctx, tempFilename, recFunc)
		os.Rename(tempFilename, filename)
	}
	if err != doneErr {
		return err
	}
	return nil
}

func getLoggingRequest(location string) *wssdadmin.LogRequest {
	return &wssdadmin.LogRequest{
		Location: location,
	}
}
