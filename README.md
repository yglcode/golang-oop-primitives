## Golang OOP primitives ##

### _traditional OOP thru "Embedded Interface" and Go prefered composition_ ###
   

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
   All these methods are by default early-bound and statically dispatched (not virtual). They are only dynamically dispatched when invoked thru interfaces (more on this later).

### **2. Embedding: for code reuse and delegation.** ###

   In traditional OOP, one purpose of inheritance is for code reuse: subclasses inherit (or embed a copy in layout) properties and methods of superclass. And inheritance set up "is-a" relation among two types: subclass can be used in anywhere superclass is expected.
   
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

   In traditional OOP, runtime polymorphism is achieved thru **virtual method table (VMT) and overrides**. Superclass can define set of virtual methods for abstraction while subclass can override virtual methods for extension and variation. As the core of class hierarchy based composition, VMTs is inherently bound with classes. In Java, by default methods are virtual and all classes has its VMT. 
   
   In Go, interfaces play the role (contains) the virtual method table [[Russ Cox blog](https://research.swtch.com/interfaces)][[Ian Lance Taylor blog](https://www.airs.com/blog/archives/277)]. If you invoke a object's method directly on itself, it is statically dispatched. If you assign an object to an interface value and invoke methods thru the interface, they are dynamically dispatched. 

   However Go interfaces are independent entities separate from structs or others (class-like entities) with methods. All Go methods are early bound and statically dispatched by default. So Go's interface itself doesn't enable class hierarchy based composition. Instead interface allows **consumers** to specify what polymorphic behaviors it is expecting. Totally unrelated components can satisfy/provide the same interface independently and implicitly (no need for "implements"). While in Java, all classes which provide/implement Java interface or VMT (ie. all interface **providers**) must be in the same class tree as interface.
   
   Interfaces can embed other interfaces; this interface embedding setup "is-a" relation among OuterInterface and InnerInterface: OuterInterface can be used where InnerInterface is required. So we can build hierarchy of abstractions only with interfaces, without implementation details.
 
### **4. How to use these primitives for traditional inheritance based OOP:** ###

   In traditional OOP (Java), classes integrate the above 3 OOP primitives into a inseparatable whole: methods, inheritance/embedding, virtual method table and overrides. This integration results in some class hierarchy based [design patterns](https://en.wikipedia.org/wiki/Design_Patterns) whose advatanges and disadvantages are broadly known.
   
   Go is flatly against these designs based on class hierarchy compositions. Go's disintegration of these OOP primitives also guard against these kind of designs. That make people think/complain Go is not a OOP language.
   
   **_Warning: the following are not encouraged practice, just for experiementation._**
   
   By combining these OOP primitives (matching their counterparts in Java class), we can achieve some traditional OOP designs with simple rules: 
   
   - every (class like) entity with methods which **provides** polymorphic behaviours should define these "virtual" methods in an related "base" interface:
```go
	//classic OOP example: class Shape with subclasses: Circle, Box,...etc.
	type Shape interface {
		draw()
	}   
```
   - define a base struct which embed the above "base" interface (as virtual method table in Java): common OO languages (such as Java) use single-dispatch: methods are dynamically dispatched based on virtual method table of the 1st (hidden) "self"/"this" argument. To achieve this in Go, define a base struct which embed the above "base" interface. Since the default value of interface is nil, the methods in this base struct are "abstract". "Default"/"stub" methods implementations should be provided in base struct or by embedding base struct and overriding/shadowing the methods:
```go
	type ShapeAbstract struct {
		Shape
	}
```
   - use embedding for inheritance and extension: embed "super"/"parent" struct or interface in outer "sub"/"child" structs to extend.
```go
	type Circle struct {
		*ShapeAbstract
	}
```
   - overriding involves two steps:
   - override methods: in outer struct, define methods with same signature as methods in "super"/"parent" inner types to shadow/override them:
```go
	func (c *Circle) draw() {
		fmt.Print("Circle")
	}
```
   - override embedded "base" interface (ie. update VMT): set the embedded "base" interface (Shape) with a instance of outer struct, so the embedded "base" interface will contain latest overriding methods. This can be done in constructor of outer struct: 
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

### **5. Go's typical composition: Simple Control Flow (Readability), Small Interfaces (Separation of Concerns)** ###

One issue of the above "template methods" design pattern is complicate control flow. Invoking a method may involve jumping up and down the inheritance hierarchy multiple times.
    
In above sample Java code, BlueCircleWithText.draw() call path will be:
    
BlueCircleWithText.draw() -> Shape.draw() -> Circle.drawBoundary() -> BlueCircleWithText.fillColor() -> back to Shape.draw() -> BlueCircleWithText.draw() complete.
    
It is not uncommon in OOP frameworks, some calls will go up and down inheritance hierarchy multiple times.
    
Similarly, in the above Go code implementing "template methods" design, the control flow is jumping back and forth in the delegation chain.
    
Embedding is used in many places inside Go standard packages, and control flow only goes in one direction: embedding struct -> embedded struct.
    
Go prefers simple straight-forward control flow which is consistent with the way how human read and understand (readability and maintainability). A prime example of this is how traditional epoll-based networking code is callback based and driven by IO events, which results in network app code flow broken up and jump thru different callback functions. In Go, by using channel and goroutine(coroutines), network app code flow becomes a simple sequential flow from top to bottom, which is easier to understand and maintain.
    
Another Go's design [proverbs](https://go-proverbs.github.io/) is preference for small interfaces. The prime examples are [io.Reader](https://golang.org/pkg/io/#Reader) and [io.Writer](https://golang.org/pkg/io/#Writer) which have one method. Small interfaces encourage separation of concerns and better abstraction. 
    
In Go, interfaces allow *consumer* code specify what polymorphic behaviors it expects. Reexaming above "Shape" interface, we can find it has two consumers, and Shape interface is in fact a mix of two separate method-sets:
    
1st consumer is client code which call/use the hierachy of Shape / Circle / Rectangle /..., which expects something *drawable*:
```go
    type Drawable interface {
        draw()
    }
    //so client code can draw a list of shapes:
    shapes := []Drawable{NewCircle(),NewRectanlge(),...}
    for _,s := range shapes { s.draw() }
```
2nd consumer is internal implementation of "draw()" method which need to be customized by polymorphic "drawBoundary()" and "fillColor()" methods. If we assume this customization is a valid design decision, we could have simpler implementation without embedding and overriding as following:
```go
    type DrawOperations interface {
        drawBoundary()
        fillColor()
    }
    // shared/reused logic
    func commonDraw(ops DrawOperations) {
        ...
        ops.drawBoundary()
        ...
        ops.fillColor()
        ...
    }
    // various shapes can be defined without embedding
    type RedCircle struct {}
    func (rc *RedCircle) drawBoundary() {
        fmt.Print("Circle")
    }
    func (rc *RedCircle) fillColor() {
        fmt.Print("Red")
    }
    func (rc *RedCircle) draw() {
        commonDraw(rc)
        ...other customizations...
    }
```
