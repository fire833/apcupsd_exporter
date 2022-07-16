// Command apcupsd_exporter provides a Prometheus exporter for apcupsd.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"runtime"

	"github.com/mdlayher/apcupsd"
	apcupsdexporter "github.com/mdlayher/apcupsd_exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"
)

var (
	Version   string = "unknown"
	Commit    string = "unknown"
	BuildTime string = "unknown"
	Go        string = runtime.Version()
	Os        string = runtime.GOOS
	Arch      string = runtime.GOARCH

	versionStr string = fmt.Sprintf("apcupsd_exporter - Version: %s\nBuilt on: %s\nCommit: %s\nGo version: %s\nOS: %s\nArch: %s\n\n",
		Version, BuildTime, Commit, Go, Os, Arch)

	telemetryAddr = flag.String("telemetry.addr", ":9162", "address for apcupsd exporter")
	metricsPath   = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")

	apcupsdAddr    = flag.String("apcupsd.addr", ":3551", "address of apcupsd Network Information Server (NIS)")
	apcupsdNetwork = flag.String("apcupsd.network", "tcp", `network of apcupsd Network Information Server (NIS): typically "tcp", "tcp4", or "tcp6"`)
)

func init() {
	klog.InitFlags(nil)
}

func main() {
	flag.Parse()

	fmt.Print(versionStr)

	if *apcupsdAddr == "" {
		klog.Exit("address of apcupsd Network Information Server (NIS) must be specified with '-apcupsd.addr' flag")
	}

	fn := newClient(*apcupsdNetwork, *apcupsdAddr)

	prometheus.MustRegister(apcupsdexporter.New(fn))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	klog.V(1).Infof("starting apcupsd exporter on %q for server %s://%s",
		*telemetryAddr, *apcupsdNetwork, *apcupsdAddr)

	if err := http.ListenAndServe(*telemetryAddr, nil); err != nil {
		klog.Exit("cannot start apcupsd exporter: %s", err)
	}
}

func newClient(network, addr string) apcupsdexporter.ClientFunc {
	return func(ctx context.Context) (*apcupsd.Client, error) {
		return apcupsd.DialContext(ctx, network, addr)
	}
}
