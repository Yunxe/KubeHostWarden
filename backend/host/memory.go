package host

type Memory struct {
	Id string `json:"id"`
	HostId string `json:"host_id"`
	Free float64 `json:"free"`
	Available float64 `json:"available"`
}

func (Memory) TableName() string {
	return "memory"
}

