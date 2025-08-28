diff --git a/README.md b/README.md
index 8cb827f..0337b79 100644
--- a/README.md
+++ b/README.md
@@ -113,18 +113,43 @@ Simply edit your file and run the tool again - it intelligently handles updates:
 
 ```bash
 # Basic operation
-ticketr -f stories.md
+ticketr push stories.md
 
 # Verbose output for debugging
-ticketr -f stories.md --verbose
+ticketr push stories.md --verbose
 
 # Continue on errors (CI/CD mode)
-ticketr -f stories.md --force-partial-upload
+ticketr push stories.md --force-partial-upload
+
+# Discover JIRA schema and generate configuration
+ticketr schema > .ticketr.yaml
 
-# Combine options
+# Legacy mode (backward compatibility)
 ticketr -f stories.md -v --force-partial-upload
 ```
 
+### Schema Discovery
+
+The `ticketr schema` command helps you discover available fields in your JIRA instance and generate a proper configuration file:
+
+```bash
+# Discover fields and generate configuration
+ticketr schema > .ticketr.yaml
+
+# View available fields with verbose output
+ticketr schema -v
+
+# The command will output field mappings like:
+# field_mappings:
+#   "Story Points":
+#     id: "customfield_10010"
+#     type: "number"
+#   "Sprint": "customfield_10020"
+#   "Epic Link": "customfield_10014"
+```
+
+This is especially useful when working with custom fields that vary between JIRA instances.
+
 ### Docker Usage
 
 Build and run using Docker:
