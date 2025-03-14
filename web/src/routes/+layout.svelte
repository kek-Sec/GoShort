<script>
	import '../app.css';

	// Get branding config from data store using runes syntax
	let { data, children } = $props();
	const { brandConfig } = data;

	// Apply theme colors using CSS variables
	const customStyles = `
		:root {
			--brand-primary: ${brandConfig.primaryColor};
			--brand-secondary: ${brandConfig.secondaryColor};
			--brand-theme: ${brandConfig.themeColor};
		}
	`;
</script>

<!-- Inject custom styles -->
<svelte:head>
	<!-- Update page title -->
	<title>{brandConfig.title}</title>

	<!-- Update meta tags -->
	<meta name="description" content={brandConfig.description} />
	<meta name="keywords" content={brandConfig.keywords} />
	<meta name="author" content={brandConfig.author} />
	<meta name="theme-color" content={brandConfig.themeColor} />

	<!-- Apply custom styles -->
	<style>{customStyles}</style>
</svelte:head>

{@render children()}
