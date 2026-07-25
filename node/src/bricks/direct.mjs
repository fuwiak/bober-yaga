export default {
  id: "direct",
  title: "Yandex Direct",
  description: "Campaigns + OAuth (private — hidden on public profile)",
  visibility: "owner",
  aliases: ["dir"],
  async run(ctx, args) {
    const [sub, ...rest] = args;
    if (!sub || sub === "help" || sub === "--help") {
      console.log(`yaga direct

  yaga direct campaigns [status|create-unified|…]
  yaga direct oauth [refresh|ensure|…]

Delegates to scripts/yandex-direct-*.mjs
`);
      return;
    }
    if (sub === "campaigns" || sub === "campaign") {
      await ctx.runScript("yandex-direct-campaigns.mjs", rest.length ? rest : ["status"]);
      return;
    }
    if (sub === "oauth") {
      await ctx.runScript("yandex-direct-oauth.mjs", rest);
      return;
    }
    await ctx.runScript("yandex-direct-campaigns.mjs", args);
  },
};
