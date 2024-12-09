#!/bin/bash

# 使用法を表示
if [ "$#" -ne 3 ]; then
  echo "Usage: $0 <source_directory> <line_count> <target_directory>"
  exit 1
fi

# 引数を変数に格納
source_directory=$1
line_count=$2
target_directory=$3

# ソースディレクトリが存在するか確認
if [ ! -d "$source_directory" ]; then
  echo "Error: Source directory $source_directory does not exist."
  exit 1
fi

# ターゲットディレクトリが存在するか確認し、なければ作成
if [ ! -d "$target_directory" ]; then
  mkdir -p "$target_directory"
fi

# 指定した行数と一致するファイルを移動
find "$source_directory" -type f | while read -r file; do
  actual_lines=$(wc -l < "$file")
  if [ "$actual_lines" -eq "$line_count" ]; then
    # 現在のタイムスタンプを取得
    timestamp=$(date "+%Y%m%d_%H%M%S")
    # 元のファイル名を取得
    filename=$(basename "$file")
    # 新しいファイル名を生成
    new_filename="${timestamp}_${filename}"
    # ファイルを移動しつつ名前を変更
    mv "$file" "$target_directory/$new_filename"
    echo "Moved: $file to $target_directory/$new_filename"
  fi
done
