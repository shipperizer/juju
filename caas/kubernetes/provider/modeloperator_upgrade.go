// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package provider

import (
	"github.com/juju/juju/api"
	"github.com/juju/names/v4"
	"github.com/juju/version/v2"
	"k8s.io/client-go/kubernetes"
)

type upgradeCAASModelOperatorBridge struct {
	clientFn    func() kubernetes.Interface
	namespaceFn func() string
	isLegacyFn  func() bool
}

type UpgradeCAASModelOperatorBroker interface {
	// Client returns a Kubernetes client associated with the current broker's
	// cluster
	Client() kubernetes.Interface

	// IsLegacyLabels indicates if this provider is operating on a legacy label schema
	IsLegacyLabels() bool

	// Namespace returns the targeted Kubernetes namespace for this broker
	Namespace() string
}

func (u *upgradeCAASModelOperatorBridge) Client() kubernetes.Interface {
	return u.clientFn()
}

func (u *upgradeCAASModelOperatorBridge) IsLegacyLabels() bool {
	return u.isLegacyFn()
}

func modelOperatorUpgrade(
	operatorName string,
	vers version.Number,
	broker UpgradeCAASModelOperatorBroker) error {
	return upgradeDeployment(
		operatorName,
		"",
		vers,
		broker.IsLegacyLabels(),
		api.NewClient(broker).AppsV1().Deployments(broker.Namespace()))
}

func (u *upgradeCAASModelOperatorBridge) Namespace() string {
	return u.namespaceFn()
}

func (k *kubernetesClient) upgradeModelOperator(agentTag names.Tag, vers version.Number) error {
	broker := &upgradeCAASModelOperatorBridge{
		clientFn:    k.client,
		namespaceFn: k.GetCurrentNamespace,
		isLegacyFn:  k.IsLegacyLabels,
	}
	return modelOperatorUpgrade(modelOperatorName, vers, broker)
}
