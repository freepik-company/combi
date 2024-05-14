package run

import (
	"fmt"
	"log"
	"os"
	"time"

	"gcmerge/internal/config"
	"gcmerge/internal/encoding/libconfig"
	"gcmerge/internal/globals"

	"github.com/spf13/cobra"
)

const (
	descriptionShort = `Execute synchronization process`
	descriptionLong  = `
	Run execute synchronization process`

	// Flags error messages
	logLevelFlagErrorMessage     = "unable to get flag --log-level: %s"
	disableTraceFlagErrorMessage = "unable to get flag --disable-trace: %s"
	syncTimeFlagErrorMessage     = "unable to get flag --sync-time: %s"
	durationParseErrorMessage    = "unable to parse duration: %s"
	configNameFlagErrorMessage   = "unable to get flag --config-name: %s"

	// Execution flow error messages
	// getConfigConfigMapErrorMessage     = "unable to get configuration configmap { ns: %s, name: %s }: %s"
	// configConfigMapDataKeyErrorMessage = "no key '%s' in configuration ConfigMap { ns: %s, name: %s }"
	// configParseErrorMessage            = "unable to parse configuration: %s"
	// getSourceConfigMapErrorMessage     = "unable to get source configmap { ns: %s, name: %s }: %s"
	// sourceConfigMapDataKeyErrorMessage = "no key '%s' in source ConfigMap { ns: %s, name: %s }"
	// getTagetConfigMapErrorMessage      = "unable to get target configmap { ns: %s, name: %s }: %s"
	// createTargetConfigMapErrorMessage  = "unable to create target configmap { ns: %s, name: %s }: %s"
	// targetConfigMapDataKeyErrorMessage = "no key '%s' in source ConfigMap { ns: %s, name: %s }"
	// targetUpdateErrorMessage           = "unable to update target configmap { ns: %s, name: %s }: %s"
	// workloadPatchErrorMessage          = "unable to path workload { kind: %s, ns: %s, name: %s }: %s"
	// workloadKindErrorMessage           = "workload { kind: %s, ns: %s, name: %s } resource with not supported type"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "run",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  descriptionLong,

		Run: RunCommand,
	}

	//
	cmd.Flags().String("log-level", "info", "Verbosity level for logs")
	cmd.Flags().Bool("disable-trace", true, "Disable showing traces in logs")
	cmd.Flags().String("sync-time", "15s", "Waiting time between group synchronizations (in duration type)")
	cmd.Flags().String("config-name", "example1", "Configuration name in the configuration list")

	return cmd
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {

	// Init the logger with logger flags
	logLevelFlag, err := cmd.Flags().GetString("log-level")
	if err != nil {
		log.Fatalf(logLevelFlagErrorMessage, err)
	}

	disableTraceFlag, err := cmd.Flags().GetBool("disable-trace")
	if err != nil {
		log.Fatalf(disableTraceFlagErrorMessage, err)
	}

	err = globals.SetLogger(logLevelFlag, disableTraceFlag)
	if err != nil {
		log.Fatal(err)
	}

	// TODO

	syncTime, err := cmd.Flags().GetString("sync-time")
	if err != nil {
		globals.ExecContext.Logger.Fatalf(syncTimeFlagErrorMessage, err)
	}

	duration, err := time.ParseDuration(syncTime)
	if err != nil {
		globals.ExecContext.Logger.Fatalf(durationParseErrorMessage, err)
	}

	// TODO
	configName, err := cmd.Flags().GetString("config-name")
	if err != nil {
		globals.ExecContext.Logger.Fatalf(configNameFlagErrorMessage, err)
	}

	/////////////////////////////
	// EXECUTION FLOW RELATED
	/////////////////////////////

	//TODO: Set git repository client to get the gcmerge config file
	// client := git.NewGitHubClient()

	for {
		// sync in x duration again
		globals.ExecContext.Logger.Infof("Syncing in %s", duration.String())
		time.Sleep(duration)

		// TODO: Get gcmerge config file from git repository (store in local /tmp/gcmerge/gitcongig.yaml)
		// TODO: Compare current gcmerge config (/var/lib/gcmerge/gitconfig.yaml) file with download one to decide make the changes (make it if first time)

		gcmFullConfigBytes, err := os.ReadFile("./config/gcmerge.yaml")
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to get gcmerge config file: %s", err.Error())
		}

		// TODO: Parse gcmerge config
		// TODO: Get specific config from config name flag in gcmerge config

		gcmFullConfig, err := config.Parse(gcmFullConfigBytes)
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to parse gcmerge configuration file: %s", err.Error())
			continue
		}

		gcmConfig, ok := gcmFullConfig.Configs[configName] // TODO: replace this line with specific name in flag
		if !ok {
			globals.ExecContext.Logger.Errorf("unable to get '%s' configuration in gcmerge configuration file: %s", configName)
			continue
		}

		// TODO: expand env variables in rawConfig local and global
		gcmFullConfig.Global.RawConfig = os.ExpandEnv(gcmFullConfig.Global.RawConfig)
		gcmConfig.RawConfig = os.ExpandEnv(gcmConfig.RawConfig)

		// TODO: Decode rawConfig and targetConfig
		// TODO: Merge rawConfig and targetConfig
		// TODO: Check config conditions

		switch gcmFullConfig.Global.Type {
		case "libconfig":
			{
				configDestination, err := libconfig.DecodeConfig(gcmConfig.TargetConfig)
				if err != nil {
					globals.ExecContext.Logger.Errorf("unable to decode target '%s' configuration file: %s", gcmConfig.TargetConfig, err.Error())
					continue
				}
				configSource, err := libconfig.DecodeConfigBytes([]byte(gcmConfig.RawConfig))
				if err != nil {
					globals.ExecContext.Logger.Errorf("unable to decode '%s' raw configuration field: %s", configName, err.Error())
					continue
				}
				libconfig.MergeConfigs(configDestination, configSource)
				configStr := libconfig.EncodeConfigString(configDestination)
				fmt.Println(configStr)
			}
		default:
			{
				globals.ExecContext.Logger.Errorf("unsuported configuration type: %s", gcmFullConfig.Global.Type)
				continue
			}
		}

		// TODO: Decode global rawConfig
		// TODO: Merge global rawConfig and current merged config
		// TODO: Check global config conditions

		// TODO: Update targetConfig with merged config file

		// TODO: Make config actions
		// TODO: Make global config actions
	}
}
