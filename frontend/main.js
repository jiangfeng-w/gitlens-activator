let state = { presets: [], customs: [], busy: false, progress: null };

// 内置 IDE 固定顺序
const PRESET_ORDER = [
  "vscode",
  "vscode-insiders",
  "cursor",
  "windsurf",
  "kiro",
  "antigravity",
  "trae",
  "traecn",
  "codebuddy",
  "codebuddycn",
  "qoder",
  "qodercn",
  "catpaw",
];

function goApp() {
  return window.go.main.App;
}

// ---- 图标 ----
const I = {
  box: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="M3.3 7 12 12l8.7-5M12 22V12"/></svg>',
  trash: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M8 6V4a1 1 0 0 1 1-1h6a1 1 0 0 1 1 1v2"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/><path d="M10 11v6"/><path d="M14 11v6"/></svg>',
  spinner: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><path d="M12 3a9 9 0 1 0 9 9" /></svg>',
  info: '<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><circle cx="5" cy="12" r="2"/><circle cx="12" cy="12" r="2"/><circle cx="19" cy="12" r="2"/></svg>',
  bolt: '<svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><path d="M13 2 4.09 12.11a.6.6 0 0 0 .45 1h5.05l-1.6 8.53a.3.3 0 0 0 .53.26L19.91 11.9a.6.6 0 0 0-.45-1h-5.05l1.6-8.64a.3.3 0 0 0-.53-.26z"/></svg>',
  restore: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M3 12a9 9 0 1 0 3-6.7L3 8"/><path d="M3 3v5h5"/></svg>',
  check: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m4.5 12.5 5 5 10-11"/></svg>',
};

// 状态徽章文案（已激活带勾勾图标）
function badgeContent(status) {
  if (status === "activated") return I.check + "<span>已激活</span>";
  return "<span>" + { inactive: "未激活", missing: "未安装" }[status] + "</span>";
}

// 本地官方图标文件（frontend/icons/*.svg），随应用内嵌加载，无网络依赖。
// key 不在表中的（kiro/antigravity/qoder/catpaw 等）回落到下面的手绘 SVG。
const OFFICIAL_ICONS = {
  "vscode": "icons/vscode.svg",
  "vscode-insiders": "icons/vscode-insiders.svg",
  "cursor": "icons/cursor.svg",
  "windsurf": "icons/windsurf.svg",
  "kiro": "icons/kiro.svg",
  "antigravity": "icons/antigravity.svg",
  "trae": "icons/trae.svg",
  "traecn": "icons/trae.svg",
  "codebuddy": "icons/codebuddy.svg",
  "codebuddycn": "icons/codebuddy.svg",
  "qoder": "icons/qoder.svg",
  "qodercn": "icons/qoder.svg",
  "catpaw": "icons/catpaw.svg",
};

const tile = (txt, color) =>
  '<span class="tile" style="color:' + color + ';border:2px solid ' + color + '">' + txt + '</span>';

