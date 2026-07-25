export default {
  id: "metrika",
  title: "Yandex Metrika",
  description: "Counter status, content analytics, create counter",
  visibility: "public",
  aliases: ["metrica", "ym"],
  async run(ctx, args) {
    const [sub, ...rest] = args;
    if (!sub || sub === "help" || sub === "--help") {
      console.log(`yaga metrika

  yaga metrika status       counter + content analytics snapshot
  yaga metrika counter      create/find counter
  yaga metrika ytm          Yandex Tag Manager status
`);
      return;
    }
    if (sub === "status") {
      await ctx.runScript("yandex-metrika-status.mjs", rest);
      return;
    }
    if (sub === "counter") {
      await ctx.runScript("yandex-metrika-counter.mjs", rest);
      return;
    }
    if (sub === "ytm" || sub === "tag") {
      await ctx.runScript("yandex-ytm-status.mjs", rest);
      return;
    }
    await ctx.runScript("yandex-metrika-status.mjs", args);
  },
};
