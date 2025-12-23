<script>
    import { onMount, onDestroy } from 'svelte';
    import ApexCharts from 'apexcharts';
    const { visibleSeries, data, isMobile } = $props();

    let chartEl;
    let chart;

    const options = {
        chart: {
            height: 221,
            type: 'area',
            padding: {
                top: 0,
                right: 0,
                bottom: 0,
                left: 0,
            },
            background: '#0d0d0d',
            toolbar: {
                show: true,
                tools: {
                    download: true,
                    selection: true,
                    zoom: true,
                    zoomin: true,
                    zoomout: true,
                    pan: true,
                    reset: true,
                },
                export: {
                    csv: {
                        enabled: true,
                        filename: 'machine-main-csv',
                    },
                    svg: {
                        enabled: true,
                        filename: 'machine-main-svg',
                    },
                    png: {
                        enabled: true,
                        filename: 'machine-main-png',
                    },
                },
            },
        },

        dataLabels: {
            enabled: false,
        },

        stroke: {
            curve: 'smooth',
            width: 2,
        },

        fill: {
            type: 'gradient',
            colors: ['#3b82f6', '#a855f7', '#10b981'],
            gradient: {
                shade: '',
                type: 'vertical',
                shadeIntensity: 0.5,
                gradientToColors: ['#3b82f6', '#a855f7', '#10b981'],
                opacityFrom: 0.6,
                opacityTo: 0.2,
                stops: [0, 65, 100],
            },
        },

        grid: {
            show: false,
        },

        legend: {
            show: false,
        },

        xaxis: {
            type: 'category',
            tickAmount: 5,

            axisBorder: {
                show: true,
                color: 'rgba(153, 161, 175, 0.3)',
            },

            axisTicks: {
                show: true,
                color: 'rgba(153, 161, 175, 0.5)',
                height: 4,
            },
        },
        yaxis: {
            min: 0,
            max: 100,
            labels: {
                formatter: val => Math.round(val),
            },
        },

        tooltip: {
            enabled: true,
            shared: true,
            intersect: false,
            x: {
                format: 'dd/MM/yy HH:mm',
            },
            y: {
                formatter: val => (val ? `${val} %` : `-`),
            },
        },

        colors: ['#3b82f6', '#a855f7', '#10b981'],

        theme: {
            mode: 'light',
        },
    };
    onMount(() => {
        chart = new ApexCharts(chartEl, { ...options, series: data });
        chart.render();
    });

    $effect(() => {
        if (chart && visibleSeries) {
            Object.entries(visibleSeries).forEach(([name, isVisible]) => {
                if (isVisible) {
                    chart.showSeries(name);
                } else {
                    chart.hideSeries(name);
                }
            });
        }
    });
    $effect(() => {
        if (chart) {
            chart.updateOptions({
                chart: {
                    height: isMobile ? 300 : 221,
                },
            });
        }
    });

    onDestroy(() => {
        chart?.destroy();
    });
</script>

<div bind:this={chartEl} class="w-full"></div>
