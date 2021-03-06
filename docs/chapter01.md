# 1. データベースの要件

## 1.1 なぜデータベースシステムなのか

- データベース => コンピューター上に保存されたデータの集まり
- 一般的にはデータベースのデータはレコードに分類

レコードを管理するための要件

- 永続性
  - レコードが消えない
- 共有できる
  - 複数の同時使用者がいる
- 正確性を保つ
  - データベースの中身が信用できる
- 大容量である
  - 100万、10億などのレコードを保管できる
- 可用性
  - 使用可能である

### 1.1.1 レコードストレージ

一般的なのはテキストファイル
ex. 行がレコードに対応し、値がタブ区切り

- メリット
  - テキストエディタで編集可能
- デメリット
  - 膨大なレコード数だと r/w が遅い。検索も遅い

### 1.1.2 マルチユーザーアクセス

- 同時データアクセス
  - 一定の制限がないとデータの正確性が担保できない
- データベースには変更を巻き戻す機能、変更を永続的にする能力が必要
  - いわゆる「コミット」

### 1.1.3 大惨事の扱い

データベース更新するプログラムを実行中に、
データベースシステムがクラッシュした際、
プログラムの更新を巻き戻せるか。

### 1.1.4 メモリ管理

データベースは永続メモリに保存する必要
ex. ディスクドライブ、フラッシュドライブ

どちらもRAMよりパフォーマンスが悪過ぎる。

データベースシステムは、難問に直面する:

RAMよりも大容量のレコードを管理し、
複数の同時ユーザーアクセスを許し、
以前の状態に完璧に復旧できて、
レスポンスタイムを維持する。