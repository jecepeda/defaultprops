# defaultprops

Defaultprops is a library that allows two identical go structs to be merged.

## How to

The first thing you need to do is to create two structures. Note that the structs must be passed by argument

```go
package main

import "github.com/jecepeda/defaultprops"

type MyStruct struct {
    Foo string
    Bar int
}

func main(){
    a := MyStruct{Foo: "foobar"}
    b := MyStruct{Bar: 2}
    defaultprops.SubstituteNonConfig(&a, &b)
    fmt.Println(b.Foo, b.Bar) // foobar 2
}
```
