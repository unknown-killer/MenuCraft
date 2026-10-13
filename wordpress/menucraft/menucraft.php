<?php
/**
 * Plugin Name: MenuCraft Restaurant Menu
 * Description: Persian/English restaurant menus: import a MenuCraft JSON snapshot or embed a live Go menu.
 * Version: 1.1.0
 * Requires at least: 6.4
 * Requires PHP: 8.1
 * Author: unknown-killer
 * Plugin URI: https://github.com/unknown-killer/MenuCraft
 * Text Domain: menucraft
 */
if (!defined('ABSPATH')) { exit; }

function menucraft_text($fa, $en) {
    return str_starts_with(determine_locale(), 'fa') ? $fa : $en;
}
function menucraft_settings() {
    $value = get_option('menucraft_settings', array());
    return wp_parse_args(is_array($value) ? $value : array(), array('mode'=>'native', 'url'=>'', 'height'=>900, 'title'=>'Restaurant menu / منوی رستوران'));
}
add_action('admin_menu', function () {
    add_options_page('MenuCraft', 'MenuCraft', 'manage_options', 'menucraft', 'menucraft_admin');
});
function menucraft_admin() {
    if (!current_user_can('manage_options')) { return; }
    $settings = menucraft_settings();
    $menu = get_option('menucraft_menu_data', array());
    ?>
    <div class="wrap">
      <h1>MenuCraft · منوساز</h1>
      <?php if (isset($_GET['saved'])) { ?><div class="notice notice-success"><p><?php echo esc_html(menucraft_text('ذخیره شد.', 'Saved.')); ?></p></div><?php } ?>
      <p><?php echo esc_html(menucraft_text('این شورت‌کد را در یک برگه قرار دهید:', 'Add this shortcode to a page:')); ?> <code>[menucraft_menu lang="fa"]</code> / <code>[menucraft_menu lang="en"]</code></p>
      <form method="post" action="<?php echo esc_url(admin_url('admin-post.php')); ?>">
        <input type="hidden" name="action" value="menucraft_save">
        <?php wp_nonce_field('menucraft_save'); ?>
        <table class="form-table"><tbody>
          <tr><th><label for="mc-mode"><?php echo esc_html(menucraft_text('حالت نمایش', 'Display mode')); ?></label></th><td>
            <select id="mc-mode" name="mode">
              <option value="native" <?php selected($settings['mode'],'native'); ?>><?php echo esc_html(menucraft_text('مستقل در وردپرس (JSON)', 'Native WordPress (JSON)')); ?></option>
              <option value="embed" <?php selected($settings['mode'],'embed'); ?>><?php echo esc_html(menucraft_text('نمایش زنده از سرور Go', 'Live Go embed')); ?></option>
            </select>
          </td></tr>
          <tr><th><label for="mc-url"><?php echo esc_html(menucraft_text('نشانی منوی Go', 'Go menu URL')); ?></label></th><td>
            <input id="mc-url" class="regular-text" type="url" name="url" dir="ltr" placeholder="https://menu.example.com/m/my-restaurant" value="<?php echo esc_attr($settings['url']); ?>">
            <p class="description"><?php echo esc_html(menucraft_text('در حالت زنده، دامنه همین سایت را در MENU_EMBED_ORIGINS سرور Go وارد کنید؛ هر دو سایت باید HTTPS داشته باشند.', 'For live mode, add this WordPress origin to MENU_EMBED_ORIGINS on the Go server. Both sites should use HTTPS.')); ?></p>
          </td></tr>
          <tr><th><label for="mc-title"><?php echo esc_html(menucraft_text('عنوان دسترس‌پذیر', 'Accessible title')); ?></label></th><td><input id="mc-title" name="title" class="regular-text" value="<?php echo esc_attr($settings['title']); ?>"></td></tr>
          <tr><th><label for="mc-height"><?php echo esc_html(menucraft_text('ارتفاع نمایش زنده (پیکسل)', 'Live frame height (pixels)')); ?></label></th><td><input id="mc-height" type="number" name="height" min="400" max="2400" value="<?php echo esc_attr($settings['height']); ?>"></td></tr>
        </tbody></table>
        <?php submit_button(); ?>
      </form>
      <hr>
      <h2><?php echo esc_html(menucraft_text('وارد کردن منو', 'Import menu')); ?></h2>
      <p><?php echo esc_html(menucraft_text('در پنل Go، از بخش انتقال منو فایل JSON بگیرید. ورود فایل، نسخه قبلی وردپرس را جایگزین می‌کند. طراحی بوم و قالب Go منتقل نمی‌شود. تصاویر را می‌توانید بعداً در کتابخانه وردپرس میزبانی کنید؛ نشانی‌های تصویر فعلی به سرور اصلی وابسته‌اند.', 'Export JSON from the Go panel’s Menu transfer page. Import replaces the previous WordPress snapshot. Go themes and canvas designs are not imported. Existing image URLs depend on the original server; you can host images in WordPress and edit their URLs in the JSON.')); ?></p>
      <p><?php echo esc_html(menucraft_text('نام مجموعه فعلی: ', 'Current venue: ')); echo esc_html($menu['restaurant']['name_fa'] ?? $menu['restaurant']['name_en'] ?? '—'); ?></p>
      <form method="post" enctype="multipart/form-data" action="<?php echo esc_url(admin_url('admin-post.php')); ?>">
        <input type="hidden" name="action" value="menucraft_import">
        <?php wp_nonce_field('menucraft_import'); ?>
        <label for="mc-file"><?php echo esc_html(menucraft_text('فایل JSON، حداکثر ۲ مگابایت', 'JSON file, maximum 2 MB')); ?></label>
        <input id="mc-file" type="file" name="menu_file" accept=".json,application/json" required>
        <?php submit_button(menucraft_text('وارد کردن و جایگزینی', 'Import and replace')); ?>
      </form>
    </div>
    <?php
}
function menucraft_scalar($key, $default = '') {
    return isset($_POST[$key]) && is_string($_POST[$key]) ? wp_unslash($_POST[$key]) : $default;
}
function menucraft_clean_url($raw) {
    if (!is_string($raw) || strlen($raw)>2000) { return ''; }
    $parts = wp_parse_url($raw);
    if (!$parts || !in_array($parts['scheme'] ?? '', array('http','https'), true) ||
        empty($parts['host']) || isset($parts['user']) || isset($parts['pass'])) { return ''; }
    return esc_url_raw($raw, array('http','https'));
}
add_action('admin_post_menucraft_save', function () {
    if (!current_user_can('manage_options')) { wp_die('Forbidden', '', array('response'=>403)); }
    check_admin_referer('menucraft_save');
    $mode = menucraft_scalar('mode') === 'embed' ? 'embed' : 'native';
    $url = menucraft_clean_url(menucraft_scalar('url'));
    if ($mode==='embed' && (!$url || !preg_match('~/m/[a-z0-9]+(?:-[a-z0-9]+)*$~', wp_parse_url($url,PHP_URL_PATH) ?? ''))) {
        wp_die('A complete http(s) MenuCraft /m/slug URL is required. / نشانی معتبر منو لازم است.');
    }
    update_option('menucraft_settings', array('mode'=>$mode, 'url'=>$url,
        'height'=>max(400,min(2400,(int)menucraft_scalar('height','900'))),
        'title'=>sanitize_text_field(substr(menucraft_scalar('title'),0,300))), false);
    wp_safe_redirect(admin_url('options-general.php?page=menucraft&saved=1'));
    exit;
});

