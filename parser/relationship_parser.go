package parser

import (
	"regexp"
	"strings"
)

// RelationshipParser は関連の解析を担当します
type RelationshipParser struct {
	relationPattern      *regexp.Regexp
	multiplicityPattern  *regexp.Regexp
	bidirectionalPattern *regexp.Regexp
}

// NewRelationshipParser は新しいRelationshipParserインスタンスを作成します
func NewRelationshipParser() *RelationshipParser {
	return &RelationshipParser{
		relationPattern:      regexp.MustCompile(`(\w+)\s*([<\-\|>o.*]+)\s*(\w+)`),
		multiplicityPattern:  regexp.MustCompile(`(\w+)\s*"([^"]+)"\s*([<\-\|>*o]+)\s*"([^"]+)"\s*(\w+)`),
		bidirectionalPattern: regexp.MustCompile(`(\w+)\s*"([^"]+)"\s*(<-->)\s*"([^"]+)"\s*(\w+)`),
	}
}

// ParseRelationship は関連定義を解析します
func (p *RelationshipParser) ParseRelationship(line string) *Relationship {
	if matches := p.multiplicityPattern.FindStringSubmatch(line); matches != nil {
		return &Relationship{
			Source:     matches[1],
			Target:     matches[5],
			Type:       matches[3],
			SourceMult: matches[2],
			TargetMult: matches[4],
		}
	}

	if matches := p.bidirectionalPattern.FindStringSubmatch(line); matches != nil {
		return &Relationship{
			Source:     matches[1],
			Target:     matches[5],
			Type:       "<-->",
			SourceMult: matches[2],
			TargetMult: matches[4],
		}
	}

	if matches := p.relationPattern.FindStringSubmatch(line); matches != nil {
		return &Relationship{
			Source: matches[1],
			Target: matches[3],
			Type:   matches[2],
		}
	}

	return nil
}

// ExtractClassNames は関連定義から関係するクラス名を抽出します
func (p *RelationshipParser) ExtractClassNames(line string) []string {
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
