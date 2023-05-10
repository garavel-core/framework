// Form:
// - https://www.php.net/manual/en/book.strings.php
// - https://github.com/laravel/framework/blob/9.x/src/Illuminate/Support/Str.php

package str

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gertd/go-pluralize"
	"github.com/google/uuid"
	"github.com/mozillazg/go-unidecode"

	"github.com/garavel-core/framework/internal/maps"
	"github.com/garavel-core/framework/support/slices"
)

type stringable interface {
	string | []string
}

// The callback that should be used to generate UUIDs.
var uuidFactory func() uuid.UUID

// The callback that should be used to generate random strings.
var randomStringFactory func(int) string

// The pluralization client instance.
var pluralizer *pluralize.Client

func init() {
	pluralizer = pluralize.NewClient()
}

// wrapStringable receives any value and returns a string slice
// containing the value as its only element, if the value is a string,
// or returns the value itself, if it is a string slice.
//
// This function is intended to be used with any type that is convertible
// to a string or string slice (i.e., is "stringable").
//
// Parameters:
//   - value: stringable value that can be converted to a string or string slice.
//
// Returns:
//   - A string slice containing the value, or the value itself, if it is
//     already a string slice.
func wrapStringable(value any) []string {
	if str, ok := value.(string); ok {
		return []string{str}
	}

	return value.([]string)
}

func StrIReplace(search, replace, subject string) string {
	m := len(search)
	n := len(subject)
	if m > n {
		return subject
	}
	skipTable := make([]int, 256)
	for i := range skipTable {
		skipTable[i] = m
	}
	for i := 0; i < m-1; i++ {
		skipTable[search[i]] = m - i - 1
	}
	skipTable[search[m-1]] = m

	result := strings.Builder{}
	i := 0

	for i <= n-m {
		j := m - 1
		for j >= 0 && (subject[i+j] == search[j] || subject[i+j] == search[j]+32 || subject[i+j] == search[j]-32) {
			j--
		}
		if j < 0 {
			result.WriteString(replace)
			i += m
		} else {
			result.WriteByte(subject[i])
			i += skipTable[subject[i+m-1]]
		}
	}
	for ; i < n; i++ {
		result.WriteByte(subject[i])
	}
	return result.String()
}

func Ieplace(search string, replace string, subject string, count ...*int) string {
	s := len(search)
	n := len(subject)

	if s == 0 || n == 0 || s > n {
		return subject
	}

	if s == n && strings.EqualFold(subject, search) {
		return replace
	}

	var result strings.Builder
	// result := strings.Builder{}
	result.Grow(n + -1*(len(replace)-s))

	// for i := 0; i < n; i += s {
	// 	str := subject[i : i+s]

	// 	strings.Index(str, search)
	// 	// if str == search {

	// 	// }

	// 	// if str == search { //strings.EqualFold(str, search) {
	// 	// 	result.WriteString(replace)

	// 	// 	// 记录被替换的次数
	// 	// 	if count != nil && count[0] != nil {
	// 	// 		*count[0]++
	// 	// 	}
	// 	// } else if i+s > n-s {
	// 	// 	//result.WriteString(subject[i:])
	// 	// 	break
	// 	// } else {
	// 	// result.WriteString(str)
	// 	//}
	// }

	return result.String()
}

// Case-insensitive version of Replace()
func IReplace[S stringable, R stringable, T stringable](search S, replace R, subject T, count ...*int) T {
	var n string

	subjects := wrapStringable(subject)

	searches := wrapStringable(search)

	replaces := wrapStringable(replace)

	l := len(replaces)

	for i, s := range subjects {
		if s == "" {
			continue
		}

		for j, o := range searches {
			if o == "" {
				continue
			}

			if j >= l {
				n = ""
			} else {
				n = replaces[j]
			}

			// 统计替换的次数
			if count != nil && count[0] != nil {
				*count[0] += strings.Count(s, o)
			}

			s = strings.ReplaceAll(s, o, n)
		}

		subjects[i] = s
	}

	if slice, ok := any(subjects).(T); ok {
		return slice
	}

	return any(subjects[0]).(T)
}

