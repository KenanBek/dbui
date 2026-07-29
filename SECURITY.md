# Security Policy

## Supported versions

Only the latest release receives fixes.

## Reporting a vulnerability

Use GitHub's private vulnerability reporting:
[Report a vulnerability](https://github.com/KenanBek/dbui/security/advisories/new).
Please do not open public issues for security reports.

dbui connects to databases with credentials you supply; treat your `dbui.yml`
as a secret file until env/keychain indirection ships (planned for v1.0).
