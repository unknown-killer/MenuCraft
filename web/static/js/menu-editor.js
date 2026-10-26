/* Menu manager: category/item dialogs, image picker, drag-to-reorder. */
(() => {
  'use strict';

  /* ==================================================================
     Image picker — shared by the theme editor, menu dialogs and settings
     ================================================================== */

  const picker = document.querySelector('[data-picker]');
  let pickerTarget = null;   // id of the hidden input currently being filled
  let pickerCallback = null; // set instead of pickerTarget for ad-hoc pickers

  const setImage = (fieldId, path) => {
    const input = document.getElementById(fieldId) || document.querySelector(`[name="${fieldId}"]`);
    if (input) {
      input.value = path;
      input.dispatchEvent(new Event('input', { bubbles: true }));
    }

    const preview = document.querySelector(`[data-img-preview="${fieldId}"]`);
    if (preview) {
      preview.classList.toggle('is-empty', !path);
      preview.replaceChildren();
      if (path) {
        const image = document.createElement('img');
        image.src = path;
        image.alt = '';
        preview.append(image);
      }
    }
  };

  const loadPickerGrid = async () => {
    const grid = picker?.querySelector('[data-picker-grid]');
    if (!grid) return;

    try {
      const res = await fetch('/panel/media?format=json', { headers: { Accept: 'application/json' } });
      if (!res.ok) return;

      const { items } = await res.json();
      grid.innerHTML = '';

      if (!items || !items.length) {
        const empty = document.createElement('p');
        empty.className = 'muted picker-empty';
        empty.textContent = picker.dataset.emptyText || '—';
        grid.appendChild(empty);
        return;
      }

      items.forEach((m) => {
        const btn = document.createElement('button');
        btn.type = 'button';
        btn.className = 'picker-item';
        btn.title = m.name || '';

        const img = document.createElement('img');
        img.src = m.path;
        img.alt = m.name || '';
        img.loading = 'lazy';
        btn.appendChild(img);

        btn.addEventListener('click', () => {
          if (pickerCallback) pickerCallback(m.path);
          else if (pickerTarget) setImage(pickerTarget, m.path);
          picker.close();
        });

        grid.appendChild(btn);
      });
    } catch {
      /* leave whatever is already on screen */
    }
  };

  document.addEventListener('click', (e) => {
    const open = e.target.closest('[data-pick-image]');
    if (open && picker) {
      pickerTarget = open.dataset.pickImage;
      pickerCallback = null;
      loadPickerGrid();
      picker.showModal();
      return;
    }

    const clear = e.target.closest('[data-clear-image]');
    if (clear) setImage(clear.dataset.clearImage, '');
  });

  // Exposed so the page designer can reuse the picker and previews.
  window.setImageField = setImage;

  window.pickImage = (onChoose) => {
    if (!picker) return;
    pickerTarget = null;
    pickerCallback = onChoose;
    loadPickerGrid();
    picker.showModal();
  };

  /* ---------- upload straight from the picker ---------- */
  const dropInput = picker?.querySelector('[data-drop-input]');
  const dropZone = picker?.querySelector('[data-drop]');
  const progress = picker?.querySelector('[data-drop-progress]');

  const upload = async (file) => {
    if (!file) return;
    if (progress) progress.hidden = false;

    const body = new FormData();
    body.set('file', file);
    body.set('csrf_token', window.csrfToken());

    try {
      const res = await fetch('/panel/media/upload', {
        method: 'POST',
        body,
        headers: { 'X-Requested-With': 'fetch', 'X-CSRF-Token': window.csrfToken() },
      });
      const data = await res.json();

      if (!res.ok) {
        window.toast?.(data.error || 'error', 'err');
        return;
      }
      if (pickerCallback) pickerCallback(data.path);
      else if (pickerTarget) setImage(pickerTarget, data.path);
      await loadPickerGrid();
      picker.close();
    } catch {
      window.toast?.('error', 'err');
    } finally {
      if (progress) progress.hidden = true;
    }
  };

  dropInput?.addEventListener('change', () => upload(dropInput.files?.[0]));

  ['dragenter', 'dragover'].forEach((ev) =>
    dropZone?.addEventListener(ev, (e) => { e.preventDefault(); dropZone.classList.add('is-over'); }));
  ['dragleave', 'drop'].forEach((ev) =>
    dropZone?.addEventListener(ev, (e) => { e.preventDefault(); dropZone.classList.remove('is-over'); }));

  dropZone?.addEventListener('drop', (e) => upload(e.dataTransfer?.files?.[0]));

  /* ==================================================================
     Category & item dialogs
     ================================================================== */

  const catModal = document.querySelector('[data-cat-modal]');
  const itemModal = document.querySelector('[data-item-modal]');

  const fillForm = (modal, data) => {
    modal.querySelectorAll('[data-f]').forEach((el) => {
      const key = el.dataset.f;
      const val = data?.[key];

      if (el.type === 'checkbox') {
        el.checked = Boolean(val);
      } else if (val === undefined || val === null || val === 0) {
        // A zero price or calorie count reads better as an empty field.
        el.value = '';
      } else if (el.hasAttribute('data-price')) {
        el.value = Number(val).toLocaleString('en-US');
      } else {
        el.value = String(val);
      }
    });

    // Toggles live outside [data-f]; set them from the payload explicitly.
    ['is_visible', 'is_available', 'is_featured'].forEach((name) => {
      const box = modal.querySelector(`[name="${name}"]`);
      if (box) box.checked = data ? Boolean(data[name]) : true;
    });
    if (!data) {
      const featured = modal.querySelector('[name="is_featured"]');
      if (featured) featured.checked = false;
    }

    // Tags
    modal.querySelectorAll('[data-tag]').forEach((box) => {
      box.checked = Array.isArray(data?.tags) && data.tags.includes(box.dataset.tag);
    });

    // Image preview
    const imgField = modal.querySelector('[data-f="image_path"]');
    if (imgField) setImage(imgField.id || 'image_path', data?.image_path || '');
  };

  document.addEventListener('click', (e) => {
    const catBtn = e.target.closest('[data-open-cat]');
    if (catBtn && catModal) {
      const raw = catBtn.dataset.cat;
      const data = raw ? JSON.parse(raw) : null;

      catModal.querySelector('[data-cat-id]').value = data?.id || '';
      fillForm(catModal, data);
      catModal.showModal();
      return;
    }

    const itemBtn = e.target.closest('[data-open-item]');
    if (itemBtn && itemModal) {
      const raw = itemBtn.dataset.item;
      const data = raw ? JSON.parse(raw) : null;

      itemModal.querySelector('[data-item-id]').value = data?.id || '';
      fillForm(itemModal, data);

      // Preselect the category the button belongs to when adding a new item.
      // fillForm blanks every [data-f] field, which would otherwise leave the
      // select with nothing chosen; fall back to the first category so the
      // dialog always opens in a submittable state.
      const select = itemModal.querySelector('[name="category_id"]');
      const wanted = data?.category_id || itemBtn.dataset.openItem;
      if (select) {
        if (wanted) select.value = String(wanted);
        if (!select.value && select.options.length) select.selectedIndex = 0;
      }

      itemModal.showModal();
      return;
    }

    if (e.target.closest('[data-close-modal]')) {
      e.target.closest('dialog')?.close();
    }
  });

  /* ---------- price fields: accept Persian digits and separators ---------- */
  const toAscii = (s) => s.replace(/[۰-۹]/g, (d) => '۰۱۲۳۴۵۶۷۸۹'.indexOf(d))
                          .replace(/[٠-٩]/g, (d) => '٠١٢٣٤٥٦٧٨٩'.indexOf(d));

  document.querySelectorAll('[data-price]').forEach((input) => {
    if (input.name === 'price') input.addEventListener('input', () => {
      const pending = input.form?.querySelector('[name="price_pending"]');
      if (pending && Number(toAscii(input.value).replace(/[,٬\s]/g, '')) > 0) pending.checked = false;
    });
    input.addEventListener('blur', () => {
      const digits = toAscii(input.value).replace(/\D/g, '');
      input.value = digits ? Number(digits).toLocaleString('en-US') : '';
    });
  });

  /* ==================================================================
     Drag to reorder
     ================================================================== */

  document.querySelectorAll('[data-sortable]').forEach((list) => {
    let dragged = null;

    list.querySelectorAll(':scope > [data-id]').forEach((row) => {
      const handle = row.querySelector('.drag');
      if (!handle) return;

      // Only the handle starts a drag, so text stays selectable.
      handle.addEventListener('mousedown', () => { row.draggable = true; });
      handle.addEventListener('touchstart', () => { row.draggable = true; }, { passive: true });
      row.addEventListener('dragend', () => { row.draggable = false; });

      row.addEventListener('dragstart', (e) => {
        dragged = row;
        row.classList.add('dragging');
        e.dataTransfer.effectAllowed = 'move';
        // Firefox needs data set for the drag to begin at all.
        e.dataTransfer.setData('text/plain', row.dataset.id);
      });

      row.addEventListener('dragend', () => {
        row.classList.remove('dragging');
        dragged = null;
        persistOrder(list);
      });

      row.addEventListener('dragover', (e) => {
        e.preventDefault();
        if (!dragged || dragged === row) return;

        const box = row.getBoundingClientRect();
        const after = (e.clientY - box.top) > box.height / 2;
        row.parentNode.insertBefore(dragged, after ? row.nextSibling : row);
      });
    });
  });

  const persistOrder = async (list) => {
    const url = list.dataset.reorderUrl;
    if (!url) return;

    const ids = [...list.querySelectorAll(':scope > [data-id]')].map((el) => el.dataset.id);
    const body = new FormData();
    body.set('order', ids.join(','));
    body.set('csrf_token', window.csrfToken());

    try {
      const res = await fetch(url, {
        method: 'POST',
        body,
        headers: { 'X-Requested-With': 'fetch', 'X-CSRF-Token': window.csrfToken() },
      });
      if (res.ok) window.toast?.('✓');
    } catch {
      /* the visual order still reflects the intent; a reload will correct it */
    }
  };
})();