// Replace all occurrences of the search string with the replacement string.
func Replace[S stringable, R stringable, T stringable](search S, replace R, subject T, count ...*int) T {
	var n string

	subjects := wrapStringable(subject)

	searches := wrapStringable(search)

	replaces := wrapStringable(replace)

	l := len(replaces)

	for i, s := range subjects {
		if s == "" {
			continue
		}

		for j, o := range searches {
			if o == "" {
				continue
			}

			if j >= l {
				n = ""
			} else {
				n = replaces[j]
			}

			// 统计替换的次数
			if count != nil && count[0] != nil {
				*count[0] += strings.Count(s, o)
			}

			s = strings.ReplaceAll(s, o, n)
		}

		subjects[i] = s
	}

	if slice, ok := any(subjects).(T); ok {
		return slice
	}

	return any(subjects[0]).(T)
}

// Translate characters or replace substrings
func Strtr[T string, F string | map[string]string](subject string, from F, to ...T) string {
	return ""
}

// Returns the number of substring occurrences.
func SubstrCount() int {
	return 0
}

// Replace text within a portion of a string.
func SubstrReplace() string {
	return ""
}

// Returns the portion of the string specified by the start and length parameters.
func Substr(str string) string {
	return ""
}

// Make a string's first character lowercase.
func Lcfirst(str string) string {
	if len(str) == 0 {
		return str
	}

	return strings.ToLower(str[:1]) + str[1:]
}

// Make a string's first character uppercase.
func Ucfirst(str string) string {
	if len(str) == 0 {
		return str
	}

	return strings.ToUpper(str[:1]) + str[1:]
}

// Return the remainder of a string after the first occurrence of a given value.
func After(subject string, search string) string {
	return ""
}

// Return the remainder of a string after the last occurrence of a given value.
func AfterLast(subject string, search string) string {
	return ""
}

// Transliterate a UTF-8 value to ASCII.
func Ascii(value string, language ...string) string {
	// TODO: To support language
	return unidecode.Unidecode(value)
}

// Transliterate a string to its closest ASCII representation.
func Transliterate(str string, unknown string, strict ...bool) string {
	return ""
}

// Get the portion of a string before the first occurrence of a given value.
func Before(subject string, search string) string {
	return ""
}

// Get the portion of a string before the last occurrence of a given value.
func BeforeLast(subject string, search string) string {
	return ""
}

// Get the portion of a string between two given values.
func Between(subject string, from string, to string) string {
	return ""
}

// Get the smallest possible portion of a string between two given values.
func BetweenFirst(subject string, from string, to string) string {
	return ""
}

// Convert a value to camel case.
func Camel(value string) string {
	return Lcfirst(Studly(value))
}

// Determine if a given string contains a given substring.
func Contains(haystack string, needles string, ignoreCase ...bool) bool {
	if ignoreCase != nil && ignoreCase[0] {
		haystack = strings.ToLower(haystack)
	}

	return strings.Contains(haystack, needles)
}

// Determine if a given string contains all array values.
func ContainsAll(haystack string, needles []string, ignoreCase ...bool) bool {
	if ignoreCase != nil && ignoreCase[0] {
		haystack = strings.ToLower(haystack)
	}

	for _, needle := range needles {
		if !strings.Contains(haystack, needle) {
			return false
		}
	}

	return true
}

// Determine if a given string ends with a given substring.
func EndsWith[T stringable](haystack string, needles T) bool {

	return false
}

// Extracts an excerpt from text that matches the first instance of a phrase.
func Excerpt(text string, phrase string, options ...any) string {
	return ""
}

// Cap a string with a single instance of a given value.
func Finish(value string, cap string) string {
	return ""
}

// Wrap the string with the given strings.
func Wrap(value string, before string, after ...string) string {
	return ""
}

// Determine if a given string matches a given pattern.
func Is[T stringable](patterns T, value string) bool {
	for _, pattern := range wrapStringable(patterns) {
		// If the given value is an exact match we can of course return true right
		// from the beginning. Otherwise, we will translate asterisks and do an
		// actual pattern match against the two strings to see if they match.
		if pattern == value {
			return true
		}

		pattern = regexp.QuoteMeta(pattern)

		// Asterisks are translated into zero-or-more regular expression wildcards
		// to make it convenient to check if the strings starts with the given
		// pattern such as "library/*", making any string check convenient.
		pattern = strings.ReplaceAll(pattern, `\*`, ".*")

		if regexp.MustCompile(`^` + pattern + `\z`).MatchString(value) {
			return true
		}
	}

	return false
}

