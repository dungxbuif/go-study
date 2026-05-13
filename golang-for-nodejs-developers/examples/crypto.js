import crypto from 'node:crypto'
import { Buffer } from 'node:buffer'

const hash = crypto.createHash('sha256').update(Buffer.from('hello')).digest()

console.log(hash.toString('hex'))
