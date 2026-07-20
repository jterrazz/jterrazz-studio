# SEO audit — full checklists

Companion to the `jterrazz-seo` skill. Work top-down; each section maps to the priority order in SKILL.md.

## 1. Crawlability

**robots.txt**

- No unintentional blocks on important paths; `Disallow` only for API routes, internals, redirect endpoints.
- References the sitemap (`Sitemap: https://…/sitemap.xml`).

**XML sitemap**

- Accessible, valid XML, submitted to Search Console.
- Contains canonical URLs only (no redirects, no noindexed pages, no 404s).
- `lastModified` reflects real content changes — engines learn to ignore fake freshness.
- Multi-locale: `xmlns:xhtml` namespace, `<xhtml:link>` alternates for all locales **including self** and `x-default`, absolute URLs everywhere.

**Architecture**

- Every important page ≤ 3 clicks from the homepage.
- No orphan pages (in sitemap but linked from nowhere).
- Large sites: parameterized URLs and faceted navigation under control (crawl budget).

## 2. Indexation

- Search Console coverage: indexed count ≈ expected count; investigate gaps.
- No `noindex` on important pages (check both meta tag and `X-Robots-Tag` header).
- Canonicals: unique self-reference per page; consistent HTTPS, host, and trailing-slash form across canonical + sitemap + hreflang + internal links.
- No redirect chains or loops; internal links point at final URLs.
- Soft 404s (thin "not found"-ish pages returning 200) cleaned up.
- Duplicate/near-duplicate content consolidated with canonicals or merged.

## 3. Technical foundations

**Core Web Vitals** — LCP < 2.5s, INP < 200ms, CLS < 0.1 (field data from the CWV report beats lab data).

Speed levers, in usual impact order: TTFB (caching/CDN/SSG), image sizing + modern formats, JavaScript payload, font loading (no invisible-text flash), CSS delivery, caching headers.

**Mobile** — responsive (one URL), correct viewport, adequate tap targets, no horizontal scroll, desktop feature parity (mobile-first indexing).

**Security** — HTTPS site-wide, valid cert, no mixed content, HTTP→HTTPS 301, HSTS optional.

## 4. On-page issues to hunt

- Duplicate or truncated titles/descriptions across pages (template bugs).
- Multiple H1s, skipped heading levels, headings used for styling.
- Thin pages: tag/category indexes with no unique value, doorway pages.
- Keyword cannibalization: two URLs ranking-flip for the same query.
- Broken internal links, generic anchors ("click here").

## 5. Content quality / E-E-A-T

- **Experience** — first-hand evidence: real examples, screenshots, numbers, lessons.
- **Expertise** — author identified with credentials; claims sourced.
- **Authoritativeness** — cited by others, consistent entity info (name, links, profiles).
- **Trust** — accurate, dated, contact/about pages, HTTPS, privacy policy.
- Decay: pages that lost traffic usually need a refresh (update facts, expand, re-date honestly) before any new content is written.

## Common issues by site type

| Site type | Usual suspects |
| --- | --- |
| Portfolio / personal | Missing Person/WebSite schema, weak internal linking from home to best work, no author entity consistency across profiles |
| SaaS / product | Shallow product pages, blog disconnected from product pages, missing comparison/alternative pages |
| Content / blog | Outdated posts never refreshed, cannibalization, no clustering, missing author pages |
| Multi-locale | Hreflang missing returns or self-refs, cross-locale canonicals, partially translated pages, no x-default |
| E-commerce | Thin categories, duplicated descriptions, facet-generated duplicates, missing Product schema |

## Tools

Free and sufficient for most audits: Google Search Console, PageSpeed Insights, Rich Results Test (renders JS — use it for schema), Bing Webmaster Tools. Paid, when scale demands: Screaming Frog, Ahrefs/Semrush.
