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

  $: filteredResources = resources.filter(resource =>
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
        title: 'Error Loading Resources',
        message: error.response?.data?.message || 'Failed to load resources'
      });
    } finally {
      loading = false;
    }
  }

  async function handleDelete(resource: Resource) {
    if (!confirm(`Are you sure you want to delete "${resource.name}"?`)) return;

    try {
      await resourcesAPI.delete(resource.id);
      
      notifications.add({
        type: 'success',
        title: 'Resource Deleted',
        message: `"${resource.name}" has been deleted successfully`
      });
      
      await loadResources();
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: 'Delete Failed',
        message: error.response?.data?.message || 'Failed to delete resource'
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
  <title>Resources - Resource Manager</title>
</svelte:head>

<div class="fade-in">
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-8">
    <div>
      <h1 class="text-3xl font-bold text-base-content mb-2">Resources</h1>
      <p class="text-base-content/60">Manage your organization's resources</p>
    </div>
    <button 
      class="btn btn-primary mt-4 sm:mt-0" 
      on:click={() => showCreateForm = true}
    >
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Add Resource
    </button>
  </div>

  <!-- Search -->
  <div class="mb-6">
    <div class="form-control w-full max-w-md">
      <input 
        type="text" 
        placeholder="Search resources..." 
        class="input input-bordered w-full"
        bind:value={searchTerm}
      />
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <div class="loading loading-spinner loading-lg"></div>
      <span class="ml-3 text-base-content/60">Loading resources...</span>
    </div>
  {:else}
    <!-- Resources Grid -->
    {#if filteredResources.length === 0}
      <div class="text-center py-12">
        <svg class="w-16 h-16 text-base-content/30 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
        </svg>
        <h3 class="text-lg font-semibold text-base-content mb-2">
          {searchTerm ? 'No resources found' : 'No resources yet'}
        </h3>
        <p class="text-base-content/60 mb-4">
          {searchTerm ? 'Try adjusting your search terms' : 'Get started by adding your first resource'}
        </p>
        {#if !searchTerm}
          <button class="btn btn-primary" on:click={() => showCreateForm = true}>
            Add First Resource
          </button>
        {/if}
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each filteredResources as resource, index (`resource-${index}`)}
          <div class="card bg-base-100 shadow-sm border border-base-200 card-hover">
            <div class="card-body">
              <div class="flex items-start justify-between mb-4">
                <div class="flex-1">
                  <h3 class="card-title text-lg">{resource.name}</h3>
                  {#if resource.description}
                    <p class="text-base-content/60 text-sm mt-1">{resource.description}</p>
                  {/if}
                </div>
                
                <div class="dropdown dropdown-end">
                  <div tabindex="0" role="button" class="btn btn-ghost btn-sm btn-square">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
                    </svg>
                  </div>
                  <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow-lg bg-base-100 rounded-box w-40 border border-base-200">
                    <li>
                      <button on:click={() => handleEdit(resource)} class="flex items-center">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                        Edit
                      </button>
                    </li>
                    <li>
                      <button on:click={() => handleDelete(resource)} class="text-error flex items-center">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                        Delete
                      </button>
                    </li>
                  </ul>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-4 mb-4">
                <div class="text-center p-3 bg-base-200 rounded-lg">
                  <div class="text-2xl font-bold text-base-content">{(resource?.quantity || 0).toLocaleString()}</div>
                  <div class="text-xs text-base-content/60 uppercase font-medium">{resource.unit}</div>
                </div>
                <div class="text-center p-3 bg-base-200 rounded-lg">
                  <div class="badge {(resource?.quantity || 0) < 100 ? 'badge-warning' : 'badge-success'} badge-sm">
                    {(resource?.quantity || 0) < 100 ? 'Low Stock' : 'In Stock'}
                  </div>
                  <div class="text-xs text-base-content/60 mt-1">Status</div>
                </div>
              </div>

              <div class="text-xs text-base-content/50">
                Updated {safeFormatDistanceToNow(resource.updated_at)}
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