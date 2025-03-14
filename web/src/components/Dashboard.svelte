<script>
  import { page } from '$app/stores';
  import { onMount } from 'svelte';

  // Get brand config from layout data
  const brandConfig = $page.data.brandConfig;

  // User data
  let user = { username: 'Loading...' };
  let urls = [];
  let isLoading = true;
  let error = null;

  // Pagination
  let currentPage = 1;
  let totalPages = 1;
  let pageSize = 5;

  onMount(async () => {
    await fetchUserData();
    await fetchUserUrls(currentPage);
  });

  async function fetchUserData() {
    try {
      const token = localStorage.getItem('authToken');
      if (!token) throw new Error('No auth token found');

      const response = await fetch('/api/v1/user/profile', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) throw new Error('Failed to fetch user data');

      const data = await response.json();
      user = data;
    } catch (error) {
      console.error('Error fetching user data:', error);
      error = 'Could not load user profile';
    }
  }

  async function fetchUserUrls(page) {
    isLoading = true;
    try {
      const token = localStorage.getItem('authToken');
      if (!token) throw new Error('No auth token found');

      const response = await fetch(`/api/v1/urls?page=${page}&limit=${pageSize}`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) throw new Error('Failed to fetch URLs');

      const data = await response.json();
      urls = data.urls || [];
      totalPages = data.totalPages || 1;
      isLoading = false;
    } catch (error) {
      console.error('Error fetching URLs:', error);
      error = 'Could not load your URLs';
      isLoading = false;
    }
  }

  function changePage(page) {
    if (page < 1 || page > totalPages) return;
    currentPage = page;
    fetchUserUrls(currentPage);
  }

  // Format date in a user-friendly way
  function formatDate(dateString) {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(date);
  }

  // Function to handle logout
  function handleLogout() {
    localStorage.removeItem('authToken');
    window.location.reload();
  }

  // Copy URL to clipboard
  async function copyUrl(url) {
    try {
      await navigator.clipboard.writeText(url);
      // You could add a visual feedback here if needed
    } catch (err) {
      console.error('Failed to copy URL:', err);
    }
  }
</script>

<div class="bg-white rounded-lg shadow-lg overflow-hidden">
  <!-- Dashboard Header -->
  <div class="bg-gradient-to-r from-brand-primary to-brand-secondary p-6 text-white">
    <div class="flex justify-between items-center">
      <div>
        <h2 class="text-2xl font-bold">Welcome back, {user.username || 'User'}</h2>
        <p class="opacity-90">Manage your shortened URLs</p>
      </div>
      <button
        on:click={handleLogout}
        class="bg-white text-brand-primary px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors duration-200 font-medium flex items-center gap-2"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M3 3a1 1 0 00-1 1v12a1 1 0 001 1h12a1 1 0 001-1V4a1 1 0 00-1-1H3zm9 3a1 1 0 112 0v4a1 1 0 11-2 0V6zm-6 9a1 1 0 100-2 1 1 0 000 2zm10 0a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
        </svg>
        Logout
      </button>
    </div>
  </div>

  <!-- Dashboard Content -->
  <div class="p-6">
    <h3 class="text-xl font-semibold text-gray-800 mb-4">Your Shortened URLs</h3>

    <!-- URL Table -->
    {#if isLoading}
      <div class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-brand-primary"></div>
      </div>
    {:else if error}
      <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        <p>{error}</p>
      </div>
    {:else if urls.length === 0}
      <div class="bg-gray-50 rounded-lg p-8 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
        <h4 class="text-lg font-medium text-gray-900 mt-4">No URLs yet!</h4>
        <p class="text-gray-600 mt-2">Start shortening URLs to see them here.</p>
      </div>
    {:else}
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead>
            <tr>
              <th class="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Short URL</th>
              <th class="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Original URL</th>
              <th class="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created</th>
              <th class="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Expiry</th>
              <th class="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Clicks</th>
              <th class="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            {#each urls as url}
              <tr class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <a
                      href={`${window.location.origin}/${url.shortUrl}`}
                      target="_blank"
                      class="text-brand-primary hover:underline"
                    >
                      {url.shortUrl}
                    </a>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="max-w-xs truncate" title={url.longUrl}>
                    {url.longUrl}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDate(url.createdAt)}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {url.expiresAt ? formatDate(url.expiresAt) : 'Never'}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {url.clicks || 0}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button
                    on:click={() => copyUrl(`${window.location.origin}/${url.shortUrl}`)}
                    class="text-brand-primary hover:text-opacity-75"
                  >
                    Copy
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      {#if totalPages > 1}
        <div class="flex justify-center mt-6">
          <nav class="flex items-center space-x-2">
            <button
              on:click={() => changePage(currentPage - 1)}
              class="px-3 py-1 rounded border {currentPage === 1 ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'bg-white text-gray-700 hover:bg-gray-50'}"
              disabled={currentPage === 1}
            >
              Previous
            </button>

            {#each Array(totalPages) as _, i}
              <button
                on:click={() => changePage(i + 1)}
                class="px-3 py-1 rounded border {currentPage === i + 1 ? 'bg-brand-primary text-white' : 'bg-white text-gray-700 hover:bg-gray-50'}"
              >
                {i + 1}
              </button>
            {/each}

            <button
              on:click={() => changePage(currentPage + 1)}
              class="px-3 py-1 rounded border {currentPage === totalPages ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'bg-white text-gray-700 hover:bg-gray-50'}"
              disabled={currentPage === totalPages}
            >
              Next
            </button>
          </nav>
        </div>
      {/if}
    {/if}
  </div>
</div>
