<script lang="ts">
  import '../app.css';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import Layout from '$lib/components/Layout.svelte';

  $: isAuthPage = ['/login', '/register'].includes($page.url.pathname);

  onMount(() => {
    auth.initialize();
  });

  // Redirect to login if not authenticated and not on auth page
  $: if (!$auth.loading && !$auth.isAuthenticated && !isAuthPage && typeof window !== 'undefined') {
    goto('/login', { replaceState: true });
  }

  // Redirect to dashboard if authenticated and on auth page
  $: if (!$auth.loading && $auth.isAuthenticated && isAuthPage && typeof window !== 'undefined') {
    goto('/', { replaceState: true });
  }
</script>

<svelte:head>
  <title>Resource Manager</title>
  <meta name="description" content="Modern resource management system" />
</svelte:head>

{#if !$auth.initialized}
  <div class="min-h-screen flex items-center justify-center bg-base-100">
    <div class="text-center">
      <div class="loading loading-spinner loading-lg mb-4"></div>
      <p class="text-base-content/60">Loading...</p>
    </div>
  </div>
{:else if isAuthPage}
  <slot />
{:else if $auth.isAuthenticated}
  <Layout>
    <slot />
  </Layout>
{:else}
  <!-- This will trigger the redirect to login -->
  <div class="min-h-screen flex items-center justify-center bg-base-100">
    <div class="text-center">
      <div class="loading loading-spinner loading-lg mb-4"></div>
      <p class="text-base-content/60">Redirecting...</p>
    </div>
  </div>
{/if}