<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { usersAPI } from '$lib/services/api';
  import { notifications } from '$lib/stores/notifications';
  import type { User } from '$lib/services/api';

  export let user: User | null = null;

  const dispatch = createEventDispatcher();

  let username = user?.username || '';
  let email = user?.email || '';
  let names = user?.names || '';
  let password = '';
  let loading = false;

  $: isEditing = !!user;

  async function handleSubmit() {
    if (!username.trim() || !email.trim() || (!isEditing && !password.trim())) {
      notifications.add({
        type: 'warning',
        title: 'Validation Error',
        message: 'Please fill in all required fields'
      });
      return;
    }

    if (!isEditing && password.length < 6) {
      notifications.add({
        type: 'warning',
        title: 'Validation Error',
        message: 'Password must be at least 6 characters long'
      });
      return;
    }

    loading = true;

    try {
      const userData = {
        username: username.trim(),
        email: email.trim(),
        names: names.trim() || undefined,
        ...((!isEditing || password.trim()) && { password: password.trim() })
      };

      if (isEditing && user) {
        await usersAPI.update(user.id, userData);
        notifications.add({
          type: 'success',
          title: 'User Updated',
          message: `User "${username}" has been updated successfully`
        });
      } else {
        await usersAPI.create(userData);
        notifications.add({
          type: 'success',
          title: 'User Created',
          message: `User "${username}" has been created successfully`
        });
      }

      dispatch('success');
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: isEditing ? 'Update Failed' : 'Creation Failed',
        message: error.response?.data?.message || `Failed to ${isEditing ? 'update' : 'create'} user`
      });
    } finally {
      loading = false;
    }
  }

  function handleClose() {
    dispatch('close');
  }
</script>

<div class="modal modal-open">
  <div class="modal-box w-11/12 max-w-2xl">
    <form on:submit|preventDefault={handleSubmit}>
      <h3 class="font-bold text-lg mb-6">
        {isEditing ? 'Edit User' : 'Create New User'}
      </h3>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <div class="form-control">
          <label class="label" for="username">
            <span class="label-text">Username <span class="text-error">*</span></span>
          </label>
          <input
            id="username"
            type="text"
            bind:value={username}
            placeholder="Enter username"
            class="input input-bordered w-full"
            class:input-error={!username.trim()}
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
            class="input input-bordered w-full"
            class:input-error={!email.trim()}
            disabled={loading}
            required
          />
        </div>

        <div class="form-control col-span-1 md:col-span-2">
          <label class="label" for="names">
            <span class="label-text">Full Name</span>
          </label>
          <input
            id="names"
            type="text"
            bind:value={names}
            placeholder="Enter full name (optional)"
            class="input input-bordered w-full"
            disabled={loading}
          />
        </div>

        <div class="form-control col-span-1 md:col-span-2">
          <label class="label" for="password">
            <span class="label-text">
              Password 
              {#if !isEditing}
                <span class="text-error">*</span>
              {:else}
                <span class="text-base-content/50">(leave empty to keep current)</span>
              {/if}
            </span>
          </label>
          <input
            id="password"
            type="password"
            bind:value={password}
            placeholder={isEditing ? "Enter new password (optional)" : "Enter password (min 6 characters)"}
            class="input input-bordered w-full"
            class:input-error={!isEditing && password.length > 0 && password.length < 6}
            disabled={loading}
            required={!isEditing}
          />
        </div>
      </div>

      <div class="modal-action">
        <button
          type="button"
          class="btn btn-ghost"
          on:click={handleClose}
          disabled={loading}
        >
          Cancel
        </button>
        <button
          type="submit"
          class="btn btn-primary"
          class:loading
          disabled={loading || !username.trim() || !email.trim() || (!isEditing && !password.trim())}
        >
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
            {isEditing ? 'Updating...' : 'Creating...'}
          {:else}
            {isEditing ? 'Update User' : 'Create User'}
          {/if}
        </button>
      </div>
    </form>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button on:click={handleClose}>close</button>
  </form>
</div>