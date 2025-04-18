package scoreboard

import (
	"sync"
)

type MutexScoreboardManager struct {
	l          sync.RWMutex
	scoreboard map[string]int
}

func NewMutexScoreboardManager() *MutexScoreboardManager {
	return &MutexScoreboardManager{
		scoreboard: map[string]int{},
	}
}

func (score *MutexScoreboardManager) Update(name string, val int) {
	score.l.Lock()
	defer score.l.Unlock()
	score.scoreboard[name] = val
}

func (rscore *MutexScoreboardManager) Read(name string) (int, bool) {
	rscore.l.RLock()
	defer rscore.l.RUnlock()
	val, ok := rscore.scoreboard[name]
	return val, ok
}
