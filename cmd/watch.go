package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/vault"
)

var (
	watchInterval  string
	watchMaxErrors int
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch [path]",
		Short: "Poll a Vault secret path and print changes as they occur",
		Args:  cobra.ExactArgs(1),
		RunE:  runWatch,
	}
	watchCmd.Flags().StringVar(&watchInterval, "interval", "", "polling interval (e.g. 30s, 1m)")
	watchCmd.Flags().IntVar(&watchMaxErrors, "max-errors", 0, "stop after this many consecutive errors (0 = use default)")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	path := args[0]

	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	wcfg := &config.WatchConfig{Interval: watchInterval, MaxErrors: watchMaxErrors}
	opt, err := config.BuildWatchOption(wcfg)
	if err != nil {
		return fmt.Errorf("watch option: %w", err)
	}

	client, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	fetchFn := func(ctx context.Context, p string) (map[string]string, error) {
		return vault.FetchSecrets(ctx, client, []string{p})
	}

	w := vault.NewWatcher(fetchFn, opt)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Fprintf(cmd.OutOrStdout(), "Watching %s every %s (ctrl+c to stop)\n", path, opt.Interval)

	for event := range w.Watch(ctx, path) {
		if event.Err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\n", event.Err)
			continue
		}
		fmt.Fprintf(cmd.OutOrStdout(), "[%s] change detected at %s\n",
			event.CheckAt.Format("15:04:05"), event.Path)
		for k, v := range event.New {
			old := event.Old[k]
			if old != v {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s: %q -> %q\n", k, old, v)
			}
		}
	}
	return nil
}
