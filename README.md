# Article API Assignment

## 概要
本課題では、外部 API から記事データを取得して保存する API と、保存された記事を一覧で取得する API を実装しました。実装にあたっては、AI にすべてを任せるのではなく、自分が書いたコードを言葉で説明できる状態になることを目標にしました。

1周目は API 全体の流れを掴むために AI 主導で実装し、2周目では責務分離や環境変数の扱いを意識して書き直しました。本リポジトリは、提出用として設計意図を意識しながら整理したものです。

---

## 使用技術

### 言語・フレームワーク
Go / chi (Web Framework) / GORM (ORM)

### データベース
SQLite
（1周目では PostgreSQL を試しましたが、外部 API の記事を保存・取得する用途としては過剰だと感じたため、今回は軽量で構築が容易な SQLite を選択しました）

### フロントエンド
Next.js
（記事が正しく取得・保存されていることを視覚的に確認するための簡易的な検証用です）

---

## ディレクトリ構成

### internal/model
データ構造（Article）の定義を担当します。

### internal/handler
HTTP リクエストの受け取りとレスポンスの返却を担当します。

### internal/service
外部 API の取得やデータ変換など、処理の流れ（ワークフロー）を制御します。

### internal/repository
データベースへの永続化処理のみに専念します。

---

## API 一覧

### GET /articles
保存されているすべての記事を JSON 形式で返します。

### POST /articles/import
外部 API から記事データを取得し、DB に保存します。一連の変換・保存処理は service 層で管理しています。

### POST /articles
タイトルと本文を受け取り、新しい記事を保存します。DB 用のモデルと入力専用の構造体を分離して定義しています。

### DELETE /articles/{id}
指定された ID の記事を削除します。成功時は 204 No Content を返します。

### PATCH /articles/{id}/pin
記事のピン留め状態（true / false）を切り替えます。状態変更の判断を伴うため service 層で制御しています。

---

## 環境変数
環境ごとに変わる設定値を管理するため .env ファイルを使用しています。必要な変数を共有するために .env.example を用意しています。

設定値の例：
EXTERNAL_API_URL=...

---

## 起動方法

1. 設定ファイルの準備
cd backend

# 外部 API の記事を取り込む
curl -X POST http://localhost:8080/articles/import

# 保存された記事を取得する
curl http://localhost:8080/articles

cp .env.example .env


2. サーバーの起動
go run main.go

---

## 設計で意識したこと
handler / service / repository に責務を分離することで、処理の流れと処理の詳細を分けて管理できるようにしました。handler は HTTP の入り口に専念し、repository は DB 操作にのみ責務を限定しています。その間にあるロジックや成功・失敗の判断を service 層に集約する構成をとっています。

今回は Docker やクリーンアーキテクチャは採用していません。これは、浅い理解のままツールを導入するよりも、Go の基礎を固めたうえで今後適切に導入したいと考えたためです。