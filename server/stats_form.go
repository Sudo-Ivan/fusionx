package server

import "time"

type StatsForm struct {
	TotalFeeds         int         `json:"total_feeds"`
	TotalItems         int         `json:"total_items"`
	TotalUnreadItems   int         `json:"total_unread_items"`
	TotalGroups        int         `json:"total_groups"`
	DatabaseSize       int64       `json:"database_size"`
	LastFeedUpdate     *time.Time  `json:"last_feed_update"`
	FailedFeeds        int         `json:"failed_feeds"`
}

type RespStats StatsForm
