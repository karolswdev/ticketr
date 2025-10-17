# Support

Thank you for using Ticketr! This document explains how to get help with Ticketr.

## 📚 Documentation

Before seeking support, please check our comprehensive documentation:

- **[README](README.md)** - Overview, features, and quick start
- **[Quick Start Guide](https://github.com/karolswdev/ticktr/wiki/Quick-Start-Tutorial)** - Get started in 5 minutes
- **[Documentation](docs/)** - Complete documentation suite
  - [WORKFLOW.md](docs/WORKFLOW.md) - End-to-end workflows and examples
  - [Migration Guide](docs/migration-guide.md) - Migrating from legacy format
  - [State Management](docs/state-management.md) - Understanding .ticketr.state
  - [Release Process](docs/release-process.md) - Version and release information
- **[Examples](examples/)** - Ready-to-use templates and examples
- **[FAQ](https://github.com/karolswdev/ticktr/wiki/FAQ)** - Frequently asked questions

## 🐛 Bug Reports

If you've found a bug, please [create an issue](https://github.com/karolswdev/ticktr/issues/new) with:

- **Ticketr version** (`ticketr --version`)
- **Operating system** (Linux, macOS, Windows)
- **Go version** (`go version`)
- **Steps to reproduce** the problem
- **Expected behavior** vs **actual behavior**
- **Relevant logs** (from `.ticketr/logs/`)
- **Minimal example** (if applicable)

**Important:** Review logs before sharing to ensure no sensitive data (API tokens, etc.) is included. Ticketr automatically redacts credentials, but please double-check.

## 💡 Feature Requests

Have an idea for a new feature? We'd love to hear it!

1. Check [existing feature requests](https://github.com/karolswdev/ticketr/issues?q=is%3Aissue+label%3Aenhancement) to avoid duplicates
2. [Create a feature request](https://github.com/karolswdev/ticktr/issues/new) describing:
   - **Use case**: What problem does this solve?
   - **Proposed solution**: How should it work?
   - **Alternatives**: Have you considered other approaches?
   - **Impact**: Who would benefit from this feature?

## ❓ Questions & Discussions

For general questions, usage help, or discussions:

- **[GitHub Discussions](https://github.com/karolswdev/ticktr/discussions)** - Ask questions, share tips, discuss ideas
- **[Wiki](https://github.com/karolswdev/ticktr/wiki)** - Community-driven documentation

### Good Discussion Topics
- How to implement specific workflows
- Best practices for team usage
- Integration with other tools
- Performance optimization tips
- Use case sharing

## 🔒 Security Vulnerabilities

**DO NOT** report security vulnerabilities via public GitHub issues.

Please follow our [Security Policy](SECURITY.md):
- Email: **karolswdev@gmail.com**
- Subject: `[SECURITY] Ticketr - Brief Description`
- Include: Description, impact, steps to reproduce, affected versions

We aim to:
- Acknowledge reports within **48 hours**
- Provide updates on progress
- Release fixes within **30 days** for critical issues

## 📞 Response Expectations

### Issues & Bug Reports
- **Acknowledgment:** Within 48 hours
- **Initial Response:** Within 7 days
- **Resolution Time:** Varies by severity
  - Critical bugs: Days to weeks
  - Regular bugs: Weeks to months
  - Enhancements: Triaged for future releases

### Discussions
- Best-effort community support
- Maintainer responses when available
- Community members encouraged to help

### Pull Requests
See [CONTRIBUTING.md](CONTRIBUTING.md) for the PR process.

## 🌍 Community Guidelines

All interactions must follow our [Code of Conduct](CODE_OF_CONDUCT.md).

In summary:
- Be respectful and inclusive
- Provide constructive feedback
- Focus on the issue, not the person
- Help others when you can

## 🛠️ Self-Help Resources

### Common Issues

**Authentication Errors:**
```bash
# Check your credentials
echo $JIRA_URL
echo $JIRA_EMAIL
# API key should be set but not echoed

# Test authentication
ticketr schema
```

**Field Not Found Errors:**
```bash
# Discover available fields
ticketr schema

# Check field names (case-sensitive!)
```

**State Issues:**
```bash
# Reset state to force full push
rm .ticketr.state
ticketr push tickets.md
```

**No Changes Detected:**
```bash
# Ticketr tracks changes via .ticketr.state
# If you need to force a push:
rm .ticketr.state

# Or make a real edit to the ticket
```

For more troubleshooting, see our [GitHub Discussions](https://github.com/karolswdev/ticktr/discussions).

## 📖 Additional Resources

- **[CHANGELOG](CHANGELOG.md)** - Version history and release notes
- **[ROADMAP](ROADMAP.md)** - Future plans and milestones
- **[ARCHITECTURE](ARCHITECTURE.md)** - Technical architecture details
- **[CONTRIBUTING](CONTRIBUTING.md)** - How to contribute

## 💼 Commercial Support

Ticketr is a community-driven open source project. Commercial support, training, or custom development is not currently available.

For enterprise needs:
- Consider contributing features you need
- Sponsor development via GitHub Sponsors (if available)
- Hire contributors for custom work (see git log for contacts)

## 🙏 Contributing

The best way to get support is to help improve Ticketr!

- Fix bugs you encounter
- Improve documentation
- Help others in Discussions
- Review pull requests
- Share your use cases

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📬 Contact

- **Issues:** https://github.com/karolswdev/ticktr/issues
- **Discussions:** https://github.com/karolswdev/ticktr/discussions
- **Security:** karolswdev@gmail.com (see [SECURITY.md](SECURITY.md))
- **Maintainer:** [@karolswdev](https://github.com/karolswdev)

---

**Thank you for being part of the Ticketr community!** 🎉

We appreciate your patience and understanding. Ticketr is maintained by volunteers in their spare time.
