import { setTimeout } from 'node:timers';

const callback = () => {
  console.log('called');
};

setTimeout(callback, 1000);
