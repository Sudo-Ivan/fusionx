<script lang="ts">
	import type { Item } from '$lib/api/model';
	import ItemActionBookmark from './ItemActionBookmark.svelte';
	import ItemActionGotoFeed from './ItemActionGotoFeed.svelte';
	import ItemActionUnread from './ItemActionUnread.svelte';
	import ItemActionVisitLink from './ItemActionVisitLink.svelte';
	import ItemActionShareLink from './ItemActionShareLink.svelte';
	import { render } from '$lib/render-item';
	import { ExternalLink, X } from 'lucide-svelte';

	interface Props {
		item: Item | null;
		onClose?: () => void;
		showCloseButton?: boolean;
	}
	
	let { item, onClose, showCloseButton = false }: Props = $props();

	let safeContent = $derived(item ? render(item.content, item.link) : '');
</script>

{#if item}
	<div class="flex flex-col h-full">
		<!-- Header with actions -->
		<div class="flex items-center justify-between p-4 border-b border-base-300">
			<div class="flex items-center gap-2">
				<ItemActionGotoFeed {item} />
				<ItemActionUnread bind:item />
				<ItemActionBookmark bind:item />
				<ItemActionVisitLink {item} />
				<ItemActionShareLink {item} />
			</div>
			{#if showCloseButton && onClose}
				<button onclick={onClose} class="btn btn-ghost btn-sm btn-circle">
					<X class="size-4" />
				</button>
			{/if}
		</div>

		<!-- Content -->
		<div class="flex-1 overflow-y-auto p-4">
			<article class="max-w-none">
				<div class="space-y-2 pb-6">
					<h1 class="text-2xl font-bold leading-tight">
						<a
							href={item.link}
							target="_blank"
							class="inline-flex items-center gap-2 no-underline hover:underline"
						>
							<span>{item.title || item.link}</span>
							<ExternalLink class="hidden size-4 md:block flex-shrink-0" />
						</a>
					</h1>
					<a 
						href={'/feeds/' + item.feed.id} 
						class="text-base-content/60 text-sm hover:underline block"
					>
						{item.feed.name} | {new Date(item.pub_date).toLocaleString()}
					</a>
				</div>
				<div class="prose max-w-none text-wrap break-words">
					{@html safeContent}
				</div>
			</article>
		</div>
	</div>
{:else}
	<div class="flex items-center justify-center h-full text-base-content/60">
		<div class="text-center">
			<p class="text-lg">Select an item to read</p>
			<p class="text-sm mt-2">Choose an article from the list to view it here</p>
		</div>
	</div>
{/if}
