# docs/ — LLM-Managed Wiki

This directory is an LLM-maintained knowledge base for the go-masker project.

## Rules for AI Agents

### When to update docs

- Adding a new masker type → update `masker-types.md`
- Changing public API (new method, changed signature, removed type) → update `api.md`
- Structural changes (new package, new file, changed module) → update `architecture.md`
- Design decisions or RFC-style proposals → add to `design/`

### How to update

- Write in English (code-facing docs) — keep it concise
- One topic per file, use the existing file if it covers the topic
- Tables over prose for reference material
- Include runnable Go examples where they clarify usage
- No duplicating README.md content — link to it instead
- Keep each file self-contained: a reader should not need to read other files to understand one file

### File naming

- Lowercase kebab-case: `masker-types.md`, `api.md`
- Design docs: `design/<topic>.md`
- No date prefixes — use git history for chronology

### What NOT to put here

- Content that belongs in GoDoc (per-function docs → put in source code)
- Changelog (use GitHub Releases)
- Tutorials or getting-started (that's README.md)

## Directory Structure

```
docs/
├── AGENTS.md            ← this file (governance)
├── architecture.md      ← project structure, module layout, key types
├── masker-types.md      ← reference for all masker types with behavior details
├── api.md               ← public API surface reference
└── design/              ← design docs, RFCs, proposals
    └── v3-design.md     ← v3 major version design
```
