<script lang="ts">
  import { authAPI } from '$lib/services/api';
  import { notifications } from '$lib/stores/notifications';

  let username = '';
  let email = '';
  let password = '';
  let names = '';
  let loading = false;

  async function handleRegister() {
    if (!username.trim() || !email.trim() || !password.trim()) {
      notifications.add({
        type: 'warning',
        title: 'Validation Error',
        message: 'Please fill in all required fields'
      });
      return;
    }

    if (password.length < 6) {
      notifications.add({
        type: 'warning',
        title: 'Validation Error',
        message: 'Password must be at least 6 characters long'
      });
      return;
    }

    loading = true;

    try {
      await authAPI.register({
        username: username.trim(),
        email: email.trim(),
        password,
        names: names.trim() || undefined
      });

      notifications.add({
        type: 'success',
        title: 'Account Created!',
        message: 'You can now sign in with your credentials'
      });

      window.location.href = '/login';
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: 'Registration Failed',
        message: error.response?.data?.message || 'Failed to create account'
      });
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Register - Resource Manager</title>
</svelte:head>

<div
  class="flex min-h-screen items-center justify-center bg-gradient-to-br from-primary/10 to-secondary/10 p-4"
>
  <div class="card w-full max-w-md bg-base-100 shadow-2xl">
    <div class="card-body">
      <div class="mb-8 text-center">
        <div
          class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-primary"
        >
          <svg
            class="h-8 w-8 text-primary-content"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
            />
          </svg>
        </div>
        <h2 class="text-2xl font-bold text-base-content">Create Account</h2>
        <p class="mt-2 text-base-content/60">Join us today</p>
      </div>

      <form on:submit|preventDefault={handleRegister} class="space-y-4">
        <div class="form-control">
          <label class="label" for="username">
            <span class="label-text">Username <span class="text-error">*</span></span>
          </label>
          <input
            id="username"
            type="text"
            bind:value={username}
            placeholder="Enter username"
            class="input-bordered input w-full"
            disabled={loading}
            required
          />
        </div>

        <div class="form-control">
          <label class="label" for="email">
            <span class="label-text">Email <span class="text-error">*</span></span>
          </label>
          <input
            id="email"
            type="email"
            bind:value={email}
            placeholder="Enter email address"
            class="input-bordered input w-full"
            disabled={loading}
            required
          />
        </div>

        <div class="form-control">
          <label class="label" for="names">
            <span class="label-text">Full Name</span>
          </label>
          <input
            id="names"
            type="text"
            bind:value={names}
            placeholder="Enter your full name"
            class="input-bordered input w-full"
            disabled={loading}
          />
        </div>

        <div class="form-control">
          <label class="label" for="password">
            <span class="label-text">Password <span class="text-error">*</span></span>
          </label>
          <input
            id="password"
            type="password"
            bind:value={password}
            placeholder="Enter password (min 6 characters)"
            class="input-bordered input w-full"
            disabled={loading}
            required
          />
        </div>

        <button
          type="submit"
          class="btn mt-6 w-full btn-primary"
          class:loading
          disabled={loading || !username.trim() || !email.trim() || !password.trim()}
        >
          {#if loading}
            <span class="loading loading-sm loading-spinner"></span>
            Creating account...
          {:else}
            Create Account
          {/if}
        </button>
      </form>

      <div class="divider">or</div>

      <p class="text-center text-sm text-base-content/60">
        Already have an account?
        <a href="/login" class="link font-medium link-primary">Sign in</a>
      </p>
    </div>
  </div>
</div>
