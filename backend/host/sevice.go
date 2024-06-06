package host

import (
	"context"
	"kubehostwarden/host/common"
	"kubehostwarden/host/darwin"
	"kubehostwarden/host/dispatcher"
	"kubehostwarden/host/linux"
	"kubehostwarden/utils/constant"
	"kubehostwarden/utils/log"
)

func NewHostService(ctx context.Context) {
	switch common.GetOSType() {
	case constant.DARWIN:
		log.Info("Starting host service for Darwin")
		metricsConfig := darwin.NewCollectorConfig()
		for _, metricConfig := range metricsConfig {
			collector := common.NewCollector(metricConfig.Type)
			go func(mc common.MetricConfig, col *common.Collector) {
				col.Start(ctx, mc.CollectFunc)
			}(metricConfig, &collector)

			go func(col *common.Collector) {
				dispatcher.Dispatch(ctx, col)
			}(&collector)
		}
	case constant.LINUX:
		log.Info("Starting host service for Linux")
		metricsConfig := linux.NewCollectorConfig()
		for _, metricConfig := range metricsConfig {
			collector := common.NewCollector(metricConfig.Type)
			go func(mc common.MetricConfig, col *common.Collector) {
				col.Start(ctx, mc.CollectFunc)
			}(metricConfig, &collector)

			go func(col *common.Collector) {
				dispatcher.Dispatch(ctx, col)
			}(&collector)
		}
	default:
		log.Error("Unsupported OS type", "osType", common.GetOSType())
	}
}
