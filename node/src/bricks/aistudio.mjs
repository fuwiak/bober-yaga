/**
 * AI Studio brick — optional; needs pip package.
 */
export default {
  id: "aistudio",
  title: "Yandex AI Studio",
  description: "Vector stores / RAG indexing via yandex-ai-studio (owner)",
  visibility: "owner",
  aliases: ["studio", "rag"],
  async run(ctx, args) {
    if (!args.length || args[0] === "help" || args[0] === "--help") {
      console.log(`yaga aistudio — delegates to yandex-ai-studio if installed

  pip install yandex-ai-studio-sdk
  yaga aistudio -- vector-stores local docs/*.pdf --name "Docs"

This brick is a klocek face; heavy lifting stays in the official SDK CLI.
`);
      if (!args.length) return;
    }
    const passthrough = args[0] === "--" ? args.slice(1) : args;
    try {
      await ctx.runBin("yandex-ai-studio", passthrough);
    } catch (err) {
      if (err?.code === "ENOENT" || String(err?.message || err).includes("ENOENT")) {
        console.error("yandex-ai-studio not found. pip install yandex-ai-studio-sdk");
        process.exitCode = 127;
        return;
      }
      process.exitCode = typeof err.code === "number" ? err.code : 1;
    }
  },
};
