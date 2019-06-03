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

package noobaa

import (
	"github.com/rook/rook/cmd/rook/rook"
	"github.com/rook/rook/pkg/operator/noobaa"
	"github.com/rook/rook/pkg/util/flags"
	"github.com/spf13/cobra"
)

var (
	operatorCmd = &cobra.Command{
		Use:   "operator",
		Short: "Runs rook-noobaa-operator in a kubernetes cluster",
		Long:  "See https://github.com/noobaa/noobaa-core and https://github.com/rook/rook",
	}
)

func init() {
	flags.SetFlagsFromEnv(operatorCmd.Flags(), rook.RookEnvVarPrefix)
	flags.SetLoggingFlags(operatorCmd.Flags())
	operatorCmd.RunE = startOperator
}

func startOperator(cmd *cobra.Command, args []string) error {

	rook.SetLogLevel()
	rook.LogStartupInfo(operatorCmd.Flags())

	logger.Info("Starting noobaa operator ...")

	context := rook.NewContext()

	operator, err := noobaa.NewOperator(context)
	rook.TerminateIfError("Failed to create operator. %+v", err)

	err = operator.Run()
	rook.TerminateIfError("Failed to run operator. %+v", err)

	return nil
}
