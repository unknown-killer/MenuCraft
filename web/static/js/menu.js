/* Progressive customer tools: the complete menu remains usable without JS. */
(() => {
'use strict';
const $ = s => document.querySelector(s), $$ = s => [...document.querySelectorAll(s)];
const reduced = window.matchMedia('(prefers-reduced-motion: reduce)');
const mode = $('[data-mode-toggle]'), systemMode = window.matchMedia('(prefers-color-scheme: dark)');
const currentMode = () => document.documentElement.dataset.mode || (systemMode.matches ? 'dark' : 'light');
const paintMode = () => { const icon = mode?.querySelector('[data-mode-icon]'); if (icon) icon.textContent = currentMode() === 'dark' ? '☀' : '☾'; };
paintMode(); systemMode.addEventListener?.('change', paintMode);
mode?.addEventListener('click', () => {
  const next = currentMode() === 'dark' ? 'light' : 'dark';
  document.documentElement.dataset.mode = next;
  try { localStorage.setItem('menu-mode', next); } catch {}
  paintMode();
});
const nav = $('[data-nav]'), search = $('[data-search]'), cats = $$('[data-cat]'), cards = $$('.m-card');
const noResults = $('[data-noresults]'), count = $('[data-result-count]'), sort = $('[data-sort]');
const key = 'menucraft:favorites:' + (document.body.dataset.venue || location.pathname);
let favorites = new Set();
try {
  const stored = JSON.parse(localStorage.getItem(key) || '[]');
  if (Array.isArray(stored)) favorites = new Set(stored.filter(id => typeof id === 'string').slice(0, 5000));
} catch {}
let category = 'all', diet = '', onlyFavorites = false;
const normalise = s => (s || '').toLowerCase().normalize('NFKC')
  .replace(/[يى]/g, 'ی').replace(/ك/g, 'ک').replace(/[\u064B-\u065F\u0670]/g, '').replace(/[\u200c\s]+/g, ' ').trim();
const index = new Map(cards.map(c => [c, normalise(c.dataset.name + ' ' + c.dataset.desc)]));
cards.forEach((c,i) => { c.dataset.order = String(i); c.style.setProperty('--i', String(Math.min(i,16))); });
const syncNav = () => nav?.querySelectorAll('[data-filter]').forEach(b => {
  const active = b.dataset.filter === category;
  b.classList.toggle('is-active',active); b.setAttribute('aria-pressed',String(active)); b.removeAttribute('aria-current');
});
const apply = () => {
  const q = normalise(search?.value); let visible = 0;
  cats.forEach(cat => {
    const items = [...cat.querySelectorAll('.m-card')], grid = cat.querySelector('.m-grid');
    items.sort((a,b) => {
      if ((sort?.value === 'low' || sort?.value === 'high') && (a.dataset.price === "" || b.dataset.price === "")) return Number(a.dataset.price === "") - Number(b.dataset.price === "");
      const diff = Number(a.dataset.price)-Number(b.dataset.price);
      return sort?.value === 'low' ? diff : sort?.value === 'high' ? -diff : Number(a.dataset.order)-Number(b.dataset.order);
    }).forEach(c => grid?.append(c));
    let shown = 0;
    items.forEach(c => {
      const tags = (c.dataset.tags || '').split(' ');
      const matchesDiet = !diet || tags.includes(diet) || (diet === 'vegetarian' && tags.includes('vegan'));
      c.hidden = !((category === 'all' || category === cat.dataset.cat) && (!q || index.get(c).includes(q)) &&
        matchesDiet && (!onlyFavorites || favorites.has(c.dataset.id)));
      if (!c.hidden) shown++;
    });
    cat.hidden = shown === 0; visible += shown;
  });
  if (noResults) noResults.hidden = visible > 0;
  if (count) count.textContent = new Intl.NumberFormat(document.documentElement.lang).format(visible) + ' ' + count.dataset.label;
};
nav?.addEventListener('click',e => {
  const b = e.target.closest('[data-filter]'); if (!b) return;
  category = b.dataset.filter; syncNav(); apply();
  if (category !== 'all') document.getElementById(category)?.scrollIntoView({behavior:reduced.matches?'auto':'smooth',block:'start'});
});
search?.addEventListener('input',() => { category='all'; syncNav(); apply(); });
$$('[data-diet]').forEach(b => b.addEventListener('click',() => {
  diet = diet === b.dataset.diet ? '' : b.dataset.diet;
  $$('[data-diet]').forEach(b => b.setAttribute('aria-pressed',String(b.dataset.diet===diet))); apply();
}));
const paintFavorites = () => $$('[data-favorite]').forEach(b => {
  const selected=favorites.has(b.dataset.favorite); b.hidden=false;
  b.setAttribute('aria-pressed',String(selected)); b.textContent=selected?'♥':'♡';
});
paintFavorites();
$$('[data-favorite]').forEach(b => b.addEventListener('click',() => {
  const id=b.dataset.favorite; favorites.has(id)?favorites.delete(id):favorites.add(id);
  try { localStorage.setItem(key,JSON.stringify([...favorites])); } catch {}
  paintFavorites(); apply();
}));
$('[data-favorites-only]')?.addEventListener('click',e => {
  onlyFavorites=!onlyFavorites; e.currentTarget.setAttribute('aria-pressed',String(onlyFavorites)); apply();
});
sort?.addEventListener('change',apply);
$('[data-reset]')?.addEventListener('click',() => {
  category='all'; diet=''; onlyFavorites=false; if(search)search.value=''; if(sort)sort.value='default';
  $$('[data-diet], [data-favorites-only]').forEach(b => b.setAttribute('aria-pressed','false'));
  syncNav(); apply(); search?.focus();
});
// Eager decoding avoids blank photos when printing a long, lazy-loaded menu.
const preparePrintImages = async () => {
  await document.fonts.ready;
  await Promise.allSettled($$('img').map(img => { img.loading='eager'; return img.decode(); }));
};
window.addEventListener('beforeprint',() => { $$('img').forEach(img=>img.loading='eager'); });
$('[data-print]')?.addEventListener('click',async () => { await preparePrintImages(); window.print(); });
$('[data-share]')?.addEventListener('click',async e => {
  const b=e.currentTarget,url=new URL(location.href), table=url.searchParams.get('table'), lang=url.searchParams.get('lang'); url.search=''; url.hash='';
  if(table && /^[0-9۰-۹٠-٩]{1,4}$/.test(table))url.searchParams.set('table',table);
  if(lang==='fa'||lang==='en')url.searchParams.set('lang',lang);
  try {
    if(navigator.share){ await navigator.share({title:document.title,url:url.href}); return; }
    await navigator.clipboard.writeText(url.href); if(count)count.textContent=b.dataset.copied;
  } catch(err) { if(err.name!=='AbortError'&&count)count.textContent=b.dataset.failed; }
});
if(nav && 'IntersectionObserver' in window) {
  const spy=new IntersectionObserver(entries => {
    if(category!=='all'||search?.value.trim()||diet||onlyFavorites)return;
    const top=entries.filter(e=>e.isIntersecting).sort((a,b)=>a.boundingClientRect.top-b.boundingClientRect.top)[0];
    if(!top)return;
    nav.querySelectorAll('[data-filter]').forEach(b=>{
      if(b.dataset.filter===top.target.dataset.cat)b.setAttribute('aria-current','location');else b.removeAttribute('aria-current');
    });
  },{rootMargin:'-20% 0px -70% 0px'});
  cats.forEach(c=>spy.observe(c));
}
const customerTools=$('[data-customer-tools]'); if(customerTools)customerTools.hidden=false;
apply();
})();
