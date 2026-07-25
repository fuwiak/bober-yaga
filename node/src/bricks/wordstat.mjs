export default {
  id: "wordstat",
  title: "Wordstat",
  description: "Yandex Wordstat / Search API keyword research",
  visibility: "owner",
  aliases: ["ws"],
  async run(ctx, args) {
    const [sub, ...rest] = args;
    if (!sub || sub === "help" || sub === "--help") {
      console.log(`yaga wordstat — keyword research

  yaga wordstat run [args…]              full pipeline (scripts/yandex-wordstat.mjs)
  yaga wordstat search "<phrase>"        same pipeline; phrase via WORDSTAT_PHRASE

Seeds: data/wordstat-seeds-*.txt
`);
      return;
    }
    if (sub === "run" || sub === "pipeline") {
      await ctx.runScript("yandex-wordstat.mjs", rest);
      return;
    }
    if (sub === "search") {
      const phrase = rest.join(" ").trim();
      if (!phrase) {
        console.error('usage: yaga wordstat search "автоматизация бизнеса"');
        process.exitCode = 2;
        return;
      }
      await ctx.runScript("yandex-wordstat.mjs", [], {
        env: { WORDSTAT_PHRASE: phrase, YAGA_WORDSTAT_PHRASE: phrase },
      });
      return;
    }
    await ctx.runScript("yandex-wordstat.mjs", args);
  },
};
