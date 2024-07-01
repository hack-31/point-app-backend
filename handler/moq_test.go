// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package handler

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository/entities"
	"github.com/hack-31/point-app-backend/service"
)

// Ensure, that RegisterUserServiceMock does implement RegisterUserService.
// If this is not the case, regenerate this file with moq.
var _ RegisterUserService = &RegisterUserServiceMock{}

// RegisterUserServiceMock is a mock implementation of RegisterUserService.
//
//	func TestSomethingThatUsesRegisterUserService(t *testing.T) {
//
//		// make and configure a mocked RegisterUserService
//		mockedRegisterUserService := &RegisterUserServiceMock{
//			RegisterUserFunc: func(ctx context.Context, temporaryUserId string, confirmCode string) (*entities.User, string, error) {
//				panic("mock out the RegisterUser method")
//			},
//		}
//
//		// use mockedRegisterUserService in code that requires RegisterUserService
//		// and then make assertions.
//
//	}
type RegisterUserServiceMock struct {
	// RegisterUserFunc mocks the RegisterUser method.
	RegisterUserFunc func(ctx context.Context, temporaryUserId string, confirmCode string) (*entities.User, string, error)

	// calls tracks calls to the methods.
	calls struct {
		// RegisterUser holds details about calls to the RegisterUser method.
		RegisterUser []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// TemporaryUserId is the temporaryUserId argument value.
			TemporaryUserId string
			// ConfirmCode is the confirmCode argument value.
			ConfirmCode string
		}
	}
	lockRegisterUser sync.RWMutex
}

// RegisterUser calls RegisterUserFunc.
func (mock *RegisterUserServiceMock) RegisterUser(ctx context.Context, temporaryUserId string, confirmCode string) (*entities.User, string, error) {
	if mock.RegisterUserFunc == nil {
		panic("RegisterUserServiceMock.RegisterUserFunc: method is nil but RegisterUserService.RegisterUser was just called")
	}
	callInfo := struct {
		Ctx             context.Context
		TemporaryUserId string
		ConfirmCode     string
	}{
		Ctx:             ctx,
		TemporaryUserId: temporaryUserId,
		ConfirmCode:     confirmCode,
	}
	mock.lockRegisterUser.Lock()
	mock.calls.RegisterUser = append(mock.calls.RegisterUser, callInfo)
	mock.lockRegisterUser.Unlock()
	return mock.RegisterUserFunc(ctx, temporaryUserId, confirmCode)
}

// RegisterUserCalls gets all the calls that were made to RegisterUser.
// Check the length with:
//
//	len(mockedRegisterUserService.RegisterUserCalls())
func (mock *RegisterUserServiceMock) RegisterUserCalls() []struct {
	Ctx             context.Context
	TemporaryUserId string
	ConfirmCode     string
} {
	var calls []struct {
		Ctx             context.Context
		TemporaryUserId string
		ConfirmCode     string
	}
	mock.lockRegisterUser.RLock()
	calls = mock.calls.RegisterUser
	mock.lockRegisterUser.RUnlock()
	return calls
}

// Ensure, that RegisterTemporaryUserServiceMock does implement RegisterTemporaryUserService.
// If this is not the case, regenerate this file with moq.
var _ RegisterTemporaryUserService = &RegisterTemporaryUserServiceMock{}

// RegisterTemporaryUserServiceMock is a mock implementation of RegisterTemporaryUserService.
//
//	func TestSomethingThatUsesRegisterTemporaryUserService(t *testing.T) {
//
//		// make and configure a mocked RegisterTemporaryUserService
//		mockedRegisterTemporaryUserService := &RegisterTemporaryUserServiceMock{
//			RegisterTemporaryUserFunc: func(ctx context.Context, firstName string, firstNameKana string, familyName string, familyNameKana string, email string, password string) (string, error) {
//				panic("mock out the RegisterTemporaryUser method")
//			},
//		}
//
//		// use mockedRegisterTemporaryUserService in code that requires RegisterTemporaryUserService
//		// and then make assertions.
//
//	}
type RegisterTemporaryUserServiceMock struct {
	// RegisterTemporaryUserFunc mocks the RegisterTemporaryUser method.
	RegisterTemporaryUserFunc func(ctx context.Context, firstName string, firstNameKana string, familyName string, familyNameKana string, email string, password string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// RegisterTemporaryUser holds details about calls to the RegisterTemporaryUser method.
		RegisterTemporaryUser []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// FirstName is the firstName argument value.
			FirstName string
			// FirstNameKana is the firstNameKana argument value.
			FirstNameKana string
			// FamilyName is the familyName argument value.
			FamilyName string
			// FamilyNameKana is the familyNameKana argument value.
			FamilyNameKana string
			// Email is the email argument value.
			Email string
			// Password is the password argument value.
			Password string
		}
	}
	lockRegisterTemporaryUser sync.RWMutex
}

