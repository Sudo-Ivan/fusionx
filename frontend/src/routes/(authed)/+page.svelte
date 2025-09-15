<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { refreshFeeds } from '$lib/api/feed';
	import ItemActionMarkAllasRead from '$lib/components/ItemActionMarkAllasRead.svelte';
	import AdaptiveItemLayout from '$lib/components/AdaptiveItemLayout.svelte';
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
				<div class="h-full flex flex-col">
					<div class="px-4 lg:px-8 py-6 flex-shrink-0">
						<h1 class="text-3xl font-bold">{t('common.unread')}</h1>
					</div>
					<div class="flex-1 px-4 lg:px-8 overflow-hidden">
						<AdaptiveItemLayout itemsData={data.items} highlightUnread={true} />
					</div>
				</div>
			{/snippet}
		</PullToRefresh>
	</div>
</div>
