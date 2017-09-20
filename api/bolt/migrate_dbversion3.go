package bolt

import "github.com/akash-Click2cloud/dolphin/api"

func (m *Migrator) updateEndpointsToDBVersion4() error {
	legacyEndpoints, err := m.EndpointService.Endpoints()
	if err != nil {
		return err
	}

	for _, endpoint := range legacyEndpoints {
		endpoint.TLSConfig = dockm.TLSConfiguration{}
		if endpoint.TLS {
			endpoint.TLSConfig.TLS = true
			endpoint.TLSConfig.TLSSkipVerify = false
			endpoint.TLSConfig.TLSCACertPath = endpoint.TLSCACertPath
			endpoint.TLSConfig.TLSCertPath = endpoint.TLSCertPath
			endpoint.TLSConfig.TLSKeyPath = endpoint.TLSKeyPath
		}
		err = m.EndpointService.UpdateEndpoint(endpoint.ID, &endpoint)
		if err != nil {
			return err
		}
	}

	return nil
}