function productLogo(key) {
  const file = OFFICIAL_ICONS[key];
  if (file) {
    return '<img class="official-logo" src="' + file + '" alt="" aria-hidden="true">';
  }
  switch (key) {
    case "vscode":
      return '<svg viewBox="0 0 24 24" fill="#3aa1e0"><path d="M23.15 2.587 18.21.21a1.494 1.494 0 0 0-1.705.29l-9.46 8.63-4.12-3.128a.999.999 0 0 0-1.276.057L.327 7.261A1 1 0 0 0 .326 8.74L3.897 12 .326 15.26a1 1 0 0 0 .001 1.479L1.65 17.94a.999.999 0 0 0 1.276.057l4.12-3.128 9.46 8.63a1.492 1.492 0 0 0 1.704.29l4.942-2.377A1.5 1.5 0 0 0 24 20.06V3.939a1.5 1.5 0 0 0-.85-1.352zm-5.146 14.861L10.826 12l7.178-5.448v10.896z"/></svg>';
    case "vscode-insiders":
      return '<svg viewBox="0 0 24 24" fill="#29b3a4"><path d="M23.15 2.587 18.21.21a1.494 1.494 0 0 0-1.705.29l-9.46 8.63-4.12-3.128a.999.999 0 0 0-1.276.057L.327 7.261A1 1 0 0 0 .326 8.74L3.897 12 .326 15.26a1 1 0 0 0 .001 1.479L1.65 17.94a.999.999 0 0 0 1.276.057l4.12-3.128 9.46 8.63a1.492 1.492 0 0 0 1.704.29l4.942-2.377A1.5 1.5 0 0 0 24 20.06V3.939a1.5 1.5 0 0 0-.85-1.352zm-5.146 14.861L10.826 12l7.178-5.448v10.896z"/></svg>';
    case "cursor":
      return '<svg viewBox="0 0 24 24" fill="#dfe6ef"><path d="M24 11.52v.684c0 .98-.394 1.922-1.098 2.61l-8.08 7.903a4.09 4.09 0 0 1-5.744 0L2.024 15.75A3.665 3.665 0 0 1 .926 13.14v-.683c0-.98.394-1.921 1.098-2.61l8.08-7.902a4.09 4.09 0 0 1 5.743 0l7.054 6.967A3.665 3.665 0 0 1 24 11.52Zm-1.245 10.05-11.44 4.855a.315.315 0 0 1-.328-.055.324.324 0 0 1-.08-.363l4.85-11.565a.55.55 0 0 1 .895-.183l6.263 6.126a.733.733 0 0 1-.16 1.185ZM9.83 5.535l-4.85 11.56a.55.55 0 0 0 .12.607l2.174 2.13a.55.55 0 0 0 .98-.187l3.48-10.42a.55.55 0 0 1 .895-.183l7.84 7.67a.316.316 0 0 0 .45-.005.324.324 0 0 0-.005-.46L10.727 5.35a.55.55 0 0 0-.898.185Z"/></svg>';
    case "windsurf":
      return '<svg viewBox="0 0 24 24" fill="none"><path d="M3 5.5 14.5 3l3 4.2L6 9.7 3 5.5Z" fill="#2eb886"/><path d="M4.5 12 18 9.5l3 4.2-13.5 2.5-3-4.2Z" fill="#45d6a4"/><path d="M6 18.5 19.5 16l1.5 2.1L9.5 22 6 18.5Z" fill="#7ce7c3"/></svg>';
    case "kiro":
      return '<svg viewBox="0 0 24 24" fill="#a855f7"><path d="M12 2a7 7 0 0 0-7 7v11l2.5-2 2.25 2L12 18l2.25 2 2.25-2L19 20V9a7 7 0 0 0-7-7Z"/><circle cx="9.5" cy="10" r="1.4" fill="#1b1030"/><circle cx="14.5" cy="10" r="1.4" fill="#1b1030"/></svg>';
    case "antigravity":
      return '<svg viewBox="0 0 24 24" fill="#cfd8e3"><path d="M12 2 2.5 21l9.5-4 9.5 4L12 2Zm0 6.2L17 18l-5-2.1L7 18l5-9.8Z"/></svg>';
    case "trae":
    case "traecn":
      return '<svg viewBox="0 0 24 24"><rect width="24" height="24" rx="5" fill="#e6392b"/><path d="M7 7h10l-3.2 4H7V7Zm0 6h5.5L9.3 17H7v-4Z" fill="#fff"/></svg>';
    case "codebuddy":
    case "codebuddycn":
      return '<svg viewBox="0 0 24 24"><rect width="24" height="24" rx="5" fill="#3d7bff"/><path d="M8 9.5 12.5 8l.8 7-4.6-1.6-.2-1.8 2.6.9" fill="none" stroke="#fff" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/><circle cx="15.6" cy="15.6" r="1.4" fill="#7ee0c3"/></svg>';
    case "qoder":
    case "qodercn":
      return '<svg viewBox="0 0 24 24"><rect width="24" height="24" rx="5" fill="#6c5ce7"/><path d="M9 14.5a4.5 4.5 0 1 1 .4-8.7 5.2 5.2 0 0 1 10.2 1.7v7" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round"/><path d="m13.6 15.6 2 2 3.4-3.4" fill="none" stroke="#7ee0c3" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>';
    case "catpaw":
      return '<svg viewBox="0 0 24 24" fill="none"><ellipse cx="8" cy="9.4" rx="1.7" ry="2.3" fill="#f5a623"/><ellipse cx="12.6" cy="8.4" rx="1.7" ry="2.3" fill="#f5a623"/><ellipse cx="17.2" cy="10.2" rx="1.7" ry="2.3" fill="#f5a623"/><path d="M7.3 16.8c1.4 2.2 7.9 2.2 9.3 0 1.1-1.7-.6-3.4-2-3.2-1.3.2-1.5.6-2.6.6s-1.3-.4-2.6-.6c-1.4-.2-3.1 1.5-2.1 3.2z" fill="#f5a623"/></svg>';
    default:
      return tile("?", "#93a0b4");
  }
}

