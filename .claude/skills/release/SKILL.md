---
name: release
description: "Release a new version of go-masker: update CHANGELOG.md from git log, create GitHub release. Use when the user says 「發布」「release」「發新版」「打 tag」「release new version」or similar. Also use when the user says 「更新 changelog」if they clearly intend to release afterward."
---

# Release

Publish a new go-masker version. v2 和 v3 共存於同一 repo（v3 在 `v3/` 子目錄），各自獨立版號和 tag。

## Module 結構

```
go-masker/
├── go.mod        ← module github.com/ggwhite/go-masker/v2
├── masker.go     ← v2 code
├── v3/
│   ├── go.mod    ← module github.com/ggwhite/go-masker/v3
│   └── masker.go ← v3 code
│   └── zapfield/
│       └── go.mod ← module github.com/ggwhite/go-masker/v3/zapfield
```

| Module | Tag 格式 | 範例 |
|--------|----------|------|
| v2 | `vX.Y.Z` | `v2.4.0` |
| v3 | `v3/vX.Y.Z` | `v3/v3.0.0` |
| zapfield | `v3/zapfield/vX.Y.Z` | `v3/zapfield/v0.1.0` |

## Workflow

### 0. Determine which module to release

Ask the user if not clear: "要發 v2 還是 v3？"

- If user says "release v2" / "發 v2" → release v2
- If user says "release v3" / "發 v3" → release v3
- If user says "release" without specifying, check which module has changes since its last tag and ask

Set variables for the rest of the workflow:

| | v2 | v3 | zapfield |
|---|---|---|---|
| `TAG_PREFIX` | (empty) | `v3/` | `v3/zapfield/` |
| `TAG_PATTERN` | `v2.*` | `v3/v3.*` | `v3/zapfield/v*` |
| `CODE_DIR` | `.` | `v3` | `v3/zapfield` |
| `CHANGELOG_PREFIX` | (empty) | `v3-` | `zapfield-` |

### 1. Determine version

```bash
git tag -l '${TAG_PATTERN}' --sort=-v:refname | head -1
```

Bump the patch number by default. If the user specifies `minor` or `major`, bump accordingly.

Version rules (from CLAUDE.md):
- **patch** `x.y.Z`: bug fix, docs, internal changes that don't affect API
- **minor** `x.Y.0`: new feature, new masker type, backwards-compatible API addition
- **major** `X.0.0`: breaking changes (removed API, changed behavior, module path change)

Show the user: "上個版本是 ${TAG_PREFIX}vX.Y.Z，這次準備發 ${TAG_PREFIX}vA.B.C" and confirm before continuing.

### 2. Collect changes

```bash
# v2: changes to root-level Go files only
git log <last-tag>..HEAD --oneline --no-decorate -- '*.go' ':!v3/'

# v3: changes under v3/ only
git log <last-tag>..HEAD --oneline --no-decorate -- 'v3/' ':!v3/zapfield/'

# zapfield: changes under v3/zapfield/ only
git log <last-tag>..HEAD --oneline --no-decorate -- 'v3/zapfield/'
```

If there are zero relevant commits since the last tag, stop and tell the user there's nothing to release.

### 3. Draft changelog section

Read the existing `CHANGELOG.md` to match its style. The format follows [Keep a Changelog](https://keepachangelog.com/):

```markdown
## [v2 X.Y.Z] - YYYY-MM-DD       ← v2 用 [X.Y.Z]
## [v3 X.Y.Z] - YYYY-MM-DD       ← v3 用 [v3 X.Y.Z]
## [zapfield X.Y.Z] - YYYY-MM-DD  ← zapfield 用 [zapfield X.Y.Z]

### Features

- **Short name** — description of the feature

### Fixes

- **Short name** — description of the fix
```

Categorization rules:
- `fix` / `bugfix` commits -> `### Fixes`
- `feat` / feature commits -> `### Features`
- `docs` commits -> `### Docs` (only if user-facing; skip changelog-only docs commits)
- `chore` / `refactor` / `perf` -> `### Internal` (only if noteworthy; skip trivial ones)

Each item should be a concise, user-facing description — not a raw commit message. Rewrite commit messages into the bold-name-em-dash style. Group related commits into a single item when they address the same thing.

Skip commits that are only changelog/release housekeeping.

### 4. Preview

Show the drafted changelog section to the user in full. Ask: "這樣 OK 嗎？要調整什麼？"

Do NOT proceed until the user approves.

### 5. Update CHANGELOG.md

Insert the new version section at the appropriate position:
- v2 sections and v3 sections interleave by date (newest first)
- If there's an `[Unreleased]` section covering this module's changes, replace those items with the versioned section

### 6. Commit

```bash
git add CHANGELOG.md
git commit -m "docs: ${TAG_PREFIX}vX.Y.Z changelog"
```

### 7. Test

Run tests for the target module before tagging:

```bash
# v2
go test -race ./...

# v3
cd v3 && go test -race ./...

# zapfield
cd v3/zapfield && go test -race ./...
```

If tests fail, stop and tell the user.

### 8. Tag and push

Ask the user before pushing: "要 push 了，確認？"

```bash
git push && git tag ${TAG_PREFIX}vX.Y.Z && git push --tags
```

### 9. Create GitHub release

```bash
gh release create ${TAG_PREFIX}vX.Y.Z --title "Release ${TAG_PREFIX}vX.Y.Z" --notes-file -
```

Pipe the changelog section content (just this version's notes) as the release notes.

After release, tell the user: "Release ${TAG_PREFIX}vX.Y.Z 已發布" and provide the release URL.

### 10. Cross-module reminder

After releasing one module, check if the other module also has unreleased changes:

```bash
# After v2 release, check v3
git log $(git tag -l 'v3/v*' --sort=-v:refname | head -1)..HEAD --oneline -- 'v3/' | head -5

# After v3 release, check v2
git log $(git tag -l 'v2.*' --sort=-v:refname | head -1)..HEAD --oneline -- '*.go' ':!v3/' | head -5
```

If there are changes, ask: "v{other} 也有未發布的變更，要一起發嗎？"

## Important rules

- **Never modify existing version sections** — they match published releases
- **Date format** — YYYY-MM-DD, use today's date (timezone: Asia/Taipei)
- **Commit message** — always `docs: ${TAG_PREFIX}vX.Y.Z changelog`
- **One module per run** — release v2 and v3 in separate runs (but prompt for the other)
- **v3 tag 必須帶 `v3/` 前綴** — `v3/v3.0.0` 不是 `v3.0.0`，Go modules 靠前綴找子目錄
- **zapfield tag 必須帶 `v3/zapfield/` 前綴** — 獨立 module 獨立版號
- **v3 依賴 zapfield 時用 replace** — 開發期 `v3/zapfield/go.mod` 可能有 `replace` 指令，release 前要確認移除或更新為正式版號
