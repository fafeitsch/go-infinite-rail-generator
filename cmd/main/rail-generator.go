package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/fafeitsch/go-infinite-rail-generator/standalone"
	"github.com/fafeitsch/go-infinite-rail-generator/web"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
)

var (
	BuildVersion = "Development Snapshot"
	BuildTime    = ""
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
		Version:         fmt.Sprintf("%s (%s)", BuildVersion, BuildTime),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  seedFlag,
				Usage: "The seed for generating the world. The same seed produces the same world if used on the same version.",
			},
			&cli.StringFlag{
				Name:  townNameFlag,
				Usage: "A path to a file which holds the town names to use, one town name per line. If empty, 25 default town names are used.",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "svg",
				Usage: "Generates and renders one tile as SVG onto the standard output.",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        sizeFlag,
						Usage:       "The size of the rendered SVG in pixels.",
						DefaultText: "200",
					},
					&cli.IntFlag{
						Name:        tileFlag,
						Usage:       "The requested tile to be rendered.",
						DefaultText: "0",
					},
				},
				Action: renderSingleTile,
			},
			{
				Name:  "serve",
				Usage: "Serves a simple tile-server for exploring the world conveniently.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  bindFlag,
						Usage: "The address to bind to.",
						Value: "127.0.0.1",
					},
					&cli.IntFlag{
						Name:  portFlag,
						Usage: "The port to listen on.",
						Value: 9551,
					},
					&cli.IntFlag{
						Name:  shiftFlag,
						Usage: "Shifts the whole world by adding a constant to each tile request. See description on Github for more information.",
						Value: -131068,
					},
				},
				Action: runServer,
			},
		},
	}).Run(os.Args)
	if err != nil {
		log.Fatalf("Program exited abnormally: %v", err)
	}
}

func readTownNames(townNameFile string) ([]string, error) {
	result := make([]string, 0, 0)
	if townNameFile != "" {
		file, err := os.Open(townNameFile)
		if err != nil {
			return nil, fmt.Errorf("could not open town name file for reading: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			result = append(result, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("could not read from town name file: %v", err)
		}
	}
	return result, nil
}

func renderSingleTile(context *cli.Context) error {
	seed := context.String(seedFlag)
	if seed == "" {
		seed, err := randomSeed()
		if err != nil {
			return fmt.Errorf("could not create random seed: %v", seed)
		}
	}
	townNames, err := readTownNames(context.String(townNameFlag))
	if err != nil {
		return err
	}
	options := standalone.RenderOptions{
		Seed:       context.String(seedFlag),
		TownNames:  townNames,
		Hectometer: context.Int(tileFlag),
		Size:       context.Int(sizeFlag),
	}
	return standalone.RenderSingleTile(options)
}

func runServer(context *cli.Context) error {
	seed := context.String(seedFlag)
	if seed == "" {
		seed, err := randomSeed()
		if err != nil {
			return fmt.Errorf("could not create random seed: %v", seed)
		}
	}
	townNames, err := readTownNames(context.String(townNameFlag))
	if err != nil {
		return err
	}
	options := web.ApiOptions{
		Shift:        context.Int(shiftFlag),
		Seed:         context.String(seedFlag),
		TownNames:    townNames,
		BuildTime:    BuildTime,
		BuildVersion: BuildVersion,
	}
	port := context.Int(portFlag)
	fmt.Printf("Server listening on http://127.0.0.1:%d\n", port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", context.String(bindFlag), port), web.ApiHandler(options))
}

func randomSeed() (string, error) {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
