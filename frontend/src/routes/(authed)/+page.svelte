<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { refreshFeeds } from '$lib/api/feed';
	import ItemActionMarkAllasRead from '$lib/components/ItemActionMarkAllasRead.svelte';
	import ItemList from '$lib/components/ItemList.svelte';
	import PageNavHeader from '$lib/components/PageNavHeader.svelte';
	import PullToRefresh from '$lib/components/PullToRefresh.svelte';
	import { t } from '$lib/i18n/index.js';
	import { toast } from 'svelte-sonner';

	let { data } = $props();

	async function handleRefresh() {
		try {
			await refreshFeeds({ all: true });
			await invalidateAll();
			toast.success('Refreshed successfully');
		} catch (error) {
			console.error('Refresh failed:', error);
			toast.error('Refresh failed');
		}
	}
</script>

<svelte:head>
	<title>{t('common.unread')}</title>
</svelte:head>

<div class="flex flex-col h-screen">
	<PageNavHeader showSearch={true}>
		{#await data.items}
			<ItemActionMarkAllasRead disabled />
		{:then items}
			<ItemActionMarkAllasRead items={items.items} />
		{/await}
	</PageNavHeader>
	<div class="flex-1 overflow-hidden">
		<PullToRefresh onrefresh={handleRefresh}>
			{#snippet children()}
				<div class="px-4 lg:px-8 h-full overflow-y-auto">
					<div class="py-6">
						<h1 class="text-3xl font-bold">{t('common.unread')}</h1>
					</div>
					<ItemList data={data.items} highlightUnread={true} />
				</div>
			{/snippet}
		</PullToRefresh>
	</div>
</div>
