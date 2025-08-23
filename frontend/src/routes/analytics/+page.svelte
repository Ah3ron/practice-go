<script lang="ts">
  import { onMount } from 'svelte';
  import { resourcesAPI } from '$lib/services/api';
  import { notifications } from '$lib/stores/notifications';
  import type { Resource } from '$lib/services/api';
  import Chart from '$lib/components/Chart.svelte';

  let resources: Resource[] = [];
  let loading = true;
  let selectedTimeRange = '30d';

  onMount(loadData);

  async function loadData() {
    loading = true;
    try {
      const resourcesResponse = await resourcesAPI.getAll();
      resources = resourcesResponse.data;
    } catch (error: any) {
      notifications.add({
        type: 'error',
        title: 'Ошибка загрузки аналитики',
        message: error.response?.data?.message || 'Не удалось загрузить данные аналитики'
      });
    } finally {
      loading = false;
    }
  }

  // Advanced analytics calculations
  $: totalResources = resources?.length || 0;
  $: totalQuantity = (resources || []).reduce((sum, resource) => {
    const quantity = resource?.quantity || 0;
    return sum + (typeof quantity === 'number' ? quantity : 0);
  }, 0);
  $: averageQuantity = totalResources > 0 ? Math.round(totalQuantity / totalResources) : 0;
  $: medianQuantity = getMedianQuantity(resources);
  $: lowStockCount = (resources || []).filter((r) => (r?.quantity || 0) < 100).length;
  $: criticalStockCount = (resources || []).filter((r) => (r?.quantity || 0) < 50).length;
  $: overstockCount = (resources || []).filter((r) => (r?.quantity || 0) > 1000).length;
  $: stockValue = calculateStockValue(resources);
  $: stockTurnover = calculateStockTurnover(resources);
  $: stockCoverage = totalResources > 0 ? Math.round(((totalResources - lowStockCount) / totalResources) * 100) : 0;

  function getMedianQuantity(resources: Resource[]): number {
    if (!resources || resources.length === 0) return 0;
    const sorted = resources.map(r => r.quantity || 0).sort((a, b) => a - b);
    const mid = Math.floor(sorted.length / 2);
    return sorted.length % 2 !== 0 ? sorted[mid] : Math.round((sorted[mid - 1] + sorted[mid]) / 2);
  }

  function calculateStockValue(resources: Resource[]): number {
    // Simulated stock value calculation (quantity * assumed price per unit)
    return resources.reduce((sum, resource) => {
      const quantity = resource?.quantity || 0;
      const estimatedPrice = getEstimatedPrice(resource.unit);
      return sum + (quantity * estimatedPrice);
    }, 0);
  }

  function getEstimatedPrice(unit: string): number {
    // Estimated prices based on unit types
    const priceMap: Record<string, number> = {
      'шт': 100,
      'кг': 50,
      'л': 30,
      'м': 25,
      'упак': 150,
      'коробка': 200
    };
    return priceMap[unit] || 75; // Default price
  }

  function calculateStockTurnover(resources: Resource[]): number {
    // Simulated stock turnover calculation
    if (totalQuantity === 0) return 0;
    const averageDailyUsage = Math.round(totalQuantity * 0.02); // 2% daily usage assumption
    return averageDailyUsage > 0 ? Math.round(totalQuantity / averageDailyUsage) : 0;
  }

  // Enhanced chart data
  $: resourceQuantityData = {
    labels: resources.slice(0, 10).map((r) => r.name),
    datasets: [
      {
        label: 'Количество',
        data: resources.slice(0, 10).map((r) => r.quantity),
        backgroundColor: 'rgba(59, 130, 246, 0.8)',
        borderColor: 'rgba(59, 130, 246, 1)',
        borderWidth: 2,
        borderRadius: 4
      }
    ]
  };

  $: resourcesDistributionData = {
    labels: resources.slice(0, 8).map((r) => r.name),
    datasets: [
      {
        data: resources.slice(0, 8).map((r) => r.quantity || 0),
        backgroundColor: [
          'rgba(59, 130, 246, 0.8)',
          'rgba(16, 185, 129, 0.8)',
          'rgba(245, 158, 11, 0.8)',
          'rgba(239, 68, 68, 0.8)',
          'rgba(139, 92, 246, 0.8)',
          'rgba(236, 72, 153, 0.8)',
          'rgba(34, 197, 94, 0.8)',
          'rgba(168, 85, 247, 0.8)'
        ],
        borderWidth: 2,
        borderColor: 'rgba(255, 255, 255, 1)'
      }
    ]
  };

  $: resourceTrendsData = {
    labels: resources.slice(0, 8).map((r) => r.name),
    datasets: [
      {
        label: 'Текущий запас',
        data: resources.slice(0, 8).map((r) => r.quantity),
        backgroundColor: 'rgba(59, 130, 246, 0.2)',
        borderColor: 'rgba(59, 130, 246, 1)',
        borderWidth: 2,
        fill: true
      },
      {
        label: 'Рекомендуемый запас',
        data: resources.slice(0, 8).map((r) => Math.max(r.quantity * 1.2, 150)),
        backgroundColor: 'rgba(16, 185, 129, 0.2)',
        borderColor: 'rgba(16, 185, 129, 1)',
        borderWidth: 2,
        fill: false
      }
    ]
  };

  $: unitDistributionData = {
    labels: [...new Set(resources.map((r) => r.unit))],
    datasets: [
      {
        label: 'Количество ресурсов',
        data: [...new Set(resources.map((r) => r.unit))].map(
          (unit) => resources.filter((r) => r.unit === unit).length
        ),
        backgroundColor: [
          'rgba(59, 130, 246, 0.8)',
          'rgba(16, 185, 129, 0.8)',
          'rgba(245, 158, 11, 0.8)',
          'rgba(239, 68, 68, 0.8)',
          'rgba(139, 92, 246, 0.8)',
          'rgba(236, 72, 153, 0.8)'
        ],
        borderWidth: 2
      }
    ]
  };

  const insightCards = [
    {
      title: 'Медианное количество',
      value: medianQuantity,
      description: 'Средний показатель запасов',
      icon: 'M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18'
    },
    {
      title: 'Оборачиваемость запасов',
      value: `${stockTurnover} дней`,
      description: 'Примерное время полного оборота',
      icon: 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15'
    },
    {
      title: 'Типы единиц измерения',
      value: [...new Set(resources.map((r) => r.unit))].length,
      description: 'Различных единиц измерения',
      icon: 'M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01'
    },
    {
      title: 'Избыточные запасы',
      value: overstockCount,
      description: 'Ресурсы с количеством >1000',
      icon: 'M7 11.5V14m0-2.5v-6a1.5 1.5 0 113 0m-3 6a1.5 1.5 0 00-3 0v2a7.5 7.5 0 0015 0v-5a1.5 1.5 0 00-3 0m-6-3V11m0-5.5v-1a1.5 1.5 0 013 0v1m0 0V11m0-5.5a1.5 1.5 0 013 0v3m0 0V11'
    }
  ];
