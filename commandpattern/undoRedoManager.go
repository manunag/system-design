package main

import (
	"fmt"
)

// I have implemented a slight modified version of stack

type ItemType interface{}

//Stack - Stack of items.
type Stack struct {
	items        []ItemType
	currentIndex int
}

// New - Creates a new Stack.
func (stack *Stack) New() *Stack {

	stack.items = []ItemType{}
	stack.currentIndex = -1

	return stack
}

// Push - Adds an Item to the top of the stack (Using the current Index). I also deletes the
func (stack *Stack) Push(t ItemType) {

	//Initialize items slice if not initialized
	if stack.items == nil {
		stack.items = []ItemType{}
	}

	if len(stack.items)-1 > stack.currentIndex {
		stack.items = stack.items[0:stack.currentIndex]
	}

	// Performs append operation.
	stack.items = append(stack.items, t)
	stack.currentIndex = len(stack.items) - 1

}

// Pop is slightly a different version of usual Pop operation. Instead of removing the top element,
//It moves the pointer to the previous Top
func (stack *Stack) Pop() {
	stack.currentIndex -= 1
}

func (stack *Stack) Top() ItemType {
	item := stack.items[stack.currentIndex]

	return item
}

type Command interface {
	execute()
}

type Button struct {
	command Command
}

func (b *Button) press() {
	b.command.execute()
}

type OperationType interface {
	undo()
	redo()
}

type undoOperation struct {
	operationType OperationType
}

func (c *undoOperation) execute() {
	c.operationType.undo()
}

type redoOperation struct {
	operationType OperationType
}

func (c *redoOperation) execute() {
	c.operationType.redo()
}

type TextEditor struct {
	data    string
	history *Stack
}

func NewTextEditor() *TextEditor {
	var history Stack
	var data string

	return &TextEditor{data: data, history: history.New()}
}

func (te *TextEditor) write(a string) {
	te.data = te.data + " " + a
	te.history.Push(te.data)

	fmt.Println("Writing to file")
}

func (te *TextEditor) print() {
	fmt.Println("Printing Data: ", te.data)
}

func (te *TextEditor) undo() {
	te.history.Pop()
	te.data = te.history.Top().(string)
	fmt.Println("Undo Complete")
}

func (te *TextEditor) redo() {
	if te.history.currentIndex < len(te.history.items)-1 {
		te.history.currentIndex += 1
	} else {
		fmt.Println("Redo Unsuccessful")
		return
	}

	te.data = te.history.Top().(string)
	fmt.Println("Redo complete")
}

func main() {
	te := NewTextEditor()

	undo := &undoOperation{
		operationType: te,
	}

	undoButton := &Button{
		command: undo,
	}

	redo := &redoOperation{
		operationType: te,
	}

	redoButton := &Button{
		command: redo,
	}

	te.write("a")
	te.write("b")
	te.write("c")
	te.print()
	undoButton.press()
	te.print()
	redoButton.press()
	te.print()
	te.write("d")
	te.print()
	undoButton.press()
	te.print()
}
