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
    import ProgressBar from '../../../../components/common/ProgressBar.svelte';

    let width = $state(0);
    let itemHoveredDetail = $state(null);
    let activeBar = $state(null);
    let isMobile = $derived(width < 640);
    // let visibleSeries = $state(Object.fromEntries(MACHINE.stats.map(stat => [stat.name, true])));
    // let filteredStats = $derived(MACHINE.stats.filter(status => visibleSeries[status.name]));
    let activeRange = $state('1m');
    const githubdata = {
        series: [
            [1454976000000, 50],
            [1455062400000, 30],
            [1455148800000, 80],
            [1455235200000, 30],
            // ...
        ],
    };

    // function toggle(name) {
    //     const activeCount = Object.values(visibleSeries).filter(Boolean).length;

    //     if (activeCount === 1 && visibleSeries[name]) return;

    //     visibleSeries = {
    //         ...visibleSeries,
    //         [name]: !visibleSeries[name],
    //     };
    // }
</script>

<svelte:window bind:innerWidth={width} />

<div class="w-full h-auto flex flex-col gap-4 pb-3 relative">
    <h2 class="h-12.5! text-2xl w-full flex items-center border-b border-[#e5e5e5] capitalize">
        {$page.url.pathname.split('/')[2]}
    </h2>
    <!-- 
    <div class="w-full">
        <div
            class="flex justify-center items-center gap-2 [&>button]:cursor-pointer bg-[#f5f5f5] w-fit p-2 rounded-lg [&>button]:rounded-md [&>button]:px-4 [&>button]:py-2">
            {#each Object.keys(visibleSeries) as option}
                <button class="capitalize {visibleSeries[option] ? 'bg-[#fefefe] ' : ''}" onclick={() => toggle(option)}
                    >{option}</button>
            {/each}
        </div>
    </div> -->

    <div
        class="grid grid-cols-3 justify-center items-start gap-2 w-full h-auto [&>div]:px-4 [&>div]:pt-4 [&>div]:border [&>div]:border-[#e5e5e5] [&>div]:rounded-lg">
        <div class="flex-1 h-full flex flex-col gap-1 justify-start items-start shadow-lg">
            <div class="flex w-full gap-4 text-[26px] font-semibold items-center">
                <div class="p-2 rounded-lg">
                    <img width="55" src="/icons/cpu.png" alt="cpu" />
                </div>
                <span>CPU</span>

                <div class="ms-auto pe-2 flex gap-2 justify-center items-start text-sm pb-2">
                    <span
                        class="text-[27px] {MACHINE.cpu[MACHINE.cpu.length - 1].usage_percent
                            ? MACHINE.cpu[MACHINE.cpu.length - 1].usage_percent < 65
                                ? 'text-green-700'
                                : MACHINE.cpu[MACHINE.cpu.length - 1].usage_percent < 85
                                  ? 'text-yellow-500'
                                  : 'text-red-700'
                            : 'bg-black/20'}">{MACHINE.cpu[MACHINE.cpu.length - 1].usage_percent} %</span>
                    <img width="40" height="40" src="/icons/chart.png" alt="chart" />
                </div>
            </div>
            <div class="w-full flex justify-start my-auto px-3 mt-5">
                <div class="flex gap-1.5 items-start">
                    {#each isMobile ? MACHINE.cpu.slice(-28) : MACHINE.cpu as detail}
                        <div
                            class="h-10 min-w-[8.4px] rounded-full cursor-pointer transition-all {detail.usage_percent
                                ? detail.usage_percent < 65
                                    ? 'bg-green-700 h-5!'
                                    : detail.usage_percent < 85
                                      ? 'bg-yellow-500 h-7!'
                                      : 'bg-red-700 h-10!'
                                : 'bg-black/20'}">
                            <div
                                class="opacity-0 group-hover:opacity-100 absolute top-10 lg:top-12 start-1/2 -translate-x-1/2 text-gray-400 text-xs">
                                {formatFullDate(detail.timestamp_ms)}
                            </div>
                        </div>
                    {/each}
                </div>
            </div>

            <InitialChart name="cpu" data={MACHINE.cpu.map(item => item.usage_percent)} />

            <div class="w-full flex justify-center items-center">
                <StrokedGaugeChart value={MACHINE.cpu[MACHINE.cpu.length - 1].load_1} title="Load Avg (1m)" />
                <StrokedGaugeChart value={MACHINE.cpu[MACHINE.cpu.length - 1].load_5} title="Load Avg (5m)" />
                <StrokedGaugeChart value={MACHINE.cpu[MACHINE.cpu.length - 1].load_15} title="Load Avg (15m)" />
            </div>
        </div>
        <div class="flex-1 h-full flex flex-col gap-1 justify-start items-start shadow-lg">
            <div class="flex gap-4 text-[26px] font-semibold items-center w-full">
                <div class="p-2 rounded-lg">
                    <img width="55" src="/icons/disk.png" alt="disk" />
                </div>
                <span>DISK</span>
                <div class="ms-auto pe-2 flex gap-2 justify-center items-start text-sm pb-2">
                    <span
                        class="text-[27px] {MACHINE.disk[MACHINE.disk.length - 1].usage_percent
                            ? MACHINE.disk[MACHINE.disk.length - 1].usage_percent < 65
                                ? 'text-green-700'
                                : MACHINE.disk[MACHINE.disk.length - 1].usage_percent < 85
                                  ? 'text-yellow-500'
                                  : 'text-red-700'
                            : 'bg-black/20'}">{MACHINE.disk[MACHINE.disk.length - 1].usage_percent} %</span>
                    <img width="40" height="40" src="/icons/chart.png" alt="chart" />
                </div>
            </div>
            <div class="w-full flex justify-start my-auto px-3 mt-5">
                <div class="flex gap-1.5 items-start">
                    {#each isMobile ? MACHINE.disk.slice(-28) : MACHINE.disk as detail}
                        <div
                            class="h-10 min-w-[8.4px] rounded-full cursor-pointer transition-all {detail.usage_percent
                                ? detail.usage_percent < 65
                                    ? 'bg-green-700 h-5!'
                                    : detail.usage_percent < 85
                                      ? 'bg-yellow-500 h-7!'
                                      : 'bg-red-700 h-10!'
                                : 'bg-black/20'}">
                            <div
                                class="opacity-0 group-hover:opacity-100 absolute top-10 lg:top-12 start-1/2 -translate-x-1/2 text-gray-400 text-xs">
                                {formatFullDate(detail.timestamp_ms)}
                            </div>
                        </div>
                    {/each}
                </div>
            </div>

            <InitialChart name="disk" data={MACHINE.disk.map(item => item.usage_percent)} />

            <div class="w-full grid grid-cols-2 justify-center items-center">
                <StrokedGaugeChart value={MACHINE.disk[MACHINE.disk.length - 1].usage_percent} title="Usage" />
                <div class="w-full h-full flex justify-center items-center">
                    <div class="justify-center items-center h-full flex flex-col gap-3 text-sm font-semibold">
                        <div class="flex justify-start items-center gap-2">
                            <span>Total :</span>
                            <span>{MACHINE.disk[MACHINE.disk.length - 1].total_gb} GB</span>
                        </div>

                        <div class="flex justify-start items-center gap-2">
                            <span>Used :</span>
                            <span>{MACHINE.disk[MACHINE.disk.length - 1].used_gb} GB</span>
                        </div>
                        <div class="flex justify-start items-center gap-2">
                            <span>Available :</span>
                            <span
                                >{MACHINE.disk[MACHINE.disk.length - 1].total_gb -
                                    MACHINE.disk[MACHINE.disk.length - 1].used_gb} GB</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="flex-1 h-full flex flex-col gap-1 justify-start items-start shadow-lg">
            <div class="w-full flex gap-4 text-[26px] font-semibold items-center">
                <div class="p-2 rounded-lg">
                    <img width="50" src="/icons/memory.png" alt="memory" />
                </div>
                <span>MEMORY</span>
                <div class="ms-auto pe-2 flex gap-2 justify-center items-start text-sm pb-2">
                    <span
                        class="text-[27px] {MACHINE.memory[MACHINE.memory.length - 1].usage_percent
                            ? MACHINE.memory[MACHINE.memory.length - 1].usage_percent < 65
                                ? 'text-green-700'
                                : MACHINE.memory[MACHINE.memory.length - 1].usage_percent < 85
                                  ? 'text-yellow-500'
                                  : 'text-red-700'
                            : 'bg-black/20'}">{MACHINE.memory[MACHINE.memory.length - 1].usage_percent} %</span>
                    <img width="40" height="40" src="/icons/chart.png" alt="chart" />
                </div>
            </div>
            <div class="w-full flex justify-start my-auto px-3 mt-5">
                <div class="flex gap-1.5 items-start">
                    {#each isMobile ? MACHINE.memory.slice(-28) : MACHINE.memory as detail}
                        <div
                            class="h-10 min-w-[8.4px] rounded-full cursor-pointer transition-all {detail.usage_percent
                                ? detail.usage_percent < 65
                                    ? 'bg-green-700 h-5!'
                                    : detail.usage_percent < 85
                                      ? 'bg-yellow-500 h-7!'
                                      : 'bg-red-700 h-10!'
                                : 'bg-black/20'}">
                            <div
                                class="opacity-0 group-hover:opacity-100 absolute top-10 lg:top-12 start-1/2 -translate-x-1/2 text-gray-400 text-xs">
                                {formatFullDate(detail.timestamp_ms)}
                            </div>
                        </div>
                    {/each}
                </div>
            </div>

            <InitialChart name="memory" data={MACHINE.memory.map(item => item.usage_percent)} />

            <div class="w-full grid grid-cols-2 justify-center items-center">
                <StrokedGaugeChart value={MACHINE.memory[MACHINE.memory.length - 1].usage_percent} title="Usage" />
                <div class="w-full h-full flex justify-center items-center">
                    <div class="justify-center items-center h-full flex flex-col gap-3 text-sm font-semibold">
                        <div class="flex justify-start items-center gap-2">
                            <span>Total :</span>
                            <span
                                >{Number(MACHINE.memory[MACHINE.memory.length - 1].total_mb).toLocaleString()} Mb</span>
                        </div>

                        <div class="flex justify-start items-center gap-2">
                            <span>Used :</span>
                            <span>{Number(MACHINE.memory[MACHINE.memory.length - 1].used_mb).toLocaleString()} Mb</span>
                        </div>
                        <div class="flex justify-start items-center gap-2">
                            <span>Available :</span>
                            <span>
                                {(
                                    Number(MACHINE.memory[MACHINE.memory.length - 1].total_mb) -
                                    Number(MACHINE.memory[MACHINE.memory.length - 1].used_mb)
                                ).toLocaleString()}
                                Mb
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- <UptimeChart
        height={260}
        prices={[8107, 8128, 8122, 8165, 8340, 8423]}
        dates={[
            '2018-09-19T00:00:00.000Z',
            '2018-09-20T00:00:00.000Z',
            '2018-09-21T00:00:00.000Z',
            '2018-09-22T00:00:00.000Z',
            '2018-09-23T00:00:00.000Z',
            '2018-09-24T00:00:00.000Z',
        ]} /> -->
    <!--   {#if Object.values(visibleSeries).some(item => item)}
        <Chart
            {isMobile}
            {visibleSeries}
            data={MACHINE.stats.map(stat => ({
                name: stat.name,
                data: isMobile ? stat.detail.slice(-15).map(d => d.loaded ?? 0) : stat.detail.map(d => d.loaded ?? 0),
            }))} />
    {/if} -->
</div>
