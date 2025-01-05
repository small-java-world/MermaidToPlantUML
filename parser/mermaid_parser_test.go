package parser

import (
	"strings"
	"testing"
)

func TestMermaidParser_ParseToPlantUML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name: "シンプルなクラス定義",
			input: `classDiagram
    class Order {
        +String orderId
        +Date orderDate
        +placeOrder()
    }`,
			expected: `@startuml
class Order {
    +orderId: String
    +orderDate: Date
    +placeOrder()
}
@enduml`,
			wantErr: false,
		},
		{
			name: "クラス間の関連",
			input: `classDiagram
    Order *-- Customer
    class Order {
        +String orderId
    }
    class Customer {
        +String name
    }`,
			expected: `@startuml
Order *-- Customer
class Customer {
    +name: String
}
class Order {
    +orderId: String
}
@enduml`,
			wantErr: false,
		},
		{
			name: "継承関係",
			input: `classDiagram
    Animal <|-- Dog
    Animal <|-- Cat
    class Animal {
        +String name
        +makeSound()
    }`,
			expected: `@startuml
Animal <|-- Dog
Animal <|-- Cat
class Animal {
    +name: String
    +makeSound()
}
@enduml`,
			wantErr: false,
		},
		{
			name: "enumeration定義",
			input: `classDiagram
    class Status {
        <<enumeration>>
        ACTIVE
        INACTIVE
        SUSPENDED
    }`,
			expected: `@startuml
class Status {
    <<enumeration>>
    ACTIVE
    INACTIVE
    SUSPENDED
}
@enduml`,
			wantErr: false,
		},
		{
			name: "enumeration関連",
			input: `classDiagram
    class Order {
        +String id
        +Status status
    }
    class Status {
        <<enumeration>>
        PENDING
        COMPLETED
    }
    Order .. Status`,
			expected: `@startuml
Order .. Status
class Order {
    +id: String
    +status: Status
}
class Status {
    <<enumeration>>
    PENDING
    COMPLETED
}
@enduml`,
			wantErr: false,
		},
		{
			name:  "空の入力",
			input: "",
			expected: `@startuml
@enduml`,
			wantErr: false,
		},
		{
			name: "不正なクラス定義",
			input: `classDiagram
    class Order {
        invalid syntax here
    }`,
			wantErr: true,
		},
		{
			name: "複数のenumeration",
			input: `classDiagram
    class OrderStatus {
        <<enumeration>>
        PENDING
        PROCESSING
        COMPLETED
    }
    class PaymentStatus {
        <<enumeration>>
        UNPAID
        PAID
        REFUNDED
    }
    Order .. OrderStatus
    Order .. PaymentStatus`,
			expected: `@startuml
Order .. OrderStatus
Order .. PaymentStatus
class OrderStatus {
    <<enumeration>>
    PENDING
    PROCESSING
    COMPLETED
}
class PaymentStatus {
    <<enumeration>>
    UNPAID
    PAID
    REFUNDED
}
@enduml`,
			wantErr: false,
		},
		{
			name: "インターフェース定義",
			input: `classDiagram
    class IPaymentProcessor {
        <<interface>>
        +process()
        +refund()
    }
    class StripeProcessor {
        -String apiKey
        +process()
        +refund()
    }
    IPaymentProcessor <|.. StripeProcessor`,
			expected: `@startuml
IPaymentProcessor <|.. StripeProcessor
class IPaymentProcessor {
    <<interface>>
    +process()
    +refund()
}
class StripeProcessor {
    -apiKey: String
    +process()
    +refund()
}
@enduml`,
			wantErr: false,
		},
		{
			name: "複雑な関連",
			input: `classDiagram
    class Order {
        +String id
        +Customer customer
        +List~Product~ products
    }
    class Customer {
        +String name
        +Address address
    }
    class Product {
        +String sku
        +Double price
    }
    class Address {
        +String street
        +String city
    }
    Order "1" *-- "1" Customer
    Order "1" o-- "many" Product
    Customer "1" -- "1" Address`,
			expected: `@startuml
Order "1" *-- "1" Customer
Order "1" o-- "many" Product
Customer "1" -- "1" Address
class Address {
    +street: String
    +city: String
}
class Customer {
    +name: String
    +address: Address
}
class Order {
    +id: String
    +customer: Customer
    +products: List~Product~
}
class Product {
    +sku: String
    +price: Double
}
@enduml`,
			wantErr: false,
		},
		{
			name: "可視性修飾子のバリエーション",
			input: `classDiagram
    class Example {
        +String publicField
        -String privateField
        #String protectedField
        ~String packageField
        +publicMethod()
        -privateMethod()
        #protectedMethod()
        ~packageMethod()
    }`,
			expected: `@startuml
class Example {
    +publicField: String
    -privateField: String
    #protectedField: String
    ~packageField: String
    +publicMethod()
    -privateMethod()
    #protectedMethod()
    ~packageMethod()
}
@enduml`,
			wantErr: false,
		},
		{
			name: "抽象クラス",
			input: `classDiagram
    class Shape {
        <<abstract>>
        +Double area
        +calculate()
    }
    class Circle {
        +Double radius
        +calculate()
    }
    class Rectangle {
        +Double width
        +Double height
        +calculate()
    }
    Shape <|-- Circle
    Shape <|-- Rectangle`,
			expected: `@startuml
Shape <|-- Circle
Shape <|-- Rectangle
class Circle {
    +radius: Double
    +calculate()
}
class Rectangle {
    +width: Double
    +height: Double
    +calculate()
}
class Shape {
    <<abstract>>
    +area: Double
    +calculate()
}
@enduml`,
			wantErr: false,
		},
		{
			name: "パラメータ付きメソッド",
			input: `classDiagram
    class Calculator {
        +add(Double a, Double b)
        +subtract(Double x, Double y)
        +multiply(Double p, Double q)
        +divide(Double m, Double n)
    }`,
			expected: `@startuml
class Calculator {
    +add(Double a, Double b)
    +subtract(Double x, Double y)
    +multiply(Double p, Double q)
    +divide(Double m, Double n)
}
@enduml`,
			wantErr: false,
		},
		{
			name: "複数の継承関係",
			input: `classDiagram
    class Vehicle {
        +String model
        +start()
        +stop()
    }
    class Engine {
        +Integer power
        +start()
        +stop()
    }
    class ElectricCar {
        +Integer batteryLevel
        +charge()
    }
    Vehicle <|-- ElectricCar
    Engine <|-- ElectricCar`,
			expected: `@startuml
Vehicle <|-- ElectricCar
Engine <|-- ElectricCar
class ElectricCar {
    +batteryLevel: Integer
    +charge()
}
class Engine {
    +power: Integer
    +start()
    +stop()
}
class Vehicle {
    +model: String
    +start()
    +stop()
}
@enduml`,
			wantErr: false,
		},
		{
			name: "双方向の関連",
			input: `classDiagram
    class Student {
        +String name
        +List~Course~ courses
    }
    class Course {
        +String code
        +List~Student~ students
    }
    Student "many" <--> "many" Course`,
			expected: `@startuml
Student "many" <--> "many" Course
class Course {
    +code: String
    +students: List~Student~
}
class Student {
    +name: String
    +courses: List~Course~
}
@enduml`,
			wantErr: false,
		},
		{
			name: "関連クラス",
			input: `classDiagram
    class Order {
        +String orderId
    }
    class Product {
        +String sku
    }
    class OrderItem {
        +Integer quantity
        +Double price
    }
    Order "1" -- "many" Product
    OrderItem .. Order
    OrderItem .. Product`,
			expected: `@startuml
Order "1" -- "many" Product
OrderItem .. Order
OrderItem .. Product
class Order {
    +orderId: String
}
class OrderItem {
    +quantity: Integer
    +price: Double
}
class Product {
    +sku: String
}
@enduml`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewMermaidParser()
			got, err := p.ParseToPlantUML(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToPlantUML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 改行コードを統一して比較
				got = strings.ReplaceAll(got, "\r\n", "\n")
				expected := strings.ReplaceAll(tt.expected, "\r\n", "\n")

				if got != expected {
					t.Errorf("ParseToPlantUML() = %v, want %v", got, expected)
				}
			}
		})
	}
}
