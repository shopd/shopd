package main

import (
	"context"
	"os"

	"github.com/shopd/shopd/www/content/login"
)

func main() {
	ctx := context.Background()
	login.Index().Render(ctx, os.Stdout)
}
