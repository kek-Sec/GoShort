import containerQueries from '@tailwindcss/container-queries';
import forms from '@tailwindcss/forms';
import typography from '@tailwindcss/typography';

/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],

	theme: {
		extend: {}
	},
	
	// Add a safelist for dynamic classes
	safelist: [
		'text-brand-primary',
		'bg-brand-primary',
		'text-brand-secondary',
		'bg-brand-secondary',
		'ring-brand-primary',
		'ring-brand-secondary',
	],

	plugins: [typography, forms, containerQueries]
};
