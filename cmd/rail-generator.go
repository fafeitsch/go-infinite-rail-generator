package main

import (
	"bufio"
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/noise"
	"github.com/fafeitsch/go-infinite-rail-generator/renderer"
	"github.com/fafeitsch/go-infinite-rail-generator/version"
	"github.com/fafeitsch/go-infinite-rail-generator/web"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
)

const sizeFlag = "size"
const tileFlag = "tile"
const bindFlag = "bind"
const portFlag = "port"
const shiftFlag = "shift"
const seedFlag = "seed"
const townNameFlag = "towns"

func main() {
	err := (&cli.App{
		Name:            "rail-generator",
		Usage:           "A procedural rail line generator",
		UsageText:       "rail-generator [global options] command [command options]",
		HideHelpCommand: true,
		Copyright:       "2021, Fabian Feitsch (info@fafeitsch.de), Licensed under MIT",
		Authors:         []*cli.Author{{Name: "Fabian Feitsch"}},
		Version:         fmt.Sprintf("%s (%s)", version.BuildVersion, version.BuildTime),
		Flags: []cli.Flag{
			&cli.StringFlag{Name: seedFlag, Usage: "The seed for generating the world. The same seed produces the same world if used on the same version."},
			&cli.StringFlag{Name: townNameFlag, Usage: "A path to a file which holds the town names to use, one town name per line. If empty, 25 default town names are used."},
		},
		Commands: []*cli.Command{
			{
				Name:  "svg",
				Usage: "Generates and renders one tile as SVG onto the standard output.",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: sizeFlag, Usage: "The size of the rendered SVG in pixels.", DefaultText: "200"},
					&cli.IntFlag{Name: tileFlag, Usage: "The requested tile to be rendered.", DefaultText: "0"},
				},
				Action: renderSingleTile,
			},
			{
				Name:  "serve",
				Usage: "Serves a simple tile-server for exploring the world conveniently.",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: bindFlag, Usage: "The address to bind to.", Value: "127.0.0.1"},
					&cli.IntFlag{Name: portFlag, Usage: "The port to listen on.", Value: 9551},
					&cli.IntFlag{Name: shiftFlag, Usage: "Shifts the whole world by adding a constant to each tile request. See description on Github for more information.", Value: -131068},
				},
				Action: runServer,
			},
		},
	}).Run(os.Args)
	if err != nil {
		log.Fatalf("Program exited abnormally: %v", err)
	}
}

func renderSingleTile(context *cli.Context) error {
	defaultNoise, err := getGenerator(context)
	if err != nil {
		return err
	}
	hectometer := context.Int(tileFlag)
	size := context.Int(sizeFlag)
	if size == 0 {
		size = 200
	}
	tile := defaultNoise.Generate(hectometer)
	rn := renderer.New(os.Stdout, size)
	return fmt.Errorf("could not render tile %d: %v", hectometer, rn.Render(tile))
}

func getGenerator(context *cli.Context) (*noise.Generator, error) {
	seed := context.String(seedFlag)
	var err error
	if seed == "" {
		seed, err = noise.RandomSeed()
	}
	if err != nil {
		return nil, fmt.Errorf("could not generate random seed: %v", err)
	}
	defaultGenerator := noise.New(seed)
	townNameFile := context.String(townNameFlag)
	if townNameFile != "" {
		file, err := os.Open(townNameFile)
		if err != nil {
			return nil, fmt.Errorf("could not open town name file for reading: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			defaultGenerator.TownNames = append(defaultGenerator.TownNames, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("could not read from town name file: %v", err)
		}
	}
	return defaultGenerator, nil
}

func runServer(context *cli.Context) error {
	defaultNoise, err := getGenerator(context)
	if err != nil {
		return err
	}
	bind := context.String(bindFlag)
	port := context.Int(portFlag)
	shift := context.Int(shiftFlag)
	fmt.Printf("Server listening on http://127.0.0.1:%d\n", port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", bind, port), web.ApiHandler(defaultNoise, shift))
}