// ---- 数据获取 ----
async function refresh() {
  try {
    const data = await goApp().DetectAll();
    state.presets = data.presets || [];
    state.customs = data.customs || [];
    render();
  } catch (err) {
    showToast("加载列表失败: " + err.message, "error");
  }
}

function orderedCandidates() {
  const byKey = new Map(state.presets.map(p => [p.key, p]));
  const presets = PRESET_ORDER.map(k => byKey.get(k)).filter(Boolean);
  return [...presets, ...state.customs];
}

function rowKey(ed) {
  return ed.custom ? "custom|" + ed.extensionsDir : "preset|" + ed.key;
}

// 激活状态判定：后端检测所有扩展都已激活才算已激活
function isActivated(ed) {
  if (!ed.extensions || ed.extensions.length === 0) return false;
  return ed.extensions.every(ext => ext.activated);
}

// 卡片状态：activated / inactive / missing
function cardStatus(ed) {
  const hasGitLens = ed.installed && ed.extensions && ed.extensions.length > 0;
  if (!hasGitLens) return "missing";
  return isActivated(ed) ? "activated" : "inactive";
}

// ---- 渲染 ----
// 复用已有卡片节点，避免每次操作后所有卡片重播入场动画导致闪烁
function render() {
  const grid = document.getElementById("ideList");
  const all = orderedCandidates();

  if (all.length === 0) {
    grid.innerHTML = '<div class="empty-msg">无数据</div>';
    updateSummary();
    return;
  }

  const cards = Array.from(grid.querySelectorAll(".card"));
  const existing = new Map(cards.map(c => [c.dataset.key, c]));
  const keys = new Set(all.map(rowKey));

  // 1. 移除已消失的卡片（自定义目录被删除等）
  existing.forEach((card, key) => {
    if (!keys.has(key)) card.remove();
  });

  // 2. 原地更新已有卡片；新卡片按顺序插入。
  //    注意：绝不能把已有节点摘下再挂回（detach + reattach 会重置 CSS 动画，导致所有卡片重播入场动画闪烁）
  const emptyMsg = grid.querySelector(".empty-msg");
  if (emptyMsg) emptyMsg.remove();
  all.forEach(ed => {
    const key = rowKey(ed);
    const old = existing.get(key);
    if (old) {
      updateCard(old, ed);
    } else {
      grid.appendChild(createCard(ed));
    }
  });
  updateSummary();
}

