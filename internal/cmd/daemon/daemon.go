package daemon

import (
	"os"
	"time"

	"gcmerge/internal/actions"
	"gcmerge/internal/conditions"
	"gcmerge/internal/config"
	"gcmerge/internal/encoding/libconfig"
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
	cmd.Flags().String(flags.SourceTypeFlagName, "git", "Source where find cmerged configuration")
	cmd.Flags().String(flags.SourcePathFlagName, "config/gcmerge.yaml", "Source path where find cmerged configuration")
	cmd.Flags().String(flags.SourceFieldFlagName, "example1", "Field in cmerged configuration map to find mergeble config")

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

		gcmFullConfigBytes, err := src.GetConfig()
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to get git gcmerge configuration: %s", err.Error())
			continue
		}

		if !src.NeedUpdate() {
			continue
		}

		// Parse gcmerge config
		gcmFullConfig, err := config.Parse(gcmFullConfigBytes)
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to parse gcmerge config file: %s", err.Error())
			continue
		}

		gcmGlobalConfig := gcmFullConfig.Global
		gcmLocalConfig, ok := gcmFullConfig.Configs[cmdFlags.SourceField]
		if !ok {
			globals.ExecContext.Logger.Errorf("unable to get '%s' local configuration in gcmerge config file: %s", cmdFlags.SourceField)
			continue
		}

		// Expand env variables in local and global rawConfig
		gcmGlobalConfig.RawConfig = os.ExpandEnv(gcmGlobalConfig.RawConfig)
		gcmLocalConfig.RawConfig = os.ExpandEnv(gcmLocalConfig.RawConfig)

		var mergedConfigStr string
		var evalLocalConditionsSuccess bool
		var evalGlobalConditionsSuccess bool
		switch gcmFullConfig.Kind {
		case "libconfig":
			{
				targetConfig, err := libconfig.DecodeConfig(gcmLocalConfig.TargetConfig)
				if err != nil {
					globals.ExecContext.Logger.Errorf("unable to decode '%s' target config file: %s", gcmLocalConfig.TargetConfig, err.Error())
					continue
				}
				localRawConfig, err := libconfig.DecodeConfigBytes([]byte(gcmLocalConfig.RawConfig))
				if err != nil {
					globals.ExecContext.Logger.Errorf("unable to decode '%s' local raw config field: %s", cmdFlags.SourceField, err.Error())
					continue
				}
				libconfig.MergeConfigs(targetConfig, localRawConfig)

				globalRawConfig, err := libconfig.DecodeConfigBytes([]byte(gcmGlobalConfig.RawConfig))
				if err != nil {
					globals.ExecContext.Logger.Errorf("unable to decode global raw config field: %s", err.Error())
					continue
				}
				libconfig.MergeConfigs(targetConfig, globalRawConfig)

				// Check local config conditions
				targetConfigMap := libconfig.ConfigToMap(targetConfig)
				evalLocalConditionsSuccess, err = conditions.EvalConditions(&gcmLocalConfig.Conditions, &targetConfigMap)
				if err != nil {
					globals.ExecContext.Logger.Errorf("unable to evaluate local conditions: %s", err.Error())
				}

				// Check global config conditions
				evalGlobalConditionsSuccess, err = conditions.EvalConditions(&gcmGlobalConfig.Conditions, &targetConfigMap)
				if err != nil {
					globals.ExecContext.Logger.Errorf("unable to evaluate global conditions: %s", err.Error())
				}

				mergedConfigStr = libconfig.EncodeConfigString(targetConfig)
			}
		default:
			{
				globals.ExecContext.Logger.Errorf("unsuported configuration type: %s", gcmFullConfig.Kind)
				continue
			}
		}

		// Execute local+global actions
		if evalLocalConditionsSuccess && evalGlobalConditionsSuccess {
			// Update targetConfig with merged config file
			err = os.WriteFile(gcmLocalConfig.MergedConfig, []byte(mergedConfigStr), 0644)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to create '%s' merged config file: %s", gcmLocalConfig.MergedConfig, err.Error())
				continue
			}

			err = actions.RunActions(&gcmLocalConfig.Actions.OnSuccess)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to execute local success actions: %s", err.Error())
			}

			err = actions.RunActions(&gcmGlobalConfig.Actions.OnSuccess)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to execute global success actions: %s", err.Error())
			}
		} else {
			err = actions.RunActions(&gcmLocalConfig.Actions.OnFailure)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to execute local failure actions: %s", err.Error())
			}

			err = actions.RunActions(&gcmGlobalConfig.Actions.OnFailure)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to execute global failure actions: %s", err.Error())
			}
		}
	}
}
