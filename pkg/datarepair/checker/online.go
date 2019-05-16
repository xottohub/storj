// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package checker

import (
	"context"
	"time"

	"storj.io/storj/pkg/overlay"
	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/storj"
)

// TODO: figure out a better name
type ReliableCache struct {
	overlay    *overlay.Cache
	staleness  time.Duration
	lastUpdate time.Time
	reliable   map[storj.NodeID]struct{}
}

func NewReliableCache(overlay *overlay.Cache, staleness time.Duration) *ReliableCache {
	return &ReliableCache{
		overlay:   overlay,
		staleness: staleness,
		reliable:  map[storj.NodeID]struct{}{},
	}
}

func (cache *ReliableCache) LastUpdate() time.Time { return cache.lastUpdate }

func (cache *ReliableCache) MissingPieces(ctx context.Context, created time.Time, pieces []*pb.RemotePiece) ([]int32, error) {
	if created.After(cache.lastUpdate) || time.Since(cache.lastUpdate) > cache.staleness {
		err := cache.refresh(ctx)
		if err != nil {
			return nil, err
		}
	}

	var unreliable []int32
	for _, piece := range pieces {
		if _, ok := cache.reliable[piece.NodeId]; !ok {
			unreliable = append(unreliable, piece.PieceNum)
		}
	}
	return unreliable, nil
}

func (cache *ReliableCache) refresh(ctx context.Context) error {
	for id := range cache.reliable {
		delete(cache.reliable, id)
	}

	cache.lastUpdate = time.Now()

	nodes, err := cache.overlay.Reliable(ctx)
	if err != nil {
		return Error.Wrap(err)
	}

	for _, id := range nodes {
		cache.reliable[id] = struct{}{}
	}

	return nil
}
