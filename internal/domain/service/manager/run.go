package manager




// If run returns error in the middle of running task, it means task is not successfully finished
func (m *Manager) Run() error {

	// Here are procedure.

	// 1. Connect to Worker Manager as a Websocket.
	// Next, wait for any task to be allocated.

	// 2. Create stock data receiver as a channel.
	// If event type is real, connect to kafka
	// Otherwise if it is past, connect to mongodb

	// 3. Create model with it's initializer
	// this is made up of sub procedure.
	// 3-1. Get model as a file accessing to MiniO
	// 3-2. Compile c++ model file to a binary.
	// 3-3. Run a binary and create input channel and output channel

	// 4. Put the data receiver channel to the model input

	// 5. Put the model output to the kafka message broker

	// 6. Wait til the task is finished and free all resources

	return nil
}