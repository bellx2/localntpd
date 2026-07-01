# localntpd

[English version](README.md)

自分のPCの時刻を提供するシンプルなNTPサーバーです。Go製でマルチプラットフォーム
(Windows / macOS / Linux) に対応し、OSのシステムサービスとして登録できます。

## 特徴

- ローカルマシンの時刻をNTPプロトコル (UDP) で提供
- マルチプラットフォーム対応の単一バイナリ (Windows / macOS / Linux)
- OSのシステムサービスとして登録可能 ([kardianos/service](https://github.com/kardianos/service) を利用)
- 上位NTPサーバーとは同期しない、スタンドアロンなローカルクロック (stratum / reference ID `LOCL`)

## ビルド

```bash
go build -o localntpd .
```

## 使い方

```bash
localntpd [コマンド] [オプション]
```

### コマンド

| コマンド    | 説明                             |
|-------------|----------------------------------|
| `run`       | フォアグラウンドで起動 (デフォルト) |
| `install`   | システムサービスとして登録       |
| `uninstall` | サービスを削除                   |
| `start`     | サービスを開始                   |
| `stop`      | サービスを停止                   |
| `restart`   | サービスを再起動                 |
| `status`    | サービスの状態を表示             |
| `help`      | ヘルプを表示                     |

### オプション

| オプション        | デフォルト | 説明                                 |
|-------------------|-----------|--------------------------------------|
| `-addr string`    | `:123`    | リッスンアドレス (例: `:123`, `0.0.0.0:123`) |
| `-stratum uint`   | `2`       | ストラタム (1〜15)                   |

### 例

```bash
localntpd run -addr :12345    # 非特権ポートで起動
localntpd install             # サービス登録 (管理者権限が必要)
localntpd start
```

## 注意

- ポート123 (標準NTPポート) を使用するには管理者/root権限が必要です。`install` や
  デフォルトの `:123` での起動が該当します。
- 開発・テスト時は `-addr :12345` など非特権ポートを使用してください。

## ライセンス

ライセンス情報はリポジトリを参照してください。
