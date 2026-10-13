/*
  Page designer: drag the sections into order, and edit custom blocks.

  Reordering saves immediately and refreshes the preview iframe, so the layout
  the owner sees is always the layout that is stored.
*/
(() => {
  'use strict';

  const root = document.querySelector('[data-designer]');
  if (!root) return;

  const list = root.querySelector('[data-section-list]');
  const frame = root.querySelector('[data-preview-frame]');
  const stage = root.querySelector('.preview-stage');
  const modal = document.querySelector('[data-block-modal]');

  const refreshPreview = () => {
    if (!frame) return;
    // Cache-bust so the iframe re-fetches rather than restoring from cache.
    const base = frame.src.split('?')[0];
    frame.src = `${base}?t=${Date.now()}`;
  };

  /* ================= drag to reorder sections ================= */

  let dragging = null;

  const rows = () => [...list.querySelectorAll('.sect-row')];

  const saveOrder = async () => {
    const order = rows().map((r) => r.dataset.key).join(',');
    const body = new FormData();
    body.set('order', order);
    body.set('csrf_token', window.csrfToken());

    try {
      const res = await fetch(list.dataset.reorderUrl, {
        method: 'POST',
        body,
        headers: { 'X-Requested-With': 'fetch', 'X-CSRF-Token': window.csrfToken() },
      });
      if (res.ok) {
        window.toast?.('✓');
        refreshPreview();
      } else {
        window.toast?.('!', 'err');
      }
    } catch {
      window.toast?.('!', 'err');
    }
  };

  const wireRow = (row) => {
    const handle = row.querySelector('.drag');
    if (!handle) return;

    // Only the handle starts a drag, so buttons and text stay usable.
    handle.addEventListener('mousedown', () => { row.draggable = true; });
    handle.addEventListener('touchstart', () => { row.draggable = true; }, { passive: true });

    row.addEventListener('dragstart', (e) => {
      dragging = row;
      row.classList.add('dragging');
      e.dataTransfer.effectAllowed = 'move';
      e.dataTransfer.setData('text/plain', row.dataset.key);
    });

    row.addEventListener('dragend', () => {
      row.classList.remove('dragging');
      row.draggable = false;
      dragging = null;
      saveOrder();
    });

    row.addEventListener('dragover', (e) => {
      e.preventDefault();
      if (!dragging || dragging === row) return;
      const box = row.getBoundingClientRect();
      const after = (e.clientY - box.top) > box.height / 2;
      row.parentNode.insertBefore(dragging, after ? row.nextSibling : row);
    });
  };

  rows().forEach(wireRow);

  /* Keyboard reordering — drag-and-drop alone is not accessible. */
  list.addEventListener('keydown', (e) => {
    if (!e.altKey || (e.key !== 'ArrowUp' && e.key !== 'ArrowDown')) return;
    const row = e.target.closest('.sect-row');
    if (!row) return;

    e.preventDefault();
    if (e.key === 'ArrowUp' && row.previousElementSibling) {
      row.parentNode.insertBefore(row, row.previousElementSibling);
    } else if (e.key === 'ArrowDown' && row.nextElementSibling) {
      row.parentNode.insertBefore(row.nextElementSibling, row);
    }
    row.querySelector('.drag')?.focus();
    saveOrder();
  });

  // Make the handles reachable by keyboard.
  list.querySelectorAll('.drag').forEach((h) => {
    h.tabIndex = 0;
    h.setAttribute('role', 'button');
  });

  /* ================= block dialog ================= */

  const KIND_TITLES = {
    add: modal.dataset.titleNew,
    edit: modal.dataset.titleEdit,
  };

  // Show only the fields that apply to the selected block kind.
  const syncKindFields = (kind) => {
    modal.querySelectorAll('[data-when]').forEach((el) => {
      el.hidden = !el.dataset.when.split(' ').includes(kind);
    });
  };

  const syncColorField = () => {
    const color = modal.querySelector('[name="s_color"]')?.value;
    modal.querySelectorAll('[data-when-color]').forEach((el) => {
      el.hidden = el.dataset.whenColor !== color;
    });
  };

  const setField = (name, value) => {
    const el = modal.querySelector(`[name="${name}"]`);
    if (!el) return;
    if (el.type === 'checkbox') el.checked = Boolean(value);
    else el.value = value === undefined || value === null ? '' : String(value);

    if (el.type === 'range') {
      const out = modal.querySelector(`[data-out-for="${name}"]`);
      if (out) out.textContent = el.value + (el.dataset.unit || '');
    }
  };

  /* --- gallery strip --- */
  const galleryInput = () => modal.querySelector('[data-b="images"]');
  const galleryStrip = () => modal.querySelector('[data-gallery-strip]');

  const renderGallery = () => {
    const strip = galleryStrip();
    const input = galleryInput();
    if (!strip || !input) return;

    const paths = input.value.split(',').map((p) => p.trim()).filter(Boolean);
    strip.innerHTML = '';

    paths.forEach((p, i) => {
      const cell = document.createElement('span');
      cell.className = 'gallery-cell';

      const img = document.createElement('img');
      img.src = p;
      img.alt = '';
      cell.appendChild(img);

      const del = document.createElement('button');
      del.type = 'button';
      del.className = 'gallery-x';
      del.textContent = '×';
      del.addEventListener('click', () => {
        const next = paths.filter((_, j) => j !== i);
        input.value = next.join(',');
        renderGallery();
      });
      cell.appendChild(del);

      strip.appendChild(cell);
    });
  };

  modal.querySelector('[data-gallery-add]')?.addEventListener('click', () => {
    // Reuse the shared media picker, appending rather than replacing.
    window.pickImage?.((path) => {
      const input = galleryInput();
      const paths = input.value.split(',').map((p) => p.trim()).filter(Boolean);
      paths.push(path);
      input.value = paths.join(',');
      renderGallery();
    });
  });

  const openBlock = (kind, data) => {
    modal.querySelector('[data-block-kind]').value = kind;
    modal.querySelector('[data-block-id]').value = data?.id || '';
    modal.querySelector('[data-block-title]').textContent =
      data ? (KIND_TITLES.edit || 'Edit') : (KIND_TITLES.add || 'New');

    // Content fields
    ['title_fa', 'title_en', 'body_fa', 'body_en', 'link_url', 'link_label_fa', 'link_label_en']
      .forEach((n) => setField(n, data?.[n] || ''));

    const imgPath = data?.image_path || '';
    setField('image_path', imgPath);
    window.setImageField?.('blk_image_path', imgPath);

    const images = Array.isArray(data?.images) ? data.images.join(',') : '';
    setField('images', images);
    renderGallery();

    setField('is_visible', data ? Boolean(data.is_visible) : true);

    // Style fields — fall back to whatever the server rendered as defaults.
    const st = data?.style;
    if (st) {
      setField('s_size', st.size);
      setField('s_align', st.align);
      setField('s_color', st.color);
      setField('s_customCol', st.customCol);
      setField('s_bg', st.bg);
      setField('s_width', st.width);
      setField('s_padding', st.padding);
      setField('s_effect', st.effect);
      setField('s_ratio', st.ratio);
      setField('s_divider', st.divider);
      setField('s_weight', st.weight);
      setField('s_radius', st.radius);
      setField('s_shadow', st.shadow);
      setField('s_height', st.height);
      setField('s_columns', st.columns);
      setField('s_overlay', st.overlay);
      setField('s_italic', st.italic);
      setField('s_uppercase', st.uppercase);
      setField('s_border', st.border);
    }

    syncKindFields(kind);
    syncColorField();
    modal.showModal();
  };

  document.addEventListener('click', (e) => {
    const btn = e.target.closest('[data-open-block]');
    if (!btn) return;
    const raw = btn.dataset.block;
    openBlock(btn.dataset.openBlock, raw ? JSON.parse(raw) : null);
  });

  modal.addEventListener('change', (e) => {
    if (e.target.name === 's_color') syncColorField();
  });

  modal.addEventListener('input', (e) => {
    if (e.target.type === 'range') {
      const out = modal.querySelector(`[data-out-for="${e.target.name}"]`);
      if (out) out.textContent = e.target.value + (e.target.dataset.unit || '');
    }
  });

  /* ================= device switcher ================= */
  root.querySelectorAll('[data-device]').forEach((radio) => {
    radio.addEventListener('change', () => {
      if (radio.checked) stage.dataset.stage = radio.value;
    });
  });
})();
