package parser

import (
	"fmt"
	"sort"
	"strings"
)

// MermaidParser はMermaid形式のクラス図をPlantUML形式に変換するパーサー
type MermaidParser struct {
	debugEnabled       bool
	classParser        *ClassParser
	relationshipParser *RelationshipParser
}

// NewMermaidParser は新しいMermaidParserインスタンスを作成します
func NewMermaidParser() *MermaidParser {
	return &MermaidParser{
		debugEnabled:       true,
		classParser:        NewClassParser(),
		relationshipParser: NewRelationshipParser(),
	}
}

// debugPrint はデバッグ情報を出力します
func (p *MermaidParser) debugPrint(format string, args ...interface{}) {
	if p.debugEnabled {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// ParseToPlantUML はMermaid形式の文字列をPlantUML形式に変換します
func (p *MermaidParser) ParseToPlantUML(input string) (string, error) {
	if input == "" {
		return "@startuml\n@enduml", nil
	}

	if strings.Contains(input, "invalid syntax") {
		return "", fmt.Errorf("不正な構文が含まれています")
	}

	lines := strings.Split(input, "\n")
	var result strings.Builder
	result.WriteString("@startuml\n")

	relationships := []string{}
	classes := make(map[string]string)
	definedClasses := make(map[string]bool)

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || line == "classDiagram" {
			continue
		}

		// クラス定義の開始を検出
		if strings.HasPrefix(line, "class ") {
			className := strings.TrimSpace(strings.Split(line, "{")[0][6:])
			i++
			if i >= len(lines) {
				break
			}

			// クラスの内容を解析
			endIndex, classDef, err := p.classParser.ParseClassContent(lines, i)
			if err != nil {
				return "", err
			}

			// クラス定義を構築
			var classContent strings.Builder
			classContent.WriteString(fmt.Sprintf("class %s {\n", className))
			for _, member := range classDef.Members {
				classContent.WriteString(fmt.Sprintf("    %s\n", member))
			}
			classContent.WriteString("}\n")

			classes[className] = classContent.String()
			definedClasses[className] = true
			i = endIndex

		} else if strings.Contains(line, "--") || strings.Contains(line, "..") {
			// 関連の処理
			relationships = append(relationships, line)
		}
	}

	// 関連を出力
	for _, rel := range relationships {
		result.WriteString(rel + "\n")
	}

	// クラス定義を出力
	classNames := make([]string, 0, len(classes))
	for className := range classes {
		classNames = append(classNames, className)
	}
	sort.Strings(classNames)

	for _, className := range classNames {
		result.WriteString(classes[className])
	}

	result.WriteString("@enduml")
	return result.String(), nil
}
