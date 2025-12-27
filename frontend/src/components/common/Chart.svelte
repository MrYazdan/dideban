<script>
    import { onMount, onDestroy } from 'svelte';
    import ApexCharts from 'apexcharts';
    const { visibleSeries, data, isMobile } = $props();

    let chartEl;
    const hexToRgba = (hex, opacity) => {
        const r = parseInt(hex.slice(1, 3), 16);
        const g = parseInt(hex.slice(3, 5), 16);
        const b = parseInt(hex.slice(5, 7), 16);
        return `rgba(${r}, ${g}, ${b}, ${opacity})`;
    };
    let chart;

    const options = {
        chart: {
            height: 221,
            type: 'area',
            zoom: {
                enabled: false,
            },
            padding: {
                top: 0,
                right: 0,
                bottom: 0,
                left: 0,
            },
            background: 'transparent',
            toolbar: {
                show: false,
                tools: {
                    download: false,
                    selection: false,
                    zoom: false,
                    zoomin: false,
                    zoomout: false,
                    pan: false,
                    reset: false,
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
                opacityTo: 0.1,
                stops: [0, 90, 100],
            },
        },

        grid: {
            show: false,
            padding: {
                top: 0,
                right: 0,
                bottom: 0,
                left: 0,
            },
        },

        legend: {
            show: false,
        },

        xaxis: {
            show: false,
            floating: true,
            type: 'numeric',
            tickAmount: 5,
            labels: { show: false },

            axisBorder: {
                show: false,
                color: 'rgba(153, 161, 175, 0.3)',
            },

            axisTicks: {
                show: false,
                color: 'rgba(153, 161, 175, 0.5)',
                height: 4,
            },
        },
        yaxis: {
            show: false,
            floating: true,
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
        chart = new ApexCharts(chartEl, {
            ...options,
            series: data,
            annotations: {
                points: data.flatMap((series, seriesIndex) =>
                    series.data
                        .map((y, pointIndex) => {
                            if (y < 70) return null;

                            const baseColor = options.colors[seriesIndex];
                            const opacity = y >= 80 ? 1 : 0.5;

                            return {
                                x: pointIndex + 1,
                                y,
                                seriesIndex,
                                marker: {
                                    size: 3,
                                    fillColor: hexToRgba(baseColor, opacity),
                                    strokeWidth: 7,
                                    strokeColor: hexToRgba(baseColor, opacity * 0.5),
                                },
                            };
                        })
                        .filter(Boolean),
                ),
            },
        });
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
                    height: isMobile ? 300 : 200,
                },
            });
        }
    });

    onDestroy(() => {
        chart?.destroy();
    });
</script>

<div bind:this={chartEl} class="w-full"></div>

<style>
    :global(.apexcharts-tooltip) {
        background: rgba(0, 0, 0, 0.4) !important;
        color: #ffffff !important;
        border-color: rgb(255, 255, 255, 0.2) !important;
        box-shadow: 0 0px 10px rgba(168, 167, 167, 0.1) !important;
        border-radius: 12px !important;
        padding: 10px 8px !important;
        backdrop-filter: blur(10px);
        -webkit-backdrop-filter: blur(10px);
    }

    :global(.apexcharts-tooltip-title) {
        display: none;
    }

    :global(.apexcharts-tooltip-text) {
        color: #ffffff !important;
    }
</style>
