import { loadConfig } from '$lib/config';

/** @type {import('./$types').LayoutLoad} */
export async function load({ url }) {
  const config = await loadConfig();

  // Add path info to help with active route detection
  const path = url.pathname;

  return {
    brandConfig: config,
    currentPath: path
  };
}
