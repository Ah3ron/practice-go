<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { resourcesAPI, validateResource, validateResourceId } from '$lib/services/api';
  import { notifications } from '$lib/stores/notifications';
  import type { Resource } from '$lib/services/api';

  export let resource: Resource | null = null;

  const dispatch = createEventDispatcher();

  let name = resource?.name || '';
  let description = resource?.description || '';
  let unit = resource?.unit || '';
  let quantity = resource?.quantity || 0;
  let loading = false;

  $: isEditing = !!resource;

  async function handleSubmit() {
    if (!name.trim() || !unit.trim() || quantity < 0) {
      notifications.add({
        type: 'warning',
        title: 'Validation Error',
        message: 'Please fill in all required fields with valid values'
      });
      return;
    }

    // Additional validation for editing mode with detailed debugging
    if (isEditing) {
      if (!resource) {
        notifications.add({
          type: 'error',
          title: 'Invalid Resource',
          message: 'Cannot update resource: Resource object is null'
        });
        return;
      }
      
      // Use validation helper function
      if (!validateResourceId(resource.id)) {
        console.error('Resource validation failed:', {
          resource,
          id: resource.id,
          idType: typeof resource.id,
          keys: Object.keys(resource),
          isValidResource: validateResource(resource)
        });
        notifications.add({
          type: 'error',
          title: 'Invalid Resource',
          message: `Cannot update resource: ID is invalid (${resource.id}, type: ${typeof resource.id})`
        });
        return;
      }
    }

    loading = true;

    try {
      const resourceData = {
        name: name.trim(),
        description: description.trim() || undefined,
        unit: unit.trim(),
        quantity: Math.max(0, Math.floor(quantity))
      };

      if (isEditing && resource && validateResourceId(resource.id)) {
        await resourcesAPI.update(resource.id, resourceData);
        notifications.add({
          type: 'success',
          title: 'Resource Updated',
          message: `"${name}" has been updated successfully`
        });
      } else {
        await resourcesAPI.create(resourceData);
        notifications.add({
          type: 'success',
          title: 'Resource Created',
          message: `"${name}" has been created successfully`
        });
      }

      dispatch('success');
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: isEditing ? 'Update Failed' : 'Creation Failed',
        message: error.response?.data?.message || `Failed to ${isEditing ? 'update' : 'create'} resource`
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
        {isEditing ? 'Edit Resource' : 'Create New Resource'}
      </h3>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <div class="form-control col-span-1 md:col-span-2">
          <label class="label" for="name">
            <span class="label-text">Resource Name <span class="text-error">*</span></span>
          </label>
          <input
            id="name"
            type="text"
            bind:value={name}
            placeholder="Enter resource name"
            class="input input-bordered w-full"
            class:input-error={!name.trim()}
            disabled={loading}
            required
          />
        </div>

        <div class="form-control">
          <label class="label" for="unit">
            <span class="label-text">Unit <span class="text-error">*</span></span>
          </label>
          <input
            id="unit"
            type="text"
            bind:value={unit}
            placeholder="e.g., kg, pcs, liters"
            class="input input-bordered w-full"
            class:input-error={!unit.trim()}
            disabled={loading}
            required
          />
        </div>

        <div class="form-control">
          <label class="label" for="quantity">
            <span class="label-text">Quantity <span class="text-error">*</span></span>
          </label>
          <input
            id="quantity"
            type="number"
            bind:value={quantity}
            placeholder="Enter quantity"
            class="input input-bordered w-full"
            class:input-error={quantity < 0}
            min="0"
            step="1"
            disabled={loading}
            required
          />
        </div>

        <div class="form-control col-span-1 md:col-span-2">
          <label class="label" for="description">
            <span class="label-text">Description</span>
          </label>
          <textarea
            id="description"
            bind:value={description}
            placeholder="Enter resource description (optional)"
            class="textarea textarea-bordered w-full h-20"
            disabled={loading}
          ></textarea>
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
          disabled={loading || !name.trim() || !unit.trim() || quantity < 0}
        >
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
            {isEditing ? 'Updating...' : 'Creating...'}
          {:else}
            {isEditing ? 'Update Resource' : 'Create Resource'}
          {/if}
        </button>
      </div>
    </form>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button on:click={handleClose}>close</button>
  </form>
</div>