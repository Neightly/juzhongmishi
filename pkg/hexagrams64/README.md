# `hexagrams64`

### `base64`

标准版的`base64`采用`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/`或者`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_`作为编码表。而且还允许指定定制的编码表，但是必须是64个不同的`byte`。

### `hexagrams64`

提到64，中国传统文化中刚好有八卦图与之对应。

很可惜，64个卦图都是3个字节的`rune`，无法直接使用标准版的`base64`。

参考`base64`的设计思想，那么使用64个`rune`又会怎么样呢？

#### 编码表不再需要

由于编码表是固定的64个卦图，不允许指定其它。因此`Encoding`无需保留编码表的属性。

#### 编码映射表不再需要

编码映射表用于指定某个`byte`在编码表中的索引，而且巧妙地利用了`[256]uint8`结构。

由于`hexagrams64`采用了`rune`而不是`byte`，那么就需要`[4294967296]int32`，实际上只用到了64个位置，造成巨大的浪费。

而且64个卦图的`Unicode`是连续的，我们只要把编码表设计成顺序（或者倒序），那么只要知道`Unicode`的最大和最小边界就够了。

#### 严格模式不需要支持

严格模式是为了在特定场景下忽略`\r`和`\n`等特殊符号的。

由于历史原因，某些场景下一行只能容纳少量字节（譬如76），那么大段的`base64`编码就必须分成多行，在解码的时候就需要把拆分的内容合并起来。

但是`hexagrams64`编码的`rune`都是3-byte的，如果`\r`或`\n`出现在`rune`和`rune`之间还好理解，如果出现在`rune`的3个`byte`之间，那就变成非法的编码值了。

为了降低解码难度，我们不处理`\r`和`\n`等特定含义的字节，任何`1-byte rune`、`2-byte rune`、`4-byte rune`全部认定为非法编码值。

为此，我们要同时要求如果有`padding`，`padding`也必须是`3-byte rune`。