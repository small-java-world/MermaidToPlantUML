package parser

import (
	"strings"
	"testing"
)

func TestMermaidParser_ParseToPlantUML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "空の入力",
			input: "",
			want:  "@startuml\n@enduml",
		},
		{
			name:    "不正な構文",
			input:   "invalid syntax",
			wantErr: true,
		},
		{
			name: "単純なクラス定義",
			input: `classDiagram
class User {
    +name: String
    +age: Integer
    +getName()
}`,
			want: "@startuml\nclass User {\n    +name: String\n    +age: Integer\n    +getName()\n}\n@enduml",
		},
		{
			name: "列挙型の定義",
			input: `classDiagram
class Status {
    <<enumeration>>
    PENDING
    ACTIVE
    COMPLETED
}`,
			want: "@startuml\nclass Status {\n    <<enumeration>>\n    PENDING\n    ACTIVE\n    COMPLETED\n}\n@enduml",
		},
		{
			name: "クラス間の関連",
			input: `classDiagram
class Order {
    +id: Integer
}
class OrderItem {
    +quantity: Integer
}
Order "1" *-- "0..*" OrderItem`,
			want: "@startuml\nOrder \"1\" *-- \"0..*\" OrderItem\nclass Order {\n    +id: Integer\n}\nclass OrderItem {\n    +quantity: Integer\n}\n@enduml",
		},
		{
			name: "インターフェースとクラス",
			input: `classDiagram
class IProcessor {
    <<interface>>
    +process()
}
class DataProcessor {
    +processData()
}
DataProcessor ..|> IProcessor`,
			want: "@startuml\nDataProcessor ..|> IProcessor\nclass DataProcessor {\n    +processData()\n}\nclass IProcessor {\n    <<interface>>\n    +process()\n}\n@enduml",
		},
		{
			name: "複雑な関連とジェネリック型",
			input: `classDiagram
class ShoppingCart {
    +items: List~Product~
    +addItem()
}
class Product {
    +name: String
    +price: Double
}
ShoppingCart "1" o-- "*" Product`,
			want: "@startuml\nShoppingCart \"1\" o-- \"*\" Product\nclass Product {\n    +name: String\n    +price: Double\n}\nclass ShoppingCart {\n    +items: List~Product~\n    +addItem()\n}\n@enduml",
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
				// 改行コードを正規化して比較
				got = strings.ReplaceAll(got, "\r\n", "\n")
				want := strings.ReplaceAll(tt.want, "\r\n", "\n")
				if got != want {
					t.Errorf("ParseToPlantUML() got = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestMermaidParser_DebugPrint(t *testing.T) {
	tests := []struct {
		name      string
		format    string
		args      []interface{}
		debugMode bool
	}{
		{
			name:      "デバッグ有効",
			format:    "Test message: %s",
			args:      []interface{}{"debug"},
			debugMode: true,
		},
		{
			name:      "デバッグ無効",
			format:    "Test message: %s",
			args:      []interface{}{"debug"},
			debugMode: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewMermaidParser()
			p.debugEnabled = tt.debugMode
			// デバッグ出力の検証は実装依存のため、エラーが発生しないことだけを確認
			p.debugPrint(tt.format, tt.args...)
		})
	}
}
