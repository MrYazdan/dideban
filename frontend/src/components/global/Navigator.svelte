<script>
    import { page } from '$app/stores';
    import { onMount } from 'svelte';
    import { fly, fade } from 'svelte/transition';
    let isMenu = $state(false);
    let width = $state(0);
    let isMobile = $derived(width < 640);

    const MENU_ITEMS = [
        { name: 'Home', path: '/', iconHref: '/icons/home.png' },
        { name: 'Alerts', path: '/alerts', iconHref: '/icons/alert.png' },
        { name: 'Settings', path: '/setting', iconHref: '/icons/setting.png' },
    ];

    $effect(() => {
        if (!isMobile) {
            isMenu = false;
        }
    });
</script>

<svelte:window bind:innerWidth={width} />
<aside
    class="{isMenu
        ? 'flex w-full! h-full fixed start-0 z-50 xl:static xl:w-[255px]!'
        : 'hidden'} border-e bg-[#f5f5f5] xl:flex flex-col w-[255px]">
    <div class="h-[50px] w-full mx-auto border-b border-[#e5e5e5] flex justify-start items-center ps-4">
        <img width="22" height="22" src="/icons/monitoring.png" alt="monitoring" />
        <h3 class="ms-3 text-xl tracking-widest">Monitoring</h3>
    </div>

    <div class="px-4 py-6 flex flex-col justify-start items-start h-fit">
        <ul class="flex flex-col gap-4 text-base xl:text-sm">
            {#each MENU_ITEMS as item}
                <li>
                    <a
                        onclick={e => {
                            isMenu = false;
                            item.path === $page.url.pathname ? e.preventDefault() : null;
                        }}
                        class="flex justify-start gap-2 items-center {item.path === $page.url.pathname
                            ? 'font-bold'
                            : ''}"
                        href={item.path}>
                        <img width="18" height="18" src={item.iconHref} alt={item.name} />
                        <span>{item.name}</span></a>
                </li>
            {/each}
        </ul>
    </div>
</aside>

<button
    title="menu"
    type="button"
    onclick={() => {
        isMenu = !isMenu;
    }}
    class="fixed xl:hidden size-6 top-4 end-4 z-50">
    <div class="relative w-full h-full flex flex-col gap-1.5">
        <div
            class="{isMenu
                ? '-rotate-45 absolute top-1/2 -translate-y-1/2'
                : ''} h-0.5 w-full bg-black rounded-full transition-all">
        </div>
        <div class="{isMenu ? 'hidden' : ''} h-0.5 w-full bg-black rounded-full transition-all"></div>

        <div
            class="{isMenu
                ? 'rotate-45 absolute  top-1/2 -translate-y-1/2'
                : ''} h-0.5 w-full bg-black rounded-full transition-all">
        </div>
    </div>
</button>
