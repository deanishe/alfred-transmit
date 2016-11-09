//
// Copyright (c) 2016 Dean Jackson <deanishe@deanishe.net>
//
// MIT Licence. See http://opensource.org/licenses/MIT
//
// Created on 2016-11-06
//

package main

import (
	"fmt"
	"log"

	"os"

	"strings"

	"github.com/deanishe/alfred-transmit"
	"github.com/docopt/docopt-go"
	"gogs.deanishe.net/deanishe/awgo"
)

var usage = `transmit [options] [<query>]

Search and open Transmit favourites.

Usage:
    transmit [<query>]
    transmit -u
    transmit --distname

Options:
    -u, --update    Check for workflow update.
    --distname      Print name of workflow file (for building).
    -h, --help      Show this help message.
    --version       Show version.
`

var (
	iconUpdate = &aw.Icon{Value: "update.png"}
	favesPath  string
	// Command-line flags
	doUpdate bool
	doDist   bool
	query    string
	// Workflow configuration
	opts *aw.Options
	wf   *aw.Workflow
)

func init() {
	// Command-line flags
	opts = &aw.Options{
		GitHub: "deanishe/alfred-transmit",
	}
	wf = aw.NewWorkflow(opts)
}

func parseArgs() {
	args, _ := docopt.Parse(usage, wf.Args(), true, wf.Version(), false)
	// log.Printf("args=%#v", args)
	if args["--update"] != false {
		doUpdate = true
	}
	if args["--distname"] != false {
		doDist = true
	}
	if args["<query>"] != nil {
		query = args["<query>"].(string)
	}
}

// Check for updates
func runUpdate() {
	wf.TextErrors = true
	log.Println("Checking for update...")
	if err := wf.CheckForUpdate(); err != nil {
		wf.FatalError(err)
	}
}

// Print workflow file name
func runDist() {
	name := strings.Replace(
		fmt.Sprintf("%s-%s.alfredworkflow", wf.Name(), wf.Version()),
		" ", "-", -1)
	fmt.Print(name)
}

func run() {
	parseArgs()

	// ----------------------------------------------------------------
	// Check for updates if --update was passed
	if doUpdate {
		runUpdate()
		return
	}

	// ----------------------------------------------------------------
	// Print distname if --distname was passed
	if doDist {
		runDist()
		return
	}

	// ----------------------------------------------------------------
	// Updates
	var noUID bool
	if wf.UpdateCheckDue() {
		wf.Var("check_update", "1")
	}
	// Inform of update if query is empty
	if wf.UpdateAvailable() && query == "" {
		wf.NewItem("An update is available").
			Subtitle("↩ or ⇥ to install update").
			Valid(false).
			Autocomplete("workflow:update").
			Icon(iconUpdate)
		noUID = true
	}

	// ----------------------------------------------------------------
	// Load data
	favesPath = strings.Replace(os.Getenv("FAVES_PATH"), "~", os.Getenv("HOME"), 1)
	transmit.BookmarksPath = favesPath
	log.Printf("Favorites.xml: %s", favesPath)
	log.Printf("query=%s", query)

	faves, err := transmit.Favourites()
	if err != nil {
		log.Printf("Error fetching favourites: %s", err)
		wf.FatalError(err)
	}
	log.Printf("%d favourites.", len(faves))

	// ----------------------------------------------------------------
	// Display results
	for _, fav := range faves {
		it := wf.NewItem(fav.Name).
			Subtitle(fav.Server).
			Arg(fav.ID).
			Valid(true).
			SortKey(fmt.Sprintf("%s %s", fav.Name, fav.Server))
		if !noUID {
			it.UID(fav.ID)
		}
	}

	if query != "" {
		res := wf.Filter(query)
		for i, r := range res {
			log.Printf("%02d. %4.2f %s", i+1, r.Score, r.SortKey)
		}
	}

	wf.WarnEmpty("No favourites found", "Try a different query?")
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
