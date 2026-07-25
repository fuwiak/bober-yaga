import { createContext } from "./context.mjs";
import { loadBricks } from "./registry.mjs";
import { PROFILES } from "./config.mjs";

export async function run(argv) {
  const ctx = await createContext(argv);
  const bricks = await loadBricks(ctx.cfg);

  if (argv.includes("--version") || argv.includes("-V")) {
    console.log("yaga 0.1.0");
    return;
  }

  if (argv.length === 0 || argv[0] === "--help" || argv[0] === "-h" || argv[0] === "help") {
    printHelp(ctx, bricks);
    return;
  }

  const [cmd, ...rest] = argv;

  // Built-in meta (always available)
  if (cmd === "bricks" || cmd === "plugins" || cmd === "blocks") {
    await cmdBricks(ctx, rest);
    return;
  }
  if (cmd === "profile") {
    await cmdProfile(ctx, rest);
    return;
  }
  if (cmd === "doctor") {
    await cmdDoctor(ctx, bricks);
    return;
  }
  if (cmd === "help") {
    printHelp(ctx, bricks);
    return;
  }

  const brick = bricks.find((b) => b.id === cmd || b.aliases?.includes(cmd));
  if (!brick) {
    console.error(`yaga: unknown command '${cmd}' (profile=${ctx.cfg.profile})`);
    console.error(`Try: yaga bricks   or   yaga help`);
    process.exitCode = 2;
    return;
  }

  await brick.run(ctx, rest);
}

function printHelp(ctx, bricks) {
  const { cfg } = ctx;
  console.log(`yaga — modular Yandex CLI (bricks you can snap on/off)

Profile: ${cfg.profile}  ·  config: ${cfg.path}
Profiles: ${Object.keys(PROFILES).join(", ")}  (YAGA_PROFILE=public to hide owner bricks)

Usage:
  yaga <brick> [args…]
  yaga bricks                 list bricks (on/off for this profile)
  yaga profile [name]         show / set profile (owner|public|custom)
  yaga doctor                 check tokens / binaries
  yaga help

Bricks:`);
  for (const b of bricks) {
    const vis = b.visibility === "public" ? "public" : "owner ";
    console.log(`  ${b.id.padEnd(12)} [${vis}]  ${b.title}`);
    if (b.description) console.log(`               ${b.description}`);
  }
  console.log(`
Examples:
  yaga wordstat search "автоматизация бизнеса"
  yaga direct campaigns status
  yaga metrika status
  yaga webmaster status
  yaga cloud -- help
  YAGA_PROFILE=public yaga help`);
}

async function cmdBricks(ctx, args) {
  const { listAllBrickMeta } = await import("./registry.mjs");
  const { saveConfig } = await import("./config.mjs");
  const sub = args[0] || "list";

  if (sub === "list" || sub === "ls") {
    const all = await listAllBrickMeta(ctx.cfg);
    console.log(`profile=${ctx.cfg.profile}\n`);
    for (const b of all) {
      const mark = b.visible ? "●" : "○";
      console.log(
        `  ${mark} ${b.id.padEnd(12)} vis=${(b.visibility || "owner").padEnd(6)}  ${b.title}`,
      );
    }
    console.log(`\n● visible  ○ hidden for this profile`);
    console.log(`Enable/disable: yaga bricks enable <id> | yaga bricks disable <id>`);
    console.log(`Or set profile: yaga profile public`);
    return;
  }

  if (sub === "enable" || sub === "disable") {
    const id = args[1];
    if (!id) {
      console.error(`usage: yaga bricks ${sub} <id>`);
      process.exitCode = 2;
      return;
    }
    const cfg = { ...ctx.cfg };
    cfg.profile = "custom";
    cfg.enabled = [...(cfg.enabled || [])];
    cfg.disabled = [...(cfg.disabled || [])];
    if (sub === "enable") {
      cfg.disabled = cfg.disabled.filter((x) => x !== id);
      if (!cfg.enabled.includes(id)) cfg.enabled.push(id);
      // if enabled list was empty before, seed with currently visible + this one
    } else {
      cfg.enabled = cfg.enabled.filter((x) => x !== id);
      if (!cfg.disabled.includes(id)) cfg.disabled.push(id);
    }
    const path = await saveConfig(cfg);
    console.log(`saved ${sub} ${id} → ${path} (profile=custom)`);
    return;
  }

  if (sub === "hide" || sub === "unhide") {
    const id = args[1];
    if (!id) {
      console.error(`usage: yaga bricks ${sub} <id>`);
      process.exitCode = 2;
      return;
    }
    const cfg = { ...ctx.cfg };
    cfg.hidden = [...(cfg.hidden || [])];
    if (sub === "hide") {
      if (!cfg.hidden.includes(id)) cfg.hidden.push(id);
    } else {
      cfg.hidden = cfg.hidden.filter((x) => x !== id);
    }
    const path = await saveConfig(cfg);
    console.log(`saved ${sub} ${id} → ${path}`);
    return;
  }

  console.error(`usage: yaga bricks [list|enable|disable|hide|unhide]`);
  process.exitCode = 2;
}

async function cmdProfile(ctx, args) {
  const { saveConfig, PROFILES: P } = await import("./config.mjs");
  const name = args[0];
  if (!name) {
    console.log(`current: ${ctx.cfg.profile}`);
    for (const [k, v] of Object.entries(P)) {
      const mark = k === ctx.cfg.profile ? "→" : " ";
      console.log(`  ${mark} ${k.padEnd(8)} ${v.description}`);
    }
    return;
  }
  if (!P[name]) {
    console.error(`unknown profile '${name}'. Use: ${Object.keys(P).join("|")}`);
    process.exitCode = 2;
    return;
  }
  const cfg = { ...ctx.cfg, profile: name };
  const path = await saveConfig(cfg);
  console.log(`profile=${name} → ${path}`);
}

async function cmdDoctor(ctx, bricks) {
  const checks = [
    ["YANDEX_OAUTH_TOKEN / METRIKA", !!(process.env.YANDEX_OAUTH_TOKEN || process.env.YANDEX_METRIKA_OAUTH_TOKEN)],
    ["YANDEX_WEBMASTER_OAUTH_TOKEN", !!process.env.YANDEX_WEBMASTER_OAUTH_TOKEN],
    ["YANDEX_DIRECT_OAUTH_TOKEN / REFRESH", !!(process.env.YANDEX_DIRECT_OAUTH_TOKEN || process.env.YANDEX_DIRECT_REFRESH_TOKEN)],
    ["YANDEX_DIRECT_CLIENT_ID", !!process.env.YANDEX_DIRECT_CLIENT_ID],
    ["YANDEX_SEARCH_API_KEY", !!process.env.YANDEX_SEARCH_API_KEY],
    ["yc on PATH", await binOk("yc")],
    ["yandex-disk on PATH", await binOk("yandex-disk")],
  ];
  console.log(`yaga doctor · profile=${ctx.cfg.profile} · bricks=${bricks.length}\n`);
  for (const [label, ok] of checks) {
    console.log(`  ${ok ? "✓" : "○"} ${label}`);
  }
  console.log(`\nVisible bricks: ${bricks.map((b) => b.id).join(", ")}`);
}

function binOk(name) {
  return new Promise((resolve) => {
    import("node:child_process").then(({ spawn }) => {
      const c = spawn("which", [name], { stdio: "ignore" });
      c.on("close", (code) => resolve(code === 0));
      c.on("error", () => resolve(false));
    });
  });
}
