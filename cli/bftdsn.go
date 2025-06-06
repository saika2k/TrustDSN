package cli

import (
	"fmt"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-state-types/big"
	lapi "github.com/filecoin-project/lotus/api"
	"github.com/ipfs/go-cid"
	"github.com/klauspost/reedsolomon"
	"github.com/urfave/cli/v2"
	"github.com/zhuaiballl/homohash"
	"golang.org/x/xerrors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var BftDsnCmd = &cli.Command{
	Name:  "bftdsn",
	Usage: "Interact with BFT-DSN functions",
	Flags: []cli.Flag{},
	Subcommands: []*cli.Command{
		BftDsnEncodeCmd,
		BftDsnDecodeCmd,
		BftDsnDealCmd,
		BftDsnRetrieveCmd,
	},
}

var BftDsnEncodeCmd = &cli.Command{
	Name:      "encode",
	Usage:     "EC encode file",
	ArgsUsage: "[inputPath]",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "k",
			Value: 10,
			Usage: "parameter K of RS-code",
		},
		&cli.IntFlag{
			Name:  "m",
			Value: 3,
			Usage: "parameter M of RS-code",
		},
		&cli.BoolFlag{
			Name:  "keep-chunks",
			Value: false,
			Usage: "keep chunks produced during making deals",
		},
		&cli.BoolFlag{
			Name:  "hash",
			Value: false,
			Usage: "output hashes of chunks",
		},
	},
	Action: func(cctx *cli.Context) error {
		keepChunks := cctx.Bool("keep-chunks")
		outputHash := cctx.Bool("hash")
		if cctx.NArg() != 1 {
			return IncorrectNumArgs(cctx)
		}

		// Read file
		absPath, err := filepath.Abs(cctx.Args().First())
		if err != nil {
			return err
		}
		dataShards := cctx.Int("k")
		parShards := cctx.Int("m")

		err = encodeWithPath(absPath, dataShards, parShards, keepChunks, outputHash)
		if err != nil {
			return err
		}

		return nil
	},
}

var BftDsnDecodeCmd = &cli.Command{
	Name:      "decode",
	Usage:     "EC decode file",
	ArgsUsage: "[inputPath]",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "k",
			Value: 10,
			Usage: "parameter K of RS code",
		},
		&cli.IntFlag{
			Name:  "m",
			Value: 3,
			Usage: "parameter M of RS code",
		},
		&cli.StringFlag{
			Name:  "out",
			Usage: "Alternative output path",
		},
	},
	Action: func(cctx *cli.Context) error {
		if cctx.NArg() != 1 {
			return IncorrectNumArgs(cctx)
		}

		// Create encoding matrix
		dataShards := cctx.Int("k")
		parShards := cctx.Int("m")
		absPath, err := filepath.Abs(cctx.Args().First())
		if err != nil {
			return err
		}
		outFile := cctx.String("out")
		if outFile == "" {
			outFile = absPath
		}

		err = decodeWithPath(absPath, outFile, dataShards, parShards)
		if err != nil {
			return err
		}

		return nil
	},
}

