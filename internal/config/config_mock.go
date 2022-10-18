// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package config

import (
	"sync"
)

// Ensure, that ConfigMock does implement Config.
// If this is not the case, regenerate this file with moq.
var _ Config = &ConfigMock{}

// ConfigMock is a mock implementation of Config.
//
//	func TestSomethingThatUsesConfig(t *testing.T) {
//
//		// make and configure a mocked Config
//		mockedConfig := &ConfigMock{
//			AuthTokenFunc: func() (string, error) {
//				panic("mock out the AuthToken method")
//			},
//			LeaguesFunc: func() []string {
//				panic("mock out the Leagues method")
//			},
//			SetAuthTokenFunc: func(s string) error {
//				panic("mock out the SetAuthToken method")
//			},
//			SetLeaguesFunc: func(strings []string) error {
//				panic("mock out the SetLeagues method")
//			},
//		}
//
//		// use mockedConfig in code that requires Config
//		// and then make assertions.
//
//	}
type ConfigMock struct {
	// AuthTokenFunc mocks the AuthToken method.
	AuthTokenFunc func() string

	// LeaguesFunc mocks the Leagues method.
	LeaguesFunc func() []string

	// SetAuthTokenFunc mocks the SetAuthToken method.
	SetAuthTokenFunc func(s string) error

	// SetLeaguesFunc mocks the SetLeagues method.
	SetLeaguesFunc func(strings []string) error

	// calls tracks calls to the methods.
	calls struct {
		// AuthToken holds details about calls to the AuthToken method.
		AuthToken []struct {
		}
		// Leagues holds details about calls to the Leagues method.
		Leagues []struct {
		}
		// SetAuthToken holds details about calls to the SetAuthToken method.
		SetAuthToken []struct {
			// S is the s argument value.
			S string
		}
		// SetLeagues holds details about calls to the SetLeagues method.
		SetLeagues []struct {
			// Strings is the strings argument value.
			Strings []string
		}
	}
	lockAuthToken    sync.RWMutex
	lockLeagues      sync.RWMutex
	lockSetAuthToken sync.RWMutex
	lockSetLeagues   sync.RWMutex
}

// AuthToken calls AuthTokenFunc.
func (mock *ConfigMock) AuthToken() string {
	if mock.AuthTokenFunc == nil {
		panic("ConfigMock.AuthTokenFunc: method is nil but Config.AuthToken was just called")
	}
	callInfo := struct {
	}{}
	mock.lockAuthToken.Lock()
	mock.calls.AuthToken = append(mock.calls.AuthToken, callInfo)
	mock.lockAuthToken.Unlock()
	return mock.AuthTokenFunc()
}

// AuthTokenCalls gets all the calls that were made to AuthToken.
// Check the length with:
//
//	len(mockedConfig.AuthTokenCalls())
func (mock *ConfigMock) AuthTokenCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockAuthToken.RLock()
	calls = mock.calls.AuthToken
	mock.lockAuthToken.RUnlock()
	return calls
}

// Leagues calls LeaguesFunc.
func (mock *ConfigMock) Leagues() []string {
	if mock.LeaguesFunc == nil {
		panic("ConfigMock.LeaguesFunc: method is nil but Config.Leagues was just called")
	}
	callInfo := struct {
	}{}
	mock.lockLeagues.Lock()
	mock.calls.Leagues = append(mock.calls.Leagues, callInfo)
	mock.lockLeagues.Unlock()
	return mock.LeaguesFunc()
}

// LeaguesCalls gets all the calls that were made to Leagues.
// Check the length with:
//
//	len(mockedConfig.LeaguesCalls())
func (mock *ConfigMock) LeaguesCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockLeagues.RLock()
	calls = mock.calls.Leagues
	mock.lockLeagues.RUnlock()
	return calls
}

// SetAuthToken calls SetAuthTokenFunc.
func (mock *ConfigMock) SetAuthToken(s string) error {
	if mock.SetAuthTokenFunc == nil {
		panic("ConfigMock.SetAuthTokenFunc: method is nil but Config.SetAuthToken was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockSetAuthToken.Lock()
	mock.calls.SetAuthToken = append(mock.calls.SetAuthToken, callInfo)
	mock.lockSetAuthToken.Unlock()
	return mock.SetAuthTokenFunc(s)
}

// SetAuthTokenCalls gets all the calls that were made to SetAuthToken.
// Check the length with:
//
//	len(mockedConfig.SetAuthTokenCalls())
func (mock *ConfigMock) SetAuthTokenCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockSetAuthToken.RLock()
	calls = mock.calls.SetAuthToken
	mock.lockSetAuthToken.RUnlock()
	return calls
}

// SetLeagues calls SetLeaguesFunc.
func (mock *ConfigMock) SetLeagues(strings []string) error {
	if mock.SetLeaguesFunc == nil {
		panic("ConfigMock.SetLeaguesFunc: method is nil but Config.SetLeagues was just called")
	}
	callInfo := struct {
		Strings []string
	}{
		Strings: strings,
	}
	mock.lockSetLeagues.Lock()
	mock.calls.SetLeagues = append(mock.calls.SetLeagues, callInfo)
	mock.lockSetLeagues.Unlock()
	return mock.SetLeaguesFunc(strings)
}

// SetLeaguesCalls gets all the calls that were made to SetLeagues.
// Check the length with:
//
//	len(mockedConfig.SetLeaguesCalls())
func (mock *ConfigMock) SetLeaguesCalls() []struct {
	Strings []string
} {
	var calls []struct {
		Strings []string
	}
	mock.lockSetLeagues.RLock()
	calls = mock.calls.SetLeagues
	mock.lockSetLeagues.RUnlock()
	return calls
}
