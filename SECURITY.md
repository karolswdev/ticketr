# Security Policy

## Supported Versions

Ticketr is currently in pre-release (v0.x). Security updates are provided for the latest version only.

| Version | Supported          |
| ------- | ------------------ |
| 0.x.x   | :white_check_mark: |
| < 0.1.0 | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue in Ticketr, please follow responsible disclosure practices:

### How to Report

**DO NOT** open a public GitHub issue for security vulnerabilities.

Instead, please email security reports to:

- **Email:** karolswdev@gmail.com
- **Subject:** [SECURITY] Ticketr - Brief Description

### What to Include

Please include the following information in your report:

1. **Description** - A clear description of the vulnerability
2. **Impact** - What an attacker could achieve by exploiting this
3. **Steps to Reproduce** - Detailed steps to reproduce the issue
4. **Affected Versions** - Which versions of Ticketr are affected
5. **Suggested Fix** - If you have a recommendation (optional)
6. **Your Contact Info** - So we can follow up with questions

### What to Expect

- **Acknowledgment:** We will acknowledge receipt of your report within 48 hours
- **Updates:** We will provide regular updates on our progress
- **Timeline:** We aim to release a fix within 30 days for critical issues
- **Credit:** We will credit you in the security advisory (unless you prefer to remain anonymous)

### Disclosure Policy

- We request that you do not publicly disclose the vulnerability until we have released a fix
- We will coordinate with you on the disclosure timeline
- Once fixed, we will publish a security advisory crediting the reporter (if desired)

## Security Considerations for Users

### JIRA API Tokens

Ticketr requires JIRA API credentials to function. Please follow these best practices:

1. **Never commit credentials to version control**
   - Use `.env` file (already in `.gitignore`)
   - Use environment variables in CI/CD
   - Never hardcode tokens in configuration files

2. **Rotate tokens regularly**
   - Generate new API tokens periodically
   - Revoke old tokens after rotation
   - Use different tokens for different environments

3. **Minimum permissions**
   - Grant only the permissions needed for Ticketr to function
   - Use project-specific tokens when possible
   - Avoid using admin-level API tokens

4. **Secure storage**
   - Protect `.env` file with appropriate file permissions (`chmod 600 .env`)
   - Use secrets management in CI/CD (GitHub Secrets, etc.)
   - Never share tokens via unsecured channels

### Execution Logs

Ticketr automatically logs operations to `.ticketr/logs/`. These logs:

- Automatically redact sensitive data (API keys, emails, passwords)
- Should be added to `.gitignore` (already included)
- May contain ticket content and metadata
- Are rotated automatically (last 10 files kept)

If sharing logs for debugging:
- Review logs before sharing to ensure no sensitive data is present
- Redact any additional sensitive information manually if needed

### Local State Files

The `.ticketr.state` file tracks ticket hashes and should be:

- Added to `.gitignore` (environment-specific)
- Not shared between environments
- Treated as sensitive if it contains ticket identifiers

## Known Security Limitations

As a pre-1.0 project, Ticketr has the following known limitations:

1. **API Token Storage:** Credentials stored in plaintext `.env` file (user's responsibility to secure)
2. **Network Security:** No built-in HTTPS certificate validation override protection
3. **Rate Limiting:** No built-in protection against API rate limiting (relies on JIRA's limits)

We are actively working on improving security for the 1.0 release.

## Security Best Practices for Development

If you are contributing to Ticketr:

1. Never commit test credentials or API tokens
2. Use `.env.example` as template, never commit actual `.env`
3. Review code for potential injection vulnerabilities
4. Validate all user input appropriately
5. Follow secure coding practices for Go
6. Run `go vet` and `staticcheck` before submitting PRs

## Additional Resources

- [JIRA API Token Management](https://id.atlassian.com/manage-profile/security/api-tokens)
- [GitHub Security Best Practices](https://docs.github.com/en/code-security)
- [OWASP Secure Coding Practices](https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/)

---

**Last Updated:** October 2025
