<script>
	export let longUrl = '';
	export let shortUrl = '';
	export let customUrl = '';
	export let errorMessage = '';
	export let expiry = '';
	export let validationError = '';
	export let isCopied = false;
	export let showAccordion = false;

	// Function to validate the custom URL input
	const validateCustomUrl = (url) => {
		const customUrlRegex = /^[a-zA-Z0-9_-]+$/;
		return customUrlRegex.test(url);
	};

	// Function to ensure URL has https://
	const ensureHttps = (url) => {
		if (!url.startsWith('http')) {
			return `https://${url}`;
		}
		return url;
	};

	const handleInputChange = (event) => {
		longUrl = ensureHttps(event.target.value);
	};

	export let shortenUrl = async () => {
		if (!longUrl) {
			validationError = 'Please enter a valid URL.';
			return;
		}

		if (customUrl && !validateCustomUrl(customUrl)) {
			validationError =
				'Custom URL contains spaces or invalid characters. Only alphanumeric, dashes, and underscores are allowed.';
			return;
		}
		validationError = '';

		try {
			errorMessage = '';
			shortUrl = '';

			shortUrl = `${window.location.origin}/${customUrl || 'generated-url'}`;

			showAccordion = false;
		} catch (error) {
			console.error('Error shortening URL:', error);
		}
	};

	export let copyToClipboard = async () => {
		if (shortUrl) {
			try {
				await navigator.clipboard.writeText(shortUrl);
				isCopied = true;
				setTimeout(() => (isCopied = false), 2000);
			} catch (error) {
				console.error('Failed to copy URL:', error);
			}
		}
	};
</script>

<div class="flex flex-col items-center justify-center bg-gray-100 p-4">
	<h1 class="text-4xl font-bold text-blue-600 mb-6">GoShort - URL Shortener</h1>

	<div class="w-full max-w-md space-y-4">
		<!-- Long URL Input -->
		<div>
			<label for="long-url" class="block text-gray-700 font-medium mb-2">Enter URL</label>
			<input
				id="long-url"
				type="url"
				bind:value={longUrl}
				on:input={handleInputChange}
				placeholder="Enter your URL"
				class="w-full p-3 border border-gray-300 rounded-lg shadow-sm focus:ring focus:ring-blue-300 focus:outline-none"
				required
			/>
			{#if validationError}
				<p class="text-red-500 text-sm mt-1">{validationError}</p>
			{/if}
		</div>

		<!-- Accordion for Optional Inputs -->
		<div>
			<button
				type="button"
				on:click={() => (showAccordion = !showAccordion)}
				class="w-full bg-gray-200 text-gray-700 font-semibold py-2 px-4 rounded-lg hover:bg-gray-300 focus:ring focus:ring-gray-300 focus:outline-none transition flex items-center justify-between"
			>
				<span>Additional Options (Optional)</span>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-5 w-5 transform transition-transform"
					viewBox="0 0 20 20"
					fill="currentColor"
					class:rotate-180={showAccordion}
				>
					<path
						fill-rule="evenodd"
						d="M10 3a1 1 0 01.707.293l6 6a1 1 0 11-1.414 1.414L10 5.414 4.707 10.707A1 1 0 013.293 9.293l6-6A1 1 0 0110 3z"
						clip-rule="evenodd"
					/>
				</svg>
			</button>
			{#if showAccordion}
				<div class="mt-2 space-y-4">
					<!-- Custom URL Input -->
					<div>
						<label for="custom-url" class="block text-gray-700 font-medium mb-2">Custom URL</label>
						<input
							id="custom-url"
							type="text"
							bind:value={customUrl}
							placeholder="Enter custom URL (e.g., my-short-url)"
							class="w-full p-3 border border-gray-300 rounded-lg shadow-sm focus:ring focus:ring-blue-300 focus:outline-none"
						/>
					</div>

					<!-- Expiry Date Input -->
					<div>
						<label for="expiry" class="block text-gray-700 font-medium mb-2">Expiry Date & Time</label>
						<input
							id="expiry"
							type="datetime-local"
							bind:value={expiry}
							class="w-full p-3 border border-gray-300 rounded-lg shadow-sm focus:ring focus:ring-blue-300 focus:outline-none"
						/>
					</div>
				</div>
			{/if}
		</div>

		<!-- Shorten Button -->
		<button
			on:click={shortenUrl}
			class="w-full bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg hover:bg-blue-700 focus:ring focus:ring-blue-300 focus:outline-none transition"
		>
			Shorten
		</button>
	</div>

	<!-- Shortened URL Section -->
	{#if shortUrl}
		<div class="mt-6 text-center bg-white p-6 rounded-lg shadow-lg">
			<p class="text-lg font-medium mb-4">Your shortened URL is ready! 🚀</p>
			<div class="flex flex-col items-center space-y-4">
				<div
					class="bg-gray-100 text-gray-800 px-4 py-2 rounded-md border border-gray-300 font-mono text-sm break-all w-full max-w-lg"
				>
					<a href={shortUrl} target="_blank" class="hover:text-blue-500">{shortUrl}</a>
				</div>
				<button
					on:click={copyToClipboard}
					class="bg-green-500 text-white font-semibold py-2 px-4 rounded-lg hover:bg-green-600 focus:ring focus:ring-green-300 focus:outline-none transition flex items-center space-x-2"
				>
					{#if isCopied}
						<span class="animate-pulse">Copied! ✅</span>
					{:else}
						<span>📋 Copy URL</span>
					{/if}
				</button>
			</div>
		</div>
	{/if}

	<!-- Error Message -->
	{#if errorMessage}
		<div class="mt-4 text-red-500 font-medium">{errorMessage}</div>
	{/if}
</div>
