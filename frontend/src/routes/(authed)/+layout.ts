import { getConfig } from '$lib/api/config';
import { listFeeds } from '$lib/api/feed';
import { allGroups } from '$lib/api/group';
import { setDemoMode, setGlobalFeeds, setGlobalGroups } from '$lib/state.svelte';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ depends }) => {
	depends('app:feeds', 'app:groups', 'app:config');

	await Promise.all([
		getConfig().then((config) => {
			setDemoMode(config.demo_mode);
		}),
		allGroups().then((groups) => {
			groups.sort((a, b) => a.id - b.id);
			setGlobalGroups(groups);
		}),
		listFeeds().then((feeds) => {
			setGlobalFeeds(feeds);
		})
	]);

	return {};
};
