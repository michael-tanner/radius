/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package aws

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	armrpc_rest "github.com/radius-project/radius/pkg/armrpc/rest"
	"github.com/radius-project/radius/pkg/armrpc/rpctest"
	"github.com/radius-project/radius/pkg/components/database"
	"github.com/radius-project/radius/pkg/components/secret"
	"github.com/radius-project/radius/pkg/to"
	"github.com/radius-project/radius/pkg/ucp/api/v20231001preview"
	"github.com/radius-project/radius/test/testutil"

	armrpc_controller "github.com/radius-project/radius/pkg/armrpc/frontend/controller"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_AWS_Credential(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatabaseClient := database.NewMockClient(mockCtrl)
	mockSecretClient := secret.NewMockClient(mockCtrl)

	credentialCtrl, err := NewCreateOrUpdateAWSCredential(armrpc_controller.Options{DatabaseClient: mockDatabaseClient}, mockSecretClient)
	require.NoError(t, err)

	tests := []struct {
		name       string
		filename   string
		headerfile string
		url        string
		expected   armrpc_rest.Response
		fn         func(mockDatabaseClient database.MockClient, mockSecretClient secret.MockClient)
		err        error
	}{
		{
			name:       "test_credential_creation",
			filename:   "aws-credential.json",
			headerfile: testHeaderFile,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=2023-10-01-preview",
			expected:   getAwsResponse(),
			fn:         setupCredentialSuccessMocks,
			err:        nil,
		},
		{
			name:       "test_invalid_version_credential_resource",
			filename:   "aws-credential.json",
			headerfile: testHeaderFileWithBadAPIVersion,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=bad",
			expected:   nil,
			fn:         setupEmptyMocks,
			err:        v1.ErrUnsupportedAPIVersion,
		},
		{
			name:       "test_invalid_credential_request",
			filename:   "invalid-request-aws-credential.json",
			headerfile: testHeaderFile,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=2023-10-01-preview",
			expected:   nil,
			fn:         setupEmptyMocks,
			err: &v1.ErrModelConversion{
				PropertyName: "$.properties",
				ValidValue:   "not nil",
			},
		},
		{
			name:       "test_credential_created",
			filename:   "aws-credential.json",
			headerfile: testHeaderFile,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=2023-10-01-preview",
			expected:   getAwsResponse(),
			fn:         setupCredentialNotFoundMocks,
			err:        nil,
		},
		{
			name:       "test_credential_notFoundError",
			filename:   "aws-credential.json",
			headerfile: testHeaderFile,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=2023-10-01-preview",
			fn:         setupCredentialNotFoundErrorMocks,
			err:        errors.New("Error"),
		},
		{
			name:       "test_credential_get_failure",
			filename:   "aws-credential.json",
			headerfile: testHeaderFile,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=2023-10-01-preview",
			fn:         setupCredentialGetFailMocks,
			err:        errors.New("Failed Get"),
		},
		{
			name:       "test_credential_secret_save_failure",
			filename:   "aws-credential.json",
			headerfile: testHeaderFile,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=2023-10-01-preview",
			fn:         setupCredentialSecretSaveFailMocks,
			err:        errors.New("Secret Save Failure"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn(*mockDatabaseClient, *mockSecretClient)

			credentialVersionedInput := &v20231001preview.AwsCredentialResource{}
			credentialInput := testutil.ReadFixture(tt.filename)
			err = json.Unmarshal(credentialInput, credentialVersionedInput)
			require.NoError(t, err)

			request, err := rpctest.NewHTTPRequestFromJSON(context.Background(), http.MethodPut, tt.headerfile, credentialVersionedInput)
			require.NoError(t, err)

			ctx := rpctest.NewARMRequestContext(request)

			response, err := credentialCtrl.Run(ctx, nil, request)
			if tt.err != nil {
				require.Equal(t, tt.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, response)
			}
		})
	}

}

func getAwsResponse() armrpc_rest.Response {
	return armrpc_rest.NewOKResponseWithHeaders(&v20231001preview.AwsCredentialResource{
		Location: to.Ptr("West US"),
		ID:       to.Ptr("/planes/aws/awscloud/providers/System.AWS/credentials/default"),
		Name:     to.Ptr("default"),
		Type:     to.Ptr("System.AWS/credentials"),
		Tags: map[string]*string{
			"env": to.Ptr("dev"),
		},
		Properties: &v20231001preview.AwsAccessKeyCredentialProperties{
			AccessKeyID: to.Ptr("00000000-0000-0000-0000-000000000000"),
			Kind:        to.Ptr(v20231001preview.AWSCredentialKindAccessKey),
			Storage: &v20231001preview.InternalCredentialStorageProperties{
				Kind:       to.Ptr(v20231001preview.CredentialStorageKindInternal),
				SecretName: to.Ptr("aws-awscloud-default"),
			},
		},
	}, map[string]string{"ETag": ""})
}

func setupCredentialSuccessMocks(mockDatabaseClient database.MockClient, mockSecretClient secret.MockClient) {
	mockDatabaseClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).DoAndReturn(func(ctx context.Context, id string, _ ...database.GetOptions) (*database.Object, error) {
		return nil, &database.ErrNotFound{ID: id}
	})
	mockSecretClient.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockDatabaseClient.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
}

func setupEmptyMocks(mockDatabaseClient database.MockClient, mockSecretClient secret.MockClient) {
}

func setupCredentialNotFoundMocks(mockDatabaseClient database.MockClient, mockSecretClient secret.MockClient) {
	mockDatabaseClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id string, options ...database.GetOptions) (*database.Object, error) {
			return nil, &database.ErrNotFound{ID: id}
		}).Times(1)
	mockSecretClient.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockDatabaseClient.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
}

func setupCredentialNotFoundErrorMocks(mockDatabaseClient database.MockClient, mockSecretClient secret.MockClient) {
	mockDatabaseClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id string, options ...database.GetOptions) (*database.Object, error) {
			return nil, errors.New("Error")
		}).Times(1)
}

func setupCredentialGetFailMocks(mockDatabaseClient database.MockClient, mockSecretClient secret.MockClient) {
	mockDatabaseClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id string, options ...database.GetOptions) (*database.Object, error) {
			return nil, errors.New("Failed Get")
		}).Times(1)
}

func setupCredentialSecretSaveFailMocks(mockDatabaseClient database.MockClient, mockSecretClient secret.MockClient) {
	mockDatabaseClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(ctx context.Context, id string, options ...database.GetOptions) (*database.Object, error) {
			return nil, &database.ErrNotFound{ID: id}
		}).Times(1)
	mockSecretClient.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Secret Save Failure")).Times(1)
}