// RegisterTemporaryUser calls RegisterTemporaryUserFunc.
func (mock *RegisterTemporaryUserServiceMock) RegisterTemporaryUser(ctx context.Context, firstName string, firstNameKana string, familyName string, familyNameKana string, email string, password string) (string, error) {
	if mock.RegisterTemporaryUserFunc == nil {
		panic("RegisterTemporaryUserServiceMock.RegisterTemporaryUserFunc: method is nil but RegisterTemporaryUserService.RegisterTemporaryUser was just called")
	}
	callInfo := struct {
		Ctx            context.Context
		FirstName      string
		FirstNameKana  string
		FamilyName     string
		FamilyNameKana string
		Email          string
		Password       string
	}{
		Ctx:            ctx,
		FirstName:      firstName,
		FirstNameKana:  firstNameKana,
		FamilyName:     familyName,
		FamilyNameKana: familyNameKana,
		Email:          email,
		Password:       password,
	}
	mock.lockRegisterTemporaryUser.Lock()
	mock.calls.RegisterTemporaryUser = append(mock.calls.RegisterTemporaryUser, callInfo)
	mock.lockRegisterTemporaryUser.Unlock()
	return mock.RegisterTemporaryUserFunc(ctx, firstName, firstNameKana, familyName, familyNameKana, email, password)
}

// RegisterTemporaryUserCalls gets all the calls that were made to RegisterTemporaryUser.
// Check the length with:
//
//	len(mockedRegisterTemporaryUserService.RegisterTemporaryUserCalls())
func (mock *RegisterTemporaryUserServiceMock) RegisterTemporaryUserCalls() []struct {
	Ctx            context.Context
	FirstName      string
	FirstNameKana  string
	FamilyName     string
	FamilyNameKana string
	Email          string
	Password       string
} {
	var calls []struct {
		Ctx            context.Context
		FirstName      string
		FirstNameKana  string
		FamilyName     string
		FamilyNameKana string
		Email          string
		Password       string
	}
	mock.lockRegisterTemporaryUser.RLock()
	calls = mock.calls.RegisterTemporaryUser
	mock.lockRegisterTemporaryUser.RUnlock()
	return calls
}

// Ensure, that SigninServiceMock does implement SigninService.
// If this is not the case, regenerate this file with moq.
var _ SigninService = &SigninServiceMock{}

// SigninServiceMock is a mock implementation of SigninService.
//
//	func TestSomethingThatUsesSigninService(t *testing.T) {
//
//		// make and configure a mocked SigninService
//		mockedSigninService := &SigninServiceMock{
//			SigninFunc: func(ctx context.Context, email string, password string) (string, error) {
//				panic("mock out the Signin method")
//			},
//		}
//
//		// use mockedSigninService in code that requires SigninService
//		// and then make assertions.
//
//	}
type SigninServiceMock struct {
	// SigninFunc mocks the Signin method.
	SigninFunc func(ctx context.Context, email string, password string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// Signin holds details about calls to the Signin method.
		Signin []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Email is the email argument value.
			Email string
			// Password is the password argument value.
			Password string
		}
	}
	lockSignin sync.RWMutex
}

