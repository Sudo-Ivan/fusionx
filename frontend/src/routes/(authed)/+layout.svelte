<script lang="ts">
	import { beforeNavigate } from '$app/navigation';
	import FeedActionImport from '$lib/components/FeedActionImport.svelte';
	import ShortcutHelpModal from '$lib/components/ShortcutHelpModal.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import { globalState } from '$lib/state.svelte';

	let { children } = $props();
	let showSidebar = $state(false);
	beforeNavigate(() => {
		showSidebar = false;
	});
</script>

<div class="drawer lg:drawer-open">
	<input id="sidebar-toggle" type="checkbox" bind:checked={showSidebar} class="drawer-toggle" />
	<div class="drawer-content bg-base-100 relative z-10 min-h-screen overflow-x-clip">
		{#if globalState.demoMode}
			<div class="alert alert-info shadow-lg">
				<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
				<span><strong>Demo Mode:</strong> This is a read-only demonstration. You cannot add, edit, or delete feeds.</span>
			</div>
		{/if}
		<div class="mx-auto flex h-full max-w-6xl flex-col pb-4">
			<svelte:boundary>
				{@render children()}
				{#snippet failed(error, reset)}
					<p>{error}</p>
					<button onclick={reset} class="btn w-fit">oops! try again</button>
				{/snippet}
			</svelte:boundary>
		</div>
	</div>
	<div class="drawer-side z-10">
		<label for="sidebar-toggle" aria-label="close sidebar" class="drawer-overlay"></label>
		<div
			class="text-base-content bg-base-200 z-50 h-full min-h-full w-[80%] overflow-x-hidden px-2 py-4 lg:w-72"
		>
			<Sidebar />
		</div>
	</div>
</div>

<!-- put these outside the drawer because when its inner modal is placed inside the drawer sidebar, the underlying dialog won't close properly -->
<FeedActionImport />
<ShortcutHelpModal />