function createCard(ed) {
  const key = rowKey(ed);
  const status = cardStatus(ed);
  const hasGitLens = status !== "missing";
  const activated = status === "activated";

  const card = document.createElement("div");
  card.className = "card" + (activated ? " selected" : "") + (!hasGitLens ? " locked" : "");
  card.dataset.key = key;
  card.setAttribute("role", "option");
  card.setAttribute("aria-selected", activated ? "true" : "false");
  card.tabIndex = hasGitLens ? 0 : -1;

  // 头部：logo + 名称 + 状态徽章 + 详情按钮
  const head = document.createElement("div");
  head.className = "card-head";

  const logo = document.createElement("div");
  logo.className = "logo";
  logo.innerHTML = productLogo(ed.key);
  head.appendChild(logo);

  const name = document.createElement("div");
  name.className = "name";
  name.textContent = ed.name;
  head.appendChild(name);

  const badge = document.createElement("span");
  badge.className = "badge status-" + status;
  badge.innerHTML = badgeContent(status);
  head.appendChild(badge);

  // 设置（详情）按钮
  const info = document.createElement("button");
  info.className = "card-info";
  info.title = "详情";
  info.setAttribute("aria-label", "详情");
  info.innerHTML = I.info;
  info.addEventListener("click", e => {
    e.stopPropagation();
    showDetail(ed);
  });
  head.appendChild(info);

  card.appendChild(head);

  // 底部提示
  const hint = document.createElement("div");
  hint.className = "card-hint" + (hasGitLens ? "" : " off");
  if (hasGitLens) {
    hint.innerHTML = activated
      ? I.restore + "<span>点击恢复</span>"
      : I.bolt + "<span>点击激活</span>";
  } else {
    hint.textContent = "无 GitLens";
  }
  card.appendChild(hint);

  // 点击卡片：已激活 -> 恢复；未激活 -> 激活
  // 注意：卡片节点会被 render() 复用（updateCard 只更新视觉），
  // 必须在点击时从 state 实时解析状态，不能用创建时闭包快照，
  // 否则激活/恢复动作永远停留在卡片首次创建时的状态。
  card.addEventListener("click", () => {
    if (state.busy) return;
    const cur = orderedCandidates().find(e => rowKey(e) === key);
    if (!cur || !cur.installed || !cur.extensions || cur.extensions.length === 0) return;
    if (isActivated(cur)) restoreOne(cur);
    else activateOne(cur);
  });

  return card;
}

// 复用已有卡片节点，仅更新动态内容（激活态、状态徽章、版本行、提示），不重建 DOM
function updateCard(card, ed) {
  const status = cardStatus(ed);
  const hasGitLens = status !== "missing";
  const activated = status === "activated";

  card.classList.toggle("selected", activated);
  card.classList.toggle("locked", !hasGitLens);
  card.setAttribute("aria-selected", activated ? "true" : "false");
  card.tabIndex = hasGitLens ? 0 : -1;

  // 状态徽章
  const badge = card.querySelector(".badge");
  if (badge) {
    badge.className = "badge status-" + status;
    badge.innerHTML = badgeContent(status);
  }

  // 底部提示
  const hint = card.querySelector(".card-hint");
  if (hint) {
    hint.classList.toggle("off", !hasGitLens);
    if (hasGitLens) {
      hint.innerHTML = activated
        ? I.restore + "<span>点击恢复</span>"
        : I.bolt + "<span>点击激活</span>";
    } else {
      hint.textContent = "无 GitLens";
    }
  }
}

function updateSummary() {
  const all = orderedCandidates();
  const found = all.filter(ed => ed.extensions && ed.extensions.length > 0).length;
  const act = all.filter(ed => isActivated(ed)).length;
  document.getElementById("summary").textContent =
    "已找到 " + all.length + " 个 IDE · 已激活 " + act + " 个 · 点击卡片切换激活 / 恢复";
}

// ---- 操作 ----
function extDirs(ed) {
  return (ed.extensions || []).map(ext => ext.dirPath);
}

async function activateOne(ed) {
  const dirs = extDirs(ed);
  const ok = await runAction(dirs, "激活");
  if (ok) showToast(ed.name + " 已激活", "success");
}

async function restoreOne(ed) {
  const dirs = extDirs(ed);
  const ok = await runAction(dirs, "恢复");
  if (ok) showToast(ed.name + " 已恢复", "success");
}

