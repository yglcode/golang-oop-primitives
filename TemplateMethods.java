package cmd.test;

import static java.lang.System.*;

class Shape {
    void draw() {
        drawBoundary();
        out.print("-");
        fillColor();
    }
    void drawBoundary() {
        //noop
        out.print("draw nothing");
    }
    void fillColor() {
        //noop
        out.print("fill nothing");
    }
}

class Circle extends Shape {
    void drawBoundary() {
        out.print("Circle");
    }
}

class RedRectangle extends Shape {
    void fillColor() {
        out.print("Red");
    }
    void drawBoundary() {
        out.print("Rectangle");
    }
}

class BlueCircleWithText extends Circle {
    void fillColor() {
        out.print("Blue");
    }
    void draw() {
        super.draw();
        out.print("-TextAnnotation");
    }
}

public class TemplateMethods {
    public static void main(String[] args) {
        Shape[] ss={
            new Shape(),
            new Circle(),
            new RedRectangle(),
            new BlueCircleWithText()
        };
        for(Shape s: ss) {
            s.draw();
            out.println();
        }
    }
}
