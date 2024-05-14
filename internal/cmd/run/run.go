package run

import (
	"log"
	"os"
	"os/exec"
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

		// Parse gcmerge config
		gcmFullConfig, err := config.Parse(gcmFullConfigBytes)
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to parse gcmerge configuration file: %s", err.Error())
			continue
		}

		gcmGlobalConfig := gcmFullConfig.Global
		gcmConfig, ok := gcmFullConfig.Configs[configName]
		if !ok {
			globals.ExecContext.Logger.Errorf("unable to get '%s' configuration in gcmerge configuration file: %s", configName)
			continue
		}

		// Expand env variables in local and global rawConfig
		gcmGlobalConfig.RawConfig = os.ExpandEnv(gcmGlobalConfig.RawConfig)
		gcmConfig.RawConfig = os.ExpandEnv(gcmConfig.RawConfig)

		var mergedConfigStr string
		switch gcmFullConfig.Kind {
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
				// TODO: Check config conditions

				configGlobal, err := libconfig.DecodeConfigBytes([]byte(gcmGlobalConfig.RawConfig))
				if err != nil {
					globals.ExecContext.Logger.Errorf("unable to decode global raw configuration field: %s", err.Error())
					continue
				}
				libconfig.MergeConfigs(configDestination, configGlobal)
				// TODO: Check global config conditions

				mergedConfigStr = libconfig.EncodeConfigString(configDestination)
			}
		default:
			{
				globals.ExecContext.Logger.Errorf("unsuported configuration type: %s", gcmFullConfig.Kind)
				continue
			}
		}

		// Update targetConfig with merged config file
		err = os.WriteFile(gcmConfig.MergedConfig, []byte(mergedConfigStr), 0644)
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to create '%s' merged configuration file: %s", gcmConfig.MergedConfig, err.Error())
			continue
		}

		// Execute config actions
		for _, action := range gcmConfig.Actions {
			command := exec.Command(action.Command[0], action.Command[1:]...)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			if err = command.Run(); err != nil {
				globals.ExecContext.Logger.Errorf("unable to execute config action '%s': %s", action.Name, err.Error())
			}
		}
		// Execute global actions
		for _, action := range gcmGlobalConfig.Actions {
			command := exec.Command(action.Command[0], action.Command[1:]...)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			if err = command.Run(); err != nil {
				globals.ExecContext.Logger.Errorf("unable to execute global action '%s': %s", action.Name, err.Error())
			}
		}
	}
}
