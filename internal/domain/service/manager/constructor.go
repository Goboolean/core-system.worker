package manager

import "context"




type Manager struct {

}


func New(ctx context.Context) *Manager {
	return &Manager{}
}

