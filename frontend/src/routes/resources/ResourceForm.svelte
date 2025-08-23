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
        title: 'Ошибка валидации',
        message: 'Пожалуйста, заполните все обязательные поля корректными значениями'
      });
      return;
    }

    // Additional validation for editing mode with detailed debugging
    if (isEditing) {
      if (!resource) {
        notifications.add({
          type: 'error',
          title: 'Неверный ресурс',
          message: 'Не удалось обновить ресурс: Объект ресурса нулевой'
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
          title: 'Неверный ресурс',
          message: `Не удалось обновить ресурс: Неверный ID (${resource.id}, тип: ${typeof resource.id})`
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
          title: 'Ресурс обновлён',
          message: `"${name}" был успешно обновлён`
        });
      } else {
        await resourcesAPI.create(resourceData);
        notifications.add({
          type: 'success',
          title: 'Ресурс создан',
          message: `"${name}" был успешно создан`
        });
      }

      dispatch('success');
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: isEditing ? 'Ошибка обновления' : 'Ошибка создания',
        message:
          error.response?.data?.message || `Не удалось ${isEditing ? 'обновить' : 'создать'} ресурс`
      });
    } finally {
      loading = false;
    }
  }

  function handleClose() {
    dispatch('close');
  }
</script>

<div class="modal-open modal">
  <div class="modal-box w-11/12 max-w-2xl">
    <form on:submit|preventDefault={handleSubmit}>
      <h3 class="mb-6 text-lg font-bold">
        {isEditing ? 'Редактировать ресурс' : 'Создать новый ресурс'}
      </h3>

      <div class="mb-6 grid grid-cols-1 gap-4 md:grid-cols-2">
        <div class="form-control col-span-1 md:col-span-2">
          <label class="label" for="name">
            <span class="label-text">Название ресурса <span class="text-error">*</span></span>
          </label>
          <input
            id="name"
            type="text"
            bind:value={name}
            placeholder="Введите название ресурса"
            class="input-bordered input w-full"
            class:input-error={!name.trim()}
            disabled={loading}
            required
          />
        </div>

        <div class="form-control">
          <label class="label" for="unit">
            <span class="label-text">Единица измерения <span class="text-error">*</span></span>
          </label>
          <input
            id="unit"
            type="text"
            bind:value={unit}
            placeholder="например: кг, шт, литры"
            class="input-bordered input w-full"
            class:input-error={!unit.trim()}
            disabled={loading}
            required
          />
        </div>

        <div class="form-control">
          <label class="label" for="quantity">
            <span class="label-text">Количество <span class="text-error">*</span></span>
          </label>
          <input
            id="quantity"
            type="number"
            bind:value={quantity}
            placeholder="Введите количество"
            class="input-bordered input w-full"
            class:input-error={quantity < 0}
            min="0"
            step="1"
            disabled={loading}
            required
          />
        </div>

        <div class="form-control col-span-1 md:col-span-2">
          <label class="label" for="description">
            <span class="label-text">Описание</span>
          </label>
          <textarea
            id="description"
            bind:value={description}
            placeholder="Введите описание ресурса (необязательно)"
            class="textarea-bordered textarea h-20 w-full"
            disabled={loading}
          ></textarea>
        </div>
      </div>

      <div class="modal-action">
        <button type="button" class="btn btn-ghost" on:click={handleClose} disabled={loading}>
          Отмена
        </button>
        <button
          type="submit"
          class="btn btn-primary"
          class:loading
          disabled={loading || !name.trim() || !unit.trim() || quantity < 0}
        >
          {#if loading}
            <span class="loading loading-sm loading-spinner"></span>
            {isEditing ? 'Обновление...' : 'Создание...'}
          {:else}
            {isEditing ? 'Обновить ресурс' : 'Создать ресурс'}
          {/if}
        </button>
      </div>
    </form>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button on:click={handleClose}>close</button>
  </form>
</div>
