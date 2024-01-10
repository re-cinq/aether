package v1

import "context"

type Scheduler interface {

	// Schedule the scraping of the provider data
	Schedule(ctx context.Context)

	// Cancel the scraping of the provider data
	Cancel()
}
