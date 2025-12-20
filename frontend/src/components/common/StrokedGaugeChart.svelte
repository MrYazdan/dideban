<script>
    import { onMount, onDestroy } from 'svelte';
    import ApexCharts from 'apexcharts';
    const { title, value } = $props();

    let chart;
    let chartEl;
    function getColor(val) {
        if (val < 65) return '#22c55e'; // سبز
        if (val <= 85) return '#eab308'; // زرد
        return '#ef4444'; // قرمز
    }

    let options = {
        series: [value],
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
                    size: '50%',
                },
                dataLabels: {
                    name: {
                        fontSize: '14px',
                        offsetY: 120,
                        color: '#000000',
                    },
                    value: {
                        color: getColor(value),
                        offsetY: 76,
                        fontSize: '22px',
                        formatter: val => `${val}%`,
                    },
                },
            },
        },

        fill: {
            type: 'gradient',
            colors: [getColor(value)],
        },

        stroke: {
            dashArray: 4,
        },

        labels: [title],
    };

    onMount(() => {
        chart = new ApexCharts(chartEl, options);
        chart.render();
    });

    onDestroy(() => {
        chart?.destroy();
    });
</script>

<div class="w-fit max-w-47 m-auto" bind:this={chartEl}></div>
