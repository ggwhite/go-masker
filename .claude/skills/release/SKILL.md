---
name: release
description: "Release a new version of go-masker: update CHANGELOG.md from git log, create GitHub release. Use when the user says 「發布」「release」「發新版」「打 tag」「release new version」or similar. Also use when the user says 「更新 changelog」if they clearly intend to release afterward."
---

# Release

Publish a new go-masker version. Uses `gh release create` to publish, so the job here is: write accurate release notes, get user approval, then tag and push.

## Workflow

### 1. Determine version

```bash
git tag -l 'v*' --sort=-v:refname | head -1
```

Bump the patch number by default (e.g., v2.2.1 -> v2.2.2). If the user specifies `minor` or `major`, bump accordingly.

Version rules (from CLAUDE.md):
- **patch** `x.y.Z`: bug fix, docs, internal changes that don't affect API
- **minor** `x.Y.0`: new feature, new masker type, backwards-compatible API addition
- **major** `X.0.0`: breaking changes (removed API, changed behavior, module path change)

Show the user: "上個版本是 vX.Y.Z，這次準備發 vA.B.C" and confirm before continuing.

### 2. Collect changes

```bash
git log <last-tag>..HEAD --oneline --no-decorate
```

If there are zero commits since the last tag, stop and tell the user there's nothing to release.

### 3. Draft changelog section

Read the existing `CHANGELOG.md` to match its style. The format follows [Keep a Changelog](https://keepachangelog.com/):

```markdown
## [X.Y.Z] - YYYY-MM-DD

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

Skip commits that are only changelog/release housekeeping (e.g., "docs: vX.Y.Z changelog").

### 4. Preview

Show the drafted changelog section to the user in full. Ask: "這樣 OK 嗎？要調整什麼？"

Do NOT proceed until the user approves. If the user requests changes, revise and show again.

### 5. Update CHANGELOG.md

Insert the new version section **at the top**, right after the file header (the "# Changelog" heading and the format description paragraph). The existing sections must not be modified — they correspond to already-released tags.

### 6. Commit

```bash
git add CHANGELOG.md
git commit -m "docs: vX.Y.Z changelog"
```

### 7. Tag and push

Ask the user before pushing: "要 push 了，確認？"

```bash
git push && git tag vX.Y.Z && git push --tags
```

### 8. Create GitHub release

```bash
gh release create vX.Y.Z --title "Release vX.Y.Z" --notes-file -
```

Pipe the changelog section content (just this version's notes, not the full CHANGELOG.md) as the release notes.

After release, tell the user: "Release vX.Y.Z 已發布" and provide the release URL.

## Important rules

- **Never modify existing version sections** — they match published releases
- **Date format** — YYYY-MM-DD, use today's date (timezone: Asia/Taipei)
- **Commit message** — always `docs: vX.Y.Z changelog` (no emoji, no scope parentheses)
- **One version per run** — if the user wants to skip versions, they say so explicitly
- **v3 sub-module** — if the release includes v3/ changes, remind the user that v3 has its own module path (`github.com/ggwhite/go-masker/v3`) and may need a separate tag (`v3/vX.Y.Z`) when v3 is officially released
