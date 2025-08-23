<script lang="ts">
  import { page } from '$app/stores';

  $: pathname = $page.url.pathname;

  interface MenuItem {
    path: string;
    label: string;
    icon: string;
  }

  const menuItems: MenuItem[] = [
    {
      path: '/',
      label: 'Dashboard',
      icon: 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z M7 7l5 5l5-5'
    },
    {
      path: '/resources',
      label: 'Resources',
      icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4'
    },
    {
      path: '/users',
      label: 'Users',
      icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z'
    },
    {
      path: '/analytics',
      label: 'Analytics',
      icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z'
    },
  ];

  function isActive(itemPath: string): boolean {
    if (itemPath === '/') {
      return pathname === '/';
    }
    return pathname.startsWith(itemPath);
  }
</script>

<aside class="w-64 min-h-full bg-base-200 text-base-content">
  <div class="p-6">
    <div class="flex items-center space-x-2">
      <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
        <svg class="w-5 h-5 text-primary-content" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
      </div>
      <span class="text-lg font-bold">Resource Manager</span>
    </div>
  </div>

  <nav class="px-3">
    <ul class="menu menu-vertical w-full space-y-1">
      {#each menuItems as item, index (`menu-${index}`)}
        <li>
          <a 
            href={item.path}
            class="flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors duration-200 {isActive(item.path) ? 'bg-primary text-primary-content' : 'hover:bg-base-300'}"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={item.icon} />
            </svg>
            <span class="font-medium">{item.label}</span>
          </a>
        </li>
      {/each}
    </ul>
  </nav>

  <div class="absolute bottom-0 left-0 right-0 p-6">
    <div class="bg-base-300 rounded-lg p-4 text-center">
      <div class="text-sm text-base-content/60 mb-2">Version 1.0.0</div>
      <div class="text-xs text-base-content/40">© 2024 Resource Manager</div>
    </div>
  </div>
</aside>