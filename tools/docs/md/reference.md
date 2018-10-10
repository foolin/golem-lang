# The Golem Programming Language Reference

## Types

* Basic Types:
  * [Null](null.html)
  * [Bool](bool.html)
  * [Int](int.html)
  * [Float](float.html)
  * [Str](str.html)
* Composite Types:
  * [List](list.html)
  * [Range](range.html)
  * [Tuple](tuple.html)
  * [Dict](dict.html)
  * [Set](set.html)
  * [Struct](struct.html)
* Miscellaneous Types:
  * [Func](func.html)
  * [Chan](chan.html)

## Interfaces

* [Comparable](interfaces.html#comparable)
* [Hashable](interfaces.html#hashable)
* [Indexable](interfaces.html#indexable)
* [Iterable](interfaces.html#iterable)
* [Lenable](interfaces.html#lenable)
* [Sliceable](interfaces.html#sliceable)

## Builtin Functions

### Sandbox Builtins:

* [arity()](builtins.html#arity)
* [assert()](builtins.html#assert)
* [chan()](builtins.html#chan)
* [fields()](builtins.html#fields)
* [freeze()](builtins.html#freeze)
* [frozen()](builtins.html#frozen)
* [has()](builtins.html#has)
* [hashCode()](builtins.html#hashcode)
* [iter()](builtins.html#iter)
* [len()](builtins.html#len)
* [merge()](builtins.html#merge)
* [range()](builtins.html#range)
* [str()](builtins.html#str)
* [stream()](builtins.html#stream)
* [type()](builtins.html#type)

### SideEffect Builtins:

* [print()](builtins.html#print)
* [println()](builtins.html#println)

## Syntax

* [Operator Precedence](syntax.html#operator-precedence)
* [Expressions](syntax.html#expressions)
* [Statements](syntax.html#statements)

## Standard Library

### Sandbox Library:

* [encoding](lib_encoding.html)
  * [encoding.json](lib_encodingjson.html)
* [golem](lib_golem.html)
* [regexp](lib_regexp.html)

### SideEffect Library:

* [io](lib_io.html)
  * [io.ioutil](lib_ioioutil.html)
* [os](lib_os.html)
  * [os.exec](lib_osexec.html)
* [path](lib_path.html)
  * [path.filepath](lib_pathfilepath.html)