// Determine if a given string is 7 bit ASCII.
func IsAscii(value string) bool {
	return false
}

// Determine if a given string is valid JSON.
func IsJson(value any) bool {
	if str, ok := value.(string); !ok {
		return false
	} else {
		var decoded map[string]any
		return json.Unmarshal([]byte(str), &decoded) == nil
	}
}

// Determine if a given string is a valid UUID.
func IsUuid(value any) bool {
	if str, ok := value.(string); !ok {
		return false
	} else {
		return regexp.MustCompile(`(?i)^[\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12}$`).MatchString(str)
	}
}

// Determine if a given string is a valid ULID.
func IsUlid(value any) bool {
	if str, ok := value.(string); !ok {
		return false
	} else {
		// TODO 继续完善
		return len(str) != 26
	}
}

// Convert a string to kebab case.
func Kebab(value string) string {
	return Snake(value, "-")
}

// Return the length of the given string.
func Length(value string, encoding ...string) int {
	// TODO: 暂时只支持 utf-8
	if encoding == nil || strings.ToUpper(encoding[0]) == "UTF-8" {
		return utf8.RuneCountInString(value)
	}

	return len(value)
}

// Limit the number of characters in a string.
func Limit(value string, limit int, end ...string) string {
	return ""
}

// Convert the given string to lower-case.
func Lower(value string) string {
	return strings.ToLower(value)
}

// Limit the number of words in a string.
func Words(value string, words int, end ...string) string {
	return ""
}

// Converts GitHub flavored Markdown into HTML.
func Markdown(str string, options ...any) string {
	return str
}

// Converts inline Markdown into HTML.
func InlineMarkdown(str string, options ...any) string {
	return str
}

// Masks a portion of a string with a repeated character.
func Mask(str string, character string, index int, length int, encoding ...string) string {
	return str
}

// Get the string matching the given pattern.
func Match(pattern string, subject string) string {
	matches := regexp.MustCompile(pattern).FindStringSubmatch(subject)

	if matches == nil {
		return ""
	}

	if len(matches) > 1 && matches[1] != "" {
		return matches[1]
	}

	return matches[0]
}

// Get the string matching the given pattern.
func MatchAll(pattern string, subject string) any {
	return nil
}

// Pad both sides of a string with another.
func PadBoth(value string, length int, pad ...string) string {
	return ""
}

// Pad the left side of a string with another.
func PadLeft(value string, length int, pad ...string) string {
	return ""
}

// Pad the right side of a string with another.
func PadRight(value string, length int, pad ...string) string {
	pads := slices.Get(pad, 0, " ")

	if len(pads) == 0 {
		return value
	}

	count := length - utf8.RuneCountInString(value)

	if count < 1 {
		return value
	}

	return value + string([]rune(strings.Repeat(pads, count))[:count])
}

// Parse a Class[@]method style callback into class and method.
func ParseCallback(callback string, defaultValue ...string) any {
	return nil
}

// Get the plural form of an English word.
//
// Parameters:
//   - value     string       the word to pluralize
//   - count     int          how many of the word exist
//   - inclusive bool        whether to prefix with the number (e.g. 3 ducks)
func Plural(value string, count ...any) string {
	if len(value) == 0 {
		return value
	}

	// TODO count support Array and Collections
	if count == nil || (count[0] != 1 && count[0] != -1) {
		value = pluralizer.Plural(value)
	}

	if slices.Get(count, 1, false) {
		return fmt.Sprintf("%v %s", count[0], value)
	}

	return value
}

// Pluralize the last word of an English, studly caps case string.
func PluralStudly(value string, count ...int) string {
	// TODO count support Array and Collections
	if len(value) != 0 && (count == nil || (count[0] != 1 && count[0] != -1)) {
		index := strings.LastIndexFunc(value, func(r rune) bool {
			return r >= 'A' && r <= 'Z'
		})

		if index != -1 {
			value = value[:index] + Plural(value[index:])
		}
	}

	return value
}

