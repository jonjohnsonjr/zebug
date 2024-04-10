package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jonjohnsonjr/zebug/internal/flate"
	"github.com/jonjohnsonjr/zebug/internal/gzip"
	"golang.org/x/sync/errgroup"
)

const debug = false

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

type Line struct {
	Type    string         `json:"type,omitempty"`
	In      int64          `json:"in"`
	Out     int64          `json:"out"`
	Bits    string         `json:"bits,omitempty"`
	Size    int64          `json:"size,omitempty"`
	Final   bool           `json:"final,omitempty"`
	Header  *flate.Header  `json:"header,omitempty"`
	Trailer *flate.Trailer `json:"trailer,omitempty"`
}

func run(args []string) error {
	updates := make(chan *flate.Checkpoint, 100)
	zr, err := gzip.NewReaderWithSpans(os.Stdin, 0, updates)
	if err != nil {
		return err
	}

	var eg errgroup.Group
	eg.Go(func() error {
		// TODO: Figure out why we see the first checkpoint twice.

		enc := json.NewEncoder(os.Stdout)

		var line *Line = nil
		in, out := int64(0), int64(0)
		for checkpoint := range updates {
			nextLine := &Line{
				Type:  checkpoint.Block.Kind,
				Out:   out,
				Final: checkpoint.Block.Final,
				In:    in,
			}
			if checkpoint.NB != 0 {
				nextLine.Bits = fmt.Sprintf("%b", checkpoint.B)
			}

			in = checkpoint.In
			out = checkpoint.Out + int64(checkpoint.WrPos)

			if debug {
				if checkpoint.ToRead != 0 {
					fmt.Printf("ToRead: %d\n", checkpoint.ToRead)
				}
				if line != nil {
					fmt.Printf("%t: (%d, %d, %d) (%d -> %d)\n", checkpoint.Full, checkpoint.WrPos, checkpoint.RdPos, checkpoint.ToRead, line.Out, nextLine.Out)
				} else {
					fmt.Printf("%t: (%d, %d) (nil -> %d) \n", checkpoint.Full, checkpoint.WrPos, checkpoint.RdPos, nextLine.Out)
				}
			}

			if line != nil && line.Header != nil {
				line.Header.Csize = nextLine.In - line.In
			}

			if checkpoint.Header != nil {
				nextLine.Header = checkpoint.Header
				if line != nil && line.Trailer != nil {
					line.Trailer.Csize = 8
					nextLine.In += line.Trailer.Csize
				}
			} else if checkpoint.Trailer != nil {
				nextLine.Trailer = checkpoint.Trailer
			}

			if line != nil {
				if nextLine.Trailer != nil {
					nextLine.Out = line.Out
				}
				if nextLine.Header != nil {
					nextLine.Out = line.Out
				}

				if line.Header == nil && line.Trailer == nil {
					line.Size = nextLine.Out - line.Out
				}

				if line.Type != "" {
					// Indent to show nesting.
					fmt.Fprintf(os.Stdout, "  ")
				}

				if err := enc.Encode(line); err != nil {
					return err
				}
			}

			line = nextLine
		}

		if err := enc.Encode(line); err != nil {
			return err
		}

		return nil
	})
	eg.Go(func() error {
		defer close(updates)

		if _, err := io.Copy(io.Discard, zr); err != nil {
			return err
		}

		return nil
	})

	return eg.Wait()
}
