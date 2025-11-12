#!/usr/bin/env bash
#
# Description: Generate a markdown summary of test results with status indicators
# Usage:       ./test-summary.sh <Category> <Test1=status> [Test2=status ...]
# Exit Code:   Always exits with 0 so that downstream steps can run regardless of failures
#

set -euo pipefail
IFS=$'\n\t'

declare -g all_passed_flag=true

write_header() {
    local file="$1"
    local category="$2"
    printf '## %s Test Results\n\n' "$category" >"$file"
}

process_tests() {
    local file="$1"
    shift
    local all_passed=true

    printf '| Test Suite | Status |\n|------------|--------|\n' >>"$file"

    for arg in "$@"; do
        if [[ "$arg" =~ ^([^=]+)=(.+)$ ]]; then
            local name="${BASH_REMATCH[1]}"
            local status="${BASH_REMATCH[2]}"
        else
            printf 'Warning: invalid argument "%s"\n' "$arg" >&2
            continue
        fi

        case "$status" in
        "success")
            symbol='PASSED'
            ;;
        "skipped")
            symbol='SKIPPED'
            ;;
        *)
            symbol='FAILED'
            all_passed=false
            ;;
        esac

        printf '| %s | %s |\n' "$name" "$symbol" >>"$file"
    done

    all_passed_flag=$all_passed
}

write_summary() {
    local file="$1"
    local lines pattern total failed passed skipped

    pattern='^(##|\| Test Suite \| Status \||\| ------------ \| -------- \||$)'
    mapfile -t lines < <(grep -vE "$pattern" "$file")

    total=${#lines[@]}
    failed=0
    skipped=0

    for line in "${lines[@]}"; do
        if [[ "$line" == *'FAILED'* ]]; then
            ((failed++))
        elif [[ "$line" == *'SKIPPED'* ]]; then
            ((skipped++))
        fi
    done

    passed=$((total - failed - skipped))

    {
        printf '\n### Summary\n\n'
        printf -- '- **Total test suites:** %d\n' "$total"
        printf -- '- **Passed:** %d\n' "$passed"
        printf -- '- **Failed:** %d\n' "$failed"
        printf -- '- **Skipped:** %d\n' "$skipped"
        printf '\n'

        if [[ $failed -gt 0 ]]; then
            printf -- '**Status:** Tests failed - review failed test suites above\n'
        elif [[ $passed -eq 0 ]]; then
            printf -- '**Status:** No tests executed\n'
        else
            printf -- '**Status:** All executed tests passed\n'
        fi
    } >>"$file"
}

main() {
    if (($# < 2)); then
        printf 'Usage: %s <Category> <Test1=status> [Test2=status ...]\n' "$0" >&2
        printf 'Example: %s "Action" "Basic Tests=success" "Advanced Tests=failure"\n' "$0" >&2
        exit 1
    fi

    local category="$1"
    shift
    local output="./${category// /-}-results.md"

    write_header "$output" "$category"
    process_tests "$output" "$@"
    write_summary "$output"

    local status_text
    if $all_passed_flag; then
        status_text='PASSED'
    else
        status_text='FAILED'
    fi

    # Update header with final status
    sed -i "1s/.*/## ${category} Test Results - ${status_text}/" "$output"

    printf 'Test summary written to: %s\n' "$output"
}

main "$@"
exit 0
