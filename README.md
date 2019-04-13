# The Knox Programming Language

Knox is an experimental language meant to help me learn Go and explore compiler design. The compiler is written in Go and generates C. The compiler consists of a lexer, recurseive-descent parser, semantic analyzer, and code emitter.

The principles behind the design of Knox are:
 - Explicitness. Explicit and unambiguous code is a priority, even over brevity. No surprises.  
 - Pass by reference. All complex types should be pass by reference and pointers should be hidden.
 - Small language. Simple and consistent syntax with few constructs, as an alternative to Go.
 - Convenient data structures. Strings, lists, and maps are first-class data structures that should be as easy as Python.
 - Operability. Use any C library and produce C libraries.
 - Automatic reference counting. Avoid manual memory management without garbage collection pauses.
 - Fast enough. Compiling time, execution time, and memory usage should be comparable to Go.  
 - Well-behaved. Contracts, error handling, and unit testing are first-class constructs.

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
 - Constructors
 - Classes must explicitly implement interfaces
 - No pointers
 - All return values must be used or explicitly thrown away
 - No goto
 - Multiple assignment is only for multiple return values
 


