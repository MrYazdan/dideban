<script>
    import { onMount, onDestroy } from 'svelte';
    import ApexCharts from 'apexcharts';

    let chart;
    let chartEl;

    const options = {
        series: [
            {
                name: 'Sales',
                data: [33, 13, 45, 95, 12, 15, 34, 61, 66, 82, 1, 36],
            },
        ],
        chart: {
            type: 'line',
            height: 150,
            fontFamily: 'IRANSans, sans-serif',
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
            axisBorder: { show: false },
            axisTicks: { show: false },
            labels: { show: false },
            tooltip: { enabled: false },
        },

        yaxis: {
            min: 0,
            max: 100,
            tickAmount: 5,
            labels: {
                formatter: val => val.toFixed(0),
            },
            title: {
                text: 'Value',
            },
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
            theme: 'dark',
            shared: false,
            intersect: true,
            y: {
                formatter: val => `${val}`,
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
                        text: '',
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
        chart = new ApexCharts(chartEl, options);
        chart.render();
    });

    onDestroy(() => {
        chart?.destroy();
    });
</script>

<div class="w-full" bind:this={chartEl}></div>
