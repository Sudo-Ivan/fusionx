import { api } from './api';

export type Config = {
	demo_mode: boolean;
};

export type AppConfig = {
	feed_refresh_interval_minutes: number;
};

export async function getConfig(): Promise<Config> {
	return await api.get('config').json<Config>();
}

export async function getAppConfig(): Promise<AppConfig> {
	return await api.get('config').json<AppConfig>();
}

export async function updateAppConfig(config: { feed_refresh_interval_minutes: number }): Promise<void> {
	await api.patch('config', { json: config });
}
