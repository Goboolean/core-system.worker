package model

// Packet은 Job과 Job 사이에서 데이터를 전달할 때 사용하는 구조체입니다.
// Job과 Job 사이에서 교환될 데이터와 그 메타 정보를 포합합니다.
type Packet struct {
	// Sequence는 데이터가 생성된 순서를 나타내는 변수입니다.
	// Job과 Job사이를 이동할 때 sequence는 변하지 않습니다.
	Sequence int64

	// Data는 원하는 데이터를 타입에 관계없이 아무거나 보관할 수 있습니다.
	// 이때 Data에 구조체를 싣고 싶은 경우 포인터 타입으로 실어야 합니다.
	Data any
}
