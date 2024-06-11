package gamematcher_test

import (
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
			name:   "1",
			args:   args{1},
			suffix: "", // 1 is ok
		},
		{
			name:   "8",
			args:   args{8},
			suffix: "", // 8 is ok
		},
		{
			name:   "26",
			args:   args{26},
			suffix: "n(26) must be power of 2", // panic("... n(26) must be power of 2")
		},
		{
			name:   "32",
			args:   args{32},
			suffix: "", // 32 is ok
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer func() {
					got, _ := recover().(string)
					if !strings.HasSuffix(got, tt.suffix) {
						t.Errorf("gamematcher.New() got recover = %q, want suffix %q", got, tt.suffix)
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
		name string
		g    gamematcher.Game
		args args
		want uint8
	}{
		{
			name: "p=q",
			g:    gamematcher.New(32),
			args: args{7, 7},
			want: 0,
		},
		{
			name: "七脉会武",
			g:    gamematcher.New(64),
			args: args{1, 33},
			want: 6,
		},
		{
			name: "紧邻不匹配",
			g:    gamematcher.New(128),
			args: args{32, 33},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.MatchCloseNext(tt.args.p, tt.args.q); got != tt.want {
				t.Errorf("Game.MatchCloseNext() = %v, want %v", got, tt.want)
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
		name string
		g    gamematcher.Game
		args args
		want uint8
	}{
		{
			name: "p=q",
			g:    gamematcher.New(8),
			args: args{5, 5},
			want: 0,
		},
		{
			name: "NBA季后赛",
			g:    gamematcher.New(8),
			args: args{2, 6},
			want: 2,
		},
		{
			name: "七脉会武",
			g:    gamematcher.New(64),
			args: args{26, 47},
			want: 4,
		},
		{
			name: "网球大满贯",
			g:    gamematcher.New(128),
			args: args{125, 128},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.MatchHeadTail(tt.args.p, tt.args.q); got != tt.want {
				t.Errorf("Game.MatchHeadTail() = %v, want %v", got, tt.want)
			}
		})
	}
}
