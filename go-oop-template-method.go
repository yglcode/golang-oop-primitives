package main

import (
	"fmt"
)

//implement traditional OOP template methods design pattern

//Use interface to play the role of
//virtual method table of Java class
//contains methods which need to be overriden/virtual
type Shape interface {
	drawBoundary()
	fillColor()
	draw()
}

//embed interface to define abstract base class in OOP
//1. outer structs (embedding this) will "inherit" these interface methods.
//2. the interface value is nil here, so methods are "abstract".
type ShapeAbstract struct {
	Shape
}

//define logic reused in child classes
func (sa ShapeAbstract) draw() {
	//following template methods design pattern
	//invoke virtual/"abstract" methods (defined in interface)
	sa.drawBoundary()
	fmt.Print("-")
	sa.fillColor()
}

//extends "abstract class" with placeholder methods implementations
type ShapeBase struct {
	*ShapeAbstract
}

//common constructor pattern:
//override embedded Shape interface value with itself - newly created object.
//so interface will take latest overriding methods, exactly how OOP overrides works
func NewShapeBase() *ShapeBase {
	sb := &ShapeBase{&ShapeAbstract{}}
	sb.Shape = sb
	return sb
}

//override abstract method
func (sb *ShapeBase) drawBoundary() {
	//no-op
	fmt.Print("draw nothing")
}

//override abstract method
func (sb *ShapeBase) fillColor() {
	//no-op
	fmt.Print("fill nothing")
}

//extend base class thru embedding
type Circle struct {
	*ShapeBase
}

//in constructor, assign itself - newly created object to embedded Shape interface value.
//so interface will take latest overriding methods.
func NewCircle() *Circle {
	c := &Circle{NewShapeBase()}
	c.Shape = c
	return c
}

//override base method
func (c *Circle) drawBoundary() {
	fmt.Print("Circle")
}

//embed base struct for extension
type RedRectangle struct {
	*ShapeBase
}

//in constructor, assign itself - newly created object to embedded Shape interface value.
//so interface will take latest overriding methods.
func NewRedRectangle() *RedRectangle {
	rr := &RedRectangle{NewShapeBase()}
	rr.Shape = rr
	return rr
}

//override base method
func (rt *RedRectangle) drawBoundary() {
	fmt.Print("Rectangle")
}

//override base method
func (rt *RedRectangle) fillColor() {
	fmt.Print("Red")
}

//embed Circle for extension
type BlueCircleWithText struct {
	*Circle
}

//in constructor, assign itself - newly created object to embedded Shape interface value.
//so interface will take latest overriding methods.
func NewBlueCircleWithText() *BlueCircleWithText {
	bct := &BlueCircleWithText{NewCircle()}
	bct.Shape = bct
	return bct
}

//override
func (bct *BlueCircleWithText) fillColor() {
	fmt.Print("Blue")
}

//override and extend
func (bct *BlueCircleWithText) draw() {
	//extend superclass's draw()
	//since we can embed multiple InnerTypes,
	//it is in fact multiple-inheritance: have to name super explicitly
	bct.Circle.draw()
	//extend with text annotation
	fmt.Print("-TextAnnotation")
}

func main() {
	//create array of shapes and invoke its draw()
	shapes := []Shape{
		NewShapeBase(),
		NewCircle(),
		NewRedRectangle(),
		NewBlueCircleWithText(),
	}
	for _, s := range shapes {
		s.draw()
		fmt.Println()
	}
}
