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

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

type Line struct {
	Type    string         `json:"type,omitempty"`
	In      int64          `json:"in,omitempty"`
	Out     int64          `json:"out,omitempty"`
	Size    int            `json:"Size,omitempty"`
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

		for checkpoint := range updates {
			line := Line{
				Type:  checkpoint.Block.Kind,
				In:    checkpoint.In,
				Out:   checkpoint.Out,
				Size:  checkpoint.Block.Size,
				Final: checkpoint.Block.Final,
			}
			if checkpoint.Header != nil {
				line.Header = checkpoint.Header
			} else if checkpoint.Trailer != nil {
				line.Trailer = checkpoint.Trailer
			} else if line.Type != "" {
				// Indent to show nesting.
				fmt.Fprintf(os.Stdout, "  ")
			}

			if err := enc.Encode(line); err != nil {
				return err
			}
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
