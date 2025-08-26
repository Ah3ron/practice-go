<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { page } from '$app/stores';

  $: pageTitle = getPageTitle($page.url.pathname);

  function getPageTitle(pathname: string): string {
    const routes: Record<string, string> = {
      '/': 'Панель управления',
      '/resources': 'Ресурсы',
      '/analytics': 'Аналитика',
      '/profile': 'Профиль'
    };

    return routes[pathname] || 'Приложение';
  }

  function handleLogout() {
    auth.logout();
    window.location.href = '/login';
  }
</script>

<div class="navbar border-b border-base-200 bg-base-100 px-6 shadow-sm">
  <div class="flex-1">
    <label for="drawer-toggle" class="btn btn-ghost lg:hidden">
      <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M4 6h16M4 12h16M4 18h16"
        />
      </svg>
    </label>

    <div class="ml-2">
      <h1 class="text-xl font-semibold text-base-content">{pageTitle}</h1>
      <div class="flex items-center text-sm text-base-content/60">
        <svg class="mr-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"
          />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m7 7l5 5l5-5" />
        </svg>
        Система управления ресурсами
      </div>
    </div>
  </div>

  <div class="flex-none">
    {#if $auth.isAuthenticated}
      <div class="dropdown dropdown-end">
        <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
          <div
            class="flex w-10 items-center justify-center rounded-full bg-primary text-primary-content"
          >
            <span class="text-sm font-medium">
              {$auth.user?.names?.charAt(0) || $auth.user?.username?.charAt(0) || 'U'}
            </span>
          </div>
        </div>

        <ul
          tabindex="0"
          class="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow border border-base-200"
        >
          <li class="menu-title">
            <span class="text-base-content/60">
              {$auth.user?.names || $auth.user?.username}
            </span>
          </li>
          <li><div class="divider my-1"></div></li>
          <li>
            <a href="/profile" class="flex items-center">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                />
              </svg>
              Профиль
            </a>
          </li>
          <li><div class="divider my-1"></div></li>
          <li>
            <button on:click={handleLogout} class="flex items-center text-error">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                />
              </svg>
              Выйти
            </button>
          </li>
        </ul>
      </div>
    {/if}
  </div>
</div>
