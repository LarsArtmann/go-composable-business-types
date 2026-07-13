# Domain Language

A **Unified Language** for go-composable-business-types — shared across developers, consumers, and AI.
Inspired by Domain-Driven Design (DDD) Ubiquitous Language.

Every term below should mean the **same thing** to everyone who reads it.

## Glossary

| Term          | Definition                                                          | Context                         |
| ------------- | ------------------------------------------------------------------- | ------------------------------- |
| DataPoint     | Self-contained data unit with complete audit trail (who, when, why) | Core abstraction of the library |
| Actor         | Entity that caused a data event (User, Bot, Service, System)        | Audit trail, authorization      |
| Bitemporal    | Time tracking with separate valid-time and transaction-time         | Temporal correctness            |
| NanoID        | URL-safe, cryptographically random identifier (21 chars default)    | Entity identification           |
| BoundedString | String validated against min/max length constraints                 | Input validation                |
| Importance    | Priority classification (0-100) with named levels                   | Task/project prioritization     |

## Value Objects

Immutable objects defined by attributes.

| Term       | Definition                                                   | Context                            |
| ---------- | ------------------------------------------------------------ | ---------------------------------- |
| Email      | A validated email address                                    | User identity, contact             |
| URL        | A validated URL with scheme and host                         | Web references                     |
| Percentage | A 0-100 value with float conversion (clamps overflow to 100) | Rates, proportions                 |
| Cents      | Monetary amount in smallest currency unit (no float errors)  | Payments, pricing                  |
| Timestamp  | Domain-wrapped time.Time for business timestamps             | Temporal tracking                  |
| Duration   | Domain-wrapped time.Duration for business intervals          | Time spans                         |
| Money      | ISO 4217 currency amount with formatting                     | International payments             |
| Locale     | BCP 47 language tag for internationalization                 | i18n, formatting                   |
| Tag        | Validated string label with alphanumeric+hyphen constraint   | Categorization, metadata           |
| Address    | Validated postal address (line1, city, postalCode, country)  | Shipping, billing                  |
| Contact    | Contact info (name, email, phone, website, address)          | Person/organization representation |

## Entities

Objects with identity and lifecycle.

| Term        | Definition                                                                                                 | Context                     |
| ----------- | ---------------------------------------------------------------------------------------------------------- | --------------------------- |
| DataPoint   | Wraps any payload with full metadata: actor, temporal, trigger, context, references, causes, tags, version | Event sourcing, audit trail |
| ProjectCore | Composite project metadata (name, path, languages, importance, tags)                                       | Project ecosystem           |

## Concepts

| Term       | Definition                                                       | Context                    |
| ---------- | ---------------------------------------------------------------- | -------------------------- |
| ActorChain | Ordered chain of actors showing delegation path                  | Audit trail, authorization |
| Context    | Execution environment (environment, session, request, source)    | DataPoint metadata         |
| Reference  | Type-safe reference to another entity with relationship metadata | DataPoint lineage          |
| Cause      | Entry in a causal chain for building lineage graphs              | DataPoint lineage          |
| Correction | Flag indicating a bitemporal record is a correction              | Temporal tracking          |
| Validator  | Interface for self-validating types                              | Input validation pattern   |

## Enums

| Enum      | Values                                                            | Context                     |
| --------- | ----------------------------------------------------------------- | --------------------------- |
| ActorKind | User, Bot, System, Service                                        | Who triggered an action     |
| Priority  | Low, Medium, High, Critical                                       | Task urgency                |
| Status    | Draft, Active, Paused, Archived, Deleted                          | Entity lifecycle            |
| Trigger   | Manual, Scheduled, Webhook, Import, Migration, System, Correction | What caused a DataPoint     |
| CauseKind | Direct, Command, Event                                            | Causal chain classification |

---

> **How to use this file:**
>
> - Keep terms concise — one clear sentence per definition
> - Update when new domain concepts emerge
> - Use these terms consistently in code, docs, and conversations
> - When in doubt about a word's meaning, check here first
