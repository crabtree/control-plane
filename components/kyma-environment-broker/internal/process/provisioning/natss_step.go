package provisioning

import (
	"time"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/runtime/components"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/process"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/storage"
	"github.com/kyma-project/control-plane/components/provisioner/pkg/gqlschema"
	"github.com/sirupsen/logrus"
)

type NatsStreamingStep struct {
	operationManager *process.ProvisionOperationManager
}

// ensure the interface is implemented
var _ Step = (*NatsStreamingStep)(nil)

func NewNatsStreamingOverridesStep(os storage.Operations, log logrus.FieldLogger) *NatsStreamingStep {
	return &NatsStreamingStep{
		operationManager: process.NewProvisionOperationManager(os, log),
	}
}

func (s *NatsStreamingStep) Name() string {
	return "Provision Nats Streaming"
}

func (s *NatsStreamingStep) Run(operation internal.ProvisioningOperation, opLog logrus.FieldLogger) (internal.ProvisioningOperation, time.Duration, error) {
	parameters, err := operation.GetProvisioningParameters()
	if err != nil {
		opLog.Errorf("cannot fetch provisioning parameters from operation: %s", err)
		return s.operationManager.OperationFailed(operation, "invalid operation provisioning parameters")
	}
	opLog.Infof("Provisioning for PlanID: %s", parameters.PlanID)
	operation.InputCreator.AppendOverrides(components.NatsStreaming, getNatsStreamingOverrides())
	return operation, 0, nil
}

func getNatsStreamingOverrides() []*gqlschema.ConfigEntryInput {
	return []*gqlschema.ConfigEntryInput{
		{
			Key:   "global.natsStreaming.persistence.enabled",
			Value: "false",
		},
	}
}
