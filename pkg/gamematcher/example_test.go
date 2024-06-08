package gamematcher_test

import (
	"fmt"

	"github.com/Neightly/juzhongmishi/pkg/gamematcher"
)

func ExampleGame_MatchCloseNext() {
	game := gamematcher.New(64)             // 64人紧邻连续匹配的比赛
	l, ps, qs := game.MatchCloseNext(1, 33) // 1号张小凡和33号曾书书只能在决赛（第6轮）相遇
	_, _ = ps, qs                           // 备注：非verbose模式的ps和qs无意义。
	fmt.Println(l)
	// Output: 6
}

func ExampleGame_MatchHeadTail() {
	game := gamematcher.New(128)          // 128人首尾对称匹配的比赛
	l, ps, qs := game.MatchHeadTail(1, 2) // 1号种子和2号种子只能在决赛（第7轮）相遇
	_, _ = ps, qs                         // 备注：非verbose模式的ps和qs无意义。
	fmt.Println(l)
	// Output: 7
}

func ExampleGame_Verbose_closeNext() {
	// 备注：非verbose模式的ps和qs无意义。
	game := gamematcher.New(64).Verbose()    // 64人紧邻连续匹配的比赛，verbose模式
	l, ps, qs := game.MatchCloseNext(13, 58) // 13号和58号相遇在第6轮
	fmt.Println(l)
	// 第一轮战胜14, 第二轮战胜15, 第三轮战胜9， 第四轮战胜1, 第五轮战胜17， 第六轮遇上58
	fmt.Println(ps)
	// 第一轮战胜57, 第二轮战胜59, 第三轮战胜61， 第四轮战胜49, 第五轮战胜33， 第六轮遇上13
	fmt.Println(qs)
	// Output:
	// 6
	// [13:->14 13:->15 13:<-9 9:<-1 1:->17 1:33(58)]
	// [58:<-57 57:->59 57:->61 57:<-49 49:<-33 33:1(13)]
}

func ExampleGame_Verbose_headTail() {
	// 备注：非verbose模式的ps和qs无意义。
	game := gamematcher.New(128).Verbose()  // 128人首尾对称匹配的比赛，verbose模式
	l, ps, qs := game.MatchHeadTail(45, 74) // 45号和74号相遇在第7轮
	fmt.Println(l)
	// 第一轮战胜84, 第二轮战胜20, 第三轮战胜13， 第四轮战胜4, 第五轮战胜5， 第六轮战胜1, 第七轮遇上74
	fmt.Println(ps)
	// 第一轮战胜55, 第二轮战胜10, 第三轮战胜23， 第四轮战胜7, 第五轮战胜2， 第六轮战胜3, 第七轮遇上45
	fmt.Println(qs)
	// Output:
	// 7
	// [45:->84 45:<-20 20:<-13 13:<-4 4:->5 4:<-1 1:2(74)]
	// [74:<-55 55:<-10 10:->23 10:<-7 7:<-2 2:->3 2:1(45)]
}
