package str_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/garavel-core/framework/support/arr"
	"github.com/garavel-core/framework/support/slices"
	"github.com/garavel-core/framework/support/str"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestReplace(t *testing.T) {
	t.Run("replace single string", func(t *testing.T) {
		input := "hello world, hello"
		expected := "hi world, hi"
		actual := str.Replace("hello", "hi", input)
		assert.Equal(t, expected, actual)
	})

	t.Run("replace single string multiple times", func(t *testing.T) {
		input := "apple apple banana apple pear"
		expected := "orange orange banana orange pear"
		actual := str.Replace("apple", "orange", input)
		assert.Equal(t, expected, actual)
	})

	t.Run("replace multiple strings with a single replacement", func(t *testing.T) {
		input := "one two three four"
		expected := "1   "
		actual := str.Replace([]string{"one", "two", "three", "four"}, "1", input)
		assert.Equal(t, expected, actual)
	})

	t.Run("replace multiple strings with multiple replacements", func(t *testing.T) {
		input := "blue sky red flower yellow sun"
		expected := "green sky pink flower white sun"
		actual := str.Replace([]string{"blue", "red", "yellow"}, []string{"green", "pink", "white"}, input)
		assert.Equal(t, expected, actual)
	})

	t.Run("replace empty string with non-empty string", func(t *testing.T) {
		input := "The quick brown fox jumps over the lazy dog."
		expected := "The quick brown fox jumps over the lazy dog."
		actual := str.Replace("", " ", input)
		assert.Equal(t, expected, actual)
	})

	t.Run("replace non-empty string with empty string", func(t *testing.T) {
		input := "Hello, world!"
		expected := "!"
		actual := str.Replace("Hello, world", "", input)
		assert.Equal(t, expected, actual)
	})

	t.Run("replace with count, no matches", func(t *testing.T) {
		input := "hello world"
		expected := "hello world"
		count := 0
		actual := str.Replace("foo", "bar", input, &count)
		assert.Equal(t, expected, actual)
		assert.Equal(t, 0, count)
	})

	t.Run("replace with count, single match", func(t *testing.T) {
		input := "apple apple banana"
		expected := "mango mango banana"
		count := 0
		actual := str.Replace("apple", "mango", input, &count)
		assert.Equal(t, expected, actual)
		assert.Equal(t, 2, count)
	})

	t.Run("replace with count, multiple matches", func(t *testing.T) {
		input := "blue sky blue sea blue ocean"
		expected := "red sky red sea red ocean"
		count := 0
		actual := str.Replace([]string{"blue"}, []string{"red"}, input, &count)
		assert.Equal(t, expected, actual)
		assert.Equal(t, 3, count)
	})

	t.Run("replace with count, unequal length", func(t *testing.T) {
		input := "I love ice cream!"
		expected := "I love mango!"
		count := 0
		actual := str.Replace("ice cream", "mango", input, &count)
		assert.Equal(t, expected, actual)
		assert.Equal(t, 1, count)
	})
}

func TestUcfirst(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		input := ""
		expected := ""
		actual := str.Ucfirst(input)
		assert.Equal(t, expected, actual)
	})

	t.Run("single letter", func(t *testing.T) {
		input := "a"
		expected := "A"
		actual := str.Ucfirst(input)
		assert.Equal(t, expected, actual)
	})

	t.Run("lowercase string", func(t *testing.T) {
		input := "hello world"
		expected := "Hello world"
		actual := str.Ucfirst(input)
		assert.Equal(t, expected, actual)
	})

	t.Run("uppercase string", func(t *testing.T) {
		input := "HELLO WORLD"
		expected := "HELLO WORLD"
		actual := str.Ucfirst(input)
		assert.Equal(t, expected, actual)
	})

	t.Run("mixed case string", func(t *testing.T) {
		input := "hElLo WorLD"
		expected := "HElLo WorLD"
		actual := str.Ucfirst(input)
		assert.Equal(t, expected, actual)
	})

	t.Run("string with whitespace", func(t *testing.T) {
		input := " hello world"
		expected := " hello world"
		actual := str.Ucfirst(input)
		assert.Equal(t, expected, actual)
	})
}

