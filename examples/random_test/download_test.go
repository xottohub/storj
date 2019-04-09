// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package inspector_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/storj/internal/memory"
	"storj.io/storj/internal/testcontext"
	"storj.io/storj/internal/testplanet"
	"storj.io/storj/storagenode/piecestore"
	"storj.io/storj/uplink"
)

func TestDownloadWithErrorByDeletingPieces(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	planet, err := testplanet.New(t, 1, 6, 1)
	require.NoError(t, err)
	defer ctx.Check(planet.Shutdown)

	planet.Start(ctx)

	uplink0 := planet.Uplinks[0]
	testData := make([]byte, 50*memory.MiB)
	_, err = rand.Read(testData)
	require.NoError(t, err)

	planet.Satellites[0].Repair.Checker.Loop.Stop()
	planet.Satellites[0].Repair.Repairer.Loop.Stop()

	config := &uplink.RSConfig{
		MinThreshold:     2,
		RepairThreshold:  3,
		SuccessThreshold: 4,
		MaxThreshold:     5,
	}

	err = uplink0.UploadWithConfig(ctx, planet.Satellites[0], config, "testbucket", "test/path", testData)
	require.NoError(t, err)

	// Delete some uploaded pieces
	var deletedPieceCount int
	for storagenode := 0; storagenode < 6; storagenode++ {
		orderInfo, err := planet.StorageNodes[storagenode].DB.Orders().ListUnsent(ctx, 100)
		require.NoError(t, err)

		if len(orderInfo) > 0 {
			err = planet.StorageNodes[storagenode].Storage2.Store.Delete(ctx, planet.Satellites[0].ID(), orderInfo[0].Limit.PieceId)
			require.NoError(t, err)
			deletedPieceCount += 1

			// Delete no more than 2 pieces
			if deletedPieceCount == 2 {
				break
			}
		}
	}

	require.True(t, deletedPieceCount > 1)

	data, err := uplink0.Download(ctx, planet.Satellites[0], "testbucket", "test/path")
	require.NoError(t, err)
	require.True(t, bytes.Equal(data, testData))
}

func TestDownloadWithErrorByChangingStorageNodesEndpoint(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	planet, err := testplanet.New(t, 1, 6, 1)
	require.NoError(t, err)
	defer ctx.Check(planet.Shutdown)

	planet.Start(ctx)

	uplink0 := planet.Uplinks[0]
	testData := make([]byte, 50*memory.MiB)
	_, err = rand.Read(testData)
	require.NoError(t, err)

	planet.Satellites[0].Repair.Checker.Loop.Stop()
	planet.Satellites[0].Repair.Repairer.Loop.Stop()

	config := &uplink.RSConfig{
		MinThreshold:     2,
		RepairThreshold:  3,
		SuccessThreshold: 4,
		MaxThreshold:     5,
	}

	err = uplink0.UploadWithConfig(ctx, planet.Satellites[0], config, "testbucket", "test/path", testData)
	require.NoError(t, err)

	planet.StorageNodes[0].Storage2.Endpoint = &piecestore.Endpoint{}

	data, err := uplink0.Download(ctx, planet.Satellites[0], "testbucket", "test/path")
	require.NoError(t, err)
	require.True(t, bytes.Equal(data, testData))
}
