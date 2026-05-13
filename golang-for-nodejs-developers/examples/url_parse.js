const urlstr = 'http://bob:secret@sub.example.com:8080/somepath?foo=bar';

const parsed = new URL(urlstr);
console.log(parsed.protocol);
console.log(`${parsed.username}:${parsed.password}`);
console.log(parsed.port);
console.log(parsed.hostname);
console.log(parsed.pathname);

const query = Object.fromEntries(parsed.searchParams);
console.log(query);
