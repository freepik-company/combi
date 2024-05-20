package flags

import (
	"log"

	"combi/internal/globals"

	"github.com/spf13/cobra"
)

const (
	LogLevelFlagName     = `log-level`
	DisableTraceFlagName = `disable-trace`
	TmpDirFlagName       = `tmp-dir`

	SyncTimeFlagName    = `sync-time`
	SourceTypeFlagName  = `source-type`
	SourcePathFlagName  = `source-path`
	SourceFieldFlagName = `source-field`

	GitSshUrlFlagName         = `git-ssh-url`
	GitSshKeyFilepathFlagName = `git-sshkey-filepath`
	GitBranchFlagName         = `git-branch`

	// Flags error messages
	logLevelFlagErrMsg     = "unable to get flag --log-level: %s"
	disableTraceFlagErrMsg = "unable to get flag --disable-trace: %s"
	setLoggerErrMsg        = "unable to set logger: %s"
)

type DaemonFlagsT struct {
	LogLevel     string
	DisableTrace bool
	TmpDir       string

	SyncTime    string
	SourceType  string
	SourcePath  string
	SourceField string

	GitSshUrl         string
	GitSshKeyFilepath string
	GitBranch         string
}

func GetSyncRunFlags(cmd *cobra.Command) (srFlags DaemonFlagsT, err error) {

	// Get root command flags
	srFlags.LogLevel, err = cmd.Flags().GetString(LogLevelFlagName)
	if err != nil {
		log.Fatalf(logLevelFlagErrMsg, err.Error())
	}

	srFlags.DisableTrace, err = cmd.Flags().GetBool(DisableTraceFlagName)
	if err != nil {
		log.Fatalf(disableTraceFlagErrMsg, err.Error())
	}

	err = globals.SetLogger(srFlags.LogLevel, srFlags.DisableTrace)
	if err != nil {
		log.Fatalf(setLoggerErrMsg, err.Error())
	}

	// Get synchronizer parent command flags
	srFlags.TmpDir, err = cmd.Flags().GetString(TmpDirFlagName)
	if err != nil {
		return srFlags, err
	}

	srFlags.SyncTime, err = cmd.Flags().GetString(SyncTimeFlagName)
	if err != nil {
		return srFlags, err
	}

	srFlags.SourceType, err = cmd.Flags().GetString(SourceTypeFlagName)
	if err != nil {
		return srFlags, err
	}

	srFlags.SourcePath, err = cmd.Flags().GetString(SourcePathFlagName)
	if err != nil {
		return srFlags, err
	}

	srFlags.SourceField, err = cmd.Flags().GetString(SourceFieldFlagName)
	if err != nil {
		return srFlags, err
	}

	// Get run child command flags
	srFlags.GitSshUrl, err = cmd.Flags().GetString(GitSshUrlFlagName)
	if err != nil {
		return srFlags, err
	}

	srFlags.GitSshKeyFilepath, err = cmd.Flags().GetString(GitSshKeyFilepathFlagName)
	if err != nil {
		return srFlags, err
	}

	srFlags.GitBranch, err = cmd.Flags().GetString(GitBranchFlagName)
	if err != nil {
		return srFlags, err
	}

	return srFlags, err
}
