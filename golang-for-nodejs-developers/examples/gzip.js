import zlib from 'node:zlib'
import { Buffer } from 'node:buffer'

const data = Buffer.from('hello world', 'utf-8')

zlib.gzip(data, (err, compressed) => {
  if (err) {
    console.error(err)
  }

  console.log(compressed)

  zlib.unzip(compressed, (err, decompressed) => {
    if (err) {
      console.error(err)
    }

    console.log(decompressed.toString())
  })
})
