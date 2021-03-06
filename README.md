# pstore

## これなに

* AWS System Manager のパラメータストアを操作する小さなコマンドラインツールです
* direnv や jq 等と組み合わせて利用してください

## 使い方

### ヘルプ

```sh
$ pstore -h
  -csv
        CSV 形式で出力する
  -del
        パラメータを削除する
  -endpoint string
        AWS API のエンドポイントを指定.
  -get
        パラメータの値を取得する
  -insecure
        SecureString を出力する
  -json
        JSON 形式で出力する
  -list
        StringList でパラメータを追加する
  -name string
        パラメータの名前を指定する
  -overwrite
        パラメータを上書きする
  -profile string
        Profile 名を指定.
  -put
        パラメータを追加する
  -region string
        Region 名を指定. (default "ap-northeast-1")
  -role string
        Role ARN を指定.
  -secure
        SecureString でパラメータを追加する
  -value string
        パラメータ名を値を指定する
  -version
        バージョンを出力.
```

### パラメータの一覧を取得

```sh
$ pStore
+-------------------------------------+--------------------------+--------------+---------------------+
|                NAME                 |          VALUE           |     TYPE     |  LASTMODIFIEDDATE   |
+-------------------------------------+--------------------------+--------------+---------------------+
| /123456/88888                       | kawahara-test            | StringList   | 2018-09-29 08:09:43 |
| test.test1                          | ******************       | SecureString | 2018-09-28 22:42:23 |
+-------------------------------------+--------------------------+--------------+---------------------+
```

### パラメータの追加

```sh
$ pStore -put -name="foooooon" -value="baaaaaaarn"

$ pStore
+-------------------------------------+--------------------------+--------------+---------------------+
|                NAME                 |          VALUE           |     TYPE     |  LASTMODIFIEDDATE   |
+-------------------------------------+--------------------------+--------------+---------------------+
| /123456/88888                       | kawahara-test            | StringList   | 2018-09-29 08:09:43 |
| foooooon                            | baaaaaaarn               | String       | 2018-09-29 08:37:53 |
| test.test1                          | ******************       | SecureString | 2018-09-28 22:42:23 |
+-------------------------------------+--------------------------+--------------+---------------------+
```

### パラメータの上書き

```sh
$ pStore -put -name="foooooon" -value="bazooooooon" -overwrite
$ pStore
+-------------------------------------+--------------------------+--------------+---------------------+
|                NAME                 |          VALUE           |     TYPE     |  LASTMODIFIEDDATE   |
+-------------------------------------+--------------------------+--------------+---------------------+
| /123456/88888                       | kawahara-test            | StringList   | 2018-09-29 08:09:43 |
| foooooon                            | bazooooooon              | String       | 2018-09-29 08:38:51 |
| test.test1                          | ******************       | SecureString | 2018-09-28 22:42:23 |
+-------------------------------------+--------------------------+--------------+---------------------+
```

### パラメータの削除

```sh
$ pStore -del -name="foooooon"
$ pStore
+-------------------------------------+--------------------------+--------------+---------------------+
|                NAME                 |          VALUE           |     TYPE     |  LASTMODIFIEDDATE   |
+-------------------------------------+--------------------------+--------------+---------------------+
| /123456/88888                       | kawahara-test            | StringList   | 2018-09-29 08:09:43 |
| test.test1                          | ******************       | SecureString | 2018-09-28 22:42:23 |
+-------------------------------------+--------------------------+--------------+---------------------+
```

## todo

* 色々