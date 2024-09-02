package daemon

import (
	"os"
	"reflect"
	"time"

	"combi/internal/actions"
	"combi/internal/config"
	"combi/internal/encoding"
	"combi/internal/logger"
	"combi/internal/source"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "daemon",
		DisableFlagsInUseLine: true,
		Short:                 `Execute synchronization process`,
		Long: `
	Run execute synchronization process`,

		Run: RunCommand,
	}

	cmd.Flags().String(logLevelFlagName, "info", "Verbosity level for logs")

	cmd.Flags().String(tmpDirFlagName, "/tmp/combi", "temporary directoty to store temporary objects like remote repos, scripts, etc")
	cmd.Flags().Duration(syncTimeFlagName, 15*time.Second, "waiting time between source synchronizations (in duration type)")
	cmd.Flags().String(configFlagName, "combi.yaml", "combi configuration file")

	return cmd
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {

	daemon := daemonT{}
	err := daemon.getFlags(cmd)
	if err != nil {
		logger.Log.Fatalf("unable to parse flags: %s", err.Error())
	}

	if err = os.MkdirAll(daemon.flags.tmpDir, 0744); err != nil {
		logger.Log.Fatalf("unable to create '%s' tmp dir: %s", daemon.flags.tmpDir, err)
	}

	configBytes, err := os.ReadFile(daemon.flags.config)
	if err != nil {
		logger.Log.Fatalf("unable to get combi configuration: %s", err.Error())
	}

	if daemon.config, err = config.Parse(configBytes); err != nil {
		logger.Log.Fatalf("unable to parse combi configuration: %s", err.Error())
	}

	/////////////////////////////
	// EXECUTION FLOW RELATED
	/////////////////////////////

	if _, encoderOk := encoding.GetEncoders()[daemon.config.Kind]; !encoderOk {
		logger.Log.Fatalf("unsuported config kind '%s'", daemon.config.Kind)
	}

	logger.Log.Infof("init configs from defined sources")

	if reflect.ValueOf(daemon.config.Global.Source).IsZero() {
		daemon.config.Global.Source.Type = "raw"
		daemon.config.Global.Source.RawConfig = ""
	}

	globalSource, ok := source.GetSources()[daemon.config.Global.Source.Type]
	if !ok {
		logger.Log.Fatalf("unsupported global source type '%s'", daemon.config.Global.Source.Type)
	}
	globalSource.Init(daemon.config.Global.Source)

	configSource, ok := source.GetSources()[daemon.config.Config.Source.Type]
	if !ok {
		logger.Log.Fatalf("unsupported config source type '%s'", daemon.config.Config.Source.Type)
	}
	configSource.Init(daemon.config.Config.Source)

	logger.Log.Infof("init synchronization loop")

	firstLoop := true
	for {
		// sync in x duration again
		if !firstLoop {
			logger.Log.Infof("Syncing in %s", daemon.flags.syncTime.String())
			time.Sleep(daemon.flags.syncTime)
		}
		firstLoop = false

		logger.Log.Infof("get configurations from sources")

		globalConfigBytes, globalConfigUpdated, err := globalSource.GetConfig()
		if err != nil {
			logger.Log.Errorf("unable to get global source config: %s", err.Error())
			continue
		}

		configBytes, configUpdated, err := configSource.GetConfig()
		if err != nil {
			logger.Log.Errorf("unable to get source config: %s", err.Error())
			continue
		}

		if !globalConfigUpdated && !configUpdated {
			continue
		}

		configEncoder, err := daemon.mergeConfigurations(globalConfigBytes, configBytes)
		if err != nil {
			logger.Log.Errorf("unable to merge configurations: %s", err.Error())
			continue
		}

		logger.Log.Infof("evaluate confitions")
		targetConfigMap := configEncoder.ConfigToMap()
		success, err := daemon.evaluateConditions(targetConfigMap)
		if err != nil {
			logger.Log.Errorf("unable to evaluate conditions: %s", err.Error())
			continue
		}

		if !success {
			logger.Log.Infof("execute failure actions")

			err = actions.RunActions(&daemon.config.Global.Actions.OnFailure)
			if err != nil {
				logger.Log.Errorf("unable to execute global failure actions: %s", err.Error())
			}

			err = actions.RunActions(&daemon.config.Config.Actions.OnFailure)
			if err != nil {
				logger.Log.Errorf("unable to execute local failure actions: %s", err.Error())
			}
			continue
		}

		// Execute local+global actions
		logger.Log.Infof("create final merged configuration file")

		mergedConfigStr := configEncoder.EncodeConfigString()
		// Update targetConfig with merged config file
		err = os.WriteFile(daemon.config.Config.MergedConfig, []byte(mergedConfigStr), 0744)
		if err != nil {
			logger.Log.Errorf("unable to create '%s' merged config file: %s", daemon.config.Config.MergedConfig, err.Error())
			continue
		}

		logger.Log.Infof("execute success actions")
		err = actions.RunActions(&daemon.config.Config.Actions.OnSuccess)
		if err != nil {
			logger.Log.Errorf("unable to execute local success actions: %s", err.Error())
		}

		err = actions.RunActions(&daemon.config.Global.Actions.OnSuccess)
		if err != nil {
			logger.Log.Errorf("unable to execute global success actions: %s", err.Error())
		}
	}
}
