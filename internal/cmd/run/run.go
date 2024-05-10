package run

import (
	"log"
	"time"

	"gcu/internal/git"
	"gcu/internal/globals"

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
	_ = duration

	// TODO

	/////////////////////////////
	// EXECUTION FLOW RELATED
	/////////////////////////////
	client := git.NewGitHubClient()
}
