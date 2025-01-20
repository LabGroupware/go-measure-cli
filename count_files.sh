#!/bin/bash

TARGET_DIR=${1:-"."}

THREAD_RESULT=("Thread" "")
CASE_RESULT=("Case" "")
COUNT_RESULT=("Count" "")

while read -r dir; do
    THREAD_RESULT+=($(echo "$dir" | grep -oP "thread_\K[0-9]+"))
    CASE_RESULT+=($(basename "$dir"))
    COUNT_RESULT+=($(find "$dir" -type f | wc -l))

    # printf "%s  %s  %s\n" "${THREAD_RESULT[-1]}" "${CASE_RESULT[-1]}" "${COUNT_RESULT[-1]}"
done < <(find "$TARGET_DIR" -type d -regex ".*/thread_[0-9]+/[^/]+")

max_thread_len=$(printf "%s\n" "${THREAD_RESULT[@]}" | awk '{ if (length > max) max = length } END { print max }')
max_case_len=$(printf "%s\n" "${CASE_RESULT[@]}" | awk '{ if (length > max) max = length } END { print max }')
max_count_len=$(printf "%s\n" "${COUNT_RESULT[@]}" | awk '{ if (length > max) max = length } END { print max }')

printf "| %-${max_thread_len}s | %-${max_case_len}s | %-${max_count_len}s |\n" "$(printf "%-${max_thread_len}s" "-" | sed "s/ /-/g")" "$(printf "%-${max_case_len}s" "-" | sed "s/ /-/g")" "$(printf "%-${max_count_len}s" "-" | sed "s/ /-/g")"
printf "| %-${max_thread_len}s | %-${max_case_len}s | %-${max_count_len}s |\n" "${THREAD_RESULT[0]}" "${CASE_RESULT[0]}" "${COUNT_RESULT[0]}"
printf "| %-${max_thread_len}s | %-${max_case_len}s | %-${max_count_len}s |\n" "$(printf "%-${max_thread_len}s" "-" | sed "s/ /-/g")" "$(printf "%-${max_case_len}s" "-" | sed "s/ /-/g")" "$(printf "%-${max_count_len}s" "-" | sed "s/ /-/g")"

# # データ出力
for ((i = 2; i < ${#THREAD_RESULT[@]}; i++)); do
    printf "| %-${max_thread_len}s | %-${max_case_len}s | %-${max_count_len}s |\n" "${THREAD_RESULT[i]}" "${CASE_RESULT[i]}" "${COUNT_RESULT[i]}"
done

printf "| %-${max_thread_len}s | %-${max_case_len}s | %-${max_count_len}s |\n" "$(printf "%-${max_thread_len}s" "-" | sed "s/ /-/g")" "$(printf "%-${max_case_len}s" "-" | sed "s/ /-/g")" "$(printf "%-${max_count_len}s" "-" | sed "s/ /-/g")"