import { openSync, writeSync, readSync, closeSync, unlinkSync } from 'node:fs'
import { Buffer } from 'node:buffer'

// create file
const createFd = openSync('test.txt', 'w')
closeSync(createFd)

// open file (returns file descriptor)
const fd = openSync('test.txt', 'r+')

const wbuf = Buffer.from('hello world.')
const rbuf = Buffer.alloc(12)
const off = 0
const len = 12
const pos = 0

// write file
writeSync(fd, wbuf, 0, wbuf.length, pos)

// read file
readSync(fd, rbuf, off, len, pos)

console.log(rbuf.toString())

// close file
closeSync(fd)

// delete file
unlinkSync('test.txt')
