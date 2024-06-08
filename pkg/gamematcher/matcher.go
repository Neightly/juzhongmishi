package gamematcher

import (
	"fmt"
	"math/bits"
	"runtime"
)

// Game 表示若干选手参加的淘汰赛
type Game struct {
	players uint64 // 数量必须是2的幂次
	verbose bool   // verbose模式：不仅给出两个选手相遇的轮次，也包括晋级过程中遇到的对手
}

// New 创建n个选手参加的淘汰赛，n必须是2的幂次且不能是0和1。
func New(n uint64) Game {
	assertf(n > 1 && (n-1)&n == 0, "n(%d) out of range: {2,4,8,16,32,64,128,...}", n)
	return Game{players: n}
}

// Verbose 设置verbose模式：不仅给出两个选手相遇的轮次，也包括晋级过程中遇到的对手。
func (g Game) Verbose() Game {
	g.verbose = true
	return g
}

// Reset 清除verbose模式：仅给出两个选手相遇的轮次，不计算遇到的对手可简化操作。
func (g Game) Reset() Game {
	g.verbose = false
	return g
}

// assertPlayersInRange 保证p和q不超出范围：譬如128人的比赛不允许出现150号选手。
func (g Game) assertPlayersInRange(p, q uint64) {
	assertf(1 <= p && p <= g.players, "p(%d) out of range: [1, %d]", p, g.players)
	assertf(1 <= q && q <= g.players, "q(%d) out of range: [1, %d]", q, g.players)
}

func (g Game) debug(p, q uint64) {
	fmt.Printf("%0[4]*[1]b(%[1]d) ^ %0[4]*[2]b(%[2]d) = %0[4]*[3]b\n",
		p, q, p^q, bits.TrailingZeros64(g.players))
}

// MatchCloseNext 给出紧邻连续匹配规则下选手p和q相遇的轮次。
// 如果是verbose模式，ps和qs描述了各自晋级过程中遇到的对手。
// 自己和自己原则上不允许相遇，在此按第0轮处理。
func (g Game) MatchCloseNext(p, q uint64) (l uint8, ps, qs []string) {
	g.assertPlayersInRange(p, q)
	if p == q {
		return
	}
	if !g.verbose {
		return g.optimizedCloseNext(p-1, q-1)
	}
	return g.verboseCloseNext(p-1, q-1)
}

func (g Game) optimizedCloseNext(p, q uint64) (l uint8, ps, qs []string) {
	for power := uint64(1); ; power <<= 1 {
		if xor := p ^ q; (xor-1)&xor == 0 { // is p^q one of [1 2 4 8 16 32 64 ...]?
			l = uint8(bits.TrailingZeros64(xor)) + 1 // 1-based
			return
		}
		p &^= power
		q &^= power
	}
}

func (g Game) verboseCloseNext(p, q uint64) (l uint8, ps, qs []string) {
	var p0, q0 = p, q // 记录原始值备用
	var p1, q1 uint64 // 对手值
	for power := uint64(1); ; power <<= 1 {
		l++
		p1, q1 = p^power, q^power
		if p^q == power {
			ps = append(ps, fmt.Sprintf("%d:%d(%d)", p+1, p1+1, q0+1))
			qs = append(qs, fmt.Sprintf("%d:%d(%d)", q+1, q1+1, p0+1))
			return
		}
		if p1 < p { // 下克上
			ps = append(ps, fmt.Sprintf("%d:<-%d", p+1, p1+1))
			p = p1
		} else {
			ps = append(ps, fmt.Sprintf("%d:->%d", p+1, p1+1))
		}
		if q1 < q { // 下克上
			qs = append(qs, fmt.Sprintf("%d:<-%d", q+1, q1+1))
			q = q1
		} else {
			qs = append(qs, fmt.Sprintf("%d:->%d", q+1, q1+1))
		}
	}
}

// MatchHeadTail 给出首尾对称匹配规则下选手p和q相遇的轮次。
// 如果是verbose模式，ps和qs描述了各自晋级过程中遇到的对手。
// 自己和自己原则上不允许相遇，在此按第0轮处理。
func (g Game) MatchHeadTail(p, q uint64) (l uint8, ps, qs []string) {
	g.assertPlayersInRange(p, q)
	if p == q {
		return
	}
	if !g.verbose {
		return g.optimizedHeadTail(p-1, q-1)
	}
	return g.verboseHeadTail(p-1, q-1)
}

func (g Game) optimizedHeadTail(p, q uint64) (l uint8, ps, qs []string) {
	for mask := g.players - 1; ; mask >>= 1 {
		if xor := p ^ q; xor&(xor+1) == 0 { // is p^q one of [..., 127, 63, 31, 15, 7, 3, 1]?
			l = uint8(bits.TrailingZeros64(g.players)-bits.TrailingZeros64(xor+1)) + 1 // 1-based
			return
		}
		p = min(p, mask^p)
		q = min(q, mask^q)
	}
}

func (g Game) verboseHeadTail(p, q uint64) (l uint8, ps, qs []string) {
	var p0, q0 = p, q // 记录原始值备用
	var p1, q1 uint64 // 对手值
	for mask := g.players - 1; ; mask >>= 1 {
		l++
		p1, q1 = mask^p, mask^q
		if p^q == mask {
			ps = append(ps, fmt.Sprintf("%d:%d(%d)", p+1, p1+1, q0+1))
			qs = append(qs, fmt.Sprintf("%d:%d(%d)", q+1, q1+1, p0+1))
			return
		}
		if p1 < p { // 下克上
			ps = append(ps, fmt.Sprintf("%d:<-%d", p+1, p1+1))
			p = p1
		} else {
			ps = append(ps, fmt.Sprintf("%d:->%d", p+1, p1+1))
		}
		if q1 < q { // 下克上
			qs = append(qs, fmt.Sprintf("%d:<-%d", q+1, q1+1))
			q = q1
		} else {
			qs = append(qs, fmt.Sprintf("%d:->%d", q+1, q1+1))
		}
	}
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
