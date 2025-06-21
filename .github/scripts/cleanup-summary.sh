#!/usr/bin/env bash
#
# Description: Generate a markdown summary of container registry cleanup operations
# Usage:       ./cleanup-summary.sh <Registry1=status> [Registry2=status ...]
# Exit Code:   Always exits with 0 for workflow continuation
#

set -euo pipefail

# Determine script directory for proper relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
UTILS_DIR="${SCRIPT_DIR}/utils"

# Add custom footer information about retention policies
add_policy_info() {
    local output_file="$1"

    {
        printf -- '\n### Retention Policies Applied\n\n'
        printf -- '- Commit hash tags (.*-[0-9a-f]{7}): 60 days\n'
        printf -- '- Untagged images (<none>): 60 days\n'
        printf -- '- Skipped tags: latest, main, version tags (vX.Y.Z)\n'
        printf -- '\n### Next Scheduled Run\n\n'
        printf -- '- Next cleanup scheduled for the first day of next month at 02:00 UTC\n'
    } >>"$output_file"
}

# Execute the shared summary generator with cleanup-specific parameters
"${UTILS_DIR}/generate-summary.sh" \
    --title "Container Registry Cleanup Report" \
    --subtitle "Cleanup executed on: $(date -u)" \
    --output "./cleanup-report.md" \
    "$@"

# Add the registry-specific policy information
add_policy_info "./cleanup-report.md"

echo "Cleanup summary generated at ./cleanup-report.md"
exit 0
