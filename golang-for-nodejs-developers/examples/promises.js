import { setTimeout } from 'node:timers/promises';

async function asyncMethod(value) {
  await setTimeout(1000);
  return `resolved: ${value}`;
}

async function main() {
  asyncMethod('foo')
    .then(result => console.log(result))
    .catch(err => console.error(err));

  try {
    const results = await Promise.all([
      asyncMethod('A'),
      asyncMethod('B'),
      asyncMethod('C')
    ]);
    console.log(results);
  } catch (err) {
    console.error(err);
  }
}

main();
