<script lang="ts">
  import { onMount } from 'svelte';
  import { resourcesAPI, usersAPI } from '$lib/services/api';
  import { notifications } from '$lib/stores/notifications';
  import type { Resource, User } from '$lib/services/api';
  import Chart from '$lib/components/Chart.svelte';

  let resources: Resource[] = [];
  let users: User[] = [];
  let loading = true;

  onMount(loadData);

  async function loadData() {
    loading = true;
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
        title: 'Error Loading Analytics',
        message: error.response?.data?.message || 'Failed to load analytics data'
      });
    } finally {
      loading = false;
    }
  }

  // Analytics calculations with safe handling
  $: totalResources = resources?.length || 0;
  $: totalQuantity = (resources || []).reduce((sum, resource) => {
    const quantity = resource?.quantity || 0;
    return sum + (typeof quantity === 'number' ? quantity : 0);
  }, 0);
  $: averageQuantity = totalResources > 0 ? Math.round(totalQuantity / totalResources) : 0;
  $: lowStockCount = (resources || []).filter(r => (r?.quantity || 0) < 100).length;
  $: highStockCount = (resources || []).filter(r => (r?.quantity || 0) >= 500).length;

  // Chart data
  $: resourceQuantityData = {
    labels: resources.slice(0, 10).map(r => r.name),
    datasets: [{
      label: 'Quantity',
      data: resources.slice(0, 10).map(r => r.quantity),
      backgroundColor: 'rgba(59, 130, 246, 0.8)',
      borderColor: 'rgba(59, 130, 246, 1)',
      borderWidth: 1
    }]
  };

  $: resourceUnitsData = {
    labels: [...new Set(resources.map(r => r.unit))],
    datasets: [{
      label: 'Count by Unit',
      data: [...new Set(resources.map(r => r.unit))].map(unit => 
        resources.filter(r => r.unit === unit).length
      ),
      backgroundColor: [
        'rgba(59, 130, 246, 0.8)',
        'rgba(16, 185, 129, 0.8)',
        'rgba(245, 158, 11, 0.8)',
        'rgba(239, 68, 68, 0.8)',
        'rgba(139, 92, 246, 0.8)',
        'rgba(236, 72, 153, 0.8)'
      ],
      borderWidth: 0
    }]
  };

  $: stockStatusData = {
    labels: ['Low Stock (<100)', 'Normal Stock (100-499)', 'High Stock (≥500)'],
    datasets: [{
      data: [
        lowStockCount,
        resources.filter(r => r.quantity >= 100 && r.quantity < 500).length,
        highStockCount
      ],
      backgroundColor: [
        'rgba(245, 158, 11, 0.8)',
        'rgba(16, 185, 129, 0.8)',
        'rgba(59, 130, 246, 0.8)'
      ],
      borderWidth: 0
    }]
  };

  const kpiCards = [
    {
      title: 'Total Resources',
      value: totalResources,
      change: '+12%',
      changeType: 'positive',
      icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4'
    },
    {
      title: 'Total Quantity',
      value: typeof totalQuantity === 'number' ? totalQuantity.toLocaleString() : '0',
      change: '+8%',
      changeType: 'positive',
      icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z'
    },
    {
      title: 'Average Quantity',
      value: typeof averageQuantity === 'number' ? averageQuantity.toLocaleString() : '0',
      change: '+5%',
      changeType: 'positive',
      icon: 'M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z'
    },
    {
      title: 'Low Stock Items',
      value: lowStockCount,
      change: '-15%',
      changeType: 'negative',
      icon: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16c-.77.833.192 2.5 1.732 2.5z'
    }
  ];
</script>

<svelte:head>
  <title>Analytics - Resource Manager</title>
</svelte:head>

