import { writable } from 'svelte/store';

export const theme = writable('dark');

theme.subscribe(value => {
  if (typeof document !== 'undefined') {
    document.documentElement.classList.toggle('dark', value === 'dark');
  }
});