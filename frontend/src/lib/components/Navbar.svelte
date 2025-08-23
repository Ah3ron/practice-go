<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { page } from '$app/stores';

  $: pageTitle = getPageTitle($page.url.pathname);

  function getPageTitle(pathname: string): string {
    const routes: Record<string, string> = {
      '/': 'Dashboard',
      '/resources': 'Resources',
      '/users': 'Users',
      '/analytics': 'Analytics',
      '/profile': 'Profile',
    };
    
    return routes[pathname] || 'Application';
  }

  function handleLogout() {
    auth.logout();
    window.location.href = '/login';
  }
</script>

<div class="navbar bg-base-100 shadow-sm border-b border-base-200 px-6">
  <div class="flex-1">
    <label for="drawer-toggle" class="btn btn-ghost lg:hidden">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
      </svg>
    </label>
    
    <div class="ml-2">
      <h1 class="text-xl font-semibold text-base-content">{pageTitle}</h1>
      <div class="text-sm text-base-content/60 flex items-center">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m7 7l5 5l5-5" />
        </svg>
        Resource Management System
      </div>
    </div>
  </div>
  
  <div class="flex-none">
    {#if $auth.isAuthenticated}
      <div class="dropdown dropdown-end">
        <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
          <div class="w-10 rounded-full bg-primary text-primary-content flex items-center justify-center">
            <span class="text-sm font-medium">
              {$auth.user?.names?.charAt(0) || $auth.user?.username?.charAt(0) || 'U'}
            </span>
          </div>
        </div>
        
        <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow-lg bg-base-100 rounded-box w-52 border border-base-200">
          <li class="menu-title">
            <span class="text-base-content/60">
              {$auth.user?.names || $auth.user?.username}
            </span>
          </li>
          <div class="divider my-1"></div>
          <li><a href="/profile" class="flex items-center">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            Profile
          </a></li>
          <div class="divider my-1"></div>
          <li>
            <button on:click={handleLogout} class="text-error flex items-center">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              Sign out
            </button>
          </li>
        </ul>
      </div>
    {/if}
  </div>
</div>