package daemon

import (
	"os"
	"time"

	"gcmerge/internal/actions"
	"gcmerge/internal/config"
	"gcmerge/internal/encoding"
	"gcmerge/internal/flags"
	"gcmerge/internal/globals"
	"gcmerge/internal/source"

	"github.com/spf13/cobra"
)

const (
	descriptionShort = `Execute synchronization process`
	descriptionLong  = `
	Run execute synchronization process`

	getFlagsErrMsg      = "unable to get flags: %s"
	durationParseErrMsg = "unable to parse duration: %s"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "daemon",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  descriptionLong,

		Run: RunCommand,
	}

	cmd.Flags().String(flags.LogLevelFlagName, "info", "Verbosity level for logs")
	cmd.Flags().Bool(flags.DisableTraceFlagName, true, "Disable showing traces in logs")

	cmd.Flags().String(flags.TmpDirFlagName, "/tmp/combi", "Verbosity level for logs")
	cmd.Flags().String(flags.SyncTimeFlagName, "15s", "Waiting time between group synchronizations (in duration type)")
	cmd.Flags().String(flags.SourceTypeFlagName, "git", "Source where find source config")
	cmd.Flags().String(flags.SourcePathFlagName, "config/gcmerge.yaml", "Source path where find source config")
	cmd.Flags().String(flags.SourceFieldFlagName, "example1", "Field in source config map to find mergeble config")

	//
	cmd.Flags().String(flags.GitSshUrlFlagName, "git@github.com:sebastocorp/gcmerge.git", "Git repository ssh url")
	cmd.Flags().String(flags.GitSshKeyFilepathFlagName, "/home/svargas/.ssh/id_rsa_github", "Ssh private key filepath for git repository")
	cmd.Flags().String(flags.GitBranchFlagName, "main", "Git branch repository")

	return cmd
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {

	cmdFlags, err := flags.GetSyncRunFlags(cmd)
	if err != nil {
		globals.ExecContext.Logger.Fatalf(getFlagsErrMsg, err.Error())
	}

	duration, err := time.ParseDuration(cmdFlags.SyncTime)
	if err != nil {
		globals.ExecContext.Logger.Fatalf(durationParseErrMsg, err)
	}

	/////////////////////////////
	// EXECUTION FLOW RELATED
	/////////////////////////////

	if err = os.MkdirAll(cmdFlags.TmpDir, 0744); err != nil {
		globals.ExecContext.Logger.Fatalf("unable to create '%s' tmp dir: %s", cmdFlags.TmpDir, err)
	}

	src, ok := source.GetSources()[cmdFlags.SourceType]
	if !ok {
		globals.ExecContext.Logger.Fatalf("unsuported source type: %s", cmdFlags.SourceType)
	}

	src.Init(cmdFlags)

	firstLoop := true
	for {
		// sync in x duration again
		if !firstLoop {
			globals.ExecContext.Logger.Infof("Syncing in %s", duration.String())
			time.Sleep(duration)
		}
		firstLoop = false

		combiFullConfigBytes, err := src.GetConfig()
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to get source config: %s", err.Error())
			continue
		}

		if !src.NeedUpdate() {
			continue
		}

		// Parse config
		combiFullConfig, err := config.Parse(combiFullConfigBytes)
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to parse source config: %s", err.Error())
			continue
		}

		combiLocalConfig, ok := combiFullConfig.Configs[cmdFlags.SourceField]
		if !ok {
			globals.ExecContext.Logger.Errorf("unable to get '%s' local config in source config file: %s", cmdFlags.SourceField)
			continue
		}

		// Expand env variables in local and global rawConfig
		combiFullConfig.Global.RawConfig = os.ExpandEnv(combiFullConfig.Global.RawConfig)
		combiLocalConfig.RawConfig = os.ExpandEnv(combiLocalConfig.RawConfig)

		_, ok = encoding.GetEncoders()[combiFullConfig.Kind]
		if !ok {
			globals.ExecContext.Logger.Errorf("unsuported config type: %s", combiFullConfig.Kind)
			continue
		}

		targetEncoder, err := mergeConfigurations(combiFullConfig, cmdFlags.SourceField) // TODO: fix this function
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to merge configs: %s", err.Error())
			continue
		}

		targetConfigMap := targetEncoder.ConfigToMap()
		result, err := evaluateConditions(combiFullConfig, cmdFlags.SourceField, targetConfigMap)
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to evaluate conditions: %s", err.Error())
			continue
		}

		mergedConfigStr := targetEncoder.EncodeConfigString()

		// Execute local+global actions
		if result {
			// Update targetConfig with merged config file
			err = os.WriteFile(combiLocalConfig.MergedConfig, []byte(mergedConfigStr), 0644)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to create '%s' merged config file: %s", combiLocalConfig.MergedConfig, err.Error())
				continue
			}

			err = actions.RunActions(&combiLocalConfig.Actions.OnSuccess)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to execute local success actions: %s", err.Error())
			}

			err = actions.RunActions(&combiFullConfig.Global.Actions.OnSuccess)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to execute global success actions: %s", err.Error())
			}
		} else {
			err = actions.RunActions(&combiLocalConfig.Actions.OnFailure)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to execute local failure actions: %s", err.Error())
			}

			err = actions.RunActions(&combiFullConfig.Global.Actions.OnFailure)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to execute global failure actions: %s", err.Error())
			}
		}
	}
}