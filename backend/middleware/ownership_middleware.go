package middleware

import (
	"backend/config"
	"backend/models"
)

func EnsureMitraOwnsOutlet(userID uint, outletID uint) (bool, error) {
	var mitra models.Mitra
	if err := config.DB.Where("user_id = ?", userID).First(&mitra).Error; err != nil {
		return false, err
	}

	var outlet models.Outlet
	if err := config.DB.First(&outlet, outletID).Error; err != nil {
		return false, err
	}

	return outlet.MitraID == mitra.ID, nil
}
