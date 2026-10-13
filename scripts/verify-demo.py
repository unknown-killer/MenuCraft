from pathlib import Path
from html.parser import HTMLParser
from urllib.parse import urlsplit, unquote
import json

root=Path('_site').resolve()
class Links(HTMLParser):
    def __init__(self): super().__init__(); self.urls=[]
    def handle_starttag(self,tag,attrs):
        for key,value in attrs:
            if key in ('src','href') and value: self.urls.append(value)
pages=list(root.rglob('*.html'))
errors=[]
for path in pages:
    parser=Links(); parser.feed(path.read_text(encoding='utf-8'))
    for url in parser.urls:
        part=urlsplit(url)
        if part.scheme or part.netloc or not part.path: continue
        target=(path.parent/unquote(part.path)).resolve()
        if not target.is_relative_to(root) or not target.exists(): errors.append((str(path.relative_to(root)),url))
catalogue=json.loads((root/'themes.json').read_text(encoding='utf-8'))
assert len(catalogue)==114 and len(pages)==229
assert not list(root.rglob('*.db')) and not list(root.rglob('.secret'))
assert next(p for p in catalogue if p['key']=='saffron')['colors'][0]=='#c1440e'
assert next(p for p in catalogue if p['key']=='minimal')['mode']=='light'
assert not errors, errors[:20]
print(f'PASS: {len(pages)} HTML pages, {len(catalogue)} themes, all local links, intended palettes, no database or signing key.')
