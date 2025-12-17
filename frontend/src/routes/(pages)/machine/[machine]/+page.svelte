<script>
    import UptimeChart from './../../../../components/common/UptimeChart.svelte';
    import { fade } from 'svelte/transition';
    import formatTimeAgo from '../../../../utils/formatTimeAgo.svelte';
    import formatFullDate from '../../../../utils/formatFullDate.js';
    import { MACHINE } from './constant.svelte.js';
    import { page } from '$app/stores';
    import Chart from '../../../../components/common/Chart.svelte';
    import InitialChart from '../../../../components/common/InitialChart.svelte';
    import StrokedGaugeChart from '../../../../components/common/StrokedGaugeChart.svelte';

    let width = $state(0);
    let itemHoveredDetail = $state(null);
    let activeBar = $state(null);
    let isMobile = $derived(width < 640);
    let visibleSeries = $state(Object.fromEntries(MACHINE.stats.map(stat => [stat.name, true])));
    let filteredStats = $derived(MACHINE.stats.filter(status => visibleSeries[status.name]));
    let activeRange = $state('1m');
    const githubdata = {
        series: [
            [1454976000000, 5],
            [1455062400000, 3],
            [1455148800000, 6],
            [1455235200000, 2],
            // ...
        ],
    };

    function toggle(name) {
        const activeCount = Object.values(visibleSeries).filter(Boolean).length;

        if (activeCount === 1 && visibleSeries[name]) return;

        visibleSeries = {
            ...visibleSeries,
            [name]: !visibleSeries[name],
        };
    }
</script>

<svelte:window bind:innerWidth={width} />

