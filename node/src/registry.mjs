import { readdir } from "node:fs/promises";
import { join, basename } from "node:path";
import { pathToFileURL } from "node:url";
import { isBrickVisible, yagaRoot } from "./config.mjs";

/**
 * Brick contract:
 * {
 *   id: string,
 *   title: string,
 *   description?: string,
 *   visibility?: "public" | "owner",  // default owner
 *   async run(ctx, args): Promise<void>
 * }
 */
export async function loadBricks(cfg) {
  const dir = join(yagaRoot(), "src", "bricks");
  const files = (await readdir(dir))
    .filter((f) => f.endsWith(".mjs") && !f.startsWith("_"))
    .sort();

  const bricks = [];
  for (const file of files) {
    const mod = await import(pathToFileURL(join(dir, file)).href);
    const brick = mod.default;
    if (!brick?.id || typeof brick.run !== "function") {
      console.warn(`yaga: skip invalid brick ${file}`);
      continue;
    }
    brick.visibility = brick.visibility || "owner";
    brick.file = basename(file);
    bricks.push(brick);
  }

  return bricks.filter((b) => isBrickVisible(b, cfg));
}

export function listAllBrickMeta(cfg) {
  // used by `yaga bricks` — show visible + hidden with flags
  return loadAllBricksRaw().then((all) =>
    all.map((b) => ({
      ...b,
      visible: isBrickVisible(b, cfg),
    })),
  );
}

async function loadAllBricksRaw() {
  const dir = join(yagaRoot(), "src", "bricks");
  const files = (await readdir(dir))
    .filter((f) => f.endsWith(".mjs") && !f.startsWith("_"))
    .sort();
  const bricks = [];
  for (const file of files) {
    const mod = await import(pathToFileURL(join(dir, file)).href);
    const brick = mod.default;
    if (!brick?.id) continue;
    brick.visibility = brick.visibility || "owner";
    brick.file = basename(file);
    bricks.push(brick);
  }
  return bricks;
}
