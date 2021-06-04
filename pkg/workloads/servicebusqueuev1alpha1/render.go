// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package servicebusqueuev1alpha1

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/servicebus/mgmt/servicebus"
	"github.com/Azure/radius/pkg/curp/armauth"
	"github.com/Azure/radius/pkg/curp/components"
	"github.com/Azure/radius/pkg/curp/handlers"
	"github.com/Azure/radius/pkg/curp/resources"
	"github.com/Azure/radius/pkg/workloads"
)

// Renderer is the WorkloadRenderer implementation for the service bus workload.
type Renderer struct {
	Arm armauth.ArmConfig
}

// Allocate is the WorkloadRenderer implementation for servicebus workload.
func (r Renderer) AllocateBindings(ctx context.Context, workload workloads.InstantiatedWorkload, resources []workloads.WorkloadResourceProperties) (map[string]components.BindingState, error) {
	if len(workload.Workload.Bindings) > 0 {
		return nil, fmt.Errorf("component of kind %s does not support user-defined bindings", Kind)
	}

	if len(resources) != 1 || resources[0].Type != workloads.ResourceKindAzureServiceBusQueue {
		return nil, fmt.Errorf("cannot fulfill binding - expected properties for %s", workloads.ResourceKindAzureServiceBusQueue)
	}

	properties := resources[0].Properties
	namespaceName := properties[handlers.ServiceBusNamespaceNameKey]
	queueName := properties[handlers.ServiceBusQueueNameKey]

	sbClient := servicebus.NewNamespacesClient(r.Arm.SubscriptionID)
	sbClient.Authorizer = r.Arm.Auth
	accessKeys, err := sbClient.ListKeys(ctx, r.Arm.ResourceGroup, namespaceName, "RootManageSharedAccessKey")

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve connection strings: %w", err)
	}

	if accessKeys.PrimaryConnectionString == nil && accessKeys.SecondaryConnectionString == nil {
		return nil, fmt.Errorf("failed to retrieve connection strings")
	}

	cs := accessKeys.PrimaryConnectionString

	bindings := map[string]components.BindingState{
		"default": {
			Component: workload.Name,
			Binding:   "default",
			Kind:      "azure.com/ServiceBusQueue",
			Properties: map[string]interface{}{
				"connectionString": *cs,
				"namespace":        namespaceName,
				"queue":            queueName,
			},
		},
	}

	return bindings, nil
}

// Render is the WorkloadRenderer implementation for servicebus workload.
func (r Renderer) Render(ctx context.Context, w workloads.InstantiatedWorkload) ([]workloads.WorkloadResource, error) {
	component := ServiceBusQueueComponent{}
	err := w.Workload.AsRequired(Kind, &component)
	if err != nil {
		return nil, err
	}

	if component.Config.Managed {
		if component.Config.Queue == "" {
			return nil, errors.New("the 'topic' field is required when 'managed=true'")
		}

		if component.Config.Resource != "" {
			return nil, workloads.ErrResourceSpecifiedForManagedResource
		}

		// generate data we can use to manage a servicebus queue

		resource := workloads.WorkloadResource{
			Type: workloads.ResourceKindAzureServiceBusQueue,
			Resource: map[string]string{
				handlers.ManagedKey:             "true",
				handlers.ServiceBusQueueNameKey: component.Config.Queue,
			},
		}

		// It's already in the correct format
		return []workloads.WorkloadResource{resource}, nil
	} else {
		if component.Config.Resource == "" {
			return nil, workloads.ErrResourceMissingForUnmanagedResource
		}

		queueID, err := workloads.ValidateResourceID(component.Config.Resource, QueueResourceType, "ServiceBus Queue")
		if err != nil {
			return nil, err
		}

		// generate data we can use to connect to a servicebus queue
		resource := workloads.WorkloadResource{
			Type: workloads.ResourceKindAzureServiceBusQueue,
			Resource: map[string]string{
				handlers.ManagedKey: "false",

				// Truncate the queue part of the ID to make an ID for the namespace
				handlers.ServiceBusNamespaceIDKey:   resources.MakeID(queueID.SubscriptionID, queueID.ResourceGroup, queueID.Types[0]),
				handlers.ServiceBusQueueIDKey:       queueID.ID,
				handlers.ServiceBusNamespaceNameKey: queueID.Types[0].Name,
				handlers.ServiceBusQueueNameKey:     queueID.Types[1].Name,
			},
		}

		// It's already in the correct format
		return []workloads.WorkloadResource{resource}, nil
	}
}
