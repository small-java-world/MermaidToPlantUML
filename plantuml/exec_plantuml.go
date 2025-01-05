package plantuml

import (
	"fmt"
	"os"
	"os/exec"
)

// PlantUMLExecutor はPlantUMLコマンドの実行を管理します
type PlantUMLExecutor struct {
	plantumlPath string
}

// NewPlantUMLExecutor は新しいPlantUMLExecutorインスタンスを作成します
func NewPlantUMLExecutor() *PlantUMLExecutor {
	return &PlantUMLExecutor{
		plantumlPath: "plantuml", // デフォルトではPATHから検索
	}
}

// GenerateImage はPlantUMLファイルから画像を生成します
func (e *PlantUMLExecutor) GenerateImage(pumlFile string, format string) error {
	// 入力ファイルの存在確認
	if _, err := os.Stat(pumlFile); os.IsNotExist(err) {
		return fmt.Errorf("入力ファイルが存在しません: %s", pumlFile)
	}

	// フォーマットの検証
	switch format {
	case "png", "svg", "pdf":
		// サポートされているフォーマット
	default:
		return fmt.Errorf("サポートされていないフォーマット: %s", format)
	}

	// PlantUMLコマンドが利用可能か確認
	if _, err := exec.LookPath(e.plantumlPath); err != nil {
		return fmt.Errorf("PlantUMLが利用できません。Java / Docker が実装されているか確認してください: %v", err)
	}

	// 出力フォーマットに応じたオプションを設定
	args := []string{"-t" + format}
	args = append(args, pumlFile)

	// PlantUMLコマンドを実行
	cmd := exec.Command(e.plantumlPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("PlantUMLの実行に失敗しました: %v", err)
	}

	return nil
}

// SetPlantUMLPath はPlantUMLコマンドのパスを設定します
func (e *PlantUMLExecutor) SetPlantUMLPath(path string) {
	e.plantumlPath = path
}
