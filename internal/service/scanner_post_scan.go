package service

import (
	"context"
	"time"

	"go.uber.org/zap"
)

func (s *ScannerService) invalidateMediaCache(ctx context.Context) {
	if s != nil && s.cache != nil {
		s.cache.DeletePrefix(ctx, "media:")
		s.cache.DeletePrefix(ctx, "stats:")
	}
}

func (s *ScannerService) startAutoScrape(ctx context.Context, libraryID string) {
	scrapeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Minute)
	go func() {
		defer cancel()
	defer func() {
		s.autoScrapeMu.Lock()
		delete(s.autoScrapeActive, libraryID)
		s.autoScrapeMu.Unlock()
	}()
	if _, err := s.scraper.EnrichLibraryDetailedWithOptions(scrapeCtx, libraryID, ScrapeOptions{}); err != nil {
		s.log.Warn("scraper enrich failed", zap.Error(err))
		return
	}
		// In-place metadata enrich only (candidates: rows still pending).
		// Deliberately NO organizer reclassify here — reclassify relocates
		// files whose stored location disagrees with the computed category,
		// which must never happen as a side effect of an automatic scan.
		_, _ = s.scraper.EnrichLibraryDetailedWithOptions(scrapeCtx, libraryID, ScrapeOptions{})
	}()
}
