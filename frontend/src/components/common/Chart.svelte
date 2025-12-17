<script>
    import { onMount, onDestroy } from 'svelte';
    import ApexCharts from 'apexcharts';
    const { visibleSeries, data, isMobile } = $props();

    let chartEl;
    let chart;

    const options = {
        chart: {
            height: 410,
            type: 'area',
            padding: {
                top: 0,
                right: 0,
                bottom: 0,
                left: 0,
            },
            background: '#fefefe',
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

            animations: {
                enabled: true,
                speed: 1000,
            },
        },

        dataLabels: {
            enabled: false,
        },

        stroke: {
            curve: 'smooth',
            width: 1,
        },
        markers: {
            size: 3,
            colors: ['#fefefe'],
            strokeColors: '#012a4a',
            strokeWidth: 1,
            hover: {
                size: 5,
            },
        },

        fill: {
            type: 'gradient',
            colors: ['#61a5c2', '#2ec4b6', '#84afe6'],
            gradient: {
                shade: '',
                type: 'vertical',
                shadeIntensity: 0.5,
                gradientToColors: ['#61a5c2', '#2ec4b6', '#84afe6'],
                opacityFrom: 1,
                opacityTo: 0.3,
                stops: [0, 65, 100],
            },
        },

        grid: {
            show: true,
            borderColor: '#e5e5e5',
            strokeDashArray: 3,
        },

        legend: {
            show: false,
        },

        xaxis: {
            type: 'category',
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

        colors: ['#61a5c2', '#2ec4b6', '#84afe6'],

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
                    height: isMobile ? 300 : 410,
                },
            });
        }
    });

    onDestroy(() => {
        chart?.destroy();
    });
</script>

<div bind:this={chartEl} class="w-full border overflow-hidden rounded-lg border-[#e5e5e5] pe-4 py-4 shadow-lg"></div>
