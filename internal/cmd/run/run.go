package run

import (
	"log"
	"os"
	"time"

	"gcmerge/internal/actions"
	"gcmerge/internal/conditions"
	"gcmerge/internal/config"
	"gcmerge/internal/encoding/libconfig"
	"gcmerge/internal/git"
	"gcmerge/internal/globals"

	"github.com/spf13/cobra"
)

const (
	descriptionShort = `Execute synchronization process`
	descriptionLong  = `
	Run execute synchronization process`

	// Flags error messages
	logLevelFlagErrorMessage        = "unable to get flag --log-level: %s"
	disableTraceFlagErrorMessage    = "unable to get flag --disable-trace: %s"
	syncTimeFlagErrorMessage        = "unable to get flag --sync-time: %s"
	durationParseErrorMessage       = "unable to parse duration: %s"
	configNameFieldFlagErrorMessage = "unable to get flag --config-name: %s"

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
	cmd.Flags().String("config-field", "example1", "Configuration name in the gcmerge configuration map")
	cmd.Flags().String("source-type", "git", "Source where find gcmerge configuration")
	cmd.Flags().String("source-config-filepath", "config/gcmerge.yaml", "Source gcmerge configuration filepath")
	cmd.Flags().String("git-ssh-url", "git@github.com:sebastocorp/gcmerge.git", "Git repository ssh url")
	cmd.Flags().String("git-sshkey-filepath", "/home/svargas/.ssh/id_rsa_github", "Ssh private key filepath for git repository")
	cmd.Flags().String("git-branch", "main", "Git branch repository")

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
	configNameField, err := cmd.Flags().GetString("config-field")
	if err != nil {
		globals.ExecContext.Logger.Fatalf(configNameFieldFlagErrorMessage, err)
	}

	sourceType, err := cmd.Flags().GetString("source-type")
	if err != nil {
		globals.ExecContext.Logger.Fatalf("%s", err)
	}

	sourceConfigFilepath, err := cmd.Flags().GetString("source-config-filepath")
	if err != nil {
		globals.ExecContext.Logger.Fatalf("%s", err)
	}

	gitSshUrl, err := cmd.Flags().GetString("git-ssh-url")
	if err != nil {
		globals.ExecContext.Logger.Fatalf("%s", err)
	}

	gitSshKeyFilepath, err := cmd.Flags().GetString("git-sshkey-filepath")
	if err != nil {
		globals.ExecContext.Logger.Fatalf("%s", err)
	}

	gitBranch, err := cmd.Flags().GetString("git-branch")
	if err != nil {
		globals.ExecContext.Logger.Fatalf("%s", err)
	}

	/////////////////////////////
	// EXECUTION FLOW RELATED
	/////////////////////////////

	gcmTmpPath := "/tmp/gcmerge"
	if err = os.MkdirAll(gcmTmpPath, 0744); err != nil {
		globals.ExecContext.Logger.Fatalf("unable to create '%s' tmp dir: %s", gcmTmpPath, err)
	}

	source := git.Git{
		SshKeyFilepath:     gitSshKeyFilepath,
		RepoSshUrl:         gitSshUrl,
		RepoBranch:         gitBranch,
		RepoPath:           gcmTmpPath + "/repo",
		RepoConfigFilepath: sourceConfigFilepath,
	}

	firstLoop := true
	for {
		// sync in x duration again
		if !firstLoop {
			globals.ExecContext.Logger.Infof("Syncing in %s", duration.String())
			time.Sleep(duration)
		}
		firstLoop = false

		var gcmFullConfigBytes []byte
		switch sourceType {
		case "git":
			{
				gcmFullConfigBytes, err = source.GetConfig()
				if err != nil {
					globals.ExecContext.Logger.Errorf("unable to get git gcmerge configuration: %s", err.Error())
					continue
				}

				if !source.NeedUpdate() {
					continue
				}
			}
		case "local":
			{
				gcmFullConfigBytes, err = os.ReadFile(sourceConfigFilepath)
				if err != nil {
					globals.ExecContext.Logger.Errorf("unable to get local gcmerge configuration: %s", err.Error())
					continue
				}
			}
		default:
			{
				globals.ExecContext.Logger.Errorf("unsuported source type: %s", sourceType)
				continue
			}
		}

		// Parse gcmerge config
		gcmFullConfig, err := config.Parse(gcmFullConfigBytes)
		if err != nil {
			globals.ExecContext.Logger.Errorf("unable to parse gcmerge config file: %s", err.Error())
			continue
		}

		gcmGlobalConfig := gcmFullConfig.Global
		gcmLocalConfig, ok := gcmFullConfig.Configs[configNameField]
		if !ok {
			globals.ExecContext.Logger.Errorf("unable to get '%s' local configuration in gcmerge config file: %s", configNameField)
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
					globals.ExecContext.Logger.Errorf("unable to decode '%s' local raw config field: %s", configNameField, err.Error())
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
