# Schema types — properties per type

Companion to `jterrazz-structured-data`. "Required" = required for the rich result / core meaning; "recommended" = strengthens the entity. Only include properties you can populate truthfully from visible content.

## Person

```json
{
  "@type": "Person",
  "@id": "https://site.com/#person",
  "name": "…",
  "url": "https://site.com",
  "image": "https://site.com/portrait.jpg",
  "jobTitle": "…",
  "sameAs": ["https://github.com/…", "https://x.com/…", "https://linkedin.com/in/…"]
}
```

`sameAs` is the identity-consolidation lever — list every canonical profile, keep the set identical everywhere the Person appears.

## WebSite

Required: `name`, `url`. Recommended: `description`, `inLanguage`, `publisher`/`author` → `@id` of the Person/Organization. Add `potentialAction` (SearchAction) only if the site has real search.

## Organization

Required: `name`, `url`. Recommended: `logo` (ImageObject), `sameAs`, `contactPoint`. Use instead of Person when the brand isn't an individual.

## Article / BlogPosting

Required: `headline` (≤110 chars), `image`, `datePublished`.
Recommended: `dateModified` (only when genuinely updated), `author` → Person with `url`, `publisher`, `description`, `inLanguage`, `mainEntityOfPage`.
Use `BlogPosting` for blog content, `Article` for editorial/news-like pieces.

## BreadcrumbList

`itemListElement` of `ListItem`s with `position`, `name`, `item` (absolute URL). Last item may omit `item`.

## SoftwareApplication / MobileApplication

Required: `name`, `offers` (with `price`, `priceCurrency` — `"price": "0"` for free) **or** `aggregateRating` for the rich result.
Recommended: `applicationCategory`, `operatingSystem` (`"iOS"`, `"iOS, Android"`), `screenshot`, `downloadUrl`/`installUrl` (App Store / Play Store URLs), `author`.
`aggregateRating` only with real review data (store ratings qualify — keep them fresh).

## FAQPage

`mainEntity`: array of `Question` (`name`) each with `acceptedAnswer` → `Answer` (`text`). Every Q&A must be visible on the page. Google now shows FAQ rich results only for well-known authoritative sites — implement for the entity clarity, not the SERP real estate.

## HowTo

Required: `name`, `step` (array of `HowToStep` with `text`, optionally `name`, `image`, `url` anchors).
Recommended: `totalTime` (ISO 8601 duration, `PT30M`), `tool`, `supply`.

## ProfilePage

`mainEntity` → the `Person` `@id`. Signals "this page is about this person" — useful for about pages and author pages.

## LocalBusiness

Required: `name`, `address` (PostalAddress). Recommended: `geo` (GeoCoordinates), `openingHoursSpecification`, `telephone`, `priceRange`, `image`. NAP (name, address, phone) must match Google Business Profile exactly.

## ImageObject (photography pages)

`contentUrl`, `creator` → Person, `creditText`, `copyrightNotice`, `license`. Google surfaces `creator`/`credit`/`license` in image search ("Licensable" badge).

## @graph composition example

```json
{
  "@context": "https://schema.org",
  "@graph": [
    { "@type": "WebSite", "@id": "https://site.com/#website", "name": "…", "url": "https://site.com", "publisher": { "@id": "https://site.com/#person" } },
    { "@type": "Person", "@id": "https://site.com/#person", "name": "…" },
    { "@type": "BlogPosting", "headline": "…", "author": { "@id": "https://site.com/#person" }, "mainEntityOfPage": "https://site.com/articles/slug" }
  ]
}
```
