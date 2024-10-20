package user

import (
	"fmt"

	"rush/attendance"
)

type updater struct {
	userRepo userRepo
	// Used to update the user data in attendance.
	attendanceRepo attendanceRepo
}

func NewUpdater(userRepo userRepo, attendanceRepo attendanceRepo) *updater {
	return &updater{
		userRepo:       userRepo,
		attendanceRepo: attendanceRepo,
	}
}

func (u *updater) Update(id string, updateForm UpdateForm) error {
	if updateForm.ExternalName != nil || updateForm.Generation != nil {
		updateAttendanceForm := attendance.UpdateUserAttendanceForm{
			UserExternalName: updateForm.ExternalName,
			UserGeneration:   updateForm.Generation,
		}
		if err := u.attendanceRepo.UpdateUserAttendance(id, updateAttendanceForm); err != nil {
			return fmt.Errorf("failed to update user's attendance: %w", err)
		}
	}

	return u.userRepo.Update(id, updateForm)
}

type userRepo interface {
	Update(id string, updateForm UpdateForm) error
}

type attendanceRepo interface {
	UpdateUserAttendance(userId string, updateForm attendance.UpdateUserAttendanceForm) error
}