func TestLcfirst(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		input := ""
		expected := ""
		actual := str.Lcfirst(input)
		assert.Equal(t, expected, actual)
	})

	t.Run("single letter", func(t *testing.T) {
		input := "A"
		expected := "a"
		actual := str.Lcfirst(input)
		assert.Equal(t, expected, actual)
	})

	t.Run("lowercase string", func(t *testing.T) {
		input := "hello world"
		expected := "hello world"
		actual := str.Lcfirst(input)
		assert.Equal(t, expected, actual)
	})

	t.Run("uppercase string", func(t *testing.T) {
		input := "HELLO WORLD"
		expected := "hELLO WORLD"
		actual := str.Lcfirst(input)
		assert.Equal(t, expected, actual)
	})

	t.Run("mixed case string", func(t *testing.T) {
		input := "HeLLo WorLD"
		expected := "heLLo WorLD"
		actual := str.Lcfirst(input)
		assert.Equal(t, expected, actual)
	})

	t.Run("string with whitespace", func(t *testing.T) {
		input := "  Hello world"
		expected := "  Hello world"
		actual := str.Lcfirst(input)
		assert.Equal(t, expected, actual)
	})
}

func BenchmarkPadRight(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.PadRight("❤MultiByte☆", 16)
	}
}

func TestPadRight(t *testing.T) {
	assert.Equal(t, "Alien-=-=-", str.PadRight("Alien", 10, "-="))
	assert.Equal(t, "Alien     ", str.PadRight("Alien", 10))
	assert.Equal(t, "❤MultiByte☆     ", str.PadRight("❤MultiByte☆", 16))
	assert.Equal(t, "❤MultiByte☆❤☆❤☆❤", str.PadRight("❤MultiByte☆", 16, "❤☆"))
}

// 12840 ns/op，相较于其他函数较慢
func BenchmarkPlural(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Plural("apple", 3, true)
	}
}

func TestPlural(t *testing.T) {
	assert.Equal(t, "apples", str.Plural("apple"))
	assert.Equal(t, "apple", str.Plural("apple", 1))
	assert.Equal(t, "apples", str.Plural("apple", 0))
	assert.Equal(t, "apples", str.Plural("apple", 2))
	// Support inclusive parameter
	assert.Equal(t, "1 apple", str.Plural("apple", 1, true))
	assert.Equal(t, "3 apples", str.Plural("apple", 3, true))
	// Ensure consistency with laravel
	assert.Equal(t, "apple", str.Plural("apple", -1))
}

// 14799 ns/op，相较于其他函数较慢
func BenchmarkPluralStudly(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.PluralStudly("MultipleWordsInOneString")
	}
}

func TestPluralStudly(t *testing.T) {
	assert.Equal(t, "RealHumans", str.PluralStudly("RealHuman"))
	assert.Equal(t, "Models", str.PluralStudly("Model"))
	assert.Equal(t, "VortexFields", str.PluralStudly("VortexField"))
	assert.Equal(t, "MultipleWordsInOneStrings", str.PluralStudly("MultipleWordsInOneString"))

	// With count
	assert.Equal(t, "RealHuman", str.PluralStudly("RealHuman", 1))
	assert.Equal(t, "RealHumans", str.PluralStudly("RealHuman", 2))
	assert.Equal(t, "RealHuman", str.PluralStudly("RealHuman", -1))
	assert.Equal(t, "RealHumans", str.PluralStudly("RealHuman", -2))
}

func BenchmarkRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Random()
	}
}

func TestRandom(t *testing.T) {
	assert.Equal(t, 16, str.Length(str.Random()))
	randomInteger := int(rand.Int63() % 100)
	assert.Equal(t, randomInteger, str.Length(str.Random(randomInteger)))

	t.Run("Whether the number of generated characters is equally distributed", func(t *testing.T) {
		results := make(map[string]int, 620000)

		for i := 0; i < len(results); i++ {
			random := str.Random(1)
			results[random]++
		}

		// each character should occur 100.000 times with a variance of 5%.
		for _, result := range results {
			assert.InDelta(t, 10000, result, 500)
		}
	})
}

func BenchmarkCreateRandomStringsUsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.CreateRandomStringsUsing(func(length int) string {
			return "length:" + strconv.Itoa(length)
		})
	}
}

