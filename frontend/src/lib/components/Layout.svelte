<script lang="ts">
  import { page } from '$app/stores';
  import { auth } from '$lib/stores/auth';
  import { notifications } from '$lib/stores/notifications';
  import Sidebar from './Sidebar.svelte';
  import Navbar from './Navbar.svelte';
  import NotificationToast from './NotificationToast.svelte';
</script>

<div class="drawer min-h-screen bg-base-100 lg:drawer-open">
  <input id="drawer-toggle" type="checkbox" class="drawer-toggle" />

  <div class="drawer-content flex flex-col">
    <!-- Navbar -->
    <Navbar />

    <!-- Page content -->
    <main class="custom-scrollbar flex-1 overflow-y-auto p-6">
      <div class="mx-auto max-w-7xl">
        <slot />
      </div>
    </main>
  </div>

  <!-- Sidebar -->
  <div class="drawer-side">
    <label for="drawer-toggle" class="drawer-overlay"></label>
    <Sidebar />
  </div>
</div>

<!-- Notifications -->
<div class="toast toast-top toast-end z-50">
  {#each $notifications as notification (notification.id)}
    <NotificationToast {notification} />
  {/each}
</div>

<style>
  :global(.drawer-content) {
    scroll-behavior: smooth;
  }
</style>