// Generate a more truly "random" alpha-numeric string.
func Random(length ...int) string {
	// 设置默认参数
	if length == nil {
		length = append(length, 16)
	}

	if randomStringFactory != nil {
		return randomStringFactory(length[0])
	}

	return (func(length int) string {
		var result strings.Builder

		result.Grow(length)

		for size := length; size > 0; size = length - result.Len() {
			bytesSize := int(math.Ceil(float64(size)/3) * 3)

			bytes := make([]byte, bytesSize)

			// RawStdEncoding 编码器不会追加填充字符 =
			if _, err := rand.Read(bytes); err != nil {
				continue
			}

			str := Replace([]string{"/", "+"}, "", base64.RawStdEncoding.EncodeToString(bytes))

			// 确保索引安全
			if len(str) > size {
				str = str[:size]
			}

			result.WriteString(str)
		}
		return result.String()
	})(length[0])
}

// Set the callable that will be used to generate random strings.
func CreateRandomStringsUsing(factory ...func(int) string) {
	if factory == nil {
		randomStringFactory = nil
	} else {
		randomStringFactory = factory[0]
	}
}

// Set the sequence that will be used to generate random strings.
func CreateRandomStringsUsingSequence(sequence map[int]string, whenMissing ...func(int) string) {
	next := 0

	if whenMissing == nil {
		whenMissing = append(whenMissing, func(length int) string {
			factoryCache := randomStringFactory

			randomStringFactory = nil

			randomString := Random(length)

			randomStringFactory = factoryCache

			next++

			return randomString
		})
	}

	CreateRandomStringsUsing(func(length int) string {
		if str, exists := sequence[next]; exists {
			next++
			return str
		}

		return whenMissing[0](length)
	})
}

// Indicate that random strings should be created normally and not using a custom factory.
func CreateRandomStringsNormally() {
	randomStringFactory = nil
}

// Repeat the given string.
func Repeat(str string, times int) string {
	return strings.Repeat(str, times)
}

// Replace a given value in the string sequentially with an array.
func ReplaceArray(search string, replace []string, subject string) string {
	// 避免内存分配
	if len(search) == 0 || len(subject) == 0 || len(replace) == 0 {
		return subject
	}

	segments := strings.Split(subject, search)

	if len(segments) == 1 {
		return subject
	}

	var result strings.Builder

	result.Grow(len(subject))

	for i, segment := range segments {
		if i > 0 {
			result.WriteString(slices.Get(replace, i-1, search))
		}

		result.WriteString(segment)
	}

	return result.String()
}

// Replace the first occurrence of a given value in the string.
func ReplaceFirst(search string, replace string, subject string) string {
	// 避免内存分配
	if len(search) == 0 || len(subject) == 0 {
		return subject
	}

	if position := strings.Index(subject, search); position != -1 {
		return subject[:position] + replace + subject[position+len(search):]
	}

	return subject
}

// Replace the last occurrence of a given value in the string.
func ReplaceLast(search string, replace string, subject string) string {
	// 避免内存分配
	if len(search) == 0 || len(subject) == 0 {
		return subject
	}

	if position := strings.LastIndex(subject, search); position != -1 {
		return subject[:position] + replace + subject[position+len(search):]
	}

	return subject
}

// Remove any occurrence of the given string in the subject.
func Remove[T stringable](search T, subject string, caseSensitive ...bool) string {
	// 默认是大小写敏感的
	if caseSensitive != nil && !caseSensitive[0] {
		return IReplace(search, "", subject)
	}

	return Replace(search, "", subject)
}

