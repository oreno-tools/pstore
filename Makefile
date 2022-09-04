default: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test-setup: ## テストのセットアップ
	@cd tests; ./_setup.sh

test-setup: ## テストの実行
	@go test -v

test-teardown: ## テストの環境の破棄
	@cd tests; ./_teardown.sh

build: ## バイナリをビルドする
	@./build.sh pstore.go

release: ## バイナリをリリースする. 引数に `_VER=バージョン番号` を指定する.
	@ghr -u inokappa -r pstore v${_VER} ./pkg/