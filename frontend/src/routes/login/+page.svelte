<script lang="ts">
  import { goto } from '$app/navigation';
  import { authAPI } from '$lib/services/api';
  import { auth } from '$lib/stores/auth';
  import { notifications } from '$lib/stores/notifications';

  let username = '';
  let password = '';
  let loading = false;

  async function handleLogin() {
    if (!username.trim() || !password.trim()) {
      notifications.add({
        type: 'warning',
        title: 'Ошибка валидации',
        message: 'Пожалуйста, заполните все поля'
      });
      return;
    }

    loading = true;

    try {
      const response = await authAPI.login({ username: username.trim(), password });

      // Update auth state
      auth.login(response.data.user, response.data.token);

      notifications.add({
        type: 'success',
        title: 'Добро пожаловать!',
        message: `Привет, ${response.data.user.names || response.data.user.username}`
      });

      // Use SvelteKit navigation with a small delay to ensure state propagation
      setTimeout(() => {
        goto('/', { replaceState: true });
      }, 100);
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: 'Ошибка входа',
        message: error.response?.data?.message || 'Неверные данные для входа'
      });
    } finally {
      loading = false;
    }
  }

  function handleKeyPress(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      handleLogin();
    }
  }
</script>

<svelte:head>
  <title>Вход - Менеджер ресурсов</title>
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
        <h2 class="text-2xl font-bold text-base-content">Добро пожаловать</h2>
        <p class="mt-2 text-base-content/60">Войдите в свою учётную запись</p>
      </div>

      <form on:submit|preventDefault={handleLogin} class="space-y-4">
        <div class="form-control">
          <label class="label" for="username">
            <span class="label-text">Имя пользователя</span>
          </label>
          <input
            id="username"
            type="text"
            bind:value={username}
            on:keypress={handleKeyPress}
            placeholder="Введите ваше имя пользователя"
            class="input-bordered input w-full"
            class:input-error={!username.trim() && loading}
            disabled={loading}
            required
          />
        </div>

        <div class="form-control">
          <label class="label" for="password">
            <span class="label-text">Пароль</span>
          </label>
          <input
            id="password"
            type="password"
            bind:value={password}
            on:keypress={handleKeyPress}
            placeholder="Введите ваш пароль"
            class="input-bordered input w-full"
            class:input-error={!password.trim() && loading}
            disabled={loading}
            required
          />
        </div>

        <button
          type="submit"
          class="btn mt-6 w-full btn-primary"
          class:loading
          disabled={loading || !username.trim() || !password.trim()}
        >
          {#if loading}
            <span class="loading loading-sm loading-spinner"></span>
            Вход...
          {:else}
            Войти
          {/if}
        </button>
      </form>

      <div class="divider">или</div>

      <p class="text-center text-sm text-base-content/60">
        Нет учётной записи?
        <a href="/register" class="link font-medium link-primary">Зарегистрироваться</a>
      </p>
    </div>
  </div>
</div>
