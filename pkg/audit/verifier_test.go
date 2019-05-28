// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package audit

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vivint/infectious"
	"storj.io/storj/internal/testcontext"
	"storj.io/storj/internal/testidentity"
	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/peertls/tlsopts"
	"storj.io/storj/pkg/storj"
	"storj.io/storj/pkg/transport"
)

func TestFailingAudit(t *testing.T) {
	const (
		required = 8
		total    = 14
	)

	f, err := infectious.NewFEC(required, total)
	if err != nil {
		panic(err)
	}

	shares := make([]infectious.Share, total)
	output := func(s infectious.Share) {
		shares[s.Number] = s.DeepCopy()
	}

	// the data to encode must be padded to a multiple of required, hence the
	// underscores.
	err = f.Encode([]byte("hello, world! __"), output)
	if err != nil {
		panic(err)
	}

	modifiedShares := make([]infectious.Share, len(shares))
	for i := range shares {
		modifiedShares[i] = shares[i].DeepCopy()
	}

	modifiedShares[0].Data[1] = '!'
	modifiedShares[2].Data[0] = '#'
	modifiedShares[3].Data[1] = '!'
	modifiedShares[4].Data[0] = 'b'

	badPieceNums := []int{0, 2, 3, 4}

	ctx := context.Background()
	auditPkgShares := make(map[int]Share, len(modifiedShares))
	for i := range modifiedShares {
		auditPkgShares[modifiedShares[i].Number] = Share{
			PieceNum: modifiedShares[i].Number,
			Data:     append([]byte(nil), modifiedShares[i].Data...),
		}
	}

	pieceNums, correctedShares, err := auditShares(ctx, 8, 14, auditPkgShares)
	if err != nil {
		panic(err)
	}

	for i, num := range pieceNums {
		if num != badPieceNums[i] {
			t.Fatal("expected nums in pieceNums to be same as in badPieceNums")
		}
	}

	require.Equal(t, shares, correctedShares)
}

func TestNotEnoughShares(t *testing.T) {
	const (
		required = 8
		total    = 14
	)

	f, err := infectious.NewFEC(required, total)
	if err != nil {
		panic(err)
	}

	shares := make([]infectious.Share, total)
	output := func(s infectious.Share) {
		shares[s.Number] = s.DeepCopy()
	}

	// the data to encode must be padded to a multiple of required, hence the
	// underscores.
	err = f.Encode([]byte("hello, world! __"), output)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	auditPkgShares := make(map[int]Share, len(shares))
	for i := range shares {
		auditPkgShares[shares[i].Number] = Share{
			PieceNum: shares[i].Number,
			Data:     append([]byte(nil), shares[i].Data...),
		}
	}
	_, _, err = auditShares(ctx, 20, 40, auditPkgShares)
	require.Contains(t, err.Error(), "infectious: must specify at least the number of required shares")
}

func TestVerifier_getNodeConnection(t *testing.T) {
	t.Run("error: node is offline", func(t *testing.T) {
		listener, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)

		ident, err := testidentity.PregeneratedIdentity(0, storj.LatestIDVersion())
		require.NoError(t, err)

		opts, err := tlsopts.NewOptions(ident, tlsopts.Config{
			PeerIDVersions: "*",
		})
		require.NoError(t, err)

		var (
			nodeAddr = &pb.NodeAddress{
				Transport: pb.NodeTransport_TCP_TLS_GRPC,
				Address:   listener.Addr().String(),
			}
			verifier = &Verifier{
				transport: transport.NewClient(opts),
			}
		)

		// Close Node server for simulating a connection to a offline node
		require.NoError(t, listener.Close())

		_, err = verifier.getNodeConnection(storj.NodeID{123}, nodeAddr)
		require.Error(t, err)
		require.Equal(t, errStorageNodeOffline, err)
	})

	t.Run("error: node dialing accept conn and closed it", func(t *testing.T) {
		listener, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, listener.Close()) }()

		ctx := testcontext.NewWithTimeout(t, 30*time.Second)
		ctx.Go(func() (err error) {
			for i := 0; i < 2; i++ {
				conn, err := listener.Accept()
				if err != nil {
					return err
				}

				err = conn.Close()
				if err != nil {
					return err
				}
			}

			return nil
		})

		defer ctx.Cleanup()

		ident, err := testidentity.PregeneratedIdentity(0, storj.LatestIDVersion())
		require.NoError(t, err)

		opts, err := tlsopts.NewOptions(ident, tlsopts.Config{
			PeerIDVersions: "*",
		})
		require.NoError(t, err)

		var (
			nodeAddr = &pb.NodeAddress{
				Transport: pb.NodeTransport_TCP_TLS_GRPC,
				Address:   listener.Addr().String(),
			}
			verifier = &Verifier{
				transport: transport.NewClient(opts),
			}
		)

		_, err = verifier.getNodeConnection(storj.NodeID{123}, nodeAddr)
		require.Error(t, err)
		require.Equal(t, errStorageNodeDialUnexpected, err)
	})

	t.Run("error: node dialing send garbage", func(t *testing.T) {
		listener, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, listener.Close()) }()

		ctx := testcontext.NewWithTimeout(t, 30*time.Second)
		ctx.Go(func() (err error) {
			for i := 0; i < 2; i++ {
				conn, err := listener.Accept()
				if err != nil {
					return err
				}

				_, err = conn.Write([]byte("garbage"))
				if err != nil {
					return err
				}

				err = conn.Close()
				if err != nil {
					return err
				}
			}

			return nil
		})

		defer ctx.Cleanup()

		ident, err := testidentity.PregeneratedIdentity(0, storj.LatestIDVersion())
		require.NoError(t, err)

		opts, err := tlsopts.NewOptions(ident, tlsopts.Config{
			PeerIDVersions: "*",
		})
		require.NoError(t, err)

		var (
			nodeAddr = &pb.NodeAddress{
				Transport: pb.NodeTransport_TCP_TLS_GRPC,
				Address:   listener.Addr().String(),
			}
			verifier = &Verifier{
				transport: transport.NewClient(opts),
			}
		)

		_, err = verifier.getNodeConnection(storj.NodeID{123}, nodeAddr)
		require.Error(t, err)
		require.Equal(t, errStorageNodeDialUnexpected, err)
	})
}
