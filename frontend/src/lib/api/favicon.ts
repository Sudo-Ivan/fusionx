/**
 *  RSSHub paths mapped to their respective hostname.
 *  Sorted by hostname first then by RSSHub path.
 */
const rssHubMap = {
	'/papers/category/arxiv': 'arxiv.org',
	'/trendingpapers/papers': 'arxiv.org',
	'/github': 'github.com',
	'/google': 'google.com',
	'/dockerhub': 'hub.docker.com',
	'/imdb': 'imdb.com',
	'/hackernews': 'news.ycombinator.com',
	'/phoronix': 'phoronix.com',
	'/rsshub': 'rsshub.app',
	'/twitch': 'twitch.tv',
	'/youtube': 'youtube.com'
};

function extractHostname(feedLink: string): string {
	const url = new URL(feedLink);
	let hostname = url.hostname.replace(/^www\./, '');

	if (hostname.includes('rsshub')) {
		for (const prefix in rssHubMap) {
			if (url.pathname.startsWith(prefix)) {
				hostname = rssHubMap[prefix as keyof typeof rssHubMap];
				break;
			}
		}
	}

	return hostname;
}

function getCacheKey(hostname: string): string {
	let h1 = 0xdeadbeef >>> 0, h2 = 0x41c6ce57 >>> 0;
	for (let i = 0; i < hostname.length; i++) {
		const ch = hostname.charCodeAt(i);
		h1 = (h1 ^ ch) >>> 0;
		h1 = Math.imul(h1, 2654435761) >>> 0;
		h2 = (h2 ^ ch) >>> 0;
		h2 = Math.imul(h2, 1597334677) >>> 0;
	}
	h1 = (Math.imul((h1 ^ (h1 >>> 16)), 2246822507) ^ Math.imul((h2 ^ (h2 >>> 13)), 3266489909)) >>> 0;
	h2 = (Math.imul((h2 ^ (h2 >>> 16)), 2246822507) ^ Math.imul((h1 ^ (h1 >>> 13)), 3266489909)) >>> 0;
	const hash = ((h2 & 0x1fffff) * 4294967296) + (h1 >>> 0);
	return ((hash & 0xffffffff) >>> 0).toString(16).slice(0, 8).padStart(8, '0');
}

export function getFavicon(feedLink: string): string {
	try {
		const hostname = extractHostname(feedLink);
		const cacheKey = getCacheKey(hostname);
		return `/api/favicons/${cacheKey}.png`;
	} catch (error) {
		console.warn('Failed to generate favicon URL:', error);
		return '/api/favicons/default.png';
	}
}
