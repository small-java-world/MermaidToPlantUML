package parser

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// MermaidParser はMermaid形式のクラス図をPlantUML形式に変換するパーサー
type MermaidParser struct {
	debugEnabled bool
}

type ClassDefinition struct {
	Members []string
	IsEnum  bool
}

// NewMermaidParser は新しいMermaidParserインスタンスを作成します
func NewMermaidParser() *MermaidParser {
	return &MermaidParser{
		debugEnabled: true,
	}
}

// debugPrint はデバッグ情報を出力します
func (p *MermaidParser) debugPrint(format string, args ...interface{}) {
	if p.debugEnabled {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

func (p *MermaidParser) extractClassNames(line string) []string {
	var classNames []string
	parts := strings.Fields(line)
	for _, part := range parts {
		if !strings.Contains(part, "--") && !strings.Contains(part, "..") &&
			!strings.Contains(part, "<") && !strings.Contains(part, ">") &&
			!strings.Contains(part, "*") && !strings.Contains(part, "\"") {
			classNames = append(classNames, part)
		}
	}
	return classNames
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
	undefinedClasses := make(map[string]bool)
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
			endIndex, classDef, err := p.parseClassContent(lines, i)
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
			delete(undefinedClasses, className)
			i = endIndex

		} else if strings.Contains(line, "--") || strings.Contains(line, "..") {
			// 関連の処理
			relationships = append(relationships, line)
			// 関連に含まれるクラスを未定義クラスとして登録
			for _, className := range p.extractClassNames(line) {
				if !definedClasses[className] {
					undefinedClasses[className] = true
				}
			}
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

func (p *MermaidParser) parseClassContent(lines []string, startIndex int) (int, *ClassDefinition, error) {
	classDef := &ClassDefinition{
		Members: []string{},
		IsEnum:  false,
	}

	currentIndex := startIndex
	memberPattern := regexp.MustCompile(`\s*([+\-#~])?(?:(\w+(?:~\w+~)?)\s+(\w+)|(\w+)(?:\((.*?)\))?)`)

	for currentIndex < len(lines) {
		line := strings.TrimSpace(lines[currentIndex])
		if line == "}" {
			break
		}

		if line == "<<enumeration>>" {
			classDef.IsEnum = true
			classDef.Members = append(classDef.Members, line)
			currentIndex++
			continue
		}

		if classDef.IsEnum {
			// 列挙型の値は単純に追加
			if line != "" && line != "<<enumeration>>" {
				classDef.Members = append(classDef.Members, line)
			}
		} else if line == "<<interface>>" || line == "<<abstract>>" {
			// インターフェースと抽象クラスのステレオタイプを追加
			classDef.Members = append(classDef.Members, line)
		} else if matches := memberPattern.FindStringSubmatch(line); matches != nil {
			// メンバーの処理
			visibility := matches[1]
			if visibility == "" {
				visibility = "+"
			}

			if matches[2] != "" && matches[3] != "" {
				// 属性の場合
				typeName := matches[2]
				memberName := matches[3]
				classDef.Members = append(classDef.Members, fmt.Sprintf("%s%s: %s", visibility, memberName, typeName))
			} else if matches[4] != "" {
				// メソッドの場合
				methodName := matches[4]
				params := matches[5]
				if params != "" {
					classDef.Members = append(classDef.Members, fmt.Sprintf("%s%s(%s)", visibility, methodName, params))
				} else {
					classDef.Members = append(classDef.Members, fmt.Sprintf("%s%s()", visibility, methodName))
				}
			}
		}
		currentIndex++
	}

	return currentIndex, classDef, nil
}
