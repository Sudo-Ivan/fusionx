<script lang="ts">
	import type { Item } from '$lib/api/model';
	import { globalState } from '$lib/state.svelte';
	import ItemList from './ItemList.svelte';
	import ThreePaneLayout from './ThreePaneLayout.svelte';
	import DrawerLayout from './DrawerLayout.svelte';
	
	interface Props {
		itemsData: Promise<{
			total: number;
			items: Item[];
		}>;
		highlightUnread?: boolean;
	}
	
	let { itemsData, highlightUnread = true }: Props = $props();
</script>

{#if globalState.readingPaneMode === '3pane'}
	<ThreePaneLayout {itemsData} {highlightUnread} />
{:else if globalState.readingPaneMode === 'drawer'}
	<DrawerLayout {itemsData} {highlightUnread} />
{:else}
	<!-- Default mode - use regular ItemList with navigation -->
	<ItemList data={itemsData} {highlightUnread} />
{/if}
