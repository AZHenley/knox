# The Knox Programming Language

Knox is an experimental language meant to help me learn Go and explore compiler design. The compiler is written in Go and generates Go. The compiler consists of a lexer, recurseive-descent parser, semantic analyzer, and code emitter.

The principles behind the design of Knox are:
 - Explicitness. Explicit and unambiguous code is a priority, even over brevity. No surprises.  
 - Pass by reference. All complex types should be pass by reference and pointers should be hidden.
 - Small language. Simple and consistent syntax with few constructs, as an alternative to Go.
 - Convenient data structures. Lists and maps are first-class data structures that should be as easy as Python.
 - Operability. Use any Go library and produce Go libraries.
 - Fast enough. Compiling time, execution time, and memory usage should be less than a 25% increase over Go.  
 - Well-behaved. Contracts, error handling, and unit testing are first-class constructs.

```
func main() void {
    fizzbuzz(300);
}

func fizzbuzz(n : int) void {
    var i : int = 1;
    while i <= n {
        if i%15 == 0 {
            print("FizzBuzz");
        } else if i%3 == 0 {
            print("Fizz");
        } else if i%5 == 0 {
            print("Buzz");
        } else {
            print(i);
        }
        i = i + 1;
    }
}
```

## Comparison to Go:
 - No type inference
 - No short form of variable declarations
 - No variable declaration blocks
 - No implicit casting
 - Variables must be initialized
 - Semicolons required
 - Different syntax for variable and function declarations
 - C-style For and While loop syntax
 - Allows whitespace between if and elseif blocks
 - Different type and struct syntax
 - Enum support
 - Classes instead of structs
 - Classes are pass-by-reference
 - Ada-style type constraints
 - Constructors
 - Explicit interfaces only
 - No pointers
 - All return values must be used or explicitly thrown away
 - No goto
 - Multiple assignment is only for multiple return values
 


