<script lang="ts">
	import { getAppConfig, updateAppConfig } from '$lib/api/config';
	import { globalState } from '$lib/state.svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import Section from './Section.svelte';
	import { t } from '$lib/i18n';

	let feedRefreshInterval = $state(30);
	let loading = $state(false);

	onMount(async () => {
		try {
			const config = await getAppConfig();
			feedRefreshInterval = config.feed_refresh_interval_minutes;
		} catch (e) {
			console.error('Failed to load app config:', e);
			feedRefreshInterval = 30;
		}
	});

	async function handleSave() {
		if (loading || globalState.demoMode) return;
		
		loading = true;
		try {
			await updateAppConfig({
				feed_refresh_interval_minutes: feedRefreshInterval
			});
			toast.success(t('state.success'));
		} catch (e) {
			toast.error((e as Error).message);
		} finally {
			loading = false;
		}
	}

	function handleIntervalChange(event: Event) {
		const input = event.target as HTMLInputElement;
		feedRefreshInterval = parseInt(input.value) || 30;
	}
</script>

<Section
	id="system"
	title="System Settings"
	description="Configure system-wide settings for feed management"
>
	<div class="flex flex-col space-y-4">
		<fieldset class="fieldset">
			<legend class="fieldset-legend">Feed Refresh Interval</legend>
			<div class="flex items-center gap-4">
				<input
					type="number"
					min="1"
					max="10080"
					value={feedRefreshInterval}
					onchange={handleIntervalChange}
					disabled={globalState.demoMode || loading}
					class="input w-24"
				/>
				<span class="text-sm text-gray-600">minutes</span>
				<button
					onclick={handleSave}
					disabled={globalState.demoMode || loading}
					class="btn btn-primary btn-sm"
				>
					{loading ? 'Saving...' : t('common.save')}
				</button>
			</div>
			<div class="text-xs text-gray-500 mt-2">
				How often feeds are checked for updates. Minimum: 1 minute, Maximum: 1 week (10080 minutes)
			</div>
		</fieldset>
	</div>
</Section>
