# Windows Docs Index

Use this page as the entry point for Windows build, smoke, and real-host validation.

## Start Here

- [`docs/windows-build.md`](windows-build.md): build harness and bundle shape
- [`docs/windows-validation.md`](windows-validation.md): real Windows host checklist
- [`docs/windows-validation-quickstart.md`](windows-validation-quickstart.md): shortest real-host path
- [`docs/windows-validation-report-template.md`](windows-validation-report-template.md): fill-in report template

## Recommended Reading Order

1. read `docs/windows-build.md`
2. read `docs/windows-validation.md`
3. use `docs/windows-validation-quickstart.md` when you are on a real Windows machine
4. fill `docs/windows-validation-report-template.md` only for real host runs

## Quick Links

- `bash scripts/windows_toolchain_smoke.sh`: local smoke coverage on non-Windows hosts
- `just windows-smoke`: repo-level smoke wrapper
- `just windows-preflight`: Windows compiler/header/loader check
- `just windows-validate`: full real-host validation

If you are comparing Windows guidance with the EQ8 flow, start with `docs/eq8-index.md` and then return here for the platform-specific pieces.
