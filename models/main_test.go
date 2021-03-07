package models

func init() {
	_ = setupTestDb()
	_ = testDb.AutoMigrate(&User{}, &Application{}, &Key{})
}
