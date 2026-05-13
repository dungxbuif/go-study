import { exec } from 'node:child_process';
import process from 'node:process';

exec("echo 'hello world'", (error, stdout, stderr) => {
	if (error) {
		console.error(error);
	}

	if (stderr) {
		console.error(stderr);
	}

	if (stdout) {
		process.stdout.write(stdout);
	}
});
