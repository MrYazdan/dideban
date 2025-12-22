<script>
    import { fade } from 'svelte/transition';
    import { page } from '$app/stores';
    import { onMount } from 'svelte';

    const error = $page.error;
    const status = $page.status;
    let mount = $state(false);
    onMount(() => {
        mount = true;
    });
</script>

{#if mount}
    <div
        in:fade={{ duration: 1000 }}
        class="fixed top-1/2 start-1/2 -translate-1/2 flex flex-col justify-center items-center [&>p]:text-base [&>p]:tracking-wide [&>p]:text-gray-500">
        <span class="text-7xl font-bold tracking-widest mb-4 pb-4 border-b border-[#bdbdbd] w-full text-center">
            {status}</span>
        {#if status === 404}
            <p class="mb-1">Oops, it looks like the page you're looking for</p>
            <p>doesn't exist.</p>
        {:else}
            <p class="mb-1">Oops, we couldn’t process your request at this time.</p>
            <p>Please try again later</p>
        {/if}

        <a
            href="/"
            class="w-2/3 mt-10 h-10 flex justify-center items-center bg-black text-white font-semibold rounded-md cursor-pointer hover:bg-black/70">
            Go to home page
        </a>
    </div>
{/if}
