package gamematcher_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/Neightly/juzhongmishi/pkg/gamematcher"
)

func TestNew(t *testing.T) {
	type args struct {
		n uint64
	}
	tests := []struct {
		name   string
		args   args
		suffix string
	}{
		{
			name:   "0",
			args:   args{0},
			suffix: "n(0) must be power of 2", // panic("... n(0) must be power of 2")
		},
		{
			name: "2**5",
			args: args{32},
		},
		{
			name: "2**7",
			args: args{128},
		},
		{
			name:   "even",
			args:   args{13},
			suffix: "n(13) must be power of 2", // panic("... n(13) must be power of 2")
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer func() {
					got, _ := recover().(string)
					if !strings.HasSuffix(got, tt.suffix) {
						t.Errorf("gamematcher.New() got = %q, want suffix %q", got, tt.suffix)
					}
				}()
				_ = gamematcher.New(tt.args.n)
			}()
		})
	}
}

func TestGame_MatchCloseNext(t *testing.T) {
	type args struct {
		p uint64
		q uint64
	}
	tests := []struct {
		name   string
		g      gamematcher.Game
		args   args
		wantL  uint8
		wantPs []string
		wantQs []string
	}{
		{
			name:  "64",
			g:     gamematcher.New(64),
			args:  args{16, 17},
			wantL: 5,
			// 第一轮战胜15, 第二轮战胜13, 第三轮战胜9， 第四轮战胜1, 第五轮遇上17
			wantPs: []string{"16:<-15", "15:<-13", "13:<-9", "9:<-1", "1:17(17)"},
			// 第一轮战胜18, 第二轮战胜19, 第三轮战胜21, 第四轮战胜25, 第五轮遇上16
			wantQs: []string{"17:->18", "17:->19", "17:->21", "17:->25", "17:1(16)"},
		},
		{
			name:  "七脉会武",
			g:     gamematcher.New(64),
			args:  args{1, 33},
			wantL: 6,
			// 第一轮战胜2, 第二轮战胜3, 第三轮战胜5， 第四轮战胜9, 第五轮战胜17, 第六轮遇上33
			wantPs: []string{"1:->2", "1:->3", "1:->5", "1:->9", "1:->17", "1:33(33)"},
			// 第一轮战胜34, 第二轮战胜35, 第三轮战胜37， 第四轮战胜41, 第五轮战胜49， 第六轮遇上1
			wantQs: []string{"33:->34", "33:->35", "33:->37", "33:->41", "33:->49", "33:1(1)"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotL0, _, _ := tt.g.MatchCloseNext(tt.args.p, tt.args.q)
			gotL1, gotPs, gotQs := tt.g.Verbose().MatchCloseNext(tt.args.p, tt.args.q)
			if gotL0 != tt.wantL {
				t.Errorf("Game.MatchCloseNext() gotL = %v, want %v", gotL0, tt.wantL)
			}
			if gotL1 != tt.wantL {
				t.Errorf("Game.Verbose().MatchCloseNext() gotL = %v, want %v", gotL1, tt.wantL)
			}
			if tt.wantPs != nil && !slices.Equal(gotPs, tt.wantPs) {
				t.Errorf("Game.MatchCloseNext() gotPs = %v, want %v", gotPs, tt.wantPs)
			}
			if tt.wantQs != nil && !slices.Equal(gotQs, tt.wantQs) {
				t.Errorf("Game.MatchCloseNext() gotQs = %v, want %v", gotQs, tt.wantQs)
			}
		})
	}
}

func TestGame_MatchHeadTail(t *testing.T) {
	type args struct {
		p uint64
		q uint64
	}
	tests := []struct {
		name   string
		g      gamematcher.Game
		args   args
		wantL  uint8
		wantPs []string
		wantQs []string
	}{
		{
			name:  "大满贯",
			g:     gamematcher.New(128),
			args:  args{125, 128},
			wantL: 6,
			// 第一轮战胜4, 第二轮战胜61, 第三轮战胜29， 第四轮战胜13, 第五轮战胜5， 第六轮遇上128
			wantPs: []string{"125:<-4", "4:->61", "4:->29", "4:->13", "4:->5", "4:1(128)"},
			// 第一轮战胜1, 第二轮战胜64, 第三轮战胜32， 第四轮战胜16, 第五轮战胜8， 第六轮遇上125
			wantQs: []string{"128:<-1", "1:->64", "1:->32", "1:->16", "1:->8", "1:4(125)"},
		},
		{
			name:  "七脉会武",
			g:     gamematcher.New(64),
			args:  args{26, 47},
			wantL: 4,
			// 第一轮战胜39, 第二轮战胜7, 第三轮战胜10， 第四轮遇上47
			wantPs: []string{"26:->39", "26:<-7", "7:->10", "7:2(47)"},
			// 第一轮战胜18, 第二轮战胜15, 第三轮战胜2， 第四轮遇上26
			wantQs: []string{"47:<-18", "18:<-15", "15:<-2", "2:7(26)"},
		},
		{
			name:  "NBA季后赛",
			g:     gamematcher.New(8),
			args:  args{2, 6},
			wantL: 2,
			// 第一轮战胜7, 第二轮遇上6
			wantPs: []string{"2:->7", "2:3(6)"},
			// 第一轮战胜3, 第二轮遇上2
			wantQs: []string{"6:<-3", "3:2(2)"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotL0, _, _ := tt.g.MatchHeadTail(tt.args.p, tt.args.q)
			gotL1, gotPs, gotQs := tt.g.Verbose().MatchHeadTail(tt.args.p, tt.args.q)
			if gotL0 != tt.wantL {
				t.Errorf("Game.MatchHeadTail() gotL = %v, want %v", gotL0, tt.wantL)
			}
			if gotL1 != tt.wantL {
				t.Errorf("Game.Verbose().MatchHeadTail() gotL = %v, want %v", gotL1, tt.wantL)
			}
			if tt.wantPs != nil && !slices.Equal(gotPs, tt.wantPs) {
				t.Errorf("Game.MatchHeadTail() gotPs = %v, want %v", gotPs, tt.wantPs)
			}
			if tt.wantQs != nil && !slices.Equal(gotQs, tt.wantQs) {
				t.Errorf("Game.MatchHeadTail() gotQs = %v, want %v", gotQs, tt.wantQs)
			}
		})
	}
}

func TestGame(t *testing.T) {
	// 测试在所有的场景中，optimized模式和verbose模式给出的轮次是否一致
	optimized := gamematcher.New(64)
	verbose := optimized.Verbose()
	var l0, l uint8
	var ps, qs []string
	for p := uint64(1); p <= 64; p++ {
		for q := uint64(1); q <= 64; q++ {
			l0, _, _ = optimized.MatchCloseNext(p, q)
			l, ps, qs = verbose.MatchCloseNext(p, q)
			if l0 != l {
				t.Fatalf("MatchCloseNext(%d, %d) got %d, want %d. details: %v and %v", p, q, l0, l, ps, qs)
			}

			l0, _, _ = optimized.MatchHeadTail(p, q)
			l, ps, qs = verbose.MatchHeadTail(p, q)
			if l0 != l {
				t.Fatalf("MatchHeadTail(%d, %d) got %d, want %d. details: %v and %v", p, q, l0, l, ps, qs)
			}
		}
	}
}
