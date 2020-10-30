package deprovisioning

import (
	"time"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/process"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/avs"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/storage"
	"github.com/sirupsen/logrus"
)

type AvsEvaluationRemovalStep struct {
	delegator             *avs.Delegator
	operationsStorage     storage.Operations
	externalEvalAssistant avs.EvalAssistant
	internalEvalAssistant avs.EvalAssistant
	deProvisioningManager *process.DeprovisionOperationManager
}

func NewAvsEvaluationsRemovalStep(delegator *avs.Delegator, operationsStorage storage.Operations, externalEvalAssistant,
	internalEvalAssistant avs.EvalAssistant, log logrus.FieldLogger) *AvsEvaluationRemovalStep {
	return &AvsEvaluationRemovalStep{
		delegator:             delegator,
		operationsStorage:     operationsStorage,
		externalEvalAssistant: externalEvalAssistant,
		internalEvalAssistant: internalEvalAssistant,
		deProvisioningManager: process.NewDeprovisionOperationManager(operationsStorage, log),
	}
}

func (ars *AvsEvaluationRemovalStep) Name() string {
	return "De-provision_AVS_Evaluations"
}

func (ars *AvsEvaluationRemovalStep) Run(deProvisioningOperation internal.DeprovisioningOperation, opLog logrus.FieldLogger) (internal.DeprovisioningOperation, time.Duration, error) {
	opLog.Infof("Avs lifecycle %+v", deProvisioningOperation.Avs)
	if deProvisioningOperation.Avs.AVSExternalEvaluationDeleted && deProvisioningOperation.Avs.AVSInternalEvaluationDeleted {
		opLog.Infof("Both internal and external evaluations have been deleted")
		return deProvisioningOperation, 0, nil
	}

	deProvisioningOperation, err := ars.delegator.DeleteAvsEvaluation(deProvisioningOperation, opLog, ars.internalEvalAssistant)
	if err != nil {
		return ars.deProvisioningManager.RetryOperationWithoutFail(deProvisioningOperation, err.Error(), 10*time.Second, 10*time.Minute)
	}

	deProvisioningOperation, err = ars.delegator.DeleteAvsEvaluation(deProvisioningOperation, opLog, ars.externalEvalAssistant)
	if err != nil {
		return ars.deProvisioningManager.RetryOperationWithoutFail(deProvisioningOperation, err.Error(), 10*time.Second, 10*time.Minute)
	}
	return deProvisioningOperation, 0, nil

}
