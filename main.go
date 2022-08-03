/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package main

import (
	"context"
	"log"

	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/SomusHQ/tacostand/core"
)

func main() {
	client, err := core.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx = contextutils.WithConfig(ctx, client.Config)
	ctx = contextutils.WithLogger(ctx, client.Logger)
	ctx = contextutils.WithScheduler(ctx, client.Scheduler)
	ctx = contextutils.WithSlackApi(ctx, client.SlackAPI)
	ctx = contextutils.WithSlackSocketClient(ctx, client.SlackSocketClient)
	ctx = contextutils.WithDatabase(ctx, client.Database)

	err = client.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
