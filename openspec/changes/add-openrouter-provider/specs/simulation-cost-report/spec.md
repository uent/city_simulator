## ADDED Requirements

### Requirement: Cost accumulator
The `internal/llm` package SHALL provide a `CostAccumulator` struct that aggregates `Usage` values across multiple LLM calls. It SHALL be safe for concurrent use.

Methods:
- `Add(u Usage)` — add usage from one call to the running totals
- `Total() Usage` — return the cumulative totals

#### Scenario: Accumulation across multiple calls
- **WHEN** `Add` is called three times with distinct `Usage` values
- **THEN** `Total()` SHALL return the sum of all prompt tokens, completion tokens, and estimated cost

#### Scenario: Zero total when no calls made
- **WHEN** a new `CostAccumulator` is created and `Total()` is called before any `Add`
- **THEN** all fields of the returned `Usage` SHALL be zero

#### Scenario: Concurrent Add is safe
- **WHEN** multiple goroutines call `Add` concurrently
- **THEN** no data race SHALL occur and `Total()` SHALL reflect all additions

---

### Requirement: Cost report rendered in simulation summary
The `internal/summary` package's `GenerateSummary` function SHALL accept a `*llm.CostAccumulator` parameter. When `accumulator.Total().PromptTokens > 0` or `accumulator.Total().CompletionTokens > 0`, a `## Cost Report` section SHALL be appended to the summary after the character cards block.

The cost report SHALL include:
- Prompt tokens (integer)
- Completion tokens (integer)
- Total tokens (prompt + completion)
- Estimated cost in USD formatted to 6 decimal places (e.g. `$0.000123`)

When all totals are zero (Ollama or no LLM calls), the cost report section SHALL NOT be appended.

#### Scenario: Cost report appended for OpenRouter simulation
- **WHEN** `GenerateSummary` is called with an accumulator whose totals are non-zero
- **THEN** the returned summary string SHALL contain a `## Cost Report` section with token counts and formatted USD cost

#### Scenario: Cost report omitted for Ollama simulation
- **WHEN** `GenerateSummary` is called with an accumulator whose all totals are zero
- **THEN** the returned summary string SHALL NOT contain any `## Cost Report` section

#### Scenario: Cost report omitted when accumulator is nil
- **WHEN** `GenerateSummary` is called with a nil accumulator pointer
- **THEN** the function SHALL not panic and SHALL produce a summary without a cost report section

#### Scenario: Cost formatted to 6 decimal places
- **WHEN** the accumulated estimated cost is `0.000123456`
- **THEN** the cost report SHALL display `$0.000123` (truncated/rounded to 6 decimal places)
