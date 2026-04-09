## 1. Implement character card renderer

- [x] 1.1 Add `renderCharacterCards(chars []*character.Character) string` function in `internal/summary/summary.go` that skips game directors and returns empty string when no eligible characters exist
- [x] 1.2 Render each character card as a `### <Name>` Markdown section with labeled fields, omitting zero/empty/nil values
- [x] 1.3 Include Voice sub-fields (Formality, Verbal Tics, Response Length, Humor Type, Communication Style) under a "Voice" sub-heading when any are non-empty
- [x] 1.4 Include Relational Defaults (Strangers, Authority, Vulnerable) under a "Relational Defaults" sub-heading when any are non-empty
- [x] 1.5 Include Cover Identity (Alias, Role, Backstory, Weaknesses) under a "Cover Identity" sub-heading when `CoverIdentity` is non-nil
- [x] 1.6 Include Formative Events as a bullet list when non-empty
- [x] 1.7 Include Dialogue Examples as a bullet list when non-empty
- [x] 1.8 Include Relationships (Judgments) as a sub-section with trust/interest/threat levels and impression per known character

## 2. Wire cards into GenerateSummary

- [x] 2.1 In `GenerateSummary`, after receiving the LLM narrative text, call `renderCharacterCards(chars)` and append the result to the narrative string before returning

## 3. Verify output

- [x] 3.1 Run a simulation (e.g. `default` scenario) and confirm the saved `.md` file and terminal output contain the `## Character Cards` section after the narrative
- [x] 3.2 Confirm game director characters do not appear in the cards section
- [x] 3.3 Confirm characters with sparse attributes (missing fear, cover identity, etc.) render cleanly with only populated fields shown
