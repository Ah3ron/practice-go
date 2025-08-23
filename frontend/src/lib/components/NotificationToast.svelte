<script lang="ts">
  import { notifications, type Notification } from '$lib/stores/notifications';

  export let notification: Notification;

  const typeClasses = {
    success: 'alert-success',
    error: 'alert-error',
    warning: 'alert-warning',
    info: 'alert-info',
  };

  const typeIcons = {
    success: 'M5 13l4 4L19 7',
    error: 'M6 18L18 6M6 6l12 12',
    warning: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16c-.77.833.192 2.5 1.732 2.5z',
    info: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
  };

  function handleClose() {
    notifications.remove(notification.id);
  }
</script>

<div class="alert {typeClasses[notification.type]} shadow-lg mb-2 animate-slide-up">
  <svg class="w-6 h-6 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={typeIcons[notification.type]} />
  </svg>
  
  <div class="flex-1">
    <h3 class="font-bold">{notification.title}</h3>
    {#if notification.message}
      <div class="text-xs opacity-75">{notification.message}</div>
    {/if}
  </div>
  
  <button class="btn btn-sm btn-ghost" on:click={handleClose}>
    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
    </svg>
  </button>
</div>