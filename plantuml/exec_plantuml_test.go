package plantuml

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// モック用のPlantUMLエグゼキューター
type mockPlantUMLExecutor struct {
	plantumlPath string
	mockOutput   bool
}

func (e *mockPlantUMLExecutor) GenerateImage(pumlFile string, format string) error {
	// 入力ファイルの存在確認
	if _, err := os.Stat(pumlFile); os.IsNotExist(err) {
		return err
	}

	// フォーマットの検証
	switch format {
	case "png", "svg", "pdf":
		// サポートされているフォーマット
	default:
		return fmt.Errorf("サポートされていないフォーマット: %s", format)
	}

	if e.mockOutput {
		// テスト用の空の出力ファイルを作成
		outputFile := pumlFile[:len(pumlFile)-len(".puml")] + "." + format
		return os.WriteFile(outputFile, []byte("mock output"), 0644)
	}
	return nil
}

func TestPlantUMLExecutor_GenerateImage(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "plantuml_test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のPlantUMLファイルを作成
	testPuml := `@startuml
class TestClass {
    +String name
}
@enduml`
	pumlFile := filepath.Join(tempDir, "test.puml")
	if err := os.WriteFile(pumlFile, []byte(testPuml), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	tests := []struct {
		name      string
		pumlFile  string
		format    string
		wantErr   bool
		setupFunc func(*mockPlantUMLExecutor)
	}{
		{
			name:     "PNG形式での出力",
			pumlFile: pumlFile,
			format:   "png",
			wantErr:  false,
			setupFunc: func(e *mockPlantUMLExecutor) {
				e.mockOutput = true
			},
		},
		{
			name:     "SVG形式での出力",
			pumlFile: pumlFile,
			format:   "svg",
			wantErr:  false,
			setupFunc: func(e *mockPlantUMLExecutor) {
				e.mockOutput = true
			},
		},
		{
			name:     "存在しないファイル",
			pumlFile: "nonexistent.puml",
			format:   "png",
			wantErr:  true,
		},
		{
			name:     "不正なフォーマット",
			pumlFile: pumlFile,
			format:   "invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &mockPlantUMLExecutor{
				plantumlPath: "mock-plantuml",
				mockOutput:   false,
			}
			if tt.setupFunc != nil {
				tt.setupFunc(executor)
			}

			err := executor.GenerateImage(tt.pumlFile, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateImage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.format == "png" && executor.mockOutput {
				// 出力ファイルが生成されているか確認
				outFile := tt.pumlFile[:len(tt.pumlFile)-len(".puml")] + ".png"
				if _, err := os.Stat(outFile); os.IsNotExist(err) {
					t.Errorf("出力ファイルが生成されていません: %s", outFile)
				}
			}
		})
	}
}