func TestCreateRandomStringsUsing(t *testing.T) {
	str.CreateRandomStringsUsing(func(length int) string {
		return "length:" + strconv.Itoa(length)
	})

	assert.Equal(t, "length:7", str.Random(7))
	assert.Equal(t, "length:7", str.Random(7))

	str.CreateRandomStringsNormally()

	assert.NotEqual(t, "length:7", str.Random())
}

func BenchmarkCreateRandomStringsUsingSequence(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.CreateRandomStringsUsingSequence(map[int]string{
			0: "x",
		})
	}
}

func TestCreateRandomStringsUsingSequence(t *testing.T) {
	t.Run("Can specify a sequence of random strings to utilise", func(t *testing.T) {
		str.CreateRandomStringsUsingSequence(map[int]string{
			0: "x",
			// 1 => just generate a random one here...
			2: "y",
			3: "z",
			// ... => continue to generate random strings...
		})

		assert.Equal(t, "x", str.Random())
		assert.Equal(t, 16, str.Length(str.Random()))
		assert.Equal(t, "y", str.Random())
		assert.Equal(t, "z", str.Random())
		assert.Equal(t, 16, str.Length(str.Random()))
		assert.Equal(t, 16, str.Length(str.Random()))

		str.CreateRandomStringsNormally()
	})

	t.Run("Can specify a fallback for a random string sequence", func(t *testing.T) {
		str.CreateRandomStringsUsingSequence(map[int]string{
			0: str.Random(),
			1: str.Random(),
		}, func(_ int) string {
			panic("Out of random strings.")
		})

		str.Random()
		str.Random()

		assert.Panics(t, func() {
			str.Random()
		})

		str.CreateRandomStringsNormally()
	})
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Repeat("a", 5)
	}
}

func TestRepeat(t *testing.T) {
	assert.Equal(t, "aaaaa", str.Repeat("a", 5))
	assert.Equal(t, "", str.Repeat("", 5))

	// Ensure unnecessary memory allocation
	allocs := testing.AllocsPerRun(1, func() {
		str.Repeat("", 5)
		str.Repeat("", 0)
	})

	assert.Equal(t, 0, int(allocs))
}

func BenchmarkReplaceArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.ReplaceArray("?", []string{"foo", "bar", "baz"}, "?/?/?")
	}
}

func TestReplaceArray(t *testing.T) {
	assert.Equal(t, "?/?/?", str.ReplaceArray("?", nil, "?/?/?"))
	assert.Equal(t, "?/?/?", str.ReplaceArray("?", []string{}, "?/?/?"))
	assert.Equal(t, "foo/bar/baz", str.ReplaceArray("?", []string{"foo", "bar", "baz"}, "?/?/?"))
	assert.Equal(t, "foo/bar/baz/?", str.ReplaceArray("?", []string{"foo", "bar", "baz"}, "?/?/?/?"))
	assert.Equal(t, "foo/bar", str.ReplaceArray("?", []string{"foo", "bar", "baz"}, "?/?"))
	assert.Equal(t, "?/?/?", str.ReplaceArray("x", []string{"foo", "bar", "baz"}, "?/?/?"))
	// Ensure recursive replacements are avoided
	assert.Equal(t, "foo?/bar/baz", str.ReplaceArray("?", []string{"foo?", "bar", "baz"}, "?/?/?"))

	// Ensure unnecessary memory allocation
	allocs := testing.AllocsPerRun(1, func() {
		str.ReplaceArray("", []string{"foo", "bar", "baz"}, "?/?/?")
		str.ReplaceArray("?", nil, "?/?/?")
		str.ReplaceArray("?", []string{}, "?/?/?")
		str.ReplaceArray("?", []string{"foo", "bar", "baz"}, "")
	})

	assert.Equal(t, 0, int(allocs))
}

func BenchmarkReplaceFirst(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.ReplaceFirst("bar?", "qux?", "foo/bar? foo/bar?")
	}
}

