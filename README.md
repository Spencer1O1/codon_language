# Codon Language Manifesto
## What Codon Is

**Codon is a domain-specific language (DSL) for describing software architecture.**

A Codon genome defines the **semantic structure of a system** — domains, modules, models, capabilities, and relationships — independently of any implementation language or framework.

Codon does **not** describe code directly.
Instead, Codon describes architecture, which is then **expressed** into implementation artifacts.

The genome is the source of **truth** for system architecture.

---

# Core Philosophy

Codon is built on several guiding principles:
- **Architecture first** — system structure is defined before implementation.
- **Semantics over files** — the genome describes meaning, not directories.
- **Structure before detail** — fields and types can be inferred from higher-level structure.
- **Explicit dependencies** — architectural boundaries must be declared.
- **Deterministic expression** — the same genome and expression rules must produce the same output.
- **Evolution through mutation** — architecture evolves through explicit patches.

---

# The Genome Model
Codon organizes architecture using biological metaphors.

| Concept | Meaning |
|---------|---------|
| Genome | complete architecture of a system |
| Chromosome | architectural domain or bounded context |
| Gene | module-level unit within a chromosome |
| Codon	| architectural instruction (entity, capability, relation, etc.) |
| Nucleotide | primitive data (strings, flags, field types, etc.) |

This hierarchy defines the **semantic structure** of the architecture.

---

# Core Codon Types

Genes may contain the following architectural codons:

### Entities

Persistent domain models owned by a gene.

```
entities:
  User:
    fields:
      email: string
```

---

### Capabilities

Semantic domain actions.

```
capabilities:
  - create-user
  - authenticate-user
```

Capabilities describe **what the system does**, not how it is implemented.

---

### Relations

Structural links between entities within the same gene.

```
relations:
  - from: Comment
    to: Issue
    type: many-to-one
    name: issue
```

---

### References

Cross-gene links to entities owned by other genes.

```
references:
  - from: Issue
    to: identity.user.User
```

References require explicit dependencies.

---

### Traits

Reusable architectural programs that expand into genome structure.

Traits can introduce:
- chromosomes
- genes
- entities
- capabilities
- relations
- references

Traits allow common architecture patterns to be reused.

---

# Addressing Model
Codon uses a canonical addressing system.
```
chromosome.gene.entity
```
Example:
```
identity.user.User
```
This addressing model guarantees deterministic reference resolution.

---

# Dependencies

Genes must declare dependencies on other genes they reference.
```
dependencies:
  - identity.user
```
Rules:
- dependencies must reference existing genes
- cross-gene references require a declared dependency
- circular dependencies are invalid unless explicitly allowed

Dependencies enforce **architectural boundaries**.

---

Genome Composition

A genome is composed from:
```
genome.yaml
chromosomes/
```
Example structure:
```
genome.yaml
chromosomes/
  identity/
    user.yaml
  tracking/
    issues.yaml
```
The **composed genome** is the semantic model formed by combining the manifest with all gene files.

---

# Traits and Expansion

Traits are **architectural macros** that expand into genome structure.

Two attachment scopes exist:

| Scope |	Behavior |
|-------|----------|
| project |	may introduce chromosomes or genes |
| gene | expands only within the owning gene |

Trait expansion must be deterministic.

---

# Expression

The genome itself does **not** describe implementation structure.

Instead, **expression rules** project architecture into code.
```
genome → expression → implementation
```
Examples of expression artifacts:
- services
- controllers
- repositories
- API routes
- database schemas

Expression rules are defined in the `expression/` directory.

---

# Mutation

Genomes evolve through **explicit mutations**.

Mutations are applied using structured patches.

AI proposes mutation
User approves mutation
Codon applies mutation
Git records evolution

Mutations guarantee that architectural evolution is:
- explicit
- reviewable
- reproducible

---

# Validation

The genome must satisfy several rule classes:

### Structural validation
- schema compliance
- allowed codon types
- valid file structure

### Addressing validation
- reference resolution
- identifier uniqueness

### Architectural validation
- dependency correctness
- relation validity
- trait expansion consistency

Validation ensures the genome remains a coherent architecture model.

---

# Codon Is a Language
Codon is not just configuration.

It defines:
- **grammar** — valid structures
- **semantics** — architectural meaning
- **resolution rules** — addressing and references
- **transformation rules** — traits and mutation
- **projection rules** — expression into code

The YAML files are only the **syntax**.

The language itself is the **Codon genome model**.

---

# Long-Term Goal

The ultimate goal of Codon is:

> A deterministic architectural language from which complete software systems can be derived, evolved, validated, and regenerated.

Codon genomes are intended to be:
-readable by humans
-analyzable by machines
-evolvable over time
- expressive enough to describe complex architectures


---

# Layers

```
Codon Language
│
├── syntax (YAML structure)
├── grammar (allowed constructs)
├── semantics (meaning of constructs)
├── addressing (name resolution)
├── transformation (traits + mutation)
└── projection (expression → code)
```