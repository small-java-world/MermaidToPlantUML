package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"mermaid2plantuml/parser"
	"mermaid2plantuml/plantuml"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// コマンドライン引数の解析
	format := flag.String("format", "png", "出力フォーマット (png|svg|pdf)")
	output := flag.String("o", "", "出力ファイルパス")
	flag.Parse()

	// 入力ファイルの確認
	if flag.NArg() < 1 {
		return fmt.Errorf("使用方法: mmd2img [-format=<png|svg|pdf>] [-o output_file] input.mmd")
	}

	inputFile := flag.Arg(0)
	if filepath.Ext(inputFile) != ".mmd" {
		return fmt.Errorf("入力ファイルは.mmd拡張子である必要があります")
	}

	// 入力ファイルの読み込み
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("入力ファイルの読み込みに失敗: %v", err)
	}

	// Mermaid → PlantUML変換
	p := parser.NewMermaidParser()
	pumlContent, err := p.ParseToPlantUML(string(input))
	if err != nil {
		return fmt.Errorf("Mermaid形式の解析に失敗: %v", err)
	}

	// 出力ファイル名の決定
	var outputPuml string
	if *output != "" {
		// 出力先が指定されている場合
		outputDir := filepath.Dir(*output)
		baseName := filepath.Base(*output)
		outputPuml = filepath.Join(outputDir, baseName[:len(baseName)-len(filepath.Ext(baseName))]+".puml")
	} else {
		// 出力先が指定されていない場合
		inputBase := filepath.Base(inputFile)
		outputPuml = filepath.Join(
			filepath.Dir(inputFile),
			inputBase[:len(inputBase)-len(filepath.Ext(inputBase))]+".puml",
		)
	}

	// PlantUMLファイルの保存
	if err := ioutil.WriteFile(outputPuml, []byte(pumlContent), 0644); err != nil {
		return fmt.Errorf("PlantUMLファイルの保存に失敗: %v", err)
	}

	// PlantUML実行
	executor := plantuml.NewPlantUMLExecutor()
	if err := executor.GenerateImage(outputPuml, *format); err != nil {
		return fmt.Errorf("画像生成に失敗: %v", err)
	}

	fmt.Printf("変換が完了しました:\n")
	fmt.Printf("- PlantUMLファイル: %s\n", outputPuml)
	fmt.Printf("- 画像ファイル: %s\n", outputPuml[:len(outputPuml)-len(".puml")]+"."+*format)

	return nil
}