func TestReplaceFirst(t *testing.T) {
	assert.Equal(t, "fooqux foobar", str.ReplaceFirst("bar", "qux", "foobar foobar"))
	assert.Equal(t, "foo/qux? foo/bar?", str.ReplaceFirst("bar?", "qux?", "foo/bar? foo/bar?"))
	assert.Equal(t, "foo foobar", str.ReplaceFirst("bar", "", "foobar foobar"))
	assert.Equal(t, "foobar foobar", str.ReplaceFirst("xxx", "yyy", "foobar foobar"))
	assert.Equal(t, "foobar foobar", str.ReplaceFirst("", "yyy", "foobar foobar"))

	// Test for multibyte string support
	assert.Equal(t, "Jxxxnköping Malmö", str.ReplaceFirst("ö", "xxx", "Jönköping Malmö"))
	assert.Equal(t, "Jönköping Malmö", str.ReplaceFirst("", "yyy", "Jönköping Malmö"))
}

func BenchmarkReplaceLast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.ReplaceLast("ö", "xxx", "Malmö Jönköping")
	}
}

func TestReplaceLast(t *testing.T) {
	assert.Equal(t, "foobar fooqux", str.ReplaceLast("bar", "qux", "foobar foobar"))
	assert.Equal(t, "foo/bar? foo/qux?", str.ReplaceLast("bar?", "qux?", "foo/bar? foo/bar?"))
	assert.Equal(t, "foobar foo", str.ReplaceLast("bar", "", "foobar foobar"))
	assert.Equal(t, "foobar foobar", str.ReplaceLast("xxx", "yyy", "foobar foobar"))
	assert.Equal(t, "foobar foobar", str.ReplaceLast("", "yyy", "foobar foobar"))
	// Test for multibyte string support
	assert.Equal(t, "Malmö Jönkxxxping", str.ReplaceLast("ö", "xxx", "Malmö Jönköping"))
	assert.Equal(t, "Malmö Jönköping", str.ReplaceLast("", "yyy", "Malmö Jönköping"))
}

func BenchmarkRemove(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Remove([]string{"f", "b"}, "Foobar", false)
	}
}

func TestRemove(t *testing.T) {
	assert.Equal(t, "Fbar", str.Remove("o", "Foobar"))
	assert.Equal(t, "Foo", str.Remove("bar", "Foobar"))
	assert.Equal(t, "oobar", str.Remove("F", "Foobar"))
	assert.Equal(t, "Foobar", str.Remove("f", "Foobar"))
	assert.Equal(t, "oobar", str.Remove("f", "Foobar", false))

	assert.Equal(t, "Fbr", str.Remove([]string{"o", "a"}, "Foobar"))
	assert.Equal(t, "Fooar", str.Remove([]string{"f", "b"}, "Foobar"))
	assert.Equal(t, "ooar", str.Remove([]string{"f", "b"}, "Foobar", false))
	assert.Equal(t, "Foobar", str.Remove([]string{"f", "|"}, "Foo|bar"))
}

func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Reverse("☆etyBitluM❤")
	}
}

func TestReverse(t *testing.T) {
	assert.Equal(t, "FooBar", str.Reverse("raBooF"))
	assert.Equal(t, "Teniszütő", str.Reverse("őtüzsineT"))
	assert.Equal(t, "❤MultiByte☆", str.Reverse("☆etyBitluM❤"))
}

func BenchmarkStart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Start("//test/string", "/")
	}
}

func TestStart(t *testing.T) {
	assert.Equal(t, "/test/string", str.Start("test/string", "/"))
	assert.Equal(t, "/test/string", str.Start("/test/string", "/"))
	assert.Equal(t, "/test/string", str.Start("//test/string", "/"))
}

func BenchmarkUpper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Upper("żółtałódka")
	}
}

func TestUpper(t *testing.T) {
	assert.Equal(t, "FOO BAR BAZ", str.Upper("foo bar baz"))
	assert.Equal(t, "FOO BAR BAZ", str.Upper("foO bAr BaZ"))
	// Test for multibyte string support
	assert.Equal(t, "ŻÓŁTAŁÓDKA", str.Upper("żółtałódka"))
}

func BenchmarkTitle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Title("jefferson costella")
	}
}

