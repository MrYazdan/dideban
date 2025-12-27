<script>
    import Sidebar from '../../components/global/sidebar/index.svelte';
    import MainLayout from '../../components/layouts/MainLayout.svelte';
    import { theme } from '../../stores/theme.svelte';
    let { children } = $props();
</script>

<Sidebar />
<MainLayout>
    <div class="px-4 py-6">
        <div class="w-full h-full flex px-7">
            <div
                class="py-2 w-25.75 flex justify-center items-center gap-2 bg-[#0D0D0D]/5 dark:bg-white/5 border border-[#e5e7eb] dark:border-white/5 rounded-[14px]">
                <span class="text-xs text-[#99a1af]">
                    {$theme === 'dark' ? 'Dark' : 'Light'}
                </span>
                <button
                    aria-label="theme toggle"
                    onclick={() => theme.update(t => (t === 'dark' ? 'light' : 'dark'))}
                    class="w-11 h-6 bg-[#00bc7d]/20 border border-[#00bc7d]/30 rounded-full relative cursor-pointer">
                    <div
                        style="box-shadow: 0 5px 30px #00bc7d;"
                        class="absolute top-1/2 -translate-y-1/2 left-px size-5 rounded-full bg-[#00bc7d] transition-transform duration-300 ease-in-out dark:translate-x-5 translate-x-0">
                    </div>
                </button>
            </div>
            <div
                class="ms-3 size-10.5 p-2.5 rounded-[14px] dark:bg-white/5 bg-[#0D0D0D]/5 border border-[#e5e7eb] dark:border-white/5 flex justify-center items-center">
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

                <div class="h-full w-px bg-[#0D0D0D]/20 dark:bg-[#FFFFFF]/20"></div>

                <div class="flex gap-2 justify-center items-center">
                    <div class="flex flex-col justify-center items-end">
                        <span class="text-sm dark:text-white">Admin</span>

                        <span class="text-xs text-[#99A1AF]">System Administrator</span>
                    </div>

                    <div class="w-10.5 h-10 rounded-[10px] flex justify-center items-center bg-[#00b478]">
                        <img src="/icons/user.svg" alt="user" />
                    </div>
                </div>
            </div>
        </div>
    </div>
    {@render children()}
</MainLayout>
