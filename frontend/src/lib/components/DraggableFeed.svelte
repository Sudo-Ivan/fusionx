<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { updateFeed } from '$lib/api/feed';
	import type { Feed } from '$lib/api/model';
	import { getFavicon } from '$lib/api/favicon';
	import { toast } from 'svelte-sonner';

	interface Props {
		feed: Feed & { indexInList: number };
		groupId: number;
		isHighlighted: boolean;
		onclick?: () => void;
	}

	let { feed, groupId, isHighlighted, onclick }: Props = $props();

	let isDragging = $state(false);
	let dragElement: HTMLElement | undefined;

	function handleDragStart(e: DragEvent) {
		if (!e.dataTransfer) return;
		
		isDragging = true;
		e.dataTransfer.setData('text/plain', JSON.stringify({
			feedId: feed.id,
			sourceGroupId: groupId,
			feedName: feed.name
		}));
		e.dataTransfer.effectAllowed = 'move';
		
		if (dragElement) {
			dragElement.style.opacity = '0.5';
		}
	}

	function handleDragEnd() {
		isDragging = false;
		if (dragElement) {
			dragElement.style.opacity = '1';
		}
	}

	async function handleMoveToGroup(targetGroupId: number) {
		if (targetGroupId === groupId) return;
		
		try {
			await updateFeed(feed.id, { group_id: targetGroupId });
			await invalidateAll();
			toast.success(`Moved ${feed.name} to new group`);
		} catch (error) {
			console.error('Failed to move feed:', error);
			toast.error('Failed to move feed');
		}
	}

	const textColor = $derived(
		feed.suspended
			? 'text-neutral-content/60'
			: feed.failure
				? 'text-error'
				: ''
	);
</script>

<li>
	<a
		bind:this={dragElement}
		id="sidebar-feed-{feed.indexInList}"
		data-group-id={groupId}
		href="/feeds/{feed.id}"
		class={`${isHighlighted ? 'menu-active' : ''} focus:ring-2 cursor-grab active:cursor-grabbing transition-all duration-200 hover:bg-base-200`}
		class:opacity-50={isDragging}
		class:scale-105={isDragging}
		draggable="true"
		ondragstart={handleDragStart}
		ondragend={handleDragEnd}
		{onclick}
	>
		<div class="avatar">
			<div class="size-4 rounded-full">
				<img src={getFavicon(feed.link)} alt={feed.name} loading="lazy" />
			</div>
		</div>
		<span class={`line-clamp-1 grow ${textColor}`}>{feed.name}</span>
		{#if feed.unread_count > 0}
			<span class="text-base-content/60 text-xs">{feed.unread_count}</span>
		{/if}
	</a>
</li>
