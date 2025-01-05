package parser

import (
	"reflect"
	"testing"
)

func TestClassParser_ParseClassContent(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		startLine int
		want      *ClassDefinition
		wantIndex int
		wantErr   bool
	}{
		{
			name: "通常のクラス定義",
			lines: []string{
				"+name: String",
				"+age: Integer",
				"+getName()",
				"}",
			},
			startLine: 0,
			want: &ClassDefinition{
				Members: []string{
					"+name: String",
					"+age: Integer",
					"+getName()",
				},
				IsEnum: false,
			},
			wantIndex: 3,
			wantErr:   false,
		},
		{
			name: "列挙型の定義",
			lines: []string{
				"<<enumeration>>",
				"PENDING",
				"ACTIVE",
				"COMPLETED",
				"}",
			},
			startLine: 0,
			want: &ClassDefinition{
				Members: []string{
					"<<enumeration>>",
					"PENDING",
					"ACTIVE",
					"COMPLETED",
				},
				IsEnum: true,
			},
			wantIndex: 4,
			wantErr:   false,
		},
		{
			name: "インターフェース定義",
			lines: []string{
				"<<interface>>",
				"+process()",
				"+cancel()",
				"}",
			},
			startLine: 0,
			want: &ClassDefinition{
				Members: []string{
					"<<interface>>",
					"+process()",
					"+cancel()",
				},
				IsEnum: false,
			},
			wantIndex: 3,
			wantErr:   false,
		},
		{
			name: "可視性修飾子の組み合わせ",
			lines: []string{
				"+public: String",
				"-private: Integer",
				"#protected: Double",
				"~package: Boolean",
				"}",
			},
			startLine: 0,
			want: &ClassDefinition{
				Members: []string{
					"+public: String",
					"-private: Integer",
					"#protected: Double",
					"~package: Boolean",
				},
				IsEnum: false,
			},
			wantIndex: 4,
			wantErr:   false,
		},
		{
			name: "パラメータ付きメソッド",
			lines: []string{
				"+calculate(Double x, Double y)",
				"-process(String data)",
				"}",
			},
			startLine: 0,
			want: &ClassDefinition{
				Members: []string{
					"+calculate(Double x, Double y)",
					"-process(String data)",
				},
				IsEnum: false,
			},
			wantIndex: 2,
			wantErr:   false,
		},
		{
			name: "ジェネリック型",
			lines: []string{
				"+items: List~String~",
				"+counts: Map~String,Integer~",
				"}",
			},
			startLine: 0,
			want: &ClassDefinition{
				Members: []string{
					"+items: List~String~",
					"+counts: Map~String,Integer~",
				},
				IsEnum: false,
			},
			wantIndex: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewClassParser()
			gotIndex, got, err := p.ParseClassContent(tt.lines, tt.startLine)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseClassContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotIndex != tt.wantIndex {
				t.Errorf("ParseClassContent() gotIndex = %v, want %v", gotIndex, tt.wantIndex)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseClassContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClassParser_ParseMember(t *testing.T) {
	tests := []struct {
		name    string
		matches []string
		want    *ClassMember
	}{
		{
			name:    "属性（デフォルト可視性）",
			matches: []string{"name: String", "", "name", "String", "", ""},
			want: &ClassMember{
				Visibility: "+",
				Name:       "name",
				Type:       "String",
				IsMethod:   false,
			},
		},
		{
			name:    "属性（private）",
			matches: []string{"-count: Integer", "-", "count", "Integer", "", ""},
			want: &ClassMember{
				Visibility: "-",
				Name:       "count",
				Type:       "Integer",
				IsMethod:   false,
			},
		},
		{
			name:    "メソッド（パラメータなし）",
			matches: []string{"+process()", "+", "", "", "process", ""},
			want: &ClassMember{
				Visibility: "+",
				Name:       "process",
				Parameters: "",
				IsMethod:   true,
			},
		},
		{
			name:    "メソッド（パラメータあり）",
			matches: []string{"+calculate(Double x, Double y)", "+", "", "", "calculate", "Double x, Double y"},
			want: &ClassMember{
				Visibility: "+",
				Name:       "calculate",
				Parameters: "Double x, Double y",
				IsMethod:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewClassParser()
			got := p.parseMember(tt.matches)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseMember() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClassParser_FormatMember(t *testing.T) {
	tests := []struct {
		name   string
		member *ClassMember
		want   string
	}{
		{
			name: "属性",
			member: &ClassMember{
				Visibility: "+",
				Name:       "name",
				Type:       "String",
				IsMethod:   false,
			},
			want: "+name: String",
		},
		{
			name: "メソッド（パラメータなし）",
			member: &ClassMember{
				Visibility: "-",
				Name:       "process",
				IsMethod:   true,
			},
			want: "-process()",
		},
		{
			name: "メソッド（パラメータあり）",
			member: &ClassMember{
				Visibility: "#",
				Name:       "calculate",
				Parameters: "Double x, Double y",
				IsMethod:   true,
			},
			want: "#calculate(Double x, Double y)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewClassParser()
			got := p.formatMember(tt.member)

			if got != tt.want {
				t.Errorf("formatMember() got = %v, want %v", got, tt.want)
			}
		})
	}
}
