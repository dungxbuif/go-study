const array = [1, 2, 3, 4, 5]
console.log(array)

const clone = [...array]
console.log(clone)

const sub = array.slice(2, 4)
console.log(sub)

const concatenated = [...clone, 6, 7]
console.log(concatenated)

const prepended = [-2, -1, 0, ...concatenated]
console.log(prepended)
