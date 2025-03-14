<script>
	import '../app.css';

	// Get branding config from data store using runes syntax
	let { data, children } = $props();
	const { brandConfig } = data;

	$effect(() => {
		// Apply theme colors using CSS variables directly to document root
		document.documentElement.style.setProperty('--brand-primary', brandConfig.primaryColor);
		document.documentElement.style.setProperty('--brand-secondary', brandConfig.secondaryColor);
		document.documentElement.style.setProperty('--brand-theme', brandConfig.themeColor);
	});
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
</svelte:head>

{@render children()}
