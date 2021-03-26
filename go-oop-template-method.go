package main

import (
	"fmt"
)

//self referential interface to play the role of
//virtual method table of Java class
//contains methods which need to be overriden/virtual
type Shape interface {
	drawBoundary()
	fillColor()
	draw(Shape)
}

//base struct and methods to be embedded/inherited
type ShapeBase struct{}
//to be overriden
func (sb *ShapeBase) drawBoundary() {
	//no-op
	fmt.Print("draw nothing")
}
//to be overriden
func (sb *ShapeBase) fillColor() {
	//no-op
	fmt.Print("fill nothing")
}
//draw() consumes polymorphic behaviours,
//so make it accept interface as 1st argument
func (_ *ShapeBase) draw(sb Shape) {
	//call methods thru interface for polymorphism
	sb.drawBoundary()
	fmt.Print("-")
	sb.fillColor()
}

//extend base class thru embedding
type Circle struct {
	*ShapeBase
}
func NewCircle() *Circle {
	return &Circle{&ShapeBase{}}
}
//override
func (c *Circle) drawBoundary() {
	fmt.Print("Circle")
}

//embed base struct for extension
type RedRectangle struct {
	*ShapeBase
}
//override
func NewRedRectangle() *RedRectangle {
	return &RedRectangle{&ShapeBase{}}
}
//override
func (rt *RedRectangle) drawBoundary() {
	fmt.Print("Rectangle")
}
//override
func (rt *RedRectangle) fillColor() {
	fmt.Print("Red")
}

//embed base struct for extension
type BlueCircleWithText struct {
	*Circle
}
//override
func NewBlueCircleWithText() *BlueCircleWithText {
	return &BlueCircleWithText{&Circle{}}
}
//override
func (bct *BlueCircleWithText) fillColor() {
	fmt.Print("Blue")
}
//override and extend
func (bct *BlueCircleWithText) draw(s Shape) {
	//extend superclass's draw()
	//since we can embed multiple InnerTypes,
	//it is in fact multiple-dispatch: have to name super explicitly
	bct.Circle.draw(s)
	//extend with text annotation
	fmt.Print("-TextAnnotation")
}

func main() {
	//create array of shapes and invoke its draw()
	shapes := []Shape{
		&ShapeBase{},
		NewCircle(),
		NewRedRectangle(),
		NewBlueCircleWithText(),
	}
	for _, s := range shapes {
		//note: we have to pass in thru interface value 
		//for polymorphism
		s.draw(s)
		fmt.Println()
	}
}
