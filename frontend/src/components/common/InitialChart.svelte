<script>
    import { onMount, onDestroy } from 'svelte';
    import ApexCharts from 'apexcharts';
    const { data, name } = $props();

    let chart;
    let chartEl;

    const options = {
        chart: {
            type: 'line',
            height: 150,
            background: '#fefefe',
            zoom: {
                enabled: false,
            },
            toolbar: {
                show: false,
            },
            animations: {
                enabled: true,
                easing: 'easeinout',
                speed: 800,
            },
        },

        theme: {
            mode: 'light',
            palette: 'palette1',
        },

        xaxis: {
            axisBorder: { show: true },
            axisTicks: { show: true },
            labels: { show: true },
            tooltip: { enabled: false },
        },
        yaxis: {
            show: false,
        },
        grid: {
            show: false,
            borderColor: '#e0e0e0',
            strokeDashArray: 4,
            position: 'back',
            xaxis: {
                lines: {
                    show: false,
                },
            },
            yaxis: {
                lines: {
                    show: false,
                },
            },
        },

        stroke: {
            width: 5,
            curve: 'monotoneCubic',
            colors: ['#0088ee'],
        },

        markers: {
            size: 5,
            strokeWidth: 2,
            strokeColors: '#fff',
            hover: {
                size: 7,
            },
        },

        dataLabels: {
            enabled: false,
        },

        legend: {
            show: true,
            position: 'top',
            horizontalAlign: 'right',
        },

        tooltip: {
            enabled: true,
            theme: 'light',
            shared: false,
            intersect: true,
            y: {
                formatter: val => `${val} %`,
            },
        },

        fill: {
            type: 'solid',
            opacity: 1,
        },

        annotations: {
            yaxis: [
                {
                    y: 85,
                    borderColor: '#ff0000',
                    strokeDashArray: 4,
                    label: {
                        style: {
                            color: '#fff',
                            background: '#ff0000',
                        },
                    },
                },
                {
                    y: 65,
                    borderColor: '#f0b100',
                    strokeDashArray: 4,
                },
            ],
        },

        states: {
            normal: {
                filter: {
                    type: 'none',
                },
            },
            hover: {
                filter: {
                    type: 'lighten',
                },
            },
            active: {
                filter: {
                    type: 'darken',
                },
            },
        },

        responsive: [
            {
                breakpoint: 768,
                options: {
                    chart: {
                        height: 300,
                    },
                    xaxis: {
                        labels: {
                            rotate: -90,
                        },
                    },
                },
            },
        ],
    };

    onMount(() => {
        chart = new ApexCharts(chartEl, {
            ...options,
            series: [
                {
                    name,
                    data: data.slice(-10),
                },
            ],
        });
        chart.render();
    });

    onDestroy(() => {
        chart?.destroy();
    });
</script>

<div class="w-full" bind:this={chartEl}></div>
