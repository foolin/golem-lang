Golem has a set of standard types.  Each type has a set of values,
together with operators and fields specific to those values.

Basic Types:

  * [Null](#null)
  * [Bool](#bool)
  * [Int](#int)
  * [Float](#float)
  * [Str](#str)

Composite Types:

  * [List](#list)
  * [Range](#range)
  * [Tuple](#tuple)
  * [Dict](#dict)
  * [Set](#set)
  * [Struct](#struct)

Other Types:

  * [Func](#func)
  * [Chan](#chan)

## Null

Null represents the absence of a value. The only instance of Null is `null`.

Null has no valid operators, and no fields.

## Bool

Bool represents boolean truth values.  The only instances of
Bool are `true` and `false`.

Valid operators for Bool are:
* The equality operators `==`, `!=`
* The boolean operators `||` `&&`
* The unary negation operator `!`

Bool has no fields.

## Int

Int is the set of all signed 64-bit integers (-9223372036854775808 to 9223372036854775807).

Valid operators for Int are:
* The equality operators `==`, `!=`
* The comparison operators `>`, `>=`, `<`, `<=`, `<=>`
* The arithmetic operators `+`, `-`, `*`, `/`
* The integer arithmetic operators <code>&#124;</code>, `^`, `%`, `&`, `<<`, `>>`
* The unary integer complement operator `~`
* The postfix operators `++`, `--`

When applying an arithmetic operator `+`, `-`, `*`, `/`to an Int, if the other
operand is a Float, then the result will be a Float,
otherwise the result will be an Int.

Int has no fields.

## Float

Float is the set of all IEEE-754 64-bit floating-point numbers.

Valid operators for Float are:
* The equality operators `==`, `!=`
* The comparison operators `>`, `>=`, `<`, `<=`, `<=>`
* The arithmetic operators `+`, `-`, `*`, `/`
* The postfix operators `++`, `--`

Applying an arithmetic operator to a Float always returns a Float.

Float has no fields.

## Str

Str is the set of all valid sequences of UTF-8-encoded "code points", otherwise
known as "runes".  Strs are immutable.

Valid operators for Str are:
* The equality operators `==`, `!=`,
* The comparison operators `>`, `>=`, `<`, `<=`, `<=>`
* The index operator `a[x]`
* The slice operators `a[x:y]`, `a[x:]`, `a[:y]`

The index operator and slice operators always return a Str.

/*doc
Str has the following fields, all of which are [Funcs](#func):

#### `contains`
`contains` reports whether a substring is within a string.

* signature: `contains(substr <Str>) <Bool>`
* example: `'abcdef'.contains('de')`

#### `hasPrefix`
`hasPrefix` tests whether a string begins with a prefix.

* signature: `hasPrefix(prefix <Str>) <Bool>`
* example: `'abcdef'.hasPrefix('ab')`

#### `hasSuffix`
`hasSuffix` tests whether a string ends with a suffix.

* signature: `hasSuffix(suffix <Str>) <Bool>`
* example: `'abcdef'.hasSuffix('ab')`

#### `index`
`index` returns the index of the first instance of a substring in a string.
or -1 if the substring is not present.

* signature: `index(substr <Str>) <Int>`
* example: `'abcab'.index('ab')`

#### `lastIndex`
`lastIndex` returns the index of the last instance of a substring in a string,
or -1 if the substring is not present.

* signature: `lastIndex(substr <Str>) <Int>`
* example: `'abcab'.lastIndex('ab')`

#### `map`
`map` returns a copy of a string with all its characters modified according to
the mapping function.  The mapping function must accept one Str parameter,
and must return a Str.

* signature: `map(mapping <Func>) <Str>`
* example:

```
let s = 'abc(def)[x,y,z]'
let t = s.map(fn(c) {
    return c >= 'a' && c <= 'z' ? c : ''
})
println(t)
```

#### `replace`
`replace` returns a copy of a string with the first n non-overlapping instances
of `old` replaced by `new`. If `old` is empty, it matches at the beginning of a string
and after each UTF-8 sequence, yielding up to k+1 replacements for a k-rune string.
If `n` < 0, there is no limit on the number of replacements.  The parameter `n` is
optional, and defaults to -1.

* signature: `replace(old <Str>, new <Str>, n = -1 <Int>) <Int>`
* example: `'abcab'.replace('a', 'x')`

#### `split`
`split` slices a string into all substrings separated by sep and returns a list of the substrings between those separators.

If the string does not contain sep and sep is not empty, `split` returns a list of length 1 whose only element is the string.

If sep is empty, `split` splits after each UTF-8 sequence. If both the string and sep are empty, `split` returns an empty list.

* signature: `split(sep <Str>) <List>`
* example: `'a,b,c'.split(',')`

#### `toChars`
`toChars` splits a string into a list of single-rune Strs.

* signature: `toChars() <List>`
* example: `'xyz'.toChars()`

#### `trim`
`trim` returns a new string with all leading and trailing runes contained in cutset removed.

* signature: `trim(<Str>) <Str>`
* example: `'\t\tabc\n'.trim('\t\n')`

_This document uses documentation from [go](https://github.com/golang/go), which
is licensed under the BSD-3-Clause license._
