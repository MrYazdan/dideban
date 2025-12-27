<script>
    import { page } from '$app/stores';
    import Support from './deps/Support.svelte';
</script>

<div
    class="fixed start-0 top-0 h-screen w-[256px] bg-white dark:bg-[#0D0D0D] flex flex-col border-e border-[#0D0D0D]/5 dark:border-[#222222]">
    <div class="p-6 h-25.75 flex justify-start items-center border-b border-[#0D0D0D]/5 dark:border-[#222222]">
        <img src="/icons/torbu.png" alt="turbu" />
        <div class="flex flex-col ms-3 justify-center">
            <h2 class="dark:text-white font-semibold tracking-wider">Dideban</h2>
            <span class="text-[#00bc7d] text-sm">V 1.0.0</span>
        </div>
    </div>

    <div class="p-4 flex flex-col gap-4">
        <span class="text-[#6a7282] text-sm">MAIN</span>
        <ul class="flex flex-col items-start justify-start gap-2">
            <li class="w-full rounded-[14px]">
                <a
                    class="flex w-full justify-start items-center gap-3 py-2 px-2.5 rounded-[14px] {$page.url
                        .pathname === '/'
                        ? 'outline outline-[#00d492]/15 outline-offset-1 bg-[#00d492]/10 cursor-default'
                        : 'cursor-pointer'}"
                    href="/">
                    <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path
                            d="M12.5 17.5V10.8333C12.5 10.6123 12.4122 10.4004 12.2559 10.2441C12.0996 10.0878 11.8877 10 11.6667 10H8.33333C8.11232 10 7.90036 10.0878 7.74408 10.2441C7.5878 10.4004 7.5 10.6123 7.5 10.8333V17.5"
                            stroke={$page.url.pathname === '/' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M2.5 8.33322C2.49994 8.09078 2.55278 7.85124 2.65482 7.63132C2.75687 7.4114 2.90566 7.21639 3.09083 7.05989L8.92417 2.06073C9.22499 1.80648 9.60613 1.66699 10 1.66699C10.3939 1.66699 10.775 1.80648 11.0758 2.06073L16.9092 7.05989C17.0943 7.21639 17.2431 7.4114 17.3452 7.63132C17.4472 7.85124 17.5001 8.09078 17.5 8.33322V15.8332C17.5 16.2753 17.3244 16.6992 17.0118 17.0117C16.6993 17.3243 16.2754 17.4999 15.8333 17.4999H4.16667C3.72464 17.4999 3.30072 17.3243 2.98816 17.0117C2.67559 16.6992 2.5 16.2753 2.5 15.8332V8.33322Z"
                            stroke={$page.url.pathname === '/' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                    </svg>
                    <span class="text-sm {$page.url.pathname === '/' ? 'text-[#00d492]' : 'text-[#99a1af]'}"
                        >Dashboard</span
                    ></a>
            </li>
            <li class="w-full">
                <a
                    href="/servers"
                    class="flex w-full justify-start items-center gap-3 py-2 px-2.5 rounded-[14px] {$page.url
                        .pathname === '/servers'
                        ? 'outline outline-[#00d492]/15 outline-offset-1 bg-[#00d492]/10 cursor-default'
                        : 'cursor-pointer'}">
                    <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path
                            d="M16.667 1.6665H3.33366C2.41318 1.6665 1.66699 2.4127 1.66699 3.33317V6.6665C1.66699 7.58698 2.41318 8.33317 3.33366 8.33317H16.667C17.5875 8.33317 18.3337 7.58698 18.3337 6.6665V3.33317C18.3337 2.4127 17.5875 1.6665 16.667 1.6665Z"
                            stroke={$page.url.pathname === '/servers' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M16.667 11.6665H3.33366C2.41318 11.6665 1.66699 12.4127 1.66699 13.3332V16.6665C1.66699 17.587 2.41318 18.3332 3.33366 18.3332H16.667C17.5875 18.3332 18.3337 17.587 18.3337 16.6665V13.3332C18.3337 12.4127 17.5875 11.6665 16.667 11.6665Z"
                            stroke={$page.url.pathname === '/servers' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M5 5H5.00833"
                            stroke={$page.url.pathname === '/servers' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M5 15H5.00833"
                            stroke={$page.url.pathname === '/servers' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                    </svg>
                    <span class="text-sm {$page.url.pathname === '/servers' ? 'text-[#00d492]' : 'text-[#99a1af]'}"
                        >Servers</span>
                </a>
            </li>
            <li class="w-full">
                <a
                    href="/network"
                    class="flex w-full justify-start items-center gap-3 py-2 px-2.5 rounded-[14px] {$page.url
                        .pathname === '/network'
                        ? 'outline outline-[#00d492]/15 outline-offset-1 bg-[#00d492]/10 cursor-default'
                        : 'cursor-pointer'}">
                    <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path
                            d="M17.4997 13.3335H14.1663C13.7061 13.3335 13.333 13.7066 13.333 14.1668V17.5002C13.333 17.9604 13.7061 18.3335 14.1663 18.3335H17.4997C17.9599 18.3335 18.333 17.9604 18.333 17.5002V14.1668C18.333 13.7066 17.9599 13.3335 17.4997 13.3335Z"
                            stroke={$page.url.pathname === '/servers' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M5.83366 13.3335H2.50033C2.04009 13.3335 1.66699 13.7066 1.66699 14.1668V17.5002C1.66699 17.9604 2.04009 18.3335 2.50033 18.3335H5.83366C6.2939 18.3335 6.66699 17.9604 6.66699 17.5002V14.1668C6.66699 13.7066 6.2939 13.3335 5.83366 13.3335Z"
                            stroke={$page.url.pathname === '/network' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M11.6667 1.6665H8.33333C7.8731 1.6665 7.5 2.0396 7.5 2.49984V5.83317C7.5 6.29341 7.8731 6.6665 8.33333 6.6665H11.6667C12.1269 6.6665 12.5 6.29341 12.5 5.83317V2.49984C12.5 2.0396 12.1269 1.6665 11.6667 1.6665Z"
                            stroke={$page.url.pathname === '/network' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M4.16699 13.3333V10.8333C4.16699 10.6123 4.25479 10.4004 4.41107 10.2441C4.56735 10.0878 4.77931 10 5.00033 10H15.0003C15.2213 10 15.4333 10.0878 15.5896 10.2441C15.7459 10.4004 15.8337 10.6123 15.8337 10.8333V13.3333"
                            stroke={$page.url.pathname === '/network' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M10 9.99984V6.6665"
                            stroke={$page.url.pathname === '/network' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                    </svg>
                    <span class="text-sm {$page.url.pathname === '/network' ? 'text-[#00d492]' : 'text-[#99a1af]'}"
                        >Network</span
                    ></a>
            </li>
            <li class="w-full">
                <a
                    href="/database"
                    class="flex w-full justify-start items-center gap-3 py-2 px-2.5 rounded-[14px] {$page.url
                        .pathname === '/database'
                        ? 'outline outline-[#00d492]/15 outline-offset-1 bg-[#00d492]/10 cursor-default'
                        : 'cursor-pointer'}">
                    <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path
                            d="M10 6.6665C14.1421 6.6665 17.5 5.54722 17.5 4.1665C17.5 2.78579 14.1421 1.6665 10 1.6665C5.85786 1.6665 2.5 2.78579 2.5 4.1665C2.5 5.54722 5.85786 6.6665 10 6.6665Z"
                            stroke={$page.url.pathname === '/database' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M2.5 4.1665V15.8332C2.5 16.4962 3.29018 17.1321 4.6967 17.6009C6.10322 18.0698 8.01088 18.3332 10 18.3332C11.9891 18.3332 13.8968 18.0698 15.3033 17.6009C16.7098 17.1321 17.5 16.4962 17.5 15.8332V4.1665"
                            stroke={$page.url.pathname === '/database' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M2.5 10C2.5 10.663 3.29018 11.2989 4.6967 11.7678C6.10322 12.2366 8.01088 12.5 10 12.5C11.9891 12.5 13.8968 12.2366 15.3033 11.7678C16.7098 11.2989 17.5 10.663 17.5 10"
                            stroke={$page.url.pathname === '/database' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                    </svg>
                    <span class="text-sm {$page.url.pathname === '/database' ? 'text-[#00d492]' : 'text-[#99a1af]'}"
                        >Database</span
                    ></a>
            </li>
            <li class="w-full">
                <a
                    href="/analytics"
                    class="flex w-full justify-start items-center gap-3 py-2 px-2.5 rounded-[14px] {$page.url
                        .pathname === '/analytics'
                        ? 'outline outline-[#00d492]/15 outline-offset-1 bg-[#00d492]/10 cursor-default'
                        : 'cursor-pointer'}">
                    <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path
                            d="M2.5 2.5V15.8333C2.5 16.2754 2.67559 16.6993 2.98816 17.0118C3.30072 17.3244 3.72464 17.5 4.16667 17.5H17.5"
                            stroke={$page.url.pathname === '/analytics' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M15 14.1667V7.5"
                            stroke={$page.url.pathname === '/analytics' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M10.833 14.1665V4.1665"
                            stroke={$page.url.pathname === '/analytics' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M6.66699 14.1665V11.6665"
                            stroke={$page.url.pathname === '/analytics' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                    </svg>
                    <span class="text-sm {$page.url.pathname === '/analytics' ? 'text-[#00d492]' : 'text-[#99a1af]'}"
                        >Analytics</span
                    ></a>
            </li>

            <li class="w-full">
                <a
                    href="/security"
                    class="flex w-full justify-start items-center gap-3 py-2 px-2.5 rounded-[14px] {$page.url
                        .pathname === '/security'
                        ? 'outline outline-[#00d492]/15 outline-offset-1 bg-[#00d492]/10 cursor-default'
                        : 'cursor-pointer'}">
                    <svg width="15" height="19" viewBox="0 0 15 19" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path
                            d="M14.1663 10.0005C14.1663 14.1672 11.2497 16.2505 7.78301 17.4589C7.60148 17.5204 7.40429 17.5174 7.22467 17.4505C3.74967 16.2505 0.833008 14.1672 0.833008 10.0005V4.1672C0.833008 3.94619 0.920805 3.73423 1.07709 3.57795C1.23337 3.42167 1.44533 3.33387 1.66634 3.33387C3.33301 3.33387 5.41634 2.33387 6.86634 1.0672C7.04289 0.91637 7.26747 0.833496 7.49967 0.833496C7.73188 0.833496 7.95646 0.91637 8.13301 1.0672C9.59134 2.3422 11.6663 3.33387 13.333 3.33387C13.554 3.33387 13.766 3.42167 13.9223 3.57795C14.0785 3.73423 14.1663 3.94619 14.1663 4.1672V10.0005Z"
                            stroke={$page.url.pathname === '/security' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                    </svg>
                    <span class="text-sm {$page.url.pathname === '/security' ? 'text-[#00d492]' : 'text-[#99a1af]'}"
                        >Security</span
                    ></a>
            </li>
            <li class="w-full">
                <a
                    href="/optimize"
                    class="flex w-full justify-start items-center gap-3 py-2 px-2.5 rounded-[14px] {$page.url
                        .pathname === '/optimize'
                        ? 'outline outline-[#00d492]/15 outline-offset-1 bg-[#00d492]/10 cursor-default'
                        : 'cursor-pointer'}">
                    <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path
                            d="M18.3337 9.99984H16.267C15.9028 9.99906 15.5484 10.1176 15.2579 10.3373C14.9675 10.557 14.757 10.8658 14.6587 11.2165L12.7003 18.1832C12.6877 18.2264 12.6614 18.2645 12.6253 18.2915C12.5893 18.3186 12.5454 18.3332 12.5003 18.3332C12.4552 18.3332 12.4114 18.3186 12.3753 18.2915C12.3393 18.2645 12.3129 18.2264 12.3003 18.1832L7.70033 1.8165C7.6877 1.77323 7.66139 1.73522 7.62533 1.70817C7.58926 1.68112 7.5454 1.6665 7.50033 1.6665C7.45525 1.6665 7.41139 1.68112 7.37533 1.70817C7.33926 1.73522 7.31295 1.77323 7.30033 1.8165L5.34199 8.78317C5.24405 9.13247 5.03481 9.44027 4.74605 9.65985C4.45728 9.87943 4.10476 9.9988 3.74199 9.99984H1.66699"
                            stroke={$page.url.pathname === '/optimize' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                    </svg>
                    <span class="text-sm {$page.url.pathname === '/optimize' ? 'text-[#00d492]' : 'text-[#99a1af]'}"
                        >Optimize</span>
                </a>
            </li>
            <li class="w-full">
                <a
                    href="/settings"
                    class="flex w-full justify-start items-center gap-3 py-2 px-2.5 rounded-[14px] {$page.url
                        .pathname === '/setting'
                        ? 'outline outline-[#00d492]/15 outline-offset-1 bg-[#00d492]/10 cursor-default'
                        : 'cursor-pointer'}">
                    <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path
                            d="M10.1833 1.6665H9.81667C9.37464 1.6665 8.95072 1.8421 8.63816 2.15466C8.3256 2.46722 8.15 2.89114 8.15 3.33317V3.48317C8.1497 3.77544 8.07255 4.0625 7.92628 4.31553C7.78002 4.56857 7.56978 4.7787 7.31667 4.92484L6.95834 5.13317C6.70497 5.27945 6.41756 5.35646 6.125 5.35646C5.83244 5.35646 5.54503 5.27945 5.29167 5.13317L5.16667 5.0665C4.78422 4.84589 4.32987 4.78604 3.90334 4.90009C3.47681 5.01415 3.11296 5.29278 2.89167 5.67484L2.70833 5.9915C2.48772 6.37395 2.42787 6.82831 2.54192 7.25484C2.65598 7.68137 2.93461 8.04521 3.31667 8.2665L3.44167 8.34984C3.69356 8.49526 3.90302 8.70408 4.04921 8.95553C4.1954 9.20698 4.27325 9.49231 4.275 9.78317V10.2082C4.27617 10.5019 4.19971 10.7906 4.05337 11.0453C3.90703 11.2999 3.69601 11.5113 3.44167 11.6582L3.31667 11.7332C2.93461 11.9545 2.65598 12.3183 2.54192 12.7448C2.42787 13.1714 2.48772 13.6257 2.70833 14.0082L2.89167 14.3248C3.11296 14.7069 3.47681 14.9855 3.90334 15.0996C4.32987 15.2136 4.78422 15.1538 5.16667 14.9332L5.29167 14.8665C5.54503 14.7202 5.83244 14.6432 6.125 14.6432C6.41756 14.6432 6.70497 14.7202 6.95834 14.8665L7.31667 15.0748C7.56978 15.221 7.78002 15.4311 7.92628 15.6841C8.07255 15.9372 8.1497 16.2242 8.15 16.5165V16.6665C8.15 17.1085 8.3256 17.5325 8.63816 17.845C8.95072 18.1576 9.37464 18.3332 9.81667 18.3332H10.1833C10.6254 18.3332 11.0493 18.1576 11.3618 17.845C11.6744 17.5325 11.85 17.1085 11.85 16.6665V16.5165C11.8503 16.2242 11.9275 15.9372 12.0737 15.6841C12.22 15.4311 12.4302 15.221 12.6833 15.0748L13.0417 14.8665C13.295 14.7202 13.5824 14.6432 13.875 14.6432C14.1676 14.6432 14.455 14.7202 14.7083 14.8665L14.8333 14.9332C15.2158 15.1538 15.6701 15.2136 16.0967 15.0996C16.5232 14.9855 16.887 14.7069 17.1083 14.3248L17.2917 13.9998C17.5123 13.6174 17.5721 13.163 17.4581 12.7365C17.344 12.31 17.0654 11.9461 16.6833 11.7248L16.5583 11.6582C16.304 11.5113 16.093 11.2999 15.9466 11.0453C15.8003 10.7906 15.7238 10.5019 15.725 10.2082V9.7915C15.7238 9.49782 15.8003 9.20904 15.9466 8.95441C16.093 8.69978 16.304 8.48834 16.5583 8.3415L16.6833 8.2665C17.0654 8.04521 17.344 7.68137 17.4581 7.25484C17.5721 6.82831 17.5123 6.37395 17.2917 5.9915L17.1083 5.67484C16.887 5.29278 16.5232 5.01415 16.0967 4.90009C15.6701 4.78604 15.2158 4.84589 14.8333 5.0665L14.7083 5.13317C14.455 5.27945 14.1676 5.35646 13.875 5.35646C13.5824 5.35646 13.295 5.27945 13.0417 5.13317L12.6833 4.92484C12.4302 4.7787 12.22 4.56857 12.0737 4.31553C11.9275 4.0625 11.8503 3.77544 11.85 3.48317V3.33317C11.85 2.89114 11.6744 2.46722 11.3618 2.15466C11.0493 1.8421 10.6254 1.6665 10.1833 1.6665Z"
                            stroke={$page.url.pathname === '/settings' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                        <path
                            d="M10 12.5C11.3807 12.5 12.5 11.3807 12.5 10C12.5 8.61929 11.3807 7.5 10 7.5C8.61929 7.5 7.5 8.61929 7.5 10C7.5 11.3807 8.61929 12.5 10 12.5Z"
                            stroke={$page.url.pathname === '/settings' ? '#00D492' : '#99a1af'}
                            stroke-width="1.66667"
                            stroke-linecap="round"
                            stroke-linejoin="round" />
                    </svg>
                    <span class="text-sm {$page.url.pathname === '/settings' ? 'text-[#00d492]' : 'text-[#99a1af]'}"
                        >Settings</span
                    ></a>
            </li>
        </ul>
    </div>

    <!-- Help section in slider-->
    <!-- <Support /> -->
</div>

