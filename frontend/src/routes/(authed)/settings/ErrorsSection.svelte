<script lang="ts">
	import { getFeedErrors, type FeedError } from '$lib/api/errors';
	import { refreshFeeds, deleteFeed } from '$lib/api/feed';
	import { getFavicon } from '$lib/api/favicon';
	import { globalState } from '$lib/state.svelte';
	import { t } from '$lib/i18n';
	import { AlertTriangle, ExternalLink, RefreshCw, Clock, Trash2 } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import Section from './Section.svelte';

	let errors = $state<FeedError[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	async function loadErrors() {
		try {
			loading = true;
			errors = await getFeedErrors();
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	onMount(loadErrors);

	async function retryFeed(feedId: number) {
		try {
			await refreshFeeds({ id: feedId });
			toast.success('Feed refresh initiated');
			setTimeout(loadErrors, 2000);
		} catch (e) {
			toast.error('Failed to refresh feed');
		}
	}

	async function retryAllFailed() {
		if (errors.length === 0) return;
		
		try {
			await refreshFeeds({ all: true });
			toast.success('All feeds refresh initiated');
			setTimeout(loadErrors, 2000);
		} catch (e) {
			toast.error('Failed to refresh feeds');
		}
	}

	async function removeFeed(feedId: number, feedName: string) {
		if (!confirm(`Remove feed "${feedName}"? This will permanently delete the feed and all its items.`)) return;
		
		try {
			await deleteFeed(feedId);
			toast.success(`Feed "${feedName}" removed successfully`);
			await loadErrors();
		} catch (e) {
			toast.error('Failed to remove feed');
		}
	}

	async function removeAllProblemFeeds() {
		if (errors.length === 0) return;
		
		const feedsToRemove = errors.filter(e => e.consecutive_failures >= 3);
		if (feedsToRemove.length === 0) {
			toast.info('No feeds with 3+ consecutive failures to remove');
			return;
		}
		
		if (!confirm(`Remove ${feedsToRemove.length} feed${feedsToRemove.length === 1 ? '' : 's'} with 3+ consecutive failures? This will permanently delete these feeds and all their items.`)) return;
		
		try {
			for (const feedError of feedsToRemove) {
				await deleteFeed(feedError.feed.id);
			}
			toast.success(`Removed ${feedsToRemove.length} problematic feed${feedsToRemove.length === 1 ? '' : 's'}`);
			await loadErrors();
		} catch (e) {
			toast.error('Failed to remove some feeds');
			await loadErrors();
		}
	}

	function getErrorSeverity(consecutiveFailures: number): 'warning' | 'error' {
		return consecutiveFailures >= 3 ? 'error' : 'warning';
	}

	function formatTimeAgo(date: Date): string {
		const now = new Date();
		const diff = now.getTime() - new Date(date).getTime();
		const hours = Math.floor(diff / (1000 * 60 * 60));
		const days = Math.floor(hours / 24);
		
		if (days > 0) return `${days}d ago`;
		if (hours > 0) return `${hours}h ago`;
		return 'Just now';
	}
</script>

<Section id="errors" title="Feed Errors" description="Feeds that are currently experiencing issues">
	{#if loading}
		<div class="flex items-center justify-center p-8">
			<div class="loading loading-spinner loading-md"></div>
		</div>
	{:else if error}
		<div class="alert alert-error">
			<AlertTriangle class="size-4" />
			<span>Failed to load errors: {error}</span>
		</div>
	{:else if errors.length === 0}
		<div class="alert alert-success">
			<svg xmlns="http://www.w3.org/2000/svg" class="size-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
			</svg>
			<span>All feeds are working correctly!</span>
		</div>
	{:else}
		<div class="space-y-4">
			<div class="flex items-center justify-between">
				<p class="text-sm text-base-content/60">
					{errors.length} feed{errors.length === 1 ? '' : 's'} with errors
				</p>
				<div class="flex gap-2">
					<button 
						onclick={retryAllFailed} 
						disabled={globalState.demoMode}
						class="btn btn-sm btn-outline"
					>
						<RefreshCw class="size-4" />
						Retry All
					</button>
					{#if errors.some(e => e.consecutive_failures >= 3)}
						<button 
							onclick={removeAllProblemFeeds} 
							disabled={globalState.demoMode}
							class="btn btn-sm btn-outline btn-error"
						>
							<Trash2 class="size-4" />
							Remove Problem Feeds
						</button>
					{/if}
				</div>
			</div>

			<div class="space-y-3">
				{#each errors as feedError}
					{@const severity = getErrorSeverity(feedError.consecutive_failures)}
					<div class={`alert alert-${severity} border border-opacity-30`}>
						<div class="flex items-start gap-3 w-full">
							<AlertTriangle class="size-5 shrink-0 mt-0.5 text-current" />
							
							<div class="flex-1 min-w-0">
								<div class="flex items-center gap-2 mb-1">
									<div class="avatar">
										<div class="size-4 rounded-full">
											<img src={getFavicon(feedError.feed.link)} alt={feedError.feed.name} loading="lazy" />
										</div>
									</div>
									<span class="font-semibold text-current truncate">{feedError.feed.name}</span>
									<span class="badge badge-sm badge-outline">
										{feedError.consecutive_failures} failure{feedError.consecutive_failures === 1 ? '' : 's'}
									</span>
								</div>
								
								<p class="text-sm text-current opacity-90 mb-2 break-words font-medium">
									{feedError.error_message}
								</p>
								
								<div class="flex items-center gap-4 text-xs text-current opacity-75">
									<div class="flex items-center gap-1">
										<Clock class="size-3" />
										Last attempt: {formatTimeAgo(feedError.last_attempt)}
									</div>
									<a href={feedError.feed.link} target="_blank" class="flex items-center gap-1 hover:underline text-current">
										<ExternalLink class="size-3" />
										View Feed
									</a>
								</div>
							</div>
							
							<div class="flex gap-2 shrink-0">
								<a href="/feeds/{feedError.feed.id}" class="btn btn-xs btn-ghost">
									View
								</a>
								<button 
									onclick={() => retryFeed(feedError.feed.id)} 
									disabled={globalState.demoMode}
									class="btn btn-xs btn-outline"
								>
									<RefreshCw class="size-3" />
									Retry
								</button>
								<button 
									onclick={() => removeFeed(feedError.feed.id, feedError.feed.name)} 
									disabled={globalState.demoMode}
									class="btn btn-xs btn-outline btn-error"
								>
									<Trash2 class="size-3" />
									Remove
								</button>
							</div>
						</div>
					</div>
				{/each}
			</div>
		</div>
	{/if}
</Section>
