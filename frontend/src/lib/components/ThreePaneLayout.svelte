<script lang="ts">
	import type { Item } from '$lib/api/model';
	import ItemList from './ItemList.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	
	interface Props {
		itemsData: Promise<{
			total: number;
			items: Item[];
		}>;
		highlightUnread?: boolean;
	}
	
	let { itemsData, highlightUnread = true }: Props = $props();

	function handleItemClick(item: Item) {
		// Navigate to the item URL - the layout will handle showing it in the reading pane
		goto(`/items/${item.id}`);
	}

	// Get current selected item ID from URL for highlighting
	let selectedItemId = $state<number | undefined>(undefined);
	
	$effect(() => {
		if ($page?.url?.pathname) {
			const url = $page.url;
			const itemMatch = url.pathname.match(/\/items\/(\d+)/);
			selectedItemId = itemMatch ? parseInt(itemMatch[1], 10) : undefined;
		} else {
			selectedItemId = undefined;
		}
	});
</script>

<div class="h-full">
	<ItemList 
		data={itemsData} 
		{highlightUnread}
		onItemClick={handleItemClick}
		selectedItemId={selectedItemId}
	/>
</div>
