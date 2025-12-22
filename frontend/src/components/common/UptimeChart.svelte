<script>
    import { onMount, onDestroy } from 'svelte';
    import ApexCharts from 'apexcharts';

    const { data, name, height } = $props();

    let chartEl;
    let chart;

    onMount(() => {
        const options = {
            chart: {
                type: 'area',
                height,
                zoom: {
                    enabled: false,
                },
                toolbar: {
                    show: false,
                },
            },

            grid: {
                show: false,
            },

            dataLabels: {
                enabled: true,
                formatter: val => `${val} ms`,
                offsetY: -10,
                style: {
                    fontSize: '10px',
                    fontWeight: '600',
                    colors: ['#0088ee'],
                },
                background: {
                    enabled: false,
                },
            },

            stroke: {
                curve: 'straight',
                width: 2,
            },

            fill: {
                type: 'gradient',
                gradient: {
                    shadeIntensity: 1,
                    opacityFrom: 0.4,
                    opacityTo: 0.7,
                    stops: [0, 90, 100],
                },
            },

            xaxis: {
                axisBorder: { show: true },
                axisTicks: { show: true },
                tickAmount: 10,
            },

            yaxis: {
                show: false,
                min: 0,
                max: 5000, // 5 ثانیه = 5000 میلی‌ثانیه
                tickAmount: 5,
                axisBorder: { show: false },
                axisTicks: { show: false },
                labels: {
                    formatter: val => (val >= 1000 ? `${(val / 1000).toFixed(1)} s` : `${val} ms`),
                },
            },

            legend: {
                show: false,
            },
        };

        chart = new ApexCharts(chartEl, {
            ...options,
            series: [
                {
                    name,
                    data,
                },
            ],
        });
        chart.render();
    });

    onDestroy(() => {
        chart?.destroy();
    });
</script>

<div
    class="w-full border flex flex-col justify-start items-start border-[#e5e5e5] rounded-lg overflow-hidden shadow-lg [&>div]:px-2 sm:[&>div]:px-4 [&>div]:pt-2 sm:[&>div]:pt-4">
    <div class="flex gap-2 justify-start items-center w-full">
        <img class="w-10 sm:w-13.75" width="50" height="50" src="/icons/uptime.png" alt="uptime" />
        <h3 class="text-xl sm:text-[26px] font-semibold">Uptime</h3>

        <div class="ms-auto flex gap-2 justify-center items-center text-sm sm:pb-2">
            <span
                class={`flex justify-center items-center text-xl sm:text-[27px] ${
                    data[data.length - 1]
                        ? data[data.length - 1] < 65
                            ? 'text-green-700 h-5!'
                            : data[data.length - 1] < 85
                              ? 'text-yellow-500 h-7!'
                              : 'text-red-700 h-10!'
                        : 'bg-black/20'
                }`}>{data[data.length - 1]} ms</span>
            <img class="size-6 sm:size-9" width="35" height="35" src="/icons/time.png" alt="chart" />
        </div>
    </div>

    <div class="w-full" bind:this={chartEl}></div>
</div>
