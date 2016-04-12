package cli

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/depman/pathutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/releaseman/git"
	"github.com/bitrise-tools/releaseman/releaseman"
	"github.com/codegangsta/cli"
)

//=======================================
// Utility
//=======================================

func collectPrepareParams(config releaseman.Config, c *cli.Context) (releaseman.Config, error) {
	var err error

	//
	// Fill development branch
	if config, err = fillDevelopmetnBranch(config, c); err != nil {
		return releaseman.Config{}, err
	}

	//
	// Fill release branch
	if config, err = fillReleaseBranch(config, c); err != nil {
		return releaseman.Config{}, err
	}

	//
	// Fill release version
	if config, err = fillVersion(config, c); err != nil {
		return releaseman.Config{}, err
	}

	//
	// Fill changelog path
	if config, err = fillChangelogPath(config, c); err != nil {
		return releaseman.Config{}, err
	}

	return config, nil
}

func diffChanges(allChanges, beforeChanges []string) []string {
	changeMap := map[string]bool{}
	for _, beforeChange := range beforeChanges {
		changeMap[beforeChange] = true
	}

	diff := []string{}
	for _, change := range allChanges {
		_, found := changeMap[change]
		if !found {
			diff = append(diff, change)
		}
	}

	return diff
}

//=======================================
// Main
//=======================================

func prepare(c *cli.Context) {
	changes, err := git.GetChangedFiles()
	if err != nil {
		log.Fatalf("Failed to get changes, error: %s", err)
	}

	if len(changes) > 0 {
		log.Warn("There are uncommitted git changes:")
		for _, change := range changes {
			log.Warnf(" * %s", change)
		}

		if releaseman.IsCIMode {
			log.Fatal("Aborting preparing a release...")
		}

		fmt.Println()
		answer, err := goinp.AskForBoolWithDefault("Uncommitted git changes will be skipped by release generation.\nContinue preparing a release?", false)
		if err != nil {
			log.Fatalf("Failed to ask for bool, error: %s", err)
		}

		if !answer {
			log.Fatal("Aborted preparing a release...")
		}
	}
	fmt.Println()

	//
	// Build config
	config := releaseman.Config{}
	configPath := ""
	if c.IsSet(ConfigKey) {
		configPath = c.String(ConfigKey)
	} else {
		configPath = releaseman.DefaultConfigPth
	}

	if exist, err := pathutil.IsPathExists(configPath); err != nil {
		log.Warnf("Failed to check if path exist, error: %s", err)
	} else if exist {
		config, err = releaseman.NewConfigFromFile(configPath)
		if err != nil {
			log.Fatalf("Failed to parse release config at (%s), error: %s", configPath, err)
		}
	}

	config, err = collectPrepareParams(config, c)
	if err != nil {
		log.Fatalf("Failed to collect config params, error: %s", err)
	}

	//
	// Validate config
	config.Print(releaseman.FullMode)

	if !releaseman.IsCIMode {
		ok, err := goinp.AskForBoolWithDefault("Are you ready to continue the prepare?", true)
		if err != nil {
			log.Fatalf("Failed to ask for input, error: %s", err)
		}
		if !ok {
			log.Fatal("Aborted preparing a release...")
		}
	}

	//
	// Run set version script
	if c.IsSet(SetVersionScriptKey) {
		setVersionScript := c.String(SetVersionScriptKey)
		if err := runSetVersionScript(setVersionScript, config.Release.Version); err != nil {
			log.Fatalf("Failed to run set version script, error: %s", err)
		}
	}

	//
	// Generate Changelog
	generateChangelog(config)
	fmt.Println()

	newChanges, err := git.GetChangedFiles()
	if err != nil {
		log.Fatalf("Failed to get changes, error: %s", err)
	}

	diff := diffChanges(newChanges, changes)

	log.Warn("Following changes were made during the prepare:")
	for _, change := range diff {
		log.Warnf(" * %s", change)
	}

	fmt.Println()
	log.Info("Please check out the changes and if you are ready to finish the release, call:")
	fmt.Println("releaseman finish")

	config.Release.Changes = diff

	pth := releaseman.GetPreparedReleaseConfigPath()
	if err := releaseman.WriteConfigToFile(config, pth); err != nil {
		log.Fatalf("Failed to save config to file, error: %s", err)
	}
}
