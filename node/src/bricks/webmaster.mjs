export default {
  id: "webmaster",
  title: "Yandex Webmaster",
  description: "Hosts, SEO checklist, feed, mirrors, recrawl",
  visibility: "public",
  aliases: ["wm"],
  async run(ctx, args) {
    const [sub, ...rest] = args;
    if (!sub || sub === "help" || sub === "--help") {
      console.log(`yaga webmaster

  yaga webmaster status     overall webmaster + metrika snapshot
  yaga webmaster seo        ranking checklist (SQI, diagnostics, index, queries)
  yaga webmaster boost      recrawl important URLs + feed/region checklist
  yaga webmaster feed       upload performers feed (--all / --microsites)
  yaga webmaster mirrors    mirror settings
  yaga webmaster repair     feed repair helper
  yaga webmaster recrawl    submit URL for recrawl (or --quota)
  yaga webmaster oauth      ClientID/secret → access token (browser once)
`);
      return;
    }
    if (sub === "status") {
      await ctx.runScript("yandex-status.mjs", rest);
      return;
    }
    if (sub === "seo" || sub === "positions" || sub === "rank") {
      await ctx.runScript("yandex-webmaster-seo.mjs", rest);
      return;
    }
    if (sub === "boost" || sub === "index") {
      await ctx.runScript("yandex-webmaster-boost.mjs", rest);
      return;
    }
    if (sub === "feed") {
      await ctx.runScript("yandex-webmaster-feed.mjs", rest);
      return;
    }
    if (sub === "mirrors") {
      await ctx.runScript("yandex-webmaster-mirrors.mjs", rest);
      return;
    }
    if (sub === "repair") {
      await ctx.runScript("yandex-feed-repair.mjs", rest);
      return;
    }
    if (sub === "recrawl" || sub === "crawl") {
      await ctx.runScript("yandex-webmaster-recrawl.mjs", rest);
      return;
    }
    if (sub === "oauth" || sub === "login" || sub === "auth") {
      await ctx.runScript("yandex-webmaster-oauth.mjs", rest);
      return;
    }
    await ctx.runScript("yandex-status.mjs", args);
  },
};
