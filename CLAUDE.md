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

## アーキテクチャ

2層構成。サービスライフサイクルとNTPプロトコル実装を分離している。

- `main.go` — CLIエントリポイントとサービス管理。`github.com/kardianos/service` でOS横断のサービス登録を行う。`os.Args[1]` を見て `install`/`start` 等のサービス制御コマンドと通常起動を振り分ける。サービスの `Start`/`Stop` は `context.CancelFunc` でサーバーgoroutineの停止を制御する（ブロックしない設計）。
- `internal/ntp/server.go` — NTPプロトコルの中核。`Server.Run(ctx)` がUDP 48バイトパケットを受信し、`buildResponse` で応答を組み立てる。`ctx` のキャンセルを検知するため `SetReadDeadline` で1秒ごとにループを回している（ブロッキングReadを避ける）。stratum=2 / refID=`LOCL` の固定応答で、上位NTPサーバーとは同期しないローカルクロック専用サーバー。

## NTPパケット実装の注意点

- タイムスタンプはNTPエポック (1900-01-01) 基準。Unix時刻に `epochOffset` (2208988800) を加算して変換する。小数部は `(nanosec << 32) / 1e9` で32bit固定小数に変換。
- 応答の origin timestamp はクライアントのtransmit timestamp (リクエストの `[40:48]`) をそのままエコーする。ここを変更するとクライアント側の往復遅延計算が壊れる。
- パケットのバイトオフセット (LI/VN/Mode=`resp[0]`, stratum=`resp[1]`, transmit=`resp[40:48]` 等) はRFC 5905準拠。`server_test.go` がこのオフセットを検証しているので、`buildResponse` を変更する際はテストも合わせて更新する。
