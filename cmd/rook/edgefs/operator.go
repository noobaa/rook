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

package edgefs

import (
	"fmt"

	"github.com/rook/rook/cmd/rook/rook"
	edgefsoperator "github.com/rook/rook/pkg/operator/edgefs"
	"github.com/rook/rook/pkg/operator/k8sutil"
	"github.com/rook/rook/pkg/util/flags"
	"github.com/spf13/cobra"
)

const containerName = "rook-edgefs-operator"

var operatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Runs the EdgeFS operator to deploy and manage in kubernetes clusters",
	Long: `Runs the EdgeFS operator to deploy and manage in kubernetes clusters.
https://github.com/rook/rook`,
}

func init() {
	flags.SetFlagsFromEnv(operatorCmd.Flags(), rook.RookEnvVarPrefix)

	operatorCmd.RunE = startOperator
}

func startOperator(cmd *cobra.Command, args []string) error {
	rook.SetLogLevel()
	rook.LogStartupInfo(operatorCmd.Flags())

	logger.Infof("Starting EdgeFS operator")
	context := rook.NewContext()

	// Using the current image version to deploy other rook pods
	pod, err := k8sutil.GetRunningPod(context.Clientset)
	if err != nil {
		rook.TerminateFatal(fmt.Errorf("failed to get pod. %+v\n", err))
	}

	rookImage, err := k8sutil.GetContainerImage(pod, containerName)
	if err != nil {
		rook.TerminateFatal(fmt.Errorf("failed to get container image. %+v\n", err))
	}

	op := edgefsoperator.New(context, rookImage, pod.Spec.ServiceAccountName)
	err = op.Run()
	if err != nil {
		rook.TerminateFatal(fmt.Errorf("failed to run operator. %+v\n", err))
	}

	return nil
}
