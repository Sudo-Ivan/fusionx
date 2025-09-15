<script lang="ts">
	import { beforeNavigate } from '$app/navigation';
	import FeedActionImport from '$lib/components/FeedActionImport.svelte';
	import ShortcutHelpModal from '$lib/components/ShortcutHelpModal.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import ReadingPane from '$lib/components/ReadingPane.svelte';
	import { globalState } from '$lib/state.svelte';
	import { getItem } from '$lib/api/item';
	import { page } from '$app/stores';
	import type { Item } from '$lib/api/model';

	let { children } = $props();
	let showSidebar = $state(false);
	let selectedItem = $state<Item | null>(null);
	let loading = $state(false);
	let lastLoadedItemId = $state<number | null>(null);

	// Check if we're on an item page and load it for 3-pane view
	$effect(() => {
		// Only track the specific values we need to avoid unnecessary re-runs
		const readingPaneMode = globalState.readingPaneMode;
		const pathname = $page?.url?.pathname;
		
		if (readingPaneMode === '3pane' && pathname) {
			const itemMatch = pathname.match(/\/items\/(\d+)/);
			if (itemMatch) {
				const itemId = parseInt(itemMatch[1], 10);
				if (itemId > 0 && itemId !== lastLoadedItemId && !loading) {
					loadItem(itemId);
				} else if (itemId <= 0) {
					selectedItem = null;
					lastLoadedItemId = null;
				}
			} else {
				selectedItem = null;
				lastLoadedItemId = null;
			}
		} else {
			selectedItem = null;
			lastLoadedItemId = null;
		}
	});

	async function loadItem(itemId: number) {
		if (loading || itemId === lastLoadedItemId) return;
		loading = true;
		lastLoadedItemId = itemId;
		try {
			selectedItem = await getItem(itemId);
		} catch (error) {
			console.error('Failed to load item:', error);
			selectedItem = null;
			lastLoadedItemId = null;
			// If item doesn't exist, navigate back to list view
			if ($page?.url?.pathname?.startsWith('/items/')) {
				history.back();
			}
		} finally {
			loading = false;
		}
	}

	function closeReadingPane() {
		selectedItem = null;
		lastLoadedItemId = null;
		// Navigate back to the current list view
		if ($page?.url?.pathname) {
			const currentPath = $page.url.pathname;
			if (currentPath.startsWith('/items/')) {
				history.back();
			}
		}
	}

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
		<div class="flex h-full">
			<!-- Main content area -->
			<div class={`flex flex-col pb-4 ${globalState.readingPaneMode === '3pane' && (selectedItem || loading) ? 'flex-1' : 'mx-auto max-w-6xl flex-1'}`}>
				<svelte:boundary>
					{@render children()}
					{#snippet failed(error, reset)}
						<p>{error}</p>
						<button onclick={reset} class="btn w-fit">oops! try again</button>
					{/snippet}
				</svelte:boundary>
			</div>
			
			<!-- Reading pane sidebar for 3-pane mode -->
			{#if globalState.readingPaneMode === '3pane' && (selectedItem || loading)}
				<div class="w-1/2 border-l border-base-300 bg-base-50">
					{#if loading}
						<div class="flex items-center justify-center h-full">
							<span class="loading loading-spinner loading-lg"></span>
						</div>
					{:else}
						<ReadingPane item={selectedItem} onClose={closeReadingPane} showCloseButton={true} />
					{/if}
				</div>
			{/if}
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
