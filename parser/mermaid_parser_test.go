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
        +void placeOrder()
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
        +void makeSound()
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
