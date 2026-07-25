export default {
  id: "disk",
  title: "Yandex Disk",
  description: "Passthrough to yandex-disk sync client (if installed)",
  visibility: "public",
  aliases: ["ydisk"],
  async run(ctx, args) {
    if (!args.length || args[0] === "help" || args[0] === "--help") {
      console.log(`yaga disk — wrapper for official yandex-disk

  yaga disk status
  yaga disk start | stop | sync
  yaga disk -- setup

Install: https://repo.yandex.ru/yandex-disk/
`);
      if (!args.length) return;
    }
    const passthrough = args[0] === "--" ? args.slice(1) : args;
    try {
      await ctx.runBin("yandex-disk", passthrough);
    } catch (err) {
      if (err?.code === "ENOENT" || String(err?.message || err).includes("ENOENT")) {
        console.error("yandex-disk not found. Install the official Linux client, or use rclone.");
        process.exitCode = 127;
        return;
      }
      process.exitCode = typeof err.code === "number" ? err.code : 1;
    }
  },
};
