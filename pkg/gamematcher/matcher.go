package gamematcher

import (
	"fmt"
	"math/bits"
	"runtime"
)

// Game 表示若干选手参加的淘汰赛
type Game struct {
	players uint64 // 数量必须是2的幂次
}

// New 创建n个选手参加的淘汰赛，n必须是2的幂次。
func New(n uint64) Game {
	assertf(bits.OnesCount64(n) == 1, "n(%d) must be power of 2", n)
	return Game{players: n}
}

// assertPlayersInRange 保证p和q不超出范围：譬如128人的比赛不允许出现150号选手。
func (g Game) assertPlayersInRange(p, q uint64) {
	assertf(1 <= p && p <= g.players, "p(%d) out of range: [1, %d]", p, g.players)
	assertf(1 <= q && q <= g.players, "q(%d) out of range: [1, %d]", q, g.players)
}

// MatchCloseNext 给出紧邻连续匹配规则下选手p和q相遇的轮次。
// p和q必须具有合法编号：譬如128人的比赛不允许出现150号选手。
// 自己和自己原则上不允许相遇，在此按第0轮处理。
func (g Game) MatchCloseNext(p, q uint64) uint8 {
	g.assertPlayersInRange(p, q)
	if p == q {
		return 0
	}
	return g.optimizedCloseNext(p-1, q-1) // p and q => 0-based
}

// optimizedCloseNext 计算p^q可以用几个bit表示，个数越大轮次越大。
// 譬如0b100/0b101/0b110/0b111都可以用3个bit表示。3代表第3轮。
func (g Game) optimizedCloseNext(p, q uint64) uint8 {
	return uint8(bits.Len64(p ^ q))
}

// MatchHeadTail 给出首尾对称匹配规则下选手p和q相遇的轮次。
// p和q必须具有合法编号：譬如128人的比赛不允许出现150号选手。
// 自己和自己原则上不允许相遇，在此按第0轮处理。
func (g Game) MatchHeadTail(p, q uint64) uint8 {
	g.assertPlayersInRange(p, q)
	if p == q {
		return 0
	}
	return g.optimizedHeadTail(p-1, q-1) // p and q => 0-based
}

// optimizedHeadTail 计算p^q尾部连续1的个数或者连续0的个数，个数越大轮次越小。
// 譬如0b110111尾部有3个1，将其转换为0b111000便于计算。3代表倒数第3轮。
func (g Game) optimizedHeadTail(p, q uint64) uint8 {
	xor := p ^ q
	if xor&0b1 != 0 {
		xor++
	}
	return uint8(bits.Len64(g.players) - bits.TrailingZeros64(xor))
}

func assertf(p bool, format string, args ...any) {
	if !p {
		msg := fmt.Sprintf(format, args...)
		// Include information about the assertion location.
		// Due to panic recovery, this location is otherwise buried in the middle of the panicking stack.
		if _, file, line, ok := runtime.Caller(1); ok {
			msg = fmt.Sprintf("%s:%d: %s", file, line, msg)
		}
		panic(msg)
	}
}
