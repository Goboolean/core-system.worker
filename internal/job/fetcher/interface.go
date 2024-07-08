package fetcher

// fetcher는 pipeline의 trade data fetch 단계를 수행할 수 있는 Job을 구현합니다.
// 모든 fetch job은 Fetcher 인터페이스를 구현해야 합니다.

import (
	"github.com/Goboolean/core-system.worker/internal/job"
)

// Fetcher represents a job fetcher that retrieves trade data.
type Fetcher interface {
	job.Common

	// Output returns the data channel for the fetched trade data.
	Output() job.DataChan

	// Stop stops the fetcher and releases any allocated resources.
	NotifyStop()
}
