<script>
    import './layout.css';
    let { children } = $props();
</script>

<svelte:head>
    <link rel="icon" href={'/icons/monitoring.png'} />
</svelte:head>

{@render children()}
