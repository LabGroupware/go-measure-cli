# Go CLI For Measure

## Setup
``` sh
asdf plugin add golang
asdf plugin add golangci-lint
asdf install
```

ログイン
``` sh
go run main.go login --config configs/config.yaml
```

リフレッシュ
```sh
go run main.go refresh --config configs/config.yaml
```

実行
```sh
go run main.go batch test --config configs/config.yaml

init/main.yaml
```

データの移動
```sh
./move_matching.sh bench/batch/
```

データ数
``` sh
./count_files.sh datasets/saga/
```