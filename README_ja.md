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

### Windowsビルド (アイコン・バージョン情報付き)

`.exe` にアイコンとバージョン情報を埋め込むには、ビルド前に
[`go-winres`](https://github.com/tc-hib/go-winres) でWindowsリソースを生成します。
生成された `rsrc_windows_*.syso` は `go build` が自動でリンクします。

```bash
go install github.com/tc-hib/go-winres@v0.3.3   # 初回のみ
go-winres make --arch amd64,arm64                # winres/ から .syso を生成
GOOS=windows GOARCH=amd64 go build -o localntpd.exe .
```

Go 1.26以上が必要です。リソース定義は `winres/winres.json`、アイコンは
`winres/icon16.png` / `icon32.png` / `icon48.png` / `icon256.png`
（元データ: `winres/icon.svg`）です。差し替えて `go-winres make` を再実行すれば
カスタマイズできます。`.syso` はWindowsビルド時のみリンクされるため、
macOS/Linuxビルドには影響しません。

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
localntpd run -addr :12345              # 非特権ポートで起動
localntpd install                       # サービス登録 (管理者権限が必要)
localntpd install -addr :12345          # カスタムアドレスでサービス登録
localntpd start
```

## 注意

- ポート123 (標準NTPポート) を使用するには管理者/root権限が必要です。`install` や
  デフォルトの `:123` での起動が該当します。
- 開発・テスト時は `-addr :12345` など非特権ポートを使用してください。
- `install` に渡したオプション (`-addr`, `-stratum`) はサービス起動引数として保存されます。

## ライセンス

ライセンス情報はリポジトリを参照してください。
