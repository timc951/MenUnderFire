#!/bin/bash
set -euo pipefail

# Parse security scan results from Trivy and Semgrep JSON files into a single CSV
# Usage: ./tools/parse-scan-results.sh

SCAN_DIR="$(cd "$(dirname "$0")/.." && pwd)/securityScan"
OUTPUT="$SCAN_DIR/security_issues.csv"

# Write CSV header
echo 'Scanner,Source,Severity,CVE,CWE,Package,Installed Version,Fixed Version,Title,Description' > "$OUTPUT"

# Helper: find the latest file matching a prefix
latest_file() {
  local prefix="$1"
  ls -t "$SCAN_DIR"/${prefix}*.json 2>/dev/null | head -1
}

# Helper: escape a field for CSV (double quotes inside, wrap in quotes)
csv_escape() {
  local val="$1"
  # Replace double quotes with two double quotes, then wrap in quotes
  val="${val//\"/\"\"}"
  printf '"%s"' "$val"
}

# Helper: write a CSV row
write_row() {
  local scanner="$1" source="$2" severity="$3" cve="$4" cwe="$5" package="$6" installed="$7" fixed="$8" title="$9" description="${10}"
  printf '%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n' \
    "$(csv_escape "$scanner")" \
    "$(csv_escape "$source")" \
    "$(csv_escape "$severity")" \
    "$(csv_escape "$cve")" \
    "$(csv_escape "$cwe")" \
    "$(csv_escape "$package")" \
    "$(csv_escape "$installed")" \
    "$(csv_escape "$fixed")" \
    "$(csv_escape "$title")" \
    "$(csv_escape "$description")" >> "$OUTPUT"
}

# Parse Trivy JSON files
parse_trivy() {
  local file="$1"
  local source="$2"

  if [[ ! -f "$file" ]]; then
    echo "Warning: $file not found, skipping"
    return
  fi

  echo "Parsing Trivy: $source ($file)"

  # Use jq to extract vulnerabilities
  jq -r '
    .Results[]? |
    select(.Vulnerabilities != null) |
    .Vulnerabilities[] |
    [
      .VulnerabilityID // "",
      .PkgName // "",
      .InstalledVersion // "",
      .FixedVersion // "",
      .Severity // "",
      .Title // "",
      (.Description // "" | gsub("\n"; " ") | gsub("\r"; ""))
    ] | @tsv
  ' "$file" | while IFS=$'\t' read -r vuln_id pkg_name installed fixed severity title description; do
    write_row "Trivy" "$source" "$severity" "$vuln_id" "" "$pkg_name" "$installed" "$fixed" "$title" "$description"
  done
}

# Parse Semgrep JSON files
parse_semgrep() {
  local file="$1"
  local source="$2"

  if [[ ! -f "$file" ]]; then
    echo "Warning: $file not found, skipping"
    return
  fi

  echo "Parsing Semgrep: $source ($file)"

  jq -r '
    .results[]? |
    [
      .check_id // "",
      .path // "",
      (.start.line // "" | tostring),
      .extra.severity // "",
      (.extra.metadata.cwe // [] | if type == "array" then join("; ") else . end),
      (.extra.message // "" | gsub("\n"; " ") | gsub("\r"; ""))
    ] | @tsv
  ' "$file" | while IFS=$'\t' read -r check_id path line severity cwes message; do
    local location="${path}:${line}"
    write_row "Semgrep" "$source" "$severity" "" "$cwes" "$location" "" "" "$check_id" "$message"
  done
}

# Verify jq is available
if ! command -v jq &> /dev/null; then
  echo "Error: jq is required. Install: sudo apt-get install jq"
  exit 1
fi

echo "=== Parsing scan results ==="
echo ""

# Parse latest Trivy scans
parse_trivy "$(latest_file frontend_)" "Frontend"
parse_trivy "$(latest_file backend_)" "Backend"
parse_trivy "$(latest_file docker_app_)" "Docker - App"
parse_trivy "$(latest_file docker_postgres_)" "Docker - Postgres"

# Parse latest Semgrep scans
parse_semgrep "$(latest_file semgrep_frontend_)" "Frontend"
parse_semgrep "$(latest_file semgrep_backend_)" "Backend"
parse_semgrep "$(latest_file semgrep_docker_)" "Docker"

# Sort by severity (keep header, sort the rest)
# Severity order: CRITICAL > HIGH > ERROR > MEDIUM > WARNING > LOW > INFO
TEMP="$SCAN_DIR/.sort_temp.csv"
head -1 "$OUTPUT" > "$TEMP"
tail -n +2 "$OUTPUT" | awk -F',' '{
    sev = $3
    gsub(/"/, "", sev)
    if (sev == "CRITICAL") order = 1
    else if (sev == "HIGH") order = 2
    else if (sev == "ERROR") order = 3
    else if (sev == "MEDIUM") order = 4
    else if (sev == "WARNING") order = 5
    else if (sev == "LOW") order = 6
    else order = 7
    print order "|||" $0
  }' | sort -t'|' -k1,1n | sed 's/^[0-9]|||//' >> "$TEMP"

mv "$TEMP" "$OUTPUT"

TOTAL=$(tail -n +2 "$OUTPUT" | wc -l)
echo ""
echo "=== Done ==="
echo "Total issues: $TOTAL"
echo "Output: $OUTPUT"
