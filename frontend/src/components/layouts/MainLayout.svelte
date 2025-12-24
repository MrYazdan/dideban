<script>
    let { children } = $props();
</script>

<main class="fixed end-0 top-0 w-416 h-screen overflow-y-auto custom-scroll">
    {@render children()}
</main>

<style>
    .custom-scroll::-webkit-scrollbar {
        width: 3px;
    }

    .custom-scroll::-webkit-scrollbar-track {
        background: transparent;
    }

    .custom-scroll::-webkit-scrollbar-thumb {
        background: #99a1af;
        border-radius: 9999px;
    }
</style>