var BftDsnDealCmd = &cli.Command{
	Name:      "deal",
	Usage:     "Make BFT-DSN deals",
	ArgsUsage: "[inputPath]",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "k",
			Value: 10,
			Usage: "parameter K of RS code",
		},
		&cli.IntFlag{
			Name:  "m",
			Value: 3,
			Usage: "parameter M of RS code",
		},
		&cli.BoolFlag{
			Name:  "keep-chunks",
			Value: false,
			Usage: "keep chunks produced during making deals",
		},
		&cli.BoolFlag{
			Name:  "hash",
			Value: false,
			Usage: "output hashes of chunks",
		},
	},
	Action: func(cctx *cli.Context) error {
		dataShards := cctx.Int("k")
		parShards := cctx.Int("m")
		outputHash := cctx.Bool("hash")

		// prepare chunks
		if cctx.NArg() != 1 {
			return IncorrectNumArgs(cctx)
		}

		path := cctx.Args().First()
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		beginTime := time.Now()
		fmt.Println("Start preparing deals")
		err = encodeWithPath(absPath, dataShards, parShards, true, outputHash)
		if err != nil {
			return err
		}

		// make deal
		api, closer, err := GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		ctx := ReqContext(cctx)
		afmt := NewAppFmt(cctx.App)
		wa, err := api.WalletDefaultAddress(ctx)
		if err != nil {
			return err
		}

		ts, err := LoadTipSet(ctx, cctx, api)
		if err != nil {
			return err
		}

		miners, err := api.StateListMiners(ctx, ts.Key())
		if err != nil {
			return err
		}
		n := len(miners)

		encoder, err := GetCidEncoder(cctx)
		if err != nil {
			return err
		}

		dir, file := filepath.Split(absPath)
		for i := 0; i < dataShards+parShards; i++ {
			outfn := fmt.Sprintf("%s.%d", file, i)
			pathI := filepath.Join(dir, outfn)

			fileRef := lapi.FileRef{
				Path:  pathI,
				IsCAR: false, //cctx.Bool("car"),
			}
			c, err := api.ClientImport(ctx, fileRef)
			if err != nil {
				return err
			}
			// send shards[i] to m
			ref := &storagemarket.DataRef{
				TransferType: storagemarket.TTGraphsync,
				Root:         c.Root, //cid
			}
			sdParams := &lapi.StartDealParams{
				Data:               ref, //shards[i%n]
				Wallet:             wa,
				Miner:              miners[i%n],
				EpochPrice:         big.NewInt(2600000000000000), //0.0026
				MinBlocksDuration:  uint64(518400),
				DealStartEpoch:     -1,
				FastRetrieval:      true,
				VerifiedDeal:       false,
				ProviderCollateral: big.Int{},
			}
			proposal, err := api.ClientStartDeal(ctx, sdParams)
			if err != nil {
				return err
			}
			afmt.Println("Transaction", i, encoder.Encode(*proposal))
			if !cctx.Bool("keep-chunks") {
				err = os.Remove(pathI)
				if err != nil {
					return err
				}
			}
		}
		fmt.Println("Deals all sent. Took", time.Now().Sub(beginTime).Truncate(time.Millisecond))
		return nil
	},
}

