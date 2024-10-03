package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSession_CanUpdateMetadata(t *testing.T) {
	session := Session{
		AttendanceStatus: AttendanceStatusNotAppliedYet,
	}
	assert.True(t, session.CanUpdateMetadata())

	session.AttendanceStatus = AttendanceStatusIgnored
	assert.True(t, session.CanUpdateMetadata())

	session.AttendanceStatus = AttendanceStatusApplied
	assert.False(t, session.CanUpdateMetadata())
}