async function runAction(dirs, verb, confirmText) {
  if (state.busy) return false;
  if (confirmText) {
    const confirmed = await askConfirm(confirmText);
    if (!confirmed) return false;
  }
  setBusy(true);
  try {
    const results = await (verb === "激活" ? goApp().Activate(dirs) : goApp().Restore(dirs));
    const list = results || [];
    const fail = list.filter(r => !r.ok);
    if (fail.length > 0) {
      const msgs = fail.map(r => r.dirPath + ": " + r.message).join("\n");
      showToast(verb + "失败（" + fail.length + " 个）\n" + msgs, "error");
    }
    await refresh();
    return fail.length === 0;
  } catch (err) {
    showToast(verb + "失败: " + err.message, "error");
    await refresh();
    return false;
  } finally {
    setBusy(false);
  }
}

// 全部激活 / 全部恢复：对所有检测到 GitLens 的 IDE 执行
async function activateAll() {
  const targets = orderedCandidates().filter(ed => ed.extensions && ed.extensions.length > 0 && !isActivated(ed));
  if (targets.length === 0) {
    showToast("没有需要激活的 IDE", "error");
    return;
  }
  const dirs = targets.flatMap(extDirs);
  const ok = await runAction(dirs, "激活", "确定要激活以下 " + dirs.length + " 个 GitLens 吗？\n\n" + dirs.join("\n"));
  if (ok) showToast("已激活 " + targets.length + " 个 IDE", "success");
}

async function restoreAll() {
  const targets = orderedCandidates().filter(ed => ed.extensions && ed.extensions.length > 0 && isActivated(ed));
  if (targets.length === 0) {
    showToast("没有需要恢复的 IDE", "error");
    return;
  }
  const dirs = targets.flatMap(extDirs);
  const ok = await runAction(dirs, "恢复", "确定要恢复以下 " + dirs.length + " 个 GitLens 吗？\n\n" + dirs.join("\n"));
  if (ok) showToast("已恢复 " + targets.length + " 个 IDE", "success");
}

function setBusy(b) {
  state.busy = b;
  document.body.classList.toggle("is-busy", b);
  if (!b) {
    // 清除进度显示并还原 summary
    const summary = document.getElementById("summary");
    summary.classList.remove("has-progress");
    state.progress = null;
    updateSummary();
  }
}

// 批量操作进度：Go 端每完成一个目录 emit 一次 action:progress
function setupProgressListener() {
  if (!window.runtime || !window.runtime.EventsOn) return;
  window.runtime.EventsOn("action:progress", data => {
    if (!data || !state.busy) return;
    state.progress = data;
    const summary = document.getElementById("summary");
    summary.classList.add("has-progress");
    const failTxt = data.failures > 0 ? " · 失败 " + data.failures : "";
    summary.textContent = "处理中 " + data.done + " / " + data.total + failTxt;
  });
}

// ---- 自定义目录 ----
async function addCustom() {
  if (state.busy) return;
  const dir = document.getElementById("customDir").value.trim();
  if (!dir) {
    showToast("请输入目录路径", "error");
    return;
  }
  try {
    const data = await goApp().AddCustomDir(dir);
    document.getElementById("customDir").value = "";
    state.presets = data.presets || [];
    state.customs = data.customs || [];
    render();
    showToast("已添加自定义目录", "success");
  } catch (err) {
    showToast(err.message, "error");
  }
}

async function removeCustom(dir) {
  const confirmed = await askConfirm("确定要从自定义列表中删除该目录吗？\n\n" + dir);
  if (!confirmed) return;
  try {
    const data = await goApp().RemoveCustomDir(dir);
    state.presets = data.presets || [];
    state.customs = data.customs || [];
    render();
    showToast("已删除", "success");
  } catch (err) {
    showToast(err.message, "error");
  }
}

// ---- 确认弹窗（Promise 化） / Toast ----
function escapeHtml(s) {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;");
}

let confirmResolve = null;

function askConfirm(text) {
  return new Promise(resolve => {
    document.getElementById("confirmText").innerHTML = escapeHtml(text).replace(/\n/g, "<br>");
    document.getElementById("confirmOverlay").classList.remove("hidden");
    confirmResolve = resolve;
  });
}

function resolveConfirm(v) {
  document.getElementById("confirmOverlay").classList.add("hidden");
  if (confirmResolve) {
    confirmResolve(v);
    confirmResolve = null;
  }
}

