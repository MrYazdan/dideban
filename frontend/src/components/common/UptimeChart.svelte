<script>
    import { onMount, onDestroy } from 'svelte';
    import ApexCharts from 'apexcharts';

    // props
    export let prices = [];
    export let dates = [];
    export let height = 395;

    let chartEl;
    let chart;

    onMount(() => {
        const options = {
            series: [
                {
                    name: 'Price',
                    data: prices,
                },
            ],

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
                show: false, // حذف خطوط زمینه
            },

            dataLabels: {
                enabled: false,
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
                    opacityTo: 0.05,
                    stops: [0, 90, 100],
                },
            },

            labels: dates,

            xaxis: {
                type: 'datetime',
                axisBorder: { show: true },
                axisTicks: { show: true },
            },

            yaxis: {
                opposite: true,
                axisBorder: { show: false },
                axisTicks: { show: false },
            },

            legend: {
                show: false,
            },
        };

        chart = new ApexCharts(chartEl, options);
        chart.render();
    });

    onDestroy(() => {
        chart?.destroy();
    });
</script>

<div
    class="w-full border flex flex-col justify-start items-start border-[#e5e5e5] rounded-lg overflow-hidden shadow-lg p-4">
    <div class="flex gap-2 justify-start items-center w-full">
        <img width="50" height="50" src="/icons/uptime.png" alt="uptime" />
        <h3 class="text-[26px] font-semibold">Uptime</h3>

        <div class="ms-auto pe-2 flex gap-2 justify-center items-center text-sm pb-2">
            <span class="text-[27px] text-green-700">1000ms</span>
            <img width="35" height="35" src="/icons/time.png" alt="chart" />
        </div>
    </div>

    <div class="w-full" bind:this={chartEl}></div>
</div>
