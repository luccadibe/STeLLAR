// MIT License
//
// Copyright (c) 2020 Theodor Amariucai and EASE Lab
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package benchhttp

import (
	"fmt"
	"net/http"
	"stellar/setup"
	"stellar/setup/deployment/connection"
	"stellar/setup/deployment/connection/amazon"
	"testing"

	"github.com/stretchr/testify/require"
)

const randomGatewayID = "uicnaywo3rb3nsci"

func TestCreateAWSRequest(t *testing.T) {
	randomPayloadLength := 7
	randomEndpoint := setup.EndpointInfo{
		ID:                   randomGatewayID,
		DataTransferChainIDs: []string{},
	}

	connection.Initialize("aws", "", "../../../setup/deployment/raw-code/functions/producer-consumer/api-template.json")

	randomAssignedIncrement := int64(1482911482)
	req := CreateRequest("aws", randomPayloadLength, randomEndpoint, randomAssignedIncrement, false, "route1")

	expectedHostname := fmt.Sprintf("%s.execute-api.%s.amazonaws.com", randomEndpoint.ID, amazon.AWSRegion)
	require.Equal(t, expectedHostname, req.Host)
	require.Equal(t, expectedHostname, req.URL.Host)
	require.Equal(t, http.MethodGet, req.Method)
	require.Equal(t, "https", req.URL.Scheme)
}

func TestCreateExternalRequest(t *testing.T) {
	randomPayloadLength := 7
	randomAssignedIncrement := int64(1482911482)
	req := CreateRequest("www.google.com", randomPayloadLength, setup.EndpointInfo{}, randomAssignedIncrement, false, "route1")

	require.Equal(t, "www.google.com", req.Host)
	require.Equal(t, "www.google.com", req.URL.Host)
	require.Equal(t, http.MethodGet, req.Method)
	require.Equal(t, "https", req.URL.Scheme)
}

func TestCreateFlyioRequest(t *testing.T) {
	randomPayloadLength := 7
	randomEndpoint := setup.EndpointInfo{
		ID:                   randomGatewayID,
		DataTransferChainIDs: []string{},
	}

	connection.Initialize("fly.io", "", "../../../setup/deployment/raw-code/functions/producer-consumer/api-template.json")

	randomAssignedIncrement := int64(1482911482)
	req := CreateRequest("fly.io", randomPayloadLength, randomEndpoint, randomAssignedIncrement, false, "route1")

	expectedHostname := fmt.Sprintf("%s.fly.dev", randomEndpoint.ID)
	require.Equal(t, expectedHostname, req.Host)
	require.Equal(t, expectedHostname, req.URL.Host)
	require.Equal(t, http.MethodGet, req.Method)
	require.Equal(t, "https", req.URL.Scheme)
}
