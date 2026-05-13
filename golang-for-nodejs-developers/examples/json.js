let jsonstr = '{"foo":"bar"}'

const parsed = JSON.parse(jsonstr)
console.log(parsed)

jsonstr = JSON.stringify(parsed)
console.log(jsonstr)
