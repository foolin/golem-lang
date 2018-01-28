# Golem Language Reference

## Types

### Null

    Null represents the absence of a value.  

    Intrinsic Functions: none

### Bool

    Bool represents truth or falsity.

    Intrinsic Functions: none

### Str

    A Str is a sequence of unicode characters, a.k.a 'runes'.  

    Multiline strings can be created by using the backtick charcter '\`' as the delimiter.

    Unicode escape sequences are declared using a `\u`.  For example, to encode "💖",
    you would use `\u{1F496}`.

    Intrinsic Functions: 

        contains(str)
            Return whether this string contains the given string.

        index(str)
            Return the starting index of the given string in this string, or -1.

        startsWith(str)
            Return whether this string starts with the given string.

        endsWith(str)
            Return whether this string ends with the given string.

        replace(old, new, <max>)
            Replace all the occurences of the old string with the new string.  If
            the optional 'max' paramter is supplied, then a maximum of that many 
            replacements are done.

### Int

    An Int is a signed, 64-bit integer.

    Intrinsic Functions: 
        **TODO** none yet! 

### Float

    A Float is a 64-bit floating point value.

    Intrinsic Functions: 
        **TODO** none yet! 

### Func

    A Func is a function -- a sequence of instructions that perform a task and return a value.

    Intrinsic Functions: none

### List

    A List is an ordered sequence of values.

    Intrinsic Functions: 
        **TODO** create a bunch more

        add(elem)
            Append an element to the end of the list.

        addAll(coll)
            Append all of the elements in an iterable collection to the end of the list.

        clear()
            Remove all the elements from the list.

        isEmpty()
            Return whether or not the list is empty.

        contains(elem)
            Return whether or not the list contains the given element.

        indexOf(elem)
            Return the index of the given element in the list, or -1

        join(<sep>)
            Return a string which is the result of concatenating together the string
            representation of every element in the list.  If the optional 'separator'
            parameter (which must be a string) is passed in, it is used to separate
            each pair of elements in the string.

        map(func)
            Create a new list which is the result of applying the given function
            to each element of the existing list.

        reduce(initial, func)
            Reduce the list to a single value, which is the accumulated result of 
            applying the given function to each element in the existing list.  Each 
            successive invocation of the function is supplied the return value 
            of the previous invocation.  The given function must take 2 parameters.
            For example:
                let ls = [1, 2, 3, 4, 5]
                let sum = ls.reduce(0, |acc, x| => acc + x)

        filter(predFunc)
            Create a new list which is the result of applying the given predicate function
            to each element of the existing list, to see if we want the element in the
            new list. Note that 'predFunc' must always return a Bool.

        remove(index)
            Return the element at the given index from the list.

### Range
    
    A Range is an immutable representation of a sequence of integers.

    Intrinsic Functions: 

        from()
            Return the first integer of the range (inclusive).

        to()
            Return the last integer of the range (exclusive).

        step()
            Return how far apart each pair of intergers are from each other.

        count()
            Return the total number of intergers in the range.

### Tuple

    A Tuple is an immutable ordered sequence of values.

    Intrinsic Functions: none

### Dict

    A Dict is a hash table, also known as an associative array.

    The keys can be any value that supports hashing (currently str, int, float, or bool). 
    **TODO** A future version of Golem will allow for structs to act as a dict key.

    Intrinsic Functions: 
        **TODO** create a bunch more

        addAll(coll)
            Append all of the elements in an iterable collection to the dict.  Each
            element in the list must be a 2-tuple representing (key, value). 

        clear()
            Remove all the elements from the dict.

        isEmpty()
            Return whether or not the dict is empty.

        containsKey(key)
            Return whether or not the dict contains the given key.

        remove(key)
            Return the entry having the given key from the dict, if it is present.
            Return whether it was present or not.

### Set

    A Set is a collection of unique hashable values.

    **TODO** A future version of Golem will allow for structs to act as an element
    in a set.

    Intrinsic Functions: 
        **TODO** create a bunch more

        add(elem)
            Add an element to the set, if its not already there.

        addAll(coll)
            Add all of the elements in an iterable collection to the set.

        clear()
            Remove all the elements from the set.

        isEmpty()
            Return whether or not the set is empty.

        contains(elem)
            Return whether or not the set contains the given element.

        remove(elem)
            Return the given element from the set, if it is present.
            Return whether it was present or not.

### Struct

    A Struct is a hashtable-like data structure that is used to aggregate values together
    into a collection of named key-value pairs.  Each key in a struct must be a valid
    Golem identifier.  Structs can be merged together via the merge() builtin function.

    Intrinsic Functions: none

### Chan

    A Chan is a channel that can be used to pass messages between goroutines.  Chans
    are created via the chan() builtin function.

    Intrinsic Functions: 

        send(val)
            Send a value to the channel.

        recv()
            Receive a value from the channel.

## Builtin Functions

    assert(val)
        Assert that the given value (which must be boolean) is true.  If it is false,
        an error is thrown

    chan(<size>)
        Create a new channel.  If the optional 'size' parameter (which must be 
        an int) is passed in, it is used to create a buffered channel.

    len(val)
        Return the size of any collection or string.
        
    merge(...)
        Merge together an arbitrary number of structs.  See the tutorial for a detailed
        explanation of how this works. **TODO* explain it here too.

    println(...)
        Print an arbitrary number of values, followed by a newline.

    print(...)
        Print an arbitrary number of values.

    range(from, to, <step>)
        Create a Range from the 'from' paramter (inclusive), to the 'to' parameter
        (exclusive).  If the optional 'step' parameter (which must be 
        an int) is passed in, it is used to step further that just 1 integer at a time.
        If 'to' is less than 'from', and a negative value for 'step' is passed in,
        then the range will count backwards.

    str(val)
        Return the string representation of any value.

    type(val)
        Return a string describing the type of any value

    freeze(val)
        Make a value become immutable, if it isn't already.  Return the value that
        was passed in.

    frozen(val)
        Return whether a value is immutable or not.

## Standard Modules

**TODO** work on these has barely begun.

    sys
        exit(val)
            Exit the interpreter's process with the given exit code.

    io
        File(name) Create a file:
            isDir()
            items()
            readLines()
            name
            
    regex
        compile(str) Create a regex:
            match(str)

