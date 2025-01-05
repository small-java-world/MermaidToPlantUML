package parser

// ClassDefinition はクラス定義の内容を表現します
type ClassDefinition struct {
	Members []string
	IsEnum  bool
}

// Relationship はクラス間の関連を表現します
type Relationship struct {
	Source     string
	Target     string
	Type       string
	SourceMult string
	TargetMult string
}

// ClassMember はクラスのメンバー（属性やメソッド）を表現します
type ClassMember struct {
	Visibility string
	Name       string
	Type       string
	Parameters string
	IsMethod   bool
}
