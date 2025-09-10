<script lang="ts">
	import { getStats, type Stats } from '$lib/api/stats';
	import { t } from '$lib/i18n';
	import { Database, Rss, FileText, Folder, AlertTriangle, Clock } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import Section from './Section.svelte';

	let stats = $state<Stats | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);

	onMount(async () => {
		try {
			stats = await getStats();
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	});

	function formatBytes(bytes: number): string {
		if (bytes === 0) return '0 B';
		const k = 1024;
		const sizes = ['B', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
	}

	function formatDate(date: Date | null): string {
		if (!date) return 'Never';
		return new Date(date).toLocaleString();
	}
</script>

<Section id="stats" title="Statistics" description="Overview of your Fusion instance data">
	{#if loading}
		<div class="flex items-center justify-center p-8">
			<div class="loading loading-spinner loading-md"></div>
		</div>
	{:else if error}
		<div class="alert alert-error">
			<AlertTriangle class="size-4" />
			<span>Failed to load statistics: {error}</span>
		</div>
	{:else if stats}
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
			<div class="stat bg-base-200 rounded-lg p-4">
				<div class="stat-figure text-primary">
					<Rss class="size-8" />
				</div>
				<div class="stat-title text-sm">Total Feeds</div>
				<div class="stat-value text-2xl">{stats.total_feeds}</div>
				{#if stats.failed_feeds > 0}
					<div class="stat-desc text-error">
						{stats.failed_feeds} with errors
					</div>
				{/if}
			</div>

			<div class="stat bg-base-200 rounded-lg p-4">
				<div class="stat-figure text-secondary">
					<FileText class="size-8" />
				</div>
				<div class="stat-title text-sm">Total Items</div>
				<div class="stat-value text-2xl">{stats.total_items.toLocaleString()}</div>
				<div class="stat-desc">
					{stats.total_unread_items.toLocaleString()} unread
				</div>
			</div>

			<div class="stat bg-base-200 rounded-lg p-4">
				<div class="stat-figure text-accent">
					<Folder class="size-8" />
				</div>
				<div class="stat-title text-sm">Groups</div>
				<div class="stat-value text-2xl">{stats.total_groups}</div>
			</div>

			<div class="stat bg-base-200 rounded-lg p-4">
				<div class="stat-figure text-info">
					<Database class="size-8" />
				</div>
				<div class="stat-title text-sm">Database Size</div>
				<div class="stat-value text-2xl">{formatBytes(stats.database_size)}</div>
			</div>

			<div class="stat bg-base-200 rounded-lg p-4 md:col-span-2">
				<div class="stat-figure text-warning">
					<Clock class="size-8" />
				</div>
				<div class="stat-title text-sm">Last Feed Update</div>
				<div class="stat-value text-lg">{formatDate(stats.last_feed_update)}</div>
			</div>
		</div>
	{/if}
</Section>
