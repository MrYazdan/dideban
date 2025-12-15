<script>
    import { fade } from 'svelte/transition';
    import { MACHINES } from '../constant.svelte';
    import formatTimeAgo from '../../utils/formatFullDate';
    import formatFullDate from '../../utils/formatFullDate';
    let width = $state(0);
    let itemHoveredDetail = $state(null);
    let activeBar = $state(null);
    let isMobile = $derived(width < 640);

    $effect(() => {
        const media = window.matchMedia('(max-width: 640px)');
        const handler = e => (isMobile = e.matches);

        media.addEventListener('change', handler);
        return () => media.removeEventListener('change', handler);
    });
</script>

<svelte:window bind:innerWidth={width} />

<div class="w-full h-auto flex flex-col">
    <h2 class="h-[50px]! text-2xl xl:text-3xl w-full flex items-center border-b border-[#e5e5e5] mb-4">Machines</h2>

    {#each MACHINES as item (item.name)}
        <div class="w-full h-full pb-8">
            <h2 class="text-2xl mb-4"># {item.name}</h2>
            <div class="w-full flex flex-col gap-5">
                {#each item.stats as stat (stat.name)}
                    <div
                        class="relative flex flex-col lg:flex-row h-[110px] rounded-lg md:border md:border-[#e5e5e5] md:px-5 md:py-3">
                        <div class="absolute top-3.5 end-0 md:end-5 md:top-2 flex gap-1 justify-start items-start">
                            <div class="text-gray-500 text-xs flex items-baseline gap-2">
                                <div
                                    class="size-2.5 rounded-full {stat.detail[stat.detail.length - 1]?.loaded < 65
                                        ? 'bg-green-700'
                                        : stat.detail[stat.detail.length - 1]?.loaded < 85
                                          ? 'bg-yellow-500'
                                          : 'bg-red-600'}">
                                </div>
                                <span style="word-spacing: -2px;">{formatTimeAgo(stat.updateAt)}</span>
                            </div>
                        </div>

                        <h4
                            class="justify-start my-auto text-xl h-fit border-s-3 ps-2 w-[120px] {stat.detail[
                                stat.detail.length - 1
                            ]?.loaded < 65
                                ? 'border-s-green-700'
                                : stat.detail[stat.detail.length - 1]?.loaded < 85
                                  ? 'border-s-yellow-500'
                                  : 'border-s-red-600'} flex justify-center items-center">
                            {stat.name}
                        </h4>

                        <div
                            class="absolute lg:static lg:w-full top-27 sm:top-3 max-lg:start-1/2 max-lg:-translate-x-1/2 lg:start-0 justify-center items-center text-xs sm:text-sm flex lg:flex-col gap-7 lg:gap-1">
                            {#if itemHoveredDetail?.id === item?.id && itemHoveredDetail?.name === stat?.name}
                                <div class="flex [&>div]:text-nowrap">
                                    <div in:fade={{ duration: 500 }}>Total</div>
                                    <div in:fade={{ duration: 500 }}>: {stat.total}</div>
                                </div>
                                <div class="flex [&>div]:text-nowrap">
                                    <div in:fade={{ duration: 1000 }}>
                                        {itemHoveredDetail?.status?.usage ? 'Usage' : null}
                                    </div>
                                    <div in:fade={{ duration: 1000 }}>
                                        {itemHoveredDetail?.status?.usage
                                            ? ': ' + itemHoveredDetail?.status?.usage
                                            : null}
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

                        <div class="w-full lg:w-fit ms-auto flex flex-col justify-center gap-3 my-auto">
                            <div class="flex gap-1 justify-center items-center w-full py-1 mx-auto">
                                {#each isMobile ? stat.detail.slice(-28) : stat.detail as status}
                                    <div
                                        role="presentation"
                                        onclick={() => {
                                            activeBar = activeBar === status ? null : status;

                                            itemHoveredDetail = activeBar
                                                ? { status: { ...status }, id: item.id, name: stat.name }
                                                : null;
                                        }}
                                        onmouseenter={() =>
                                            (itemHoveredDetail = {
                                                status: { ...status },
                                                id: item.id,
                                                name: stat.name,
                                            })}
                                        onmouseleave={() => (itemHoveredDetail = null)}
                                        class="h-10 w-[2%] lg:w-2 min-w-2 rounded-full hover:scale-110 hover:opacity-90 transition-all cursor-pointer relative group {status.loaded
                                            ? status.loaded < 65
                                                ? 'bg-green-700'
                                                : status.loaded < 85
                                                  ? 'bg-yellow-500'
                                                  : 'bg-red-700'
                                            : 'bg-black/20'}">
                                        <div
                                            style="word-spacing: -2px"
                                            class="opacity-0 group-hover:opacity-100 transition-opacity absolute top-10 lg:top-12 start-1/2 -translate-x-1/2 text-gray-400 text-xs text-nowrap">
                                            {formatFullDate(stat.updateAt)}
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        </div>
                    </div>{/each}
            </div>
            <a title="{item.name} detail" href="/machines/{item.name.toLocaleLowerCase()}" class="">
                <img width="30" class="rtl:rotate-180 ms-auto mt-4" src="/icons/direction.png" alt="direction" />
            </a>
        </div>{/each}
</div>
