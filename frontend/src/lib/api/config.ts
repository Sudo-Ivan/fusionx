import { api } from './api';

export type Config = {
	demo_mode: boolean;
};

export type AppConfig = {
	feed_refresh_interval_minutes: number;
	reading_pane_mode: 'default' | '3pane' | 'drawer';
	demo_mode: boolean;
};

export async function getConfig(): Promise<Config> {
	return await api.get('config').json<Config>();
}

export async function getAppConfig(): Promise<AppConfig> {
	return await api.get('config').json<AppConfig>();
}

export async function updateAppConfig(config: Partial<AppConfig>): Promise<void> {
	await api.patch('config', { json: config });
}
