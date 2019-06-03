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
	noobaav1 "github.com/rook/rook/pkg/apis/noobaa.rook.io/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// State keeps the desired/current state of all the NooBaa resources
type State struct {
	System       *noobaav1.System
	Account      *corev1.ServiceAccount
	Role         *rbacv1.Role
	RoleBinding  *rbacv1.RoleBinding
	CoreApp      *appsv1.StatefulSet
	AgentsApp    *appsv1.StatefulSet
	S3Service    *corev1.Service
	MgmtService  *corev1.Service
	StorageClass *storagev1.StorageClass
}

func NewState(ds *noobaav1.System) *State {
	ownerIsController := true
	ownerRefs := []metav1.OwnerReference{{
		APIVersion: noobaav1.Version,
		Kind:       SystemCRD.Spec.Names.Kind,
		Name:       ds.ObjectMeta.Name,
		UID:        ds.UID,
		Controller: &ownerIsController,
	}}
	labels := map[string]string{"app": "noobaa"}
	makeObjectMeta := func(name string) metav1.ObjectMeta {
		return metav1.ObjectMeta{
			Name:            name,
			Namespace:       ds.Namespace,
			OwnerReferences: ownerRefs,
			Labels:          labels,
		}
	}
	makeInt32Ptr := func(x int32) *int32 {
		return &x
	}
	s := &State{
		System: ds.DeepCopy(),

		Account: &corev1.ServiceAccount{
			ObjectMeta: makeObjectMeta("noobaa-account"),
		},

		Role: &rbacv1.Role{
			ObjectMeta: makeObjectMeta("noobaa-role"),
			Rules: []rbacv1.PolicyRule{{
				APIGroups: []string{"apps"},
				Resources: []string{"statefulsets"},
				Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
			}},
		},

		RoleBinding: &rbacv1.RoleBinding{
			ObjectMeta: makeObjectMeta("noobaa-role-binding"),
			Subjects: []rbacv1.Subject{{
				Kind: "ServiceAccount",
				Name: "noobaa-account",
			}},
			RoleRef: rbacv1.RoleRef{
				APIGroup: rbacv1.GroupName,
				Kind:     "Role",
				Name:     "noobaa-role",
			},
		},

		CoreApp: &appsv1.StatefulSet{
			ObjectMeta: makeObjectMeta("noobaa-core"),
			Spec: appsv1.StatefulSetSpec{
				Replicas: makeInt32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"noobaa-module": "noobaa-core",
					},
				},
				ServiceName:          "noobaa-mgmt",
				VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{}},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{},
					Spec: corev1.PodSpec{
						ServiceAccountName: "noobaa-account",
						Hostname:           "noobaa-core",
						Containers: []corev1.Container{{
							Name: "noobaa-core",
						}},
					},
				},
			},
		},

		AgentsApp: &appsv1.StatefulSet{
			ObjectMeta: makeObjectMeta("noobaa-agents"),
			Spec: appsv1.StatefulSetSpec{
				Replicas: makeInt32Ptr(3),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"noobaa-module": "noobaa-agent"},
				},
			},
		},

		S3Service: &corev1.Service{
			ObjectMeta: makeObjectMeta("s3"),
			Spec:       corev1.ServiceSpec{},
		},

		MgmtService: &corev1.Service{
			ObjectMeta: makeObjectMeta("noobaa-mgmt"),
			Spec:       corev1.ServiceSpec{},
		},

		// StorageClass: &storagev1.StorageClass{
		// 	ObjectMeta: makeObjectMeta("noobaa-storage-class"),
		// },
	}
	return s
}

// func (c *State) load() error {
// 	logger.Println(nbvers.Version)

// 	// Read config yaml packaged in the operator image
// 	data, err := ioutil.ReadFile("noobaa_core.yaml")
// 	if err != nil {
// 		return err
// 	}

// 	textDefs := strings.Split(string(data), "---")
// 	c.subs = make(map[string]*subState)

// 	for _, textDef := range textDefs {
// 		// Decode text (yaml/json) to kube api object
// 		def, group, err := scheme.Codecs.UniversalDeserializer().Decode([]byte(textDef), nil, nil)
// 		if err != nil {
// 			return err
// 		}

// 		// not sure if really needed, but set it anyway
// 		def.GetObjectKind().SetGroupVersionKind(*group)

// 		err = c.addSub(def)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (c *noobaaState) addSub(def runtime.Object) error {

// 	// convert the runtime.Object to unstructured.Unstructured
// 	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(def)
// 	if err != nil {
// 		return err
// 	}
// 	unstructuredDef := &unstructured.Unstructured{Object: unstructuredMap}

// 	gvk := def.GetObjectKind().GroupVersionKind()
// 	mapping, err := c.controller.mapper.RESTMapping(schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}, gvk.Version)
// 	if err != nil {
// 		return err
// 	}

// 	var resourceIfc dynamic.ResourceInterface
// 	nsResourceIfc := c.controller.context.DynamicClient.Resource(mapping.Resource)
// 	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
// 		resourceIfc = nsResourceIfc.Namespace(c.def.Namespace)
// 	} else {
// 		resourceIfc = nsResourceIfc
// 	}

// 	name := mapping.Resource.String()

// 	c.subs[name] = &subState{
// 		name:            name,
// 		def:             def,
// 		unstructuredDef: unstructuredDef,
// 		dynamicClient:   resourceIfc,
// 	}

// 	logger.Info("addSub(): ", c.subs[name])

// 	return nil
// }
