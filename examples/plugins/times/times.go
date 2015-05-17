package times

import (
	"encoding/json"
	"github.com/bluele/datahose/lib/plugins"
	"time"
)

type TimeHose struct {
	plugins.BasePlugin
}

func New() plugins.Plugin {
	return &TimeHose{}
}

func (ho *TimeHose) Name() string {
	return "time"
}

func (ho *TimeHose) Execute() {
	ch := ho.DataChannel()
	per := 10

	for {
		pack := make([][]byte, per)
		for i := 0; i < per; i++ {
			pack[i] = getCurrentTime()
		}
		ch <- pack
		time.Sleep(time.Second)
	}
}

func (ho *TimeHose) IsRunning() bool {
	return len(ho.Pool) > 0
}

func getCurrentTime() []byte {
	v, _ := json.Marshal(map[string]string{
		"time": time.Now().String(),
	})
	return v
}