func TestTitle(t *testing.T) {
	assert.Equal(t, "Jefferson Costella", str.Title("jefferson costella"))
	assert.Equal(t, "Jefferson Costella", str.Title("jefFErson coSTella"))
	assert.Equal(t, "Jefferson-Costella", str.Title("jefFErson-coSTella"))
	assert.Equal(t, "Jefferson中Costella", str.Title("jefFErson中coSTella"))
	assert.Equal(t, "Jefferson Ö Costella", str.Title("jefFErson ö coSTella"))
	assert.Equal(t, "Jeffersonöcostella", str.Title("jefFErsonöcoSTella"))
}

func BenchmarkHeadline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Headline("jefferson costella")
	}
}

func TestHeadline(t *testing.T) {
	assert.Equal(t, "Jefferson Costella", str.Headline("jefferson costella"))
	assert.Equal(t, "Jefferson Costella", str.Headline("jefFErson coSTella"))
	assert.Equal(t, "Jefferson Costella Uses Garavel", str.Headline("jefferson_costella uses-_Garavel"))
	assert.Equal(t, "Jefferson Costella Uses Garavel", str.Headline("jefferson_costella uses__Garavel"))

	assert.Equal(t, "Garavel G O Lang Framework", str.Headline("garavel_g_o_lang_framework"))
	assert.Equal(t, "Garavel G O Lang Framework", str.Headline("garavel _g _o _lang _framework"))
	assert.Equal(t, "Garavel Golang Framework", str.Headline("garavel_golang_framework"))
	assert.Equal(t, "Garavel Go Lang Framework", str.Headline("garavel-goLang-framework"))
	assert.Equal(t, "Garavel Golang Framework", str.Headline("   garavel  -_-  golang   -_-   framework   "))

	assert.Equal(t, "Foo Bar", str.Headline("fooBar"))
	assert.Equal(t, "Foo Bar", str.Headline("foo_bar"))
	assert.Equal(t, "Foo Bar Baz", str.Headline("foo-barBaz"))
	assert.Equal(t, "Foo Bar Baz", str.Headline("foo-bar_baz"))

	assert.Equal(t, "Öffentliche Überraschungen", str.Headline("öffentliche-überraschungen"))
	assert.Equal(t, "Öffentliche Überraschungen", str.Headline("-_öffentliche_überraschungen_-"))
	assert.Equal(t, "Öffentliche Überraschungen", str.Headline("-öffentliche überraschungen"))

	assert.Equal(t, "Sind Öde Und So", str.Headline("sindÖdeUndSo"))

	assert.Equal(t, "Orwell 1984", str.Headline("orwell 1984"))
	assert.Equal(t, "Orwell 1984", str.Headline("orwell   1984"))
	assert.Equal(t, "Orwell 1984", str.Headline("-orwell-1984 -"))
	assert.Equal(t, "Orwell 1984", str.Headline(" orwell_- 1984 "))
}

// 17104 ns/op，相较于其他函数性能偏低
func BenchmarkSingular(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Singular("apples")
	}
}

func TestSingular(t *testing.T) {
	assert.Equal(t, "apple", str.Singular("apples"))
}

func BenchmarkSlug(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Slug("hello world")
	}
}

func TestSlug(t *testing.T) {
	assert.Equal(t, "hello-world", str.Slug("hello   world"))
	assert.Equal(t, "hello-world", str.Slug("hello-world"))
	assert.Equal(t, "hello-world", str.Slug("hello_world"))
	assert.Equal(t, "hello_world", str.Slug("hello_world", "_"))
	assert.Equal(t, "user-at-host", str.Slug("user@host"))
	assert.Equal(t, "سلام-دنیا", str.Slug("سلام دنیا", "-", nil))
	assert.Equal(t, "sometext", str.Slug("some text", ""))
	assert.Equal(t, "", str.Slug("", ""))
	assert.Equal(t, "", str.Slug(""))
	// TODO ascii 三方库翻译结果有误暂时先不测试
	// assert.Equal(t, "bsm-allah", str.Slug("بسم الله", "-", "en", map[string]string{"allh": "allah"}))
	assert.Equal(t, "500-dollar-bill", str.Slug("500$ bill", "-", "en", map[string]string{"$": "dollar"}))
	assert.Equal(t, "500-dollar-bill", str.Slug("500--$----bill", "-", "en", map[string]string{"$": "dollar"}))
	assert.Equal(t, "500--dollar--bill", str.Slug("500-$-bill   ", "--", "en", map[string]string{"$": "dollar"}))
	assert.Equal(t, "500-dollar-bill", str.Slug("500$--bill", "-", "en", map[string]string{"$": "dollar"}))
	assert.Equal(t, "500-dollar-bill", str.Slug("500-$--bill", "-", "en", map[string]string{"$": "dollar"}))
	assert.Equal(t, "أحمد-في-المدرسة", str.Slug("أحمد@المدرسة", "-", nil, map[string]string{"@": "في"}))
}

