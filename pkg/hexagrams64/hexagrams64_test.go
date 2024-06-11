package hexagrams64

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type testpair struct {
	decoded  string
	encoded0 string // base64 encoded
	encoded1 string // hexagrams64 encoded
}

var pairs = []testpair{
	// RFC 3548 examples
	{"\x14\xfb\x9c\x03\xd9\x7e", "FPucA9l+", "䷅䷏䷮䷜䷀䷽䷥䷾"},
	{"\x14\xfb\x9c\x03\xd9", "FPucA9k=", "䷅䷏䷮䷜䷀䷽䷤〇"},
	{"\x14\xfb\x9c\x03", "FPucAw==", "䷅䷏䷮䷜䷀䷰〇〇"},

	// RFC 4648 examples
	{"", "", ""},
	{"f", "Zg==", "䷙䷠〇〇"},
	{"fo", "Zm8=", "䷙䷦䷼〇"},
	{"foo", "Zm9v", "䷙䷦䷽䷯"},
	{"foob", "Zm9vYg==", "䷙䷦䷽䷯䷘䷠〇〇"},
	{"fooba", "Zm9vYmE=", "䷙䷦䷽䷯䷘䷦䷄〇"},
	{"foobar", "Zm9vYmFy", "䷙䷦䷽䷯䷘䷦䷅䷲"},

	// Wikipedia examples
	{"sure.", "c3VyZS4=", "䷜䷷䷕䷲䷙䷒䷸〇"},
	{"sure", "c3VyZQ==", "䷜䷷䷕䷲䷙䷐〇〇"},
	{"sur", "c3Vy", "䷜䷷䷕䷲"},
	{"su", "c3U=", "䷜䷷䷔〇"},
	{"leasure.", "bGVhc3VyZS4=", "䷛䷆䷕䷡䷜䷷䷕䷲䷙䷒䷸〇"},
	{"easure.", "ZWFzdXJlLg==", "䷙䷖䷅䷳䷝䷗䷉䷥䷋䷠〇〇"},
	{"asure.", "YXN1cmUu", "䷘䷗䷍䷵䷜䷦䷔䷮"},
	{"sure.", "c3VyZS4=", "䷜䷷䷕䷲䷙䷒䷸〇"},
}

func TestMain(m *testing.M) {
	// 首先保证pairs条目里的编码有效
	for _, p := range pairs {
		if !compitables(p.encoded0, p.encoded1) {
			fmt.Fprintf(os.Stderr, "%q and %q are not compitable for decoded:%q",
				p.encoded0, p.encoded1, p.decoded)
			os.Exit(1)
		}
	}
	os.Exit(m.Run())
}

func TestEncode(t *testing.T) {
	for _, p := range pairs {
		got := StdEncoding.EncodeToString([]byte(p.decoded))
		if got != p.encoded1 {
			t.Errorf("Encode(%q) = %q, want %q", p.decoded, got, p.encoded1)
		}

		dst := StdEncoding.AppendEncode([]byte("lead"), []byte(p.decoded))
		if string(dst) != "lead"+p.encoded1 {
			t.Errorf(`AppendEncode("lead", %q) = %q, want %q`, p.decoded, string(dst), "lead"+p.encoded1)
		}
	}
}

func TestDecode(t *testing.T) {
	for _, p := range pairs {
		encoded := p.encoded1
		dbuf := make([]byte, StdEncoding.DecodedLen(len(encoded)))
		count, err := StdEncoding.Decode(dbuf, []byte(encoded))
		testEqual(t, "Decode(%q) = error %v, want %v", encoded, err, error(nil))
		testEqual(t, "Decode(%q) = length %v, want %v", encoded, count, len(p.decoded))
		testEqual(t, "Decode(%q) = %q, want %q", encoded, string(dbuf[0:count]), p.decoded)

		dbuf, err = StdEncoding.DecodeString(encoded)
		testEqual(t, "DecodeString(%q) = error %v, want %v", encoded, err, error(nil))
		testEqual(t, "DecodeString(%q) = %q, want %q", encoded, string(dbuf), p.decoded)

		dst, err := StdEncoding.AppendDecode([]byte("lead"), []byte(encoded))
		testEqual(t, "AppendDecode(%q) = error %v, want %v", p.encoded1, err, error(nil))
		testEqual(t, `AppendDecode("lead", %q) = %q, want %q`, p.encoded1, string(dst), "lead"+p.decoded)
	}
}

func testEqual(t *testing.T, msg string, args ...any) bool {
	t.Helper()
	if args[len(args)-2] != args[len(args)-1] {
		t.Errorf(msg, args...)
		return false
	}
	return true
}

func compitables(encoded0, encoded1 string) bool {
	const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	for i, hex := range encoded1 {
		std := encoded0[i/3]
		switch {
		case hex == StdPadding:
			if std != '=' {
				return false
			}
		case '\u4DC0' <= hex && hex <= '\u4DFF':
			if int(hex-'\u4DC0') != strings.IndexByte(encodeStd, std) {
				return false
			}
		default:
			panic("unreachable")
		}

	}
	return true
}