</script>

<svelte:head>
  <title>Аналитика - Менеджер ресурсов</title>
</svelte:head>

<div class="fade-in">
  <div class="mb-8 flex flex-col sm:flex-row sm:items-center sm:justify-between">
    <div>
      <h1 class="mb-2 text-3xl font-bold text-base-content">Панель аналитики</h1>
      <p class="text-base-content/60">Полный анализ управления вашими ресурсами</p>
    </div>
    <div class="mt-4 sm:mt-0">
      <select bind:value={selectedTimeRange} class="select select-bordered w-full max-w-xs">
        <option value="7d">7 дней</option>
        <option value="30d" selected>30 дней</option>
        <option value="90d">90 дней</option>
        <option value="1y">1 год</option>
      </select>
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <div class="loading loading-lg loading-spinner"></div>
      <span class="ml-3 text-base-content/60">Загрузка аналитики...</span>
    </div>
  {:else}
    <!-- Charts Grid -->
    <div class="mb-8 grid grid-cols-1 gap-8 lg:grid-cols-2">
      <!-- Resource Trends Chart -->
      <div class="card border border-base-200 bg-base-100 shadow-sm">
        <div class="card-body">
          <h2 class="mb-4 card-title">Тренды запасов ресурсов</h2>
          {#if resources.length > 0}
            <Chart type="line" data={resourceTrendsData} />
          {:else}
            <div class="py-8 text-center">
              <p class="text-base-content/60">Нет доступных ресурсов</p>
            </div>
          {/if}
        </div>
      </div>

      <!-- Resources Distribution Chart -->
      <div class="card border border-base-200 bg-base-100 shadow-sm">
        <div class="card-body">
          <h2 class="mb-4 card-title">Распределение ресурсов</h2>
          {#if resources.length > 0}
            <Chart type="pie" data={resourcesDistributionData} />
          {:else}
            <div class="py-8 text-center">
              <p class="text-base-content/60">Нет доступных ресурсов</p>
            </div>
          {/if}
        </div>
      </div>
    </div>

    <!-- Secondary Charts -->
    <div class="mb-8 grid grid-cols-1 gap-8 lg:grid-cols-2">
      <!-- Top Resources Chart -->
      <div class="card border border-base-200 bg-base-100 shadow-sm">
        <div class="card-body">
          <h2 class="mb-4 card-title">Топ ресурсов по количеству</h2>
          {#if resources.length > 0}
            <Chart type="bar" data={resourceQuantityData} />
          {:else}
            <div class="py-8 text-center">
              <p class="text-base-content/60">Нет доступных ресурсов</p>
            </div>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>