/** Validate the complete document before replacing the saved snapshot. */
function menucraft_validate_menu($raw) {
    if (!is_string($raw) || strlen($raw)>2*1024*1024) { return new WP_Error('menu_size','Menu exceeds 2 MB.'); }
    $data = json_decode($raw,true,32);
    if (!is_array($data) || ($data['version'] ?? null)!==1 || !is_array($data['restaurant'] ?? null) ||
        !is_array($data['categories'] ?? null) || !array_is_list($data['categories']) || count($data['categories'])>200) {
        return new WP_Error('menu_format','Invalid MenuCraft v1 document.');
    }
    $currency = $data['restaurant']['currency'] ?? '';
    if (!in_array($currency,array('IRT','IRR','USD','EUR','AED','TRY','IQD'),true)) { return new WP_Error('currency','Invalid currency.'); }
    $clean = array('version'=>1, 'restaurant'=>array('currency'=>$currency), 'categories'=>array());
    foreach (array('name_fa','name_en','description_fa','description_en') as $key) {
        $value=$data['restaurant'][$key] ?? '';
        if (!is_string($value) || strlen($value)>8000) { return new WP_Error('text','Invalid restaurant text.'); }
        $clean['restaurant'][$key]=sanitize_textarea_field($value);
    }
    $source=menucraft_clean_url($data['source_url'] ?? '');
    $source_parts=$source ? wp_parse_url($source) : array();
    $origin=$source ? $source_parts['scheme'].'://'.$source_parts['host'].(isset($source_parts['port'])?':'.$source_parts['port']:'') : '';
    $total=0;
    foreach ($data['categories'] as $cat) {
        if (!is_array($cat) || !is_array($cat['items'] ?? null) || !array_is_list($cat['items']) || !is_bool($cat['visible'] ?? null)) {
            return new WP_Error('category','Invalid category.');
        }
        $out=array('visible'=>$cat['visible'],'items'=>array());
        foreach (array('name_fa','name_en','description_fa','description_en','icon') as $key) {
            $value=$cat[$key] ?? '';
            if (!is_string($value) || strlen($value)>8000) { return new WP_Error('text','Invalid category text.'); }
            $out[$key]=sanitize_textarea_field($value);
        }
        if (!trim($out['name_fa'].$out['name_en'])) { return new WP_Error('name','Category requires a name.'); }
        $total+=count($cat['items']);
        if ($total>5000) { return new WP_Error('count','Too many dishes.'); }
        foreach ($cat['items'] as $item) {
            if (!is_array($item) || !is_bool($item['available'] ?? null)) { return new WP_Error('item','Invalid dish.'); }
            $dish=array('available'=>$item['available']);
            foreach (array('name_fa','name_en','description_fa','description_en') as $key) {
                $value=$item[$key] ?? '';
                if (!is_string($value) || strlen($value)>8000) { return new WP_Error('text','Invalid dish text.'); }
                $dish[$key]=sanitize_textarea_field($value);
            }
            if (!trim($dish['name_fa'].$dish['name_en'])) { return new WP_Error('name','Dish requires a name.'); }
            foreach (array('price','discount_price') as $key) {
                $value=$item[$key] ?? 0;
                if (!is_int($value) || $value<0 || $value>1000000000000) { return new WP_Error('price','Invalid price.'); }
                $dish[$key]=$value;
            }
            if ($dish['discount_price']>0 && $dish['discount_price']>=$dish['price']) { return new WP_Error('price','Invalid discount.'); }
            $tags=$item['tags'] ?? array();
            if (!is_array($tags) || count($tags)>20) { return new WP_Error('tags','Invalid tags.'); }
            $dish['tags']=array_values(array_intersect(array_filter($tags,'is_string'),array('spicy','vegetarian','vegan','new','popular','chef','glutenfree','cold','hot')));
            $path=$item['image_path'] ?? '';
            if (!is_string($path)) { return new WP_Error('image','Invalid image URL.'); }
            if (preg_match('~^/uploads/[0-9]+/[a-zA-Z0-9_-]+\.(?:png|jpe?g|gif|webp)$~D',$path)) { $path=$origin ? $origin.$path : ''; }
            $dish['image_url']=menucraft_clean_url($path);
            $out['items'][]=$dish;
        }
        $clean['categories'][]=$out;
    }
    return $clean;
}
add_action('admin_post_menucraft_import', function () {
    if (!current_user_can('manage_options')) { wp_die('Forbidden','',array('response'=>403)); }
    check_admin_referer('menucraft_import');
    $file=$_FILES['menu_file'] ?? null;
    if (!is_array($file) || ($file['error'] ?? 1)!==UPLOAD_ERR_OK || !is_string($file['tmp_name'] ?? null) ||
        !is_uploaded_file($file['tmp_name']) || filesize($file['tmp_name'])>2*1024*1024) { wp_die('Invalid upload. / فایل نامعتبر است.'); }
    $clean=menucraft_validate_menu(file_get_contents($file['tmp_name']));
    if (is_wp_error($clean)) { wp_die(esc_html($clean->get_error_message())); }
    update_option('menucraft_menu_data',$clean,false);
    wp_safe_redirect(admin_url('options-general.php?page=menucraft&saved=1'));
    exit;
});
add_action('wp_enqueue_scripts', function () {
    wp_enqueue_style('menucraft',plugins_url('assets/menu.css',__FILE__),array(),'1.1.0');
    wp_enqueue_script('menucraft',plugins_url('assets/menu.js',__FILE__),array(),'1.1.0',true);
});
function menucraft_pick($data,$field,$lang) {
    $primary=$data[$field.'_'.$lang] ?? '';
    return $primary!=='' ? $primary : ($data[$field.'_'.($lang==='fa'?'en':'fa')] ?? '');
}
function menucraft_price($amount,$currency,$lang) {
    $number=number_format($amount,0,'.',',');
    $labels=array('IRT'=>'تومان','IRR'=>'ریال','USD'=>'دلار','EUR'=>'یورو','AED'=>'درهم','TRY'=>'لیر','IQD'=>'دینار');
    if ($lang==='fa') { $number=strtr($number,array('0'=>'۰','1'=>'۱','2'=>'۲','3'=>'۳','4'=>'۴','5'=>'۵','6'=>'۶','7'=>'۷','8'=>'۸','9'=>'۹',','=>'٬')); }
    return $number.' '.($lang==='fa' ? ($labels[$currency] ?? $currency) : ($currency==='IRT'?'Toman':$currency));
}
add_shortcode('menucraft_menu', function ($attributes) {
    $attributes=shortcode_atts(array('lang'=>str_starts_with(determine_locale(),'fa')?'fa':'en'),$attributes,'menucraft_menu');
    $lang=$attributes['lang']==='fa'?'fa':'en';
    $settings=menucraft_settings();
    if ($settings['mode']==='embed') {
        $url=menucraft_clean_url($settings['url']);
        if (!$url) { return ''; }
        return '<div class="menucraft"><iframe class="menucraft-frame" src="'.esc_url(add_query_arg('lang',$lang,$url)).
            '" title="'.esc_attr($settings['title']).'" height="'.esc_attr(max(400,min(2400,(int)$settings['height']))).
            '" loading="lazy" referrerpolicy="strict-origin-when-cross-origin"></iframe><p><a href="'.esc_url($url).
            '" target="_blank" rel="noopener noreferrer">'.esc_html($lang==='fa'?'باز کردن منو در صفحه جدا':'Open menu in a new tab').'</a></p></div>';
    }
    $data=get_option('menucraft_menu_data',array());
    if (empty($data['restaurant'])) {
        return current_user_can('manage_options') ? '<p>'.esc_html(menucraft_text('ابتدا فایل JSON را در تنظیمات MenuCraft وارد کنید.','Import a JSON menu in Settings → MenuCraft first.')).'</p>' : '';
    }
    $fa=$lang==='fa'; $currency=$data['restaurant']['currency'];
    $uid=wp_unique_id('menucraft-');
    ob_start(); ?>
    <section class="menucraft" lang="<?php echo esc_attr($lang); ?>" dir="<?php echo $fa?'rtl':'ltr'; ?>">
      <header class="menucraft-header"><span class="menucraft-eyebrow"><?php echo esc_html($fa?'منوی دیجیتال':'THE MENU'); ?></span>
        <h2><?php echo esc_html(menucraft_pick($data['restaurant'],'name',$lang)); ?></h2>
        <p><?php echo esc_html(menucraft_pick($data['restaurant'],'description',$lang)); ?></p>
      </header>
      <div class="menucraft-controls" hidden>
        <label for="<?php echo esc_attr($uid); ?>search"><?php echo esc_html($fa?'جست‌وجوی غذا':'Find a dish'); ?></label>
        <input id="<?php echo esc_attr($uid); ?>search" type="search" data-mc-search placeholder="<?php echo esc_attr($fa?'نام یا توضیح غذا…':'Name or description…'); ?>">
        <label for="<?php echo esc_attr($uid); ?>category"><?php echo esc_html($fa?'دسته‌بندی':'Category'); ?></label>
        <select id="<?php echo esc_attr($uid); ?>category" data-mc-category>
          <option value="all"><?php echo esc_html($fa?'همه دسته‌ها':'All categories'); ?></option>
          <?php foreach ($data['categories'] as $i=>$cat) { if (!$cat['visible']) { continue; } ?>
          <option value="<?php echo esc_attr($i); ?>"><?php echo esc_html(menucraft_pick($cat,'name',$lang)); ?></option>
          <?php } ?>
        </select>
      </div>
      <?php foreach ($data['categories'] as $i=>$cat) {
        if (!$cat['visible']) { continue; }
        $items=array_values(array_filter($cat['items'],fn($it)=>$it['available']));
        if (!$items) { continue; }
      ?>
        <section class="menucraft-category" data-mc-cat="<?php echo esc_attr($i); ?>">
          <h3><?php echo esc_html($cat['icon'].' '.menucraft_pick($cat,'name',$lang)); ?></h3>
          <p><?php echo esc_html(menucraft_pick($cat,'description',$lang)); ?></p>
          <div class="menucraft-grid">
          <?php foreach ($items as $item) { ?>
            <article class="menucraft-dish" data-mc-text="<?php echo esc_attr($item['name_fa'].' '.$item['name_en'].' '.$item['description_fa'].' '.$item['description_en']); ?>">
              <?php if ($item['image_url']) { ?><img src="<?php echo esc_url($item['image_url']); ?>" alt="<?php echo esc_attr(menucraft_pick($item,'name',$lang)); ?>" loading="lazy"><?php } ?>
              <div><h4><?php echo esc_html(menucraft_pick($item,'name',$lang)); ?></h4>
                <p><?php echo esc_html(menucraft_pick($item,'description',$lang)); ?></p>
                <div class="menucraft-price"><?php if ($item['discount_price']>0) { ?><s><?php echo esc_html(menucraft_price($item['price'],$currency,$lang)); ?></s> <?php } ?>
                  <strong><?php echo esc_html(menucraft_price($item['discount_price']>0?$item['discount_price']:$item['price'],$currency,$lang)); ?></strong>
                </div>
              </div>
            </article>
          <?php } ?>
          </div>
        </section>
      <?php } ?>
      <p data-mc-status role="status" aria-live="polite" data-label="<?php echo esc_attr($fa?'غذا پیدا شد':'dishes found'); ?>"></p>
      <p data-mc-empty hidden><?php echo esc_html($fa?'غذایی پیدا نشد.':'No dishes found.'); ?></p>
    </section>
    <?php return ob_get_clean();
});

