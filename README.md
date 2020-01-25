# The Knox Programming Language

Knox is an experimental language meant to help me learn Go and explore compiler design. It acts as a systems language with high-level constructs for convenience. The compiler is written in Go and generates C. It is very early in development.

The principles behind the design of Knox are:
 - Explicitness. Explicit and unambiguous code is a priority, even over brevity. No surprises.  
 - Pass by reference. All complex types should be pass by reference and pointers should be hidden, like Java and C#.
 - Small language. Simple and consistent syntax with few constructs, as an alternative to Go.
 - Convenient data structures. Strings, lists, and maps are first-class data structures that should be as easy as Python.
 - Operability. Use any C library and produce C libraries.
 - Automatic reference counting. Avoid manual memory management without garbage collection pauses.
 - Easy to setup and use. No massive installation like C# or Java and no annoying configuration like Go's gopath.
 - Fast enough. Compiling time, execution time, and memory usage should be comparable to Go (but probably better!).  
 - Well-behaved. Contracts, error handling, and unit tests are first-class constructs.

```
func main() void {
    fizzbuzz(300);
}

func fizzbuzz(n : int) void {
    for i : int in stl.range(1,10,1) {
        if i%15 == 0 {
            stl.print("FizzBuzz");
        } else if i%3 == 0 {
            stl.print("Fizz");
        } else if i%5 == 0 {
            stl.print("Buzz");
        } else {
            stl.print(i);
        }
    }
}
```

## Comparison to Go:
 - Classes instead of structs
 - Objects are pass-by-reference
 - Ada-style type constraints
 - No type inference
 - No short form of variable declarations
 - No variable declaration blocks
 - No implicit casting
 - Variables must be initialized
 - Semicolons required
 - Different syntax for variable and function declarations
 - Python-style While and For loops
 - Allows whitespace between if and elseif blocks
 - Enum support
 - Sum types
 - Constructors
 - Classes must explicitly implement interfaces
 - No pointers
 - All return values must be used or explicitly thrown away
 - No goto
 - Multiple assignment is only for multiple return values
 


