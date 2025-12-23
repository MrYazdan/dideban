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

            <div class="w-80.75 h-10 flex justify-between items-center ms-auto">
                <div class="w-31.5 h-10 rounded-lg bg-[#10b981] flex justify-center items-center gap-2">
                    <img src="/icons/plus-mark.png" alt="plus mark" />

                    <span class="text-white text-sm">Add Agent</span>
                </div>
                <div class="h-full w-px bg-[#FFFFFF]/20"></div>

                <div class="flex gap-2 justify-center items-center">
                    <div class="flex flex-col justify-center items-end">
                        <span class="text-sm text-white">Andrew Smith</span>

                        <span class="text-xs text-[#99A1AF]">System Admin</span>
                    </div>

                    <div class="w-10.5 h-10 rounded-[10px] flex justify-center items-center bg-[#00b478]">
                        <img src="/icons/user.svg" alt="user" />
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Content of dashboard page -->

    <div class="w-full gap-11.5 p-10 pt-3">
        <div class="w-full h-103.25 flex gap-7.75">
            <div class="w-213 h-full p-6 rounded-[14px] bg-[#0D0D0D] border border-white/5">
                <div class="flex flex-col gap-4 items-start justify-around">
                    <div class="w-full flex flex-col justify-start items-start">
                        <span class="text-xl text-white">24-Hour Performance Overview</span>
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
                                <span class="text-sm text-[#6a7282]">Network Average</span>
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
                                    ? MACHINE.cpu.slice(-15).map(d => d.usage_percent ?? 0)
                                    : MACHINE.cpu.map(d => d.usage_percent ?? 0),
                            },
                            {
                                name: 'Memory',
                                data: isMobile
                                    ? MACHINE.memory.slice(-15).map(d => d.usage_percent ?? 0)
                                    : MACHINE.memory.map(d => d.usage_percent ?? 0),
                            },
                            {
                                name: 'Disk',
                                data: isMobile
                                    ? MACHINE.disk.slice(-15).map(d => d.usage_percent ?? 0)
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
    </div>
</section>
