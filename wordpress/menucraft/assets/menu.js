(() => {'use strict';document.querySelectorAll('.menucraft').forEach(root=>{
const search=root.querySelector('[data-mc-search]'),category=root.querySelector('[data-mc-category]');if(!search)return;
const fold=s=>(s||'').toLowerCase().normalize('NFKC').replace(/[يى]/g,'ی').replace(/ك/g,'ک').replace(/[\u064b-\u065f\u0670]/g,'').replace(/[\u200c\s]+/g,' ').trim();
const update=()=>{let count=0;root.querySelectorAll('[data-mc-cat]').forEach(cat=>{let shown=0;cat.querySelectorAll('[data-mc-text]').forEach(item=>{item.hidden=!((category.value==='all'||category.value===cat.dataset.mcCat)&&fold(item.dataset.mcText).includes(fold(search.value)));if(!item.hidden)shown++;});cat.hidden=shown===0;count+=shown;});root.querySelector('[data-mc-empty]').hidden=count>0;const status=root.querySelector('[data-mc-status]');status.textContent=new Intl.NumberFormat(root.lang).format(count)+' '+status.dataset.label;};
search.addEventListener('input',update);category.addEventListener('change',update);root.querySelector('.menucraft-controls').hidden=false;update();
});})();
