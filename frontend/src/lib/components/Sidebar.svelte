<script lang="ts">
  import { page } from '$app/stores';

  interface MenuItem {
    path: string;
    label: string;
    icon: string;
  }

  const menuItems: MenuItem[] = [
    {
      path: '/',
      label: 'Панель управления',
      icon: 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z M7 7l5 5l5-5'
    },
    {
      path: '/resources',
      label: 'Ресурсы',
      icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4'
    },
    {
      path: '/analytics',
      label: 'Аналитика',
      icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z'
    }
  ];

  const isActive = (itemPath: string, currentPath: string): boolean => {
    if (itemPath === '/') {
      return currentPath === '/';
    }
    return currentPath.startsWith(itemPath);
  };
</script>

<aside class="relative min-h-full w-64 bg-base-200 text-base-content">
  <!-- Header -->
  <div class="p-6">
    <div class="flex items-center space-x-2">
      <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
        <svg
          class="h-5 w-5 text-primary-content"
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
      <span class="text-lg font-bold">Менеджер ресурсов</span>
    </div>
  </div>

  <nav class="px-3">
    <ul class="menu w-full space-y-1">
      {#each menuItems as item, i (i)}
        <li>
          <a
            href={item.path}
            class="flex items-center space-x-3 rounded-lg px-4 py-3 transition-colors duration-200"
            class:bg-primary={isActive(item.path, $page.url.pathname)}
            class:text-primary-content={isActive(item.path, $page.url.pathname)}
            class:hover:bg-base-300={!isActive(item.path, $page.url.pathname)}
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={item.icon} />
            </svg>
            <span class="font-medium">{item.label}</span>
          </a>
        </li>
      {/each}
    </ul>
  </nav>

  <div class="absolute right-0 bottom-0 left-0 p-6">
    <div class="rounded-lg bg-base-300 p-4 text-center">
      <div class="mb-2 text-sm text-base-content/60">Version 1.0.0</div>
      <div class="text-xs text-base-content/40">© 2025 Менеджер ресурсов</div>
    </div>
  </div>
</aside>
