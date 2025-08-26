<script lang="ts">
  import { onMount } from 'svelte';
  import { resourcesAPI } from '$lib/services/api';
  import { notifications } from '$lib/stores/notifications';
  import type { Resource } from '$lib/services/api';
  import ResourceForm from './ResourceForm.svelte';
  import { safeFormatDistanceToNow } from '$lib/utils/date';

  let resources: Resource[] = [];
  let loading = true;
  let showCreateForm = false;
  let editingResource: Resource | null = null;
  let searchTerm = '';

  $: filteredResources = resources.filter(
    (resource) =>
      resource.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      resource.description?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      resource.unit.toLowerCase().includes(searchTerm.toLowerCase())
  );

  onMount(loadResources);

  async function loadResources() {
    loading = true;
    try {
      const response = await resourcesAPI.getAll();
      resources = response.data;
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: 'Ошибка загрузки ресурсов',
        message: error.response?.data?.message || 'Не удалось загрузить ресурсы'
      });
    } finally {
      loading = false;
    }
  }

  async function handleDelete(resource: Resource) {
    if (!confirm(`Вы уверены, что хотите удалить "${resource.name}"?`)) return;

    try {
      await resourcesAPI.delete(resource.id);

      notifications.add({
        type: 'success',
        title: 'Ресурс удалён',
        message: `"${resource.name}" был успешно удалён`
      });

      await loadResources();
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: 'Ошибка удаления',
        message: error.response?.data?.message || 'Не удалось удалить ресурс'
      });
    }
  }

  function handleEdit(resource: Resource) {
    editingResource = resource;
    showCreateForm = true;
  }

  function handleFormClose() {
    showCreateForm = false;
    editingResource = null;
  }

  function handleFormSuccess() {
    handleFormClose();
    loadResources();
  }
</script>

<svelte:head>
  <title>Ресурсы - Менеджер ресурсов</title>
</svelte:head>

<div class="fade-in">
  <div class="mb-8 flex flex-col sm:flex-row sm:items-center sm:justify-between">
    <div>
      <h1 class="mb-2 text-3xl font-bold text-base-content">Ресурсы</h1>
      <p class="text-base-content/60">Управляйте ресурсами вашей организации</p>
    </div>
    <button class="btn mt-4 btn-primary sm:mt-0" on:click={() => (showCreateForm = true)}>
      <svg class="mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Добавить ресурс
    </button>
  </div>

  <!-- Search -->
  <div class="mb-6">
    <label class="form-control w-full max-w-md">
      <input
        type="text"
        placeholder="Поиск ресурсов..."
        class="input-bordered input w-full"
        bind:value={searchTerm}
      />
    </label>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <div class="loading loading-lg loading-spinner"></div>
      <span class="ml-3 text-base-content/60">Загрузка ресурсов...</span>
    </div>
  {:else}
    <!-- Resources Grid -->
    {#if filteredResources.length === 0}
      <div class="py-12 text-center">
        <svg
          class="mx-auto mb-4 h-16 w-16 text-base-content/30"
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
        <h3 class="mb-2 text-lg font-semibold text-base-content">
          {searchTerm ? 'Ресурсы не найдены' : 'Пока нет ресурсов'}
        </h3>
        <p class="mb-4 text-base-content/60">
          {searchTerm
            ? 'Попробуйте изменить поисковые термины'
            : 'Начните с добавления вашего первого ресурса'}
        </p>
        {#if !searchTerm}
          <button class="btn btn-primary" on:click={() => (showCreateForm = true)}>
            Добавить первый ресурс
          </button>
        {/if}
      </div>
    {:else}
      <div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        {#each filteredResources as resource, index (`resource-${index}`)}
          <div class="card-hover card border border-base-200 bg-base-100 shadow-sm">
            <div class="card-body">
              <div class="mb-4 flex items-start justify-between">
                <div class="flex-1">
                  <h3 class="card-title text-lg">{resource.name}</h3>
                  {#if resource.description}
                    <p class="mt-1 text-sm text-base-content/60">{resource.description}</p>
                  {/if}
                </div>

                <div class="dropdown dropdown-end">
                  <div tabindex="0" role="button" class="btn btn-square btn-ghost btn-sm">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
                      />
                    </svg>
                  </div>
                  <ul
                    tabindex="0"
                    class="dropdown-content menu z-[1] w-40 rounded-box border border-base-200 bg-base-100 p-2 shadow"
                  >
                    <li>
                      <button on:click={() => handleEdit(resource)} class="flex items-center">
                        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                          />
                        </svg>
                        Редактировать
                      </button>
                    </li>
                    <li>
                      <button
                        on:click={() => handleDelete(resource)}
                        class="flex items-center text-error"
                      >
                        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                          />
                        </svg>
                        Удалить
                      </button>
                    </li>
                  </ul>
                </div>
              </div>

              <div class="mb-4 grid grid-cols-2 gap-4">
                <div class="rounded-lg bg-base-200 p-3 text-center">
                  <div class="text-2xl font-bold text-base-content">
                    {(resource?.quantity || 0).toLocaleString()}
                  </div>
                  <div class="text-xs font-medium text-base-content/60 uppercase">
                    {resource.unit}
                  </div>
                </div>
                <div class="rounded-lg bg-base-200 p-3 text-center">
                  <div
                    class="badge badge-sm {(resource?.quantity || 0) < 100
                      ? 'badge-warning'
                      : 'badge-success'}"
                  >
                    {(resource?.quantity || 0) < 100 ? 'Низкий остаток' : 'В наличии'}
                  </div>
                  <div class="mt-1 text-xs text-base-content/60">Статус</div>
                </div>
              </div>

              <div class="text-xs text-base-content/50">
                Обновлено {safeFormatDistanceToNow(resource.updated_at)}
              </div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  {/if}
</div>

<!-- Create/Edit Form Modal -->
{#if showCreateForm}
  <ResourceForm
    resource={editingResource}
    on:close={handleFormClose}
    on:success={handleFormSuccess}
  />
{/if}
