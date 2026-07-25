/**
 * GPT / Foundation Models brick — chat via Yandex Cloud AI.
 * Not Claude Code; thin owner-facing klocek.
 */
import fetch from "../../../../scripts/lib/fetch.mjs";

export default {
  id: "gpt",
  title: "YandexGPT",
  description: "Chat completion via Foundation Models (owner)",
  visibility: "owner",
  aliases: ["yandexgpt", "llm"],
  async run(ctx, args) {
    const [sub, ...rest] = args;
    if (!sub || sub === "help" || sub === "--help") {
      console.log(`yaga gpt

  yaga gpt chat <prompt…>

Env (one of):
  YANDEX_GPT_API_KEY / YANDEX_IAM_TOKEN
  YANDEX_GPT_FOLDER_ID | YANDEX_FOLDER_ID
  YANDEX_GPT_MODEL   (default yandexgpt-lite)
`);
      return;
    }
    if (sub === "chat") {
      const prompt = rest.join(" ").trim();
      if (!prompt) {
        console.error("usage: yaga gpt chat <prompt>");
        process.exitCode = 2;
        return;
      }
      await chat(prompt);
      return;
    }
    console.error(`unknown: yaga gpt ${sub}`);
    process.exitCode = 2;
  },
};

async function chat(prompt) {
  const folderId =
    process.env.YANDEX_GPT_FOLDER_ID?.trim() ||
    process.env.YANDEX_FOLDER_ID?.trim() ||
    process.env.YANDEX_CLOUD_FOLDER_ID?.trim();
  const apiKey = process.env.YANDEX_GPT_API_KEY?.trim();
  const iam = process.env.YANDEX_IAM_TOKEN?.trim();
  const model = process.env.YANDEX_GPT_MODEL?.trim() || "yandexgpt-lite";

  if (!folderId || (!apiKey && !iam)) {
    console.error("Need YANDEX_GPT_FOLDER_ID and YANDEX_GPT_API_KEY (or YANDEX_IAM_TOKEN)");
    process.exitCode = 1;
    return;
  }

  const auth = apiKey ? `Api-Key ${apiKey}` : `Bearer ${iam}`;
  const modelUri = `gpt://${folderId}/${model}`;
  const res = await fetch("https://llm.api.cloud.yandex.net/foundationModels/v1/completion", {
    method: "POST",
    headers: {
      Authorization: auth,
      "Content-Type": "application/json",
      "x-folder-id": folderId,
    },
    body: JSON.stringify({
      modelUri,
      completionOptions: { stream: false, temperature: 0.3, maxTokens: 1000 },
      messages: [{ role: "user", text: prompt }],
    }),
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    console.error(JSON.stringify(data, null, 2) || `HTTP ${res.status}`);
    process.exitCode = 1;
    return;
  }
  const text = data?.result?.alternatives?.[0]?.message?.text;
  console.log(text || JSON.stringify(data, null, 2));
}
