/*
  Free-form canvas editor.

  Design note on where rendering happens: dragging and resizing patch the
  element's inline geometry directly, so pointer movement stays at 60fps with
  no round-trip. Everything else (colour, font, shape, effects) is re-rendered
  by the server from the same Go code the public page uses, which keeps one
  source of truth for appearance instead of a JS copy that can drift.
*/
(() => {
  'use strict';

  const root = document.querySelector('[data-canvas-editor]');
  if (!root) return;

  const stageWrap = root.querySelector('[data-cv-stage-wrap]');
  const propsPanel = root.querySelector('[data-cv-props]');
  const layerList = root.querySelector('[data-cv-layers]');
  const jsonField = root.querySelector('[data-cv-json]');
  const emptyHint = root.querySelector('[data-cv-empty]');

  /** @type {{ratio:number,bg:string,bgImage:string,bgFit:string,overlay:number,grid:number,elements:Array}} */
  let canvas;
  try {
    canvas = JSON.parse(jsonField.value || '{}');
  } catch {
    canvas = {};
  }
  if (!Array.isArray(canvas.elements)) canvas.elements = [];
  if (!canvas.ratio) canvas.ratio = 1.9;
  if (canvas.grid === undefined) canvas.grid = 1;

  let selectedId = null;
  let dirty = false;

  const el = (id) => canvas.elements.find((e) => e.id === id);
  const stage = () => stageWrap.querySelector('.cv-stage');

  /* ================= server-side render ================= */

  let renderTimer = null;
  let rendering = false;
  let renderQueued = false;

  const pushState = () => {
    jsonField.value = JSON.stringify(canvas);
    dirty = true;
  };

  const renderStage = async () => {
    if (rendering) { renderQueued = true; return; }
    rendering = true;

    try {
      const body = new FormData();
      body.set('canvas', JSON.stringify(canvas));
      body.set('csrf_token', window.csrfToken());

      const res = await fetch('/panel/design/canvas/render', {
        method: 'POST',
        body,
        headers: { 'X-Requested-With': 'fetch', 'X-CSRF-Token': window.csrfToken() },
      });
      if (res.ok) {
        stageWrap.innerHTML = await res.text();
        decorateStage();
      }
    } catch {
      /* keep the current stage on a failed refresh */
    } finally {
      rendering = false;
      if (renderQueued) { renderQueued = false; renderStage(); }
    }
  };

  const scheduleRender = () => {
    clearTimeout(renderTimer);
    renderTimer = setTimeout(renderStage, 90);
  };

  /* ================= selection & handles ================= */

  const HANDLES = ['nw', 'n', 'ne', 'e', 'se', 's', 'sw', 'w'];

  const clearDecorations = () => {
    stageWrap.querySelectorAll('.cv-handle, .cv-guide, .cv-readout').forEach((n) => n.remove());
    stageWrap.querySelectorAll('.cv-el.is-selected').forEach((n) => n.classList.remove('is-selected'));
  };

  const placeHandles = (node) => {
    const s = stage();
    if (!s) return;

    const pos = {
      nw: [0, 0], n: [50, 0], ne: [100, 0], e: [100, 50],
      se: [100, 100], s: [50, 100], sw: [0, 100], w: [0, 50],
    };
    const e = el(node.dataset.id);
    if (!e) return;

    HANDLES.forEach((h) => {
      const dot = document.createElement('span');
      dot.className = 'cv-handle';
      dot.dataset.h = h;
      dot.dataset.for = e.id;
      const [px, py] = pos[h];
      dot.style.left = `calc(${e.x + (e.w * px) / 100}% - 5.5px)`;
      dot.style.top = `calc(${e.y + (e.h * py) / 100}% - 5.5px)`;
      s.appendChild(dot);
    });
  };

  const select = (id) => {
    selectedId = id;
    clearDecorations();

    const node = id && stageWrap.querySelector(`.cv-el[data-id="${id}"]`);
    if (node) {
      node.classList.add('is-selected');
      placeHandles(node);
    }
    renderProps();
    renderLayers();
  };

  // Re-apply selection chrome after a server re-render.
  const decorateStage = () => {
    const s = stage();
    if (!s) return;
    s.style.setProperty('--cv-grid-op', canvas.grid > 0 ? '.18' : '0');

    if (emptyHint) emptyHint.hidden = canvas.elements.length > 0;

    if (selectedId) {
      const node = stageWrap.querySelector(`.cv-el[data-id="${selectedId}"]`);
      if (node) {
        node.classList.add('is-selected');
        placeHandles(node);
      } else {
        selectedId = null;
      }
    }
  };

  /* ================= drag & resize ================= */

  const snap = (v) => {
    const g = canvas.grid;
    if (!g) return v;
    return Math.round(v / g) * g;
  };

  const rectPct = (ev) => {
    const s = stage();
    const box = s.getBoundingClientRect();
    return {
      x: ((ev.clientX - box.left) / box.width) * 100,
      y: ((ev.clientY - box.top) / box.height) * 100,
      box,
    };
  };

  const applyGeometry = (e) => {
    const node = stageWrap.querySelector(`.cv-el[data-id="${e.id}"]`);
    if (!node) return;
    node.style.left = e.x + '%';
    node.style.top = e.y + '%';
    node.style.width = e.w + '%';
    node.style.height = e.h + '%';
  };

  const showGuides = (e) => {
    const s = stage();
    s.querySelectorAll('.cv-guide').forEach((g) => g.remove());

    const centreX = e.x + e.w / 2;
    const centreY = e.y + e.h / 2;
    const near = (a, b) => Math.abs(a - b) < 0.8;

    const addGuide = (vertical, at) => {
      const g = document.createElement('span');
      g.className = 'cv-guide ' + (vertical ? 'cv-guide-v' : 'cv-guide-h');
      if (vertical) g.style.left = at + '%';
      else g.style.top = at + '%';
      s.appendChild(g);
    };

    // Canvas centre lines.
    if (near(centreX, 50)) addGuide(true, 50);
    if (near(centreY, 50)) addGuide(false, 50);

    // Edges shared with any other element.
    canvas.elements.forEach((o) => {
      if (o.id === e.id) return;
      [o.x, o.x + o.w, o.x + o.w / 2].forEach((v) => {
        if (near(e.x, v) || near(e.x + e.w, v) || near(centreX, v)) addGuide(true, v);
      });
      [o.y, o.y + o.h, o.y + o.h / 2].forEach((v) => {
        if (near(e.y, v) || near(e.y + e.h, v) || near(centreY, v)) addGuide(false, v);
      });
    });
  };

  const showReadout = (e, text) => {
    const s = stage();
    let r = s.querySelector('.cv-readout');
    if (!r) {
      r = document.createElement('span');
      r.className = 'cv-readout';
      s.appendChild(r);
    }
    r.textContent = text;
    r.style.left = e.x + '%';
    r.style.top = `calc(${e.y}% - 24px)`;
  };

  const endGesture = () => {
    const s = stage();
    if (s) {
      s.querySelectorAll('.cv-guide, .cv-readout').forEach((n) => n.remove());
    }
    pushState();
    if (selectedId) {
      clearDecorations();
      const node = stageWrap.querySelector(`.cv-el[data-id="${selectedId}"]`);
      if (node) { node.classList.add('is-selected'); placeHandles(node); }
    }
    renderProps();
  };

  stageWrap.addEventListener('pointerdown', (ev) => {
    const handle = ev.target.closest('.cv-handle');
    const node = ev.target.closest('.cv-el');

    if (!handle && !node) { select(null); return; }

    ev.preventDefault();
    const target = handle ? el(handle.dataset.for) : el(node.dataset.id);
    if (!target) return;

    if (!handle) select(target.id);

    const start = rectPct(ev);
    const origin = { x: target.x, y: target.y, w: target.w, h: target.h };
    const mode = handle ? handle.dataset.h : 'move';

    ev.target.setPointerCapture?.(ev.pointerId);

    const onMove = (mv) => {
      const now = rectPct(mv);
      const dx = now.x - start.x;
      const dy = now.y - start.y;

      if (mode === 'move') {
        target.x = snap(origin.x + dx);
        target.y = snap(origin.y + dy);
        showReadout(target, `${Math.round(target.x)}% , ${Math.round(target.y)}%`);
      } else {
        // Each handle moves its own edges; the opposite edge stays anchored.
        if (mode.includes('e')) target.w = Math.max(1, snap(origin.w + dx));
        if (mode.includes('s')) target.h = Math.max(1, snap(origin.h + dy));
        if (mode.includes('w')) {
          const w = Math.max(1, snap(origin.w - dx));
          target.x = snap(origin.x + (origin.w - w));
          target.w = w;
        }
        if (mode.includes('n')) {
          const h = Math.max(1, snap(origin.h - dy));
          target.y = snap(origin.y + (origin.h - h));
          target.h = h;
        }
        showReadout(target, `${Math.round(target.w)}% × ${Math.round(target.h)}%`);
      }

      applyGeometry(target);
      showGuides(target);

      // Keep the handles glued to the element while it moves.
      stageWrap.querySelectorAll(`.cv-handle[data-for="${target.id}"]`).forEach((dot) => {
        const pos = {
          nw: [0, 0], n: [50, 0], ne: [100, 0], e: [100, 50],
          se: [100, 100], s: [50, 100], sw: [0, 100], w: [0, 50],
        }[dot.dataset.h];
        dot.style.left = `calc(${target.x + (target.w * pos[0]) / 100}% - 5.5px)`;
        dot.style.top = `calc(${target.y + (target.h * pos[1]) / 100}% - 5.5px)`;
      });
    };

    const onUp = () => {
      window.removeEventListener('pointermove', onMove);
      window.removeEventListener('pointerup', onUp);
      endGesture();
    };

    window.addEventListener('pointermove', onMove);
    window.addEventListener('pointerup', onUp);
  });

  /* ================= toolbar: add / duplicate / delete ================= */

  const newId = () => 'e' + Math.random().toString(36).slice(2, 9);

  const addElement = async (kind) => {
    const body = new FormData();
    body.set('kind', kind);
    body.set('id', newId());
    body.set('csrf_token', window.csrfToken());

    try {
      // The server owns the per-kind defaults, so a new element always matches
      // what Go would produce rather than a second set of numbers here.
      const res = await fetch('/panel/design/canvas/element', {
        method: 'POST',
        body,
        headers: { 'X-Requested-With': 'fetch', 'X-CSRF-Token': window.csrfToken() },
      });
      if (!res.ok) return;

      const e = await res.json();
      e.z = canvas.elements.reduce((m, x) => Math.max(m, x.z || 0), 0) + 1;
      canvas.elements.push(e);
      pushState();
      await renderStage();
      select(e.id);
    } catch {
      window.toast?.('!', 'err');
    }
  };

  root.querySelectorAll('[data-cv-add]').forEach((btn) => {
    btn.addEventListener('click', () => addElement(btn.dataset.cvAdd));
  });

  const duplicateSelected = () => {
    const e = el(selectedId);
    if (!e) return;
    const copy = { ...e, id: newId(), x: e.x + 3, y: e.y + 3, z: (e.z || 1) + 1 };
    canvas.elements.push(copy);
    pushState();
    renderStage().then(() => select(copy.id));
  };

  const deleteSelected = () => {
    if (!selectedId) return;
    canvas.elements = canvas.elements.filter((e) => e.id !== selectedId);
    selectedId = null;
    pushState();
    renderStage().then(() => { renderProps(); renderLayers(); });
  };

  root.querySelector('[data-cv-duplicate]')?.addEventListener('click', duplicateSelected);
  root.querySelector('[data-cv-delete]')?.addEventListener('click', deleteSelected);

  /* ================= z-order ================= */

  const bump = (dir) => {
    const e = el(selectedId);
    if (!e) return;
    const zs = canvas.elements.map((x) => x.z || 0);
    e.z = dir > 0 ? Math.max(...zs) + 1 : Math.max(0, Math.min(...zs) - 1);
    pushState();
    scheduleRender();
    renderLayers();
  };

  root.querySelector('[data-cv-front]')?.addEventListener('click', () => bump(1));
  root.querySelector('[data-cv-back]')?.addEventListener('click', () => bump(-1));

  /* ================= alignment ================= */

  root.querySelectorAll('[data-cv-align]').forEach((btn) => {
    btn.addEventListener('click', () => {
      const e = el(selectedId);
      if (!e) return;
      switch (btn.dataset.cvAlign) {
        case 'left':   e.x = 0; break;
        case 'hcenter': e.x = (100 - e.w) / 2; break;
        case 'right':  e.x = 100 - e.w; break;
        case 'top':    e.y = 0; break;
        case 'vcenter': e.y = (100 - e.h) / 2; break;
        case 'bottom': e.y = 100 - e.h; break;
        case 'stretch': e.x = 0; e.w = 100; break;
      }
      applyGeometry(e);
      endGesture();
    });
  });

  /* ================= keyboard ================= */

  document.addEventListener('keydown', (ev) => {
    if (!selectedId) return;
    // Never hijack typing in the property panel.
    if (ev.target.matches('input, textarea, select')) return;

    const e = el(selectedId);
    if (!e) return;
    const step = ev.shiftKey ? 5 : 1;

    switch (ev.key) {
      case 'ArrowLeft':  e.x -= step; break;
      case 'ArrowRight': e.x += step; break;
      case 'ArrowUp':    e.y -= step; break;
      case 'ArrowDown':  e.y += step; break;
      case 'Delete':
      case 'Backspace':  ev.preventDefault(); deleteSelected(); return;
      case 'Escape':     select(null); return;
      default:
        if ((ev.ctrlKey || ev.metaKey) && ev.key.toLowerCase() === 'd') {
          ev.preventDefault();
          duplicateSelected();
        }
        return;
    }

    ev.preventDefault();
    applyGeometry(e);
    endGesture();
  });

  /* ================= property panel ================= */

  const FIELDS = [
    { k: 'textFa', t: 'text', kinds: 'heading text button icon price' },
    { k: 'textEn', t: 'text', kinds: 'heading text button icon price' },
    { k: 'link', t: 'text', kinds: 'heading text button image icon price shape' },
    { k: 'x', t: 'num', min: -50, max: 100, step: 0.5 },
    { k: 'y', t: 'num', min: -50, max: 100, step: 0.5 },
    { k: 'w', t: 'num', min: 1, max: 100, step: 0.5 },
    { k: 'h', t: 'num', min: 1, max: 100, step: 0.5 },
    { k: 'rotate', t: 'range', min: -180, max: 180, step: 1 },
    { k: 'fontSize', t: 'range', min: 0.5, max: 30, step: 0.25, kinds: 'heading text button icon price' },
    { k: 'weight', t: 'range', min: 100, max: 900, step: 100, kinds: 'heading text button price' },
    { k: 'lineHeight', t: 'range', min: 0.8, max: 3, step: 0.05, kinds: 'heading text' },
    { k: 'letter', t: 'range', min: -2, max: 8, step: 0.1, kinds: 'heading text button price' },
    { k: 'opacity', t: 'range', min: 0.05, max: 1, step: 0.05 },
    { k: 'radius', t: 'range', min: 0, max: 100, step: 1 },
    { k: 'borderW', t: 'range', min: 0, max: 12, step: 0.5 },
    { k: 'shadow', t: 'range', min: 0, max: 5, step: 1 },
    { k: 'blur', t: 'range', min: 0, max: 30, step: 1 },
  ];

  const labelFor = (key) => propsPanel.dataset['l' + key.charAt(0).toUpperCase() + key.slice(1)] || key;

  const renderProps = () => {
    const e = el(selectedId);
    const body = propsPanel.querySelector('[data-cv-props-body]');
    const none = propsPanel.querySelector('[data-cv-props-none]');

    if (!e) {
      body.hidden = true;
      none.hidden = false;
      return;
    }
    body.hidden = false;
    none.hidden = true;

    // Show only the controls that apply to this element kind.
    body.querySelectorAll('[data-for-kinds]').forEach((n) => {
      n.hidden = !n.dataset.forKinds.split(' ').includes(e.kind);
    });

    body.querySelectorAll('[data-p]').forEach((input) => {
      const key = input.dataset.p;
      const val = e[key];
      if (input.type === 'checkbox') input.checked = Boolean(val);
      else input.value = val === undefined || val === null ? '' : val;

      const out = body.querySelector(`[data-out-for="${key}"]`);
      if (out) out.textContent = input.value;
    });

    // The label sits in the panel heading, outside the body element.
    const kindLabel = propsPanel.querySelector('[data-cv-kind-label]');
    if (kindLabel) kindLabel.textContent = kindName(e.kind);
  };

  // Kind labels are carried as data-ktext, data-kheading … so the dataset key
  // is the kind appended verbatim, with no camel-casing.
  const kindName = (kind) => propsPanel.dataset['k' + kind] || kind;

  propsPanel.addEventListener('input', (ev) => {
    const input = ev.target.closest('[data-p]');
    if (!input) return;

    const e = el(selectedId);
    if (!e) return;

    const key = input.dataset.p;
    let val;
    if (input.type === 'checkbox') val = input.checked;
    else if (input.type === 'range' || input.type === 'number') val = parseFloat(input.value) || 0;
    else val = input.value;

    e[key] = val;

    const out = propsPanel.querySelector(`[data-out-for="${key}"]`);
    if (out) out.textContent = input.value;

    pushState();

    // Geometry is cheap to patch locally; anything else needs the Go renderer.
    if (['x', 'y', 'w', 'h'].includes(key)) {
      applyGeometry(e);
      clearDecorations();
      const node = stageWrap.querySelector(`.cv-el[data-id="${e.id}"]`);
      if (node) { node.classList.add('is-selected'); placeHandles(node); }
    } else {
      scheduleRender();
    }
  });

  // Image picker for the selected element and for the canvas background.
  propsPanel.querySelector('[data-cv-pick-image]')?.addEventListener('click', () => {
    window.pickImage?.((path) => {
      const e = el(selectedId);
      if (!e) return;
      e.image = path;
      pushState();
      scheduleRender();
    });
  });

  root.querySelector('[data-cv-pick-bg]')?.addEventListener('click', () => {
    window.pickImage?.((path) => {
      canvas.bgImage = path;
      pushState();
      scheduleRender();
    });
  });

  root.querySelector('[data-cv-clear-bg]')?.addEventListener('click', () => {
    canvas.bgImage = '';
    pushState();
    scheduleRender();
  });

  /* ================= canvas-level settings ================= */

  root.querySelectorAll('[data-c]').forEach((input) => {
    const key = input.dataset.c;
    if (canvas[key] !== undefined) {
      if (input.type === 'checkbox') input.checked = Boolean(canvas[key]);
      else input.value = canvas[key];
    }

    input.addEventListener('input', () => {
      let val;
      if (input.type === 'checkbox') val = input.checked;
      else if (input.type === 'range' || input.type === 'number') val = parseFloat(input.value) || 0;
      else val = input.value;

      canvas[key] = val;
      const out = root.querySelector(`[data-out-for="${key}"]`);
      if (out) out.textContent = input.value;

      pushState();
      scheduleRender();
    });
  });

  /* ================= layer list ================= */

  const renderLayers = () => {
    if (!layerList) return;
    layerList.innerHTML = '';

    // Topmost first, matching how a layers panel is normally read.
    [...canvas.elements]
      .sort((a, b) => (b.z || 0) - (a.z || 0))
      .forEach((e) => {
        const row = document.createElement('button');
        row.type = 'button';
        row.className = 'cv-layer' + (e.id === selectedId ? ' is-active' : '');
        row.dataset.id = e.id;

        const label = kindName(e.kind);
        const text = (e.textFa || e.textEn || '').trim();
        row.textContent = text ? `${label} — ${text.slice(0, 22)}` : label;

        row.addEventListener('click', () => select(e.id));
        layerList.appendChild(row);
      });
  };

  /* ================= save ================= */

  root.querySelector('[data-cv-save]')?.addEventListener('click', async () => {
    const form = root.closest('form') || document.querySelector('[data-canvas-form]');
    jsonField.value = JSON.stringify(canvas);
    if (form) {
      dirty = false;
      form.submit();
    }
  });

  window.addEventListener('beforeunload', (ev) => {
    if (!dirty) return;
    ev.preventDefault();
    ev.returnValue = '';
  });

  /* ================= boot ================= */

  decorateStage();
  renderProps();
  renderLayers();
})();
