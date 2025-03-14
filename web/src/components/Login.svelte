<script>
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { browser } from '$app/environment';

  // Get the brand config from the layout data
  const brandConfig = $page.data.brandConfig;

  // Form data
  let username = '';
  let password = '';
  let isLoading = false;
  let errorMessage = '';

  // Form validation
  let touched = {
    username: false,
    password: false
  };

  $: usernameError = touched.username && !username ? 'Username is required' : '';
  $: passwordError = touched.password && !password ? 'Password is required' : '';
  $: formIsValid = username && password && !usernameError && !passwordError;

  function handleBlur(field) {
    touched[field] = true;
  }

  async function handleLogin() {
    if (!formIsValid) {
      // Mark all fields as touched to show validation errors
      Object.keys(touched).forEach(key => touched[key] = true);
      return;
    }

    isLoading = true;
    errorMessage = '';

    try {
      // For static sites, using fetch might not work as expected without a server
      // This is a client-side only operation
      if (browser) {
        const response = await fetch('/api/v1/auth/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ username, password })
        });

        if (!response.ok) {
          const data = await response.json();
          throw new Error(data.message || 'Login failed');
        }

        const data = await response.json();

        // Store the JWT token
        localStorage.setItem('authToken', data.token);

        // Redirect to dashboard or home page
        goto('/');
      }
    } catch (error) {
      errorMessage = error.message || 'An error occurred during login';
      console.error('Login error:', error);
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="w-full max-w-md space-y-8 bg-white p-10 rounded-xl shadow-lg">
  <div class="text-center">
    <h2 class="mt-6 text-3xl font-extrabold text-gray-900">
      Sign in to your account
    </h2>
    <p class="mt-2 text-sm text-gray-600">
      Welcome back to {brandConfig.title}
    </p>
  </div>

  {#if errorMessage}
    <div class="bg-red-50 border-l-4 border-red-500 p-4 mb-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <p class="text-sm text-red-700">{errorMessage}</p>
        </div>
      </div>
    </div>
  {/if}

  <form class="mt-8 space-y-6" on:submit|preventDefault={handleLogin}>
    <div class="rounded-md shadow-sm space-y-4">
      <div>
        <label for="username" class="sr-only">Username</label>
        <input
          id="username"
          name="username"
          type="text"
          autocomplete="username"
          bind:value={username}
          on:blur={() => handleBlur('username')}
          class="appearance-none rounded-lg relative block w-full px-3 py-2 border {usernameError ? 'border-red-300' : 'border-gray-300'} placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-2 focus:ring-brand-primary focus:border-transparent"
          placeholder="Username"
        />
        {#if usernameError}
          <p class="mt-1 text-sm text-red-600">{usernameError}</p>
        {/if}
      </div>

      <div>
        <label for="password" class="sr-only">Password</label>
        <input
          id="password"
          name="password"
          type="password"
          autocomplete="current-password"
          bind:value={password}
          on:blur={() => handleBlur('password')}
          class="appearance-none rounded-lg relative block w-full px-3 py-2 border {passwordError ? 'border-red-300' : 'border-gray-300'} placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-2 focus:ring-brand-primary focus:border-transparent"
          placeholder="Password"
        />
        {#if passwordError}
          <p class="mt-1 text-sm text-red-600">{passwordError}</p>
        {/if}
      </div>
    </div>

    <div>
      <button
        type="submit"
        disabled={isLoading}
        class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-lg text-white bg-brand-primary hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-primary disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
      >
        {#if isLoading}
          <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Signing in...
        {:else}
          Sign in
        {/if}
      </button>
    </div>
  </form>

  <div class="text-center">
    <p class="mt-2 text-sm text-gray-600">
      Don't have an account?
      <a href="#" class="font-medium text-brand-primary hover:opacity-80">
        Sign up
      </a>
    </p>
  </div>
</div>
