package gamematcher_test

import (
	"fmt"

	"github.com/Neightly/juzhongmishi/pkg/gamematcher"
)

func ExampleGame_MatchCloseNext() {
	game := gamematcher.New(64)     // 七脉会武是64人紧邻连续匹配的比赛
	l := game.MatchCloseNext(1, 33) // 1号张小凡和33号曾书书只能在决赛（第6轮）相遇
	fmt.Println(l)
	// Output: 6
}

func ExampleGame_MatchHeadTail() {
	game := gamematcher.New(128)  // 网球大满贯是128人首尾对称匹配的比赛
	l := game.MatchHeadTail(1, 2) // 1号种子和2号种子只能在决赛（第7轮）相遇
	fmt.Println(l)
	// Output: 7
}