<div class="fade-in">
  <div class="mb-8">
    <h1 class="text-3xl font-bold text-base-content mb-2">Analytics Dashboard</h1>
    <p class="text-base-content/60">Comprehensive insights into your resource management</p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <div class="loading loading-spinner loading-lg"></div>
      <span class="ml-3 text-base-content/60">Loading analytics...</span>
    </div>
  {:else}
    <!-- KPI Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      {#each kpiCards as kpi, index (`kpi-${index}`)}
        <div class="card bg-base-100 shadow-sm border border-base-200 card-hover">
          <div class="card-body p-6">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-base-content/60 text-sm font-medium">{kpi.title}</p>
                <p class="text-2xl font-bold text-base-content mt-1">{kpi.value}</p>
                <div class="flex items-center mt-2">
                  <span class="text-xs font-medium {kpi.changeType === 'positive' ? 'text-success' : 'text-error'}">
                    {kpi.change}
                  </span>
                  <span class="text-xs text-base-content/50 ml-1">vs last month</span>
                </div>
              </div>
              <div class="w-12 h-12 bg-primary bg-opacity-10 rounded-lg flex items-center justify-center">
                <svg class="w-6 h-6 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={kpi.icon} />
                </svg>
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>

    <!-- Charts Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
      <!-- Resource Quantities Chart -->
      <div class="card bg-base-100 shadow-sm border border-base-200">
        <div class="card-body">
          <h2 class="card-title mb-4">Top Resources by Quantity</h2>
          {#if resources.length > 0}
            <Chart type="bar" data={resourceQuantityData} />
          {:else}
            <div class="text-center py-8">
              <p class="text-base-content/60">No resources available</p>
            </div>
          {/if}
        </div>
      </div>

      <!-- Stock Status Chart -->
      <div class="card bg-base-100 shadow-sm border border-base-200">
        <div class="card-body">
          <h2 class="card-title mb-4">Stock Status Distribution</h2>
          {#if resources.length > 0}
            <Chart type="doughnut" data={stockStatusData} />
          {:else}
            <div class="text-center py-8">
              <p class="text-base-content/60">No resources available</p>
            </div>
          {/if}
        </div>
      </div>
    </div>

    <!-- Resource Units Chart -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <div class="card bg-base-100 shadow-sm border border-base-200">
        <div class="card-body">
          <h2 class="card-title mb-4">Resources by Unit Type</h2>
          {#if resources.length > 0}
            <Chart type="pie" data={resourceUnitsData} />
          {:else}
            <div class="text-center py-8">
              <p class="text-base-content/60">No resources available</p>
            </div>
          {/if}
        </div>
      </div>

      <!-- Summary Statistics -->
      <div class="card bg-base-100 shadow-sm border border-base-200">
        <div class="card-body">
          <h2 class="card-title mb-4">Summary Statistics</h2>
          
          <div class="space-y-4">
            <div class="flex justify-between items-center p-4 bg-base-200 rounded-lg">
              <div>
                <div class="font-medium text-base-content">Total Users</div>
                <div class="text-sm text-base-content/60">Registered users</div>
              </div>
              <div class="text-2xl font-bold text-primary">{users.length}</div>
            </div>

            <div class="flex justify-between items-center p-4 bg-base-200 rounded-lg">
              <div>
                <div class="font-medium text-base-content">Unique Units</div>
                <div class="text-sm text-base-content/60">Different measurement units</div>
              </div>
              <div class="text-2xl font-bold text-secondary">{[...new Set(resources.map(r => r.unit))].length}</div>
            </div>

            <div class="flex justify-between items-center p-4 bg-base-200 rounded-lg">
              <div>
                <div class="font-medium text-base-content">High Stock Items</div>
                <div class="text-sm text-base-content/60">Items with ≥500 quantity</div>
              </div>
              <div class="text-2xl font-bold text-success">{highStockCount}</div>
            </div>

            <div class="flex justify-between items-center p-4 bg-base-200 rounded-lg">
              <div>
                <div class="font-medium text-base-content">Stock Coverage</div>
                <div class="text-sm text-base-content/60">Percentage well-stocked</div>
              </div>
              <div class="text-2xl font-bold text-accent">
                {totalResources > 0 ? Math.round(((totalResources - lowStockCount) / totalResources) * 100) : 0}%
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>