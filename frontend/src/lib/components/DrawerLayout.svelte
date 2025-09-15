<script lang="ts">
	import type { Item } from '$lib/api/model';
	import ItemList from './ItemList.svelte';
	import ReadingPane from './ReadingPane.svelte';
	import { getItem } from '$lib/api/item';
	import { page } from '$app/state';
	
	interface Props {
		itemsData: Promise<{
			total: number;
			items: Item[];
		}>;
		highlightUnread?: boolean;
	}
	
	let { itemsData, highlightUnread = true }: Props = $props();
	
	let selectedItem = $state<Item | null>(null);
	let loading = $state(false);
	let drawerOpen = $state(false);

	// Check if we're on an item page and preselect it
	$effect(() => {
		if ($page?.url?.pathname) {
			const url = $page.url;
			const itemMatch = url.pathname.match(/\/items\/(\d+)/);
			if (itemMatch) {
				const itemId = parseInt(itemMatch[1], 10);
				if (itemId > 0) {
					loadItem(itemId);
				} else {
					selectedItem = null;
					drawerOpen = false;
				}
			} else {
				selectedItem = null;
				drawerOpen = false;
			}
		} else {
			selectedItem = null;
			drawerOpen = false;
		}
	});

	async function loadItem(itemId: number) {
		if (loading) return;
		loading = true;
		try {
			selectedItem = await getItem(itemId);
			drawerOpen = true;
			// Update URL without full navigation
			const newUrl = `/items/${itemId}`;
			if ($page?.url?.pathname !== newUrl) {
				history.pushState({}, '', newUrl);
			}
		} catch (error) {
			console.error('Failed to load item:', error);
			selectedItem = null;
			drawerOpen = false;
			// If item doesn't exist, navigate back to list view
			if ($page?.url?.pathname?.startsWith('/items/')) {
				history.back();
			}
		} finally {
			loading = false;
		}
	}

	function handleItemClick(item: Item) {
		loadItem(item.id);
	}

	function closeDrawer() {
		drawerOpen = false;
		selectedItem = null;
		// Navigate back to list view
		if ($page?.url?.pathname) {
			const currentPath = $page.url.pathname;
			if (currentPath.startsWith('/items/')) {
				// Navigate to the parent route
				if (currentPath.includes('/feeds/')) {
					history.back();
				} else if (currentPath.includes('/groups/')) {
					history.back();
				} else {
					history.pushState({}, '', '/');
				}
			}
		}
	}

	// Handle escape key to close drawer
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && drawerOpen) {
			closeDrawer();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="relative h-full">
	<!-- Main content - Item List -->
	<div class="h-full">
		<ItemList 
			data={itemsData} 
			{highlightUnread}
			onItemClick={handleItemClick}
			selectedItemId={selectedItem?.id}
		/>
	</div>

	<!-- Drawer overlay -->
	{#if drawerOpen}
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div 
			class="fixed inset-0 bg-black/50 z-40"
			onclick={closeDrawer}
			onkeydown={(e) => {
				if (e.key === 'Enter' || e.key === ' ') {
					e.preventDefault();
					closeDrawer();
				}
			}}
		></div>
	{/if}

	<!-- Drawer -->
	<div class={`fixed top-0 right-0 h-full w-full sm:w-4/5 md:w-3/4 lg:w-3/5 xl:w-1/2 2xl:w-2/5 bg-base-100 shadow-xl z-50 transform transition-transform duration-300 ease-in-out ${drawerOpen ? 'translate-x-0' : 'translate-x-full'}`}>
		<!-- Resize handle -->
		<div class="absolute left-0 top-0 w-1 h-full bg-base-300 hover:bg-primary cursor-col-resize opacity-0 hover:opacity-100 transition-opacity"></div>
		{#if loading}
			<div class="flex items-center justify-center h-full">
				<span class="loading loading-spinner loading-lg"></span>
			</div>
		{:else}
			<ReadingPane item={selectedItem} onClose={closeDrawer} showCloseButton={true} />
		{/if}
	</div>
</div>
