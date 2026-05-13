const map = new Map()
map.set('foo', 'bar')

let found = map.has('foo')
console.log(found)

let item = map.get('foo')
console.log(item)

map.delete('foo')

found = map.has('foo')
console.log(found)

item = map.get('foo')
console.log(item)

const map2 = {}
map2['foo'] = 'bar'
item = map2['foo']
delete map2['foo']

const map3 = new Map([
  ['foo', 100],
  ['bar', 200],
  ['baz', 300]
])

for (const [key, value] of map3) {
  console.log(key, value)
}
