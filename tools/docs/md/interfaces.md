# The Golem Programming Language Reference

## Interfaces

Interfaces define the various intrinsic capabilities of Golem's standard types.

### Comparable

A value is comparable if it supports the comparison operators 
`>`, `>=`, `<`, `<=`, `<=>`.  

[Str](str.html), [Int](int.html), [Float](float.html), and [Bool](bool.html) are comparable.

### Hashable

A value is hashable if it can be a key in a [Dict](dict.html) or 
[Set](set.html).  The builtin function [hashCode()](builtins.html#hashcode) 
returns the hashCode of a hashable value.  

[Str](str.html), [Int](int.html), [Float](float.html), [Bool](bool.html), 
and [Tuple](tuple.html) are hashable. 

### Indexable

A value is indexable if it supports the index operator `a[x]`.   

[Str](str.html), [List](list.html), [Range](range.html), [Tuple](tuple.html) 
and [Dict](dict.html) are indexable.

### Iterable

A value is iterable if its entries can be iterated over.  Iterable values can be the
subject of a `for` loop.  

The builtin function [iter()](builtins.html#iter) returns a [Struct](struct.html) 
that can be used to iterate over the entries in an iterable value.  

The builtin function [stream()](builtins.html#stream) returns a [Struct](struct.html) 
that can be used to perform a series of transforms on a sequence of iterated values, 
and then collects the values into a final result.

[Str](str.html), [List](list.html), [Tuple](tuple.html), [Range](range.html), [Dict](dict.html) 
and [Set](set.html) are iterable.

### Lenable

A value is lenabale if it has a length. The builtin function [len()](builtins.html#len) 
returns the length of a lenable value.  

[Str](str.html), [List](list.html), [Range](range.html), [Tuple](tuple.html), 
[Dict](dict.html) and [Set](set.html) are lenable.

### Sliceable

A value is sliceable if it supports the slice operators `a[x:y]`, `a[x:]`, `a[:y]`.  

[Str](str.html) and [List](list.html) are sliceable.

