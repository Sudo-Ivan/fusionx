<script lang="ts">
	interface Props {
		disabled?: boolean;
		threshold?: number;
		resistance?: number;
		onrefresh?: () => void | Promise<void>;
		children?: any;
	}

	let { disabled = false, threshold = 60, resistance = 2.5, onrefresh, children }: Props = $props();

	let container: HTMLDivElement;
	let pullDistance = $state(0);
	let isRefreshing = $state(false);
	let startY = 0;
	let isDragging = $state(false);

	function handleTouchStart(e: TouchEvent) {
		if (disabled || isRefreshing) return;
		
		const scrollTop = container.scrollTop;
		if (scrollTop > 0) return;
		
		startY = e.touches[0].clientY;
		isDragging = true;
	}

	function handleTouchMove(e: TouchEvent) {
		if (!isDragging || disabled || isRefreshing) return;
		
		const currentY = e.touches[0].clientY;
		const deltaY = currentY - startY;
		
		if (deltaY > 0) {
			e.preventDefault();
			pullDistance = Math.min(deltaY / resistance, threshold * 1.5);
		}
	}

	async function handleTouchEnd() {
		if (!isDragging || disabled || isRefreshing) return;
		
		isDragging = false;
		
		if (pullDistance >= threshold) {
			isRefreshing = true;
			if (onrefresh) {
				await onrefresh();
			}
			
			setTimeout(() => {
				isRefreshing = false;
				pullDistance = 0;
			}, 1500);
		} else {
			pullDistance = 0;
		}
	}

	function handleMouseDown(e: MouseEvent) {
		if (disabled || isRefreshing) return;
		
		const scrollTop = container.scrollTop;
		if (scrollTop > 0) return;
		
		startY = e.clientY;
		isDragging = true;
	}

	function handleMouseMove(e: MouseEvent) {
		if (!isDragging || disabled || isRefreshing) return;
		
		const currentY = e.clientY;
		const deltaY = currentY - startY;
		
		if (deltaY > 0) {
			e.preventDefault();
			pullDistance = Math.min(deltaY / resistance, threshold * 1.5);
		}
	}

	async function handleMouseUp() {
		if (!isDragging || disabled || isRefreshing) return;
		
		isDragging = false;
		
		if (pullDistance >= threshold) {
			isRefreshing = true;
			if (onrefresh) {
				await onrefresh();
			}
			
			setTimeout(() => {
				isRefreshing = false;
				pullDistance = 0;
			}, 1500);
		} else {
			pullDistance = 0;
		}
	}

	$effect(() => {
		if (typeof document !== 'undefined') {
			document.addEventListener('mousemove', handleMouseMove);
			document.addEventListener('mouseup', handleMouseUp);
			
			return () => {
				document.removeEventListener('mousemove', handleMouseMove);
				document.removeEventListener('mouseup', handleMouseUp);
			};
		}
	});
</script>

<div 
	bind:this={container}
	class="pull-to-refresh-container"
	style:transform={`translateY(${pullDistance}px)`}
	style:transition={isDragging ? 'none' : 'transform 0.3s ease'}
	role="application"
	aria-label="Pull to refresh content"
	tabindex="0"
	onmousedown={handleMouseDown}
	ontouchstart={handleTouchStart}
	ontouchmove={handleTouchMove}
	ontouchend={handleTouchEnd}
>
	{#if pullDistance > 0 || isRefreshing}
		<div 
			class="pull-indicator"
			style:opacity={Math.min(pullDistance / threshold, 1)}
		>
			<div class="flex items-center justify-center gap-2 py-4">
				{#if isRefreshing}
					<div class="loading loading-spinner loading-sm"></div>
					<span class="text-sm">Refreshing...</span>
				{:else if pullDistance >= threshold}
					<span class="text-sm">Release to refresh</span>
				{:else}
					<span class="text-sm">Pull to refresh</span>
				{/if}
			</div>
		</div>
	{/if}
	
	{@render children?.()}
</div>

<style>
	.pull-to-refresh-container {
		overflow-y: auto;
		height: 100%;
		-webkit-overflow-scrolling: touch;
	}
	
	.pull-indicator {
		position: absolute;
		top: -60px;
		left: 0;
		right: 0;
		height: 60px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: hsl(var(--b1));
		border-bottom: 1px solid hsl(var(--b3));
		z-index: 10;
	}
</style>
