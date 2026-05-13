import test from 'node:test';
import assert from 'node:assert/strict';

function sum(a, b) {
	return a + b;
}

test('sum', (t) => {
	const tt = [
		{ a: 1, b: 1, ret: 2 },
		{ a: 2, b: 3, ret: 5 },
		{ a: 5, b: 5, ret: 10 }
	];

	tt.forEach(item => {
		assert.strictEqual(sum(item.a, item.b), item.ret);
	});
});
