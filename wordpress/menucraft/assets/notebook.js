/* Private guest notebook. No checkout, order submission or remote storage. */
(() => {
  'use strict';
  const digits = value => String(value ?? '').replace(/[۰-۹٠-٩]/g, c => String(c.charCodeAt(0) - (c >= '۰' ? 1776 : 1632)));
  const tableNumber = value => { const s = digits(value).trim(); return /^\d{1,4}$/.test(s) && Number(s) > 0 ? String(Number(s)) : ''; };
  function restore(raw, available) {
    const result = {items: [], note: '', table: ''};
    if (!raw || typeof raw !== 'object' || Array.isArray(raw)) return result;
    result.note = typeof raw.note === 'string' ? raw.note.slice(0,500) : '';
    result.table = tableNumber(raw.table);
    const seen = new Set();
    if (Array.isArray(raw.items)) raw.items.slice(0,200).forEach(item => {
      if (!item || typeof item.id !== 'string' || !available.has(item.id) || seen.has(item.id) || !Number.isInteger(item.qty) || item.qty < 1 || item.qty > 99) return;
      seen.add(item.id); result.items.push({id:item.id,qty:item.qty,note:typeof item.note === 'string' ? item.note.slice(0,200) : ''});
    });
    return result;
  }
  function total(items, catalog) {
    let amount = 0n;
    for (const item of items) {
      const price = catalog.get(item.id)?.price;
      if (typeof price !== 'string' || !/^\d{1,13}$/.test(price)) return null;
      amount += BigInt(price) * BigInt(item.qty);
    }
    return amount;
  }
  if (typeof module !== 'undefined' && module.exports) module.exports = {restore,tableNumber,total};
  if (typeof document === 'undefined') return;
  document.querySelectorAll('[data-notebook]').forEach((root, instance) => {
    const fa = (root.lang || document.documentElement.lang).toLowerCase().split('-')[0] === 'fa';
    const words = fa ? {
      title:'📝 دفترچهٔ انتخاب من',add:'＋ افزودن به دفترچه',open:'📝 انتخاب‌های من',close:'بستن',empty:'هنوز چیزی انتخاب نکرده‌اید. از کنار غذاها «افزودن به دفترچه» را بزنید.',
      hint:'این فهرست سفارش ثبت نمی‌کند؛ آن را به گارسون نشان دهید.',table:'شمارهٔ میز',tableHint:'اختیاری؛ ۱ تا ۹۹۹۹',note:'یادداشت کلی',itemNote:'توضیح این مورد',placeholder:'مثلاً بدون پیاز یا سس جدا',
      more:'افزایش تعداد',less:'کاهش تعداد',remove:'حذف',show:'👀 نمایش به گارسون',edit:'✏️ ادامهٔ انتخاب',copy:'📋 کپی فهرست',print:'🖨️ چاپ فهرست',clear:'پاک‌کردن دفترچه',confirm:'همهٔ انتخاب‌ها و یادداشت‌های این دفترچه پاک شوند؟',
      pending:'قیمت به‌زودی',sumPending:'قیمت بعضی انتخاب‌ها هنوز مشخص نیست؛ جمع کل پس از تعیین قیمت محاسبه می‌شود.',sum:'جمع تقریبی',sumHint:'مبلغ نهایی را با رستوران بررسی کنید.',saved:'در همین مرورگر ذخیره شد.',unsaved:'ذخیرهٔ مرورگر در دسترس نیست؛ پیش از بستن صفحه فهرست را کپی کنید.',copied:'فهرست کپی شد.',manual:'کپی خودکار در دسترس نیست؛ متن زیر را انتخاب و کپی کنید.',
      updated:'مواردی که دیگر در این منو در دسترس نیستند از دفترچه حذف شدند.',limit:'حداکثر ۲۰۰ مورد و ۹۹ عدد از هر مورد.',invalidTable:'شمارهٔ میز باید از ۱ تا ۹۹۹۹ باشد.',added:'به دفترچه اضافه شد.'
    } : {
      title:'📝 My menu notebook',add:'＋ Add to notebook',open:'📝 My selections',close:'Close',empty:'Nothing selected yet. Use “Add to notebook” beside a dish.',
      hint:'This list does not place an order. Show it to your waiter.',table:'Table number',tableHint:'Optional; 1–9999',note:'General note',itemNote:'Item note',placeholder:'For example, no onions or sauce on the side',
      more:'Increase quantity',less:'Decrease quantity',remove:'Remove',show:'👀 Show to waiter',edit:'✏️ Continue choosing',copy:'📋 Copy list',print:'🖨️ Print list',clear:'Clear notebook',confirm:'Clear all selections and notes in this notebook?',
      pending:'Price coming soon',sumPending:'Some selections have no price yet. The total will be available once all prices are set.',sum:'Estimated total',sumHint:'Confirm the final amount with the restaurant.',saved:'Saved in this browser.',unsaved:'Browser storage is unavailable. Copy the list before closing this page.',copied:'List copied.',manual:'Automatic copying is unavailable. Select and copy the text below.',
      updated:'Items no longer available on this menu were removed from the notebook.',limit:'Up to 200 different items and 99 of each.',invalidTable:'The table number must be between 1 and 9999.',added:'Added to notebook.'
    };
    const number = value => new Intl.NumberFormat(fa?'fa':'en').format(value);
    const catalog = new Map();
    root.querySelectorAll('[data-note-item]').forEach(card => {
      catalog.set(card.dataset.noteItem,{name:card.dataset.noteName || '',price:card.dataset.notePrice || '',pending:card.dataset.notePending==='1',card});
    });
    const key = 'menucraft:notebook:v1:' + root.dataset.notebook;
    let raw, storageOK = true;
    try { raw = JSON.parse(localStorage.getItem(key) || '{}'); } catch { storageOK = false; }
    let state = restore(raw,catalog), presentation = false;
    const scanned = tableNumber(new URL(location.href).searchParams.get('table'));
    if (scanned) state.table = scanned;
    const currency = root.dataset.noteCurrency || '';
    const money = amount => number(amount) + ' ' + currency;
    const el = (tag,className,text) => { const e = document.createElement(tag); if(className)e.className=className; if(text!==undefined)e.textContent=text; return e; };
    const button = (text,action,className='mc-n-button') => { const b=el('button',className,text); b.type='button'; b.addEventListener('click',action); return b; };
    const dialog = el('dialog','mc-notebook'); dialog.lang=fa?'fa':'en'; dialog.dir=fa?'rtl':'ltr';
    const title=el('h2','',words.title); title.id='mc-notebook-title-'+instance; title.tabIndex=-1; dialog.setAttribute('aria-labelledby',title.id);
    const top=el('header','mc-n-top'); top.append(title,button('×',()=>dialog.close(),'mc-n-close')); top.lastChild.setAttribute('aria-label',words.close);
    const hint=el('p','mc-n-hint',words.hint), venue=el('p','mc-n-venue',root.dataset.noteVenue || document.title);
    const tableLabel=el('label','mc-n-field',words.table), table=el('input'); table.type='text'; table.inputMode='numeric'; table.maxLength=4; table.placeholder=words.tableHint; table.value=state.table; tableLabel.append(table);
    const rows=el('div','mc-n-rows');
    const generalLabel=el('label','mc-n-field',words.note), general=el('textarea'), generalText=el('p','mc-n-readnote',state.note); general.maxLength=500; general.rows=2; general.value=state.note; generalLabel.append(general,generalText);
    general.setAttribute('aria-label',words.note);
    const estimate=el('p','mc-n-total'), estimateHint=el('small','',words.sumHint);
    const status=el('p','mc-n-status'); status.setAttribute('role','status'); status.setAttribute('aria-live','polite');
    const announcement=el('span','mc-n-sr'); announcement.setAttribute('role','status'); root.append(announcement);
    const copyFallback=el('textarea','mc-n-copy'); copyFallback.readOnly=true; copyFallback.hidden=true; copyFallback.setAttribute('aria-label',words.copy);
    const actions=el('div','mc-n-actions');
    const toggle=button(words.show,()=> { if(!table.reportValidity())return; presentation=!presentation; paint(); dialog.scrollTop=0; title.focus({preventScroll:true}); });
    const launcher=button('',()=> { paint(); dialog.showModal(); },'mc-n-launcher'); launcher.setAttribute('aria-haspopup','dialog');
    const clear=button(words.clear,()=> { if(!confirm(words.confirm))return; state={items:[],note:'',table:state.table}; general.value=''; presentation=false; save(); paint(); });
    const copy=button(words.copy,async()=> { if(!table.reportValidity())return; const value=asText(); try { await navigator.clipboard.writeText(value); status.textContent=words.copied; copyFallback.hidden=true; } catch { copyFallback.hidden=false; copyFallback.value=value; copyFallback.focus(); copyFallback.select(); status.textContent=words.manual; } });
    const print=button(words.print,()=> { if(!table.reportValidity())return; document.body.classList.add('mc-print-notebook'); dialog.classList.add('mc-n-print-target'); window.print(); });
    window.addEventListener('afterprint',()=> { document.body.classList.remove('mc-print-notebook'); dialog.classList.remove('mc-n-print-target'); });
    dialog.addEventListener('close',()=> { document.body.classList.remove('mc-print-notebook'); dialog.classList.remove('mc-n-print-target'); launcher.focus(); });
    actions.append(toggle,copy,print,clear);
    dialog.append(top,venue,hint,tableLabel,rows,generalLabel,estimate,estimateHint,actions,copyFallback,status);
    document.body.append(dialog); root.append(launcher);
    function save() { try { localStorage.setItem(key,JSON.stringify(state)); storageOK=true; } catch { storageOK=false; } status.textContent=storageOK?words.saved:words.unsaved; copyFallback.hidden=true; }
    table.addEventListener('input',()=> { const value=tableNumber(table.value); table.setCustomValidity(table.value.trim()&&!value?words.invalidTable:''); if(table.value.trim()&&!value){status.textContent=words.invalidTable; return;} state.table=value; save(); });
    general.addEventListener('input',()=> { state.note=general.value.slice(0,500); generalText.textContent=state.note; save(); });
    function asText() {
      const lines=[root.dataset.noteVenue || document.title,words.title,words.hint];
      if(state.table)lines.push(words.table+': '+number(Number(state.table)));
      state.items.forEach(item=> { const dish=catalog.get(item.id); lines.push(number(item.qty)+' × '+dish.name+(dish.price?' — '+money(BigInt(dish.price)*BigInt(item.qty)):dish.pending?' — '+words.pending:'')); if(item.note)lines.push('  '+item.note); });
      if(state.note)lines.push(words.note+': '+state.note);
      const amount=total(state.items,catalog); if(state.items.length && amount!==null)lines.push(words.sum+': '+money(amount),words.sumHint);
      if(state.items.some(item=>catalog.get(item.id).pending))lines.push(words.sumPending);
      return lines.join('\n');
    }
    function paint() {
      const quantity=state.items.reduce((n,item)=>n+item.qty,0);
      launcher.textContent=words.open+' · '+number(quantity)+(state.table?' · '+words.table+' '+number(Number(state.table)):'');
      dialog.classList.toggle('mc-n-presentation',presentation); toggle.textContent=presentation?words.edit:words.show;
      table.readOnly=presentation; general.readOnly=presentation;
      generalText.textContent=state.note;
      clear.hidden=presentation; toggle.disabled=copy.disabled=print.disabled=state.items.length===0 && !state.note;
      const amount=total(state.items,catalog); estimate.hidden=estimateHint.hidden=!state.items.length || amount===null;
      if(state.items.some(item=>catalog.get(item.id).pending)){estimateHint.hidden=false;estimateHint.textContent=words.sumPending;}else estimateHint.textContent=words.sumHint;
      estimate.textContent=amount===null?'':words.sum+': '+money(amount);
      rows.replaceChildren();
      if(!state.items.length)rows.append(el('p','mc-n-empty',words.empty));
      state.items.forEach(item=> {
        const dish=catalog.get(item.id), row=el('section','mc-n-row'), heading=el('h3','',number(item.qty)+' × '+dish.name);
        const tools=el('div','mc-n-quantity');
        const change=delta=> { item.qty+=delta; save(); paint(); const current=rows.querySelector(`[data-quantity-id="${CSS.escape(item.id)}"]`); (current?.querySelector('button:not(:disabled)') || toggle).focus(); };
        const minus=button('−',()=>change(-1)), plus=button('+',()=>change(1));
        minus.disabled=item.qty<=1; plus.disabled=item.qty>=99;
        minus.setAttribute('aria-label',words.less+': '+dish.name); plus.setAttribute('aria-label',words.more+': '+dish.name);
        tools.dataset.quantityId=item.id;
        tools.append(minus,el('span','',number(item.qty)),plus,button(words.remove,()=> {state.items=state.items.filter(value=>value.id!==item.id);save();paint();toggle.focus();}));
        const label=el('label','mc-n-field',words.itemNote), note=el('textarea'); note.maxLength=200; note.rows=2; note.value=item.note; note.placeholder=words.placeholder; note.readOnly=presentation;
        note.setAttribute('aria-label',words.itemNote+': '+dish.name);
        const readNote=el('p','mc-n-readnote',item.note);
        note.addEventListener('input',()=> {item.note=note.value.slice(0,200);readNote.textContent=item.note;save();}); label.append(note,readNote);
        row.append(heading); if(dish.price)row.append(el('p','mc-n-price',money(BigInt(dish.price)*BigInt(item.qty))));
        if(dish.pending)row.append(el('p','mc-n-price',words.pending));
        row.append(tools,label); rows.append(row);
      });
    }
    catalog.forEach((dish,id)=> {
      // No innerHTML: dish names and guest notes are untrusted text.
      const add=button(words.add,()=> {
        const current=state.items.find(item=>item.id===id);
        if(current ? current.qty>=99 : state.items.length>=200){announcement.textContent=words.limit;return;}
        if(current)current.qty++; else state.items.push({id,qty:1,note:''});
        save(); paint(); announcement.textContent=dish.name+' — '+words.added+(storageOK?'':' '+words.unsaved);
      },'mc-n-add');
      add.setAttribute('aria-label',words.add+': '+dish.name);
      (dish.card.querySelector('[data-note-actions]') || dish.card).append(add);
    });
    paint();
    if(raw?.items?.length>state.items.length)status.textContent=words.updated;
    else if(!storageOK)status.textContent=words.unsaved;
    if(scanned)save();
  });
})();
