package run

import (
	"log"
	"time"

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

	/////////////////////////////
	// EXECUTION FLOW RELATED
	/////////////////////////////

	//TODO: Set git repository client to get the gcmerge config file
	// client := git.NewGitHubClient()

	for {
		// TODO: Get gcmerge config file from git repository
		// TODO: Compare current gcmerge config file with download one to decide make the changes (make it if first time)
		// TODO: Parse gcmerge config
		// TODO: Get specific config from config name flag in gcmerge config

		// TODO: expand env variables in rawConfig

		// TODO: Decode rawConfig and targetConfig
		// TODO: Merge rawConfig and targetConfig
		// TODO: Check config conditions

		// TODO: Decode global rawConfig
		// TODO: Merge global rawConfig and current merged config
		// TODO: Check global config conditions

		// TODO: Update targetConfig with merged config file

		// TODO: Make config actions
		// TODO: Make global config actions

		// sync in x duration again
		globals.ExecContext.Logger.Infof("Syncing again in %s", duration.String())
		time.Sleep(duration)
	}
}
