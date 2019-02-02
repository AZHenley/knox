# The Knox Programming Language

Knox is an experimental language as I learn Go and explore compiler design. The compiler is written in Go and generates Go. The compiler consists of a lexer, recurseive-descent parser, semantic analyzer, and code emitter.

The principles and major features behind the design of Knox are:
 - Explicitness. Explicit and unambiguous code is a priority, even over brevity. No surprises.
 - Small language. Simple and consistent syntax with few constructs, as an alternative to Go.
 - Pass by reference. All complex types should be pass by reference and pointers should be hidden.
 - Convenient data structures. Lists and maps are first-class data structures that should be as easy as Python.
 - Operability. Use any Go library and produce Go libraries.
 - Well-behaved. Contracts, error handling, and unit testing are first-class constructs.

```
Example goes here.
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
 
 
 - Add assertions
 - Method overloading?
 - Parameter generics?
 - Optional parameters?
 - Named parameters?
 - Custom operators
 - Remove varargs?
 - defer?
 - Constructor? Destructor?
 - Change syntax for make/new
 - Add implements syntax for interfaces
 - Nested functions?
 - Ranges syntax, not a range function
 - No bitwise operators
 - Pass by reference
 - Design by contracts?
 - Give a variable/type properties, like an int must be even or in some range


## Decisions we made:
  - No optional, default, or named parameters
  - No method overloading
  - No implicit casting
  - Generics for parameters only (and return)
  - let name:type = value
  - func name(a:int, b:int, c:int) out (x:int, y:int) { }
  - No return values, only out 
  - _ throws away a return value
  - When calling a function, all return values must be assigned to a variable (or ignored with _)
  - Use Go's standard library
  - Maps and lists are builtin
  - Change type to alias
  - Change struct and interface syntax to struct X {} and interface Y {}
  - No constructor or destructor
  - Class fields can have defaults
  - Change struct to class, methods go inside?, no access modifiers, no inheritance 
  - No goto
  - Static variables
  - No custom or overloaded operators
  - Extension methods by doing func type.method


## Examples:

```
let d:int, e:string = doSomething(1, 2);
let foo:int
let bar:string
foo, bar = doSomething()

let bizz:int, bazz:string = doSomething();

func doSomething (a:int, b:int, c?:int) out (d:int, e:string) {
	d = a + b + (c ?? 0);
	e = i_to_s(d);
}
```
## Notes on unit testing
```
@assert in(1, 2, 3) out(6, "6")
@assert in(1, 2) out(3, "3")
func doSomething (a:int, b:int, c?:int) out (d:int, e:string) {
	d = a + b + (c ?? 0);
	e = i_to_s(d);
}
```
 - `knoxc filename.kx` will fail if any assert fails
 - Use `knoxc -coverage filename.kx` to verify the coverage of the file
 - Use `knoxc -coverage filename.kx:doSomething` to verify the coverage of a function

## Notes on generic parameters
```
generic T, U
func DoStuff(mom:T, dad:T, homie:U) {}

DoStuff<int, string>(5, 3, "lol")
DoStuff<bool, float>(true, false, 3.5)
DoStuff<bool, float>(false, false, 0.0)
```
so the compiler will generate 2 separate versions... `DoStuff_int_string` and `DoStuff_bool_float`
```
generic T
func DuplicateItemsAndReturnList(item:T, count:int) out (myList T[]) {
  for i in range count {
     myList.append(item)
  }
}

let lolz:string[] = DuplicateItemsAndReturnList<string>("lol", 18)
```

## Misc notes

This Knox code...
```
func doSomething(a:int) out (b:string, c:int) {
  // do a little logic here
  b = "some words"
  // do some more logic here
  c = 0
}
```
Compiles to...
```
type doSomethingRetStruct struct {b string; c int}
func doSomething(a int) : doSomethingRetStruct {
  doSomethingRetStruct {}
  // do a little logic here
  doSomethingRetStruct.b = "some words"
  // do some more logic here
  doSomethingRetStruct.c = 0
}
```

