# tic-tac-toe-ai-backend
AIと対戦できる三目並べのbackendの実装です。goで記述されています。

## 実行方法
```shell
go run main.go
```

[`http://localhost:8080/`](http://localhost:8080/)からアクセスできます。

## API
- `api/state_action?state=10`
  - response
    - `{"action": 5}`
  - description
    - 三目並べの盤面を整数値に変換してクエリパラメータ`state`に渡すと、AIがコマを置く位置を整数値で返します。

## frontend
frntendの実装は以下のリポジトリにあります。

[https://github.com/yutatech/tic-tac-toe-ai-frontend](https://github.com/yutatech/tic-tac-toe-ai-frontend)

## AI
強化学習で学習したQ値のtableを読み込みます。強化学習の実装は以下のリポジトリにあります。

[https://github.com/yutatech/reinforcement-learning-tic-tac-toe](https://github.com/yutatech/reinforcement-learning-tic-tac-toe)

上記リポジトリのコードを実行して生成されるバイナリファイル`q_table_go`を`main.go`と同じディレクトリに配置します。