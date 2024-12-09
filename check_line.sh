#!/bin/bash

# 使用法を表示
if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <directory> <line_count>"
  exit 1
fi

# 引数を変数に格納
directory=$1
line_count=$2

# ディレクトリが存在するか確認
if [ ! -d "$directory" ]; then
  echo "Error: Directory $directory does not exist."
  exit 1
fi

# 検索して行数が一致しないファイルを表示
find "$directory" -type f | while read -r file; do
  actual_lines=$(wc -l < "$file")
  if [ "$actual_lines" -ne "$line_count" ]; then
    echo "$file"
  fi
done
