package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// ClassParser はクラス定義の解析を担当します
type ClassParser struct {
	memberPattern *regexp.Regexp
}

// NewClassParser は新しいClassParserインスタンスを作成します
func NewClassParser() *ClassParser {
	return &ClassParser{
		memberPattern: regexp.MustCompile(`\s*([+\-#~])?(?:(\w+):\s*(\w+(?:~[^~]+~)?)|(\w+)(?:\((.*?)\))?)`)}
}

// ParseClassContent はクラス定義の内容を解析します
func (p *ClassParser) ParseClassContent(lines []string, startIndex int) (int, *ClassDefinition, error) {
	classDef := &ClassDefinition{
		Members: []string{},
		IsEnum:  false,
	}

	currentIndex := startIndex

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
		} else if matches := p.memberPattern.FindStringSubmatch(line); matches != nil {
			member := p.parseMember(matches)
			if member != nil {
				classDef.Members = append(classDef.Members, p.formatMember(member))
			}
		}
		currentIndex++
	}

	return currentIndex, classDef, nil
}

// parseMember はメンバー定義を解析します
func (p *ClassParser) parseMember(matches []string) *ClassMember {
	visibility := matches[1]
	if visibility == "" {
		visibility = "+"
	}

	if matches[2] != "" && matches[3] != "" {
		// 属性の場合
		return &ClassMember{
			Visibility: visibility,
			Name:       matches[2],
			Type:       matches[3],
			IsMethod:   false,
		}
	} else if matches[4] != "" {
		// メソッドの場合
		return &ClassMember{
			Visibility: visibility,
			Name:       matches[4],
			Parameters: matches[5],
			IsMethod:   true,
		}
	}

	return nil
}

// formatMember はメンバーを文字列形式にフォーマットします
func (p *ClassParser) formatMember(member *ClassMember) string {
	if member.IsMethod {
		if member.Parameters != "" {
			return fmt.Sprintf("%s%s(%s)", member.Visibility, member.Name, member.Parameters)
		}
		return fmt.Sprintf("%s%s()", member.Visibility, member.Name)
	}
	return fmt.Sprintf("%s%s: %s", member.Visibility, member.Name, member.Type)
}
