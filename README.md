## Golang OOP primitives ##

### _and combine them to do traditional OOP thru "Embedded Interface"_ ###
   

Golang is not a traditional object oriented programming language. Instead, it distilled a few OO programming primitives and allow you to compose them to achieve different OO designs.

### **1. Methods (or method-set): for "abstract data types"** ###

   In traditional OOP, methods are inherently bound with class and objects.
   In Go, methods can be defined for any ["named"/"defined"](https://golang.org/ref/spec#Type_definitions) types. Instead of everything is an object as in some OO language, everything (almost) can be attached methods.

   So we can have methods defined for integers:
```go
	type MyInt int
	func (mi MyInt) addMore(more MyInt) MyInt {
		return mi+more
	}
```   
Please note there is no "wrapping" objects needed for primitive types as in Java.

Or methods defined for functions:
```go   
	type HTTPHandler func(req *http.Request, resp http.Response)
	func (hh HTTPHanlder) handle(req *http.Request, resp http.Response) {
		hh(req,resp);
	}
```
   For more traditional OOP:
```go
	type Node struct {
	   value string
	   edges []*Edge
	}
	func (p *Node) AddEdge(e *Edge) {...}
```
   All these methods are by default (or by itself) early-bound and statically dispatched (not virtual). They are only dynamically dispatched when invoked thru interfaces (more on this later).

### **2. Embedding: for code reuse and delegation.** ###

   In traditional OOP, one purpose of inheritance is for code reuse: subclasses inherit(or embed in layout) properties and methods of superclass. And inheritance set up "is-a" relation among two types: subclass can be used in anywhere superclass is expected.
   
   In Go, a "outer" struct type can [embed](https://golang.org/ref/spec#Struct_types) another "inner" type to reuse inner's code and logic :
```go
	type OuterType struct {
		InnerType1
		*InnerType2
		...
	}
```
   InnerTypes' fields and methods are _**promoted**_ (accessible by selector) at OuterType; however it is more delegation (has-a relation) than subtyping: OuterType is not a subtype of InnerType, they are independent types:
   - OuterType cannot be used where InnerType is expected.
   - OuterType doesnot contain(embed) InnerTypes' properties directly; when constructing OuterType, embedded InnerType has to be constructed explicitly.
   - Although InnerTypes' method are promoted and can be invoked at OuterType, its target is still InnerType.
   - So we cannot build type hierarchy in Go thru embedding as in Java thru inheritance.

   Shadowing: If OuterType defines a method with same signature as InnerType, this method at OutType will hide its counterpart of InnerType at invocation.

### **3. Interface: for polymorphism.** ###

   In traditional OOP, runtime polymorphism is achieved thru **virtual method table (VMT) and overrides**. Superclass can define set of virtual methods for abstraction while subclass can override virtual methods for extension and variation. VMTs is inherently bound with classes. In Java, all methods are virtual and all classes has its VMT. 
   
   In Go, interfaces play the role (contains) the virtual method table [[Ian Lance Taylor blog](https://www.airs.com/blog/archives/277)]. If you invoke a object's method directly on itself, it is statically dispatched. If you assign an object to an interface value and invoke methods thru the interface, they are dynamically dispatched. 

   However Go interfaces are independent entities separate from structs or others (class-like entities) with methods. All Go methods are early bound and statically dispatched by default. Interface allows **consumers** to specify what polymorphic behaviors it is expecting. Totally unrelated components can satisfy/provide the same interface independently and implicitly (no need for "implements"). While in Java, all classes which provide/implement Java interface or VMT (ie. all interface **providers**) must be in the same class tree as interface.
   
   Interfaces can embed other interfaces; this interface embedding setup "is-a" relation among OuterInterface and InnerInterface: OuterInterface can be used where InnerInterface is required. So we can build hierarchy of abstractions only with interfaces, without implementation details.
 
### **4. How to use these primitives for traditional inheritance based OOP:** ###

   In traditional OOP (Java), classes integrate the above 3 OOP primitives into a unseparatable whole: methods, inheritance/embedding, virtual method table and overrides. This integration results in some class hierarchy based [design patterns](https://en.wikipedia.org/wiki/Design_Patterns) whose advatanges and disadvantages are broadly known.
   
   Go is flatly against these designs based on class hierarch compositions. Go's disintegration of these OOP primitives also guard against these kind of designs. That make people think/complain Go is not a OOP language.
   
   **_Warning: the following are not encouraged practice, just for experiementation._**
   
   By combining these OOP primitives (matching their counterparts in Java class), we can achieve some traditional OOP designs with simple rules: 
   
   - every (class like) entity with methods which **provides** polymorphic behaviours should define these "virtual" methods in an related "base" interface (corresponding to virtual method table in Java):
```go
	//classic OOP example: class Shape with subclasses: Circle, Box,...etc.
	type Shape interface {
		draw()
	}   
```
   - define a "abstract" base struct which embed the above "base" interface: common OO language(such as Java) use single-dispatch: methods are dynamically dispacthed based on virtual method table of the 1st (hidden) "self"/"this" argument. To achieve this in Go, define a "abstract" base struct which embed the above "base" interface. Since the default value of interface is nil, the methods in this base struct are "abstract":
```go
	type ShapeAbstract struct {
		Shape
	}
```
   - use embedding for inheritance and extension: embed the above interface or abstract base struct in outer structs.
```go
	type Circle struct {
		*ShapeAbstract
	}
```
   - overriding involves two steps:
     * method override: in outer struct, define methods with same signature as methods in "base" interface to override them:
```go
	func (c *Circle) draw() {
		fmt.Print("Circle")
	}
```
     * embedded "base" interface override: set the embedded "base" interface (Shape) with a instance of outer struct, so the embedded "base" interface will contain latest overriding methods. This is normally done in constructor of outer struct:
```go
	func NewCircle() *Cirlce {
		rc := &Circle{&ShapeAbstract{}}
		rc.Shape = rc
		return rc
	}
```

   Let's implement the ["template methods"](https://en.wikipedia.org/wiki/Template_method_pattern) design pattern using Go.
   
   In the following Java class Shape, we have three (virtual) methods "drawBoundary(), fillColor()" for extension in subclasses, define reused logic in draw():
```java
	class Shape {
	   //extension point
	   void drawBoundary() { 
           	//no-op
           	out.print("draw nothing");
       	   }
           //extension point
           void fillColor() { 
           	//no-op
           	out.print("fill nothing");
           }
       	   //logic reused in subclasses
           void draw() {
	      drawBoundary();
	      fillColor();
	   }
	}
```
   In Go, define these three virtual methods in a "base" interface and define a "abstract" struct to embed this "base" interface. And we can define reused logic in draw() method with this "abstract" struct following "template methods" design pattern.
```go
	//interface to replace virtual method table in related Java class
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
		//invoke "abstract" methods (defined in interface)
		sa.drawBoundary()
		fmt.Print("-")
		sa.fillColor()
	}
```
   Then define a base struct to extend/embed this "abstract" struct and define placeholder methods. Please note the "constructor pattern" which overrides embedded "Shape" interface value with itself - newly created object.
```go
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
```
   In Java, we can extends the above Shape class with variance:
```java
	class RedRectangle extends Shape {
	   void drawBoundary() {
	   	out.print("Rectangle");
	   }
	   void fillColor() {
		out.print("Red");
	   }
	}
	//create array of shapes and call draw() method on each
	Shape[] shapes = {new Shape(),new RedRectangle()};
	for(Shape s: shapes) { s.draw(); }
```
   In Go, use embedding for inheritance and please note the "constructor pattern" which overrides the embedded Shape interface value with itself - newly created object.
```go
        //embed base struct for inheritance
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
	func(rr RedRectangle) drawBoundary() {
	     fmt.Print("Rectangle")
	}
	//override base method
	func(rr RedRectangle) fillColor() {
	     fmt.Print("Red")
	}
	//create array of shapes and call draw() method on each
	shapes := []Shape{NewRedRectangle(),NewCircle(),...}
	for _,s := range shapes { s.draw() }
```
   Finally, all methods in Java are virtual, so we can override draw() itself for extended behaviour:
```java
	class BlueCircleWithText extends Circle {
	   void fillColor() {
		out.print("Blue");
	   }
	   //override draw() to add text annotation
	   void draw() {
		  //extend superclass's draw()
		  super.draw();
		  //add text
		  out.print("-TextAnnotation");
	   }
	}
```
   In Go, since an outer struct can embed multiple inner types, it is in fact multiple inheritance. So when override and extend draw() method, we have to name the "super" or InnerType explicitly to invoke its draw() method.
```go
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
		bct.Circle.draw()
		//extend with text annotation
		fmt.Print("-TextAnnotation")
	}
```
   Java code creates a 3 level type hierarchy: BlueCircleWithText <= Circle <= Shape, where BlueCircleWithText is subclass of Circle which is subclass of Shape.

   Go code creates a 3 parts delegation chain: BlueCircleWithText -> Circle -> ShapeBase, where all 3 are indepedent types and they all satisfy the Shape interface.

   Again, although we can simulate traditional OOP by combining Go's OOP primitives, it is not encouraged practice.

   [Java](https://github.com/yglcode/golang-oop-primitives/blob/main/TemplateMethods.java) and [Go code](https://github.com/yglcode/golang-oop-primitives/blob/main/go-oop-template-method.go) can be found at [https://github.com/yglcode/golang-oop-primitives](https://github.com/yglcode/golang-oop-primitives).


