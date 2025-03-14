<script>
  import { onMount } from 'svelte';
  import Main from '../components/Main.svelte';
  import Footer from '../components/Footer.svelte';
  import Login from '../components/Login.svelte';
  import Dashboard from '../components/Dashboard.svelte';

  let longUrl = '';
  let shortUrl = '';
  let customUrl = '';
  let errorMessage = '';
  let expiry = '';
  let validationError = '';
  let isCopied = false;
  let showAccordion = false;
  let showLoginModal = false;
  let isLoggedIn = false;
  let userData = null;

  onMount(() => {
    // Check if user is logged in on page load
    const token = localStorage.getItem('authToken');
    if (token) {
      isLoggedIn = true;
      // You could fetch user data here or let the Dashboard component handle it
    }
  });

  const shortenUrl = async () => {
    if (!longUrl) {
      validationError = 'Please enter a valid URL.';
      return;
    }
    validationError = '';

    try {
      const response = await fetch('/api/v1/shorten', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          // Include auth token if available
          ...( localStorage.getItem('authToken') ? 
               { 'Authorization': `Bearer ${localStorage.getItem('authToken')}` } : 
               {} )
        },
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

  const handleLoginSuccess = (event) => {
    userData = event.detail.user;
    isLoggedIn = true;
    showLoginModal = false;
  };

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    isLoggedIn = false;
    userData = null;
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
        <Login isModal={true} on:close={closeLoginModal} on:loginsuccess={handleLoginSuccess} />
      </div>
    </div>
  {/if}

  <!-- Login/User button at the top right -->
  <div class="absolute top-4 right-4 z-10">
    {#if isLoggedIn}
      <button
        on:click={handleLogout}
        class="group bg-brand-secondary hover:bg-opacity-90 text-white p-2.5 rounded-full shadow-md transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-brand-secondary focus:ring-offset-2"
        aria-label="Logout"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M3 3a1 1 0 00-1 1v12a1 1 0 001 1h12a1 1 0 001-1V4a1 1 0 00-1-1H3zm9 3a1 1 0 112 0v4a1 1 0 11-2 0V6zm-6 9a1 1 0 100-2 1 1 0 000 2zm10 0a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
        </svg>
        <span class="invisible opacity-0 group-hover:visible group-hover:opacity-100 absolute -left-2 top-full mt-2 w-max bg-gray-800 text-white text-sm rounded px-3 py-1 transition-opacity duration-300 pointer-events-none">
          Logout
        </span>
      </button>
    {:else}
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
    {/if}
  </div>
  
  <main class="flex-grow flex items-center justify-center p-4">
    {#if isLoggedIn}
      <!-- Show Dashboard when logged in -->
      <Dashboard />
    {:else}
      <!-- Show URL Shortener when not logged in -->
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
    {/if}
  </main>
  <Footer />
</div>
