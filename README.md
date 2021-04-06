# Gocalc
[![godoc](https://godoc.org/github.com/Rollylni/Gocalc?status.svg)](https://pkg.go.dev/github.com/Rollylni/Gocalc)
[![license](https://img.shields.io/github/license/rollylni/gocalc?style=flat-square)](https://en.wikipedia.org/wiki/MIT_License)

Simple calculator in Go!

## Install
```bash
go get github.com/Rollylni/Gocalc
```

## Example
```go
package example

import (
     "fmt"
     gocalc "github.com/Rollylni/Gocalc"
)

func calc(input string) {
     calc := gocalc.NewCalculator(input)
     res, err := calc.Process()
     if err != nil {
         fmt.Println(err)
     } else {
         fmt.Println(res) 
     }
}

func main() {
     calc("2 + 2 * 2)
     calc("PI + 1.0")
     calc("1 + 24.00")
}

```
Output
```bash
> 6
> 4.141593
> Operation error: Type Error!
```
