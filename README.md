# The Quetzal Profarting Language

This was a weekend project because what the heck. Actually, this was a weekend project because an incompetent Concurrent Programming professor made us do a syntactic analyser as a final project... but I decided to be really stupid and make an interpreter that actually works. Sorta.

I've never had a class on compilers, so I don't know anything about interpreters and all that fancy stuff, except for what I've learned by poking at other people's code (that, is almost nothing).

Strings are **REALLY** broken, and if you ever try to do anything useful with it, it will probably eat you alive.

Also, _never, ever, ever_ read the source. Seriously, don't do it. It's the worst piece of shit you will never see, because you'll be a well–behaved person and you _will not look at it_, because I told you not to. Ok?

Okay, I am assuming you're back to reading this, after you went and read the source. Look at you, you subversive little shit. Don't flame me if you did—I _did_ warn you.

Oh yeah, error messages are in spanish because I'm Mexican.

## Syntax

There are only three keywords. There are no comments. There are no operators.

* `let`: Defines variables
* `proc`: Defines procedures
* `return`: Returns values from procedures

All instructions must end in a semicolon. There is no fancy "everything is an expression" concept here; the only expressions are calling procedures that return something and assigning that to variables.

Procedures are defined like this:

    proc do_some_shit(param1 int, param2 string) int {
        return param1;
    }
    
    proc no_return_no_params() {
    }

Variables are defined like this:

    let a int = 5;
    let b string = "foo";
    let c int = do_some_shit(a, "bar");
    
You can pass literals or variables to procedures. <-- Yeah, it's that awesome.

You can't nest procedure calls, so everything is super verbose.
    
## Type System

There's only two types: `int` and `string`. There's no variable re–assignment. Also, you can't return literals from functions because I'm a lazy fuck and didn't implement that.

There's no type inference; all variables must specify their types.

There is _some_ type safety: You can't pass variables of mismatched types to procedures, and you can't assign results from procedures to variables of mismatched types.

All variables are scoped locally to their parent procedure.

## Runtime

Programs must have a `main` procedure, and all procedures and variables must be defined before they are used.

## Built–ins

There's a few built–in procedures:

* `suma`: Add
* `sub`: Subtract
* `mult`: Multiply
* `print`: Print (duh)

All of these take two parameters except for `print` which takes one, and is the _only_ type–flexible procedure (it takes either `int` or `string`). All of these return `int` except for `print` which returns nothing. I don't know what happens if you try to assign the value of `print` to a variable and then try to use it, so your computer might explode if you do that.

## Usage

Compile it, and run it passing the resulting binary a path to a file as the first argument.

Example:

```bash
go build -o quetzal main.go
./quetzal file.quetzal
```

## Example

    proc square(a int) int {
        let b int = mult(a, a);
        return b;
    }
    
    proc main() {
        let c int = square(5);
        print("This should return 25:\n");
        print(c);
    }
    
**FIN**