// Reverse the given string.
func Reverse(value string) string {
	// 避免内存分配
	if len(value) == 0 {
		return value
	}

	runes := []rune(value)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// Begin a string with a single instance of a given value.
func Start(value string, prefix string) string {
	if len(value) == 0 || len(prefix) == 0 {
		return value
	}

	// 比正则表达式快
	for l := len(prefix); strings.HasPrefix(value, prefix); {
		value = value[l:]
	}

	return prefix + value
}

// Convert the given string to upper-case.
func Upper(value string) string {
	// Go 语言中的字符串标准库能够正确处理多字节字符
	return strings.ToUpper(value)
}

// Convert the given string to title case.
func Title(value string) string {
	if len(value) == 0 {
		return value
	}

	var result strings.Builder
	var prev rune

	result.Grow(len(value))

	for i, r := range value {
		// unicode.IsLetter 判断中文也返回真，这和 php 中的 mb_convert_case 函数不一致
		// 所以在这里使用 ToUpper 方法，能转换的都属于字母
		if i == 0 || (!unicode.IsUpper(prev) && unicode.ToUpper(prev) == prev) {
			result.WriteRune(unicode.ToUpper(r))
		} else {
			result.WriteRune(unicode.ToLower(r))
		}

		prev = r
	}

	return result.String()
}

// Convert the given string to title case for each word.
func Headline(value string) string {
	// 避免分配内存
	if len(value) == 0 {
		return value
	}

	if !strings.Contains(value, " ") {
		value = strings.Join(Ucsplit(value), " ")
	}

	var result strings.Builder
	var separable bool

	result.Grow(len(value))

	for i, r := range value {
		// 跳过多个分隔符，最终只保留一个
		if r == ' ' || r == '-' || r == '_' {
			separable = true
			continue
		} else if i == 0 || separable {
			if result.Len() != 0 {
				result.WriteByte(' ')
			}

			result.WriteRune(unicode.ToUpper(r))
		} else {
			result.WriteRune(unicode.ToLower(r))
		}

		separable = false
	}

	return result.String()
}

// Get the singular form of an English word.
func Singular(value string) string {
	return pluralizer.Singular(value)
}

// Generate a URL friendly "slug" from a given string.
func Slug(title string, args ...any) string {
	// 避免内存分配
	if len(title) == 0 {
		return title
	}

	if language, ok := slices.Get(args, 1, any("en")).(string); ok && len(language) != 0 {
		title = Ascii(title, language)
	}

	// 获取参数并配置默认值
	separator := slices.Get(args, 0, "-")
	dictionary := slices.Get(args, 2, map[string]string{"@": "at"})

	// Convert all dashes/underscores into separator
	flip := "-"

	if separator == "-" {
		flip = "_"
	}

	title = strings.ReplaceAll(title, flip, separator)

	// Replace dictionary words
	for key, value := range dictionary {
		title = strings.ReplaceAll(title, key, separator+value+separator)
	}

	var result strings.Builder

	result.Grow(len(title))

	n, start := 0, 0

	// Remove all characters that are not the separator, letters, numbers, or whitespace
	// Replace all separator characters and whitespace by a single separator
	for i, r := range title {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {

			if n == 0 {
				if start == -1 {
					result.WriteString(separator)
				}

				start = i
			}

			n++

		} else if n > 0 {
			result.WriteString(title[start:i])

			// 如果是空字符或者是分隔符则替换成自定义的分隔符
			if IsSpace(r) || strings.ContainsRune(separator, r) && i <= len(title) {
				start = -1
			}

			n = 0
		}
	}

	if n > 0 {
		result.WriteString(title[start:])
	}

	return result.String()
}

// IsSpace reports whether the rune is a space character as defined
// by Unicode's White Space property; in the Latin-1 space
// this is
//
//	'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
//
// Faster by around 40 ns than unicode.IsSpace()
func IsSpace(r rune) bool {
	return r == ' ' || (r >= '\t' && r <= '\r') || r == '\u00A0' || r == '\u0085'
}

// Convert a string to snake case.
func Snake(value string, delimiter ...string) string {
	if len(value) == 0 {
		return value
	}

	// 设置默认参数
	if delimiter == nil {
		delimiter = append(delimiter, "_")
	}

	var result strings.Builder
	var separable bool

	result.Grow(len(value))

	for _, r := range value {
		if IsSpace(r) {
			separable = true
			continue
		} else if r >= 'A' && r <= 'Z' {
			separable = true
		}

		if separable && result.Len() != 0 {
			result.WriteString(delimiter[0])
		}

		separable = false

		result.WriteRune(unicode.ToLower(r))
	}

	return result.String()
}

// Remove all "extra" blank space from the given string.
func Squish(value string) string {
	if len(value) == 0 {
		return value
	}

	var result strings.Builder

	result.Grow(len(value))

	// n-已缓冲的字符串数 start-字符串截取的开始索引
	n, start := 0, 0

	for i, r := range value {
		if !IsSpace(r) && r != '\u3164' && r != '\uFEFF' {
			if n == 0 {
				if result.Len() != 0 {
					result.WriteByte(' ')
				}

				start = i
			}

			n++

		} else if n > 0 {
			result.WriteString(value[start:i])

			n = 0
		}
	}

	if n > 0 {
		result.WriteString(value[start:])
	}

	return result.String()
}

// Determine if a given string starts with a given substring.
func StartsWith(haystack string, needles ...string) bool {
	for _, needle := range needles {
		if needle != "" && strings.HasPrefix(haystack, needle) {
			return true
		}
	}

	return false
}

// Convert a value to studly caps case.
func Studly(str string) string {
	// 避免分配内存
	if len(str) == 0 {
		return str
	}

	var result strings.Builder
	// 快 10ns
	result.Grow(len(str))

	// n-缓冲区的字符数量，start-截取字符串开始的索引
	n, start := 0, 0

	for i, r := range str {
		if r != '-' && r != '_' && r != ' ' {
			if n == 0 {
				result.WriteRune(unicode.ToUpper(r))

				start = i + utf8.RuneLen(r)
			}

			n++
		} else if n > 0 {
			result.WriteString(str[start:i])

			n = 0
		}
	}

	if n > 0 {
		result.WriteString(str[start:])
	}

	return result.String()
}

// Swap multiple keywords in a string with other keywords.
func Swap(pairs map[string]string, subject string) string {
	return Strtr(subject, pairs)
}

// Split a string into pieces by uppercase characters.
func Ucsplit(str string) []string {
	if len(str) == 0 {
		return nil
	}

	words := make([]string, 0)
	last := 0

	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			words = append(words, str[last:i])
			last = i
		}
	}

	return append(words, str[last:])
}