func BenchmarkSnake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Snake("GaravelGOLangFramework")
	}
}

func TestSnake(t *testing.T) {
	assert.Equal(t, "garavel_g_o_lang_framework", str.Snake("GaravelGOLangFramework"))
	assert.Equal(t, "garavel_golang_framework", str.Snake("GaravelGolangFramework"))
	assert.Equal(t, "garavel golang framework", str.Snake("GaravelGolangFramework", " "))
	assert.Equal(t, "garavel_golang_framework", str.Snake("Garavel Golang Framework"))
	assert.Equal(t, "garavel_golang_framework", str.Snake("   Garavel    Golang      Framework   "))
	// ensure cache keys don"t overlap
	assert.Equal(t, "garavel__golang__framework", str.Snake("GaravelGolangFramework", "__"))
	assert.Equal(t, "garavel_golang_framework_", str.Snake("GaravelGolangFramework_", "_"))
	assert.Equal(t, "garavel_golang_framework", str.Snake("garavel golang Framework"))
	assert.Equal(t, "garavel_golang_frame_work", str.Snake("garavel golang FrameWork"))
	// prevent breaking changes
	assert.Equal(t, "foo-bar", str.Snake("foo-bar"))
	assert.Equal(t, "foo-_bar", str.Snake("Foo-Bar"))
	assert.Equal(t, "foo__bar", str.Snake("Foo_Bar"))
	assert.Equal(t, "żółtałódka", str.Snake("ŻółtaŁódka"))
}

func BenchmarkSquish(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Squish(" garavel   go  framework ")
	}
}

func TestSquish(t *testing.T) {
	assert.Equal(t, "garavel go framework", str.Squish(" garavel   go  framework "))
	assert.Equal(t, "garavel go framework", str.Squish("garavel\t\tgo\n\nframework"))
	assert.Equal(t, "garavel go framework", str.Squish(`
        garavel
        go
        framework
    `))
	assert.Equal(t, "garavel go framework", str.Squish("   garavel   go   framework   "))
	assert.Equal(t, "123", str.Squish("   123    "))
	assert.Equal(t, "だ", str.Squish("だ"))
	assert.Equal(t, "ム", str.Squish("ム"))
	assert.Equal(t, "だ", str.Squish("   だ    "))
	assert.Equal(t, "ム", str.Squish("   ム    "))
	assert.Equal(t, "garavel go framework", str.Squish("garavelㅤㅤㅤgoㅤframework"))
}

func BenchmarkStartsWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.StartsWith("jason", "jas")
	}
}

func TestStartsWith(t *testing.T) {
	assert.True(t, str.StartsWith("jason", "jas"))
	assert.True(t, str.StartsWith("jason", "jason"))
	assert.True(t, str.StartsWith("jason", "day", "jas"))
	assert.False(t, str.StartsWith("jason", "day"))
	assert.False(t, str.StartsWith("jason", ""))
	assert.False(t, str.StartsWith("0123", ""))
	assert.True(t, str.StartsWith("0123", "0"))
	assert.False(t, str.StartsWith("jason", "J"))
	assert.False(t, str.StartsWith("jason", ""))
	assert.False(t, str.StartsWith("", ""))
	assert.False(t, str.StartsWith("7", " 7"))
	assert.True(t, str.StartsWith("7a", "7"))
	// Test for multibyte string support
	assert.True(t, str.StartsWith("Jönköping", "Jö"))
	assert.True(t, str.StartsWith("Malmö", "Malmö"))
	assert.False(t, str.StartsWith("Jönköping", "Jonko"))
	assert.False(t, str.StartsWith("Malmö", "Malmo"))
	assert.True(t, str.StartsWith("你好", "你"))
	assert.False(t, str.StartsWith("你好", "好"))
	assert.False(t, str.StartsWith("你好", "a"))
}

