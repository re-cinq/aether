package emissionfactors

// interface actions
// calculate embodied emissions
// calculate operational emissions
// I don't think this needs to be an interface, because the calculations are the
// same no matter the provider.

// I think it would be beneficial to have the emissions-data.generator be
// a callable function. So that we can both read data here, but maybe trigger
// the generation as well if wanted? unless that's going to be a cron job?

func getData() {}
func readEmissionCoefficient() (*emissionFactors, error) {

	return nil, nil
}
