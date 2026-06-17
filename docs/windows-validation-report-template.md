# Windows Validation Report Template

Use this template when running the real Windows validation flow from `docs/windows-validation.md`.

If you do not have a real Windows machine, leave this template alone and rely on CI plus `bash scripts/windows_toolchain_smoke.sh` instead.
If you do have a real Windows machine, fill out only the commands and checks you actually ran.

## Summary

- Date:
- Validator:
- Result: pass / fail / blocked
- Scope: MinGW-w64 / host load / editor reopen / save-restore

## Environment

- Windows version:
- CPU architecture:
- Go version:
- MinGW-w64 version:
- WebView2 SDK version or source:
- WebView2 runtime version:
- Host application name:
- Host application version:

## Commands

Record the command output only for the steps you ran on the real Windows machine.

### `go test -timeout=30000s ./...`

- Result:
- Notes:

```text
paste output here
```

### `just windows-preflight`

- Result:
- Notes:

```text
paste output here
```

### `just windows-build-dll`

- Result:
- Output DLL path:
- Generated header path:
- Notes:

```text
paste output here
```

### `just windows-build dist/windows/vst3go.dll`

- Result:
- Bundle path:
- Notes:

```text
paste output here
```

### `just windows-check-bundle`

- Result:
- Notes:

```text
paste output here
```

## Bundle Verification

- DLL present in `Contents/x86_64-win`: yes / no
- Generated header sidecar present: yes / no
- `WebView2Loader.dll` present: yes / no
- Bundle layout matches expectation: yes / no

## Host Validation

- Plugin discovered by host: yes / no
- Plugin inserted on track: yes / no
- Audio processing starts: yes / no
- Editor opens: yes / no
- Embedded web UI renders: yes / no
- UI parameter changes reach processor: yes / no
- Host automation reflects in UI: yes / no
- Reopen after close works: yes / no
- Multiple instances work: yes / no
- Project save and reopen restores state: yes / no
- Rescan or reload avoids file-lock issues: yes / no

## Evidence

- Screenshot paths:
- Video paths:
- Crash log paths:
- Additional notes:

Only attach evidence for failures or the host behaviors you needed to confirm manually.

## Issues

### Issue 1

- Title:
- Repro steps:
- Expected:
- Actual:
- Severity:
- Attachments:

### Issue 2

- Title:
- Repro steps:
- Expected:
- Actual:
- Severity:
- Attachments:

## Final Decision

- MinGW-w64 validation complete: yes / no
- Ready to keep Windows in release scope: yes / no
- Follow-up tasks:

If Windows validation is still blocked by missing hardware or a missing host, mark that clearly here so the next person does not have to rediscover it.
