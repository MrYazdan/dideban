<script>
    import Chart from '../../components/common/Chart.svelte';
    import { MACHINE } from './machine/[machine]/constant.svelte';

    let width = $state(0);
    let isMobile = $derived(width < 365);
    let isTablet = $derived(width > 1279 && width < 1536);
    let visibleSeries = $state({
        CPU: true,
        Memory: true,
        Disk: true,
    });

    let theme = $state('dark');
</script>

<svelte:window bind:innerWidth={width} />

<section class="w-full h-auto flex flex-col">
    <!-- Head of dashboard page -->
    <div class="px-4 py-6">
        <div class="w-full h-full flex px-7">
            <div
                class="py-2 w-25.75 flex justify-center items-center gap-2 bg-white/5 border border-white/5 rounded-[14px]">
                <span class="text-xs text-[#99a1af]">Dark</span>
                <div class="w-11 h-6 bg-[#00bc7d]/20 border border-[#00bc7d]/30 rounded-full relative">
                    <div
                        style="box-shadow: 0 5px 30px #00bc7d;"
                        class="absolute top-1/2 -translate-y-1/2 end-px size-5 rounded-full bg-[#00bc7d]">
                    </div>
                </div>
            </div>
            <div
                class="ms-3 size-10.5 p-2.5 rounded-[14px] bg-white/5 border border-white/5 flex justify-center items-center">
                <div class="relative">
                    <div class="size-1 absolute -top-1 end-0 rounded-full bg-red-500"></div>
                    <img src="/icons/bell.png" alt="bell" />
                </div>
            </div>

            <div class="h-10 flex justify-between items-center ms-auto gap-4">
                <div class="flex flex-col gap-0.5">
                    <div
                        class="w-fit px-4 text-[#ff6467]/80 text-xs flex justify-center items-center gap-2 opacity-50 scale-90">
                        <div
                            style="box-shadow: 0 0 10px 1px #ff6467;"
                            class="hidden size-1.5 bg-[#ff6467] rounded-full">
                        </div>
                        <span
                            class="hover:animate-pulse hover:[text-shadow:0_0_100px_#F87171,0_0_150px_#F87171,0_0_200px_#F87171,0_0_250px_#F87171,0_0_300px_#F87171]">
                            Lorem ipsum dolor sit amet consectetur, adipisicing elit. Deserunt, ducimus?</span>
                    </div>
                    <div
                        class="w-fit bg-[#ff6467]/10 rounded-full text-[#ff6467]/80 text-xs flex justify-center items-center gap-2">
                        <div style="box-shadow: 0 0 10px 1px #ff6467;" class="rounded-full animate-pulse">
                            <img src="/icons/error.svg" alt="error" />
                        </div>
                        <span
                            class="hover:animate-pulse hover:[text-shadow:0_0_100px_#F87171,0_0_150px_#F87171,0_0_200px_#F87171,0_0_250px_#F87171,0_0_300px_#F87171]">
                            Lorem ipsum dolor sit amet consectetur, adipisicing elit. Deserunt, ducimus?</span>
                    </div>
                </div>

                <div class="h-full w-px bg-[#FFFFFF]/20"></div>

                <div class="flex gap-2 justify-center items-center">
                    <div class="flex flex-col justify-center items-end">
                        <span class="text-sm text-white">Admin</span>

                        <span class="text-xs text-[#99A1AF]">System Administrator</span>
                    </div>

                    <div class="w-10.5 h-10 rounded-[10px] flex justify-center items-center bg-[#00b478]">
                        <img src="/icons/user.svg" alt="user" />
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Content of dashboard page -->
    <div class="w-full flex flex-col gap-7.75 p-10 pt-3">
        <div class="w-full h-103.25 flex gap-7.75">
            <div class="w-213 h-full p-6 rounded-[14px] bg-[#0D0D0D] border border-white/5">
                <div class="flex flex-col gap-4 items-start justify-around">
                    <div class="w-full flex flex-col justify-start items-start">
                        <span class="text-xl text-white">Main Performance Overview</span>
                        <span class="text-sm text-[#99a1af]">System resource utilization trends</span>
                    </div>

                    <div class="w-full h-20.5 flex justify-around gap-4.75">
                        <div
                            class="h-full w-59.25 flex flex-col justify-start items-start gap-2 px-4 py-3 rounded-[10px] bg-[#121212] border border-white/5">
                            <div class="w-full flex justify-start items-center gap-1.5">
                                <span class="size-2 rounded-full bg-[#ad46ff]"></span>
                                <span class="text-sm text-[#6a7282]">CPU Average</span>
                            </div>

                            <span class="text-white text-2xl">60.1%</span>
                        </div>
                        <div
                            class="h-full w-59.25 flex flex-col justify-start items-start gap-2 px-4 py-3 rounded-[10px] bg-[#121212] border border-white/5">
                            <div class="w-full flex justify-start items-center gap-1.5">
                                <span class="size-2 rounded-full bg-[#2b7fff]"></span>
                                <span class="text-sm text-[#6a7282]">Memory Average</span>
                            </div>

                            <span class="text-white text-2xl">65.4%</span>
                        </div>
                        <div
                            class="h-full w-59.25 flex flex-col justify-start items-start gap-2 px-4 py-3 rounded-[10px] bg-[#121212] border border-white/5">
                            <div class="w-full flex justify-start items-center gap-1.5">
                                <span class="size-2 rounded-full bg-[#00bc7d]"></span>
                                <span class="text-sm text-[#6a7282]">Disk</span>
                            </div>

                            <span class="text-white text-2xl">50.6%</span>
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
                <div class="border border-white/5 rounded-[14px] relative overflow-hidden group">
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
                        <div class="text-white text-3xl mb-1">156</div>
                        <span class="text-sm text-[#99a1af]">Total Services</span>
                    </div>
                </div>
                <div class="border border-white/5 rounded-[14px] relative overflow-hidden group">
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
                <div class="border border-white/5 rounded-[14px] relative overflow-hidden group">
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
                <div class="border border-white/5 rounded-[14px] relative overflow-hidden group">
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

        <div class="w-full p-6 rounded-[14px] bg-[#0D0D0D] border border-white/5">
            <div class="w-full flex justify-between items-start pb-8">
                <div class="flex flex-col gap-1">
                    <span class="text-xl text-white">Machine Status</span>
                    <span class="text-sm text-[#99a1af]">Infrastructure Nodes & Cluster Health</span>
                </div>
                <div
                    class="w-12 h-10 flex justify-center items-center bg-[#22c55e]/10 rounded-lg text-xl text-[#10b981]">
                    +
                </div>
            </div>
            <div class="w-full grid grid-cols-3 gap-7.5">
                <div
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
        <div class="w-full p-6 rounded-[14px] bg-[#0D0D0D] border border-white/5">
            <div class="w-full flex justify-between items-start pb-8">
                <div class="flex flex-col gap-1">
                    <span class="text-xl text-white">Application Status</span>
                    <span class="text-sm text-[#99a1af]">Microservices Health & Availability</span>
                </div>
                <div
                    class="w-12 h-10 flex justify-center items-center bg-[#22c55e]/10 rounded-lg text-xl text-[#10b981]">
                    +
                </div>
            </div>
            <div class="w-full grid grid-cols-3 gap-7.5">
                <div
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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
                    class="w-122.75 h-37.75 border border-white/5 rounded-[14px] flex flex-col py-8 gap-7 bg-[#121212]">
                    <div class="flex justify-start items-center px-4.25 gap-4">
                        <div class="flex justify-center items-center size-12 rounded-2xl bg-[#6366f1]/10">
                            <img src="/icons/container.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-lg text-white">US-East-Cluster-01</span>
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

        <div class="w-full p-6 rounded-[14px] bg-[#0D0D0D] border border-white/5">
            <div class="w-full flex justify-between items-start pb-8">
                <div class="flex flex-col gap-1">
                    <span class="text-xl text-white">Agents Status</span>
                    <span class="text-sm text-[#99a1af]">Microservices Health & Availability</span>
                </div>
                <div
                    class="w-12 h-10 flex justify-center items-center bg-[#22c55e]/10 rounded-lg text-xl text-[#10b981]">
                    <img width="17" src="/icons/refresh-green.svg" alt="refresh" />
                </div>
            </div>

            <div class="w-full grid grid-cols-3 gap-7.5">
                <div class="w-full p-5 rounded-xl bg-[#121212] border border-white/5 flex flex-col gap-4">
                    <div class="flex w-full justify-start items-center gap-4">
                        <div
                            class="flex justify-center items-center size-10 rounded-2xl bg-[#00d492]/10 border border-[#00d492]/30">
                            <img src="/icons/agent-green.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-base text-white">Agent-Alpha</span>
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
                                <span class="text-white">45%</span>
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
                                <span class="text-white">62%</span>
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
                                <span class="text-white">45%</span>
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
                            class="w-full pt-2 border-t border-white/5 flex justify-between items-center text-[#99a1af] text-xs">
                            <div class="flex items-center justify-start gap-1">
                                <img src="/icons/time.svg" alt="time" />
                                <span>Uptime</span>
                            </div>
                            <span class="text-white">14d 6h 23m</span>
                        </div>
                    </div>
                </div>
                <div class="w-full p-5 rounded-xl bg-[#121212] border border-white/5 flex flex-col gap-4">
                    <div class="flex w-full justify-start items-center gap-4">
                        <div
                            class="flex justify-center items-center size-10 rounded-2xl bg-[#00d492]/10 border border-[#00d492]/30">
                            <img src="/icons/agent-green.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-base text-white">Agent-Alpha</span>
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
                                <span class="text-white">45%</span>
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
                                <span class="text-white">62%</span>
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
                                <span class="text-white">45%</span>
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
                            class="w-full pt-2 border-t border-white/5 flex justify-between items-center text-[#99a1af] text-xs">
                            <div class="flex items-center justify-start gap-1">
                                <img src="/icons/time.svg" alt="time" />
                                <span>Uptime</span>
                            </div>
                            <span class="text-white">14d 6h 23m</span>
                        </div>
                    </div>
                </div>
                <div class="w-full p-5 rounded-xl bg-[#121212] border border-white/5 flex flex-col gap-4">
                    <div class="flex w-full justify-start items-center gap-4">
                        <div
                            class="flex justify-center items-center size-10 rounded-2xl bg-[#00d492]/10 border border-[#00d492]/30">
                            <img src="/icons/agent-green.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-base text-white">Agent-Alpha</span>
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
                                <span class="text-white">45%</span>
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
                                <span class="text-white">62%</span>
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
                                <span class="text-white">45%</span>
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
                            class="w-full pt-2 border-t border-white/5 flex justify-between items-center text-[#99a1af] text-xs">
                            <div class="flex items-center justify-start gap-1">
                                <img src="/icons/time.svg" alt="time" />
                                <span>Uptime</span>
                            </div>
                            <span class="text-white">14d 6h 23m</span>
                        </div>
                    </div>
                </div>
                <div class="w-full p-5 rounded-xl bg-[#EF4444]/5 border border-white/5 flex flex-col gap-4">
                    <div class="flex w-full justify-start items-center gap-4">
                        <div
                            class="flex justify-center items-center size-10 rounded-2xl bg-[#EF4444]/10 border border-[#EF4444]/30">
                            <img src="/icons/agent-red.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-base text-white">Agent-Alpha</span>
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
                                <span class="text-white">45%</span>
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
                                <span class="text-white">62%</span>
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
                                <span class="text-white">45%</span>
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
                            class="w-full pt-2 border-t border-white/5 flex justify-between items-center text-[#99a1af] text-xs">
                            <div class="flex items-center justify-start gap-1">
                                <img src="/icons/time.svg" alt="time" />
                                <span>Uptime</span>
                            </div>
                            <span class="text-white">14d 6h 23m</span>
                        </div>
                    </div>
                </div>
                <div class="w-full p-5 rounded-xl bg-[#121212] border border-white/5 flex flex-col gap-4">
                    <div class="flex w-full justify-start items-center gap-4">
                        <div
                            class="flex justify-center items-center size-10 rounded-2xl bg-[#00d492]/10 border border-[#00d492]/30">
                            <img src="/icons/agent-green.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-base text-white">Agent-Alpha</span>
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
                                <span class="text-white">45%</span>
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
                                <span class="text-white">62%</span>
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
                                <span class="text-white">45%</span>
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
                            class="w-full pt-2 border-t border-white/5 flex justify-between items-center text-[#99a1af] text-xs">
                            <div class="flex items-center justify-start gap-1">
                                <img src="/icons/time.svg" alt="time" />
                                <span>Uptime</span>
                            </div>
                            <span class="text-white">14d 6h 23m</span>
                        </div>
                    </div>
                </div>
                <div class="w-full p-5 rounded-xl bg-[#F0B100]/5 border border-white/5 flex flex-col gap-4">
                    <div class="flex w-full justify-start items-center gap-4">
                        <div
                            class="flex justify-center items-center size-10 rounded-2xl bg-[#F0B100]/10 border border-[#F0B100]/30">
                            <img src="/icons/agent-yellow.svg" alt="container" />
                        </div>

                        <div class="flex flex-col justify-center items-start gap-1">
                            <span class="text-base text-white">Agent-Alpha</span>
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
                                <span class="text-white">45%</span>
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
                                <span class="text-white">62%</span>
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
                                <span class="text-white">45%</span>
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
                            class="w-full pt-2 border-t border-white/5 flex justify-between items-center text-[#99a1af] text-xs">
                            <div class="flex items-center justify-start gap-1">
                                <img src="/icons/time.svg" alt="time" />
                                <span>Uptime</span>
                            </div>
                            <span class="text-white">14d 6h 23m</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="w-full grid grid-cols-4 gap-7.5">
            <div class="flex flex-col p-5 gap-3 bg-[#0D0D0D] border border-white/5 rounded-xl">
                <div class="flex items-center gap-3">
                    <span class="size-2 bg-[#22c55e] rounded-full"></span>
                    <span class="text-white text-base">Health Center</span>
                </div>
                <div class="text-xs text-[#99a1af]">Self-help tools for troubleshooting network issues.</div>
                <div class="w-full py-2.5 flex justify-center items-center gap-3 bg-[#1a1c23] rounded-lg mt-1.5">
                    <span class="text-sm text-[#d1d5db]">Health Check</span>
                    <img src="/icons/arrow-forward.svg" alt="arrow forward" />
                </div>
            </div>
            <div class="flex flex-col p-5 gap-3 bg-[#0D0D0D] border border-white/5 rounded-xl">
                <div class="flex items-center gap-3">
                    <span class="size-2 bg-[#3b82f6] rounded-full"></span>
                    <span class="text-white text-base">Performance</span>
                </div>
                <div class="text-xs text-[#99a1af]">Analyze latency and throughput bottlenecks.</div>
                <div class="w-full py-2.5 flex justify-center items-center gap-3 bg-[#1a1c23] rounded-lg mt-1.5">
                    <span class="text-sm text-[#d1d5db]">Analyze</span>
                    <img src="/icons/arrow-forward.svg" alt="arrow forward" />
                </div>
            </div>
            <div class="flex flex-col p-5 gap-3 bg-[#0D0D0D] border border-white/5 rounded-xl">
                <div class="flex items-center gap-3">
                    <span class="size-2 bg-[#a855f7] rounded-full"></span>
                    <span class="text-white text-base">Logs</span>
                </div>
                <div class="text-xs text-[#99a1af]">View real-time system and security logs.</div>
                <div class="w-full py-2.5 flex justify-center items-center gap-3 bg-[#1a1c23] rounded-lg mt-1.5">
                    <span class="text-sm text-[#d1d5db]">View Logs</span>
                    <img src="/icons/arrow-forward.svg" alt="arrow forward" />
                </div>
            </div>
            <div class="flex flex-col p-5 gap-3 bg-[#0D0D0D] border border-white/5 rounded-xl">
                <div class="flex items-center gap-3">
                    <span class="size-2 bg-[#f97316] rounded-full"></span>
                    <span class="text-white text-base">Alerts</span>
                </div>
                <div class="text-xs text-[#99a1af]">Configure notification policies and thresholds.</div>
                <div class="w-full py-2.5 flex justify-center items-center gap-3 bg-[#1a1c23] rounded-lg mt-1.5">
                    <span class="text-sm text-[#d1d5db]">Configure</span>
                    <img src="/icons/arrow-forward.svg" alt="arrow forward" />
                </div>
            </div>
        </div>
    </div>
</section>
