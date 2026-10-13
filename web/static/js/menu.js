/* Public menu: category filtering, search, and scroll-spy on the nav. */
(() => {
  'use strict';

  /* ---------- light / dark toggle ---------- */
  const modeBtn = document.querySelector('[data-mode-toggle]');
  if (modeBtn) {
    const icon = modeBtn.querySelector('[data-mode-icon]');

    const currentMode = () => {
      const set = document.documentElement.getAttribute('data-mode');
      if (set) return set;
      return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    };

    const paint = () => { if (icon) icon.textContent = currentMode() === 'dark' ? '☀' : '☾'; };
    paint();

    modeBtn.addEventListener('click', () => {
      const next = currentMode() === 'dark' ? 'light' : 'dark';
      document.documentElement.setAttribute('data-mode', next);
      try { localStorage.setItem('menu-mode', next); } catch { /* private mode */ }
      paint();
    });
  }

  const nav = document.querySelector('[data-nav]');
  const search = document.querySelector('[data-search]');
  const cats = [...document.querySelectorAll('[data-cat]')];
  const cards = [...document.querySelectorAll('.m-card')];
  const noResults = document.querySelector('[data-noresults]');

  /* Stagger the entrance animation without inline styles per card. */
  cards.forEach((card, i) => card.style.setProperty('--i', String(Math.min(i, 16))));

  let activeFilter = 'all';

  const normalise = (s) =>
    (s || '')
      .toLowerCase()
      // Fold Arabic yeh/kaf onto their Persian forms so both spellings match.
      .replace(/ي/g, 'ی')
      .replace(/ك/g, 'ک')
      .replace(/‌/g, ' ')
      .trim();

  const apply = () => {
    const q = normalise(search?.value);
    let visible = 0;

    cats.forEach((cat) => {
      const inFilter = activeFilter === 'all' || cat.dataset.cat === activeFilter;
      let shown = 0;

      cat.querySelectorAll('.m-card').forEach((card) => {
        const haystack = normalise(card.dataset.name + ' ' + card.dataset.desc);
        const match = !q || haystack.includes(q);
        card.hidden = !(inFilter && match);
        if (!card.hidden) shown++;
      });

      cat.hidden = shown === 0;
      visible += shown;
    });

    if (noResults) noResults.hidden = visible > 0;
  };

  /* ---------- category buttons ---------- */
  nav?.addEventListener('click', (e) => {
    const btn = e.target.closest('[data-filter]');
    if (!btn) return;

    activeFilter = btn.dataset.filter;
    nav.querySelectorAll('[data-filter]').forEach((b) => b.classList.toggle('is-active', b === btn));

    if (search) search.value = '';
    apply();

    if (activeFilter !== 'all') {
      document.getElementById(activeFilter)?.scrollIntoView({ behavior: 'smooth', block: 'start' });
    } else {
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }
  });

  /* ---------- search ---------- */
  let t = null;
  search?.addEventListener('input', () => {
    clearTimeout(t);
    t = setTimeout(() => {
      // Typing searches across the whole menu, not just the active category.
      if (search.value.trim()) {
        activeFilter = 'all';
        nav?.querySelectorAll('[data-filter]').forEach((b) =>
          b.classList.toggle('is-active', b.dataset.filter === 'all'));
      }
      apply();
    }, 140);
  });

  /* ---------- scroll spy ---------- */
  if (nav && cats.length && 'IntersectionObserver' in window) {
    const spy = new IntersectionObserver((entries) => {
      if (search?.value.trim() || activeFilter !== 'all') return;

      const top = entries
        .filter((en) => en.isIntersecting)
        .sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top)[0];
      if (!top) return;

      const id = top.target.dataset.cat;
      nav.querySelectorAll('[data-filter]').forEach((b) =>
        b.classList.toggle('is-active', b.dataset.filter === id));
    }, { rootMargin: '-20% 0px -70% 0px', threshold: 0 });

    cats.forEach((c) => spy.observe(c));
  }
})();
