# peruservice-auth

peruservice用の認証サービスです。

APIの仕様は`OpenApi.yml`を参照してください。

## 環境変数

`MODE`: サーバーのモード (`DEVELOPMENT`|`PRODUCTION`)<br>
`AUTH_PORT`: 認証サーバーのポート番号<br>
`AUTH_ALLOW_ORIGINS`: `PRODUCTION`モード時、サーバーとの通信を許可するオリジン (コンマ区切り) (例: `http://localhost:3000`)<br>
`PRIVATE_KEY_FILE`: JWT生成用の秘密鍵のファイルパス (例: `/path/to/key.pem``)<br>
`PUBLIC_KEY_FILE`: JWT復号用の公開鍵へのファイルパス (例: `/path/to/key.pub``)<br>
`CONFIG_JSON_FILE`: サービスの設定ファイル名とパス (例: `/path/to/config.json`)<br>
`DB_DIRECTORY`: DBのパス <br>
`DB_PORT`: DBのポート番号<br>
`DB_HOST`: DBのホスト名<br>
`DB_USER`: DBのユーザー名<br>
`DB_PASSWORD`: DBのパスワード<br>
`DB_NAME`: DBの名前<br>

## 設定項目

JSONで設定ファイルを記述します（設定ファイルは環境変数`CONFIG_JSON_FILE`で指定します）。

`refresh_token_exp_duration_hour`: リフレッシュトークンの有効期限 (時間)<br>
`access_token_exp_duration_min`: アクセストークンの有効期限 (分)<br>
