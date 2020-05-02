# undo

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/chrboe/undo)

`undo` is a very small Go package which aims to abstract the concept of executing
a function while being able to roll back its effects later. To demonstrate,
here is an illustrative example:

```go
import (
	"github.com/chrboe/undo"
)

func someProcess() error {
	i := 0

	action := undo.New(func() error {
		// This is our action, we execute the desired logic here
		i++
		return nil
	}, func() {
		// This is our rollback function; it should undo the changes
		i--
	})

	err := action.Execute()
	if err != nil {
		// Oops, our action failed. Do something about it
	}

	// After the action is executed, we queue its Undo function to be
	// called when this function returns.
	defer action.Undo()

	// ...
	// Here comes the rest of our logic...
	// If something were to return here, our action would get undone.
	// ...

	// When we are sure we want to keep the changes, we can commit the action.
	// The action will be "defused" by this.
	action.Commit()

	// At this point we can safely return, the Undo function will not be called.
	return nil
}
```

## Chaining actions

There is also a convenient helper to handle a common pattern where multiple actions
are executed one after the other, rolling back step-by-step if one of them fails.

```go
package main

import (
	"fmt"
	"errors"
	"github.com/chrboe/undo"
)

func main() {
	action1 := undo.New(func() error {
		fmt.Println("running action1")
		return nil
	}, func() {
		fmt.Println("rolling back action1")
	})

	action2 := undo.New(func() error {
		fmt.Println("running action2")
		return nil
	}, func() {
		fmt.Println("rolling back action2")
	})

	action3 := undo.New(func() error {
		fmt.Println("running action3")
		return errors.New("oops! something went wrong")
	}, func() {
		fmt.Println("rolling back action3")
	})

	err := undo.ExecuteChain(action1, action2, action3)
	fmt.Println(err)
}
```

This would print:

```
running action1
running action2
running action3
rolling back action2
rolling back action1
oops! something went wrong
```