var BftDsnRetrieveCmd = &cli.Command{
	Name:        "retrieve",
	Usage:       "Make BFT-DSN retrieval deals",
	ArgsUsage:   "[inputPath outPath]",
	Description: "",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "k",
			Value: 10,
			Usage: "parameter K of RS code",
		},
		&cli.IntFlag{
			Name:  "m",
			Value: 3,
			Usage: "parameter M of RS code",
		},
		&cli.BoolFlag{
			Name:  "keep-chunks",
			Value: false,
			Usage: "keep chunks produced during retrieval",
		},
		&cli.BoolFlag{
			Name:  "hash",
			Value: false,
			Usage: "output hashes of chunks",
		},
	},
	Action: func(cctx *cli.Context) error {
		dataShards := cctx.Int("k")
		parShards := cctx.Int("m")
		outputHash := cctx.Bool("hash")
		path := cctx.Args().First()
		outpath := cctx.Args().Get(1)

		// prepare chunks
		if cctx.NArg() != 2 {
			return IncorrectNumArgs(cctx)
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		err = encodeWithPath(absPath, dataShards, parShards, true, outputHash)
		if err != nil {
			return err
		}

		// make deal
		api, closer, err := GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		ctx := ReqContext(cctx)
		afmt := NewAppFmt(cctx.App)
		wa, err := api.WalletDefaultAddress(ctx)
		if err != nil {
			return err
		}

		ts, err := LoadTipSet(ctx, cctx, api)
		if err != nil {
			return err
		}

		miners, err := api.StateListMiners(ctx, ts.Key())
		if err != nil {
			return err
		}
		n := len(miners)

		//encoder, err := GetCidEncoder(cctx)
		//if err != nil {
		//	return err
		//}

		// prepare cid list
		cids := make([]cid.Cid, dataShards+parShards)
		dir, file := filepath.Split(absPath)
		for i := 0; i < dataShards+parShards; i++ {
			outfn := fmt.Sprintf("%s.%d", file, i)
			pathI := filepath.Join(dir, outfn)

			fileRef := lapi.FileRef{
				Path:  pathI,
				IsCAR: false, //cctx.Bool("car"),
			}
			c, err := api.ClientImport(ctx, fileRef)
			if err != nil {
				return err
			}
			// c.Root is the cid
			cids[i] = c.Root
			if !cctx.Bool("keep-chunks") {
				err = os.Remove(pathI)
				if err != nil {
					return err
				}
			}
		}
		afmt.Println("CID list obtained.")

		fapi, fcloser, err := GetFullNodeAPIV1(cctx)
		if err != nil {
			return err
		}
		defer fcloser()
		beginTime := time.Now()
		afmt.Println("Retrieve begins")
		for i := 0; i < dataShards+parShards; i++ {
			shardcid := cids[i]
			var eref *lapi.ExportRef
			var offer lapi.QueryOffer
			minerAddr := miners[i%n]
			offer, err = fapi.ClientMinerQueryOffer(ctx, minerAddr, shardcid, nil)
			if err != nil {
				return err
			}
			if offer.Err != "" {
				return fmt.Errorf("offer error: %s", offer.Err)
			}

			o := offer.Order(wa)

			subscribeEvents, err := fapi.ClientGetRetrievalUpdates(ctx)
			if err != nil {
				return xerrors.Errorf("error setting up retrieval updates: %w", err)
			}
			retrievalRes, err := fapi.ClientRetrieve(ctx, o)
			if err != nil {
				return xerrors.Errorf("error setting up retrieval: %w", err)
			}

		readEvents:
			for {
				var evt lapi.RetrievalInfo
				select {
				case <-ctx.Done():
					return xerrors.New("Retrieval Timed Out")
				case evt = <-subscribeEvents:
					if evt.ID != retrievalRes.DealID {
						// we can't check the deal ID ahead of time because:
						// 1. We need to subscribe before retrieving.
						// 2. We won't know the deal ID until after retrieving.
						continue
					}
				}

				//event := "New"
				//if evt.Event != nil {
				//	event = retrievalmarket.ClientEvents[*evt.Event]
				//}
				//
				//printf("Recv %s, Paid %s, %s (%s), %s\n",
				//	types.SizeStr(types.NewInt(evt.BytesReceived)),
				//	types.FIL(evt.TotalPaid),
				//	strings.TrimPrefix(event, "ClientEvent"),
				//	strings.TrimPrefix(retrievalmarket.DealStatuses[evt.Status], "DealStatus"),
				//	time.Now().Sub(start).Truncate(time.Millisecond),
				//)

				switch evt.Status {
				case retrievalmarket.DealStatusCompleted:
					break readEvents
				case retrievalmarket.DealStatusRejected:
					return xerrors.Errorf("Retrieval Proposal Rejected: %s", evt.Message)
				case retrievalmarket.DealStatusCancelled:
					return xerrors.Errorf("Retrieval Proposal Cancelled: %s", evt.Message)
				case
					retrievalmarket.DealStatusDealNotFound,
					retrievalmarket.DealStatusErrored:
					return xerrors.Errorf("Retrieval Error: %s", evt.Message)
				}
			}

			eref = &lapi.ExportRef{
				Root:   shardcid,
				DealID: retrievalRes.DealID,
			}
			if eref == nil {
				return xerrors.Errorf("failed to find providers")
			}

			err = fapi.ClientExport(ctx, *eref, lapi.FileRef{
				Path:  fmt.Sprintf("%s.%d", outpath, i),
				IsCAR: false,
			})
			if err != nil {
				return err
			}
			afmt.Println("Successfully retrieved one chunk")
		}
		afmt.Println("Chunks retrieved. Took", time.Now().Sub(beginTime).Truncate(time.Millisecond))

		// decode and get the output file
		err = decodeWithPath(outpath, outpath, dataShards, parShards)
		if err != nil {
			return err
		}
		afmt.Println("Success.")

		if !cctx.Bool("keep-chunks") {
			for i := 0; i < dataShards+parShards; i++ {
				chunkPath := fmt.Sprintf("%s.%d", outpath, i)
				err = os.Remove(chunkPath)
				if err != nil {
					return err
				}
			}
		}

		return nil
	},
}

