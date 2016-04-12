package cli

import (
	"github.com/bitrise-tools/releaseman/releaseman"
	"github.com/codegangsta/cli"
)

const (
	// LogLevelEnvKey ...
	LogLevelEnvKey = "LOGLEVEL"
	// LogLevelKey ...
	LogLevelKey      = "loglevel"
	logLevelKeyShort = "l"

	// HelpKey ...
	HelpKey      = "help"
	helpKeyShort = "h"

	// VersionKey ...
	VersionKey      = "version"
	versionKeyShort = "v"

	// CIKey ...
	CIKey = "ci"
	// CIModeEnvKey ...
	CIModeEnvKey = "CI"

	// StartFromKey ...
	StartFromKey = "start_from"

	// ReleaseOnKey ...
	ReleaseOnKey = "release_on"

	// StartStateKey ...
	StartStateKey = "start-state"

	// EndStateKey ...
	EndStateKey = "end-state"

	// ChangelogPathKey ...
	ChangelogPathKey = "changelog-path"

	// BumpVersionKey ...
	BumpVersionKey = "bump-version"
	// PatchKey ...
	PatchKey = "patch"
	// MinorKey ...
	MinorKey = "minor"
	// MajorKey ...
	MajorKey = "major"

	// GetVersionScriptKey ...
	GetVersionScriptKey = "get-version-script"

	// SetVersionScriptKey ...
	SetVersionScriptKey = "set-version-script"

	// ConfigKey ...
	ConfigKey = "config"
)

var changelogFlags = []cli.Flag{
	cli.StringFlag{
		Name:  ConfigKey,
		Usage: "Release configuration file path.",
		Value: releaseman.DefaultConfigPth,
	},
	cli.StringFlag{
		Name:  StartFromKey,
		Usage: "Development branch",
	},
	cli.StringFlag{
		Name:  VersionKey,
		Usage: "Release version",
	},
	cli.StringFlag{
		Name:  BumpVersionKey,
		Value: "patch",
		Usage: "Bump version (options: patch, minor, major).",
	},
	cli.StringFlag{
		Name:  GetVersionScriptKey,
		Usage: "Script for getting current version.",
	},
	cli.StringFlag{
		Name:  SetVersionScriptKey,
		Usage: "Script for setting next version.",
	},
	cli.StringFlag{
		Name:  ChangelogPathKey,
		Usage: "Change log path",
	},
}

var releaseFlags = append(changelogFlags, cli.StringFlag{
	Name:  ReleaseOnKey,
	Usage: "Release branch",
})

var (
	commands = []cli.Command{
		{
			Name:   "prepare",
			Usage:  "Prepares the release of a new version",
			Action: prepare,
			Flags:  releaseFlags,
		},
		{
			Name:   "create",
			Usage:  "Create changelog and release new version",
			Action: create,
			Flags:  releaseFlags,
		},
		{
			Name:   "create-changelog",
			Usage:  "Create changelog",
			Action: createChangelog,
			Flags:  changelogFlags,
		},
		{
			Name:   "create-release",
			Usage:  "Release new version",
			Action: createRelease,
			Flags:  releaseFlags,
		},
		{
			Name:   "init",
			Usage:  "Initialize release configuration",
			Action: initRelease,
		},
	}

	appFlags = []cli.Flag{
		cli.StringFlag{
			Name:   LogLevelKey + ", " + logLevelKeyShort,
			Value:  "info",
			Usage:  "Log level (options: debug, info, warn, error, fatal, panic).",
			EnvVar: LogLevelEnvKey,
		},
		cli.BoolFlag{
			Name:   CIKey,
			Usage:  "If true it indicates that we're used by another tool so don't require any user input!",
			EnvVar: CIModeEnvKey,
		},
	}
)

func init() {
	// Override default help and version flags
	cli.HelpFlag = cli.BoolFlag{
		Name:  HelpKey + ", " + helpKeyShort,
		Usage: "Show help.",
	}

	cli.VersionFlag = cli.BoolFlag{
		Name:  VersionKey + ", " + versionKeyShort,
		Usage: "Print the version.",
	}
}
