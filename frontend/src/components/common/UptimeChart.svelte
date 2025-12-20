<script>
    import { onMount, onDestroy } from 'svelte';
    import ApexCharts from 'apexcharts';

    // props
    export let prices = [];
    export let dates = [];

    let chartEl;
    let chart;

    onMount(() => {
        const options = {
            series: [
                {
                    data: prices,
                },
            ],

            chart: {
                type: 'area',
                height: 270,
                zoom: {
                    enabled: false,
                },
                toolbar: {
                    show: false,
                },
            },

            dataLabels: {
                enabled: false,
            },

            stroke: {
                curve: 'straight',
            },

            title: {
                text: 'Fundamental Analysis of Stocks',
                align: 'left',
            },

            subtitle: {
                text: 'Price Movements',
                align: 'left',
            },

            labels: dates,

            xaxis: {
                type: 'datetime',
            },

            yaxis: {
                opposite: true,
            },
            grid: {
                yaxis: { lines: { show: false } },
            },

            legend: {
                horizontalAlign: 'left',
            },
        };

        chart = new ApexCharts(chartEl, options);
        chart.render();
    });

    onDestroy(() => {
        chart?.destroy();
    });
</script>

<div bind:this={chartEl} class="w-full border overflow-hidden rounded-lg border-[#e5e5e5] pe-4 py-4 shadow-lg"></div>
