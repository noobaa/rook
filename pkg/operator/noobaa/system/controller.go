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

// Package system to control a NooBaa System
package system

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/pkg/capnslog"
	noobaav1 "github.com/rook/rook/pkg/apis/noobaa.rook.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

//
// NooBaa data-service CRD definitions
//

// SystemCRD is the noobaa data-service CRD
var SystemCRD = &extv1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: "systems." + noobaav1.CustomResourceGroup,
	},
	Spec: extv1.CustomResourceDefinitionSpec{
		Scope: extv1.NamespaceScoped,
		Group: noobaav1.CustomResourceGroup,
		Names: extv1.CustomResourceDefinitionNames{
			Kind:       "System",
			ListKind:   "SystemList",
			Singular:   "system",
			Plural:     "systems",
			ShortNames: []string{"nb", "noobaa", "ds"},
		},
		Version: noobaav1.Version,
		Versions: []extv1.CustomResourceDefinitionVersion{{
			Name:    noobaav1.Version,
			Served:  true,
			Storage: true,
		}},
	},
}

var logger = capnslog.NewPackageLogger("github.com/rook/rook", "noobaa-controller")

// Controller represents a controller object for noobaa custom resources
type Controller struct {
	Manager manager.Manager
	// This client, initialized using mgr.Client(), is a split client
	// that reads objects from the cache and writes to the apiserver
	Client client.Client
}

// NewController create controller for watching noobaa custom resources created
func NewController(mgr manager.Manager) *Controller {
	return &Controller{
		Manager: mgr,
		Client:  mgr.GetClient(),
	}
}

func (c *Controller) Start(stopChan <-chan struct{}) error {

	crdCreated := false
	for retries := 0; retries <= 10; retries++ {
		err := c.initCRD()
		if err != nil {
			logger.Warning("Start() Warning - still initializing CRD -", err)
			time.Sleep(time.Second)
		} else {
			crdCreated = true
		}
	}
	if !crdCreated {
		return fmt.Errorf("Start() ERROR - CRD create exhausted: %s", SystemCRD.ObjectMeta.Name)
	}

	// Create a new controller
	ctrl, err := controller.New("noobaa-controller", c.Manager, controller.Options{Reconciler: c})
	if err != nil {
		return err
	}

	// Watch for changes
	err = ctrl.Watch(&source.Kind{Type: &noobaav1.System{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO: Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Noobaa
	err = ctrl.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &noobaav1.System{},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) initCRD() error {
	ctx := context.TODO()

	err := c.Client.Create(ctx, SystemCRD.DeepCopyObject())
	if err != nil {
		if errors.IsAlreadyExists(err) {
			logger.Info("initCRD() OK - CRD already exists")
		} else {
			return err
		}
	} else {
		logger.Info("initCRD() OK - CRD created")
	}

	objkey, err := client.ObjectKeyFromObject(SystemCRD)
	if err != nil {
		return err
	}

	crd := &extv1.CustomResourceDefinition{}
	err = c.Client.Get(ctx, objkey, crd)
	logger.Info("GGG: ", err, crd)
	if err != nil {
		return err
	}
	for _, cond := range crd.Status.Conditions {
		switch cond.Type {
		case extv1.NamesAccepted:
			if cond.Status == extv1.ConditionFalse {
				return fmt.Errorf("initCRD() ERROR - name conflict: %v", cond.Reason)
			}
		case extv1.Established:
			if cond.Status == extv1.ConditionTrue {
				logger.Info("initCRD() OK - CRD is ready ", cond)
				return nil // success
			}
		}
	}
	return fmt.Errorf("initCRD() Warning - CRD not yet ready")
}

// Reconcile watches for instances of NooBaa custom resources and acts on them
func (c *Controller) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	logger.Info("Reconcile(): ", req)

	return reconcile.Result{}, nil
}
