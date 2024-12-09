#!/bin/bash

# 使用法を表示
if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <source_directory>"
  exit 1
fi

# 引数を変数に格納
source_directory=$1

# ソースディレクトリが存在するか確認
if [ ! -d "$source_directory" ]; then
  echo "Error: Source directory $source_directory does not exist."
  exit 1
fi

# ディレクトリ構造の定義と行数マッピング
base_output_dir="./datasets/wait_saga/thread_1"
declare -A categories=(
  ["team"]=5
  ["file_object"]=4
  ["organization"]=6
  ["task"]=7
  ["user_profile"]=4
)

# ネストされたすべてのファイルを処理
find "$source_directory" -type f -name "socket_subscribe_output_create_*_*.csv" | while read -r file; do
  # 現在のタイムスタンプを取得
  timestamp=$(date "+%Y%m%d_%H%M%S")

  # ファイル名からカテゴリを抽出
  filename=$(basename "$file")
  echo $filename
  category=$(echo "$filename" | sed -n 's/^socket_subscribe_output_create_\([a-z_]*\)_.*\.csv$/\1/p')
  echo $category

  # カテゴリが定義されている場合のみ処理
  if [[ -n "$category" && -n "${categories[$category]}" ]]; then
    required_lines=${categories[$category]}
    target_dir="$base_output_dir/create_$category"

    # 行数を取得
    actual_lines=$(wc -l < "$file")

    # 行数が一致する場合のみ移動
    if [ "$actual_lines" -eq "$required_lines" ]; then
      # ターゲットディレクトリが存在しない場合は作成
      mkdir -p "$target_dir"

      # 新しいファイル名を生成
      new_filename="${timestamp}_${filename}"

      # ファイルを移動
      mv "$file" "$target_dir/$new_filename"
      echo "Moved: $file -> $target_dir/$new_filename"
    else
      echo "Skipped: $file (lines: $actual_lines, required: $required_lines)"
    fi
  else
    echo $category
    echo "Skipped: $file (unknown category)"
  fi
done
