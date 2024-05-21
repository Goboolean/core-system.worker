package joiner

import "github.com/Goboolean/core-system.worker/internal/model"

type Dummy struct {
	Joinner

	refIn   chan model.Packet
	modelIn chan model.Packet
	out     chan model.Packet
}
