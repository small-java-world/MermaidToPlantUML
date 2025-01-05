# サンプル図一覧

このディレクトリには、Mermaid記法で書かれたクラス図のサンプルが含まれています。
このサンプルを使用して、Mermaid形式からPlantUMLを経由して画像を生成する方法を説明します。

## 前提条件

1. Java（JRE 8以上）がインストールされていること
2. `lib`ディレクトリに`plantuml-1.2024.8.jar`が配置されていること
3. Go言語の開発環境がインストールされていること

## ファイル構成

- `domain_model.mmd`: ECサイトのドメインモデルを表現したMermaid形式のクラス図

## 変換手順

### 1. mermaid2plantumlのビルド

```bash
# プロジェクトのルートディレクトリで実行
go build
# => mermaid2plantuml.exe（Windows）または mermaid2plantuml（Mac/Linux）が生成されます
```

### 2. Mermaid形式からPlantUML形式への変換

Windows環境:
```cmd
# Mermaid形式（.mmd）からPlantUML形式（.puml）に変換
.\mermaid2plantuml.exe samples/domain_model.mmd
# => samples/domain_model.puml が生成されます
```

Mac/Linux環境:
```bash
# Mermaid形式（.mmd）からPlantUML形式（.puml）に変換
./mermaid2plantuml samples/domain_model.mmd
# => samples/domain_model.puml が生成されます
```

### 3. PlantUML形式から画像への変換

Windows環境:
```cmd
# プロジェクトのルートディレクトリで実行
# PNGで出力（デフォルト）
java -jar lib/plantuml-1.2024.8.jar samples/domain_model.puml
# => samples/domain_model.png が生成されます

# SVGで出力
java -jar lib/plantuml-1.2024.8.jar -tsvg samples/domain_model.puml
# => samples/domain_model.svg が生成されます

# 出力先を指定（outディレクトリ）
java -jar lib/plantuml-1.2024.8.jar -o out samples/domain_model.puml
# => out/domain_model.png が生成されます
```

Mac/Linux環境:
```bash
# プロジェクトのルートディレクトリで実行
# PNGで出力（デフォルト）
java -jar lib/plantuml-1.2024.8.jar samples/domain_model.puml
# => samples/domain_model.png が生成されます

# SVGで出力
java -jar lib/plantuml-1.2024.8.jar -tsvg samples/domain_model.puml
# => samples/domain_model.svg が生成されます

# 出力先を指定（outディレクトリ）
java -jar lib/plantuml-1.2024.8.jar -o out samples/domain_model.puml
# => out/domain_model.png が生成されます
```

## 生成されるファイル

1. ツールのビルド
   - 入力: Goのソースコード
   - 出力: `mermaid2plantuml.exe`（Windows）または`mermaid2plantuml`（Mac/Linux）

2. 中間ファイル
   - 入力: `samples/domain_model.mmd`（Mermaid形式）
   - 出力: `samples/domain_model.puml`（PlantUML形式）

3. 最終出力
   - PNG形式（デフォルト）
     - 出力: `samples/domain_model.png`
     - 用途: Web表示、ドキュメント埋め込み
   - SVG形式（オプション）
     - 出力: `samples/domain_model.svg`
     - 用途: 拡大縮小可能な図、Web表示

## サンプルの内容

`domain_model.mmd`は以下のような構成のECサイトのドメインモデルを表現しています：

- User（ユーザー）クラス
  - 基本情報（ID、名前、メール、作成日時）
  - ユーザー登録・更新機能

- UserProfile（ユーザープロファイル）クラス
  - プロフィール情報（ニックネーム、アバター、誕生日）
  - プロフィール更新機能

- Order（注文）クラス
  - 注文情報（ID、ユーザーID、ステータス、注文日、合計金額）
  - 注文作成・キャンセル・完了機能

- OrderItem（注文商品）クラス
  - 商品情報（注文ID、商品ID、数量、価格）
  - 小計計算機能

- Product（商品）クラス
  - 商品情報（ID、名前、説明、価格、在庫数）
  - 商品更新・在庫調整機能 