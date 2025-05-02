#!/usr/bin/env bash
#
# Description: Generate a markdown summary of test results with symbols and include overall status in the title.
# Usage:       ./test-summary.sh <Category> <Test1=status> [Test2=status ...]
# Exit Code:   Always exits with 0 so that downstream steps can run regardless of failures.
#

set -euo pipefail
IFS=$'\n\t'

declare -g all_passed_flag=true

write_header() {
    local file="$1"
    local category="$2"
    printf -- '## %s Tests Summary\n\n' "$category" >"$file"
}

process_tests() {
    local file="$1"
    shift
    local all_passed=true

    printf -- '| Test | Status |\n| ---- | ------ |\n' >>"$file"

    for arg in "$@"; do
        if [[ "$arg" =~ ^([^=]+)=(.+)$ ]]; then
            local name="${BASH_REMATCH[1]}"
            local status="${BASH_REMATCH[2]}"
        else
            printf -- 'Warning: invalid argument "%s"\n' "$arg" >&2
            continue
        fi

        if [[ "$status" == "success" ]]; then
            symbol='✅'
        else
            symbol='❌'
            all_passed=false
        fi

        printf -- '| %s | %s |\n' "$name" "$symbol" >>"$file"
    done

    all_passed_flag=$all_passed
}

write_summary() {
    local file="$1"
    local lines pattern total failed passed
    pattern='^(##|\| Test \| Status \||\| ---- \| ------ \||$)'
    mapfile -t lines < <(grep -vE "$pattern" "$file")
    total=${#lines[@]}
    failed=0
    for line in "${lines[@]}"; do
        [[ "$line" == *'❌'* ]] && ((failed++))
    done
    passed=$((total - failed))

    {
        printf -- '### Detailed Summary\n\n'
        printf -- '- **Total tests:** %d\n' "$total"
        printf -- '- **Passed:**      %d ✅\n' "$passed"
        printf -- '- **Failed:**      %d ❌\n' "$failed"
    } >>"$file"
}

main() {
    if (($# < 2)); then
        exit 1
    fi

    local category="$1"
    shift
    local output="./${category}-results.md"

    write_header "$output" "$category"
    process_tests "$output" "$@"
    write_summary "$output"

    local status_symbol status_text
    if $all_passed_flag; then
        status_symbol='✅'
        status_text='SUCCESS'
    else
        status_symbol='❌'
        status_text='FAILED'
    fi

    sed -i "1s/.*/## ${category} Tests Summary — ${status_symbol} ${status_text}/" "$output"
}

main "$@"
exit 0
