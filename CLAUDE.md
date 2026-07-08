# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

- 自分のPCの時刻を利用したシンプルなNTPサーバーアプリ
- Go製でマルチプラットフォーム対応 (Windows/macOS/Linux)
- OSのシステムサービスとして登録・起動可能

## コマンド

```bash
go build -o localntpd .          # ビルド
go test ./...                    # 全テスト
go test ./internal/ntp -run TestBuildResponse -v   # 単一テスト実行
go run . run -addr :12345        # 非特権ポートでフォアグラウンド起動

./localntpd install              # システムサービスとして登録 (要管理者/root)
./localntpd start / stop / restart / status / uninstall
```

ポート123は特権ポートのため、`install` およびデフォルトの `:123` での起動には管理者/root権限が必要。開発・テスト時は `-addr :12345` など非特権ポートを使う。

### Windowsビルド (アイコン/バージョン情報付き)

`.exe` にアイコンとバージョン情報を埋め込むには、ビルド前に `go-winres` でリソースを生成する。生成された `rsrc_windows_*.syso` を `go build` が自動リンクする。

```bash
go run github.com/tc-hib/go-winres@v0.3.3 make --arch amd64,arm64
GOOS=windows GOARCH=amd64 go build -o localntpd.exe .
```

- Go 1.26以上が必要。`go-winres` の事前インストールは不要（`go run` で実行）。
- リソース定義は `winres/winres.json`、アイコンは `icon16/32/48/256.png`（元データ: `icon.svg`）。差し替えたら上記の `go run ... make` を再実行する。
- `*.syso` は生成物のため `.gitignore` 済み。`winres/` 一式のみコミットする。
- `.syso` はファイル名のビルド制約により **Windowsビルド時のみ** リンクされ、macOS/Linuxビルドには影響しない。

## アーキテクチャ

2層構成。サービスライフサイクルとNTPプロトコル実装を分離している。

- `main.go` — CLIエントリポイントとサービス管理。`github.com/kardianos/service` でOS横断のサービス登録を行う。`os.Args[1]` を見て `install`/`start` 等のサービス制御コマンドと通常起動を振り分ける。`install` 時の `-addr`/`-stratum` はサービス起動引数として保存する。`Start` はサーバーgoroutineを起動し、`Stop` は `context` キャンセル後に終了を待つ。
- `internal/ntp/server.go` — NTPプロトコルの中核。`Server.Run(ctx)` がUDP 48バイトパケットを受信し、`buildResponse` で応答を組み立てる。`ctx` のキャンセルを検知するため `SetReadDeadline` で1秒ごとにループを回している（ブロッキングReadを避ける）。stratum=2 / refID=`LOCL` の固定応答で、上位NTPサーバーとは同期しないローカルクロック専用サーバー。

## NTPパケット実装の注意点

- タイムスタンプはNTPエポック (1900-01-01) 基準。Unix時刻に `epochOffset` (2208988800) を加算して変換する。小数部は `(nanosec << 32) / 1e9` で32bit固定小数に変換。
- 応答の origin timestamp はクライアントのtransmit timestamp (リクエストの `[40:48]`) をそのままエコーする。ここを変更するとクライアント側の往復遅延計算が壊れる。
- パケットのバイトオフセット (LI/VN/Mode=`resp[0]`, stratum=`resp[1]`, transmit=`resp[40:48]` 等) はRFC 5905準拠。`server_test.go` がこのオフセットを検証しているので、`buildResponse` を変更する際はテストも合わせて更新する。
