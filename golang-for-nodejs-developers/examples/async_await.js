import { setTimeout } from 'node:timers/promises';

async function hello(name) {
  await setTimeout(1000);
  if (name === 'fail') {
    throw new Error('failed');
  }
  return `hello ${name}`;
}

async function main() {
  try {
    let output = await hello('bob');
    console.log(output);

    output = await hello('fail');
    console.log(output);
  } catch (err) {
    console.error(err.message);
  }
}

main();
