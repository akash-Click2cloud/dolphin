package bolt

import "github.com/akash-Click2cloud/dolphin/api"

func (m *Migrator) updateSettingsToVersion3() error {
	legacySettings, err := m.SettingsService.Settings()
	if err != nil {
		return err
	}

	legacySettings.AuthenticationMethod = dockm.AuthenticationInternal
	legacySettings.LDAPSettings = dockm.LDAPSettings{
		TLSConfig: dockm.TLSConfiguration{},
		SearchSettings: []dockm.LDAPSearchSettings{
			dockm.LDAPSearchSettings{},
		},
	}

	err = m.SettingsService.StoreSettings(legacySettings)
	if err != nil {
		return err
	}

	return nil
}