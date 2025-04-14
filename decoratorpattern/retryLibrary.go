package main

import (
	"errors"
	"fmt"
	"time"
)

type sequence interface {
	next() int
}

type even struct {
	Number int
}

func (e *even) next() int {
	e.Number += 2
	return e.Number
}

func limitedNextEven(seq *even, limit int) func() int {
	firstCall := true
	return func() int {
		if firstCall {
			firstCall = false
			return seq.Number
		}
		if limit == -1 || seq.Number+2 <= limit {
			seq.next()
		}

		return seq.Number
	}
}

func newEven() sequence {
	return &even{Number: 0}
}

type odd struct {
	Number int
}

func (o *odd) next() int {
	o.Number += 2
	return o.Number
}

func limitedNextOdd(seq *odd, limit int) func() int {
	firstCall := true

	return func() int {
		if firstCall {
			firstCall = false
			return seq.Number
		}
		if limit == -1 || seq.Number+2 <= limit {
			seq.next()
		}

		return seq.Number
	}
}

func newOdd() sequence {
	return &odd{Number: 1}
}

type fibonacci struct {
	Number   int
	Previous int
}

func newFibonacci() sequence {
	return &fibonacci{Number: 0, Previous: 0}
}

func (f *fibonacci) next() int {
	if f.Previous == 0 && f.Number == 0 {
		f.Number = 1
	} else if f.Previous == 0 && f.Number == 1 {
		f.Previous = 1
		f.Number = 1
	} else {
		temp := f.Previous
		f.Previous = f.Number
		f.Number = temp + f.Number
	}

	return f.Previous
}

func retry(fun func() error, nextFunc func() int, tries int) error {

	var err error
	for i := 0; i < tries; i++ {

		err = fun()
		if err == nil {
			return nil
		}

		fmt.Println("Retrying.......")

		current := nextFunc()
		fmt.Println("sleep", current, "seconds")

		time.Sleep(time.Duration(current) * time.Second)
	}

	fmt.Println("Max retry reached.......")
	return err
}

func main() {
	fibGen := newFibonacci()
	evenGen := newEven()
	oddGen := newOdd()

	fun := func() error {
		fmt.Println("My function that does nothing :)")
		return errors.New("throwing error for testing")
	}

	fmt.Println("\n\n Base retry")
	fmt.Println(retry(fun, func() int {
		return 0
	}, 5))

	fmt.Println("\n\n Fibonacci retry")
	fmt.Println(retry(fun, fibGen.next, 5))

	fmt.Println("\n\n Even retry")
	fmt.Println(retry(fun, evenGen.next, 5))

	fmt.Println("\n\n Odd retry")
	fmt.Println(retry(fun, oddGen.next, 5))
}
