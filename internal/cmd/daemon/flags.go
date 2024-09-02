package daemon

import (
	"log"
	"time"

	"combi/internal/globals"
	"combi/internal/logger"

	"github.com/spf13/cobra"
)

const (
	// FLAG NAMES

	logLevelFlagName = `log-level`
	tmpDirFlagName   = `tmp-dir`
	syncTimeFlagName = `sync-time`
	configFlagName   = `config`

	// ERROR MESSAGES

	logLevelFlagErrMsg = "unable to get flag --log-level: %s"
	tmpDirFlagErrMsg   = "unable to get flag --tmp-dir: %s"
	syncTimeFlagErrMsg = "unable to get flag --sync-time: %s"
	configFlagErrMsg   = "unable to get flag --config: %s"
)

type daemonFlagsT struct {
	logLevel string

	tmpDir   string
	syncTime time.Duration
	config   string
}

func (d *daemonT) getFlags(cmd *cobra.Command) (err error) {

	// Get root command flags
	d.flags.logLevel, err = cmd.Flags().GetString(logLevelFlagName)
	if err != nil {
		log.Fatalf(logLevelFlagErrMsg, err.Error())
	}

	level, err := logger.GetLevel(d.flags.logLevel)
	if err != nil {
		log.Fatalf(logLevelFlagErrMsg, err.Error())
	}

	logger.InitLogger(d.context, level)

	// Get command flags

	d.flags.tmpDir, err = cmd.Flags().GetString(tmpDirFlagName)
	if err != nil {
		return err
	}

	globals.TmpDir = d.flags.tmpDir

	d.flags.syncTime, err = cmd.Flags().GetDuration(syncTimeFlagName)
	if err != nil {
		return err
	}

	d.flags.config, err = cmd.Flags().GetString(configFlagName)

	return err
}
