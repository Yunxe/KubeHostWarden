package darwin

import "kubehostwarden/host/common"

func NewCollectorConfig() []common.MetricConfig {
    return []common.MetricConfig{
        {
            Type: "cpu",
            CollectFunc: CollectCPU, 
        },
        {
            Type: "memory",
            CollectFunc: CollectMemory, 
        },
        {
            Type: "disk",
            CollectFunc: CollectDisk,
        },
        {
            Type: "load",
            CollectFunc: CollectLoad,
        },
    }
}