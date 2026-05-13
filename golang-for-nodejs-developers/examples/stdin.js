import process from 'node:process';

const stdin = process.stdin;

process.stdout.write('Enter name: ');

stdin.on('data', data => {
	const name = data.toString().trim();
	console.log(`Your name is: ${name}`);
	process.exit();
});
