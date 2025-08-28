diff --git a/README.md b/README.md
index 0337b79..fefe1a3 100644
--- a/README.md
+++ b/README.md
@@ -150,6 +150,24 @@ ticketr schema -v
 
 This is especially useful when working with custom fields that vary between JIRA instances.
 
+### State Management
+
+Ticketr automatically tracks changes to prevent redundant updates to JIRA:
+
+```bash
+# The .ticketr.state file is created automatically
+# It stores SHA256 hashes of ticket content
+
+# Only changed tickets are pushed to JIRA
+ticketr push stories.md  # Skips unchanged tickets
+
+# The state file contains:
+# - Ticket ID to content hash mappings
+# - Automatically updated after each successful push
+```
+
+**Note**: The `.ticketr.state` file should be added to `.gitignore` as it's environment-specific.
+
 ### Docker Usage
 
 Build and run using Docker:
