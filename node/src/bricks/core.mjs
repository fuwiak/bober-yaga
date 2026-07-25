/** Meta brick — always visible; documents the klocek model. */
export default {
  id: "core",
  title: "Core",
  description: "About yaga brick architecture",
  visibility: "public",
  async run(ctx, args) {
    const sub = args[0] || "about";
    if (sub === "about" || sub === "info") {
      console.log(`yaga bricks = snap-in services

Each file in cli/yaga/src/bricks/*.mjs is a brick:
  export default { id, title, visibility, async run(ctx, args) }

visibility:
  public  → available when YAGA_PROFILE=public (safe to share)
  owner   → only on your machine (profile=owner, default)

Toggle:
  yaga bricks list
  yaga bricks hide direct
  yaga bricks enable metrika
  yaga profile public

Config: ~/.config/yaga/config.json
`);
      return;
    }
    console.error("usage: yaga core about");
    process.exitCode = 2;
  },
};
