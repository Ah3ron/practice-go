<script lang="ts">
  import { onMount } from 'svelte';
  import { auth } from '$lib/stores/auth';
  import { resourcesAPI, usersAPI } from '$lib/services/api';
  import { notifications } from '$lib/stores/notifications';
  import type { Resource, User } from '$lib/services/api';

  let resources: Resource[] = [];
  let users: User[] = [];
  let loading = true;

  $: totalResources = resources?.length || 0;
  $: totalUsers = users?.length || 0;
  $: totalQuantity = (resources || []).reduce((sum, resource) => {
    const quantity = resource?.quantity || 0;
    return sum + (typeof quantity === 'number' ? quantity : 0);
  }, 0);
  $: lowStockResources = (resources || []).filter(r => (r?.quantity || 0) < 100).length;

  onMount(async () => {
    try {
      const [resourcesResponse, usersResponse] = await Promise.all([
        resourcesAPI.getAll(),
        usersAPI.getAll()
      ]);
      
      resources = resourcesResponse.data;
      users = usersResponse.data;
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: 'Error Loading Dashboard',
        message: error.response?.data?.message || 'Failed to load dashboard data'
      });
    } finally {
      loading = false;
    }
  });

  const stats = [
    {
      title: 'Total Resources',
      value: totalResources,
      icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4',
      color: 'primary',
      change: '+12%',
      changeType: 'positive'
    },
    {
      title: 'Total Users',
      value: totalUsers,
      icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z',
      color: 'secondary',
      change: '+3%',
      changeType: 'positive'
    },
    {
      title: 'Total Quantity',
      value: typeof totalQuantity === 'number' ? totalQuantity.toLocaleString() : '0',
      icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z',
      color: 'accent',
      change: '+8%',
      changeType: 'positive'
    },
    {
      title: 'Low Stock Items',
      value: lowStockResources,
      icon: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16c-.77.833.192 2.5 1.732 2.5z',
      color: 'warning',
      change: '-5%',
      changeType: 'negative'
    }
  ];
</script>

<svelte:head>
  <title>Dashboard - Resource Manager</title>
</svelte:head>

<div class="fade-in">
  <div class="mb-8">
    <h1 class="text-3xl font-bold text-base-content mb-2">
      Welcome back, {$auth.user?.names || $auth.user?.username}!
    </h1>
    <p class="text-base-content/60">Here's what's happening with your resources today.</p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <div class="loading loading-spinner loading-lg"></div>
      <span class="ml-3 text-base-content/60">Loading dashboard...</span>
    </div>
  {:else}
    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      {#each stats as stat, index (`stat-${stat.title}-${index}`)}
        <div class="card bg-base-100 shadow-sm border border-base-200 card-hover">
          <div class="card-body p-6">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-base-content/60 text-sm font-medium">{stat.title}</p>
                <p class="text-2xl font-bold text-base-content mt-1">{stat.value}</p>
                <div class="flex items-center mt-2">
                  <span class="text-xs font-medium {stat.changeType === 'positive' ? 'text-success' : 'text-error'}">
                    {stat.change}
                  </span>
                  <span class="text-xs text-base-content/50 ml-1">vs last month</span>
                </div>
              </div>
              <div class="w-12 h-12 bg-{stat.color} bg-opacity-10 rounded-lg flex items-center justify-center">
                <svg class="w-6 h-6 text-{stat.color}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={stat.icon} />
                </svg>
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Recent Resources -->
      <div class="card bg-base-100 shadow-sm border border-base-200">
        <div class="card-body">
          <div class="flex items-center justify-between mb-4">
            <h2 class="card-title">Recent Resources</h2>
            <a href="/resources" class="btn btn-ghost btn-sm">View All</a>
          </div>
          
          {#if resources.length === 0}
            <div class="text-center py-8">
              <svg class="w-12 h-12 text-base-content/30 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
              <p class="text-base-content/60">No resources found</p>
              <a href="/resources" class="btn btn-primary btn-sm mt-2">Add Resource</a>
            </div>
          {:else}
            <div class="space-y-3">
              {#each resources.slice(0, 5) as resource, index (`resource-item-${index}`)}
                <div class="flex items-center justify-between p-3 bg-base-200 rounded-lg">
                  <div>
                    <h3 class="font-medium text-base-content">{resource?.name || 'Unknown Resource'}</h3>
                    <p class="text-sm text-base-content/60">{resource?.quantity || 0} {resource?.unit || 'units'}</p>
                  </div>
                  <div class="text-right">
                    <div class="badge {(resource?.quantity || 0) < 100 ? 'badge-warning' : 'badge-success'} badge-sm">
                      {(resource?.quantity || 0) < 100 ? 'Low Stock' : 'In Stock'}
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="card bg-base-100 shadow-sm border border-base-200">
        <div class="card-body">
          <h2 class="card-title mb-4">Quick Actions</h2>
          
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <a href="/resources" class="btn btn-primary btn-outline flex items-center justify-start">
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              Add Resource
            </a>
            
            <a href="/users" class="btn btn-secondary btn-outline flex items-center justify-start">
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
              </svg>
              Add User
            </a>
            
            <a href="/analytics" class="btn btn-accent btn-outline flex items-center justify-start">
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
              View Analytics
            </a>
            
            <a href="/resources" class="btn btn-info btn-outline flex items-center justify-start">
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
              </svg>
              View Report
            </a>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>