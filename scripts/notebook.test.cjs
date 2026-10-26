const test = require('node:test');
const assert = require('node:assert/strict');
const {restore,tableNumber,total} = require('../web/static/js/notebook.js');
test('normalizes Persian/Arabic table numbers and rejects ambiguous or hostile input', () => {
  for (const [input,expected] of [[' ۰۱۲ ','12'],['٠٠٧','7'],['9999','9999'],['0',''],['-1',''],['1.5',''],['1e3',''],['10000',''],['<img>','']]) assert.equal(tableNumber(input),expected);
});
test('restored notes cannot resurrect unavailable dishes, duplicate rows or invalid quantities', () => {
  const restored=restore({items:[{id:'1',qty:2,note:'x'.repeat(400)},{id:'1',qty:7},{id:'2',qty:1},{id:'3',qty:0},{id:'4',qty:1.5},{id:'5',qty:100}],note:'z'.repeat(900),table:'۱۲'},new Map([['1',{}],['3',{}],['4',{}],['5',{}]]));
  assert.equal(restored.items.length,1); assert.equal(restored.items[0].qty,2); assert.equal(restored.items[0].note.length,200); assert.equal(restored.note.length,500); assert.equal(restored.table,'12');
  for(const raw of [null,[],42,'text',{}]) assert.deepEqual(restore(raw,new Map()),{items:[],note:'',table:''});
});
test('uses current prices with exact arithmetic; suppresses totals when a price is hidden',()=>{
  const items=Array.from({length:200},(_,i)=>({id:String(i),qty:99}));
  const catalog=new Map(items.map(item=>[item.id,{price:'1000000000000'}]));
  assert.equal(total(items,catalog),19800000000000000n);
  catalog.get('0').price=''; assert.equal(total(items,catalog),null);
  catalog.get('0').price='900'; assert.equal(total([{id:'0',qty:2}],catalog),1800n);
});
test('WordPress ships the same notebook behavior and styling',()=>{
  const fs=require('node:fs');
  for(const [folder,file] of [['js','notebook.js'],['css','notebook.css']]) assert.equal(fs.readFileSync(`web/static/${folder}/${file}`,'utf8'),fs.readFileSync(`wordpress/menucraft/assets/${file}`,'utf8'));
});
