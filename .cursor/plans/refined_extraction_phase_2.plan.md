# # Implementation Plan - Refinements

Based on the analysis of real-world API data (MusicBrainz), we are implementing several refinements to make the generated Go code more idiomatic and less fragmented.

1.  **Micro-Embedding Prevention ([extract.go](extract.go))**:

    - Add a heuristic to prevent extracting a new shared type if it only adds a single field to an already extracted base.
    - This reduces the depth of the embedding hierarchy and prevents "over-DRYing" the code at the cost of readability.

2.  **Auto-Extract Extended Structs ([extract.go](extract.go))**:

    - Automatically promote any anonymous struct that contains an embedded field to a named top-level struct.
    - This ensures that if we are using a trait hierarchy, the "leaf" types are also named and clean.

3.  **Context-Aware Naming ([extract.go](extract.go) / [names.go](names.go))**:

    - Update the naming logic to consider the JSON keys that refer to a common structure.
    - If multiple references use similar names (e.g., `area`, `begin-area`), use that as a strong hint for the type name (e.g., `Area`).

4.  **Solid Testing Suite**:

    - **Micro-Embedding**: Test that adding 1 field to a base doesn't create a 3rd type.
    - **Auto-Promotion**: Test that `[]struct { Base; Extra }` becomes `[]ExtendedType`.
    - **Naming Hints**: Test that shared structures get names from their parent fields when applicable.

## Insights & Design Considerations

### Hierarchy Depth vs. Complexity

While deeply nested embeddings are DRY, they can be hard to follow. We should prioritize "Breadth" (naming meaningful entities) over "Depth" (naming every possible subset).