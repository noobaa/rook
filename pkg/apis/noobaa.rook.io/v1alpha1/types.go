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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ***************************************************************************
// IMPORTANT FOR CODE GENERATION
// If the types in this file are updated, you will need to run
// `make codegen` to generate the new types under the client/clientset folder.
// ***************************************************************************

// System is the custom resource describing a NooBaa System
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type System struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SystemSpec   `json:"spec,omitempty"`
	Status            SystemStatus `json:"status,omitempty"`
}

// SystemList is just a list of Systems
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SystemList struct {
	metav1.TypeMeta `json:",inline,omitempty"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []System `json:"items"`
}

///////////////////////////////////////////////////////////////////////////////

// SystemSpec is the spec of a System system
type SystemSpec struct {

	// ServerImage (optional) overrides the default image for server pods
	ServerImage string `json:"serverImage,omitempty"`

	// AgentImage (optional) overrides the default image for agent pods
	AgentImage string `json:"agentImage,omitempty"`

	// ImageVersion (optional) overrides the server and agent image tag to force upgrade
	ImageVersion string `json:"imageVersion,omitempty"`

	// Email (optional) overrides the owner account default email demo@noobaa.com
	Email string `json:"email,omitempty"`

	CloudBuckets []CloudBucketSpec `json:"cloudBuckets,omitempty"`

	// StorageClassName (optional) overrides
	StorageClassName storagev1.StorageClass `json:"storageClass,omitempty"`
}

type CloudBucketSpec struct {
	Name             string `json:"name,omitempty"`
	CloudBucketName  string `json:"cloudBucketName,omitempty"`
	StorageClassName string `json:"storageClassName,omitempty"`
	// Connection      CloudConnection `json:"connection,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////

// SystemStatus defines the observed state of System
type SystemStatus struct {

	// Readme is a message to the user with instructions and information
	Readme string `json:"readme,omitempty"`

	Accounts []AccountStatus `json:"accounts,omitempty"`

	MgmtEndpoints []MgmtEndpoint `json:"mgmtEndpoints,omitempty"`

	S3Endpoints []S3Endpoint `json:"s3Endpoints,omitempty"`
}

// AccountStatus specifies an account to be used for administration of the data service
type AccountStatus struct {
	Email           string                 `json:"email,omitempty"`
	MgmtCredentials corev1.SecretReference `json:"mgmtCredentials,omitempty"`
	S3Credentials   corev1.SecretReference `json:"s3Credentials,omitempty"`
}

// MgmtEndpoint specifies an account to be used for administration of the data service
type MgmtEndpoint struct {
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

// S3Endpoint specifies an account to be used for administration of the data service
type S3Endpoint struct {
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

// CloudConnection
type CloudConnection struct {
	CloudType   string                 `json:"cloudType,omitempty"`
	Endpoint    string                 `json:"endpoint,omitempty"`
	Credentials corev1.SecretReference `json:"credentials,omitempty"`
}
