let input = 'foobar'
const replaced = input.replace(/foo(.*)/i, 'qux$1')
console.log(replaced)

const match = /o{2}/i.test(input)
console.log(match)

input = '111-222-333'
const matches = input.match(/([0-9]+)/gi)
console.log(matches)
