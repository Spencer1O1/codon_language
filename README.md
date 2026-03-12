# Codon Language (spec + engine)
Work in progress! Contribution welcome, contact @ spencerliamsmith@gmail.com

Codon is a DSL  for describing software architecture in an implementation‑agnostic way. It’s layered:
- **Genome** – the whole spec for a system.
- **Chromosomes** – bounded contexts/domains.
- **Genes** – modules within a chromosome.
- **Codons** – architectural units (entities, capabilities, relations, traits, docs, implementation rules, expression docs, etc.).
Optionally, a genome includes **expression** assets (targets/projections/styles/templates) that guide an **expressor** to generate concrete project structure (services, APIs, UI, infra) deterministically after validation.

This repo contains:
- The self‑hosted Codon language specification (written in Codon).
- A Go loader/validator/emitter CLI to parse, validate, and emit composed genomes.

**Future use:** once projection engines are plugged in, any validated genome + expression assets can generate runnable project scaffolds (API + UI + infra), with traits enabling reusable architecture patterns.

## What’s in the repo

- `.codon/` – the language genome (spec):
  - `genome.yaml` – manifest
  - `chromosomes/` – genes documenting language constructs (entities, capabilities, relations, traits, expression docs, tooling, etc.)
  - `codon_schemas/` – codon schema definitions
  - `nucleotides/types/` – primitive & engine nucleotypes
  - `traits/` – trait specs
  - `expression/` – (optional) projection assets live next to chromosomes; their shapes are documented under `chromosomes/expression/`.
- `pkg/` – Go implementation of loader, validator, nucleotype parser.
- `cmd/codon/` – CLI entrypoint (`load`, `validate`, `emit`).
- `fixtures/` – test genomes grouped by domain (language/traits/manifest/expression) and a full example.

## Workflow

Validate a genome (runs loader + validator):
```bash
make validate ROOT=./fixtures/example        # or any genome root
# or directly:
GOCACHE=... GOMODCACHE=... go run ./cmd/codon validate ./fixtures/example
```

Emit the composed genome (post-validation) as YAML:
```bash
make emit ROOT=./fixtures/example
```
The composed artifact includes manifest, codon_schemas (exported), nucleotypes (exported), chromosomes/genes with codons, traits_applied, issues, and optional `expression` (targets/projections/styles/templates).

Run tests:
```bash
make test
```

## Expression assets (projection inputs)

Expression files live at `expression/` (sibling to `chromosomes/`). Shapes are enforced and documented in `chromosomes/expression/*.yaml`:
- `targets.yaml` – map target_name → {kind, stack, output_root?, overwrite?, …}
- `projections.yaml` – map projection_name → {target, binding, capabilities/entities/relations selectors}
- `styles.yaml` – map style_name → {version?, …}
- `templates.yaml` – map template_name → {source, checksum?, variables?, postprocess?}

Loader parses these optionally; validator enforces shapes and required fields; composed genome mirrors them (no nesting).

## Minimal example (ticketing slice)

```
genome.yaml
chromosomes/
  ticketing/
    tickets.yaml
expression/
  targets.yaml
  projections.yaml
```

`genome.yaml`
```yaml
schema_version: 1.0.0
project:
  name: ticketing-example
  type: service
```

`chromosomes/ticketing/tickets.yaml`
```yaml
gene: tickets
description: Ticketing domain
codons:
  entities:
    ticket:
      id: uuid
      reporter_id: uuid
      status: Union<"open","closed">
  capabilities:
    create_ticket:
      effects: ["create_ticket"]
      inputs:
        reporter_id: { type: uuid }
      outputs:
        ticket_id: { type: uuid }
  relations:
    reporter_to_ticket:
      from: ticket
      to: ticket   # self for example; normally another entity
      type: one-to-one
```

`expression/targets.yaml`
```yaml
targets:
  api:
    kind: api_service
    stack: go-hex-rest@1.0
    output_root: services/api
```

`expression/projections.yaml`
```yaml
projections:
  api_routes:
    target: api
    capabilities: ["*"]
    binding: rest_default
```

Run:
```bash
make validate ROOT=./fixtures/example   # or your genome root
make emit ROOT=./fixtures/example
```

## Addressing & refs

Refs resolve by shortest path: `field` (same gene), `gene.field` (same chromosome), `chromosome.gene.field` (cross-chromosome). Over‑qualification warns; missing targets error.

## Traits

Traits live under `traits/` (genome/chromosome/gene scopes). Trait-local codon_schemas/nucleotypes are loaded before injection. Trait merge policies and conflict rules are in `chromosomes/language/traits.yaml` and enforced by the validator.

## Contributing / testing changes

1. Edit the spec in `.codon/` or engine code under `pkg/`.
2. Run `make test`.
3. Optionally emit the example: `make emit ROOT=./fixtures/example` to inspect composed output.

## Status

- Loader: parses manifest, schemas, nucleotypes, traits, expression; emits loader issues.
- Validator: grouped rule sets (manifest/language/traits/expression).
- Expression: shapes and field typing enforced; capability coverage reported as info; uniqueness enforced.
- Emission: composed genome includes expression section and is closed/typed per `chromosomes/tooling/composed_genome.yaml`.
