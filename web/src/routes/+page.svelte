<script>
  import Main from '../components/Main.svelte';
  import Footer from '../components/Footer.svelte';
  import Login from '../components/Login.svelte';

  let longUrl = '';
  let shortUrl = '';
  let customUrl = '';
  let errorMessage = '';
  let expiry = '';
  let validationError = '';
  let isCopied = false;
  let showAccordion = false;
  let showLoginModal = false;

  const shortenUrl = async () => {
    if (!longUrl) {
      validationError = 'Please enter a valid URL.';
      return;
    }
    validationError = '';

    try {
      const response = await fetch('/api/v1/shorten', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ long_url: longUrl, custom_url: customUrl, expiry }),
      });

      if (!response.ok) {
        const error = await response.text();
        if (response.status === 409) {
          errorMessage = 'The custom URL is already taken. Please try another.';
        } else {
          errorMessage = 'Error: Could not shorten the URL.';
        }
        throw new Error(errorMessage);
      }

      const data = await response.json();
      shortUrl = `${window.location.origin}/${data.short_url}`;
      showAccordion = false; // Collapse accordion
    } catch (error) {
      console.error(error);
    }
  };

  const copyToClipboard = async () => {
    if (shortUrl) {
      await navigator.clipboard.writeText(shortUrl);
      isCopied = true;
      setTimeout(() => (isCopied = false), 2000);
    }
  };
  
  const toggleLoginModal = () => {
    showLoginModal = !showLoginModal;
  };

  const closeLoginModal = () => {
    showLoginModal = false;
  };
</script>

<div class="flex flex-col min-h-screen bg-gray-100">
  <!-- Login Modal -->
  {#if showLoginModal}
    <div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div class="relative w-full max-w-md">
        <!-- Close button -->
        <button 
          on:click={closeLoginModal}
          class="absolute -top-12 right-0 text-white hover:text-gray-300"
          aria-label="Close login modal"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <Login isModal={true} on:close={closeLoginModal} />
      </div>
    </div>
  {/if}

  <!-- Login button at the top right -->
  <div class="absolute top-4 right-4 z-10">
    <button
      on:click={toggleLoginModal}
      class="group bg-brand-secondary hover:bg-opacity-90 text-white p-2.5 rounded-full shadow-md transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-brand-secondary focus:ring-offset-2"
      aria-label="Login"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
        <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
      </svg>
      <span class="invisible opacity-0 group-hover:visible group-hover:opacity-100 absolute -left-2 top-full mt-2 w-max bg-gray-800 text-white text-sm rounded px-3 py-1 transition-opacity duration-300 pointer-events-none">
        Login
      </span>
    </button>
  </div>
  
  <main class="flex-grow flex items-center justify-center">
    <Main
      bind:longUrl
      bind:shortUrl
      bind:customUrl
      bind:errorMessage
      bind:expiry
      bind:validationError
      bind:isCopied
      bind:showAccordion
      bind:showLoginModal
      {shortenUrl}
      {copyToClipboard}
    />
  </main>
  <Footer />
</div>