// Signin calls SigninFunc.
func (mock *SigninServiceMock) Signin(ctx context.Context, email string, password string) (string, error) {
	if mock.SigninFunc == nil {
		panic("SigninServiceMock.SigninFunc: method is nil but SigninService.Signin was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Email    string
		Password string
	}{
		Ctx:      ctx,
		Email:    email,
		Password: password,
	}
	mock.lockSignin.Lock()
	mock.calls.Signin = append(mock.calls.Signin, callInfo)
	mock.lockSignin.Unlock()
	return mock.SigninFunc(ctx, email, password)
}

// SigninCalls gets all the calls that were made to Signin.
// Check the length with:
//
//	len(mockedSigninService.SigninCalls())
func (mock *SigninServiceMock) SigninCalls() []struct {
	Ctx      context.Context
	Email    string
	Password string
} {
	var calls []struct {
		Ctx      context.Context
		Email    string
		Password string
	}
	mock.lockSignin.RLock()
	calls = mock.calls.Signin
	mock.lockSignin.RUnlock()
	return calls
}

// Ensure, that GetUsersServiceMock does implement GetUsersService.
// If this is not the case, regenerate this file with moq.
var _ GetUsersService = &GetUsersServiceMock{}

// GetUsersServiceMock is a mock implementation of GetUsersService.
//
//	func TestSomethingThatUsesGetUsersService(t *testing.T) {
//
//		// make and configure a mocked GetUsersService
//		mockedGetUsersService := &GetUsersServiceMock{
//			GetUsersFunc: func(ctx context.Context) (service.GetUsersResponse, error) {
//				panic("mock out the GetUsers method")
//			},
//		}
//
//		// use mockedGetUsersService in code that requires GetUsersService
//		// and then make assertions.
//
//	}
type GetUsersServiceMock struct {
	// GetUsersFunc mocks the GetUsers method.
	GetUsersFunc func(ctx context.Context) (service.GetUsersResponse, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetUsers holds details about calls to the GetUsers method.
		GetUsers []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
	}
	lockGetUsers sync.RWMutex
}

// GetUsers calls GetUsersFunc.
func (mock *GetUsersServiceMock) GetUsers(ctx context.Context) (service.GetUsersResponse, error) {
	if mock.GetUsersFunc == nil {
		panic("GetUsersServiceMock.GetUsersFunc: method is nil but GetUsersService.GetUsers was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockGetUsers.Lock()
	mock.calls.GetUsers = append(mock.calls.GetUsers, callInfo)
	mock.lockGetUsers.Unlock()
	return mock.GetUsersFunc(ctx)
}

// GetUsersCalls gets all the calls that were made to GetUsers.
// Check the length with:
//
//	len(mockedGetUsersService.GetUsersCalls())
func (mock *GetUsersServiceMock) GetUsersCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockGetUsers.RLock()
	calls = mock.calls.GetUsers
	mock.lockGetUsers.RUnlock()
	return calls
}

// Ensure, that UpdatePasswordServiceMock does implement UpdatePasswordService.
// If this is not the case, regenerate this file with moq.
var _ UpdatePasswordService = &UpdatePasswordServiceMock{}

// UpdatePasswordServiceMock is a mock implementation of UpdatePasswordService.
//
//	func TestSomethingThatUsesUpdatePasswordService(t *testing.T) {
//
//		// make and configure a mocked UpdatePasswordService
//		mockedUpdatePasswordService := &UpdatePasswordServiceMock{
//			UpdatePasswordFunc: func(ctx *gin.Context, oldPassword string, newPassword string) error {
//				panic("mock out the UpdatePassword method")
//			},
//		}
//
//		// use mockedUpdatePasswordService in code that requires UpdatePasswordService
//		// and then make assertions.
//
//	}
type UpdatePasswordServiceMock struct {
	// UpdatePasswordFunc mocks the UpdatePassword method.
	UpdatePasswordFunc func(ctx *gin.Context, oldPassword string, newPassword string) error

	// calls tracks calls to the methods.
	calls struct {
		// UpdatePassword holds details about calls to the UpdatePassword method.
		UpdatePassword []struct {
			// Ctx is the ctx argument value.
			Ctx *gin.Context
			// OldPassword is the oldPassword argument value.
			OldPassword string
			// NewPassword is the newPassword argument value.
			NewPassword string
		}
	}
	lockUpdatePassword sync.RWMutex
}

// UpdatePassword calls UpdatePasswordFunc.
func (mock *UpdatePasswordServiceMock) UpdatePassword(ctx *gin.Context, oldPassword string, newPassword string) error {
	if mock.UpdatePasswordFunc == nil {
		panic("UpdatePasswordServiceMock.UpdatePasswordFunc: method is nil but UpdatePasswordService.UpdatePassword was just called")
	}
	callInfo := struct {
		Ctx         *gin.Context
		OldPassword string
		NewPassword string
	}{
		Ctx:         ctx,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}
	mock.lockUpdatePassword.Lock()
	mock.calls.UpdatePassword = append(mock.calls.UpdatePassword, callInfo)
	mock.lockUpdatePassword.Unlock()
	return mock.UpdatePasswordFunc(ctx, oldPassword, newPassword)
}

// UpdatePasswordCalls gets all the calls that were made to UpdatePassword.
// Check the length with:
//
//	len(mockedUpdatePasswordService.UpdatePasswordCalls())
func (mock *UpdatePasswordServiceMock) UpdatePasswordCalls() []struct {
	Ctx         *gin.Context
	OldPassword string
	NewPassword string
} {
	var calls []struct {
		Ctx         *gin.Context
		OldPassword string
		NewPassword string
	}
	mock.lockUpdatePassword.RLock()
	calls = mock.calls.UpdatePassword
	mock.lockUpdatePassword.RUnlock()
	return calls
}

// Ensure, that UpdateAccountServiceMock does implement UpdateAccountService.
// If this is not the case, regenerate this file with moq.
var _ UpdateAccountService = &UpdateAccountServiceMock{}

// UpdateAccountServiceMock is a mock implementation of UpdateAccountService.
//
//	func TestSomethingThatUsesUpdateAccountService(t *testing.T) {
//
//		// make and configure a mocked UpdateAccountService
//		mockedUpdateAccountService := &UpdateAccountServiceMock{
//			UpdateAccountFunc: func(ctx *gin.Context, familyName string, familyNameKana string, firstName string, firstNameKana string) error {
//				panic("mock out the UpdateAccount method")
//			},
//		}
//
//		// use mockedUpdateAccountService in code that requires UpdateAccountService
//		// and then make assertions.
//
//	}
type UpdateAccountServiceMock struct {
	// UpdateAccountFunc mocks the UpdateAccount method.
	UpdateAccountFunc func(ctx *gin.Context, familyName string, familyNameKana string, firstName string, firstNameKana string) error

	// calls tracks calls to the methods.
	calls struct {
		// UpdateAccount holds details about calls to the UpdateAccount method.
		UpdateAccount []struct {
			// Ctx is the ctx argument value.
			Ctx *gin.Context
			// FamilyName is the familyName argument value.
			FamilyName string
			// FamilyNameKana is the familyNameKana argument value.
			FamilyNameKana string
			// FirstName is the firstName argument value.
			FirstName string
			// FirstNameKana is the firstNameKana argument value.
			FirstNameKana string
		}
	}
	lockUpdateAccount sync.RWMutex
}

// UpdateAccount calls UpdateAccountFunc.
func (mock *UpdateAccountServiceMock) UpdateAccount(ctx *gin.Context, familyName string, familyNameKana string, firstName string, firstNameKana string) error {
	if mock.UpdateAccountFunc == nil {
		panic("UpdateAccountServiceMock.UpdateAccountFunc: method is nil but UpdateAccountService.UpdateAccount was just called")
	}
	callInfo := struct {
		Ctx            *gin.Context
		FamilyName     string
		FamilyNameKana string
		FirstName      string
		FirstNameKana  string
	}{
		Ctx:            ctx,
		FamilyName:     familyName,
		FamilyNameKana: familyNameKana,
		FirstName:      firstName,
		FirstNameKana:  firstNameKana,
	}
	mock.lockUpdateAccount.Lock()
	mock.calls.UpdateAccount = append(mock.calls.UpdateAccount, callInfo)
	mock.lockUpdateAccount.Unlock()
	return mock.UpdateAccountFunc(ctx, familyName, familyNameKana, firstName, firstNameKana)
}

// UpdateAccountCalls gets all the calls that were made to UpdateAccount.
// Check the length with:
//
//	len(mockedUpdateAccountService.UpdateAccountCalls())
func (mock *UpdateAccountServiceMock) UpdateAccountCalls() []struct {
	Ctx            *gin.Context
	FamilyName     string
	FamilyNameKana string
	FirstName      string
	FirstNameKana  string
} {
	var calls []struct {
		Ctx            *gin.Context
		FamilyName     string
		FamilyNameKana string
		FirstName      string
		FirstNameKana  string
	}
	mock.lockUpdateAccount.RLock()
	calls = mock.calls.UpdateAccount
	mock.lockUpdateAccount.RUnlock()
	return calls
}

// Ensure, that ResetPasswordServiceMock does implement ResetPasswordService.
// If this is not the case, regenerate this file with moq.
var _ ResetPasswordService = &ResetPasswordServiceMock{}

// ResetPasswordServiceMock is a mock implementation of ResetPasswordService.
//
//	func TestSomethingThatUsesResetPasswordService(t *testing.T) {
//
//		// make and configure a mocked ResetPasswordService
//		mockedResetPasswordService := &ResetPasswordServiceMock{
//			ResetPasswordFunc: func(ctx context.Context, email string) error {
//				panic("mock out the ResetPassword method")
//			},
//		}
//
//		// use mockedResetPasswordService in code that requires ResetPasswordService
//		// and then make assertions.
//
//	}
type ResetPasswordServiceMock struct {
	// ResetPasswordFunc mocks the ResetPassword method.
	ResetPasswordFunc func(ctx context.Context, email string) error

	// calls tracks calls to the methods.
	calls struct {
		// ResetPassword holds details about calls to the ResetPassword method.
		ResetPassword []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Email is the email argument value.
			Email string
		}
	}
	lockResetPassword sync.RWMutex
}

// ResetPassword calls ResetPasswordFunc.
func (mock *ResetPasswordServiceMock) ResetPassword(ctx context.Context, email string) error {
	if mock.ResetPasswordFunc == nil {
		panic("ResetPasswordServiceMock.ResetPasswordFunc: method is nil but ResetPasswordService.ResetPassword was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Email string
	}{
		Ctx:   ctx,
		Email: email,
	}
	mock.lockResetPassword.Lock()
	mock.calls.ResetPassword = append(mock.calls.ResetPassword, callInfo)
	mock.lockResetPassword.Unlock()
	return mock.ResetPasswordFunc(ctx, email)
}

// ResetPasswordCalls gets all the calls that were made to ResetPassword.
// Check the length with:
//
//	len(mockedResetPasswordService.ResetPasswordCalls())
func (mock *ResetPasswordServiceMock) ResetPasswordCalls() []struct {
	Ctx   context.Context
	Email string
} {
	var calls []struct {
		Ctx   context.Context
		Email string
	}
	mock.lockResetPassword.RLock()
	calls = mock.calls.ResetPassword
	mock.lockResetPassword.RUnlock()
	return calls
}

// Ensure, that SendPointServiceMock does implement SendPointService.
// If this is not the case, regenerate this file with moq.
var _ SendPointService = &SendPointServiceMock{}

// SendPointServiceMock is a mock implementation of SendPointService.
//
//	func TestSomethingThatUsesSendPointService(t *testing.T) {
//
//		// make and configure a mocked SendPointService
//		mockedSendPointService := &SendPointServiceMock{
//			SendPointFunc: func(ctx *gin.Context, toUserId int, sendPoint int) error {
//				panic("mock out the SendPoint method")
//			},
//		}
//
//		// use mockedSendPointService in code that requires SendPointService
//		// and then make assertions.
//
//	}
type SendPointServiceMock struct {
	// SendPointFunc mocks the SendPoint method.
	SendPointFunc func(ctx *gin.Context, toUserId int, sendPoint int) error

	// calls tracks calls to the methods.
	calls struct {
		// SendPoint holds details about calls to the SendPoint method.
		SendPoint []struct {
			// Ctx is the ctx argument value.
			Ctx *gin.Context
			// ToUserId is the toUserId argument value.
			ToUserId int
			// SendPoint is the sendPoint argument value.
			SendPoint int
		}
	}
	lockSendPoint sync.RWMutex
}

// SendPoint calls SendPointFunc.
func (mock *SendPointServiceMock) SendPoint(ctx *gin.Context, toUserId int, sendPoint int) error {
	if mock.SendPointFunc == nil {
		panic("SendPointServiceMock.SendPointFunc: method is nil but SendPointService.SendPoint was just called")
	}
	callInfo := struct {
		Ctx       *gin.Context
		ToUserId  int
		SendPoint int
	}{
		Ctx:       ctx,
		ToUserId:  toUserId,
		SendPoint: sendPoint,
	}
	mock.lockSendPoint.Lock()
	mock.calls.SendPoint = append(mock.calls.SendPoint, callInfo)
	mock.lockSendPoint.Unlock()
	return mock.SendPointFunc(ctx, toUserId, sendPoint)
}

// SendPointCalls gets all the calls that were made to SendPoint.
// Check the length with:
//
//	len(mockedSendPointService.SendPointCalls())
func (mock *SendPointServiceMock) SendPointCalls() []struct {
	Ctx       *gin.Context
	ToUserId  int
	SendPoint int
} {
	var calls []struct {
		Ctx       *gin.Context
		ToUserId  int
		SendPoint int
	}
	mock.lockSendPoint.RLock()
	calls = mock.calls.SendPoint
	mock.lockSendPoint.RUnlock()
	return calls
}

// Ensure, that SignoutServiceMock does implement SignoutService.
// If this is not the case, regenerate this file with moq.
var _ SignoutService = &SignoutServiceMock{}

// SignoutServiceMock is a mock implementation of SignoutService.
//
//	func TestSomethingThatUsesSignoutService(t *testing.T) {
//
//		// make and configure a mocked SignoutService
//		mockedSignoutService := &SignoutServiceMock{
//			SignoutFunc: func(ctx *gin.Context) error {
//				panic("mock out the Signout method")
//			},
//		}
//
//		// use mockedSignoutService in code that requires SignoutService
//		// and then make assertions.
//
//	}
type SignoutServiceMock struct {
	// SignoutFunc mocks the Signout method.
	SignoutFunc func(ctx *gin.Context) error

	// calls tracks calls to the methods.
	calls struct {
		// Signout holds details about calls to the Signout method.
		Signout []struct {
			// Ctx is the ctx argument value.
			Ctx *gin.Context
		}
	}
	lockSignout sync.RWMutex
}

// Signout calls SignoutFunc.
func (mock *SignoutServiceMock) Signout(ctx *gin.Context) error {
	if mock.SignoutFunc == nil {
		panic("SignoutServiceMock.SignoutFunc: method is nil but SignoutService.Signout was just called")
	}
	callInfo := struct {
		Ctx *gin.Context
	}{
		Ctx: ctx,
	}
	mock.lockSignout.Lock()
	mock.calls.Signout = append(mock.calls.Signout, callInfo)
	mock.lockSignout.Unlock()
	return mock.SignoutFunc(ctx)
}

// SignoutCalls gets all the calls that were made to Signout.
// Check the length with:
//
//	len(mockedSignoutService.SignoutCalls())
func (mock *SignoutServiceMock) SignoutCalls() []struct {
	Ctx *gin.Context
} {
	var calls []struct {
		Ctx *gin.Context
	}
	mock.lockSignout.RLock()
	calls = mock.calls.Signout
	mock.lockSignout.RUnlock()
	return calls
}

// Ensure, that GetAccountServiceMock does implement GetAccountService.
// If this is not the case, regenerate this file with moq.
var _ GetAccountService = &GetAccountServiceMock{}

// GetAccountServiceMock is a mock implementation of GetAccountService.
//
//	func TestSomethingThatUsesGetAccountService(t *testing.T) {
//
//		// make and configure a mocked GetAccountService
//		mockedGetAccountService := &GetAccountServiceMock{
//			GetAccountFunc: func(ctx *gin.Context) (service.GetAccountResponse, error) {
//				panic("mock out the GetAccount method")
//			},
//		}
//
//		// use mockedGetAccountService in code that requires GetAccountService
//		// and then make assertions.
//
//	}
type GetAccountServiceMock struct {
	// GetAccountFunc mocks the GetAccount method.
	GetAccountFunc func(ctx *gin.Context) (service.GetAccountResponse, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetAccount holds details about calls to the GetAccount method.
		GetAccount []struct {
			// Ctx is the ctx argument value.
			Ctx *gin.Context
		}
	}
	lockGetAccount sync.RWMutex
}

// GetAccount calls GetAccountFunc.
func (mock *GetAccountServiceMock) GetAccount(ctx *gin.Context) (service.GetAccountResponse, error) {
	if mock.GetAccountFunc == nil {
		panic("GetAccountServiceMock.GetAccountFunc: method is nil but GetAccountService.GetAccount was just called")
	}
	callInfo := struct {
		Ctx *gin.Context
	}{
		Ctx: ctx,
	}
	mock.lockGetAccount.Lock()
	mock.calls.GetAccount = append(mock.calls.GetAccount, callInfo)
	mock.lockGetAccount.Unlock()
	return mock.GetAccountFunc(ctx)
}

// GetAccountCalls gets all the calls that were made to GetAccount.
// Check the length with:
//
//	len(mockedGetAccountService.GetAccountCalls())
func (mock *GetAccountServiceMock) GetAccountCalls() []struct {
	Ctx *gin.Context
} {
	var calls []struct {
		Ctx *gin.Context
	}
	mock.lockGetAccount.RLock()
	calls = mock.calls.GetAccount
	mock.lockGetAccount.RUnlock()
	return calls
}

// Ensure, that UpdateTemporaryEmailServiceMock does implement UpdateTemporaryEmailService.
// If this is not the case, regenerate this file with moq.
var _ UpdateTemporaryEmailService = &UpdateTemporaryEmailServiceMock{}

// UpdateTemporaryEmailServiceMock is a mock implementation of UpdateTemporaryEmailService.
//
//	func TestSomethingThatUsesUpdateTemporaryEmailService(t *testing.T) {
//
//		// make and configure a mocked UpdateTemporaryEmailService
//		mockedUpdateTemporaryEmailService := &UpdateTemporaryEmailServiceMock{
//			UpdateTemporaryEmailFunc: func(ctx *gin.Context, email string) (string, error) {
//				panic("mock out the UpdateTemporaryEmail method")
//			},
//		}
//
//		// use mockedUpdateTemporaryEmailService in code that requires UpdateTemporaryEmailService
//		// and then make assertions.
//
//	}
type UpdateTemporaryEmailServiceMock struct {
	// UpdateTemporaryEmailFunc mocks the UpdateTemporaryEmail method.
	UpdateTemporaryEmailFunc func(ctx *gin.Context, email string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// UpdateTemporaryEmail holds details about calls to the UpdateTemporaryEmail method.
		UpdateTemporaryEmail []struct {
			// Ctx is the ctx argument value.
			Ctx *gin.Context
			// Email is the email argument value.
			Email string
		}
	}
	lockUpdateTemporaryEmail sync.RWMutex
}

// UpdateTemporaryEmail calls UpdateTemporaryEmailFunc.
func (mock *UpdateTemporaryEmailServiceMock) UpdateTemporaryEmail(ctx *gin.Context, email string) (string, error) {
	if mock.UpdateTemporaryEmailFunc == nil {
		panic("UpdateTemporaryEmailServiceMock.UpdateTemporaryEmailFunc: method is nil but UpdateTemporaryEmailService.UpdateTemporaryEmail was just called")
	}
	callInfo := struct {
		Ctx   *gin.Context
		Email string
	}{
		Ctx:   ctx,
		Email: email,
	}
	mock.lockUpdateTemporaryEmail.Lock()
	mock.calls.UpdateTemporaryEmail = append(mock.calls.UpdateTemporaryEmail, callInfo)
	mock.lockUpdateTemporaryEmail.Unlock()
	return mock.UpdateTemporaryEmailFunc(ctx, email)
}

// UpdateTemporaryEmailCalls gets all the calls that were made to UpdateTemporaryEmail.
// Check the length with:
//
//	len(mockedUpdateTemporaryEmailService.UpdateTemporaryEmailCalls())
func (mock *UpdateTemporaryEmailServiceMock) UpdateTemporaryEmailCalls() []struct {
	Ctx   *gin.Context
	Email string
} {
	var calls []struct {
		Ctx   *gin.Context
		Email string
	}
	mock.lockUpdateTemporaryEmail.RLock()
	calls = mock.calls.UpdateTemporaryEmail
	mock.lockUpdateTemporaryEmail.RUnlock()
	return calls
}

// Ensure, that GetNotificationServiceMock does implement GetNotificationService.
// If this is not the case, regenerate this file with moq.
var _ GetNotificationService = &GetNotificationServiceMock{}

// GetNotificationServiceMock is a mock implementation of GetNotificationService.
//
//	func TestSomethingThatUsesGetNotificationService(t *testing.T) {
//
//		// make and configure a mocked GetNotificationService
//		mockedGetNotificationService := &GetNotificationServiceMock{
//			GetNotificationFunc: func(ctx *gin.Context, notificationID model.NotificationID) (service.GetNotificationResponse, error) {
//				panic("mock out the GetNotification method")
//			},
//		}
//
//		// use mockedGetNotificationService in code that requires GetNotificationService
//		// and then make assertions.
//
//	}
type GetNotificationServiceMock struct {
	// GetNotificationFunc mocks the GetNotification method.
	GetNotificationFunc func(ctx *gin.Context, notificationID model.NotificationID) (service.GetNotificationResponse, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetNotification holds details about calls to the GetNotification method.
		GetNotification []struct {
			// Ctx is the ctx argument value.
			Ctx *gin.Context
			// NotificationID is the notificationID argument value.
			NotificationID model.NotificationID
		}
	}
	lockGetNotification sync.RWMutex
}

// GetNotification calls GetNotificationFunc.
func (mock *GetNotificationServiceMock) GetNotification(ctx *gin.Context, notificationID model.NotificationID) (service.GetNotificationResponse, error) {
	if mock.GetNotificationFunc == nil {
		panic("GetNotificationServiceMock.GetNotificationFunc: method is nil but GetNotificationService.GetNotification was just called")
	}
	callInfo := struct {
		Ctx            *gin.Context
		NotificationID model.NotificationID
	}{
		Ctx:            ctx,
		NotificationID: notificationID,
	}
	mock.lockGetNotification.Lock()
	mock.calls.GetNotification = append(mock.calls.GetNotification, callInfo)
	mock.lockGetNotification.Unlock()
	return mock.GetNotificationFunc(ctx, notificationID)
}

// GetNotificationCalls gets all the calls that were made to GetNotification.
// Check the length with:
//
//	len(mockedGetNotificationService.GetNotificationCalls())
func (mock *GetNotificationServiceMock) GetNotificationCalls() []struct {
	Ctx            *gin.Context
	NotificationID model.NotificationID
} {
	var calls []struct {
		Ctx            *gin.Context
		NotificationID model.NotificationID
	}
	mock.lockGetNotification.RLock()
	calls = mock.calls.GetNotification
	mock.lockGetNotification.RUnlock()
	return calls
}

// Ensure, that GetNotificationsServiceMock does implement GetNotificationsService.
// If this is not the case, regenerate this file with moq.
var _ GetNotificationsService = &GetNotificationsServiceMock{}

// GetNotificationsServiceMock is a mock implementation of GetNotificationsService.
//
//	func TestSomethingThatUsesGetNotificationsService(t *testing.T) {
//
//		// make and configure a mocked GetNotificationsService
//		mockedGetNotificationsService := &GetNotificationsServiceMock{
//			GetNotificationsFunc: func(ctx *gin.Context, nextToken string, size string) (service.GetNotificationsResponse, error) {
//				panic("mock out the GetNotifications method")
//			},
//		}
//
//		// use mockedGetNotificationsService in code that requires GetNotificationsService
//		// and then make assertions.
//
//	}
type GetNotificationsServiceMock struct {
	// GetNotificationsFunc mocks the GetNotifications method.
	GetNotificationsFunc func(ctx *gin.Context, nextToken string, size string) (service.GetNotificationsResponse, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetNotifications holds details about calls to the GetNotifications method.
		GetNotifications []struct {
			// Ctx is the ctx argument value.
			Ctx *gin.Context
			// NextToken is the nextToken argument value.
			NextToken string
			// Size is the size argument value.
			Size string
		}
	}
	lockGetNotifications sync.RWMutex
}

// GetNotifications calls GetNotificationsFunc.
func (mock *GetNotificationsServiceMock) GetNotifications(ctx *gin.Context, nextToken string, size string) (service.GetNotificationsResponse, error) {
	if mock.GetNotificationsFunc == nil {
		panic("GetNotificationsServiceMock.GetNotificationsFunc: method is nil but GetNotificationsService.GetNotifications was just called")
	}
	callInfo := struct {
		Ctx       *gin.Context
		NextToken string
		Size      string
	}{
		Ctx:       ctx,
		NextToken: nextToken,
		Size:      size,
	}
	mock.lockGetNotifications.Lock()
	mock.calls.GetNotifications = append(mock.calls.GetNotifications, callInfo)
	mock.lockGetNotifications.Unlock()
	return mock.GetNotificationsFunc(ctx, nextToken, size)
}

// GetNotificationsCalls gets all the calls that were made to GetNotifications.
// Check the length with:
//
//	len(mockedGetNotificationsService.GetNotificationsCalls())
func (mock *GetNotificationsServiceMock) GetNotificationsCalls() []struct {
	Ctx       *gin.Context
	NextToken string
	Size      string
} {
	var calls []struct {
		Ctx       *gin.Context
		NextToken string
		Size      string
	}
	mock.lockGetNotifications.RLock()
	calls = mock.calls.GetNotifications
	mock.lockGetNotifications.RUnlock()
	return calls
}

// Ensure, that GetUncheckedNotificationCountServiceMock does implement GetUncheckedNotificationCountService.
// If this is not the case, regenerate this file with moq.
var _ GetUncheckedNotificationCountService = &GetUncheckedNotificationCountServiceMock{}

// GetUncheckedNotificationCountServiceMock is a mock implementation of GetUncheckedNotificationCountService.
//
//	func TestSomethingThatUsesGetUncheckedNotificationCountService(t *testing.T) {
//
//		// make and configure a mocked GetUncheckedNotificationCountService
//		mockedGetUncheckedNotificationCountService := &GetUncheckedNotificationCountServiceMock{
//			GetUncheckedNotificationCountFunc: func(ctx *gin.Context) (<-chan int, error) {
//				panic("mock out the GetUncheckedNotificationCount method")
//			},
//		}
//
//		// use mockedGetUncheckedNotificationCountService in code that requires GetUncheckedNotificationCountService
//		// and then make assertions.
//
//	}
type GetUncheckedNotificationCountServiceMock struct {
	// GetUncheckedNotificationCountFunc mocks the GetUncheckedNotificationCount method.
	GetUncheckedNotificationCountFunc func(ctx *gin.Context) (<-chan int, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetUncheckedNotificationCount holds details about calls to the GetUncheckedNotificationCount method.
		GetUncheckedNotificationCount []struct {
			// Ctx is the ctx argument value.
			Ctx *gin.Context
		}
	}
	lockGetUncheckedNotificationCount sync.RWMutex
}

// GetUncheckedNotificationCount calls GetUncheckedNotificationCountFunc.
func (mock *GetUncheckedNotificationCountServiceMock) GetUncheckedNotificationCount(ctx *gin.Context) (<-chan int, error) {
	if mock.GetUncheckedNotificationCountFunc == nil {
		panic("GetUncheckedNotificationCountServiceMock.GetUncheckedNotificationCountFunc: method is nil but GetUncheckedNotificationCountService.GetUncheckedNotificationCount was just called")
	}
	callInfo := struct {
		Ctx *gin.Context
	}{
		Ctx: ctx,
	}
	mock.lockGetUncheckedNotificationCount.Lock()
	mock.calls.GetUncheckedNotificationCount = append(mock.calls.GetUncheckedNotificationCount, callInfo)
	mock.lockGetUncheckedNotificationCount.Unlock()
	return mock.GetUncheckedNotificationCountFunc(ctx)
}

// GetUncheckedNotificationCountCalls gets all the calls that were made to GetUncheckedNotificationCount.
// Check the length with:
//
//	len(mockedGetUncheckedNotificationCountService.GetUncheckedNotificationCountCalls())
func (mock *GetUncheckedNotificationCountServiceMock) GetUncheckedNotificationCountCalls() []struct {
	Ctx *gin.Context
} {
	var calls []struct {
		Ctx *gin.Context
	}
	mock.lockGetUncheckedNotificationCount.RLock()
	calls = mock.calls.GetUncheckedNotificationCount
	mock.lockGetUncheckedNotificationCount.RUnlock()
	return calls
}

// Ensure, that DeleteUserServiceMock does implement DeleteUserService.
// If this is not the case, regenerate this file with moq.
var _ DeleteUserService = &DeleteUserServiceMock{}

// DeleteUserServiceMock is a mock implementation of DeleteUserService.
//
//	func TestSomethingThatUsesDeleteUserService(t *testing.T) {
//
//		// make and configure a mocked DeleteUserService
//		mockedDeleteUserService := &DeleteUserServiceMock{
//			DeleteUserFunc: func(ctx *gin.Context, userID model.UserID) error {
//				panic("mock out the DeleteUser method")
//			},
//		}
//
//		// use mockedDeleteUserService in code that requires DeleteUserService
//		// and then make assertions.
//
//	}
type DeleteUserServiceMock struct {
	// DeleteUserFunc mocks the DeleteUser method.
	DeleteUserFunc func(ctx *gin.Context, userID model.UserID) error

	// calls tracks calls to the methods.
	calls struct {
		// DeleteUser holds details about calls to the DeleteUser method.
		DeleteUser []struct {
			// Ctx is the ctx argument value.
			Ctx *gin.Context
			// UserID is the userID argument value.
			UserID model.UserID
		}
	}
	lockDeleteUser sync.RWMutex
}

// DeleteUser calls DeleteUserFunc.
func (mock *DeleteUserServiceMock) DeleteUser(ctx *gin.Context, userID model.UserID) error {
	if mock.DeleteUserFunc == nil {
		panic("DeleteUserServiceMock.DeleteUserFunc: method is nil but DeleteUserService.DeleteUser was just called")
	}
	callInfo := struct {
		Ctx    *gin.Context
		UserID model.UserID
	}{
		Ctx:    ctx,
		UserID: userID,
	}
	mock.lockDeleteUser.Lock()
	mock.calls.DeleteUser = append(mock.calls.DeleteUser, callInfo)
	mock.lockDeleteUser.Unlock()
	return mock.DeleteUserFunc(ctx, userID)
}

// DeleteUserCalls gets all the calls that were made to DeleteUser.
// Check the length with:
//
//	len(mockedDeleteUserService.DeleteUserCalls())
func (mock *DeleteUserServiceMock) DeleteUserCalls() []struct {
	Ctx    *gin.Context
	UserID model.UserID
} {
	var calls []struct {
		Ctx    *gin.Context
		UserID model.UserID
	}
	mock.lockDeleteUser.RLock()
	calls = mock.calls.DeleteUser
	mock.lockDeleteUser.RUnlock()
	return calls
}
