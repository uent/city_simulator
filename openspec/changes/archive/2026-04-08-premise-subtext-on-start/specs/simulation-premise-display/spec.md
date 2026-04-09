## ADDED Requirements

### Requirement: Print world concept at simulation start
The system SHALL print the world concept block to stdout at the start of `Engine.Run`, before the first tick executes, if `Scenario.World.Concept.Premise` is non-empty.

The block SHALL follow this format:
```
=== World Concept ===
Premise: <premise>
Flavor:  <flavor>
Rules:
  - <rule>
=====================
```

The `Flavor:` line SHALL be omitted if `Concept.Flavor` is empty. The `Rules:` section SHALL be omitted if `Concept.Rules` is empty. If `Concept.Premise` is empty, the entire block SHALL be skipped silently.

#### Scenario: Full concept block printed
- **WHEN** `Engine.Run` is called and `Scenario.World.Concept.Premise` is `"Bears hiding among humans"`
- **THEN** stdout SHALL contain a block starting with `=== World Concept ===` and including the premise line before the first tick output

#### Scenario: Flavor line present when set
- **WHEN** `Concept.Flavor` is `"absurdist heist comedy"`
- **THEN** stdout SHALL contain `Flavor:  absurdist heist comedy` inside the concept block

#### Scenario: Rules section present when set
- **WHEN** `Concept.Rules` contains two entries
- **THEN** stdout SHALL contain `Rules:` followed by two `  - <rule>` lines

#### Scenario: Flavor line omitted when empty
- **WHEN** `Concept.Flavor` is empty string
- **THEN** stdout SHALL NOT contain a `Flavor:` line in the concept block

#### Scenario: Rules section omitted when empty
- **WHEN** `Concept.Rules` is nil or empty
- **THEN** stdout SHALL NOT contain a `Rules:` line in the concept block

#### Scenario: Block skipped when premise is empty
- **WHEN** `Concept.Premise` is empty string
- **THEN** no concept block SHALL be printed and no error SHALL occur
