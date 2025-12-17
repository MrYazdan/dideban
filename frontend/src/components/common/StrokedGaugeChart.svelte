<script>
    import { onMount, onDestroy } from 'svelte';
    import ApexCharts from 'apexcharts';

    let chart;
    let chartEl;

    const options = {
        series: [67],
        chart: {
            type: 'radialBar',
            height: 350,
            offsetY: -10,
            fontFamily: 'IRANSans, sans-serif',
        },

        plotOptions: {
            radialBar: {
                startAngle: -135,
                endAngle: 135,
                hollow: {
                    size: '60%',
                },
                dataLabels: {
                    name: {
                        fontSize: '16px',
                        offsetY: 120,
                    },
                    value: {
                        offsetY: 76,
                        fontSize: '22px',
                        formatter: val => `${val}%`,
                    },
                },
            },
        },

        fill: {
            type: 'gradient',
            gradient: {
                shade: 'dark',
                shadeIntensity: 0.15,
                inverseColors: false,
                opacityFrom: 1,
                opacityTo: 1,
                stops: [0, 50, 65, 91],
            },
        },

        stroke: {
            dashArray: 4,
        },

        labels: ['Median Ratio'],
    };

    onMount(() => {
        chart = new ApexCharts(chartEl, options);
        chart.render();
    });

    onDestroy(() => {
        chart?.destroy();
    });
</script>

<div class="w-full max-w-47 m-auto" bind:this={chartEl}></div>
