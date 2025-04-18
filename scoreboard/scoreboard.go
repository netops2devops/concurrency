package scoreboard

import "context"

func scoreboardManager(ctx context.Context, in <-chan func(map[string]int)) {
	scoreboard := map[string]int{}
	for {
		select {
		case <-ctx.Done():
			return
		case f := <-in:
			f(scoreboard)
		}
	}
}

type ChannelScoreboardManager chan func(map[string]int)

func NewChannelScoreboardManager(ctx context.Context) ChannelScoreboardManager {
	ch := make(ChannelScoreboardManager)
	go scoreboardManager(ctx, ch)
	return ch
}

func (csm ChannelScoreboardManager) Update(name string, val int) {
	csm <- func(m map[string]int) {
		m[name] = val
	}
}

func (csm ChannelScoreboardManager) Read(name string) (int, bool) {
	type Result struct {
		out int
		ok  bool
	}
	resultCh := make(chan Result)
	csm <- func(m map[string]int) {
		out, ok := m[name]
		resultCh <- Result{out, ok}
	}
	result := <-resultCh
	return result.out, result.ok
}
