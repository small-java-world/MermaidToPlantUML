package parser

import (
	"reflect"
	"testing"
)

func TestRelationshipParser_ParseRelationship(t *testing.T) {
	tests := []struct {
		name string
		line string
		want *Relationship
	}{
		{
			name: "継承関係",
			line: "Dog --|> Animal",
			want: &Relationship{
				Source: "Dog",
				Target: "Animal",
				Type:   "--|>",
			},
		},
		{
			name: "多重度付き関連",
			line: `Order "1" *-- "0..*" OrderItem`,
			want: &Relationship{
				Source:     "Order",
				Target:     "OrderItem",
				Type:       "*--",
				SourceMult: "1",
				TargetMult: "0..*",
			},
		},
		{
			name: "双方向関連",
			line: `Student "1..*" <--> "1..*" Course`,
			want: &Relationship{
				Source:     "Student",
				Target:     "Course",
				Type:       "<-->",
				SourceMult: "1..*",
				TargetMult: "1..*",
			},
		},
		{
			name: "集約関係",
			line: "Department o-- Employee",
			want: &Relationship{
				Source: "Department",
				Target: "Employee",
				Type:   "o--",
			},
		},
		{
			name: "コンポジション関係",
			line: "Car *-- Engine",
			want: &Relationship{
				Source: "Car",
				Target: "Engine",
				Type:   "*--",
			},
		},
		{
			name: "依存関係",
			line: "Controller ..> Service",
			want: &Relationship{
				Source: "Controller",
				Target: "Service",
				Type:   "..>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewRelationshipParser()
			got := p.ParseRelationship(tt.line)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRelationship() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationshipParser_ExtractClassNames(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []string
	}{
		{
			name: "単純な関連",
			line: "ClassA -- ClassB",
			want: []string{"ClassA", "ClassB"},
		},
		{
			name: "多重度付き関連",
			line: `Order "1" *-- "0..*" OrderItem`,
			want: []string{"Order", "OrderItem"},
		},
		{
			name: "双方向関連",
			line: `Student "1..*" <--> "1..*" Course`,
			want: []string{"Student", "Course"},
		},
		{
			name: "継承関係",
			line: "Dog --|> Animal",
			want: []string{"Dog", "Animal"},
		},
		{
			name: "依存関係",
			line: "Controller ..> Service",
			want: []string{"Controller", "Service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewRelationshipParser()
			got := p.ExtractClassNames(tt.line)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractClassNames() got = %v, want %v", got, tt.want)
			}
		})
	}
}
