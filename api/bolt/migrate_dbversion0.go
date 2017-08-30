package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/akash-Click2cloud/dolphin/api"
)

func (m *Migrator) updateAdminUserToDBVersion1() error {
	u, err := m.UserService.UserByUsername("admin")
	if err == nil {
		admin := &dockm.User{
			Username: "admin",
			Password: u.Password,
			Role:     dockm.AdministratorRole,
		}
		err = m.UserService.CreateUser(admin)
		if err != nil {
			return err
		}
		err = m.removeLegacyAdminUser()
		if err != nil {
			return err
		}
	} else if err != nil && err != dockm.ErrUserNotFound {
		return err
	}
	return nil
}

func (m *Migrator) removeLegacyAdminUser() error {
	return m.store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucketName))
		err := bucket.Delete([]byte("admin"))
		if err != nil {
			return err
		}
		return nil
	})
}
