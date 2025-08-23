<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend,
    ArcElement,
    PointElement,
    LineElement
  } from 'chart.js';

  // Register Chart.js components
  ChartJS.register(
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend,
    ArcElement,
    PointElement,
    LineElement
  );

  export let type: 'bar' | 'line' | 'pie' | 'doughnut' = 'bar';
  export let data: any;
  export let options: any = {};

  let canvas: HTMLCanvasElement;
  let chart: ChartJS;

  const defaultOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top' as const,
      },
      tooltip: {
        backgroundColor: 'rgba(55, 65, 81, 0.9)',
        titleColor: '#F9FAFB',
        bodyColor: '#F9FAFB',
        borderColor: '#6B7280',
        borderWidth: 1,
        cornerRadius: 8,
        displayColors: true,
      },
    },
    scales: type === 'pie' || type === 'doughnut' ? undefined : {
      y: {
        beginAtZero: true,
        grid: {
          color: 'rgba(156, 163, 175, 0.2)',
        },
        ticks: {
          color: 'rgba(107, 114, 128, 0.8)',
        },
      },
      x: {
        grid: {
          color: 'rgba(156, 163, 175, 0.2)',
        },
        ticks: {
          color: 'rgba(107, 114, 128, 0.8)',
        },
      },
    },
  };

  onMount(() => {
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    chart = new ChartJS(ctx, {
      type,
      data,
      options: { ...defaultOptions, ...options },
    });

    return () => {
      if (chart) {
        chart.destroy();
      }
    };
  });

  // Update chart when data changes
  $: if (chart && data) {
    chart.data = data;
    chart.update();
  }
</script>

<div class="chart-container">
  <canvas bind:this={canvas}></canvas>
</div>