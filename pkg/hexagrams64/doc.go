// Package hexagrams64 implements base64 encoding with 64 HEXAGRAM.
package hexagrams64

// table对应的字符串，因为[64]rune不能为const，在此记录以做标记。
const _ = "" +
	"\u4DC0\u4DC1\u4DC2\u4DC3\u4DC4\u4DC5\u4DC6\u4DC7\u4DC8\u4DC9\u4DCA\u4DCB\u4DCC\u4DCD\u4DCE\u4DCF" +
	"\u4DD0\u4DD1\u4DD2\u4DD3\u4DD4\u4DD5\u4DD6\u4DD7\u4DD8\u4DD9\u4DDA\u4DDB\u4DDC\u4DDD\u4DDE\u4DDF" +
	"\u4DE0\u4DE1\u4DE2\u4DE3\u4DE4\u4DE5\u4DE6\u4DE7\u4DE8\u4DE9\u4DEA\u4DEB\u4DEC\u4DED\u4DEE\u4DEF" +
	"\u4DF0\u4DF1\u4DF2\u4DF3\u4DF4\u4DF5\u4DF6\u4DF7\u4DF8\u4DF9\u4DFA\u4DFB\u4DFC\u4DFD\u4DFE\u4DFF"

// 易经64卦符号表
var table = [64]rune{
	'\u4DC0', '\u4DC1', '\u4DC2', '\u4DC3', '\u4DC4', '\u4DC5', '\u4DC6', '\u4DC7', '\u4DC8', '\u4DC9', '\u4DCA', '\u4DCB', '\u4DCC', '\u4DCD', '\u4DCE', '\u4DCF',
	'\u4DD0', '\u4DD1', '\u4DD2', '\u4DD3', '\u4DD4', '\u4DD5', '\u4DD6', '\u4DD7', '\u4DD8', '\u4DD9', '\u4DDA', '\u4DDB', '\u4DDC', '\u4DDD', '\u4DDE', '\u4DDF',
	'\u4DE0', '\u4DE1', '\u4DE2', '\u4DE3', '\u4DE4', '\u4DE5', '\u4DE6', '\u4DE7', '\u4DE8', '\u4DE9', '\u4DEA', '\u4DEB', '\u4DEC', '\u4DED', '\u4DEE', '\u4DEF',
	'\u4DF0', '\u4DF1', '\u4DF2', '\u4DF3', '\u4DF4', '\u4DF5', '\u4DF6', '\u4DF7', '\u4DF8', '\u4DF9', '\u4DFA', '\u4DFB', '\u4DFC', '\u4DFD', '\u4DFE', '\u4DFF',
}

// StdPadding 对应base64.StdPadding
const StdPadding rune = '\u3007'

// NoPadding 对应base64.NoPadding
const NoPadding rune = -1

// StdEncoding 对应base64.StdEncoding
var StdEncoding = Encoding{padChar: StdPadding}

// RawStdEncoding 对应base64.RawStdEncoding
var RawStdEncoding = Encoding{padChar: NoPadding}
