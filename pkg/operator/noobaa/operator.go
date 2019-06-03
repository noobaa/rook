/*
Copyright 2018 The Rook Authors. All rights reserved.

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

// Package noobaa operates data services
package noobaa

import (
	"github.com/coreos/pkg/capnslog"
	"github.com/rook/rook/pkg/client/clientset/versioned/scheme"
	"github.com/rook/rook/pkg/clusterd"
	noobaasystem "github.com/rook/rook/pkg/operator/noobaa/system"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

// Change below  to serve metrics on different host or port.
const (
	metricsHost       = "0.0.0.0"
	metricsPort int32 = 8383
)

var logger = capnslog.NewPackageLogger("github.com/rook/rook", "noobaa-operator")

// Operator type for managing NooBaa data services
type Operator struct {
	Manager    manager.Manager
	Controller *noobaasystem.Controller
}

// NewOperator creates an operator instance
func NewOperator(context *clusterd.Context) (*Operator, error) {

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(context.KubeConfig, manager.Options{
		// Namespace: os.Getenv("WATCH_NAMESPACE"),
		// MetricsBindAddress: fmt.Sprintf("%s:%d", metricsHost, metricsPort),
	})
	if err != nil {
		logger.Error("NewOperator() Create manager ERROR: ", err)
		return nil, err
	}

	err = scheme.AddToScheme(mgr.GetScheme())
	if err != nil {
		logger.Error("NewOperator() AddToScheme ERROR: ", err)
		return nil, err
	}

	// Register controller to run when the manager runs and everything has been initialized
	controller := noobaasystem.NewController(mgr)
	err = mgr.Add(controller)
	if err != nil {
		logger.Error("NewOperator() Add controller ERROR: ", err)
		return nil, err
	}

	o := &Operator{
		Manager:    mgr,
		Controller: controller,
	}

	return o, nil
}

// Run the operator
func (o *Operator) Run() error {

	logger.Info("Run() Starting the manager ...")

	err := o.Manager.Start(signals.SetupSignalHandler())
	if err != nil {
		logger.Error("Run() Manager ERROR: ", err)
		return err
	}

	logger.Info("Run() Manager exited gracefully")

	return nil
}
