#!/bin/bash

# Check available issue types in JIRA project
echo "Checking issue types for project: $JIRA_PROJECT_KEY"
echo ""

# Get project details including issue types
curl -s \
  -u "$JIRA_EMAIL:$JIRA_API_KEY" \
  -H "Accept: application/json" \
  "$JIRA_URL/rest/api/2/project/$JIRA_PROJECT_KEY" | \
  python3 -m json.tool | \
  grep -A2 '"issueTypes"' | \
  grep '"name"' | \
  cut -d'"' -f4

echo ""
echo "If you need more details, run:"
echo "curl -u \"\$JIRA_EMAIL:\$JIRA_API_KEY\" -H \"Accept: application/json\" \"\$JIRA_URL/rest/api/2/project/\$JIRA_PROJECT_KEY\" | python3 -m json.tool"