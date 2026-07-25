import { homedir } from "node:os";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { access, mkdir, readFile, writeFile } from "node:fs/promises";
import { constants as fsConstants } from "node:fs";

const __dirname = dirname(fileURLToPath(import.meta.url));

/** Repo root (bober-ai): cli/yaga/node/src → ../../../.. */
export function repoRoot() {
  return resolve(__dirname, "../../../..");
}

export function yagaRoot() {
  return resolve(__dirname, "..");
}

export function defaultConfigPath() {
  return join(homedir(), ".config", "yaga", "config.json");
}

/**
 * Profiles control which bricks are visible.
 * - owner: everything (your machine)
 * - public: only bricks marked visibility: "public"
 * - custom: use enabled/disabled lists
 */
export const PROFILES = {
  owner: {
    description: "Full access — all bricks (default for you)",
    allow: "*",
  },
  public: {
    description: "Shareable subset — hide private marketing/billing bricks",
    allow: "public",
  },
  custom: {
    description: "Explicit enable/disable lists in config",
    allow: "custom",
  },
};

export function defaultConfig() {
  return {
    profile: "owner",
    enabled: [],
    disabled: [],
    hidden: [],
  };
}

export async function loadConfig() {
  const path = process.env.YAGA_CONFIG || defaultConfigPath();
  let fileCfg = {};
  try {
    const raw = await readFile(path, "utf8");
    fileCfg = JSON.parse(raw);
  } catch {
    /* missing is fine */
  }
  const cfg = { ...defaultConfig(), ...fileCfg };
  if (process.env.YAGA_PROFILE) {
    cfg.profile = process.env.YAGA_PROFILE.trim();
  }
  return { ...cfg, path };
}

export async function saveConfig(cfg) {
  const path = cfg.path || defaultConfigPath();
  await mkdir(dirname(path), { recursive: true });
  const { path: _p, ...rest } = cfg;
  await writeFile(path, `${JSON.stringify(rest, null, 2)}\n`, "utf8");
  return path;
}

export async function pathExists(p) {
  try {
    await access(p, fsConstants.F_OK);
    return true;
  } catch {
    return false;
  }
}

/**
 * Decide if a brick is visible for the active profile/config.
 * Brick fields:
 *   visibility: "public" | "owner" (default owner)
 */
export function isBrickVisible(brick, cfg) {
  const id = brick.id;
  if (cfg.hidden?.includes(id)) return false;
  if (cfg.disabled?.includes(id)) return false;

  const profile = cfg.profile || "owner";
  if (profile === "owner") {
    if (cfg.enabled?.length && !cfg.enabled.includes(id) && brick.id !== "core") {
      // enabled list is optional allow-list when set under owner+enabled
      // Only apply if user explicitly set enabled non-empty under custom; under owner ignore empty
    }
    return true;
  }

  if (profile === "public") {
    return brick.visibility === "public";
  }

  // custom
  if (cfg.enabled?.length) {
    return cfg.enabled.includes(id) || id === "core";
  }
  return !cfg.disabled?.includes(id);
}
