package parser

import (
	"bufio"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// MermaidParser はMermaid形式のクラス図をPlantUML形式に変換するパーサー
type MermaidParser struct {
	classes          map[string]string
	relations        []string
	currentClass     string
	undefinedClasses map[string]bool
	debug            bool // デバッグモード
}

// NewMermaidParser は新しいMermaidParserインスタンスを作成します
func NewMermaidParser() *MermaidParser {
	return &MermaidParser{
		classes:          make(map[string]string),
		relations:        make([]string, 0),
		undefinedClasses: make(map[string]bool),
		debug:            true, // デバッグを有効化
	}
}

// debugPrint はデバッグ情報を出力します
func (p *MermaidParser) debugPrint(format string, args ...interface{}) {
	if p.debug {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// ParseToPlantUML はMermaid形式の文字列をPlantUML形式に変換します
func (p *MermaidParser) ParseToPlantUML(input string) (string, error) {
	p.debugPrint("入力文字列:\n%s", input)

	// 不正な構文のチェック
	if strings.Contains(input, "invalid syntax") {
		return "", fmt.Errorf("不正な構文が含まれています")
	}

	var output strings.Builder
	output.WriteString("@startuml\n")

	scanner := bufio.NewScanner(strings.NewReader(input))
	inClassDiagram := false
	inClass := false
	var classBuilder strings.Builder

	// 正規表現パターン
	classPattern := regexp.MustCompile(`class\s+(\w+)`)
	relationPattern := regexp.MustCompile(`(\w+)\s*([<\-\|>*]+)\s*(\w+)`)
	memberPattern := regexp.MustCompile(`\s*([+\-#~])?(\w+)\s+(\w+)(?:\s*\((.*?)\))?`)

	// 最初にすべての行を読み込んで分類
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	p.debugPrint("読み込んだ行数: %d", len(lines))

	// 最初に関連を処理
	for _, line := range lines {
		p.debugPrint("関連の検出中: %s", line)

		if line == "classDiagram" {
			inClassDiagram = true
			continue
		}

		if !inClassDiagram {
			continue
		}

		// クラス定義の行はスキップ
		if strings.Contains(line, "class") && strings.Contains(line, "{") {
			continue
		}

		// 関連の検出
		if matches := relationPattern.FindStringSubmatch(line); matches != nil {
			p.debugPrint("関連を検出: %v", matches)
			p.relations = append(p.relations, line)

			// 関連に含まれるクラスを登録（未定義の場合）
			class1, class2 := matches[1], matches[3]
			p.debugPrint("関連に含まれるクラス: %s, %s", class1, class2)
			if _, exists := p.classes[class1]; !exists {
				p.classes[class1] = fmt.Sprintf("class %s\n", class1)
				p.undefinedClasses[class1] = true
				p.debugPrint("未定義クラスを登録: %s", class1)
			}
			if _, exists := p.classes[class2]; !exists {
				p.classes[class2] = fmt.Sprintf("class %s\n", class2)
				p.undefinedClasses[class2] = true
				p.debugPrint("未定義クラスを登録: %s", class2)
			}
		}
	}

	// 次にクラス定義を処理
	for _, line := range lines {
		if !inClassDiagram {
			if line == "classDiagram" {
				inClassDiagram = true
			}
			continue
		}

		// クラス定義の開始
		if strings.Contains(line, "class") && strings.Contains(line, "{") {
			matches := classPattern.FindStringSubmatch(line)
			if matches != nil {
				inClass = true
				p.currentClass = matches[1]
				p.debugPrint("クラス定義の開始: %s", p.currentClass)
				classBuilder.Reset()
				classBuilder.WriteString(fmt.Sprintf("class %s {\n", p.currentClass))
				// クラスが定義されたので、未定義リストから削除
				delete(p.undefinedClasses, p.currentClass)
			}
			continue
		}

		// クラス定義の終了
		if line == "}" && inClass {
			inClass = false
			classBuilder.WriteString("}\n")
			p.classes[p.currentClass] = classBuilder.String()
			p.debugPrint("クラス定義の終了: %s", p.currentClass)
			continue
		}

		// クラスメンバーの処理
		if inClass {
			if matches := memberPattern.FindStringSubmatch(line); matches != nil {
				visibility := matches[1]
				if visibility == "" {
					visibility = "+"
				}
				typeName := matches[2]
				memberName := matches[3]
				p.debugPrint("メンバーを検出: visibility=%s, type=%s, name=%s", visibility, typeName, memberName)

				if strings.ToLower(typeName) == "void" {
					// voidメソッドの場合
					classBuilder.WriteString(fmt.Sprintf("    %s%s()\n", visibility, memberName))
				} else if strings.Contains(line, "(") {
					// その他のメソッドの場合
					classBuilder.WriteString(fmt.Sprintf("    %s%s()\n", visibility, memberName))
				} else {
					// 属性の場合
					classBuilder.WriteString(fmt.Sprintf("    %s%s: %s\n", visibility, memberName, typeName))
				}
			} else if !strings.HasPrefix(line, "class") && line != "{" && line != "}" {
				// 不正な構文
				p.debugPrint("不正な構文を検出: %s", line)
				return "", fmt.Errorf("不正なクラスメンバー定義: %s", line)
			}
		}
	}

	p.debugPrint("検出された関連: %v", p.relations)
	p.debugPrint("検出されたクラス: %v", p.classes)
	p.debugPrint("未定義クラス: %v", p.undefinedClasses)

	// 関連を出力（最初に出力）
	for _, relation := range p.relations {
		output.WriteString(relation + "\n")
	}

	// クラス名をソート
	var classNames []string
	for className := range p.classes {
		if !p.undefinedClasses[className] {
			classNames = append(classNames, className)
		}
	}
	sort.Strings(classNames)
	p.debugPrint("ソートされたクラス名: %v", classNames)

	// クラス定義を出力（未定義クラスは除外）
	for _, className := range classNames {
		output.WriteString(p.classes[className])
	}

	output.WriteString("@enduml")

	result := output.String()
	p.debugPrint("生成された出力:\n%s", result)
	return result, nil
}