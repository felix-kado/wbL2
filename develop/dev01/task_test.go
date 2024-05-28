package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type NTPClient interface {
	Time(server string) (time.Time, error)
}

func TestTimeFetchingSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNTPClient := NewMockNTPClient(ctrl)
	mockNTPClient.EXPECT().Time("pool.ntp.org").Return(time.Now(), nil)

	client = mockNTPClient
	_, err := client.Time("pool.ntp.org")
	assert.Nil(t, err, "Expected no error for NTP time fetching")
}

func TestTimeFetchingFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNTPClient := NewMockNTPClient(ctrl)
	mockNTPClient.EXPECT().Time("bad.server.com").Return(time.Time{}, fmt.Errorf("failed to reach server"))

	client = mockNTPClient
	_, err := client.Time("bad.server.com")
	assert.NotNil(t, err, "Expected an error for NTP time fetching")
}
