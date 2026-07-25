export default {
  id: "cloud",
  title: "Yandex Cloud (yc)",
  description: "Passthrough to official yc CLI",
  visibility: "public",
  aliases: ["yc"],
  async run(ctx, args) {
    if (!args.length || args[0] === "help" || args[0] === "--help") {
      console.log(`yaga cloud — thin wrapper around official yc

  yaga cloud compute instance list
  yaga cloud storage bucket list
  yaga cloud serverless function list
  yaga cloud -- --help

Requires: https://cloud.yandex.ru/docs/cli/
`);
      if (!args.length) return;
    }
    const passthrough = args[0] === "--" ? args.slice(1) : args;
    try {
      await ctx.runBin("yc", passthrough);
    } catch (err) {
      if (err?.code === "ENOENT" || String(err?.message || "").includes("ENOENT")) {
        console.error("yc not found on PATH. Install Yandex Cloud CLI first.");
        process.exitCode = 127;
        return;
      }
      process.exitCode = err.code || 1;
    }
  },
};
