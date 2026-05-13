let i = 0

const id = setInterval(callback, 1000)

function callback() {
  console.log('called', i)

  if (i === 3) {
    clearInterval(id)
  }

  i++
}
