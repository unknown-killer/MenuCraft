(() => {
  'use strict';
  const $ = id => document.getElementById(id);
  const storageKey = 'menucraft:studio:v1';
  const dict = {
    fa: {
      skip:'رفتن به استودیو',brand:'منوساز',source:'سورس پروژه ↗',eyebrow:'✨ از سلیقهٔ شما، تا منوی رستوران',hero:'منوی شما،\nبا امضای شما.',intro:'یک قالب انتخاب کنید، رنگ دلخواهتان را بسازید و دسته‌ها را جابه‌جا کنید. نتیجه را همین‌جا، زنده ببینید.',start:'🎨 شروع طراحی',download:'📦 دریافت نسخهٔ کامل',noteTitle:'هر رستوران، یک شخصیت.',note:'از کافهٔ مینیمال تا رستوران سنتی؛ قالب‌ها نقطهٔ شروع‌اند، رنگ‌ها و چیدمان را شما تعیین می‌کنید.',themesStat:'قالب قابل‌تغییر',packsStat:'منوی آماده در نسخهٔ کامل',languagesStat:'دو زبان، خوانایی درست',responsiveStat:'از موبایل تا دسکتاپ',studioEyebrow:'🛠️ امتحان کنید',studioTitle:'استودیوی زندهٔ منو',studioHint:'بدون ثبت‌نام؛ تغییرات فقط در همین مرورگر ذخیره می‌شوند.',appearance:'🎨 ظاهر دلخواه',themeLabel:'قالب شروع',loading:'در حال بارگذاری قالب‌ها…',primaryColor:'رنگ اصلی',background:'پس‌زمینه',surface:'رنگ کارت',textColor:'رنگ نوشته',layoutLabel:'چیدمان غذاها',grid:'شبکه‌ای',list:'لیستی',magazine:'مجله‌ای',compact:'فشرده',columnsLabel:'ستون‌ها در نمایش عریض',radiusLabel:'گردی کارت‌ها',reorderTitle:'🖐️ منو را خودتان بچینید',reorderHint:'دسته‌ها را بکشید و رها کنید، یا از دکمه‌های بالا و پایین استفاده کنید.',reset:'↺ بازگشت به قالب اولیه',live:'پیش‌نمایش زنده',wide:'🖥️ عریض',phone:'📱 موبایل',openMenu:'↗ منوی نمونه',scope:'🧪 این دمو قالب، رنگ، چیدمان و ترتیب دسته‌ها را نشان می‌دهد. پنل حساب‌ها، ویرایش غذا، بارگذاری تصویر و طراح آزاد بوم در نسخهٔ کامل Go هستند؛ این صفحه سرور مدیریت ندارد.',fullEyebrow:'🚀 وقتی نسخهٔ کامل را نصب کنید',fullTitle:'از یک منوی ساده، خیلی بیشتر.',f1title:'طراحی با کشیدن و رهاکردن',f1:'ترتیب بخش‌های صفحه را عوض کنید و در بوم آزاد، متن، تصویر، قیمت و دکمه را هرجا خواستید قرار دهید.',f2title:'رنگ، قلم و قالب دلخواه',f2:'۱۱۴ قالب آماده؛ پالت مستقل روشن و تاریک، پس‌زمینه، اندازه‌ها، کارت‌ها و حتی سبک اختصاصی.',f3title:'منوی واقعی، مدیریت ساده',f3:'غذا، دسته، تخفیف، موجودی، کالری، زمان آماده‌سازی و ۵۴ منوی آماده برای شروع سریع.',f4title:'تجربهٔ راحت مشتری',f4:'دفترچهٔ انتخاب با تعداد و یادداشت، نمایش به گارسون، جست‌وجو، فیلتر غذایی، علاقه‌مندی و چاپ.',f5title:'یک لینک، روی هر دستگاه',f5:'منوی واکنش‌گرا با کد QR، ساعات کاری، اطلاعات تماس و نمایش درست فارسی و انگلیسی.',f6title:'آماده برای سایت شما',f6:'نسخهٔ ویندوز و لینوکس، افزونهٔ وردپرس، ورود و خروج JSON و دریافت CSV.',footer:'با سلیقهٔ خودتان بسازید. روی سرور خودتان منتشر کنید.',docs:'📚 امکانات و راهنمای کامل',saved:'✓ تغییرات در این مرورگر ذخیره شد.',sessionOnly:'✓ تغییرات اعمال شد؛ ذخیرهٔ مرورگر در دسترس نیست.',resetDone:'✓ تنظیمات اصلی این قالب برگشت.',error:'بارگذاری قالب‌ها انجام نشد؛ صفحه را دوباره باز کنید.',up:'بالا بردن',down:'پایین بردن',preview:'پیش‌نمایش منوی رستوران'
    },
    en: {
      skip:'Skip to the studio',brand:'Menu studio',source:'Source code ↗',eyebrow:'✨ Your taste, on the menu',hero:'Your menu.\nYour signature.',intro:'Pick a theme, make the colors yours and rearrange the categories. Watch your restaurant menu take shape, live.',start:'🎨 Start designing',download:'📦 Get the full app',noteTitle:'Every restaurant has a personality.',note:'From a minimal café to a traditional kitchen: themes are a starting point. You choose the colors and the layout.',themesStat:'Customizable themes',packsStat:'Starter menus in the full app',languagesStat:'Two languages, clear reading',responsiveStat:'Mobile to desktop',studioEyebrow:'🛠️ Try it yourself',studioTitle:'Your live menu studio',studioHint:'No account needed. Changes stay in this browser.',appearance:'🎨 Make it yours',themeLabel:'Starting theme',loading:'Loading themes…',primaryColor:'Primary',background:'Background',surface:'Card color',textColor:'Text color',layoutLabel:'Dish layout',grid:'Grid',list:'List',magazine:'Magazine',compact:'Compact',columnsLabel:'Columns on wide screens',radiusLabel:'Card corner radius',reorderTitle:'🖐️ Arrange your menu',reorderHint:'Drag categories into place, or use the up and down buttons.',reset:'↺ Reset to the original theme',live:'Live preview',wide:'🖥️ Wide',phone:'📱 Mobile',openMenu:'↗ Sample menu',scope:'🧪 This demo previews themes, colors, layouts and category order. Accounts, dish editing, image uploads and the freeform canvas designer belong to the full Go app. This static page has no admin server.',fullEyebrow:'🚀 In the full app',fullTitle:'Much more than a simple menu.',f1title:'Design with drag and drop',f1:'Reorder page sections. Place text, images, prices and buttons anywhere on the freeform canvas.',f2title:'Your colors, fonts and theme',f2:'114 starting themes, separate light and dark palettes, backgrounds, sizes, card styles and custom CSS.',f3title:'A real menu, simply managed',f3:'Dishes, categories, discounts, availability, calories, preparation time and 54 starter menu packs.',f4title:'Easy for your guests',f4:'A guest notebook with quantities and notes, waiter presentation, search, dietary filters, favorites and printing.',f5title:'One link. Every screen.',f5:'A responsive menu with QR codes, business hours, contact details and proper Persian and English layouts.',f6title:'Ready for your website',f6:'Windows and Linux builds, a WordPress plugin, JSON import/export and CSV downloads.',footer:'Create it your way. Publish it on your own server.',docs:'📚 Full features & documentation',saved:'✓ Saved in this browser.',sessionOnly:'✓ Applied. Browser storage is unavailable.',resetDone:'✓ Original theme restored.',error:'Could not load the themes. Please reload the page.',up:'Move up',down:'Move down',preview:'Restaurant menu preview'
    }
  };
  let catalogue = [], state = {}, frameReady = false, dragged = null;
  try { state = JSON.parse(localStorage.getItem(storageKey) || '{}') || {}; } catch { state = {}; }
  if (typeof state !== 'object' || Array.isArray(state)) state = {};
  const queryLanguage = new URL(location.href).searchParams.get('lang');
  state.lang = ['fa','en'].includes(queryLanguage) ? queryLanguage : state.lang === 'en' ? 'en' : 'fa';
  state.device = state.device === 'phone' ? 'phone' : 'wide';
  const t = key => dict[state.lang][key] || key;
  const validColor = color => typeof color === 'string' && /^#[\da-f]{6}$/i.test(color);
  function save(message = 'saved') {
    try { localStorage.setItem(storageKey, JSON.stringify(state)); $('status').textContent = t(message); }
    catch { $('status').textContent = t('sessionOnly'); }
  }
  function translate() {
    document.documentElement.lang = state.lang;
    document.documentElement.dir = state.lang === 'fa' ? 'rtl' : 'ltr';
    document.title = state.lang === 'fa' ? 'MenuCraft · استودیوی منوی شما' : 'MenuCraft · Your menu studio';
    document.querySelectorAll('[data-i18n]').forEach(el => { el.textContent = t(el.dataset.i18n); });
    $('language').textContent = state.lang === 'fa' ? 'English' : 'فارسی';
    $('language').lang = state.lang === 'fa' ? 'en' : 'fa';
    $('preview').title = t('preview');
    const url = new URL(location.href); url.searchParams.set('lang', state.lang); history.replaceState(null, '', url);
    $('status').textContent = '';
  }
  function theme() { return catalogue.find(item => item.key === state.theme); }
  function fillThemes() {
    $('theme').replaceChildren(...catalogue.map(item => {
      const option = document.createElement('option'); option.value = item.key; option.textContent = item[state.lang]; return option;
    }));
    $('theme').value = state.theme; $('theme').disabled = false;
  }
  function fillControls() {
    const preset = theme();
    ['primary','background','surface','text'].forEach((key,index) => { $(key).value = state.colors[key] || preset.colors[index]; });
    $('layout').value = state.layout;
    ['columns','radius'].forEach(key => { $(key).value = state[key]; $(key + '-value').value = state[key]; });
    document.querySelectorAll('[data-device]').forEach(button => button.setAttribute('aria-pressed', String(button.dataset.device === state.device)));
    $('preview-stage').classList.toggle('phone', state.device === 'phone');
  }
  function loadMenu() {
    frameReady = false;
    const path = `menus/${state.theme}/${state.lang}.html`;
    $('preview').src = path; $('open-menu').href = path;
  }
  function resetState() {
    const preset = theme();
    state.colors = {}; state.layout = preset.layout; state.columns = preset.columns; state.radius = preset.radius; state.order = [];
  }
  function applySettings() {
    if (!frameReady) return;
    const doc = $('preview').contentDocument, root = doc.documentElement;
    root.dataset.mode = theme().mode === 'dark' ? 'dark' : 'light';
    const variables = {primary:['--c-primary','--c-price'],background:['--c-bg','--bg-base'],surface:['--c-surface'],text:['--c-text']};
    Object.entries(state.colors).forEach(([key,value]) => { (variables[key] || []).forEach(name => root.style.setProperty(name,value)); });
    root.style.setProperty('--cols', String(state.columns));
    root.style.setProperty('--radius', state.radius + 'px');
    root.style.setProperty('--radius-sm', Math.round(state.radius * .6) + 'px');
    root.style.setProperty('--radius-lg', Math.round(state.radius * 1.4) + 'px');
    ['grid','list','magazine','compact'].forEach(layout => doc.body.classList.toggle('layout-' + layout, layout === state.layout));
    const container = doc.querySelector('.m-main'), nav = doc.querySelector('[data-nav]');
    state.order.forEach(id => {
      const section = [...doc.querySelectorAll('[data-cat]')].find(el => el.dataset.cat === id);
      if (container && section) container.append(section);
      const button = nav && [...nav.querySelectorAll('[data-filter]')].find(el => el.dataset.filter === id);
      if (button) nav.append(button);
    });
  }
  function categoryRows() {
    if (!frameReady) return;
    const sections = [...$('preview').contentDocument.querySelectorAll('[data-cat]')];
    const ids = sections.map(el => el.dataset.cat);
    state.order = [...new Set([...state.order.filter(id => ids.includes(id)), ...ids])];
    $('category-order').replaceChildren(...state.order.map((id,index) => {
      const li = document.createElement('li'); li.draggable = true; li.dataset.id = id;
      const grip = document.createElement('span'); grip.textContent = '⠿'; grip.className = 'grip'; grip.setAttribute('aria-hidden','true');
      const name = document.createElement('span'); name.className = 'category-name'; name.textContent = sections.find(el => el.dataset.cat === id).querySelector('.m-cat-title')?.textContent.trim() || id;
      li.append(grip,name);
      [-1,1].forEach(direction => {
        const button = document.createElement('button'); button.type = 'button'; button.textContent = direction < 0 ? '↑' : '↓';
        button.setAttribute('aria-label', `${t(direction < 0 ? 'up' : 'down')} ${name.textContent}`);
        button.disabled = direction < 0 ? index === 0 : index === state.order.length - 1;
        button.addEventListener('click', () => {
          const target = index + direction; [state.order[index],state.order[target]] = [state.order[target],state.order[index]];
          categoryRows(); applySettings(); save();
          const row = [...$('category-order').children].find(el => el.dataset.id === id);
          const buttons = row.querySelectorAll('button'), preferred = buttons[direction < 0 ? 0 : 1];
          (preferred.disabled ? buttons[direction < 0 ? 1 : 0] : preferred).focus();
        }); li.append(button);
      });
      li.addEventListener('dragstart', event => { dragged = id; event.dataTransfer.effectAllowed = 'move'; event.dataTransfer.setData('text/plain',id); });
      li.addEventListener('dragover', event => { if (dragged && dragged !== id) { event.preventDefault(); li.classList.add('drag-over'); } });
      li.addEventListener('dragleave', () => li.classList.remove('drag-over'));
      li.addEventListener('dragend', () => { dragged = null; document.querySelectorAll('.drag-over').forEach(el => el.classList.remove('drag-over')); });
      li.addEventListener('drop', event => {
        event.preventDefault(); li.classList.remove('drag-over');
        if (!dragged || dragged === id || !state.order.includes(dragged)) return;
        const moved = dragged; state.order = state.order.filter(value => value !== moved);
        const after = event.clientY > li.getBoundingClientRect().top + li.offsetHeight / 2;
        state.order.splice(state.order.indexOf(id) + (after ? 1 : 0), 0, moved); dragged = null;
        categoryRows(); applySettings(); save();
      }); return li;
    }));
  }
  $('preview').addEventListener('load', () => {
    if (!catalogue.length) return;
    frameReady = true; categoryRows(); applySettings();
  });
  $('language').addEventListener('click', () => { state.lang = state.lang === 'fa' ? 'en' : 'fa'; translate(); if (catalogue.length) { fillThemes(); loadMenu(); } save(); });
  $('theme').addEventListener('change', () => { state.theme = $('theme').value; resetState(); fillControls(); loadMenu(); save(); });
  ['primary','background','surface','text'].forEach(key => $(key).addEventListener('input', () => { if (!catalogue.length) return; state.colors[key] = $(key).value; applySettings(); save(); }));
  ['layout','columns','radius'].forEach(key => $(key).addEventListener('input', () => { if (!catalogue.length) return; state[key] = key === 'layout' ? $(key).value : Number($(key).value); fillControls(); applySettings(); save(); }));
  document.querySelectorAll('[data-device]').forEach(button => button.addEventListener('click', () => { if (!catalogue.length) return; state.device = button.dataset.device; fillControls(); save(); }));
  $('reset').addEventListener('click', () => { if (!catalogue.length) return; resetState(); fillControls(); loadMenu(); save('resetDone'); });
  translate();
  fetch('themes.json').then(response => { if (!response.ok) throw new Error('catalogue'); return response.json(); }).then(items => {
    catalogue = items;
    if (!catalogue.some(item => item.key === state.theme)) { state.theme = catalogue.some(item => item.key === 'saffron') ? 'saffron' : catalogue[0].key; resetState(); }
    const preset = theme();
    state.colors = Object.fromEntries(Object.entries(state.colors || {}).filter(([key,value]) => ['primary','background','surface','text'].includes(key) && validColor(value)));
    if (!['grid','list','magazine','compact'].includes(state.layout)) state.layout = preset.layout;
    ['columns','radius'].forEach(key => { if (!Number.isInteger(state[key]) || state[key] < (key === 'columns' ? 1 : 0) || state[key] > (key === 'columns' ? 4 : 40)) state[key] = preset[key]; });
    state.order = Array.isArray(state.order) ? state.order.filter(value => typeof value === 'string') : [];
    fillThemes(); fillControls(); loadMenu();
  }).catch(() => { $('status').textContent = t('error'); });
})();
