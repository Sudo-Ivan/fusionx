<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { updateFeed } from '$lib/api/feed';
	import type { Feed, Group } from '$lib/api/model';
	import { ChevronDown, ChevronRight } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import DraggableFeed from './DraggableFeed.svelte';

	interface Props {
		group: { id: number; name: string; feeds: (Feed & { indexInList: number })[] };
		isOpen: boolean;
		onToggle: () => void;
		isHighlightedGroup: (url: string) => boolean;
		isHighlightedFeed: (url: string) => boolean;
	}

	let { group, isOpen, onToggle, isHighlightedGroup, isHighlightedFeed }: Props = $props();

	let isDropTarget = $state(false);
	let dropElement: HTMLElement | undefined;

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		if (!e.dataTransfer) return;
		
		e.dataTransfer.dropEffect = 'move';
		isDropTarget = true;
	}

	function handleDragLeave(e: DragEvent) {
		if (!dropElement?.contains(e.relatedTarget as Node)) {
			isDropTarget = false;
		}
	}

	async function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDropTarget = false;
		
		if (!e.dataTransfer) return;
		
		try {
			const data = JSON.parse(e.dataTransfer.getData('text/plain'));
			const { feedId, sourceGroupId, feedName } = data;
			
			if (sourceGroupId === group.id) return;
			
			await updateFeed(feedId, { group_id: group.id });
			await invalidateAll();
			toast.success(`Moved ${feedName} to ${group.name}`);
		} catch (error) {
			console.error('Failed to move feed:', error);
			toast.error('Failed to move feed');
		}
	}
</script>

<li class="p-0">
	<div 
		bind:this={dropElement}
		class="gap-0 p-0 transition-all duration-200 rounded-md"
		class:bg-primary={isDropTarget}
		class:bg-opacity-20={isDropTarget}
		class:ring-2={isDropTarget}
		class:ring-primary={isDropTarget}
		ondragover={handleDragOver}
		ondragleave={handleDragLeave}
		ondrop={handleDrop}
	>
		<button
			class="btn btn-ghost btn-sm btn-square"
			onclick={(event) => {
				event.preventDefault();
				onToggle();
			}}
		>
			{#if isOpen}
				<ChevronDown class="size-4" />
			{:else}
				<ChevronRight class="size-4" />
			{/if}
		</button>
		<a
			href="/groups/{group.id}"
			class="line-clamp-1 block h-full grow place-content-center text-left"
			class:menu-active={isHighlightedGroup('/groups/' + group.id)}
		>
			{group.name}
		</a>
	</div>
	<ul class:hidden={!isOpen}>
		{#each group.feeds as feed}
			<DraggableFeed 
				{feed} 
				groupId={group.id}
				isHighlighted={isHighlightedFeed('/feeds/' + feed.id)}
			/>
		{/each}
	</ul>
</li>