// RSEncode with input in filepath and write shards in corresponding paths
func encodeWithPath(path string, dataShards, parShards int, keepChunk, hash bool) error {
	fmt.Println("Opening", path)
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	shards, err := encode(f, dataShards, parShards, hash)
	if err != nil {
		return err
	}

	// Write out the resulting files.
	dir, file := filepath.Split(path)
	for i, shard := range shards {
		outfn := fmt.Sprintf("%s.%d", file, i)

		if keepChunk {
			fmt.Println("Writing to", outfn)
			err = ioutil.WriteFile(filepath.Join(dir, outfn), shard, 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// RSEncode with input in byte array and output in byte arrays
func encode(f []byte, dataShards, parShards int, hash bool) ([][]byte, error) {
	// Create encoding matrix
	enc, err := reedsolomon.New(dataShards, parShards)
	if err != nil {
		return nil, err
	}

	shards, err := enc.Split(f)
	if err != nil {
		return nil, err
	}
	fmt.Printf("File split into %d data+parity shards with %d bytes/shard.\n", len(shards), len(shards[0]))

	ho := homohash.New()
	hashes := make([][]byte, len(shards))
	for i, shard := range shards {
		ho.Reset()
		hashes[i] = make([]byte, 32)
		ho.Write(shard)
		copy(hashes[i], ho.Sum(nil))
	}
	fmt.Println()

	if hash {
		err = enc.Encode(hashes)
		if err != nil {
			return nil, err
		}
		fmt.Println("Encoded hashes: ")
		for _, hash := range hashes {
			fmt.Print(hash, " ")
		}
		fmt.Println()
	}

	// Encode parity
	fmt.Println("Encode begins")
	beginTime := time.Now()
	err = enc.Encode(shards)
	if err != nil {
		return nil, err
	}
	fmt.Println("Encode finished. Took", time.Now().Sub(beginTime).Truncate(time.Millisecond))

	if hash {
		fmt.Println("Hashes of encoded shards: ")
		for _, shard := range shards {
			ho.Reset()
			ho.Write(shard)
			fmt.Print(ho.Sum(nil), " ")
		}
		fmt.Println()
	}

	return shards, nil
}

func decodeWithPath(input, output string, dataShards, parShards int) error {
	// Create shards and load the data
	shards := make([][]byte, dataShards+parShards)
	var err error
	for i := range shards {
		infn := fmt.Sprintf("%s.%d", input, i)
		fmt.Println("Opening", infn)
		shards[i], err = ioutil.ReadFile(infn)
		if err != nil {
			fmt.Println("Error reading file", err)
			shards[i] = nil
		}
	}

	// Join the shards and write them
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	fmt.Println("Writing data to", output)
	err = decode(f, shards, dataShards, parShards)
	if err != nil {
		return err
	}

	return nil
}

func decode(f *os.File, shards [][]byte, dataShards, parShards int) error {
	// Create encoding matrix
	enc, err := reedsolomon.New(dataShards, parShards)
	if err != nil {
		return err
	}

	fmt.Println("Decode begins")
	beginTime := time.Now()
	// Verify the shards
	ok, err := enc.Verify(shards)
	if ok {
		fmt.Println("No reconstruction needed")
	} else {
		fmt.Println("Verification failed. Reconstructing data")
		err = enc.Reconstruct(shards)
		if err != nil {
			return err
		}
		ok, err = enc.Verify(shards)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrpted.")
			return err
		}
	}

	// We don't know the exact filesize. ?
	err = enc.Join(f, shards, len(shards[0])*dataShards)
	if err != nil {
		return err
	}
	fmt.Println("Decode finished. Took", time.Now().Sub(beginTime).Truncate(time.Millisecond))

	return nil
}
