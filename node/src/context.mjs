import { spawn } from "node:child_process";
import { join } from "node:path";
import { loadConfig, repoRoot, yagaRoot } from "./config.mjs";

export async function createContext(argv = []) {
  const cfg = await loadConfig();
  const root = repoRoot();
  return {
    argv,
    cfg,
    repoRoot: root,
    yagaRoot: yagaRoot(),
    scriptsDir: join(root, "scripts"),
    log: (...a) => console.log(...a),
    err: (...a) => console.error(...a),
    /** Run an existing repo script: scripts/<name> */
    runScript(name, args = [], { env = {}, inherit = true } = {}) {
      const script = join(root, "scripts", name);
      return runNode(script, args, { env: { ...process.env, ...env }, inherit });
    },
    /** Spawn any binary on PATH */
    runBin(bin, args = [], opts = {}) {
      return runCmd(bin, args, opts);
    },
  };
}

function runNode(script, args, { env, inherit }) {
  return runCmd(process.execPath, [script, ...args], { env, inherit });
}

function runCmd(cmd, args, { env = process.env, inherit = true } = {}) {
  return new Promise((resolve, reject) => {
    const child = spawn(cmd, args, {
      env,
      stdio: inherit ? "inherit" : ["ignore", "pipe", "pipe"],
      shell: false,
    });
    let stdout = "";
    let stderr = "";
    if (!inherit) {
      child.stdout?.on("data", (d) => {
        stdout += d;
      });
      child.stderr?.on("data", (d) => {
        stderr += d;
      });
    }
    child.on("error", reject);
    child.on("close", (code) => {
      if (code === 0) resolve({ code, stdout, stderr });
      else {
        const err = new Error(`${cmd} exited ${code}`);
        err.code = code;
        err.stdout = stdout;
        err.stderr = stderr;
        reject(err);
      }
    });
  });
}