func BenchmarkStudly(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Studly("garavel_g_o_lang_framework")
	}
}

func TestStudly(t *testing.T) {
	assert.Equal(t, "GaravelGOLangFramework", str.Studly("garavel_g_o_lang_framework"))
	assert.Equal(t, "GaravelGolangFramework", str.Studly("garavel_golang_framework"))
	assert.Equal(t, "GaravelGoLangFramework", str.Studly("garavel-goLang-framework"))
	assert.Equal(t, "GaravelGolangFramework", str.Studly("  garavel  -_-  golang   -_-   framework   "))

	assert.Equal(t, "FooBar", str.Studly("fooBar"))
	assert.Equal(t, "FooBar", str.Studly("foo_bar"))
	// assert.Equal(t, "FooBar", str.Studly("foo_bar")); // test cache
	assert.Equal(t, "FooBarBaz", str.Studly("foo-barBaz"))
	assert.Equal(t, "FooBarBaz", str.Studly("foo-bar_baz"))

	assert.Equal(t, "ÖffentlicheÜberraschungen", str.Studly("öffentliche-überraschungen"))
}

func BenchmarkSwap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Swap(map[string]string{"GO": "GO 1.18", "awesome": "fantastic"}, "GO is awesome")
	}
}

func TestSwap(t *testing.T) {
	assert.Equal(t, "GO 1.18 fantastic", str.Swap(map[string]string{"GO": "GO 1.18", "awesome": "fantastic"}, "GO is awesome"))
	assert.Equal(t, "foo bar baz", str.Swap(map[string]string{"ⓐⓑ": "baz"}, "foo bar ⓐⓑ"))
}

func BenchmarkUcsplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.Ucsplit("Garavel_g_o_lang_framework")
	}
}

func TestUcsplit(t *testing.T) {
	assert.Equal(t, []string{"Garavel_g_o_lang_framework"}, str.Ucsplit("Garavel_g_o_lang_framework"))
	assert.Equal(t, []string{"Garavel_", "G_o_lang_framework"}, str.Ucsplit("Garavel_G_o_lang_framework"))
	assert.Equal(t, []string{"garavel", "G", "O", "Lang", "Framework"}, str.Ucsplit("garavelGOLangFramework"))
	assert.Equal(t, []string{"Garavel-go", "Lang-framework"}, str.Ucsplit("Garavel-goLang-framework"))

	assert.Equal(t, []string{"Żółta", "Łódka"}, str.Ucsplit("ŻółtaŁódka"))
	assert.Equal(t, []string{"sind", "Öde", "Und", "So"}, str.Ucsplit("sindÖdeUndSo"))
	assert.Equal(t, []string{"Öffentliche", "Überraschungen"}, str.Ucsplit("ÖffentlicheÜberraschungen"))
}

func BenchmarkWordCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.WordCount("мама", "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ")
	}
}

func TestWordCount(t *testing.T) {
	assert.Equal(t, 2, str.WordCount("Hello, world!"))
	assert.Equal(t, 10, str.WordCount("Hi, this is my first contribution to the Garavel framework."))

	assert.Equal(t, 0, str.WordCount("мама"))
	assert.Equal(t, 0, str.WordCount("мама мыла раму"))

	assert.Equal(t, 1, str.WordCount("мама", "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"))
	assert.Equal(t, 3, str.WordCount("мама мыла раму", "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"))

	assert.Equal(t, 1, str.WordCount("МАМА", "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"))
	assert.Equal(t, 3, str.WordCount("МАМА МЫЛА РАМУ", "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"))

	// support format like go str_word_count
	assert.Equal(t, 2, str.WordCount("Hello, world!", 0))
	assert.Equal(t, []string{"Hello", "world"}, str.WordCount("Hello, world!", 1))
	assert.Equal(t, map[int]string{0: "Hello", 7: "world"}, str.WordCount("Hello, world!", 2))
}

