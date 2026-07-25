export default {
  id: "mail",
  title: "Yandex 360 Mail",
  description: "Corporate mailbox helpers (owner)",
  visibility: "owner",
  aliases: ["mail360", "y360"],
  async run(ctx, args) {
    const [sub, ...rest] = args;
    if (!sub || sub === "help" || sub === "--help") {
      console.log(`yaga mail

  yaga mail mailbox [args…]   → scripts/yandex360-mailbox.mjs
`);
      return;
    }
    if (sub === "mailbox" || sub === "box") {
      await ctx.runScript("yandex360-mailbox.mjs", rest);
      return;
    }
    await ctx.runScript("yandex360-mailbox.mjs", args);
  },
};
