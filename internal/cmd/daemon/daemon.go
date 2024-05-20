package daemon

import (
	"os"
	"time"

	"combi/internal/actions"
	"combi/internal/config"
	"combi/internal/encoding"
	"combi/internal/flags"
	"combi/internal/globals"
	"combi/internal/source"

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

	cmd.Flags().String(flags.TmpDirFlagName, "/tmp/combi", "Temporary directoty to store temporary objects like remote repos, scripts, etc")
	cmd.Flags().String(flags.SyncTimeFlagName, "15s", "Waiting time between source synchronizations (in duration type)")
	cmd.Flags().String(flags.SourceTypeFlagName, "git", "Source where consume the combi config")
	cmd.Flags().String(flags.SourcePathFlagName, "config/combi.yaml", "Path in source where find combi config")
	cmd.Flags().String(flags.SourceFieldFlagName, "example1", "Field in combi config map to find the mergeble config")

	//
	cmd.Flags().String(flags.GitSshUrlFlagName, "git@github.com:example/project.git", "Git repository ssh url")
	cmd.Flags().String(flags.GitSshKeyFilepathFlagName, "/home/example/.ssh/id_rsa_github", "Ssh private key filepath for git repository")
	cmd.Flags().String(flags.GitBranchFlagName, "main", "Git branch repository")

	return cmd
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {

	cmdFlags, err := flags.GetDaemonFlags(cmd)
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

		globals.ExecContext.Logger.Infof("get configurations from source")
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

		globals.ExecContext.Logger.Infof("merge configurations from source to target")
		targetEncoder, err := mergeConfigurations(combiFullConfig, cmdFlags.SourceField) // TODO: fix this function
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to merge configs: %s", err.Error())
			continue
		}

		globals.ExecContext.Logger.Infof("evaluate confitions")
		targetConfigMap := targetEncoder.ConfigToMap()
		result, err := evaluateConditions(combiFullConfig, cmdFlags.SourceField, targetConfigMap)
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to evaluate conditions: %s", err.Error())
			continue
		}

		// Execute local+global actions
		if result {
			globals.ExecContext.Logger.Infof("create final merged configuration file")
			mergedConfigStr := targetEncoder.EncodeConfigString()
			// Update targetConfig with merged config file
			err = os.WriteFile(combiLocalConfig.MergedConfig, []byte(mergedConfigStr), 0744)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to create '%s' merged config file: %s", combiLocalConfig.MergedConfig, err.Error())
				continue
			}

			globals.ExecContext.Logger.Infof("execute success actions")
			err = actions.RunActions(&combiLocalConfig.Actions.OnSuccess)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to execute local success actions: %s", err.Error())
			}

			err = actions.RunActions(&combiFullConfig.Global.Actions.OnSuccess)
			if err != nil {
				globals.ExecContext.Logger.Errorf("unable to execute global success actions: %s", err.Error())
			}
		} else {
			globals.ExecContext.Logger.Infof("execute failure actions")
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
