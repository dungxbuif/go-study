import { Buffer } from 'node:buffer'

let bn = 75n
console.log(bn.toString(10))

bn = BigInt('75')
console.log(bn.toString(10))

bn = BigInt(0x4b)
console.log(bn.toString(10))

bn = BigInt('0x4b')
console.log(bn.toString(10))

bn = BigInt('0x' + Buffer.from('4b', 'hex').toString('hex'))
console.log(bn.toString(10))
console.log(Number(bn))
console.log(bn.toString(16))
console.log(Buffer.from(bn.toString(16), 'hex'))

const bn2 = BigInt(100)
const isEqual = bn === bn2
console.log(isEqual)

const isGreater = bn > bn2
console.log(isGreater)

const isLesser = bn < bn2
console.log(isLesser)
