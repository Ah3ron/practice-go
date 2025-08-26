<script lang="ts">
  import { authAPI } from '$lib/services/api';
  import { notifications } from '$lib/stores/notifications';

  let username = '';
  let email = '';
  let password = '';
  let names = '';
  let loading = false;

  async function handleRegister() {
    if (!username.trim() || !email.trim() || !password.trim()) {
      notifications.add({
        type: 'warning',
        title: 'Ошибка валидации',
        message: 'Пожалуйста, заполните все обязательные поля'
      });
      return;
    }

    if (password.length < 6) {
      notifications.add({
        type: 'warning',
        title: 'Ошибка валидации',
        message: 'Пароль должен содержать не менее 6 символов'
      });
      return;
    }

    loading = true;

    try {
      await authAPI.register({
        username: username.trim(),
        email: email.trim(),
        password,
        names: names.trim() || undefined
      });

      notifications.add({
        type: 'success',
        title: 'Аккаунт создан!',
        message: 'Теперь вы можете войти в систему'
      });

      window.location.href = '/login';
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: 'Ошибка регистрации',
        message: error.response?.data?.message || 'Не удалось создать аккаунт'
      });
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Регистрация - Менеджер ресурсов</title>
</svelte:head>

<div
  class="flex min-h-screen items-center justify-center bg-gradient-to-br from-primary/10 to-secondary/10 p-4"
>
  <div class="card w-full max-w-md bg-base-100 shadow-2xl">
    <div class="card-body">
      <div class="mb-8 text-center">
        <div
          class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-primary"
        >
          <svg
            class="h-8 w-8 text-primary-content"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
            />
          </svg>
        </div>
        <h2 class="text-2xl font-bold text-base-content">Создать аккаунт</h2>
        <p class="mt-2 text-base-content/60">Присоединяйтесь к нам сегодня</p>
      </div>

      <form on:submit|preventDefault={handleRegister} class="space-y-4">
        <label class="form-control w-full">
          <div class="label">
            <span class="label-text">Имя пользователя <span class="text-error">*</span></span>
          </div>
          <input
            id="username"
            type="text"
            bind:value={username}
            placeholder="Введите имя пользователя"
            class="input input-bordered w-full"
            disabled={loading}
            required
          />
        </label>

        <label class="form-control w-full">
          <div class="label">
            <span class="label-text">Email <span class="text-error">*</span></span>
          </div>
          <input
            id="email"
            type="email"
            bind:value={email}
            placeholder="Введите адрес электронной почты"
            class="input input-bordered w-full"
            disabled={loading}
            required
          />
        </label>

        <label class="form-control w-full">
          <div class="label">
            <span class="label-text">Полное имя</span>
          </div>
          <input
            id="names"
            type="text"
            bind:value={names}
            placeholder="Введите ваше полное имя"
            class="input input-bordered w-full"
            disabled={loading}
          />
        </label>

        <label class="form-control w-full">
          <div class="label">
            <span class="label-text">Пароль <span class="text-error">*</span></span>
          </div>
          <input
            id="password"
            type="password"
            bind:value={password}
            placeholder="Введите пароль (мин. 6 символов)"
            class="input input-bordered w-full"
            disabled={loading}
            required
          />
        </label>

        <button
          type="submit"
          class="btn btn-primary w-full mt-6"
          class:loading
          disabled={loading || !username.trim() || !email.trim() || !password.trim()}
        >
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
            Создание аккаунта...
          {:else}
            Создать аккаунт
          {/if}
        </button>
      </form>

      <div class="divider">или</div>

      <p class="text-center text-sm text-base-content/60">
        Уже есть аккаунт?
        <a href="/login" class="link font-medium link-primary">Войти</a>
      </p>
    </div>
  </div>
</div>
