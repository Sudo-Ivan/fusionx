import { api } from './api';
import type { Feed } from './model';

export type FeedError = {
	feed: Feed;
	error_message: string;
	consecutive_failures: number;
	last_attempt: Date;
};

export async function getFeedErrors(): Promise<FeedError[]> {
	const feeds = await api.get('feeds').json<{ feeds: Feed[] }>();
	return feeds.feeds
		.filter(feed => feed.failure && feed.failure.trim() !== '')
		.map(feed => ({
			feed,
			error_message: feed.failure,
			consecutive_failures: feed.consecutive_failures || 0,
			last_attempt: feed.updated_at
		}));
}
