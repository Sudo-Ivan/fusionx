import { type Feed, type Group } from './api/model';

export const globalState = $state({
	groups: [] as Group[],
	feeds: [] as Feed[],
	demoMode: false,
	readingPaneMode: 'default' as 'default' | '3pane' | 'drawer'
});

export function setGlobalFeeds(feeds: Feed[]) {
	globalState.feeds = feeds;
}

export function setGlobalGroups(groups: Group[]) {
	globalState.groups = groups;
}

export function setDemoMode(demoMode: boolean) {
	globalState.demoMode = demoMode;
}

export function setReadingPaneMode(mode: 'default' | '3pane' | 'drawer') {
	globalState.readingPaneMode = mode;
}

export function updateUnreadCount(feedId: number, change: number) {
	const feed = globalState.feeds.find((f) => f.id === feedId);
	if (feed) {
		feed.unread_count = Math.max(0, (feed.unread_count || 0) + change);
	}
}
