## Golang OOP primitives ##

### _and combine them to do traditional OOP thru "self referential interface"_ ###
   

Golang is not a traditional object oriented programming language. Instead, it distilled a few OO programming primitives and allow you to compose them to achieve different OO designs.

### **1. Methods (or method-set): for "abstract data types"** ###

   In traditional OOP, methods are inherently bound with class and objects.
   In Go, methods can be defined for any "named" types. Instead of everything is an object as in some OO language, everything (almost) can be attached methods.

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
   
   In Go, interfaces play the role (contains) the virtual method table [[Ian Lance Taylor blog](https://www.airs.com/blog/archives/277)]. If you invoke a object's method directly on itself, it is statically dispatched. If you assign an object to an interface value and invoke methods thru the interface, they are dynamically dispatched. However interfaces are separate entites independent from structs or others (class-like entities) with methods.
   
   Interfaces can embed other interfaces; this interface embedding setup "is-a" relation among OuterInterface and InnerInterface: OuterInterface can be used where InnerInterface is required. So we can build hierarchy of abstractions only with interfaces, without implementation details.
 
### **4. How to use these primitives for traditional inheritance based OOP:** ###

   In traditional OOP (Java), classes integrate the above 3 OOP primitives into a unseparatable whole: methods, inheritance/embedding, virtual method table and overrides. This integration results in some advanced [design patterns](https://en.wikipedia.org/wiki/Design_Patterns) whose advatanges and disadvantages are broadly known.
   
   Go is flatly against these designs based on class hierarch compositions. Go's disintegration of these OOP primitives also guard against these kind of designs. That make people think/complain Go is not a OOP language.
   
   **_Warning: the following are discouraged practice, just for experiementation._**
   
   By combining these OOP primitives (matching their counterparts in Java class), we can achieve some traditional OOP designs with simple rules: 
   
   - every (class like) entity with methods which wants polymorphic behaviour(virtual) should define its "virtual" methods in an related interface.
   - every method which consume polymorphic behaviours should accept this interface as argument (probably the 1st argument, since many OO language(such as Java) is single-dispatch: dynamically dispacthed based on virtual method table of the 1st (hidden) "self"/"this" argument). In Go all methods are early-bound (not virtual), so we have to pass it as interface value in arguments to achieve dynamic dispactching.
   - pass objects thru interface (maybe as 1st argument) and call methods thru interface to achieve virtual/dynamic dispatching.
   - use embedding to simulate inheritance, and shadowing for overriding.

   Let's implement the ["template methods"](https://en.wikipedia.org/wiki/Template_method_pattern) design pattern using Go.
   
   In the following Java class Shape, we have two (virtual) methods "drawBoundary(), fillColor()" for extension in subclasses, define reused logic in draw():
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
   In Go, define these two virtual methods in a separate interface and make draw() in related base struct take interface as 1st argument:
```go
	type Shape interface {
	   drawBoundary()
	   fillColor()
	}
	//ShapeBase is the base struct to be embedded(inherited)
	type ShapeBase struct{}
	//to be overriden
	func(sb *ShapeBase) drawBoundary() {
	   fmt.Print("draw nothing")
	}
	//to be overriden
	func(sb *ShapeBase) fillColor() {
	   fmt.Print("fill nothing")
	}
	//draw() consume polymorphic methods, so make it 
	//take Shape interface as 1st argument
	//in Java, ShapeBase itself will be polymorphic (with VMT),
	//in Go, all methods are early-bound(not virtual). 
	//So we have to pass in interface as argument separately
	func(_ *ShapeBase) draw(sb Shape) {
	   //call methods thru interface for
	   //polymorphism and dynamic dispatching
	   sb.drawBoundary()
	   sb.fillColor()
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
   In Go, use embedding for inheritance and define methods in OuterType to override/shadow methods in InnerType:
```go
    //embed base struct for inheritance
	type RedRectangle struct {
	     *ShapeBase
	}
	func NewRedRectangle() *RedRectangle {
	     return &RedRectangle{&ShapeBase{}}
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
	//note here: we have to pass in Shape instance "s" 
	//to achieve(consume) polymorphic behaviour
	for _,s := range shapes { s.draw(s) }
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
   In Go, add draw() to interface to make draw() overridable; thus we have **_self referential interface_** whose methods refer to the interface itself; that often exist because methods which consume polymorphic behaviour need to be polymorphic (overridable/virtual) itself.
```go
	type Shape interface {
	   drawBoundary()
	   fillColor()
	   draw(Shape)
	}
	type BlueCircleWithText struct {
	     *Circle
	}
	//override
	func(bct *BlueCircleWithText) fillColor() {
	     fmt.Print("Blue")
	}
	//override draw() to add text annotation
	func(bct *BlueCircleWithText) draw(sb Shape) {
	     //extend superclass's draw()
	     //since we can embed multiple InnerTypes, 
	     //it is in fact multiple-dispatch: have to name super explicitly
	     bct.Circle.draw(s)
	     //add text
	     fmt.Print("-TextAnnotation")
	}
```
   Again, although we can simulate traditional OOP by combining Go's OOP primitives, it is not encouraged practice.
   [Java](https://github.com/yglcode/golang-oop-primitives/blob/main/TemplateMethods.java) and [Go code](https://github.com/yglcode/golang-oop-primitives/blob/main/go-oop-template-method.go) can be found at [https://github.com/yglcode/golang-oop-primitives](https://github.com/yglcode/golang-oop-primitives).


