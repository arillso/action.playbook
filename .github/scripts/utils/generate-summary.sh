#!/usr/bin/env bash
#
# Description: Generic summary generator for GitHub Actions workflows
# Usage:       ./generate-summary.sh --title "Title" --output "output.md" [Item1=status] [Item2=status ...]
# Exit Code:   Always exits with 0 to ensure workflow continuation
#

set -euo pipefail
IFS=$'\n\t'

declare -g all_passed_flag=true
declare -g output=""
declare -g title=""
declare -g subtitle=""

# Process command line options
while [[ $# -gt 0 ]]; do
    case "$1" in
    --title)
        title="$2"
        shift 2
        ;;
    --subtitle)
        subtitle="$2"
        shift 2
        ;;
    --output)
        output="$2"
        shift 2
        ;;
    *)
        break
        ;;
    esac
done

# Validate required parameters
if [[ -z "$title" ]]; then
    echo "Error: --title is required" >&2
    exit 1
fi

if [[ -z "$output" ]]; then
    echo "Error: --output is required" >&2
    exit 1
fi

if [[ $# -lt 1 ]]; then
    echo "Error: At least one status item required" >&2
    exit 1
fi

write_header() {
    printf -- '## %s\n\n' "$title" >"$output"

    if [[ -n "$subtitle" ]]; then
        printf -- '%s\n\n' "$subtitle" >>"$output"
    fi
}

process_items() {
    local all_passed=true

    printf -- '| Item | Status |\n| ---- | ------ |\n' >>"$output"

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
        elif [[ "$status" == "skipped" ]]; then
            symbol='⏭️'
        else
            symbol='❌'
            all_passed=false
        fi

        printf -- '| %s | %s |\n' "$name" "$symbol" >>"$output"
    done

    all_passed_flag=$all_passed
}

write_summary() {
    local pattern total failed passed skipped
    pattern='^(##|\| Item \| Status \||\| ---- \| ------ \||$)'
    mapfile -t lines < <(grep -vE "$pattern" "$output")
    total=${#lines[@]}
    failed=0
    skipped=0
    for line in "${lines[@]}"; do
        [[ "$line" == *'❌'* ]] && ((failed++))
        [[ "$line" == *'⏭️'* ]] && ((skipped++))
    done
    passed=$((total - failed - skipped))

    {
        printf -- '\n### Statistical Summary\n\n'
        printf -- '- **Total items:** %d\n' "$total"
        printf -- '- **Passed:**      %d ✅\n' "$passed"

        if ((skipped > 0)); then
            printf -- '- **Skipped:**     %d ⏭️\n' "$skipped"
        fi

        printf -- '- **Failed:**      %d ❌\n' "$failed"
        printf -- '\n### Environment Information\n\n'
        printf -- '- **Run Date:** %s\n' "$(date -u)"
        printf -- '- **GitHub Run ID:** %s\n' "${GITHUB_RUN_ID:-Unknown}"
    } >>"$output"
}

update_title_status() {
    local status_symbol status_text
    if $all_passed_flag; then
        status_symbol='✅'
        status_text='SUCCESS'
    else
        status_symbol='❌'
        status_text='FAILED'
    fi

    sed -i "1s/.*/## ${title} — ${status_symbol} ${status_text}/" "$output"
}

main() {
    write_header
    process_items "$@"
    write_summary
    update_title_status

    echo "Summary generated at $output"
}

main "$@"
exit 0
