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

        stroke: {
            dashArray: 4,
        },
    };

    onMount(() => {
        chart = new ApexCharts(chartEl, {
            ...options,
            labels: [title],
            fill: {
                type: 'gradient',
                colors: [getColor(value)],
            },
            series: [value],
            responsive: [
                {
                    breakpoint: 768,
                    options: {
                        plotOptions: {
                            radialBar: {
                                dataLabels: {
                                    name: {
                                        fontSize: '10px',
                                        offsetY: 55,
                                        color: '#000000',
                                    },
                                    value: {
                                        color: getColor(value),
                                        offsetY: -12,
                                        fontSize: '15px',
                                        formatter: val => `${val}%`,
                                    },
                                },
                            },
                        },
                    },
                },
                {
                    breakpoint: 1536,
                    options: {
                        chart: {
                            height: 120,
                        },
                        plotOptions: {
                            radialBar: {
                                hollow: { size: '45%' },
                                dataLabels: {
                                    name: {
                                        fontSize: '12px',
                                        offsetY: 90,
                                    },
                                    value: {
                                        fontSize: '18px',
                                        offsetY: 50,
                                    },
                                },
                            },
                        },
                    },
                },
            ],
        });
        chart.render();
    });

    onDestroy(() => {
        chart?.destroy();
    });
</script>

<div class="w-fit max-w-28 sm:max-w-47 xl:max-w-30 m-auto" bind:this={chartEl}></div>
