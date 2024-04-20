package host

import (
	"context"
	"fmt"
	"kubehostwarden/host/common"
	"kubehostwarden/host/darwin"
	"kubehostwarden/host/dispatcher"
	"kubehostwarden/utils/constant"
	"sync"
)

func NewHostService(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	switch common.GetOSType() {
	case constant.DARWIN:
		metricsConfig := darwin.NewCollectorConfig()
		for _, metricConfig := range metricsConfig {
			collector := common.NewCollector(metricConfig.Type)
			wg.Add(1)
			go func(mc common.MetricConfig, col *common.Collector) {
				defer wg.Done()
				col.Start(ctx, mc.CollectFunc)
			}(metricConfig, &collector)

			wg.Add(1)
			go func(col *common.Collector) {
				defer wg.Done()
				dispatcher.Dispatch(ctx, col)
			}(&collector)
		}
	case constant.LINUX:
		fmt.Println("Linux is not supported yet")
		// Assume similar setup for Linux
		// You might have linux.CollectCPUData or other specific functions
	default:
		fmt.Println("Unsupported OS")
		// Handle unsupported OS
	}

	wg.Wait()
}
