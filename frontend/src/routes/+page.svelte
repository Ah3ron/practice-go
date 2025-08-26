<script lang="ts">
  import { onMount } from 'svelte';
  import { auth } from '$lib/stores/auth';
  import { resourcesAPI } from '$lib/services/api';
  import { notifications } from '$lib/stores/notifications';
  import type { Resource } from '$lib/services/api';

  let resources: Resource[] = [];
  let loading = true;

  $: totalResources = resources?.length || 0;
  $: lowStockResources = (resources || []).filter((r) => (r?.quantity || 0) < 100).length;
  $: recentResources = resources.slice(-3).reverse(); // Last 3 resources added

  // Get current date and time
  const currentDate = new Date();
  const dateOptions: Intl.DateTimeFormatOptions = {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  };
  const timeOptions: Intl.DateTimeFormatOptions = {
    hour: '2-digit',
    minute: '2-digit'
  };

  onMount(async () => {
    try {
      const resourcesResponse = await resourcesAPI.getAll();
      resources = resourcesResponse.data;
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: 'Ошибка загрузки панели управления',
        message: error.response?.data?.message || 'Не удалось загрузить данные панели управления'
      });
    } finally {
      loading = false;
    }
  });
</script>

<svelte:head>
  <title>Панель управления - Менеджер ресурсов</title>
</svelte:head>

<div class="fade-in">
  <div class="mb-8">
    <h1 class="mb-2 text-3xl font-bold text-base-content">
      Добро пожаловать, {$auth.user?.names || $auth.user?.username}!
    </h1>
    <p class="text-base-content/60">Вот что происходит с вашими ресурсами сегодня.</p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <div class="loading loading-lg loading-spinner"></div>
      <span class="ml-3 text-base-content/60">Загрузка панели управления...</span>
    </div>
  {:else}
    <!-- Dashboard Overview -->
    <div class="mb-8 grid grid-cols-1 gap-6 lg:grid-cols-3">
      <!-- System Status -->
      <div class="card border border-base-200 bg-base-100 shadow-sm">
        <div class="card-body">
          <div class="flex items-center">
            <div
              class="bg-opacity-10 flex h-12 w-12 items-center justify-center rounded-lg bg-success"
            >
              <svg
                class="h-6 w-6 text-success"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <div class="ml-4">
              <h3 class="text-sm font-medium text-base-content/60">Статус системы</h3>
              <p class="text-lg font-semibold text-success">Система работает</p>
              <p class="text-xs text-base-content/50">
                {currentDate.toLocaleDateString('ru-RU', dateOptions)}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Resource Summary -->
      <div class="card border border-base-200 bg-base-100 shadow-sm">
        <div class="card-body">
          <div class="flex items-center">
            <div
              class="bg-opacity-10 flex h-12 w-12 items-center justify-center rounded-lg bg-primary"
            >
              <svg
                class="h-6 w-6 text-primary"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
                />
              </svg>
            </div>
            <div class="ml-4">
              <h3 class="text-sm font-medium text-base-content/60">Всего ресурсов</h3>
              <p class="text-lg font-semibold text-base-content">{totalResources}</p>
              <p class="text-xs text-base-content/50">активных в системе</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Alerts -->
      <div class="card border border-base-200 bg-base-100 shadow-sm">
        <div class="card-body">
          <div class="flex items-center">
            <div
              class="h-12 w-12 bg-{lowStockResources > 0
                ? 'warning'
                : 'success'} bg-opacity-10 flex items-center justify-center rounded-lg"
            >
              {#if lowStockResources > 0}
                <svg
                  class="h-6 w-6 text-warning"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16c-.77.833.192 2.5 1.732 2.5z"
                  />
                </svg>
              {:else}
                <svg
                  class="h-6 w-6 text-success"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              {/if}
            </div>
            <div class="ml-4">
              <h3 class="text-sm font-medium text-base-content/60">Уведомления</h3>
              <p class="text-lg font-semibold text-{lowStockResources > 0 ? 'warning' : 'success'}">
                {#if lowStockResources > 0}
                  {lowStockResources} с низким остатком
                {:else}
                  Все в порядке
                {/if}
              </p>
              <p class="text-xs text-base-content/50">
                {#if lowStockResources > 0}
                  требуют внимания
                {:else}
                  нет активных предупреждений
                {/if}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-8 lg:grid-cols-2">
      <!-- Recent Resources -->
      <div class="card border border-base-200 bg-base-100 shadow-sm">
        <div class="card-body">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="card-title">Последние ресурсы</h2>
            <a href="/resources" class="btn btn-ghost btn-sm">Посмотреть все</a>
          </div>

          {#if resources.length === 0}
            <div class="py-8 text-center">
              <svg
                class="mx-auto mb-4 h-12 w-12 text-base-content/30"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
                />
              </svg>
              <p class="text-base-content/60">Ресурсы не найдены</p>
              <a href="/resources" class="btn mt-2 btn-sm btn-primary">Добавить ресурс</a>
            </div>
          {:else}
            <div class="space-y-3">
              {#each resources.slice(0, 5) as resource, index (`resource-item-${index}`)}
                <div class="flex items-center justify-between rounded-lg bg-base-200 p-3">
                  <div>
                    <h3 class="font-medium text-base-content">
                      {resource?.name || 'Неизвестный ресурс'}
                    </h3>
                    <p class="text-sm text-base-content/60">
                      {resource?.quantity || 0}
                      {resource?.unit || 'единиц'}
                    </p>
                  </div>
                  <div class="text-right">
                    <div
                      class="badge badge-sm {(resource?.quantity || 0) < 100
                        ? 'badge-warning'
                        : 'badge-success'}"
                    >
                      {(resource?.quantity || 0) < 100 ? 'Низкий остаток' : 'В наличии'}
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="card border border-base-200 bg-base-100 shadow-sm">
        <div class="card-body">
          <h2 class="mb-4 card-title">Быстрые действия</h2>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <a
              href="/resources"
              class="btn flex items-center justify-start btn-outline btn-primary"
            >
              <svg class="mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 4v16m8-8H4"
                />
              </svg>
              Добавить ресурс
            </a>

            <a href="/analytics" class="btn flex items-center justify-start btn-outline btn-accent">
              <svg class="mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                />
              </svg>
              Посмотреть аналитику
            </a>

            <a href="/resources" class="btn flex items-center justify-start btn-outline btn-info">
              <svg class="mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"
                />
              </svg>
              Посмотреть отчёт
            </a>

            <a
              href="/resources"
              class="btn flex items-center justify-start btn-outline btn-secondary"
            >
              <svg class="mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M7 4V2a1 1 0 011-1h4a1 1 0 011 1v2h4a1 1 0 011 1v1a1 1 0 01-1 1H3a1 1 0 01-1-1V5a1 1 0 011-1h4zM6 7v11a2 2 0 002 2h8a2 2 0 002-2V7H6z"
                />
              </svg>
              Управление складом
            </a>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>
