package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sudo-Ivan/fusionx/repo"
	"github.com/Sudo-Ivan/fusionx/service/favicon"
	"github.com/labstack/echo/v4"
)

type faviconAPI struct {
	cacheDir   string
	faviconSvc *favicon.Service
	feedRepo   *repo.Feed
}

func newFaviconAPI(cacheDir string) *faviconAPI {
	return &faviconAPI{
		cacheDir:   cacheDir,
		faviconSvc: favicon.NewService(cacheDir),
		feedRepo:   repo.NewFeed(repo.DB),
	}
}

func (f *faviconAPI) ServeFavicon(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "filename parameter required")
	}

	faviconPath := filepath.Join(f.cacheDir, filename)
	
	// If favicon doesn't exist, try to fetch it from the actual feed
	if _, err := os.Stat(faviconPath); os.IsNotExist(err) {
		requestedHash := strings.TrimSuffix(filename, ".png")
		
		// Find all feeds and check which one matches this hash
		feeds, err := f.feedRepo.FindByFaviconHash(requestedHash)
		if err == nil {
			for _, feed := range feeds {
				if feed.Link != nil {
					// Generate hash for this feed's URL and see if it matches
					if testPath, err := f.faviconSvc.GetFaviconPath(*feed.Link); err == nil {
						testHash := strings.TrimSuffix(filepath.Base(testPath), ".png")
						if testHash == requestedHash {
							// Found the matching feed! The favicon was fetched during GetFaviconPath
							c.Response().Header().Set("Cache-Control", "public, max-age=86400")
							return c.File(testPath)
						}
					}
				}
			}
		}
		
		// If we couldn't find and fetch a real favicon, create a default as last resort
		if _, err := f.faviconSvc.CreateDefaultFavicon(faviconPath); err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "favicon not found")
		}
	}
	
	c.Response().Header().Set("Cache-Control", "public, max-age=86400")
	return c.File(faviconPath)
}
