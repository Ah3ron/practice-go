<script lang="ts">
  import '../app.css';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import Layout from '$lib/components/Layout.svelte';

  $: isAuthPage = ['/login', '/register'].includes($page.url.pathname);
  let authInitialized = false;

  onMount(() => {
    auth.initialize();
    authInitialized = true;
  });

  // Only handle redirects after auth is initialized
  $: if (authInitialized && $auth.initialized) {
    // Redirect to login if not authenticated and not on auth page
    if (!$auth.isAuthenticated && !isAuthPage && typeof window !== 'undefined') {
      goto('/login', { replaceState: true });
    }
    // Redirect to dashboard if authenticated and on auth page
    else if ($auth.isAuthenticated && isAuthPage && typeof window !== 'undefined') {
      goto('/', { replaceState: true });
    }
  }
</script>

<svelte:head>
  <title>Менеджер ресурсов</title>
  <meta name="description" content="Современная система управления ресурсами" />
</svelte:head>

{#if !$auth.initialized}
  <div class="flex min-h-screen items-center justify-center bg-base-100">
    <div class="text-center">
      <div class="loading mb-4 loading-lg loading-spinner"></div>
      <p class="text-base-content/60">Загрузка...</p>
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
  <div class="flex min-h-screen items-center justify-center bg-base-100">
    <div class="text-center">
      <div class="loading mb-4 loading-lg loading-spinner"></div>
      <p class="text-base-content/60">Перенаправление...</p>
    </div>
  </div>
{/if}
