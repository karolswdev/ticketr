# Archived Migration Guides

This directory contains legacy migration guides for historical reference. These guides are preserved for users upgrading from older Ticketr versions.

## Contents

### v1.x to v2.0 Migration
- **[migration-guide.md](migration-guide.md)** - Guide for migrating from legacy `# STORY:` format to canonical `# TICKET:` format (v1.x → v2.0)

### v2.x to v3.0 Migration
- **[v3-MIGRATION-GUIDE.md](v3-MIGRATION-GUIDE.md)** - Comprehensive guide for PathResolver migration from local paths to platform-standard global directories (v2.x → v3.0)
- **[v3-migration-guide.md](v3-migration-guide.md)** - Alternative migration guide (duplicate, retained for completeness)
- **[phase-2-workspace-migration.md](phase-2-workspace-migration.md)** - Workspace model migration from environment variables to multi-workspace SQLite database

## Archive Status

**Current Version:** v3.1.1

These migration guides are **archived** because:
1. Ticketr v3.1.1 is production-ready with no migration layers
2. Users should be on v3.x by now
3. Migration tools (PROD-010) remain available in CLI for legacy format conversion

## For Current Users

If you're running Ticketr v3.1.1, you don't need these guides unless:
- You're maintaining a legacy v2.x installation
- You need to understand historical migration paths
- You're troubleshooting an incomplete migration

## Migration Tool Reference

The CLI migration commands are still available:

```bash
# Migrate legacy STORY format to TICKET format
ticketr migrate <file> --write

# Migrate v2.x local paths to v3.x global paths
ticketr migrate-paths

# Rollback to v2.x local paths
ticketr rollback-paths
```

See [REQUIREMENTS.md](../../REQUIREMENTS.md) PROD-010 for migration tool specification.

## Support

For migration issues:
1. Check the specific guide for your upgrade path
2. Review [TROUBLESHOOTING.md](../TROUBLESHOOTING.md) for common issues
3. Open an issue at https://github.com/karolswdev/ticketr/issues

---

**Archive Date:** 2025-10-20 (Phase 6 Week 1 Day 1)
**Archived By:** Scribe Agent
**Reason:** Production-ready v3.1.1 release cleanup
