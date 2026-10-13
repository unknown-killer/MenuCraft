/*
  Theme editor.

  The live preview does not reimplement the CSS generator in JavaScript. On
  every change the form is posted to /panel/theme/css, which returns the same
  variables the server would render for real, and those are injected into the
  preview iframe. One source of truth, instant feedback.
*/
(() => {
  'use strict';

  const root = document.querySelector('[data-editor]');
  if (!root) return;

  const form = root.querySelector('[data-theme-form]');
  const frame = root.querySelector('[data-preview-frame]');
  const stage = root.querySelector('.preview-stage');
  const presetField = root.querySelector('[data-preset-field]');

  /* ---------- range outputs ---------- */
  const syncOutput = (input) => {
    const out = form.querySelector(`[data-out-for="${input.name}"]`);
    if (out) out.textContent = input.value + (input.dataset.unit || '');
  };
  form.querySelectorAll('input[type="range"]').forEach(syncOutput);

  /* ---------- colour picker <-> hex text ---------- */
  form.querySelectorAll('[data-color-sync]').forEach((picker) => {
    const hex = picker.parentElement.querySelector('[data-color-hex]');
    if (!hex) return;

    picker.addEventListener('input', () => { hex.value = picker.value.toUpperCase(); });

    hex.addEventListener('input', () => {
      const v = hex.value.trim();
      // Only push valid colours through; partial typing must not reset the picker.
      if (/^#([0-9a-f]{3}|[0-9a-f]{6})$/i.test(v)) {
        picker.value = v;
        picker.dispatchEvent(new Event('input', { bubbles: true }));
      }
    });
  });

  /* ---------- show only the fields for the chosen background type ---------- */
  const syncBackgroundFields = () => {
    const type = form.querySelector('[name="b_type"]:checked')?.value || 'solid';
    form.querySelectorAll('[data-bg-when]').forEach((el) => {
      el.hidden = !el.dataset.bgWhen.split(' ').includes(type);
    });
  };
  syncBackgroundFields();

  /* ---------- live preview ---------- */
  let timer = null;
  let inFlight = false;
  let queued = false;

  const applyToFrame = (payload) => {
    const doc = frame.contentDocument;
    if (!doc) return;

    const style = doc.getElementById('live-vars');
    if (style) style.textContent = payload.css;

    if (doc.body && payload.bodyClass) {
      doc.body.className = 'menu-body ' + payload.bodyClass;
    }
  };

  const pushPreview = async () => {
    if (inFlight) { queued = true; return; }
    inFlight = true;

    try {
      const res = await fetch('/panel/theme/css', {
        method: 'POST',
        body: new FormData(form),
        headers: { 'X-Requested-With': 'fetch', 'X-CSRF-Token': window.csrfToken() },
      });
      if (res.ok) applyToFrame(await res.json());
    } catch {
      // A dropped preview refresh is not worth interrupting the user for.
    } finally {
      inFlight = false;
      if (queued) { queued = false; pushPreview(); }
    }
  };

  const schedulePreview = () => {
    clearTimeout(timer);
    timer = setTimeout(pushPreview, 130);
  };

  form.addEventListener('input', (e) => {
    if (e.target.type === 'range') syncOutput(e.target);
    if (e.target.name === 'b_type') syncBackgroundFields();

    // Any manual edit means this is no longer a stock preset.
    if (presetField && e.target.name !== 'preset') {
      presetField.value = 'custom';
      root.querySelectorAll('.preset.is-active').forEach((p) => p.classList.remove('is-active'));
    }
    schedulePreview();
  });

  form.addEventListener('change', schedulePreview);

  frame.addEventListener('load', pushPreview);

  /* ---------- presets ---------- */
  root.querySelectorAll('[data-preset]').forEach((btn) => {
    btn.addEventListener('click', async () => {
      const key = btn.dataset.preset;

      const body = new FormData();
      body.set('preset', key);
      body.set('preset_only', '1');
      body.set('preview_locale', form.querySelector('[name="preview_locale"]')?.value || 'fa');

      try {
        const res = await fetch('/panel/theme/css', {
          method: 'POST',
          body,
          headers: { 'X-Requested-With': 'fetch', 'X-CSRF-Token': window.csrfToken() },
        });
        if (!res.ok) return;
        applyToFrame(await res.json());
      } catch {
        return;
      }

      root.querySelectorAll('.preset').forEach((p) => p.classList.toggle('is-active', p === btn));
      if (presetField) presetField.value = key;

      // The controls still show the old values; saving now would contradict the
      // preview, so the form is submitted with preset_only and the page reloads
      // with every control repopulated from the preset.
      const marker = document.createElement('input');
      marker.type = 'hidden';
      marker.name = 'preset_only';
      marker.value = '1';
      form.appendChild(marker);
      form.submit();
    });
  });

  /* ---------- device switcher ---------- */
  root.querySelectorAll('[data-device]').forEach((radio) => {
    radio.addEventListener('change', () => {
      if (radio.checked) stage.dataset.stage = radio.value;
    });
  });

  /* ---------- reset ---------- */
  root.querySelector('[data-theme-reset]')?.addEventListener('click', () => {
    document.getElementById('theme-reset-form')?.submit();
  });

  /* ---------- warn before leaving with unsaved changes ---------- */
  let dirty = false;
  form.addEventListener('input', () => { dirty = true; });
  form.addEventListener('submit', () => { dirty = false; });

  window.addEventListener('beforeunload', (e) => {
    if (!dirty) return;
    e.preventDefault();
    e.returnValue = '';
  });
})();
