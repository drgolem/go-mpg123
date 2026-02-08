package main

import (
	"fmt"
	"os"

	"github.com/drgolem/go-mpg123/mpg123"
)

func main() {
	// check command-line arguments
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: mp3dump <infile.mp3> <outfile.raw>")
		return
	}

	fmt.Printf("Go binding to mpg123 library\n")
	fmt.Printf("Supported decoders: %v\n", mpg123.SupportedDecoders())

	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	if err != nil {
		panic("could not initialize mpg123")
	}
	defer decoder.Delete()

	// open a file with decoder
	err = decoder.Open(os.Args[1])
	if err != nil {
		panic("error opening mp3 file")
	}
	defer decoder.Close()

	// get audio format information
	rate, chans, _ := decoder.GetFormat()
	fmt.Fprintln(os.Stderr, "Encoding: Signed 16bit")
	fmt.Fprintln(os.Stderr, "Sample Rate:", rate)
	fmt.Fprintln(os.Stderr, "Channels:", chans)
	fmt.Fprintln(os.Stderr, "Decoder:", decoder.CurrentDecoder())

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, chans, mpg123.ENC_SIGNED_16)

	// open output file
	o, err := os.Create(os.Args[2])
	if err != nil {
		panic("error opening output file")
	}
	defer o.Close()

	// decode mp3 file and dump output
	buf := make([]byte, 2048*16)
	for {
		n, err := decoder.Read(buf)
		o.Write(buf[0:n])
		if err != nil {
			break
		}
	}
}
