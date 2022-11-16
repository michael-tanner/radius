//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package v20220901privatepreview

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// PlanesClient contains the methods for the Planes group.
// Don't use this type directly, use NewPlanesClient() instead.
type PlanesClient struct {
	host string
	pl runtime.Pipeline
}

// NewPlanesClient creates a new instance of PlanesClient with the specified values.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewPlanesClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*PlanesClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &PlanesClient{
		host: ep,
pl: pl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update a Plane.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-09-01-privatepreview
// planeType - The type of the plane
// planeName - The name of the plane
// plane - plane details
// options - PlanesClientCreateOrUpdateOptions contains the optional parameters for the PlanesClient.CreateOrUpdate method.
func (client *PlanesClient) CreateOrUpdate(ctx context.Context, planeType string, planeName string, plane PlaneResource, options *PlanesClientCreateOrUpdateOptions) (PlanesClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, planeType, planeName, plane, options)
	if err != nil {
		return PlanesClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return PlanesClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return PlanesClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *PlanesClient) createOrUpdateCreateRequest(ctx context.Context, planeType string, planeName string, plane PlaneResource, options *PlanesClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/planes/{PlaneType}/{PlaneName}"
	if planeType == "" {
		return nil, errors.New("parameter planeType cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{PlaneType}", url.PathEscape(planeType))
	if planeName == "" {
		return nil, errors.New("parameter planeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{PlaneName}", url.PathEscape(planeName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, plane)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *PlanesClient) createOrUpdateHandleResponse(resp *http.Response) (PlanesClientCreateOrUpdateResponse, error) {
	result := PlanesClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PlaneResource); err != nil {
		return PlanesClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a Plane.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-09-01-privatepreview
// planeType - The type of the plane
// planeName - The name of the plane
// options - PlanesClientDeleteOptions contains the optional parameters for the PlanesClient.Delete method.
func (client *PlanesClient) Delete(ctx context.Context, planeType string, planeName string, options *PlanesClientDeleteOptions) (PlanesClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, planeType, planeName, options)
	if err != nil {
		return PlanesClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return PlanesClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return PlanesClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return PlanesClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *PlanesClient) deleteCreateRequest(ctx context.Context, planeType string, planeName string, options *PlanesClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/planes/{PlaneType}/{PlaneName}"
	if planeType == "" {
		return nil, errors.New("parameter planeType cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{PlaneType}", url.PathEscape(planeType))
	if planeName == "" {
		return nil, errors.New("parameter planeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{PlaneName}", url.PathEscape(planeName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the properties of a UCP Plane.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-09-01-privatepreview
// planeType - The type of the plane
// planeName - The name of the plane
// options - PlanesClientGetOptions contains the optional parameters for the PlanesClient.Get method.
func (client *PlanesClient) Get(ctx context.Context, planeType string, planeName string, options *PlanesClientGetOptions) (PlanesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, planeType, planeName, options)
	if err != nil {
		return PlanesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return PlanesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return PlanesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *PlanesClient) getCreateRequest(ctx context.Context, planeType string, planeName string, options *PlanesClientGetOptions) (*policy.Request, error) {
	urlPath := "/planes/{PlaneType}/{PlaneName}"
	if planeType == "" {
		return nil, errors.New("parameter planeType cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{PlaneType}", url.PathEscape(planeType))
	if planeName == "" {
		return nil, errors.New("parameter planeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{PlaneName}", url.PathEscape(planeName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PlanesClient) getHandleResponse(resp *http.Response) (PlanesClientGetResponse, error) {
	result := PlanesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PlaneResource); err != nil {
		return PlanesClientGetResponse{}, err
	}
	return result, nil
}

// List - List all planes.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-09-01-privatepreview
// options - PlanesClientListOptions contains the optional parameters for the PlanesClient.List method.
func (client *PlanesClient) List(ctx context.Context, options *PlanesClientListOptions) (PlanesClientListResponse, error) {
	req, err := client.listCreateRequest(ctx, options)
	if err != nil {
		return PlanesClientListResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return PlanesClientListResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return PlanesClientListResponse{}, runtime.NewResponseError(resp)
	}
	return client.listHandleResponse(resp)
}

// listCreateRequest creates the List request.
func (client *PlanesClient) listCreateRequest(ctx context.Context, options *PlanesClientListOptions) (*policy.Request, error) {
	urlPath := "/planes"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *PlanesClient) listHandleResponse(resp *http.Response) (PlanesClientListResponse, error) {
	result := PlanesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PlaneResourceList); err != nil {
		return PlanesClientListResponse{}, err
	}
	return result, nil
}

