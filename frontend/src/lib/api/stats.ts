import { api } from './api';

export type Stats = {
	total_feeds: number;
	total_items: number;
	total_unread_items: number;
	total_groups: number;
	database_size: number;
	last_feed_update: Date | null;
	failed_feeds: number;
};

export async function getStats(): Promise<Stats> {
	return await api.get('stats').json<Stats>();
}
