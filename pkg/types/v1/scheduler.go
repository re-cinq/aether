package v1

type Scheduler interface {

	// Schedule the scraping of the provider data
	Schedule()

	// Cancel the scraping of the provider data
	Cancel()
}
