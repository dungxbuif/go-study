import { parseArgs } from 'node:util';

const { values } = parseArgs({
	options: {
		foo: {
			type: 'string',
			default: 'default value',
		},
		qux: {
			type: 'boolean',
			default: false,
		},
	},
});

console.log('foo:', values.foo);
console.log('qux:', values.qux);
