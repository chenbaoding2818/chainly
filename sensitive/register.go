package sensitive

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/chenbaoding2818/chainly/config"
	"github.com/chenbaoding2818/chainly/lifecycle"
)

var (
	sensitiveComponent     *SensitiveManager
	sensitiveComponentOnce sync.Once
)

type SensitiveManager struct {
}

func (sm *SensitiveManager) Start(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config) {
	// TODO: 加入路径
	bytes, err := os.ReadFile("")
	if err != nil {
		panic(err)
	}
	st = newSensitiveTrie()
	contents := strings.Split(string(bytes), "\n")
	st.addWords(contents)
}

func (sm *SensitiveManager) Priority() int32 {
	return lifecycle.LowPriority
}

func (sm *SensitiveManager) Stop() {
}

func NewSensitiveComponent() *SensitiveManager {
	sensitiveComponentOnce.Do(func() {
		if sensitiveComponent == nil {
			sensitiveComponent = &SensitiveManager{}
		}
	})

	return sensitiveComponent
}
