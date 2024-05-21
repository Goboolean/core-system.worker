package job

import "github.com/Goboolean/core-system.worker/internal/model"

type UserParams map[string]string

func (p UserParams) IsKeyNullOrEmpty(key string) bool {
	if val, ok := p[key]; !ok || val == "" {
		return true
	} else {
		return false
	}
}

type DataChan chan model.Packet
