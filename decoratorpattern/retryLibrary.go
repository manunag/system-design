package main

import (
	"errors"
	"fmt"
	"time"
)

type sequence interface {
	next() int
	retry(func() error, int) (interface{}, error)
}

type even struct {
	Number         int
	RetryInterface RetryInterface
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

func (e *even) retry(a func() error, tries int) (interface{}, error) {
	wrappedFunc := func() error {
		fmt.Println("sleep", e.Number, "seconds")

		time.Sleep(time.Duration(e.Number) * time.Second)
		e.next()

		return a()
	}

	err := e.RetryInterface.retry(wrappedFunc, tries)
	if err != nil {
		return "failure", err
	}

	return "success", nil
}

func newEven(RetryInterface RetryInterface) sequence {
	return &even{Number: 0, RetryInterface: RetryInterface}
}

type odd struct {
	Number         int
	RetryInterface RetryInterface
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

func (o *odd) retry(a func() error, tries int) (interface{}, error) {
	wrappedFunc := func() error {
		fmt.Println("sleep", o.Number, "seconds")

		time.Sleep(time.Duration(o.Number) * time.Second)
		o.next()

		return a()
	}

	err := o.RetryInterface.retry(wrappedFunc, tries)
	if err != nil {
		return "failure", err
	}

	return "success", nil
}

func newOdd(RetryInterface RetryInterface) sequence {
	return &odd{Number: 1, RetryInterface: RetryInterface}
}

type fibonacci struct {
	Number         int
	Previous       int
	RetryInterface RetryInterface
}

func newFibonacci(RetryInterface RetryInterface) sequence {
	return &fibonacci{Number: 0, Previous: 0, RetryInterface: RetryInterface}
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

func (f *fibonacci) retry(a func() error, tries int) (interface{}, error) {
	wrappedFunc := func() error {
		fmt.Println("sleep", f.Number, "seconds")

		time.Sleep(time.Duration(f.Number) * time.Second)
		f.next()

		return a()
	}

	err := f.RetryInterface.retry(wrappedFunc, tries)
	if err != nil {
		return "failure", err
	}

	return "success", nil
}

type RetryInterface interface {
	retry(func() error, int) error
}

type retrier struct{}

func newRetrier() *retrier {
	return &retrier{}
}

func (r *retrier) retry(a func() error, tries int) error {
	var err error
	for i := 0; i < tries; i++ {
		err = a()
		if err == nil {
			return nil
		}

		fmt.Println("Retrying.......")
	}

	fmt.Println("Max retry reached.......")
	return err
}

func main() {
	retryBase := newRetrier()
	fibDec := newFibonacci(retryBase)
	evenDec := newEven(retryBase)
	oddDec := newOdd(retryBase)

	a := func() error {
		fmt.Println("My function that does nothing :)")
		return errors.New("throwing error for testing")
	}

	fmt.Println("\n\n Base retry")
	fmt.Println(retryBase.retry(a, 5))

	fmt.Println("\n\n Fibonacci retry")
	fmt.Println(fibDec.retry(a, 5))

	fmt.Println("\n\n Even retry")
	fmt.Println(evenDec.retry(a, 5))

	fmt.Println("\n\n Odd retry")
	fmt.Println(oddDec.retry(a, 5))
}
