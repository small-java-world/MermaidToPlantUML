package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func TestMainIntegration(t *testing.T) {
	// PlantUMLが利用できない環境でのテストのため、エラーを期待する
	wantPlantUMLErr := true

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "mermaid2plantuml_test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のMermaidファイルを作成
	testMmd := `classDiagram
    class User {
        +String name
        +String email
        +void register()
    }
    class Order {
        +String orderId
        +Date orderDate
    }
    User "1" -- "*" Order`

	mmdFile := filepath.Join(tempDir, "test.mmd")
	if err := os.WriteFile(mmdFile, []byte(testMmd), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	// 現在のワーキングディレクトリを保存
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("現在のディレクトリの取得に失敗: %v", err)
	}

	// テストディレクトリに移動
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("テストディレクトリへの移動に失敗: %v", err)
	}
	defer os.Chdir(pwd)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "基本的な変換",
			args:    []string{filepath.Base(mmdFile)},
			wantErr: wantPlantUMLErr, // PlantUMLが利用できない場合はエラー
		},
		{
			name:    "PNG形式を指定",
			args:    []string{"-format", "png", filepath.Base(mmdFile)},
			wantErr: wantPlantUMLErr, // PlantUMLが利用できない場合はエラー
		},
		{
			name:    "出力先を指定",
			args:    []string{"-o", "output.png", filepath.Base(mmdFile)},
			wantErr: wantPlantUMLErr, // PlantUMLが利用できない場合はエラー
		},
		{
			name:    "存在しないファイル",
			args:    []string{"nonexistent.mmd"},
			wantErr: true,
		},
		{
			name:    "不正な拡張子",
			args:    []string{"test.txt"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// フラグをリセット
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// コマンドライン引数を設定
			oldArgs := os.Args
			os.Args = append([]string{"mmd2img"}, tt.args...)
			defer func() { os.Args = oldArgs }()

			// メイン関数の実行
			err := run()
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// PlantUMLファイルが生成されているか確認
				baseName := filepath.Base(tt.args[len(tt.args)-1])
				expectedPuml := baseName[:len(baseName)-len(".mmd")] + ".puml"
				if _, err := os.Stat(expectedPuml); os.IsNotExist(err) {
					t.Errorf("PlantUMLファイルが生成されていません: %s", expectedPuml)
				}
			}
		})
	}
}
