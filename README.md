# check-lastlog

`check-lastlog` は、Linux 環境で `/var/log/lastlog` を解析し、一定期間ログインしていないユーザーを検出する [Mackerel](https://mackerel.io/) などの監視ツールで利用できるチェックプラグインです。

長期間ログインしていないユーザーアカウントは、不要なアカウントや漏洩リスクのあるアカウントの兆候となる場合があります。本ツールを使うことで、そのようなアカウントを自動的に検出し、アラートとして通知できます。

## 特徴

- `/var/log/lastlog` を直接読み取り、各ユーザーの最終ログイン日時を確認
- 警告（warning）・異常（critical）のしきい値を日数で指定可能
- システムアカウントを除外するための UID 範囲指定
- ログイン不能なシェル（`nologin` など）を持つユーザーを自動除外
- ホワイトリストで特定ユーザーを常に除外可能
- Mackerel チェックプラグインとして利用可能

## 動作環境

Linux の amd64 / arm64 のみ対応しています。

## インストール

### mkr プラグインとしてインストール

```bash
$ mkr plugin install monitoring-forge/check-lastlog
```

インストール後、`/opt/mackerel-agent/plugins/bin/check-lastlog` などに配置されます。

### GitHub リリースからダウンロード

[GitHub Releases](https://github.com/monitoring-forge/check-lastlog/releases/latest) から、最新のバイナリをダウンロードできます。

```bash
# 例: Linux amd64 の場合
$ curl -L -o check-lastlog.zip "https://github.com/monitoring-forge/check-lastlog/releases/latest/download/check-lastlog_linux_amd64.zip"
$ unzip check-lastlog.zip -d check-lastlog
$ cd check-lastlog
$ ./check-lastlog -h
```

必要に応じて、Mackerel エージェントのプラグインディレクトリなどに配置してください。

```bash
$ sudo install -m 755 check-lastlog /opt/mackerel-agent/plugins/bin/check-lastlog
```

## 使い方

```
$ ./check-lastlog -h
Usage:
  check-lastlog [OPTIONS]

Application Options:
      --before=           [Deprecated] Check for users whose login is older than DAYS
  -w, --warning=          warning if users whose login is older than DAYS (default: 60)
  -c, --critical=         critical if users whose login is older than DAYS (default: 85)
      --min-uid=          min uid to check lastlog (default: 500)
      --max-uid=          max uid to check lastlog (default: 60000)
      --white-user-names= comma separeted user names that white
      --lastlog-file=     lastlog file path (default: /var/log/lastlog)
      --passwd-file=      passwd file path (default: /etc/passwd)
  -V, --verbose           Show verbose log
  -v, --version           Show version

Help Options:
  -h, --help              Show this help message
```

### オプションの説明

| オプション | 説明 | デフォルト |
|---|---|---|
| `-w`, `--warning` | 最終ログインからこの日数を超えた場合に **WARNING** | 60 |
| `-c`, `--critical` | 最終ログインからこの日数を超えた場合に **CRITICAL** | 85 |
| `--min-uid` | チェック対象とする最小 UID | 500 |
| `--max-uid` | チェック対象とする最大 UID | 60000 |
| `--white-user-names` | 常に除外するユーザー名をカンマ区切りで指定 | （なし） |
| `--lastlog-file` | 読み込む lastlog ファイルのパス | `/var/log/lastlog` |
| `--passwd-file` | 読み込む passwd ファイルのパス | `/etc/passwd` |
| `-V`, `--verbose` | デバッグ用の詳細ログを標準エラーに出力 | false |

## 実行例

### 基本的な使い方

```bash
$ ./check-lastlog
OK: No users were found who have not logged in recently
```

```bash
$ ./check-lastlog
CRITICAL: Found users who have not logged in recently: testuser(129 days), sampleuser(106 days)
$ echo $?
2
```

### しきい値を変更する

```bash
$ ./check-lastlog -w 30 -c 60
WARNING: Found users who have not logged in recently: alice(45 days)
```

### ホワイトリストで特定ユーザーを除外する

バッチ用アカウントやメンテナンス用アカウントなど、意図的に長期間ログインしないユーザーを除外できます。

```bash
$ ./check-lastlog --white-user-names backup,batch,deploy
OK: No users were found who have not logged in recently
```

### 絞り込み対象を変更する

UID 範囲を変更して、通常ユーザーのみを対象にすることもできます。

```bash
$ ./check-lastlog --min-uid 1000 --max-uid 2000
OK: No users were found who have not logged in recently
```

## Mackerel との連携

Mackerel エージェントの設定ファイル（通常は `/etc/mackerel-agent/mackerel-agent.conf`）に `[plugin.checks]` を追加してください。

### 基本的な設定例

```ini
[plugin.checks.lastlog]
command = ["/opt/mackerel-agent/plugins/bin/check-lastlog", "--white-user-names", "backup,batch"]
```

### しきい値を変更した例

```ini
[plugin.checks.lastlog]
command = ["/opt/mackerel-agent/plugins/bin/check-lastlog", "-w", "30", "-c", "60", "--white-user-names", "deploy,ansible"]
```

### 実行間隔を変更する例

```ini
[plugin.checks.lastlog]
command = ["/opt/mackerel-agent/plugins/bin/check-lastlog", "--white-user-names", "backup,batch"]
check_interval = 30
notification_interval = 60
```

設定後、以下のコマンドで Mackerel エージェントを再起動してください。

```bash
$ sudo systemctl restart mackerel-agent
```

エージェントがプラグインを実行すると、Mackerel のホスト詳細画面の「モニタリング」→「チェック監視」に結果が表示され、WARNING / CRITICAL の状態でアラートが発報されます。

## 注意事項

- 本ツールは Linux の `lastlog` バイナリ形式を解析するため、Linux の amd64 / arm64 以外では動作しません。
- `--before` オプションは非推奨です。新規には `-w` / `-c` をご利用ください。
- ログイン不能なシェルを持つユーザー、および UID が `--min-uid` 未満・`--max-uid` を超えるユーザーは自動的に除外されます。
- 以下のシェルを持つユーザーは、ログイン不能と判定されてチェック対象から除外されます。
  - `/bin/sync`
  - `/sbin/halt`
  - `/sbin/nologin`
  - `/sbin/shutdown`
  - `/usr/sbin/nologin`
  - `/usr/bin/false`
  - `/bin/false`
- `lastlog` ファイルへのアクセスには通常 `root` 権限が必要です。Mackerel エージェントは root で動作することが前提ですが、手動実行時は `sudo` が必要になる場合があります。

## ライセンス

MIT License