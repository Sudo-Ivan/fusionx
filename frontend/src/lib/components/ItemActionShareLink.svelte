<script lang="ts">
	import type { Item } from '$lib/api/model';
	import { t } from '$lib/i18n';
	import { Share2 } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	interface Props {
		item: Item;
	}

	let { item }: Props = $props();

	const isShareSupported = !!navigator?.share;

	async function shareItem() {
		try {
			if (isShareSupported) {
				await navigator.share({
					title: item.title,
					url: item.link
				});
			} else {
				await navigator.clipboard.writeText(item.link);
				toast.success(t('item.link_copied'));
			}
		} catch (e) {
			toast.error((e as Error).message);
		}
	}
</script>

<div class="tooltip tooltip-bottom" data-tip={t('item.share')}>
	<button class="btn btn-ghost btn-square" onclick={shareItem}>
		<Share2 class="size-4" />
	</button>
</div>
