<script>
    import Sidebar from './../../components/global/Sidebar.svelte';
    import MainLayout from '../../components/layouts/MainLayout.svelte';
    let { children } = $props();
</script>

<Sidebar />
<MainLayout>
    {@render children()}
</MainLayout>
