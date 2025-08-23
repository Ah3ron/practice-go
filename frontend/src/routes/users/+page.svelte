<script lang="ts">
  import { onMount } from 'svelte';
  import { usersAPI } from '$lib/services/api';
  import { notifications } from '$lib/stores/notifications';
  import type { User } from '$lib/services/api';
  import UserForm from './UserForm.svelte';
  import { safeFormatDistanceToNow } from '$lib/utils/date';

  let users: User[] = [];
  let loading = true;
  let showCreateForm = false;
  let editingUser: User | null = null;
  let searchTerm = '';

  $: filteredUsers = users.filter(user =>
    user.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
    user.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
    user.names?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  onMount(loadUsers);

  async function loadUsers() {
    loading = true;
    try {
      const response = await usersAPI.getAll();
      users = response.data;
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: 'Error Loading Users',
        message: error.response?.data?.message || 'Failed to load users'
      });
    } finally {
      loading = false;
    }
  }

  async function handleDelete(user: User) {
    if (!confirm(`Are you sure you want to delete user "${user.username}"?`)) return;

    try {
      await usersAPI.delete(user.id);
      
      notifications.add({
        type: 'success',
        title: 'User Deleted',
        message: `User "${user.username}" has been deleted successfully`
      });
      
      await loadUsers();
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: 'Delete Failed',
        message: error.response?.data?.message || 'Failed to delete user'
      });
    }
  }

  function handleEdit(user: User) {
    editingUser = user;
    showCreateForm = true;
  }

  function handleFormClose() {
    showCreateForm = false;
    editingUser = null;
  }

  function handleFormSuccess() {
    handleFormClose();
    loadUsers();
  }
</script>

<svelte:head>
  <title>Users - Resource Manager</title>
</svelte:head>

<div class="fade-in">
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-8">
    <div>
      <h1 class="text-3xl font-bold text-base-content mb-2">Users</h1>
      <p class="text-base-content/60">Manage system users and their access</p>
    </div>
    <button 
      class="btn btn-primary mt-4 sm:mt-0" 
      on:click={() => showCreateForm = true}
    >
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
      </svg>
      Add User
    </button>
  </div>

  <!-- Search -->
  <div class="mb-6">
    <div class="form-control w-full max-w-md">
      <input 
        type="text" 
        placeholder="Search users..." 
        class="input input-bordered w-full"
        bind:value={searchTerm}
      />
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <div class="loading loading-spinner loading-lg"></div>
      <span class="ml-3 text-base-content/60">Loading users...</span>
    </div>
  {:else}
    <!-- Users Table -->
    {#if filteredUsers.length === 0}
      <div class="text-center py-12">
        <svg class="w-16 h-16 text-base-content/30 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z" />
        </svg>
        <h3 class="text-lg font-semibold text-base-content mb-2">
          {searchTerm ? 'No users found' : 'No users yet'}
        </h3>
        <p class="text-base-content/60 mb-4">
          {searchTerm ? 'Try adjusting your search terms' : 'Get started by adding your first user'}
        </p>
        {#if !searchTerm}
          <button class="btn btn-primary" on:click={() => showCreateForm = true}>
            Add First User
          </button>
        {/if}
      </div>
    {:else}
      <div class="card bg-base-100 shadow-sm border border-base-200">
        <div class="card-body p-0">
          <div class="overflow-x-auto">
            <table class="table table-zebra">
              <thead>
                <tr class="border-b border-base-200">
                  <th class="bg-base-200">User</th>
                  <th class="bg-base-200">Email</th>
                  <th class="bg-base-200">Created</th>
                  <th class="bg-base-200">Updated</th>
                  <th class="bg-base-200">Actions</th>
                </tr>
              </thead>
              <tbody>
                {#each filteredUsers as user, index (`user-${index}`)}
                  <tr>
                    <td>
                      <div class="flex items-center space-x-3">
                        <div class="avatar">
                          <div class="w-10 h-10 bg-primary text-primary-content rounded-full flex items-center justify-center">
                            <span class="font-medium text-sm">
                              {user.names?.charAt(0) || user.username.charAt(0).toUpperCase()}
                            </span>
                          </div>
                        </div>
                        <div>
                          <div class="font-medium">{user.names || user.username}</div>
                          <div class="text-sm text-base-content/60">@{user.username}</div>
                        </div>
                      </div>
                    </td>
                    <td>{user.email}</td>
                    <td class="text-sm">
                      {safeFormatDistanceToNow(user.created_at)}
                    </td>
                    <td class="text-sm">
                      {safeFormatDistanceToNow(user.updated_at)}
                    </td>
                    <td>
                      <div class="flex items-center space-x-2">
                        <button 
                          class="btn btn-ghost btn-sm"
                          on:click={() => handleEdit(user)}
                        >
                          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                          </svg>
                        </button>
                        <button 
                          class="btn btn-ghost btn-sm text-error"
                          on:click={() => handleDelete(user)}
                        >
                          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                          </svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    {/if}
  {/if}
</div>

<!-- Create/Edit Form Modal -->
{#if showCreateForm}
  <UserForm 
    user={editingUser}
    on:close={handleFormClose}
    on:success={handleFormSuccess}
  />
{/if}