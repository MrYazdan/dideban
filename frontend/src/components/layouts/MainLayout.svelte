<script>
    import Navigator from '../global/Navigator.svelte';

    let { children } = $props();
</script>

<main
    class="w-full h-full fixed top-1/2 -translate-y-1/2 start-1/2 -translate-x-1/2 2xl:rounded-2xl 2xl:border *:border-[#e5e5e5] border-[#e5e5e5] flex overflow-hidden">
    <Navigator />
    <section class="relative h-auto w-full px-4 overflow-y-auto overflow-x-hidden custom-scroll mt-12 xl:mt-0">
        {@render children()}
    </section>
</main>

<style>
    .custom-scroll {
        position: relative;
    }

    .custom-scroll::-webkit-scrollbar {
        width: 2px;
        position: absolute;
    }

    .custom-scroll::-webkit-scrollbar-track {
        background: transparent;
    }

    .custom-scroll::-webkit-scrollbar-thumb {
        background: #a79f9f;
        border-radius: 9999px;
    }
</style>