<div class="w-full h-auto flex flex-col gap-4 pb-3 relative">
    <h2 class="h-12.5! text-2xl w-full flex items-center border-b border-[#e5e5e5] capitalize">
        {$page.url.pathname.split('/')[2]}
    </h2>

    <div class="w-full">
        <div
            class="flex justify-center items-center gap-2 [&>button]:cursor-pointer bg-[#f5f5f5] w-fit p-2 rounded-lg [&>button]:rounded-md [&>button]:px-4 [&>button]:py-2">
            {#each Object.keys(visibleSeries) as option}
                <button class="capitalize {visibleSeries[option] ? 'bg-[#fefefe] ' : ''}" onclick={() => toggle(option)}
                    >{option}</button>
            {/each}
        </div>
    </div>

    <div
        class="grid grid-cols-3 justify-center items-start gap-2 w-full h-auto [&>div]:px-4 [&>div]:pt-4 [&>div]:border [&>div]:border-[#e5e5e5] [&>div]:rounded-lg">
        <div class="flex-1 flex flex-col gap-1 justify-start items-start shadow-lg">
            <div class="flex gap-4 text-4xl font-semibold items-center">
                <div class="p-2 rounded-lg">
                    <img width="50" src="/icons/cpu.png" alt="cpu" />
                </div>
                <span>CPU</span>
            </div>

            <InitialChart />

            <div class="w-full flex justify-between items-center">
                <StrokedGaugeChart />
                <StrokedGaugeChart />
                <StrokedGaugeChart />
            </div>
        </div>
        <div class="flex-1 flex flex-col gap-1 justify-start items-start shadow-lg">
            <div class="flex gap-4 text-4xl font-semibold items-center">
                <div class="p-2 rounded-lg">
                    <img width="50" src="/icons/cpu.png" alt="cpu" />
                </div>
                <span>CPU</span>
            </div>

            <InitialChart />

            <div class="w-full flex justify-between items-center">
                <StrokedGaugeChart />
                <StrokedGaugeChart />
                <StrokedGaugeChart />
            </div>
        </div>
        <div class="flex-1 flex flex-col gap-1 justify-start items-start shadow-lg">
            <div class="flex gap-4 text-4xl font-semibold items-center">
                <div class="p-2 rounded-lg">
                    <img width="50" src="/icons/cpu.png" alt="cpu" />
                </div>
                <span>CPU</span>
            </div>

            <InitialChart />

            <div class="w-full flex justify-between items-center">
                <StrokedGaugeChart />
                <StrokedGaugeChart />
                <StrokedGaugeChart />
            </div>
        </div>
    </div>

    <UptimeChart {githubdata} />
    {#if Object.values(visibleSeries).some(item => item)}
        <Chart
            {isMobile}
            {visibleSeries}
            data={MACHINE.stats.map(stat => ({
                name: stat.name,
                data: isMobile ? stat.detail.slice(-15).map(d => d.loaded ?? 0) : stat.detail.map(d => d.loaded ?? 0),
            }))} />

        {#each filteredStats as status}
            <div class="w-full flex flex-col gap-5 relative shadow-md hover: transition-all">
                <div
                    class="relative flex flex-col lg:flex-row h-27.5 rounded-lg md:border md:border-[#e5e5e5] md:px-5 md:py-3">
                    <!-- status indicator -->
                    <div class="absolute top-3.5 end-0 md:end-5 md:top-2 flex gap-1">
                        <div class="text-gray-500 text-xs flex items-baseline gap-2">
                            <div
                                class="size-2.5 rounded-full {status.detail.at(-1)?.loaded < 65
                                    ? 'bg-green-700'
                                    : status.detail.at(-1)?.loaded < 85
                                      ? 'bg-yellow-500'
                                      : 'bg-red-600'}">
                            </div>
                            <span>{formatTimeAgo(status.updateAt)}</span>
                        </div>
                    </div>

                    <!-- machine name -->
                    <h4
                        class="my-auto text-xl border-s-3 ps-2 w-30 capitalize
                        {status.detail.at(-1)?.loaded < 65
                            ? 'border-s-green-700'
                            : status.detail.at(-1)?.loaded < 85
                              ? 'border-s-yellow-500'
                              : 'border-s-red-600'}">
                        {status.name}
                    </h4>

                    <!-- hover info -->
                    <div
                        class="absolute lg:static lg:w-full top-27 sm:top-3 max-lg:start-1/2 max-lg:-translate-x-1/2 lg:start-0 justify-center items-center text-xs sm:text-sm flex lg:flex-col gap-7 lg:gap-1">
                        {#if itemHoveredDetail?.name === status?.name}
                            <div class="flex [&>div]:text-nowrap">
                                <div in:fade={{ duration: 500 }}>Total</div>
                                <div in:fade={{ duration: 500 }}>: {status.total}</div>
                            </div>
                            <div class="flex [&>div]:text-nowrap">
                                <div in:fade={{ duration: 1000 }}>
                                    {itemHoveredDetail?.status?.usage ? 'Usage' : null}
                                </div>
                                <div in:fade={{ duration: 1000 }}>
                                    {itemHoveredDetail?.status?.usage ? ': ' + itemHoveredDetail?.status?.usage : null}
                                </div>
                            </div>

                            <div class="flex [&>div]:text-nowrap">
                                <div in:fade={{ duration: 1300 }}>
                                    {itemHoveredDetail?.status?.loaded ? 'Loaded' : null}
                                </div>

                                <div in:fade={{ duration: 1300 }}>
                                    {itemHoveredDetail?.status?.loaded
                                        ? ': ' + itemHoveredDetail?.status?.loaded + ' %'
                                        : null}
                                </div>
                            </div>
                        {/if}
                    </div>

                    <!-- bars -->
                    <div class="w-full lg:w-fit ms-auto flex justify-center my-auto">
                        <div class="flex gap-1">
                            {#each isMobile ? status.detail.slice(-28) : status.detail as detail}
                                <div
                                    role="presentation"
                                    onmouseenter={() => {
                                        if (!isMobile) {
                                            itemHoveredDetail = {
                                                status: { ...detail },
                                                name: status.name,
                                            };
                                        }
                                    }}
                                    onmouseleave={() => {
                                        if (!isMobile) {
                                            itemHoveredDetail = null;
                                        }
                                    }}
                                    onclick={() => {
                                        if (isMobile) {
                                            itemHoveredDetail = {
                                                status: { ...detail },
                                                name: status.name,
                                            };
                                        }
                                    }}
                                    class="h-10 w-[2%] lg:w-2 min-w-2 rounded-full cursor-pointer transition-all {detail.loaded
                                        ? detail.loaded < 65
                                            ? 'bg-green-700'
                                            : detail.loaded < 85
                                              ? 'bg-yellow-500'
                                              : 'bg-red-700'
                                        : 'bg-black/20'}">
                                    <div
                                        class="opacity-0 group-hover:opacity-100 absolute top-10 lg:top-12 start-1/2 -translate-x-1/2 text-gray-400 text-xs">
                                        {formatFullDate(status.updateAt)}
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>
                </div>
            </div>
        {/each}
    {/if}
</div>
