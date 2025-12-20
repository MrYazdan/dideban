<script>
    const { value = 0, height = 8, showValue = true } = $props();

    const safeValue = Math.min(Math.max(value, 0), 100);

    function getColor(val) {
        if (val < 65) return 'bg-green-500';
        if (val <= 85) return 'bg-yellow-500';
        return 'bg-red-500';
    }
</script>

<div class="w-full space-y-1">
    <!-- bar -->
    <div class="w-full bg-gray-200 rounded-full overflow-hidden" style="height: {height}px">
        <div
            class="h-full rounded-full transition-all duration-500 ease-out {getColor(safeValue)}"
            style="width: {safeValue}%">
        </div>
    </div>

    <!-- label -->
    {#if showValue}
        <div
            class="text-sm font-medium text-right"
            class:text-green-600={safeValue < 65}
            class:text-yellow-600={safeValue >= 65 && safeValue <= 85}
            class:text-red-600={safeValue > 85}>
            {safeValue}%
        </div>
    {/if}
</div>
