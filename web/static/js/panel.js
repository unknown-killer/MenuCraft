/* Shared panel behaviours: navigation, flashes, confirmations, copy buttons. */
(() => {
  'use strict';

  /* ---------- mobile sidebar ---------- */
  const sidebar = document.querySelector('[data-sidebar]');
  const scrim = document.querySelector('[data-nav-close]');

  const setNav = (open) => {
    if (!sidebar) return;
    sidebar.classList.toggle('is-open', open);
    if (scrim) scrim.hidden = !open;
  };

  document.querySelector('[data-nav-toggle]')?.addEventListener('click', () => {
    setNav(!sidebar?.classList.contains('is-open'));
  });
  scrim?.addEventListener('click', () => setNav(false));

  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') setNav(false);
  });

  /* ---------- dismissible flashes ---------- */
  document.addEventListener('click', (e) => {
    const btn = e.target.closest('[data-dismiss]');
    if (btn) btn.closest('.flash')?.remove();
  });

  document.querySelectorAll('.flash').forEach((el) => {
    setTimeout(() => {
      el.style.transition = 'opacity .4s, transform .4s';
      el.style.opacity = '0';
      el.style.transform = 'translateY(-8px)';
      setTimeout(() => el.remove(), 420);
    }, 5200);
  });

  /* ---------- confirm before destructive submits ---------- */
  document.addEventListener('submit', (e) => {
    const msg = e.target.dataset?.confirm;
    if (msg && !window.confirm(msg)) e.preventDefault();
  });

  document.addEventListener('click', (e) => {
    const btn = e.target.closest('[data-confirm]:not(form)');
    if (btn && btn.tagName === 'BUTTON' && !window.confirm(btn.dataset.confirm)) {
      e.preventDefault();
      e.stopImmediatePropagation();
    }
  }, true);

  /* ---------- copy to clipboard ---------- */
  document.addEventListener('click', async (e) => {
    const btn = e.target.closest('[data-copy]');
    if (!btn) return;

    const src = document.querySelector('[data-copy-src]');
    if (!src) return;

    const text = new URL(src.textContent.trim(), location.origin).href;
    try {
      await navigator.clipboard.writeText(text);
    } catch {
      // Clipboard API needs a secure context; fall back to a temp selection.
      const ta = document.createElement('textarea');
      ta.value = text;
      ta.style.position = 'fixed';
      ta.style.opacity = '0';
      document.body.appendChild(ta);
      ta.select();
      try { document.execCommand('copy'); } catch { /* nothing else to try */ }
      ta.remove();
    }

    const original = btn.textContent;
    btn.textContent = btn.dataset.copied || 'OK';
    setTimeout(() => { btn.textContent = original; }, 1600);
  });

  /* ---------- auto-submit single file uploads ---------- */
  document.querySelectorAll('[data-auto-submit]').forEach((input) => {
    input.addEventListener('change', () => {
      if (input.files?.length) input.form?.submit();
    });
  });

  /* ---------- toasts (shared helper) ---------- */
  const host = document.querySelector('[data-toasts]');

  window.toast = (message, kind = 'ok') => {
    if (!host) return;
    const el = document.createElement('div');
    el.className = `toast toast-${kind}`;
    el.textContent = message;
    host.appendChild(el);
    setTimeout(() => {
      el.style.transition = 'opacity .35s, transform .35s';
      el.style.opacity = '0';
      el.style.transform = 'translateY(8px)';
      setTimeout(() => el.remove(), 360);
    }, 2800);
  };

  /* ---------- CSRF helper for fetch() calls ---------- */
  window.csrfToken = () => {
    const m = document.cookie.match(/(?:^|;\s*)menu_csrf=([^;]+)/);
    return m ? decodeURIComponent(m[1]) : '';
  };
})();
