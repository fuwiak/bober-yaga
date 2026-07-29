export default {
  id: "metrika",
  title: "Yandex Metrika",
  description: "Counter status, filters, organic traffic",
  visibility: "public",
  aliases: ["metrica", "ym"],
  async run(ctx, args) {
    const [sub, ...rest] = args;
    if (!sub || sub === "help" || sub === "--help") {
      console.log(`yaga metrika

  yaga metrika status       counter + goals snapshot
  yaga metrika counter      create/find counter
  yaga metrika filters      ensure "don't count my visits" + list filters
  yaga metrika organic      organic traffic for period
  yaga metrika ecommerce    enable dataLayer ecommerce + RUB
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
    if (sub === "filters" || sub === "filter") {
      await ctx.runScript("yandex-metrika-filters.mjs", rest);
      return;
    }
    if (sub === "organic") {
      await ctx.runScript("yandex-metrika-filters.mjs", ["--organic", ...rest]);
      return;
    }
    if (sub === "ecommerce" || sub === "ecom") {
      await ctx.runScript("yandex-metrika-ecommerce.mjs", rest);
      return;
    }
    if (sub === "ytm" || sub === "tag") {
      await ctx.runScript("yandex-ytm-status.mjs", rest);
      return;
    }
    await ctx.runScript("yandex-metrika-status.mjs", args);
  },
};
