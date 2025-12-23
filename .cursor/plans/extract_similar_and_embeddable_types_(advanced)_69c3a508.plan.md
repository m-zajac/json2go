---
name: Extract Similar and Embeddable Types (Advanced)
overview: Implement advanced type extraction with deep type-aware similarity, transitive embedding, and robust optionality propagation to ensure valid and idiomatic Go code.
todos:
  - id: impl-deep-similarity
    content: Implement deep similarity (name + type) in node.go
    status: completed
  - id: impl-ancestor-similarity
    content: Add ancestor-aware similarity for better recursion detection
    status: completed
  - id: impl-iterative-extract
    content: Implement iterative multi-pass extraction in extract.go
    status: completed
  - id: impl-permissive-merge
    content: Implement 'most permissive' field merging logic
    status: completed
  - id: impl-name-cap
    content: Cap naming length and improve name generation in names.go
    status: completed
  - id: test-adv-logic
    content: Test optionality propagation and transitive embedding
    status: completed
---

# # Implementation Plan

1.  **Deep Similarity Metric ([node.go](node.go))**:

    - Implement `similarity` that accounts for both key names and their type IDs (e.g., `name.string` vs `name.int`).
    - Add ancestor-based similarity check to prioritize recursion detection.

2.  **Iterative Extraction ([extract.go](extract.go))**:

    - Perform multi-pass extraction to support transitive embedding (e.g., `A` embeds `B` which embeds `C`).
    - Implement "Most Permissive" field merging: if a field is optional in any source, it becomes optional in the extracted struct.

3.  **Refined Naming ([names.go](names.go))**:

    - Update naming logic to cap field-based names at 3 fields and handle collisions gracefully using `nextName`.

4.  **AST Safety ([ast.go](ast.go))**:

    - Ensure embedded fields handle JSON tag promotion correctly (avoiding shadowed fields).

5.  **Solid Testing Suite**:

    - **Deep Similarity**: Test that `{"a": 1}` and `{"a": "s"}` are NOT merged.
    - **Optionality**: Test that a mix of required/optional fields results in an optional field in the extracted type.
    - **Transitive**: Test a 3-level hierarchy (`Person` -> `Employee` -> `Manager`).
    - **Recursion**: Test a complex recursive JSON (e.g., a file system tree).

## Insights & Design Considerations

### Deep vs. Shallow Similarity

Similarity must be type-aware. In Go, `{"id": 1}` and `{"id": "abc"}` are fundamentally different types. The metric must incorporate the `nodeType` of each field, not just the name, to prevent invalid merges.

### Most Permissive Merging (Pointer Propagation)

When merging multiple nodes into a shared type (or extracting a common subset), a field must become a pointer (`*Type` with `omitempty`) if it is optional or nullable in **any** of the source structures. This ensures the generated type is compatible with all inputs.

### JSON Tag Shadowing

In Go's `encoding/json`, a field in an outer struct shadows a field with the same JSON tag in an embedded struct. The extraction logic must ensure that extracted subsets do not contain fields that would be shadowed by the parent struct, or handle these conflicts explicitly to avoid data loss during unmarshaling.

### Naming Constraints