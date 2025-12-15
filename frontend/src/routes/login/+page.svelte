<script>
    import FormInput from './../../components/common/FormInput.svelte';
    import { onMount } from 'svelte';
    import { fade, fly } from 'svelte/transition';
    import { emailRegex, phoneRegex } from '../../validators.svelte';

    let mount = $state(false);
    let username = $state(false);
    let password = $state(false);
    let error = $state({
        username: false,
        password: false,
    });

    let loading = $state(false);

    onMount(() => {
        mount = true;
    });

    function handleInputChange(e, type) {
        const value = e.target.value;

        if (type === 'username') {
            if (!value) {
                error.username = 'Please do not leave this field empty';
            } else if (emailRegex.test(value) || phoneRegex.test(value)) {
                error.username = null;
                username = value;
            } else {
                error.username = 'Enter a valid email or phone number';
            }
        }

        if (type === 'password') {
            if (!value) {
                error.password = 'Please enter your password';
            } else {
                error.password = null;
                password = value;
            }
        }
    }
</script>

{#if mount}
    <div
        in:fly={{ duration: 1000, y: -10 }}
        class="flex w-full h-full sm:h-fit md:max-w-[370px] flex-col justify-center items-start absolute top-1/2 -translate-1/2 start-1/2 md:border md:border-[#e5e5e5] sm:shadow-xl px-6 py-8 md:rounded-xl">
        <h1 in:fade={{ duration: 1400 }} class="text-2xl mb-2 tracking-wider">Login to your account</h1>

        <p in:fade={{ duration: 1500 }} class="text-gray-400 text-sm mb-7">Enter your email or phone number</p>

        <label class="w-full mb-4">
            <span class="text-sm cursor-pointer">Email or Phone</span>

            <FormInput
                type="text"
                placeholder="Enter your email or phone"
                onChange={e => handleInputChange(e, 'username')}
                error={error.username} />
        </label>

        <label class="w-full mb-6">
            <span class="text-sm cursor-pointer">Password</span>

            <FormInput
                type="password"
                placeholder="Enter your password"
                onChange={e => handleInputChange(e, 'password')}
                error={error.password} />
        </label>

        <button
            disabled={loading || error.username || error.password || !(username && password)}
            class="w-full h-10 flex justify-center items-center bg-black text-white font-semibold rounded-md cursor-pointer disabled:cursor-not-allowed hover:bg-black/70">
            {#if loading}
                <span class="animate-spin h-5 w-5 border-2 border-white border-e-transparent rounded-full"></span>
            {:else}
                Log In
            {/if}
        </button>
    </div>
{/if}