// ---- 详情弹窗 ----
// 详情实时从 state 解析（与卡片点击同策略），展示路径 / 全部版本 / 备份状态，
// 自定义目录时提供删除入口。
function showDetail(ed) {
  const key = rowKey(ed);
  const cur = orderedCandidates().find(e => rowKey(e) === key) || ed;

  document.getElementById("detailTitle").textContent = cur.name + (cur.custom ? "（自定义）" : "");
  document.getElementById("detailPath").textContent = cur.extensionsDir;
  document.getElementById("detailPath").title = cur.extensionsDir;

  const table = document.getElementById("detailExtTable");
  table.innerHTML = "";
  const exts = cur.extensions || [];
  if (exts.length === 0) {
    const row = document.createElement("div");
    row.className = "detail-empty";
    row.textContent = "未检测到 GitLens 扩展";
    table.appendChild(row);
  } else {
    exts.forEach(ext => {
      const row = document.createElement("div");
      row.className = "detail-ext-row";

      const dot = document.createElement("span");
      dot.className = "dot" + (ext.activated ? "" : " off");
      row.appendChild(dot);

      const name = document.createElement("span");
      name.className = "detail-ext-name";
      name.textContent = "v" + ext.version + (ext.universal ? " universal" : "");
      name.title = ext.dirName;
      row.appendChild(name);

      const backup = document.createElement("span");
      backup.className = "backup-tag " + (ext.hasBackup ? "has" : "none");
      backup.textContent = ext.hasBackup ? "已备份" : "无备份";
      row.appendChild(backup);

      const stateTag = document.createElement("span");
      stateTag.className = "detail-state" + (ext.activated ? " on" : "");
      stateTag.textContent = ext.activated ? "已激活" : "未激活";
      row.appendChild(stateTag);

      table.appendChild(row);
    });
  }

  // 自定义目录：显示删除按钮
  const delBtn = document.getElementById("detailDelete");
  delBtn.classList.toggle("hidden", !cur.custom);
  delBtn.onclick = () => {
    hideDetail();
    removeCustom(cur.extensionsDir);
  };

  document.getElementById("detailOverlay").classList.remove("hidden");
}

function hideDetail() {
  document.getElementById("detailOverlay").classList.add("hidden");
}

document.getElementById("confirmBtn").addEventListener("click", () => resolveConfirm(true));
document.getElementById("confirmCancel").addEventListener("click", () => resolveConfirm(false));
document.getElementById("confirmOverlay").addEventListener("click", e => {
  if (e.target === e.currentTarget) resolveConfirm(false);
});

document.getElementById("detailClose").addEventListener("click", hideDetail);
document.getElementById("detailCloseBtn").addEventListener("click", hideDetail);
document.getElementById("detailOverlay").addEventListener("click", e => {
  if (e.target === e.currentTarget) hideDetail();
});
document.addEventListener("keydown", e => {
  if (e.key === "Escape" && !document.getElementById("detailOverlay").classList.contains("hidden")) {
    hideDetail();
  }
});

function showToast(msg, type) {
  const div = document.createElement("div");
  div.className = "toast toast-" + (type === "success" ? "success" : "error");
  div.textContent = msg;
  document.getElementById("toastContainer").appendChild(div);
  requestAnimationFrame(() => div.classList.add("show"));
  setTimeout(() => {
    div.classList.remove("show");
    setTimeout(() => div.remove(), 300);
  }, type === "success" ? 3000 : 6000);
}

// ---- 绑定 ----
document.getElementById("btnAddCustom").addEventListener("click", addCustom);
document.getElementById("btnRescan").addEventListener("click", () => {
  refresh();
  showToast("已重新扫描", "success");
});
document.getElementById("btnActivateAll").addEventListener("click", activateAll);
document.getElementById("btnRestoreAll").addEventListener("click", restoreAll);
document.getElementById("customDir").addEventListener("keydown", e => {
  if (e.key === "Enter") addCustom();
});

refresh();
setupProgressListener();
