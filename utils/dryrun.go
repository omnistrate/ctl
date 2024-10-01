package utils

const (
	DryRunEnv = "DRY_RUN"
)

func IsDryRun() bool {
	return GetEnvAsBoolean(DryRunEnv, "false")
}
