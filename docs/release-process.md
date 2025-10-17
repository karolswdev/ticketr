# Release Process

This document describes the release process for Ticketr, including versioning guidelines, release checklist, and automation details.

## Table of Contents

- [Versioning](#versioning)
- [Release Types](#release-types)
- [Pre-Release Checklist](#pre-release-checklist)
- [Release Steps](#release-steps)
- [Post-Release Tasks](#post-release-tasks)
- [Automation](#automation)
- [Emergency Releases](#emergency-releases)
- [Rollback Procedures](#rollback-procedures)

## Versioning

Ticketr follows [Semantic Versioning 2.0.0](https://semver.org/):

```
MAJOR.MINOR.PATCH
```

### Version Increment Rules

**MAJOR version** (X.0.0):
- Breaking changes to CLI interface
- Incompatible API changes
- Breaking changes to Markdown schema
- Removal of features

**MINOR version** (0.X.0):
- New features (backward compatible)
- New CLI commands
- Enhanced functionality
- Deprecations (with backward compatibility)

**PATCH version** (0.0.X):
- Bug fixes
- Security patches
- Documentation updates
- Performance improvements (no behavior change)

### Pre-1.0 Policy

During pre-1.0 development (0.x.x):
- Breaking changes MAY occur in MINOR versions
- API is not considered stable
- Use caution in production environments

### Version 1.0.0 Criteria

Release 1.0.0 when:
- [ ] All ROADMAP.md milestones complete
- [ ] Test coverage ≥ 70%
- [ ] Security audit passed
- [ ] Complete documentation
- [ ] Field testing (≥3 months)
- [ ] Stable public API defined

## Release Types

### Regular Release

Scheduled releases for new features and improvements:
- Frequency: Every 2-4 weeks (during active development)
- Branch: `main`
- Process: Full release checklist

### Patch Release

Urgent bug fixes and security patches:
- Frequency: As needed
- Branch: `main` or hotfix branch
- Process: Expedited release (subset of checklist)

### Pre-release / Release Candidate

Testing releases before major versions:
- Format: `v1.0.0-rc.1`, `v1.0.0-beta.1`
- Branch: `release/v1.0.0`
- Process: Full checklist + extended testing

## Pre-Release Checklist

### 1. Code Quality

- [ ] All tests passing: `go test ./...`
- [ ] Quality checks passing: `bash scripts/quality.sh`
- [ ] Smoke tests passing: `bash tests/smoke/smoke_test.sh`
- [ ] No unresolved high-priority issues
- [ ] Code review completed for all changes
- [ ] Static analysis clean: `staticcheck ./...`
- [ ] Formatting verified: `gofmt -l .`

### 2. Documentation

- [ ] CHANGELOG.md updated with release notes
- [ ] README.md reflects new features
- [ ] All new features documented
- [ ] Migration guide updated (if breaking changes)
- [ ] API documentation current
- [ ] Example files updated

### 3. Testing

- [ ] Unit tests passing (103/103 expected)
- [ ] Integration tests passing
- [ ] Smoke tests passing (7/7 scenarios)
- [ ] Manual testing completed
- [ ] Cross-platform testing (Linux, macOS, Windows)
- [ ] Docker image tested

### 4. Dependencies

- [ ] `go.mod` and `go.sum` are tidy
- [ ] Dependency vulnerabilities checked: `go list -m all | nancy sleuth` (if using nancy)
- [ ] License compliance verified
- [ ] No deprecated dependencies

### 5. Security

- [ ] Security audit completed (for MAJOR/MINOR releases)
- [ ] No credentials in repository
- [ ] Secrets properly redacted in logs
- [ ] SECURITY.md is current

## Release Steps

### Step 1: Prepare Release Branch

```bash
# Ensure you're on main and up to date
git checkout main
git pull origin main

# For major/minor releases, create release branch (optional)
git checkout -b release/v0.3.0
```

### Step 2: Update Version Information

Update version references in:
- `CHANGELOG.md` - Add release date and finalize notes
- Any version constants in code (if applicable)

```bash
# Update CHANGELOG.md
vim CHANGELOG.md

# Example: Change [Unreleased] section to [0.3.0] - 2025-10-16
```

### Step 3: Run Pre-Release Checks

```bash
# Run all quality checks
bash scripts/quality.sh

# Run smoke tests
bash tests/smoke/smoke_test.sh

# Run full test suite
go test ./... -v -race -coverprofile=coverage.out

# Check coverage
go tool cover -func=coverage.out | tail -1
```

### Step 4: Commit Release Changes

```bash
# Commit version bump
git add CHANGELOG.md
git commit -m "chore(release): Prepare v0.3.0

- Update CHANGELOG with release notes
- Finalize version documentation

Generated with [Claude Code](https://claude.com/claude-code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>"

# Push changes
git push origin main  # or release branch
```

### Step 5: Create and Push Tag

```bash
# Create annotated tag
git tag -a v0.3.0 -m "Release v0.3.0

Major Features:
- Repository hygiene and release readiness
- Automated release workflow
- Enhanced security practices

See CHANGELOG.md for complete details."

# Push tag (this triggers the release workflow)
git push origin v0.3.0
```

### Step 6: Monitor Release Workflow

```bash
# Watch GitHub Actions workflow
# Visit: https://github.com/karolswdev/ticktr/actions

# The workflow will:
# 1. Run tests
# 2. Build binaries for all platforms
# 3. Create release archives
# 4. Generate checksums
# 5. Extract release notes from CHANGELOG
# 6. Create GitHub release
# 7. Build and publish Docker image
```

### Step 7: Verify Release

After workflow completes:

1. **Check GitHub Release**:
   - Visit: https://github.com/karolswdev/ticktr/releases
   - Verify release notes are correct
   - Confirm all artifacts present:
     - `ticketr-VERSION-linux-amd64.tar.gz`
     - `ticketr-VERSION-linux-arm64.tar.gz`
     - `ticketr-VERSION-darwin-amd64.tar.gz`
     - `ticketr-VERSION-darwin-arm64.tar.gz`
     - `ticketr-VERSION-windows-amd64.zip`
     - `ticketr-VERSION-docker.tar.gz`
     - `checksums.txt`

2. **Test Binary Downloads**:
   ```bash
   # Download and test a binary
   wget https://github.com/karolswdev/ticktr/releases/download/v0.3.0/ticketr-0.3.0-linux-amd64.tar.gz
   tar -xzf ticketr-0.3.0-linux-amd64.tar.gz
   ./ticketr-linux-amd64 --version
   ```

3. **Test Docker Image**:
   ```bash
   # If pushed to Docker Hub
   docker pull username/ticketr:0.3.0
   docker run --rm username/ticketr:0.3.0 --version

   # Or load from release artifact
   wget https://github.com/karolswdev/ticktr/releases/download/v0.3.0/ticketr-0.3.0-docker.tar.gz
   docker load < ticketr-0.3.0-docker.tar.gz
   docker run --rm ticketr:0.3.0 --version
   ```

4. **Verify Checksums**:
   ```bash
   # Download checksums
   wget https://github.com/karolswdev/ticktr/releases/download/v0.3.0/checksums.txt

   # Verify (on Linux/macOS)
   sha256sum -c checksums.txt
   ```

## Post-Release Tasks

### 1. Update Documentation

- [ ] Update README badges if needed
- [ ] Announce release in GitHub Discussions
- [ ] Update project website (if applicable)
- [ ] Update package managers (Homebrew, etc.) - future

### 2. Communication

- [ ] Post release announcement
- [ ] Update social media (if applicable)
- [ ] Notify major users (if breaking changes)
- [ ] Update issue tracker (close resolved issues)

### 3. Prepare for Next Release

```bash
# Update CHANGELOG for next version
vim CHANGELOG.md

# Add new [Unreleased] section at top:
## [Unreleased]

### Added

### Changed

### Fixed

### Removed
```

### 4. Monitor Issues

- Watch for bug reports related to new release
- Monitor GitHub Issues and Discussions
- Check CI/CD for any issues
- Review logs if telemetry available

## Automation

### GitHub Actions Workflows

#### CI Workflow (`.github/workflows/ci.yml`)

Runs on:
- Push to `main` or `feat/**` branches
- Pull requests to `main`

Jobs:
- **Build**: Multi-OS, multi-Go version matrix
- **Test**: Unit tests with race detector
- **Coverage**: Minimum 50% threshold
- **Lint**: go vet, gofmt, staticcheck
- **Smoke Tests**: 7 end-to-end scenarios

#### Release Workflow (`.github/workflows/release.yml`)

Triggers on:
- Git tags matching `v*.*.*`

Jobs:
1. **Build and Release**:
   - Verify CHANGELOG entry exists
   - Run full test suite
   - Build binaries for all platforms:
     - Linux (amd64, arm64)
     - macOS (amd64, arm64)
     - Windows (amd64)
   - Create release archives (.tar.gz, .zip)
   - Generate SHA256 checksums
   - Extract release notes from CHANGELOG
   - Create GitHub release
   - Mark pre-1.0 releases as "pre-release"

2. **Docker Release**:
   - Build Docker image
   - Tag with version and latest
   - Push to Docker Hub (if credentials configured)
   - Save image as release artifact

### Required Secrets

For full automation, configure these GitHub Secrets:

- `GITHUB_TOKEN` - Automatically provided
- `DOCKER_USERNAME` - Docker Hub username (optional)
- `DOCKER_PASSWORD` - Docker Hub access token (optional)

Docker publishing is optional - workflow skips if credentials not configured.

## Emergency Releases

For critical security patches or severe bugs:

### Fast-Track Process

1. **Create Hotfix Branch**:
   ```bash
   git checkout -b hotfix/v0.2.1 v0.2.0
   ```

2. **Apply Fix**:
   ```bash
   # Make minimal changes
   # Add tests
   git commit -m "fix(security): Patch CVE-2025-XXXXX"
   ```

3. **Fast-Track Testing**:
   ```bash
   # Run essential tests only
   go test ./...
   bash tests/smoke/smoke_test.sh
   ```

4. **Update CHANGELOG**:
   ```markdown
   ## [0.2.1] - 2025-10-17

   ### Security
   - Fix CVE-2025-XXXXX credential exposure
   ```

5. **Release**:
   ```bash
   git tag -a v0.2.1 -m "Emergency patch for CVE-2025-XXXXX"
   git push origin v0.2.1
   ```

6. **Backport to Main**:
   ```bash
   git checkout main
   git cherry-pick <hotfix-commit>
   git push origin main
   ```

### Communication

For security releases:
- [ ] Publish security advisory on GitHub
- [ ] Update SECURITY.md if needed
- [ ] Notify users directly if severe
- [ ] Coordinate with security researchers

## Rollback Procedures

If a release has critical issues:

### Option 1: Quick Patch

Release a new patch version immediately (see Emergency Releases).

### Option 2: Mark Release as Deprecated

```bash
# Edit GitHub release
# - Check "This is a pre-release"
# - Add warning to release notes
# - Point users to previous stable version
```

### Option 3: Delete Release (Last Resort)

```bash
# Delete GitHub release (UI or API)
# Delete git tag
git push --delete origin v0.3.0
git tag -d v0.3.0

# Note: Users may have already downloaded
```

### Option 4: Yanked Release

Add to CHANGELOG.md:

```markdown
## [0.3.0] - 2025-10-16 [YANKED]

**This release has been yanked due to critical bug #123.**
**Please use v0.2.1 instead.**
```

## Release Checklist Summary

Print this checklist for each release:

```
Pre-Release:
[ ] All tests passing
[ ] Quality checks passing
[ ] CHANGELOG.md updated
[ ] Documentation current
[ ] Dependencies checked
[ ] Security reviewed

Release:
[ ] Version bumped in code (if needed)
[ ] CHANGELOG date finalized
[ ] Changes committed
[ ] Tag created and pushed
[ ] GitHub Actions completed
[ ] Release artifacts verified
[ ] Binaries tested
[ ] Docker image tested

Post-Release:
[ ] Announcement posted
[ ] Issues updated
[ ] Next version prepared
[ ] Monitoring active
```

## Troubleshooting

### Release Workflow Failed

1. Check GitHub Actions logs
2. Common issues:
   - Missing CHANGELOG entry
   - Test failures
   - Build errors
   - Docker credentials (non-blocking)

### Binary Won't Run

1. Check architecture match
2. Verify checksums
3. Check file permissions: `chmod +x ticketr-*`
4. Review error messages

### Docker Image Issues

1. Verify tag pushed correctly
2. Check Docker Hub credentials
3. Load local image for testing
4. Review Dockerfile changes

## References

- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [goreleaser](https://goreleaser.com/) - Future consideration for advanced release automation

---

**Last Updated**: October 2025
**Maintainer**: Ticketr Team
