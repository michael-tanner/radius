// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package datamodel

import (
	v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
)

// EnvironmentComputeKind is the type of compute resource.
type EnvironmentComputeKind string

const (
	// UnknownComputeKind represents kubernetes compute resource type.
	UnknownComputeKind EnvironmentComputeKind = "unknown"
	// KubernetesComputeKind represents kubernetes compute resource type.
	KubernetesComputeKind EnvironmentComputeKind = "kubernetes"
)

// Environment represents Application environment resource.
type Environment struct {
	v1.BaseResource

	// Properties is the properties of the resource.
	Properties EnvironmentProperties `json:"properties"`
}

func (e *Environment) ResourceTypeName() string {
	return "Applications.Core/environments"
}

// EnvironmentProperties represents the properties of Environment.
type EnvironmentProperties struct {
	Compute   EnvironmentCompute                     `json:"compute,omitempty"`
	Recipes   map[string]EnvironmentRecipeProperties `json:"recipes,omitempty"`
	Providers ProviderProperties                     `json:"providers,omitempty"`
}

// EnvironmentCompute represents the compute resource of Environment.
type EnvironmentCompute struct {
	Kind              EnvironmentComputeKind      `json:"kind"`
	KubernetesCompute KubernetesComputeProperties `json:"kubernetes,omitempty"`
}

// KubernetesComputeProperties represents the kubernetes compute of the environment.
type KubernetesComputeProperties struct {
	ResourceID string `json:"resourceId,omitempty"`
	Namespace  string `json:"namespace"`
}

// EnvironmentRecipeProperties represents the properties of environment's recipe.
type EnvironmentRecipeProperties struct {
	ConnectorType string `json:"connectorType,omitempty"`
	TemplatePath  string `json:"templatePath,omitempty"`
}

// ProviderProperties represents the provider, eg azure
type ProviderProperties struct {
	Azure ProviderPropertiesAzure `json:"azure,omitempty"`
}

// ProviderPropertiesAzure represents the azure provider properties
type ProviderPropertiesAzure struct {
	Scope string `json:"scope,omitempty"`
}
