<script>
    import Chart from '../../components/common/Chart.svelte';
    import { theme } from '../../stores/theme.svelte';
    import { MACHINE } from './machine/[machine]/constant.svelte';

    let width = $state(0);
    let isMobile = $derived(width < 365);
    let isTablet = $derived(width > 1279 && width < 1536);
    let visibleSeries = $state({
        CPU: true,
        Memory: true,
        Disk: true,
    });
</script>

<svelte:window bind:innerWidth={width} />

<section class="w-full h-auto flex flex-col">
    <!-- Content of dashboard page -->
    <div class="w-full flex flex-col gap-7.75 p-7.75 pt-3">
        <div class="w-full h-103.25 flex gap-7.75">
            <div
                class="w-213 h-full p-6 rounded-[14px] dark:bg-[#0D0D0D] bg-[#FFFFFF] border border-[#0D0D0D]/5 dark:border-white/5">
                <div class="flex flex-col gap-4 items-start justify-around">
                    <div class="w-full flex flex-col justify-start items-start">
                        <span class="text-xl dark:text-white">Main Performance Overview</span>
                        <span class="text-sm text-[#99a1af]">System resource utilization trends</span>
                    </div>

                    <div class="w-full h-20.5 flex justify-around gap-4.75">
                        <div
                            class="h-full w-59.25 flex flex-col justify-start items-start gap-2 px-4 py-3 rounded-[10px] bg-[#F9FAFB] dark:bg-[#121212] border border-[#0D0D0D]/5 dark:border-white/5">
                            <div class="w-full flex justify-start items-center gap-1.5">
                                <span class="size-2 rounded-full bg-[#ad46ff]"></span>
                                <span class="text-sm text-[#6a7282]">CPU Average</span>
                            </div>

                            <span class="dark:text-white text-2xl">60.1%</span>
                        </div>
                        <div
                            class="h-full w-59.25 flex flex-col justify-start items-start gap-2 px-4 py-3 rounded-[10px] bg-[#F9FAFB] dark:bg-[#121212] border border-[#0D0D0D]/5 dark:border-white/5">
                            <div class="w-full flex justify-start items-center gap-1.5">
                                <span class="size-2 rounded-full bg-[#2b7fff]"></span>
                                <span class="text-sm text-[#6a7282]">Memory Average</span>
                            </div>

                            <span class="dark:text-white text-2xl">65.4%</span>
                        </div>
                        <div
                            class="h-full w-59.25 flex flex-col justify-start items-start gap-2 px-4 py-3 rounded-[10px] bg-[#F9FAFB] dark:bg-[#121212] border border-[#0D0D0D]/5 dark:border-white/5">
                            <div class="w-full flex justify-start items-center gap-1.5">
                                <span class="size-2 rounded-full bg-[#00bc7d]"></span>
                                <span class="text-sm text-[#6a7282]">Disk</span>
                            </div>

                            <span class="dark:text-white text-2xl">50.6%</span>
                        </div>
                    </div>

                    <Chart
                        {isMobile}
                        {visibleSeries}
                        data={[
                            {
                                name: 'CPU',
                                data: isMobile
                                    ? MACHINE.cpu.slice(-50).map(d => d.usage_percent ?? 0)
                                    : MACHINE.cpu.map(d => d.usage_percent ?? 0),
                            },
                            {
                                name: 'Memory',
                                data: isMobile
                                    ? MACHINE.memory.slice(-50).map(d => d.usage_percent ?? 0)
                                    : MACHINE.memory.map(d => d.usage_percent ?? 0),
                            },
                            {
                                name: 'Disk',
                                data: isMobile
                                    ? MACHINE.disk.slice(-50).map(d => d.usage_percent ?? 0)
                                    : MACHINE.disk.map(d => d.usage_percent ?? 0),
                            },
                        ]} />
                </div>
            </div>

            <div class="w-full h-full grid grid-cols-2 grid-rows-2 gap-7.75">
                <div
                    class="bg-white dark:bg-transparent border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] relative overflow-hidden group">
                    <div
                        class="absolute -top-20 end-0 size-0 rounded-full group-hover:top-2 group-hover:end-2 transition-all duration-700"
                        style="box-shadow: 0 0 300px 60px rgba(0,102,255,1);">
                        <div class="w-full h-full bg-white/5"></div>
                    </div>

                    <div class="p-6 flex flex-col">
                        <div
                            class="size-12.5 flex justify-center items-center rounded-[14px] border border-[#51a2ff]/20 bg-[#2B7FFF]/10 mb-3">
                            <img src="/icons/total.svg" alt="total" />
                        </div>
                        <div class="dark:text-white text-3xl mb-1">156</div>
                        <span class="text-sm text-[#99a1af]">Total Services</span>
                    </div>
                </div>
                <div
                    class="bg-white dark:bg-transparent border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] relative overflow-hidden group">
                    <div
                        class="absolute -top-20 end-0 size-0 rounded-full group-hover:top-2 group-hover:end-2 transition-all duration-700"
                        style="box-shadow: 0 0 300px 60px rgb(0,212,146);">
                        <div class="w-full h-full bg-white/5"></div>
                    </div>
                    <div class="p-6 flex flex-col">
                        <div
                            class="size-12.5 flex justify-center items-center rounded-[14px] border border-[#00bc7d]/20 bg-[#00BC7D]/10 mb-3">
                            <img src="/icons/tick.svg" alt="tick" />
                        </div>
                        <div class="text-[#00d492] text-3xl mb-1">142</div>
                        <span class="text-sm text-[#99a1af]">Up Services</span>
                    </div>
                </div>
                <div
                    class="bg-white dark:bg-transparent border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] relative overflow-hidden group">
                    <div
                        class="absolute -top-20 end-0 size-0 rounded-full group-hover:top-2 group-hover:end-2 transition-all duration-700"
                        style="box-shadow: 0 0 300px 60px rgb(255,100,103);">
                        <div class="w-full h-full bg-white/5"></div>
                    </div>
                    <div class="p-6 flex flex-col">
                        <div
                            class="size-12.5 flex justify-center items-center rounded-[14px] border border-[#fb2c36]/20 bg-[#FB2C36]/10 mb-3">
                            <img src="/icons/error.svg" alt="error" />
                        </div>
                        <div class="text-[#ff6467] text-3xl mb-1">8</div>
                        <span class="text-sm text-[#99a1af]">Down Services</span>
                    </div>
                </div>
                <div
                    class="bg-white dark:bg-transparent border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] relative overflow-hidden group">
                    <div
                        class="absolute -top-20 end-0 size-0 rounded-full group-hover:top-2 group-hover:end-2 transition-all duration-700"
                        style="box-shadow: 0 0 300px 60px rgb(252,200,0);">
                        <div class="w-full h-full bg-white/5"></div>
                    </div>
                    <div class="p-6 flex flex-col">
                        <div
                            class="size-12.5 flex justify-center items-center rounded-[14px] border border-[#fbbc05]/20 bg-[#F0B100]/10 mb-3">
                            <img src="/icons/warning.svg" alt="warning" />
                        </div>
                        <div class="text-[#fdc700] text-3xl mb-1">6</div>
                        <span class="text-sm text-[#99a1af]">Warning</span>
                    </div>
                </div>
            </div>
        </div>

        <div
            class="w-full p-6 rounded-[14px] bg-[#FFFFFF] dark:bg-[#0D0D0D] border border-[#0D0D0D]/5 dark:border-white/5">
            <div class="w-full flex justify-between items-start pb-8">
                <div class="flex flex-col gap-1">
                    <span class="text-xl dark:text-white">Machine Status</span>
                    <span class="text-sm text-[#99a1af]">Infrastructure Nodes & Cluster Health</span>
                </div>
                <div
                    class="w-12 h-10 flex justify-center items-center bg-[#22c55e]/10 rounded-lg text-xl text-[#10b981]">
                    +
                </div>
            </div>
            <div class="w-full grid grid-cols-3 gap-7.5">
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">192.168.1.45</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#F97316]/15 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F97316]/5">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#F97316]/10">
                            <img src="/icons/memory.svg" alt="memory" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#F97316]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 8px 0.5px #F97316;"
                            class="size-2.5 bg-[#F97316] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#EF4444]/15 rounded-[14px] flex flex-col py-8 gap-7 bg-[#EF4444]/5">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#EF4444]/10">
                            <img src="/icons/memory-red.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#F87171]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #EF4444;"
                            class="size-2.5 bg-[#EF4444] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
            </div>
        </div>
        <div
            class="w-full p-6 rounded-[14px] bg-[#FFFFFF] dark:bg-[#0D0D0D] border border-[#0D0D0D]/5 dark:border-white/5">
            <div class="w-full flex justify-between items-start pb-8">
                <div class="flex flex-col gap-1">
                    <span class="text-xl dark:text-white">Application Status</span>
                    <span class="text-sm text-[#99a1af]">Microservices Health & Availability</span>
                </div>
                <div
                    class="w-12 h-10 flex justify-center items-center bg-[#22c55e]/10 rounded-lg text-xl text-[#10b981]">
                    +
                </div>
            </div>
            <div class="w-full grid grid-cols-3 gap-7.5">
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#F97316]/15 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F97316]/5">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#F97316]/10">
                            <img src="/icons/memory.svg" alt="memory" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#F97316]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 8px 0.5px #F97316;"
                            class="size-2.5 bg-[#F97316] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#EF4444]/15 rounded-[14px] flex flex-col py-8 gap-7 bg-[#EF4444]/5">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#EF4444]/10">
                            <img src="/icons/memory-red.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#F87171]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #EF4444;"
                            class="size-2.5 bg-[#EF4444] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
                <div
                    class="w-122.75 h-37.75 border border-[#0D0D0D]/5 dark:border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#F9FAFB] dark:bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg dark:text-white">US-East-Cluster-01</span>
                            <span class="text-xs text-[#99a1af]">Uptime: 99.99%</span>
                        </div>
                        <div
                            style="box-shadow: 0 0 10px 1px #22c55e;"
                            class="size-2.5 bg-[#22c55e] rounded-full ms-auto mb-auto mt-2">
                        </div>
                    </div>
                    <div class="w-full flex gap-1 pb-2 justify-center items-start">
                        {#each MACHINE.cpu.slice(-50) as detail}
                            <div
                                class="size-1.25 rounded-[1px] {detail.usage_percent
                                    ? detail.usage_percent < 65
                                        ? 'bg-green-700'
                                        : detail.usage_percent < 85
                                          ? 'bg-[#F97316]'
                                          : 'bg-[#EF4444]'
                                    : 'bg-[#FFFFFF]/5'}">
                            </div>
                        {/each}
                    </div>
                </div>
            </div>
        </div>

        <div class="w-full p-6 rounded-[14px] bg-white dark:bg-[#0D0D0D] border border-[#0D0D0D]/5 dark:border-white/5">
            <div class="w-full flex justify-between items-start pb-8">
                <div class="flex flex-col gap-1">
                    <span class="text-xl dark:text-white">Agents Status</span>
                    <span class="text-sm text-[#99a1af]">Microservices Health & Availability</span>
                </div>
                <div
                    class="w-12 h-10 flex justify-center items-center bg-[#22c55e]/10 rounded-lg text-xl text-[#10b981]">
                    <img width="17" src="/icons/refresh-green.svg" alt="refresh" />
                </div>
            </div>

            <div class="w-full grid grid-cols-3 gap-7.5">
                <div
                    class="w-full p-5 rounded-xl bg-[#F9FAFB] dark:bg-[#121212] border border-[#0D0D0D]/5 dark:border-white/5 flex flex-col gap-4">
                    <div class="flex w-full justify-start items-center gap-4">
                        <div
                            class="flex justify-center items-center size-10 rounded-2xl bg-[#00d492]/10 border border-[#00d492]/30">
                            <img src="/icons/agent-green.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-base dark:text-white">Agent-Alpha</span>
                            <span class="text-xs text-[#99a1af]">192.168.1.45</span>
                        </div>
                        <div
                            class="flex justify-center items-center gap-2 px-3 py-1.25 ms-auto mb-auto bg-[#00bc7d]/20 border border-[#00bc7d]/30 rounded-full">
                            <span class="size-2 bg-[#00d492]"></span>
                            <span class="text-xs text-[#00d492]">Active</span>
                        </div>
                    </div>
                    <div class="w-full h-36.25 flex flex-col gap-3">
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">CPU</span>
                                <span class="dark:text-white">45%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">RAM</span>
                                <span class="dark:text-white">62%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">Storage</span>
                                <span class="dark:text-white">45%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div
                            class="w-full pt-2 border-t border-[#0D0D0D]/5 dark:border-white/5 flex justify-between items-center text-[#99a1af] text-xs">
                            <div class="flex items-center justify-start gap-1">
                                <img src="/icons/time.svg" alt="time" />
                                <span>Uptime</span>
                            </div>
                            <span class="dark:text-white">14d 6h 23m</span>
                        </div>
                    </div>
                </div>
                <div
                    class="w-full p-5 rounded-xl bg-[#F9FAFB] dark:bg-[#121212] border border-[#0D0D0D]/5 dark:border-white/5 flex flex-col gap-4">
                    <div class="flex w-full justify-start items-center gap-4">
                        <div
                            class="flex justify-center items-center size-10 rounded-2xl bg-[#00d492]/10 border border-[#00d492]/30">
                            <img src="/icons/agent-green.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-base dark:text-white">Agent-Alpha</span>
                            <span class="text-xs text-[#99a1af]">192.168.1.45</span>
                        </div>
                        <div
                            class="flex justify-center items-center gap-2 px-3 py-1.25 ms-auto mb-auto bg-[#00bc7d]/20 border border-[#00bc7d]/30 rounded-full">
                            <span class="size-2 bg-[#00d492]"></span>
                            <span class="text-xs text-[#00d492]">Active</span>
                        </div>
                    </div>
                    <div class="w-full h-36.25 flex flex-col gap-3">
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">CPU</span>
                                <span class="dark:text-white">45%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">RAM</span>
                                <span class="dark:text-white">62%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">Storage</span>
                                <span class="dark:text-white">45%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div
                            class="w-full pt-2 border-t border-[#0D0D0D]/5 dark:border-white/5 flex justify-between items-center text-[#99a1af] text-xs">
                            <div class="flex items-center justify-start gap-1">
                                <img src="/icons/time.svg" alt="time" />
                                <span>Uptime</span>
                            </div>
                            <span class="dark:text-white">14d 6h 23m</span>
                        </div>
                    </div>
                </div>
                <div
                    class="w-full p-5 rounded-xl bg-[#F9FAFB] dark:bg-[#121212] border border-[#0D0D0D]/5 dark:border-white/5 flex flex-col gap-4">
                    <div class="flex w-full justify-start items-center gap-4">
                        <div
                            class="flex justify-center items-center size-10 rounded-2xl bg-[#00d492]/10 border border-[#00d492]/30">
                            <img src="/icons/agent-green.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-base dark:text-white">Agent-Alpha</span>
                            <span class="text-xs text-[#99a1af]">192.168.1.45</span>
                        </div>
                        <div
                            class="flex justify-center items-center gap-2 px-3 py-1.25 ms-auto mb-auto bg-[#00bc7d]/20 border border-[#00bc7d]/30 rounded-full">
                            <span class="size-2 bg-[#00d492]"></span>
                            <span class="text-xs text-[#00d492]">Active</span>
                        </div>
                    </div>
                    <div class="w-full h-36.25 flex flex-col gap-3">
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">CPU</span>
                                <span class="dark:text-white">45%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">RAM</span>
                                <span class="dark:text-white">62%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">Storage</span>
                                <span class="dark:text-white">45%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div
                            class="w-full pt-2 border-t border-[#0D0D0D]/5 dark:border-white/5 flex justify-between items-center text-[#99a1af] text-xs">
                            <div class="flex items-center justify-start gap-1">
                                <img src="/icons/time.svg" alt="time" />
                                <span>Uptime</span>
                            </div>
                            <span class="dark:text-white">14d 6h 23m</span>
                        </div>
                    </div>
                </div>
                <div
                    class="w-full p-5 rounded-xl bg-[#EF4444]/5 border border-[#0D0D0D]/5 dark:border-white/5 flex flex-col gap-4">
                    <div class="flex w-full justify-start items-center gap-4">
                        <div
                            class="flex justify-center items-center size-10 rounded-2xl bg-[#EF4444]/10 border border-[#EF4444]/30">
                            <img src="/icons/agent-red.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-base dark:text-white">Agent-Alpha</span>
                            <span class="text-xs text-[#99a1af]">192.168.1.45</span>
                        </div>
                        <div
                            class="flex justify-center items-center gap-2 px-3 py-1.25 ms-auto mb-auto bg-[#EF4444]/20 border border-[#EF4444]/30 rounded-full">
                            <span class="size-2 bg-[#EF4444]"></span>
                            <span class="text-xs text-[#EF4444]">Active</span>
                        </div>
                    </div>
                    <div class="w-full h-36.25 flex flex-col gap-3">
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">CPU</span>
                                <span class="dark:text-white">45%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">RAM</span>
                                <span class="dark:text-white">62%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">Storage</span>
                                <span class="dark:text-white">45%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div
                            class="w-full pt-2 border-t border-[#0D0D0D]/5 dark:border-white/5 flex justify-between items-center text-[#99a1af] text-xs">
                            <div class="flex items-center justify-start gap-1">
                                <img src="/icons/time.svg" alt="time" />
                                <span>Uptime</span>
                            </div>
                            <span class="dark:text-white">14d 6h 23m</span>
                        </div>
                    </div>
                </div>
                <div
                    class="w-full p-5 rounded-xl bg-[#F9FAFB] dark:bg-[#121212] border border-[#0D0D0D]/5 dark:border-white/5 flex flex-col gap-4">
                    <div class="flex w-full justify-start items-center gap-4">
                        <div
                            class="flex justify-center items-center size-10 rounded-2xl bg-[#00d492]/10 border border-[#00d492]/30">
                            <img src="/icons/agent-green.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-base dark:text-white">Agent-Alpha</span>
                            <span class="text-xs text-[#99a1af]">192.168.1.45</span>
                        </div>
                        <div
                            class="flex justify-center items-center gap-2 px-3 py-1.25 ms-auto mb-auto bg-[#00bc7d]/20 border border-[#00bc7d]/30 rounded-full">
                            <span class="size-2 bg-[#00d492]"></span>
                            <span class="text-xs text-[#00d492]">Active</span>
                        </div>
                    </div>
                    <div class="w-full h-36.25 flex flex-col gap-3">
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">CPU</span>
                                <span class="dark:text-white">45%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">RAM</span>
                                <span class="dark:text-white">62%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">Storage</span>
                                <span class="dark:text-white">45%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div
                            class="w-full pt-2 border-t border-[#0D0D0D]/5 dark:border-white/5 flex justify-between items-center text-[#99a1af] text-xs">
                            <div class="flex items-center justify-start gap-1">
                                <img src="/icons/time.svg" alt="time" />
                                <span>Uptime</span>
                            </div>
                            <span class="dark:text-white">14d 6h 23m</span>
                        </div>
                    </div>
                </div>
                <div
                    class="w-full p-5 rounded-xl bg-[#F0B100]/5 border border-[#0D0D0D]/5 dark:border-white/5 flex flex-col gap-4">
                    <div class="flex w-full justify-start items-center gap-4">
                        <div
                            class="flex justify-center items-center size-10 rounded-2xl bg-[#F0B100]/10 border border-[#F0B100]/30">
                            <img src="/icons/agent-yellow.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-base dark:text-white">Agent-Alpha</span>
                            <span class="text-xs text-[#99a1af]">192.168.1.45</span>
                        </div>
                        <div
                            class="flex justify-center items-center gap-2 px-3 py-1.25 ms-auto mb-auto bg-[#F0B100]/20 border border-[#F0B100]/30 rounded-full">
                            <span class="size-2 bg-[#F0B100]"></span>
                            <span class="text-xs text-[#F0B100]">Warning</span>
                        </div>
                    </div>
                    <div class="w-full h-36.25 flex flex-col gap-3">
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">CPU</span>
                                <span class="dark:text-white">45%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">RAM</span>
                                <span class="dark:text-white">62%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div class="flex flex-col gap-1.5 justify-center">
                            <div class="flex justify-between items-start text-xs">
                                <span class="text-[#99a1af]">Storage</span>
                                <span class="dark:text-white">45%</span>
                            </div>
                            <div class="w-full flex gap-1 justify-center items-start">
                                {#each MACHINE.cpu.slice(-50) as detail}
                                    <div
                                        class="size-1.25 rounded-[1px] {detail.usage_percent
                                            ? detail.usage_percent < 65
                                                ? 'bg-green-700'
                                                : detail.usage_percent < 85
                                                  ? 'bg-[#F97316]'
                                                  : 'bg-[#EF4444]'
                                            : 'bg-[#FFFFFF]/5'}">
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <div
                            class="w-full pt-2 border-t border-[#0D0D0D]/5 dark:border-white/5 flex justify-between items-center text-[#99a1af] text-xs">
                            <div class="flex items-center justify-start gap-1">
                                <img src="/icons/time.svg" alt="time" />
                                <span>Uptime</span>
                            </div>
                            <span class="dark:text-white">14d 6h 23m</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="w-full grid grid-cols-4 gap-7.5">
            <div
                class="flex flex-col p-5 gap-3 bg-white dark:bg-[#0D0D0D] border border-[#0D0D0D]/5 dark:border-white/5 rounded-xl">
                <div class="flex items-center gap-3">
                    <span class="size-2 bg-[#22c55e] rounded-full"></span>
                    <span class="dark:text-white text-base">Health Center</span>
                </div>
                <div class="text-xs text-[#99a1af]">Self-help tools for troubleshooting network issues.</div>
                <div
                    class="w-full py-2.5 flex justify-center items-center gap-3 bg-[#e5e7eb] dark:bg-[#1a1c23] rounded-lg mt-1.5 cursor-pointer hover:bg-zinc-300 dark:hover:bg-[#22242e]">
                    <span class="text-sm dark:text-[#d1d5db]">Health Check</span>
                    <img
                        src={$theme === 'dark' ? '/icons/arrow-forward.svg' : '/icons/chevron-light.svg'}
                        alt="arrow forward" />
                </div>
            </div>
            <div
                class="flex flex-col p-5 gap-3 bg-white dark:bg-[#0D0D0D] border border-[#0D0D0D]/5 dark:border-white/5 rounded-xl">
                <div class="flex items-center gap-3">
                    <span class="size-2 bg-[#3b82f6] rounded-full"></span>
                    <span class="dark:text-white text-base">Performance</span>
                </div>
                <div class="text-xs text-[#99a1af]">Analyze latency and throughput bottlenecks.</div>
                <div
                    class="w-full py-2.5 flex justify-center items-center gap-3 bg-[#e5e7eb] dark:bg-[#1a1c23] rounded-lg mt-1.5 cursor-pointer hover:bg-zinc-300 dark:hover:bg-[#22242e]">
                    <span class="text-sm dark:text-[#d1d5db]">Analyze</span>
                    <img
                        src={$theme === 'dark' ? '/icons/arrow-forward.svg' : '/icons/chevron-light.svg'}
                        alt="arrow forward" />
                </div>
            </div>
            <div
                class="flex flex-col p-5 gap-3 bg-white dark:bg-[#0D0D0D] border border-[#0D0D0D]/5 dark:border-white/5 rounded-xl">
                <div class="flex items-center gap-3">
                    <span class="size-2 bg-[#a855f7] rounded-full"></span>
                    <span class="dark:text-white text-base">Logs</span>
                </div>
                <div class="text-xs text-[#99a1af]">View real-time system and security logs.</div>
                <div
                    class="w-full py-2.5 flex justify-center items-center gap-3 bg-[#e5e7eb] dark:bg-[#1a1c23] rounded-lg mt-1.5 cursor-pointer hover:bg-zinc-300 dark:hover:bg-[#22242e]">
                    <span class="text-sm dark:text-[#d1d5db]">View Logs</span>
                    <img
                        src={$theme === 'dark' ? '/icons/arrow-forward.svg' : '/icons/chevron-light.svg'}
                        alt="arrow forward" />
                </div>
            </div>
            <div
                class="flex flex-col p-5 gap-3 bg-white dark:bg-[#0D0D0D] border border-[#0D0D0D]/5 dark:border-white/5 rounded-xl">
                <div class="flex items-center gap-3">
                    <span class="size-2 bg-[#f97316] rounded-full"></span>
                    <span class="dark:text-white text-base">Alerts</span>
                </div>
                <div class="text-xs text-[#99a1af]">Configure notification policies and thresholds.</div>
                <div
                    class="w-full py-2.5 flex justify-center items-center gap-3 bg-[#e5e7eb] dark:bg-[#1a1c23] rounded-lg mt-1.5 cursor-pointer hover:bg-zinc-300 dark:hover:bg-[#22242e]">
                    <span class="text-sm dark:text-[#d1d5db]">Configure</span>
                    <img
                        src={$theme === 'dark' ? '/icons/arrow-forward.svg' : '/icons/chevron-light.svg'}
                        alt="arrow forward" />
                </div>
            </div>
        </div>
    </div>
</section>
