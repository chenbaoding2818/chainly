package sensitive

import (
	"os"
	"strings"

	"github.com/chenbaoding2818/chainly/config"
	"github.com/chenbaoding2818/chainly/lifecycle"
)

func (sm *SensitiveManager) Start(cfg *config.Config) {
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
