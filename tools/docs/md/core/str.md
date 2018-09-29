## Str

Str is the set of all valid sequences of UTF-8-encoded "code points", otherwise
known as "runes".  Strs are immutable.

String literals can be surrounded by either single quotes or double quotes.  Backticks
can also be used to create mulit-line strings.

Valid operators for Str are:

* The equality operators `==`, `!=`,
* The [`comparision`](interfaces.html#comparable) operators `>`, `>=`, `<`, `<=`, `<=>`
* The [`index`](interfaces.html#indexable) operator `a[x]`
* The [`slice`](interfaces.html#sliceable) operators `a[x:y]`, `a[x:]`, `a[:y]`

The index operator always returns a single-rune Str.

The slice operators always return a Str.

Strs are
[`hashable`](interfaces.html#hashable),
[`lenable`](interfaces.html#lenable) and
[`iterable`](interfaces.html#iterable).

A Str has the following fields:

* [contains](#contains)
* [hasPrefix](#hasprefix)
* [hasSuffix](#hassuffix)
* [index](#index)
* [lastIndex](#lastindex)
* [map](#map)
* [replace](#replace)
* [split](#split)
* [toChars](#tochars)
* [trim](#trim)

### `contains`

`contains` reports whether a substring is within a string.

* signature: `contains(substr <Str>) <Bool>`
* example: `'abcdef'.contains('de')`

### `hasPrefix`

`hasPrefix` tests whether a string begins with a prefix.

* signature: `hasPrefix(prefix <Str>) <Bool>`
* example: `'abcdef'.hasPrefix('ab')`

### `hasSuffix`

`hasSuffix` tests whether a string ends with a suffix.

* signature: `hasSuffix(suffix <Str>) <Bool>`
* example: `'abcdef'.hasSuffix('ab')`

### `index`

`index` returns the index of the first instance of a substring in a string.
or -1 if the substring is not present.

* signature: `index(substr <Str>) <Int>`
* example: `'abcab'.index('ab')`

### `lastIndex`

`lastIndex` returns the index of the last instance of a substring in a string,
or -1 if the substring is not present.

* signature: `lastIndex(substr <Str>) <Int>`
* example: `'abcab'.lastIndex('ab')`

### `map`

`map` returns a copy of the string with all its characters modified according to
the mapping function.

The mapping function must accept one Str parameter, and must return a Str.

* signature: `map(mapping <Func>) <Str>`
* mapping signature: `fn(s <Str>) <Str>`
* example:

```
let s = 'abc(def)[x,y,z]'
let t = s.map(fn(c) {
    return c >= 'a' && c <= 'z' ? c : ''
})
println(t)
```

### `replace`

`replace` returns a copy of a string with the first n non-overlapping instances
of `old` replaced by `new`. If `old` is empty, it matches at the beginning of a string
and after each UTF-8 sequence, yielding up to k+1 replacements for a k-rune string.
If `n` < 0, there is no limit on the number of replacements.  The parameter `n` is
optional, and defaults to -1.

* signature: `replace(old <Str>, new <Str>, n = -1 <Int>) <Int>`
* example: `'abcab'.replace('a', 'x')`

### `split`

`split` slices a string into all substrings separated by sep and returns a list
of the substrings between those separators.

If the string does not contain sep and sep is not empty, `split` returns a list
of length 1 whose only element is the string.

If sep is empty, `split` splits after each UTF-8 sequence. If both the string
and sep are empty, `split` returns an empty list.

* signature: `split(sep <Str>) <List>`
* example: `'a,b,c'.split(',')`

### `toChars`

`toChars` splits a string into a list of single-rune Strs.

* signature: `toChars() <List>`
* example: `'xyz'.toChars()`

### `trim`

`trim` returns a new string with all leading and trailing runes contained in cutset removed.

* signature: `trim(<Str>) <Str>`
* example: `'\t\tabc\n'.trim('\t\n')`

