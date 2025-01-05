# サンプル図一覧

このディレクトリには、Mermaid記法で書かれたクラス図のサンプルが含まれています。

## 前提条件

1. Java（JRE 8以上）がインストールされていること
2. `lib`ディレクトリに`plantuml-1.2024.8.jar`が配置されていること

## ファイル構成

- `domain_model.mmd`: ECサイトのドメインモデルを表現したクラス図

## 実行例

### Windows環境での実行

```cmd
# PNGで出力（デフォルト）
> plantuml.bat domain_model.mmd

# SVGで出力
> plantuml.bat -tsvg domain_model.mmd

# 出力先を指定（outディレクトリ）
> plantuml.bat -o out domain_model.mmd
```

### Mac/Linux環境での実行

```bash
# PNGで出力（デフォルト）
$ ./plantuml.sh domain_model.mmd

# SVGで出力
$ ./plantuml.sh -tsvg domain_model.mmd

# 出力先を指定（outディレクトリ）
$ ./plantuml.sh -o out domain_model.mmd
```

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