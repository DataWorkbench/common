package getcd

import (
	"context"
	"time"

	"github.com/DataWorkbench/glog"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// RetryElection master worker election, the notifyCallback will be invoking if the election was success.
func RetryElection(ctx context.Context, cli *Client, key string, value string, offer func(ctxCancel context.Context)) {
	// new logger.
	nl := glog.FromContext(ctx).Clone()
	nl.ResetFields().AddString("key", key)

	var sleep bool
LOOP:
	for {
		if sleep {
			// Sleep to prevents died loop.
			time.Sleep(time.Millisecond * 100)
		}
		sleep = true

		nl.Info().Msg("start leader election").Fire()
		sess, err := concurrency.NewSession(cli, concurrency.WithTTL(60))
		if err != nil {
			nl.Error().Error("concurrency new session error", err).Msg("continue").Fire()
			continue LOOP
		}

		election := concurrency.NewElection(sess, key)
		if err = election.Campaign(ctx, value); err != nil {
			if err == context.Canceled {
				nl.Info().Msg("ctx canceled, stop campaign").Fire()
				break LOOP
			}
			nl.Error().Error("election campaign error", err).Msg("continue").Fire()
			continue LOOP
		}

		nl.Info().Msg("current worker is leader and start of term").Fire()

		// start and load crontab.
		ctxCancel, cancel := context.WithCancel(context.Background())

		exitC := make(chan struct{})

		go func() {
			offer(ctxCancel)
			close(exitC)
		}()

		select {
		case <-sess.Done():
			nl.Info().Msg("etcd session done and continue to re-election").Fire()

			cancel()
			// wait for notify callback func exit.
			<-exitC

			continue LOOP
		case <-ctx.Done():
			nl.Info().Msg("receive ctx done signal and end of term").Fire()

			cancel()
			// wait for notify callback func exit.
			<-exitC

			if err = election.Resign(context.Background()); err != nil {
				nl.Error().Error("etcd election resign error", err).Fire()
			}
			break LOOP
		}
	}

	_ = nl.Close()
	return
}
