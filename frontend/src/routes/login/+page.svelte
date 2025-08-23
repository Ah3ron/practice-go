<script lang="ts">
  import { goto } from '$app/navigation';
  import { authAPI } from '$lib/services/api';
  import { auth } from '$lib/stores/auth';
  import { notifications } from '$lib/stores/notifications';

  let username = '';
  let password = '';
  let loading = false;

  async function handleLogin() {
    if (!username.trim() || !password.trim()) {
      notifications.add({
        type: 'warning',
        title: 'Validation Error',
        message: 'Please fill in all fields',
      });
      return;
    }

    loading = true;
    
    try {
      const response = await authAPI.login({ username: username.trim(), password });
      
      // Update auth state
      auth.login(response.data.user, response.data.token);
      
      notifications.add({
        type: 'success',
        title: 'Welcome back!',
        message: `Hello, ${response.data.user.names || response.data.user.username}`,
      });
      
      // Use SvelteKit navigation with a small delay to ensure state propagation
      setTimeout(() => {
        goto('/', { replaceState: true });
      }, 100);
      
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: 'Login Failed',
        message: error.response?.data?.message || 'Invalid credentials',
      });
    } finally {
      loading = false;
    }
  }

  function handleKeyPress(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      handleLogin();
    }
  }
</script>

<svelte:head>
  <title>Login - Resource Manager</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary/10 to-secondary/10 p-4">
  <div class="card w-full max-w-md bg-base-100 shadow-2xl">
    <div class="card-body">
      <div class="text-center mb-8">
        <div class="w-16 h-16 bg-primary rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-primary-content" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        </div>
        <h2 class="text-2xl font-bold text-base-content">Welcome back</h2>
        <p class="text-base-content/60 mt-2">Sign in to your account</p>
      </div>

      <form on:submit|preventDefault={handleLogin} class="space-y-4">
        <div class="form-control">
          <label class="label" for="username">
            <span class="label-text">Username</span>
          </label>
          <input
            id="username"
            type="text"
            bind:value={username}
            on:keypress={handleKeyPress}
            placeholder="Enter your username"
            class="input input-bordered w-full"
            class:input-error={!username.trim() && loading}
            disabled={loading}
            required
          />
        </div>

        <div class="form-control">
          <label class="label" for="password">
            <span class="label-text">Password</span>
          </label>
          <input
            id="password"
            type="password"
            bind:value={password}
            on:keypress={handleKeyPress}
            placeholder="Enter your password"
            class="input input-bordered w-full"
            class:input-error={!password.trim() && loading}
            disabled={loading}
            required
          />
        </div>

        <button
          type="submit"
          class="btn btn-primary w-full mt-6"
          class:loading
          disabled={loading || !username.trim() || !password.trim()}
        >
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
            Signing in...
          {:else}
            Sign in
          {/if}
        </button>
      </form>

      <div class="divider">or</div>

      <p class="text-center text-sm text-base-content/60">
        Don't have an account?
        <a href="/register" class="link link-primary font-medium">Sign up</a>
      </p>
    </div>
  </div>
</div>