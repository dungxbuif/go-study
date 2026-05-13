export function* generator() {
  yield 'hello';
  yield 'world';
}

const gen = generator();

while (true) {
  const { value, done } = gen.next();
  console.log(value, done);

  if (done) {
    break;
  }
}

// alternatively
for (const value of generator()) {
  console.log(value);
}
