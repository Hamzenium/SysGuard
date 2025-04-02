package types

type ResourceUsage struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
}

type SetLimit struct {
	CPUThreshold    float64 `json:"cpu_threshold"`
	MemoryThreshold float64 `json:"memory_threshold"`
	DiskThreshold   float64 `json:"disk_threshold"`
}
