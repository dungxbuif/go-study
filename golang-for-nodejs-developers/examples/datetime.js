const nowUnix = Date.now();
console.log(nowUnix);

const datestr = '2019-01-17T09:24:23+00:00';
const date = new Date(datestr);
console.log(date.getTime());
console.log(date.toString());

const futureDate = new Date(date);
futureDate.setDate(date.getDate() + 14);
console.log(futureDate.toString());

// Modern way to format
const formatted = date.toLocaleDateString('en-US', {
  month: '2-digit',
  day: '2-digit',
  year: 'numeric',
});
console.log(formatted);
