<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { getAppConfig, updateAppConfig } from '$lib/api/config';
	import { globalState, setReadingPaneMode } from '$lib/state.svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import Section from './Section.svelte';
	import { t } from '$lib/i18n';

	let feedRefreshInterval = $state(30);
	let readingPaneMode = $state<'default' | '3pane' | 'drawer'>('default');
	let loading = $state(false);

	onMount(async () => {
		try {
			const config = await getAppConfig();
			feedRefreshInterval = config.feed_refresh_interval_minutes;
			readingPaneMode = config.reading_pane_mode;
			setReadingPaneMode(config.reading_pane_mode);
		} catch (e) {
			console.error('Failed to load app config:', e);
			feedRefreshInterval = 30;
			readingPaneMode = 'default';
		}
	});

	async function handleSave() {
		if (loading || globalState.demoMode) return;
		
		loading = true;
		try {
			await updateAppConfig({
				feed_refresh_interval_minutes: feedRefreshInterval,
				reading_pane_mode: readingPaneMode
			});
			setReadingPaneMode(readingPaneMode);
			toast.success(t('state.success'));
			
			// Force reload of layout config
			await invalidateAll();
		} catch (e) {
			toast.error((e as Error).message);
			console.error('Config update error:', e);
		} finally {
			loading = false;
		}
	}

	function handleIntervalChange(event: Event) {
		const input = event.target as HTMLInputElement;
		feedRefreshInterval = parseInt(input.value) || 30;
	}

	function handleReadingPaneModeChange(event: Event) {
		const select = event.target as HTMLSelectElement;
		readingPaneMode = select.value as 'default' | '3pane' | 'drawer';
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
				<span class="text-sm text-gray-400">minutes</span>
			</div>
			<div class="text-xs text-gray-500 mt-2">
				How often feeds are checked for updates. Minimum: 1 minute, Maximum: 1 week (10080 minutes)
			</div>
		</fieldset>

		<fieldset class="fieldset">
			<legend class="fieldset-legend">Reading Pane Layout</legend>
			<select 
				value={readingPaneMode} 
				onchange={handleReadingPaneModeChange}
				disabled={globalState.demoMode || loading}
				class="select w-full max-w-xs"
			>
				<option value="default">Default (Navigate to item page)</option>
				<option value="3pane">3-Pane View (List + Reading pane)</option>
				<option value="drawer">Drawer (Slide-out reading pane)</option>
			</select>
			<div class="text-xs text-gray-500 mt-2">
				Choose how articles are displayed when clicked. Changes apply immediately after saving.
			</div>
		</fieldset>

		<div class="flex justify-end">
			<button
				onclick={handleSave}
				disabled={globalState.demoMode || loading}
				class="btn btn-primary"
			>
				{loading ? 'Saving...' : t('common.save')}
			</button>
		</div>
	</div>
</Section>
