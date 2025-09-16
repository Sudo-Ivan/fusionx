<script lang="ts">
	import AdaptiveItemLayout from '$lib/components/AdaptiveItemLayout.svelte';
	import PageNavHeader from '$lib/components/PageNavHeader.svelte';
	import { t } from '$lib/i18n';

	let { data } = $props();
</script>

<svelte:head>
	<title>{t('common.all')}</title>
</svelte:head>

<div class="flex flex-col">
	<PageNavHeader showSearch={true}></PageNavHeader>
	<div class="px-2 sm:px-4 lg:px-8">
		<div class="py-6">
			<h1 class="text-2xl sm:text-3xl font-bold">{t('common.all')}</h1>
		</div>
		<AdaptiveItemLayout itemsData={data.items} highlightUnread={true} />
	</div>
</div>
