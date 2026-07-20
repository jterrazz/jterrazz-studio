---
name: jterrazz-seo
description: Technical and on-page SEO for web projects — audit workflow, crawlability, indexation, Core Web Vitals, metadata, canonicals, hreflang, internal linking. Use when the user mentions "SEO", "SEO audit", "not ranking", "traffic dropped", "meta tags", "sitemap", "robots.txt", "canonical", "hreflang", "Core Web Vitals", "indexing", or wants more organic traffic — even a vague "my SEO is bad" starts here with an audit. For AI/LLM search visibility see jterrazz-geo; for JSON-LD see jterrazz-structured-data; for content strategy and distribution see jterrazz-content-reach.
---

# SEO — technical & on-page

Diagnose before prescribing. An SEO problem is almost never one thing; audit in priority order and report evidence, not vibes.

## Priority order

Fix in this order — each layer is worthless without the one above it:

1. **Crawlability** — can engines reach the pages? (robots.txt, sitemap, internal links, orphan pages)
2. **Indexation** — are the right pages indexed? (noindex, canonicals, redirect chains, soft 404s, duplicates)
3. **Technical foundations** — speed, mobile, HTTPS. Core Web Vitals: LCP < 2.5s, INP < 200ms, CLS < 0.1.
4. **On-page** — titles, descriptions, headings, content-keyword alignment.
5. **Content quality & authority** — depth, E-E-A-T, links. Last, not first.

## Audit workflow

1. **Context first.** Which pages/keywords matter most? Search Console access? Recent changes or migrations? Traffic baseline? Don't audit blind.
2. **Crawl the real site**, not the repo: check `robots.txt`, `/sitemap.xml`, a few key pages' rendered HTML. In the repo, verify what generates them.
3. **Report every finding as:** Issue → Impact (High/Med/Low) → Evidence (how you found it) → Fix (specific action) → Priority.
4. **Close with:** overall health, top 3–5 priorities, quick wins (easy + immediate), long-term items.

## Measurement traps (report false findings and you lose trust)

- **Plain HTTP fetches strip `<script>` tags** — JSON-LD and JS-injected markup are invisible to `curl`/web-fetch. Never report "missing schema" from a fetch; verify with a browser (`document.querySelectorAll('script[type="application/ld+json"]')`) or Google's Rich Results Test.
- **Client-side-rendered content** may be absent from initial HTML. Check the SSR/SSG output before declaring content thin or missing.
- `site:domain.com` counts are approximate — use Search Console coverage for real indexation data.

## On-page essentials

| Element | Rule |
| --- | --- |
| Title | Unique, primary keyword near the start, 50–60 chars, compelling |
| Meta description | Unique, 150–160 chars, value proposition + reason to click |
| H1 | Exactly one, contains the keyword, logical H2/H3 hierarchy below |
| URL | Short, lowercase, hyphenated, descriptive, stable |
| Content | Keyword in first 100 words, satisfies the search intent, deeper than competitors |
| Images | Descriptive filenames, alt text, modern formats, lazy-load below the fold |
| Internal links | Important pages get the most links; descriptive anchors; no broken links |

Per-page: one clear primary keyword, no two pages targeting the same one (cannibalization). Site-wide: map keywords → pages, cluster topics.

## i18n (multi-locale sites)

- Every page: **self-referencing hreflang + reciprocal pairs + `x-default`**. One missing return link drops the whole pair.
- Codes are ISO 639-1 (+ optional region): `en`, `en-GB` — never `en-UK`.
- Each locale **self-canonicals**. Cross-locale canonicals kill indexing of the other locales.
- Only ship locales you fully translate (title + meta + body, not just chrome). Thin locales drag the whole site down — better absent than machine-thin.
- Locale in the URL path (`/fr/`), never via IP/Accept-Language negotiation.

## Next.js App Router notes

- `app/sitemap.ts`, `app/robots.ts`, and the Metadata API are the canonical mechanisms — no third-party SEO packages.
- **Caveat:** `alternates.languages` does not auto-add the self-referencing hreflang for the current locale — add it explicitly.
- Reference implementation: `jterrazz-web` → `reach.config.ts` (the single source of truth) projected by `@jterrazz/reach` (`packages/reach`): sitemap/robots/llms.txt are one-liner routes, the identity JSON-LD is a projected `@graph`, and `specs/website/` runs the `conformance(website, site)` rule pack against the built site via `specification.website()` (`@jterrazz/test`). Prefer wiring a site into reach over hand-rolling these surfaces.

## Deep checklists

For the full technical audit checklists (crawlability, indexation, speed, security, common issues by site type), read [references/audit-checklists.md](references/audit-checklists.md).

## Related skills

- **jterrazz-geo** — citations in ChatGPT/Perplexity/AI Overviews, llms.txt, AI crawlers.
- **jterrazz-structured-data** — JSON-LD implementation and validation.
- **jterrazz-content-reach** — what to write and how to distribute it.