func BenchmarkCreateUuidsUsingSequence(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.CreateUuidsUsingSequence(map[int]uuid.UUID{
			0: str.Uuid(),
		})
	}
}

func TestCreateUuidsUsingSequence(t *testing.T) {
	t.Run("Can specify a sequence of uuids to utilise", func(t *testing.T) {
		sequence := map[int]uuid.UUID{
			0: str.Uuid(),
			1: str.Uuid(),
			3: str.Uuid(),
		}

		str.CreateUuidsUsingSequence(sequence)

		retrieved := str.Uuid()
		assert.Equal(t, sequence[0], retrieved)
		assert.Equal(t, sequence[0].String(), retrieved.String())

		retrieved = str.Uuid()
		assert.Equal(t, sequence[1], retrieved)
		assert.Equal(t, sequence[1].String(), retrieved.String())

		retrieved = str.Uuid()
		assert.False(t, arr.In(retrieved, sequence))
		assert.False(t, slices.In(retrieved.String(), []string{sequence[0].String(), sequence[1].String(), sequence[3].String()}))

		retrieved = str.Uuid()
		assert.Equal(t, sequence[3], retrieved)
		assert.Equal(t, sequence[3].String(), retrieved.String())

		retrieved = str.Uuid()
		assert.False(t, arr.In(retrieved, sequence))
		assert.False(t, slices.In(retrieved.String(), []string{sequence[0].String(), sequence[1].String(), sequence[3].String()}))

		str.CreateUuidsNormally()
	})

	t.Run("Can specify a fallback for a sequence", func(t *testing.T) {
		str.CreateUuidsUsingSequence(map[int]uuid.UUID{
			0: str.Uuid(),
			1: str.Uuid(),
		}, func() uuid.UUID {
			panic("Out of Uuids.")
		})

		str.Uuid()
		str.Uuid()

		assert.Panics(t, func() {
			str.Uuid()
		})

		str.CreateUuidsNormally()
	})
}

func BenchmarkFreezeUuids(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str.FreezeUuids(func(u uuid.UUID) {

		})
	}
}

func TestFreezeUuids(t *testing.T) {
	t.Run("Can freeze uuids", func(t *testing.T) {
		assert.NotEqual(t, str.Uuid().String(), str.Uuid().String())
		assert.NotEqual(t, str.Uuid(), str.Uuid())

		u := str.FreezeUuids()

		assert.Equal(t, u, str.Uuid())
		assert.Equal(t, str.Uuid(), str.Uuid())
		assert.Equal(t, u.String(), str.Uuid().String())
		assert.Equal(t, str.Uuid().String(), str.Uuid().String())

		str.CreateUuidsNormally()

		assert.NotEqual(t, str.Uuid(), str.Uuid())
		assert.NotEqual(t, str.Uuid().String(), str.Uuid().String())
	})

	t.Run("Can freeze uuids in a closure", func(t *testing.T) {
		uuids := make([]uuid.UUID, 3)

		u := str.FreezeUuids(func(u uuid.UUID) {
			uuids[0] = u
			uuids[1] = str.Uuid()
			uuids[2] = str.Uuid()
		})

		assert.Equal(t, u, uuids[0])
		assert.Equal(t, u.String(), uuids[0].String())
		assert.Equal(t, uuids[0].String(), uuids[1].String())
		assert.Equal(t, uuids[0], uuids[1])
		assert.Equal(t, uuids[1], uuids[2])
		assert.Equal(t, uuids[1].String(), uuids[2].String())

		assert.NotEqual(t, str.Uuid(), str.Uuid())
		assert.NotEqual(t, str.Uuid().String(), str.Uuid().String())

		str.CreateUuidsNormally()
	})

	t.Run("Creates uuids normally after failure with in freeze method", func(t *testing.T) {
		u := str.Uuid()

		assert.Panics(t, func() {
			str.FreezeUuids(func(u uuid.UUID) {
				str.CreateUuidsUsing(func() uuid.UUID {
					return u
				})

				assert.Equal(t, u, str.Uuid())
				panic("Something failed.")
			})
		})

		assert.NotEqual(t, u, str.Uuid().String())
	})
}
