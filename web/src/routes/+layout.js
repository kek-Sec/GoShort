import { loadConfig } from '$lib/config';

/** @type {import('./$types').LayoutLoad} */
export async function load() {
  const config = await loadConfig();
  return { brandConfig: config };
}
