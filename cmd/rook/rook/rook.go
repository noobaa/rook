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

package rook

import (
	"fmt"
	"os"
	"strings"

	"github.com/coreos/pkg/capnslog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	rookclient "github.com/rook/rook/pkg/client/clientset/versioned"
	"github.com/rook/rook/pkg/clusterd"
	"github.com/rook/rook/pkg/operator/k8sutil"
	"github.com/rook/rook/pkg/util/exec"
	"github.com/rook/rook/pkg/util/flags"
	"github.com/rook/rook/pkg/version"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

const (
	RookEnvVarPrefix = "ROOK"
	terminationLog   = "/dev/termination-log"
)

var RootCmd = &cobra.Command{
	Use: "rook",
}

var (
	logLevelRaw string
	Cfg         = &Config{}
	logger      = capnslog.NewPackageLogger("github.com/rook/rook", "rookcmd")
)

type Config struct {
	LogLevel capnslog.LogLevel
}

// Initialize the configuration parameters. The precedence from lowest to highest is:
//  1) default value (at compilation)
//  2) environment variables (upper case, replace - with _, and rook prefix. For example, discovery-url is ROOK_DISCOVERY_URL)
//  3) command line parameter
func init() {
	RootCmd.PersistentFlags().StringVar(&logLevelRaw, "log-level", "INFO", "logging level for logging/tracing output (valid values: CRITICAL,ERROR,WARNING,NOTICE,INFO,DEBUG,TRACE)")

	// load the environment variables
	flags.SetFlagsFromEnv(RootCmd.Flags(), RookEnvVarPrefix)
	flags.SetFlagsFromEnv(RootCmd.PersistentFlags(), RookEnvVarPrefix)
}

// SetLogLevel set log level based on provided log option.
func SetLogLevel() {
	// parse given log level string then set up corresponding global logging level
	ll, err := capnslog.ParseLevel(logLevelRaw)
	if err != nil {
		logger.Warningf("failed to set log level %s. %+v", logLevelRaw, err)
	}
	Cfg.LogLevel = ll
	capnslog.SetGlobalLogLevel(Cfg.LogLevel)
}

// LogStartupInfo log the version number, arguments, and all final flag values (environment variable overrides have already been taken into account)
func LogStartupInfo(cmdFlags *pflag.FlagSet) {

	flagValues := flags.GetFlagsAndValues(cmdFlags, "secret|keyring")
	logger.Infof("starting Rook %s with arguments '%s'", version.Version, strings.Join(os.Args, " "))
	logger.Infof("flag values: %s", strings.Join(flagValues, ", "))
}

// GetClientset create the k8s client
func GetClientset() (kubernetes.Interface, apiextensionsclient.Interface, rookclient.Interface) {
	config := GetKubeConfig()
	clientset, err := kubernetes.NewForConfig(config)
	TerminateIfError("failed to create k8s clientset. %+v", err)
	apiExtClientset, err := apiextensionsclient.NewForConfig(config)
	TerminateIfError("failed to create k8s API extension clientset. %+v", err)
	rookClientset, err := rookclient.NewForConfig(config)
	TerminateIfError("failed to create rook clientset. %+v", err)
	return clientset, apiExtClientset, rookClientset
}

// GetKubeConfig supports both in-cluster and out of the cluster with config file (~/.kube/config)
func GetKubeConfig() *rest.Config {

	// Try to read config from in-cluster env
	config, err := rest.InClusterConfig()
	if err == nil {
		return config
	}

	// Get user config file
	config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{}).ClientConfig()
	TerminateIfError("failed to get k8s config. %+v", err)

	return config
}

// NewContext create the k8s clients and returns the context
func NewContext() *clusterd.Context {
	clientset, apiExtClientset, rookClientset := GetClientset()
	return &clusterd.Context{
		Clientset:             clientset,
		APIExtensionClientset: apiExtClientset,
		RookClientset:         rookClientset,
		Executor:              &exec.CommandExecutor{},
		NetworkInfo:           clusterd.NetworkInfo{},
		ConfigDir:             k8sutil.DataDir,
	}
}

// TerminateIfError terminates if err is not nil
func TerminateIfError(format string, err error) {
	if err != nil {
		TerminateFatal(fmt.Errorf(format, err))
	}

}

// TerminateFatal terminates the process with an exit code of 1
// and writes the given reason to stderr and the termination log file.
func TerminateFatal(reason error) {
	fmt.Fprintln(os.Stderr, reason)

	file, err := os.OpenFile(terminationLog, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("failed to write message to termination log: %+v", err))
	} else {
		defer file.Close()
		if _, err = file.WriteString(reason.Error()); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("failed to write message to termination log: %+v", err))
		}
	}

	os.Exit(1)
}
