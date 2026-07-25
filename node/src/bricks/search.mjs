/**
 * Search brick — Yandex Cloud Search / web search API.
 * Official dedicated CLI does not exist; this is the klocek face.
 */
import fetch from "../../../../scripts/lib/fetch.mjs";

export default {
  id: "search",
  title: "Yandex Search API",
  description: "Cloud Search / web query (owner)",
  visibility: "owner",
  aliases: ["websearch"],
  async run(ctx, args) {
    const [sub, ...rest] = args;
    if (!sub || sub === "help" || sub === "--help") {
      console.log(`yaga search

  yaga search web "<query>"     Cloud Search API (needs YANDEX_SEARCH_API_KEY + folder)

Env:
  YANDEX_SEARCH_API_KEY
  YANDEX_SEARCH_FOLDER_ID | YANDEX_FOLDER_ID
`);
      return;
    }
    if (sub === "web" || sub === "query") {
      const q = rest.join(" ").trim();
      if (!q) {
        console.error('usage: yaga search web "AI автоматизация"');
        process.exitCode = 2;
        return;
      }
      await searchWeb(q);
      return;
    }
    console.error(`unknown: yaga search ${sub}`);
    process.exitCode = 2;
  },
};

async function searchWeb(query) {
  const apiKey = process.env.YANDEX_SEARCH_API_KEY?.trim();
  const folderId =
    process.env.YANDEX_SEARCH_FOLDER_ID?.trim() ||
    process.env.YANDEX_FOLDER_ID?.trim() ||
    process.env.YANDEX_CLOUD_FOLDER_ID?.trim();
  if (!apiKey || !folderId) {
    console.error("Need YANDEX_SEARCH_API_KEY and YANDEX_SEARCH_FOLDER_ID");
    process.exitCode = 1;
    return;
  }

  // Yandex Cloud Search API (Web Search) — async operation style varies by product.
  // Use the documented sync-ish WebSearch endpoint when available.
  const url = "https://searchapi.api.cloud.yandex.net/v2/web/search";
  const body = {
    query: { searchType: "SEARCH_TYPE_RU", queryText: query },
    folderId,
    responseFormat: "FORMAT_XML",
  };

  const res = await fetch(url, {
    method: "POST",
    headers: {
      Authorization: `Api-Key ${apiKey}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });
  const text = await res.text();
  if (!res.ok) {
    console.error(`Search API HTTP ${res.status}: ${text.slice(0, 800)}`);
    process.exitCode = 1;
    return;
  }
  console.log(text.slice(0, 4000));
  if (text.length > 4000) console.log("…(truncated)");
}
