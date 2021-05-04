/*


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

package main

import (
	"flag"
	"log"

	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	controllers "cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com"
	"cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/clowder_config"
	// +kubebuilder:scaffold:imports
)

var setupLog = ctrl.Log.WithName("setup")

func main() {

	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	if clowder_config.LoadedConfig.DebugOptions.Pprof.CpuFile != "" {
		f, err := os.Create(clowder_config.LoadedConfig.DebugOptions.Pprof.CpuFile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	ctrl.SetLogger(zap.New(zap.UseDevMode(false)))

	if clowder_config.LoadedConfig.DebugOptions.Pprof.Enable {
		go http.ListenAndServe("localhost:8000", nil)
	}
	controllers.Run(metricsAddr, enableLeaderElection, ctrl.GetConfigOrDie(), ctrl.SetupSignalHandler())
}