// Get the number of words a string contains.
// Support two calling methods as follows:
// PHP：WordCount(str, format, characters)
// Laravel：WordCoun(str, characters)
func WordCount(str string, args ...any) any {
	format := 0
	characters := "a-z'-"

	// 处理参数让其支持 php 和 laravel 两种调用方式
	if args != nil {
		if f, ok := args[0].(int); ok {
			format = f
		} else if args[0] != nil {
			characters += args[0].(string)
		}

		if len(args) > 1 {
			characters += args[1].(string)
		}
	}

	var buf strings.Builder

	words := make(map[int]string, 0)

	n, start := 0, 0

	// 比正则快 210 倍，以 str.WordCount("мама", "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ") 为例
	// 正则 66787 ns/op
	// 字符串算法 317.7 ns/op
	for i, r := range str {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || strings.IndexRune(characters, r) != -1 {
			if n == 0 {
				start = i
			}

			n++
		} else if n > 0 {
			buf.WriteString(str[start:i])

			words[start] = buf.String()

			buf.Reset()

			n = 0
		}
	}

	if n > 0 {
		buf.WriteString(str[start:])

		words[start] = buf.String()
	}

	if format == 2 {
		return words
	}

	if format == 1 {
		return maps.Values(words)
	}

	return len(words)
}

// Generate a UUID (version 4).
func Uuid() uuid.UUID {
	if uuidFactory != nil {
		return uuidFactory()
	}

	return uuid.New()
}

// Generate a time-ordered UUID (version 4).
func OrderedUuid() uuid.UUID {
	// TODO 时间排序
	return uuid.New()
}

// Set the callable that will be used to generate UUIDs.
func CreateUuidsUsing(factory ...func() uuid.UUID) {
	if factory == nil {
		uuidFactory = nil
	} else {
		uuidFactory = factory[0]
	}
}

// Set the sequence that will be used to generate UUIDs.
func CreateUuidsUsingSequence(sequence map[int]uuid.UUID, whenMissign ...func() uuid.UUID) {
	next := 0

	if whenMissign == nil {
		whenMissign = append(whenMissign, func() uuid.UUID {
			factoryCache := uuidFactory

			uuidFactory = nil

			u := Uuid()

			uuidFactory = factoryCache

			next++

			return u
		})
	}

	CreateUuidsUsing(func() uuid.UUID {
		if u, exists := sequence[next]; exists {
			next++
			return u
		}

		return whenMissign[0]()
	})
}

// Always return the same UUID when generating new UUIDs.
func FreezeUuids(callback ...func(uuid.UUID)) uuid.UUID {
	u := Uuid()

	CreateUuidsUsing(func() uuid.UUID {
		return u
	})

	if callback != nil && callback[0] != nil {
		defer CreateUuidsNormally()

		callback[0](u)
	}

	return u
}

// Indicate that UUIDs should be created normally and not using a custom factory.
func CreateUuidsNormally() {
	uuidFactory = nil
}

// Generate a ULID.
func Ulid() {

}
