package host

import (
	"context"
	"kubehostwarden/host/common"
	"kubehostwarden/host/darwin"
	"kubehostwarden/host/dispatcher"
	"kubehostwarden/host/linux"
	"kubehostwarden/utils/constant"
	"kubehostwarden/utils/log"
	"sync"
)

func NewHostService(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	switch common.GetOSType() {
	case constant.DARWIN:
		log.Info("Starting host service for Darwin")
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
		log.Info("Starting host service for Linux")
		metricsConfig := linux.NewCollectorConfig()
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
	default:
		log.Error("Unsupported OS type", "osType", common.GetOSType())
	}

	wg.Wait()
}
