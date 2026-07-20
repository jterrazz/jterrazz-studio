---
name: jterrazz-structured-data
description: JSON-LD structured data (schema.org) — choosing types, required properties, validation, rich results. Use when the user mentions "schema", "structured data", "JSON-LD", "rich results", "rich snippets", or when adding machine-readable markup to pages. For the broader SEO audit see jterrazz-seo; for AI citation formatting see jterrazz-geo.
---

# Structured data — JSON-LD

Schema tells engines *what an entity is*, unambiguously. Wrong or dishonest markup is worse than none — Google issues manual actions for schema that doesn't match visible content.

## Rules

- **JSON-LD only** (never microdata/RDFa), in a `<script type="application/ld+json">` in head or body.
- **Mirror visible content.** Never mark up ratings, FAQs, or facts that aren't on the page.
- Absolute URLs, ISO 8601 dates (`2026-07-20`), fully qualified `@id`s.
- Multiple entities on one page: bundle in a single `@graph` and cross-reference via `@id` (e.g. article `author` → the site-wide `Person`).
- Stable `@id` scheme (`https://site.com/#person`, `…/#website`, `…/article-slug#article`) so entities link across pages.

## Which schema for which page

| Page | Types |
| --- | --- |
| Home (personal site) | `WebSite` + `Person` (or `Organization` for a company) |
| Article / blog post | `Article` or `BlogPosting` + `BreadcrumbList`, `author` → `Person` |
| Product / app landing | `SoftwareApplication` (or `MobileApplication`) with `offers`, `aggregateRating` only if real |
| FAQ section | `FAQPage` (only if Q&A is visibly on the page) |
| How-to / tutorial | `HowTo` with steps |
| About / profile | `ProfilePage` wrapping the `Person` |
| Local business | `LocalBusiness` with NAP + `geo` + hours |

Property details per type: read [references/schema-types.md](references/schema-types.md).

## Validation (non-negotiable, in this order)

1. **Google Rich Results Test** (search.google.com/test/rich-results) — renders JavaScript, shows rich-result eligibility.
2. **validator.schema.org** — catches spec violations beyond Google's subset.
3. **Search Console → Enhancements** — ongoing monitoring after deploy.

## Detection caveat

Plain HTTP fetches strip `<script>` tags — **you cannot audit schema with curl/web-fetch**. Check the source that generates it (in Next.js: the JSON-LD builder + component), or a real browser, or the Rich Results Test. Reference implementation: `jterrazz-web` → `src/infrastructure/seo/json-ld.ts` (typed builders, unit-tested) rendered via a dedicated script component.

## Common mistakes

- Missing required properties (each rich result has its own required set — check the reference file).
- Relative URLs, non-ISO dates, `@id` collisions.
- `Person`/`Organization` inconsistencies across pages (different names/`sameAs` sets) — pick one canonical entity definition and reuse it.
- Marking up the whole site as `FAQPage` hoping for real estate — eligibility is per-page and policed.

## Related skills

- **jterrazz-seo** — schema is layer 4 of the audit; crawlability and indexation come first.
- **jterrazz-geo** — entity clarity feeds AI attribution; `sameAs` links (GitHub, X, LinkedIn) consolidate identity.
