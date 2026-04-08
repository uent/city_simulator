## 1. Update summary prompt

- [x] 1.1 Raise `maxEvents` constant in `internal/summary/summary.go` from 100 to 200
- [x] 1.2 Update system prompt in `buildPrompt` to request at least six paragraphs instead of two to four
- [x] 1.3 Add world concept block to user message (premise, flavor, rules) — omit if all fields are empty
- [x] 1.4 Add world atmosphere and weather to user message — omit if empty
- [x] 1.5 Extend character block to include motivation, fear, and goals per character — omit missing fields silently

## 2. Update spec

- [x] 2.1 Update `openspec/specs/simulation-summary/spec.md` to reflect new event cap (200), richer prompt content, and six-paragraph narrative requirement
