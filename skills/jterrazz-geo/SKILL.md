---
name: jterrazz-geo
description: AI search & LLM visibility (GEO/AEO) — getting cited by ChatGPT, Claude, Perplexity, and Google AI Overviews. Use when the user mentions "GEO", "AEO", "AI SEO", "LLM visibility", "AI citations", "llms.txt", "AI crawlers", "ChatGPT doesn't mention my site", "AI Overviews", or wants their content/product to show up in AI answers. For classic search rankings see jterrazz-seo; for JSON-LD see jterrazz-structured-data.
---

# GEO — visibility in AI answers

AI engines don't rank pages, they **cite passages**. A page ranking #15 with a clean, self-contained, well-sourced answer gets cited over a #1 page that rambles. Optimize for extraction and evidence, not position.

## Three pillars

1. **Structure** — format content so a passage can be lifted out and quoted alone.
2. **Authority** — evidence density. The Princeton GEO study measured visibility lifts of **+40% from adding citations, +37% from statistics, +25% from quotes/authoritative tone**. Keyword stuffing measurably *hurts*.
3. **Presence** — LLMs weigh third-party sources heavily (Wikipedia, Reddit, reviews, comparison posts). Owned content alone caps your ceiling. Comparison content is the most-cited type (~1/3 of citations), then definitive guides and original research.

## Crawler access (prerequisite — blocked = never cited)

robots.txt must allow the AI crawlers you want citations from:

| Bot | Engine |
| --- | --- |
| `GPTBot`, `ChatGPT-User`, `OAI-SearchBot` | OpenAI / ChatGPT search |
| `ClaudeBot`, `anthropic-ai` | Anthropic / Claude |
| `PerplexityBot`, `Perplexity-User` | Perplexity |
| `Google-Extended` | Gemini grounding (AI Overviews use normal Googlebot) |
| `Bingbot` | Microsoft Copilot |

Also verify in server/CDN logs that they actually hit the site — a WAF or bot-protection default can silently block them.

## Extractable formatting

- **Question-shaped H2/H3s**, each followed by a **direct 40–60 word answer** before any elaboration. That block is what gets quoted.
- Keep passages **self-contained** (~100–170 words): a reader landing mid-page understands it without the surrounding context. Name the subject; avoid "it/this" referring outside the block.
- Definitions first ("X is …"), then nuance. Comparison content in real tables. FAQs for long-tail questions.
- Give every claim a number or a source where possible — that's the +40/+37% lever.
- Content must be in the **server-rendered HTML**. Client-only rendering is invisible to most AI crawlers.

## llms.txt — honest take

`/llms.txt` is a markdown index of your site (title, summary, links to key pages) for LLM consumption. It costs minutes to serve and signals AI-friendliness, but **no major engine has confirmed using it as a citation lever** — add it, don't expect miracles from it. In Next.js: a simple route handler at `app/llms.txt/route.ts` listing the site's pages/articles with one-line summaries.

## Don'ts (these backfire)

- No separate "AI version" of content, no page chunking into fragments, no AI-keyword rewriting — Google treats GEO/AEO tricks as the same spam it already polices.
- No scaled low-value generation, no inauthentic third-party mentions (fake Reddit posts, planted reviews).

## Monitoring

- Monthly: run 5–10 priority queries through ChatGPT, Perplexity, and Google (AI Overviews); log whether you're cited and who is instead.
- Analytics: segment referrals from `chatgpt.com`, `perplexity.ai`, `claude.ai`, `gemini.google.com`.
- Logs: track AI-bot crawl volume as a leading indicator.

## Related skills

- **jterrazz-seo** — the technical base (crawlable, indexed, fast) is a prerequisite; GEO doesn't rescue an uncrawlable site.
- **jterrazz-structured-data** — entity clarity (Person/Organization/Article) helps engines attribute content.
- **jterrazz-content-reach** — the presence pillar: earning third-party mentions.
