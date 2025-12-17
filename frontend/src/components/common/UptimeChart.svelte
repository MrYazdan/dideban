<script>
    import { onMount, onDestroy } from 'svelte';
    import ApexCharts from 'apexcharts';

    export let githubdata = { series: [] };

    let chartEl;
    let chart;

    function updateCommits(chart) {
        const el = chartEl?.previousElementSibling?.querySelector('.commits');
        if (!el) return;

        const { minX, maxX } = chart.w.globals;
        el.textContent = chart.getSeriesTotalXRange(minX, maxX);
    }

    onMount(() => {
        const options = {
            series: [
                {
                    name: 'Commits',
                    data: githubdata.series,
                },
            ],

            chart: {
                id: 'chartyear',
                type: 'area',
                height: 180,
                background: '#F6F8FA',
                fontFamily: 'IRANSans, sans-serif',
                toolbar: {
                    show: false,
                },
                zoom: {
                    enabled: false,
                },
                animations: {
                    enabled: true,
                    easing: 'easeinout',
                    speed: 600,
                },
                events: {
                    mounted: updateCommits,
                    updated: updateCommits,
                },
            },

            theme: {
                mode: 'light',
            },

            colors: ['#FF7F00'],

            stroke: {
                show: true,
                width: 0,
                curve: 'monotoneCubic',
                lineCap: 'butt',
            },

            fill: {
                type: 'solid',
                opacity: 1,
            },

            dataLabels: {
                enabled: false,
            },

            markers: {
                size: 0,
                hover: {
                    size: 4,
                },
            },

            grid: {
                show: false,
            },

            legend: {
                show: false,
            },

            xaxis: {
                type: 'datetime',
                axisBorder: {
                    show: false,
                },
                axisTicks: {
                    show: false,
                },
                labels: {
                    show: false,
                },
                tooltip: {
                    enabled: false,
                },
            },

            yaxis: {
                show: true,
                forceNiceScale: true,
                labels: {
                    show: true,
                },
            },

            tooltip: {
                enabled: true,
                shared: false,
                intersect: false,
                x: {
                    show: true,
                    format: 'dd MMM yyyy',
                },
            },

            states: {
                normal: {
                    filter: { type: 'none' },
                },
                hover: {
                    filter: { type: 'lighten' },
                },
                active: {
                    filter: { type: 'darken' },
                },
            },
        };

        chart = new ApexCharts(chartEl, options);
        chart.render();
    });

    onDestroy(() => {
        chart?.destroy();
    });
</script>

<div class="w-full flex flex-col">
    <h2 class="text-4xl">UPTIME</h2>

    <div bind:this={chartEl}></div>
</div>
