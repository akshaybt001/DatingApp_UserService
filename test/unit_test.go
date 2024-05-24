package userServiceTest

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/akshaybt001/DatingApp_UserService/entities"
	helperstruct "github.com/akshaybt001/DatingApp_UserService/entities/helperStruct"
	mock_adapters "github.com/akshaybt001/DatingApp_UserService/internal/adapters/mockAdapters"
	"github.com/akshaybt001/DatingApp_UserService/internal/helper"
	"github.com/akshaybt001/DatingApp_UserService/internal/service"
	mock_usecases "github.com/akshaybt001/DatingApp_UserService/internal/usecases/mockUsecase"
	"github.com/akshaybt001/DatingApp_proto_files/pb"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adapter := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	userService := service.NewUserService(adapter, usecase)
	hashedPass, err := helper.HashPassword("valid")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	testUUID := uuid.New()
	tests := []struct {
		name               string
		request            *pb.LoginRequest
		mockGetUserByEmail func(string) (entities.User, error)
		wantError          bool
		expectedResult     *pb.UserSignupResponse
	}{
		{
			name: "Success",
			request: &pb.LoginRequest{
				Email:    "valid@gmail.com",
				Password: "valid",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{
					ID:           testUUID,
					Name:         "valid",
					Email:        "valid@gmail.com",
					Phone:        "8888888888",
					Password:     hashedPass,
					IsSubscribed: true,
				}, nil
			},
			wantError: false,
			expectedResult: &pb.UserSignupResponse{
				Id:    testUUID.String(),
				Name:  "valid",
				Email: "valid@gmail.com",
				Phone: "8888888888",
			},
		},
		{
			name: "Fail",
			request: &pb.LoginRequest{
				Email:    "invalid",
				Password: "invalid",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{}, nil
			},
			wantError:      true,
			expectedResult: &pb.UserSignupResponse{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter.EXPECT().GetUserByEmail(gomock.Any()).DoAndReturn(test.mockGetUserByEmail).AnyTimes().Times(1)
			result, err := userService.UserLogin(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
				if err == nil {
					t.Errorf("expected an error but didn't find one")
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}
}

func TestAdminLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adapter := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	userService := service.NewUserService(adapter, usecase)
	hashedPass, err := helper.HashPassword("valid")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	testUUID := uuid.New()
	tests := []struct {
		name               string
		request            *pb.LoginRequest
		mockGetUserByEmail func(string) (entities.User, error)
		wantError          bool
		expectedResult     *pb.UserSignupResponse
	}{
		{
			name: "Success",
			request: &pb.LoginRequest{
				Email:    "valid@gmail.com",
				Password: "valid",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{
					ID:           testUUID,
					Name:         "valid",
					Email:        "valid@gmail.com",
					Phone:        "8888888888",
					Password:     hashedPass,
					IsSubscribed: true,
				}, nil
			},
			wantError: false,
			expectedResult: &pb.UserSignupResponse{
				Id:    testUUID.String(),
				Name:  "valid",
				Email: "valid@gmail.com",
				Phone: "8888888888",
			},
		},
		{
			name: "Fail",
			request: &pb.LoginRequest{
				Email:    "invalid",
				Password: "invalid",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{}, nil
			},
			wantError:      true,
			expectedResult: &pb.UserSignupResponse{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter.EXPECT().GetUserByEmail(gomock.Any()).DoAndReturn(test.mockGetUserByEmail).AnyTimes().Times(1)
			result, err := userService.UserLogin(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
				if err == nil {
					t.Errorf("expected an error but didn't find one")
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adapter := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	userSerive := service.NewUserService(adapter, usecase)
	tests := []struct {
		name               string
		request            *pb.UserSignupRequest
		mockGetUserByEmail func(string) (entities.User, error)
		mockGetUserByPhone func(string) (entities.User, error)
		mockUserSignup     func(entities.User) (entities.User, error)
		wantError          bool
		expectedResult     *pb.UserSignupResponse
	}{
		{
			name: "Success",
			request: &pb.UserSignupRequest{
				Email:    "valid@gmail.com",
				Name:     "valid",
				Password: "valid",
				Phone:    "8888888888",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{}, nil
			},
			mockGetUserByPhone: func(s string) (entities.User, error) {
				return entities.User{}, nil
			},
			mockUserSignup: func(u entities.User) (entities.User, error) {
				return entities.User{
					Name:  "valid",
					Email: "valid@gmail.com",
					Phone: "8888888888",
				}, nil
			},
			wantError: false,
			expectedResult: &pb.UserSignupResponse{
				Name:  "valid",
				Email: "valid@gmail.com",
				Phone: "8888888888",
			},
		},
		{
			name: "EmailNotUnique",
			request: &pb.UserSignupRequest{
				Email:    "valid@gmail.com",
				Name:     "valid",
				Password: "1234",
				Phone:    "8888888888",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{
					Name:  "valid",
					Email: "valid@gmail.com",
					Phone: "8888888888",
				}, nil
			},
			mockGetUserByPhone: func(s string) (entities.User, error) {
				return entities.User{}, nil
			},
			mockUserSignup: func(u entities.User) (entities.User, error) {
				return entities.User{}, nil
			},
			wantError: true,
			expectedResult: &pb.UserSignupResponse{
				Email: "valid@gmail.com",
				Name:  "valid",
				Phone: "8888888888",
			},
		},
		{
			name: "PhoneNotUnique",
			request: &pb.UserSignupRequest{
				Email:    "valid@gmail.com",
				Name:     "valid",
				Password: "valid",
				Phone:    "8888888888",
			},
			mockGetUserByEmail: func(s string) (entities.User, error) {
				return entities.User{}, nil
			},
			mockGetUserByPhone: func(s string) (entities.User, error) {
				return entities.User{
					Name:  "valid",
					Email: "valid@gmail.com",
					Phone: "88888888888",
				}, nil
			},
			mockUserSignup: func(u entities.User) (entities.User, error) {
				return entities.User{}, nil
			},
			wantError: true,
			expectedResult: &pb.UserSignupResponse{
				Name:  "valid",
				Email: "valid@gmail.com",
				Phone: "8888888888",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter.EXPECT().GetUserByEmail(gomock.Any()).DoAndReturn(test.mockGetUserByEmail).Times(1)
			if test.name != "EmailNotUnique" {
				adapter.EXPECT().GetUserByPhone(gomock.Any()).DoAndReturn(test.mockGetUserByPhone).AnyTimes().Times(1)
			}
			if !test.wantError {
				adapter.EXPECT().UserSignup(gomock.Any()).DoAndReturn(test.mockUserSignup).AnyTimes().Times(1)

			}
			res, err := userSerive.UserSignup(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				res.Id = ""
				assert.NotNil(t, res)
				assert.Equal(t, test.expectedResult, res)
			}
		})
	}
}

func TestAddInterstAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adapter := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	userService := service.NewUserService(adapter, usecase)
	tests := []struct {
		name                  string
		request               *pb.AddInterestRequest
		mockGetInterestByName func(string) (entities.Interests, error)
		wantError             bool
	}{
		{
			name: "Success",
			request: &pb.AddInterestRequest{
				Interest: "valid",
			},
			mockGetInterestByName: func(s string) (entities.Interests, error) {
				return entities.Interests{}, nil
			},
			wantError: false,
		},
		{
			name: "Fail",
			request: &pb.AddInterestRequest{
				Interest: "valid",
			},
			mockGetInterestByName: func(s string) (entities.Interests, error) {
				return entities.Interests{
					Interest: "valid",
				}, nil
			},
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter.EXPECT().GetInterestByName(gomock.Any()).DoAndReturn(test.mockGetInterestByName).AnyTimes().Times(1)
			if !test.wantError {
				adapter.EXPECT().AdminAddInterest(gomock.Any()).Return(nil).Times(1)
			}
			_, err := userService.AdminAddInterest(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdminUpdateInterest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adapters := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	userService := service.NewUserService(adapters, usecase)
	tests := []struct {
		name                  string
		request               *pb.InterestResponse
		mockGetInterestByName func(string) (entities.Interests, error)
		wantError             bool
	}{
		{
			name: "Success",
			request: &pb.InterestResponse{
				Id:       1,
				Interest: "valid",
			},
			mockGetInterestByName: func(s string) (entities.Interests, error) {
				return entities.Interests{}, nil
			},
			wantError: false,
		},
		{
			name: "Fail",
			request: &pb.InterestResponse{
				Id:       1,
				Interest: "valid",
			},
			mockGetInterestByName: func(s string) (entities.Interests, error) {
				return entities.Interests{
					Interest: "valid",
				}, nil
			},
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapters.EXPECT().GetInterestByName(gomock.Any()).DoAndReturn(test.mockGetInterestByName).AnyTimes().Times(1)
			if !test.wantError {
				adapters.EXPECT().AdminUpdateInterest(gomock.Any()).Return(nil).Times(1)
			}
			_, err := userService.AdminUpdateInterest(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddInterestUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapters := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	userService := service.NewUserService(adapters, usecase)
	testUUID := uuid.New()
	profileTestUUID := uuid.New()
	tests := []struct {
		name                     string
		request                  *pb.DeleteInterestRequest
		mockGetInterestById      func(int) (helperstruct.InterestHelper, error)
		mockGetProfileIdByUserId func(string) (string, error)
		mockGetUserInterestById  func(string, int) (entities.UserInterests, error)
		mockUserAddInterest      func(entities.UserInterests) error
		wantError                bool
	}{
		{
			name: "Success",
			request: &pb.DeleteInterestRequest{
				InterestId: int32(1),
				UserId:     testUUID.String(),
			},
			mockGetInterestById: func(i int) (helperstruct.InterestHelper, error) {
				return helperstruct.InterestHelper{
					InterestId: 1,
				}, nil
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockGetUserInterestById: func(s string, i int) (entities.UserInterests, error) {
				return entities.UserInterests{}, nil
			},
			mockUserAddInterest: func(ui entities.UserInterests) error {
				return nil
			},
			wantError: false,
		},
		{
			name: "Fail - Invalid InterestId",
			request: &pb.DeleteInterestRequest{
				InterestId: -1,
				UserId:     testUUID.String(),
			},
			mockGetInterestById: func(i int) (helperstruct.InterestHelper, error) {
				return helperstruct.InterestHelper{
					InterestId: 0,
				}, nil
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return "", errors.New("profile not found")
			},
			mockGetUserInterestById: func(s string, i int) (entities.UserInterests, error) {
				return entities.UserInterests{}, nil
			},
			wantError: true,
		},
		// {
		// 	name: "Fail - User interest already exist",
		// 	request: &pb.DeleteInterestRequest{
		// 		InterestId: 1,
		// 		UserId:     testUUID.String(),
		// 	},
		// 	mockGetInterestById: func(i int) (helperstruct.InterestHelper, error) {
		// 		return helperstruct.InterestHelper{
		// 			InterestId: 1,
		// 		}, nil
		// 	},
		// 	mockGetProfileIdByUserId: func(s string) (string, error) {
		// 		return profileTestUUID.String(), nil
		// 	},
		// 	mockGetUserInterestById: func(s string, i int) (entities.UserInterests, error) {
		// 		return entities.UserInterests{
		// 			Id: 1,
		// 		}, errors.New("something ")
		// 	},
		// 	wantError: true,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapters.EXPECT().GetInterestById(int(test.request.InterestId)).DoAndReturn(test.mockGetInterestById).AnyTimes()
			adapters.EXPECT().GetProfileIdByUserId(test.request.UserId).DoAndReturn(test.mockGetProfileIdByUserId).AnyTimes()
			adapters.EXPECT().GetUserInterestById(gomock.Any(), int(test.request.InterestId)).DoAndReturn(test.mockGetUserInterestById).AnyTimes()
			if !test.wantError {
				adapters.EXPECT().UserAddInterest(gomock.Any()).DoAndReturn(test.mockUserAddInterest).AnyTimes()
			}

			_, err := userService.AddInterestUser(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddGenderUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapters := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	userService := service.NewUserService(adapters, usecase)
	testUUID := uuid.New()
	profileTestUUID := uuid.New()
	tests := []struct {
		name                     string
		request                  *pb.UpdateGenderRequest
		mockGetGenderById        func(int) (helperstruct.GenderHelper, error)
		mockGetProfileIdByUserId func(string) (string, error)
		mockGetUserGenderById    func(string, int) (entities.UserGenders, error)
		mockUserAddGender        func(entities.UserGenders) error
		wantError                bool
	}{
		{
			name: "Success",
			request: &pb.UpdateGenderRequest{
				GenderId: 1,
				UserId:   testUUID.String(),
			},
			mockGetGenderById: func(i int) (helperstruct.GenderHelper, error) {
				return helperstruct.GenderHelper{
					GenderId: 1,
				}, nil
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockGetUserGenderById: func(s string, i int) (entities.UserGenders, error) {
				return entities.UserGenders{}, nil
			},
			mockUserAddGender: func(ui entities.UserGenders) error {
				return nil
			},
			wantError: false,
		},
		{
			name: "Fail - Invalid InterestId",
			request: &pb.UpdateGenderRequest{
				GenderId: 3,
				UserId:   testUUID.String(),
			},
			mockGetGenderById: func(i int) (helperstruct.GenderHelper, error) {
				return helperstruct.GenderHelper{
					GenderId: 0,
				}, nil
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return "", errors.New("profile not found")
			},
			mockGetUserGenderById: func(s string, i int) (entities.UserGenders, error) {
				return entities.UserGenders{}, nil
			},
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapters.EXPECT().GetGenderById(int(test.request.GenderId)).DoAndReturn(test.mockGetGenderById).AnyTimes()
			adapters.EXPECT().GetProfileIdByUserId(test.request.UserId).DoAndReturn(test.mockGetProfileIdByUserId).AnyTimes()
			adapters.EXPECT().GetUserGenderById(gomock.Any(), int(test.request.GenderId)).DoAndReturn(test.mockGetUserGenderById).AnyTimes()
			if !test.wantError {
				adapters.EXPECT().UserAddGender(gomock.Any()).DoAndReturn(test.mockUserAddGender).AnyTimes()
			}

			_, err := userService.AddGenderUser(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// func TestDeleteInterestUser(t *testing.T){
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	adapters := mock_adapters.NewMockAdapterInterface(ctrl)
// 	usecase := mock_usecases.NewMockUsecases(ctrl)
// 	userService := service.NewUserService(adapters, usecase)
// 	testUUID := uuid.New()
// 	profileTestUUID := uuid.New()
// 	tests := []struct {
// 		name                     string
// 		request                  *pb.DeleteInterestRequest
// 		mockGetProfileIdByUserId func(string) (string, error)
// 		mockUserDeleteInterest   func(entities.UserInterests) error
// 		wantError                bool
// 	}{
// 		{
// 			name: "Success",
// 			request: &pb.DeleteInterestRequest{
// 				InterestId: 1,
// 				UserId: testUUID.String(),
// 			},
// 			mockGetProfileIdByUserId: func(s string) (string, error) {
// 				return profileTestUUID.String(), nil
// 			},
// 			mockUserDeleteInterest: func(ui entities.UserInterests) error {
// 				return nil
// 			},
// 		},
// 		{
// 			name: "",
// 		},
// 	}
// }

func TestDeleteInterestUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapters := mock_adapters.NewMockAdapterInterface(ctrl)
	userService := service.NewUserService(mockAdapters, nil)

	testUUID := uuid.New()
	profileTestUUID := uuid.New()

	tests := []struct {
		name                     string
		request                  *pb.DeleteInterestRequest
		mockGetProfileIdByUserId func(string) (string, error)
		mockUserDeleteInterest   func(entities.UserInterests) error
		wantError                bool
	}{
		{
			name: "Success",
			request: &pb.DeleteInterestRequest{
				InterestId: 1,
				UserId:     testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockUserDeleteInterest: func(ui entities.UserInterests) error {
				return nil
			},
			wantError: false,
		},
		{
			name: "Fail - GetProfileIdByUserId error",
			request: &pb.DeleteInterestRequest{
				InterestId: 1,
				UserId:     testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return "", fmt.Errorf("profile not found")
			},
			mockUserDeleteInterest: func(ui entities.UserInterests) error {
				return nil
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAdapters.EXPECT().GetProfileIdByUserId(test.request.UserId).Return(test.mockGetProfileIdByUserId(test.request.UserId)).AnyTimes().Times(1)
			if !test.wantError {
				mockAdapters.EXPECT().UserDeleteInterest(gomock.Any()).DoAndReturn(test.mockUserDeleteInterest).AnyTimes().Times(1)
			}

			_, err := userService.DeleteInterestUser(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserAddAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	adapters := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	defer ctrl.Finish()
	userService := service.NewUserService(adapters, usecase)
	testUUID := uuid.New()
	profileTestUUID := uuid.New()
	tests := []struct {
		name                      string
		request                   *pb.AddAddressRequest
		mockGetProfileIdByUserId  func(string) (string, error)
		mockGetAddressByProfileId func(string) (entities.Address, error)
		mockUserAddAddress        func(entities.Address) error
		wantError                 bool
	}{
		{
			name: "Success",
			request: &pb.AddAddressRequest{
				Country:  "validCountry",
				State:    "validState",
				District: "validDistrict",
				City:     "validCity",
				UserId:   testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockGetAddressByProfileId: func(s string) (entities.Address, error) {
				return entities.Address{}, nil
			},
			mockUserAddAddress: func(a entities.Address) error {
				return nil
			},
			wantError: false,
		},
		{
			name: "Fail - address already exist",
			request: &pb.AddAddressRequest{
				Country:  "validCountry",
				State:    "validState",
				District: "validDistrict",
				City:     "validCity",
				UserId:   testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockGetAddressByProfileId: func(s string) (entities.Address, error) {
				return entities.Address{
					Country:  "validCountry",
					State:    "validState",
					District: "validDistrict",
					City:     "validCity",
				}, nil
			},
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapters.EXPECT().GetProfileIdByUserId(gomock.Any()).DoAndReturn(test.mockGetProfileIdByUserId).AnyTimes().Times(1)
			adapters.EXPECT().GetAddressByProfileId(gomock.Any()).DoAndReturn(test.mockGetAddressByProfileId).AnyTimes().Times(1)
			if !test.wantError {
				adapters.EXPECT().UserAddAddress(gomock.Any()).DoAndReturn(test.mockUserAddAddress).AnyTimes().Times(1)
			}
			_, err := userService.UserAddAddress(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestUserAddPreference(t *testing.T) {
	ctrl := gomock.NewController(t)
	adapters := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	defer ctrl.Finish()
	userService := service.NewUserService(adapters, usecase)
	testUUID := uuid.New()
	profileTestUUID := uuid.New()
	tests := []struct {
		name                         string
		request                      *pb.PreferenceRequest
		mockGetProfileIdByUserId     func(string) (string, error)
		mockGetPreferenceByProfileId func(string) (entities.Preference, error)
		mockUserAddPreference        func(entities.Preference) error
		wantError                    bool
	}{
		{
			name: "Success",
			request: &pb.PreferenceRequest{
				Minage:     18,
				Maxage:     26,
				Gender:     int32(1),
				Desirecity: "validCity",
				UserId:     testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockGetPreferenceByProfileId: func(s string) (entities.Preference, error) {
				return entities.Preference{}, nil
			},
			mockUserAddPreference: func(a entities.Preference) error {
				return nil
			},
			wantError: false,
		},
		{
			name: "Fail - address already exist",
			request: &pb.PreferenceRequest{
				Minage:     18,
				Maxage:     26,
				Gender:     int32(1),
				Desirecity: "validCity",
				UserId:     testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockGetPreferenceByProfileId: func(s string) (entities.Preference, error) {
				return entities.Preference{
					MinAge:     18,
					MaxAge:     26,
					GenderId:   1,
					DesireCity: "validCity",
					ProfileId:  profileTestUUID,
				}, nil
			},
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapters.EXPECT().GetProfileIdByUserId(gomock.Any()).DoAndReturn(test.mockGetProfileIdByUserId).AnyTimes().Times(1)
			adapters.EXPECT().GetPreferenceByProfileId(gomock.Any()).DoAndReturn(test.mockGetPreferenceByProfileId).AnyTimes().Times(1)
			if !test.wantError {
				adapters.EXPECT().UserAddPreference(gomock.Any()).DoAndReturn(test.mockUserAddPreference).AnyTimes().Times(1)
			}
			_, err := userService.UserAddPreference(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestAdminAddGender(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	adapters := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	userService := service.NewUserService(adapters, usecase)
	tests := []struct {
		name                string
		request             *pb.AddGenderRequest
		mockGetGenderByName func(string) (entities.Gender, error)
		mockAdminAddGender  func(entities.Gender) error
		wantError           bool
	}{
		{
			name: "Success",
			request: &pb.AddGenderRequest{
				Gender: "ValidMale",
			},
			mockGetGenderByName: func(s string) (entities.Gender, error) {
				return entities.Gender{}, nil
			},
			mockAdminAddGender: func(g entities.Gender) error {
				return nil
			},
			wantError: false,
		},
		{
			name: "Fail",
			request: &pb.AddGenderRequest{
				Gender: "ValidMale",
			},
			mockGetGenderByName: func(s string) (entities.Gender, error) {
				return entities.Gender{
					Name: "ValidMale",
				}, nil
			},
			mockAdminAddGender: func(g entities.Gender) error {
				return nil
			},
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapters.EXPECT().GetGenderByName(gomock.Any()).DoAndReturn(test.mockGetGenderByName).AnyTimes().Times(1)
			if !test.wantError {
				adapters.EXPECT().AdminAddGender(gomock.Any()).DoAndReturn(test.mockAdminAddGender).AnyTimes().Times(1)
			}
			_, err := userService.AdminAddGender(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// func TestHomePage(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	adapters := mock_adapters.NewMockAdapterInterface(ctrl)
// 	usecase := mock_usecases.NewMockUsecases(ctrl)
// 	userSerivce := service.NewUserService(adapters, usecase)
// 	testUUID := uuid.New()
// 	profileTestUUID := uuid.New()
// 	temp1 := uuid.New()
// 	temp2 := uuid.New()
// 	tests := []struct {
// 		name                     string
// 		request                  *pb.GetUserById
// 		mockGetProfileIdByUserId func(string) (string, error)
// 		mockFetchPreference      func(string) (helperstruct.FetchPreference, error)
// 		mockFetchUser            func(string) (helperstruct.FetchUser, error)
// 		mockFetchInterests       func(string) ([]string, error)
// 		mockFetchUsers           func(maxage, minage, gender int, id string) ([]helperstruct.Home, error)
// 		mockFetchImages          func(string) (string, error)
// 		wantError                bool
// 		expectedResult           *pb.HomeResponse
// 	}{
// 		{
// 			name: "Success",
// 			request: &pb.GetUserById{
// 				Id: testUUID.String(),
// 			},
// 			mockGetProfileIdByUserId: func(s string) (string, error) {
// 				return profileTestUUID.String(), nil
// 			},
// 			mockFetchPreference: func(s string) (helperstruct.FetchPreference, error) {
// 				return helperstruct.FetchPreference{
// 					MinAge:     18,
// 					MaxAge:     26,
// 					Gender:     2,
// 					DesireCity: "validCity",
// 				}, nil
// 			},
// 			mockFetchUser: func(s string) (helperstruct.FetchUser, error) {
// 				return helperstruct.FetchUser{
// 					Age: 22,
// 				}, nil
// 			},
// 			mockFetchInterests: func(s string) ([]string, error) {
// 				return []string{}, nil
// 			},
// 			mockFetchUsers: func(maxage, minage, gender int, id string) ([]helperstruct.Home, error) {
// 				return []helperstruct.Home{
// 					helperstruct.Home{
// 						Id:      temp1.String(),
// 						Name:    "Nobody",
// 						Age:     21,
// 						Gender:  "Female",
// 						City:    "validCity",
// 						Country: "validCountry",
// 					},
// 					helperstruct.Home{
// 						Id:      temp2.String(),
// 						Name:    "Nobody2",
// 						Age:     21,
// 						Gender:  "Female",
// 						City:    "validCity",
// 						Country: "validCountry",
// 					},
// 				}, nil
// 			},
// 		},
// 	}
// }

func TestUserEditAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapters := mock_adapters.NewMockAdapterInterface(ctrl)
	usecase := mock_usecases.NewMockUsecases(ctrl)
	userService := service.NewUserService(adapters, usecase)

	testUUID := uuid.New()
	profileTestUUID := uuid.New()
	addressUUID := uuid.New()

	tests := []struct {
		name                     string
		request                  *pb.AddressResponse
		mockGetProfileIdByUserId func(string) (string, error)
		mockUserEditAddress      func(entities.Address) error
		wantError                bool
	}{
		{
			name: "Success",
			request: &pb.AddressResponse{
				Id:       addressUUID.String(),
				Country:  "Country",
				State:    "State",
				District: "District",
				City:     "City",
				UserId:   testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockUserEditAddress: func(a entities.Address) error {
				return nil
			},
			wantError: false,
		},
		{
			name: "Fail",
			request: &pb.AddressResponse{
				Id:       "invalid-address-id",
				Country:  "Country",
				State:    "State",
				District: "District",
				City:     "City",
				UserId:   testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return "", errors.New("profile not found")
			},
			mockUserEditAddress: func(a entities.Address) error {
				return errors.New("edit address failed")
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapters.EXPECT().GetProfileIdByUserId(test.request.UserId).Return(test.mockGetProfileIdByUserId(test.request.UserId)).Times(1)
			if !test.wantError {
				adapters.EXPECT().UserEditAddress(gomock.Any()).DoAndReturn(test.mockUserEditAddress).Times(1)
			}

			_, err := userService.UserEditAddress(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserEditPreference(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the adapters interface
	mockAdapters := mock_adapters.NewMockAdapterInterface(ctrl)
	userService := service.NewUserService(mockAdapters, nil)

	testUUID := uuid.New()
	profileTestUUID := uuid.New()

	tests := []struct {
		name                     string
		request                  *pb.PreferenceResponse
		mockGetProfileIdByUserId func(string) (string, error)
		mockUserEditPreference   func(entities.Preference) error
		wantError                bool
	}{
		{
			name: "Success",
			request: &pb.PreferenceResponse{
				UserId:     testUUID.String(),
				Minage:     25,
				Maxage:     35,
				Gender:     1,
				Desirecity: "DesireCity",
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockUserEditPreference: func(entities.Preference) error {
				return nil
			},
			wantError: false,
		},
		{
			name: "Fail - GetProfileIdByUserId error",
			request: &pb.PreferenceResponse{
				UserId:     testUUID.String(),
				Minage:     25,
				Maxage:     35,
				Gender:     1,
				Desirecity: "DesireCity",
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return "", fmt.Errorf("profile not found")
			},
			mockUserEditPreference: func(entities.Preference) error {
				return nil
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAdapters.EXPECT().GetProfileIdByUserId(test.request.UserId).Return(test.mockGetProfileIdByUserId(test.request.UserId)).Times(1)
			if !test.wantError {
				mockAdapters.EXPECT().UserEditPreference(gomock.Any()).DoAndReturn(test.mockUserEditPreference).Times(1)
			}

			_, err := userService.UserEditPreference(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapters := mock_adapters.NewMockAdapterInterface(ctrl)
	userService := service.NewUserService(mockAdapters, nil)

	testUUID := uuid.New()
	userData := entities.User{
		ID:    testUUID,
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Phone: "1234567890",
	}

	tests := []struct {
		name            string
		request         *pb.GetUserById
		mockGetUserById func(string) (entities.User, error)
		expectedResult  *pb.UserSignupResponse
		wantError       bool
	}{
		{
			name: "Success",
			request: &pb.GetUserById{
				Id: testUUID.String(),
			},
			mockGetUserById: func(s string) (entities.User, error) {
				return userData, nil
			},
			expectedResult: &pb.UserSignupResponse{
				Id:    testUUID.String(),
				Name:  "John Doe",
				Email: "john.doe@example.com",
				Phone: "1234567890",
			},
			wantError: false,
		},
		{
			name: "Fail - GetUserById error",
			request: &pb.GetUserById{
				Id: testUUID.String(),
			},
			mockGetUserById: func(s string) (entities.User, error) {
				return entities.User{}, fmt.Errorf("user not found")
			},
			expectedResult: nil,
			wantError:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAdapters.EXPECT().GetUserById(gomock.Any()).DoAndReturn(test.mockGetUserById).AnyTimes().Times(1)

			result, err := userService.GetUser(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}
}

func TestUserUploadProfileImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAdapters := mock_adapters.NewMockAdapterInterface(ctrl)
	mockUsecases := mock_usecases.NewMockUsecases(ctrl)
	userService := service.NewUserService(mockAdapters, mockUsecases)

	testUUID := uuid.New()
	profileTestUUID := uuid.New()

	tests := []struct {
		name                     string
		request                  *pb.UserImageRequest
		mockGetProfileIdByUserId func(string) (string, error)
		mockUploadImage          func(*pb.UserImageRequest, string) (string, error)
		expectedResult           *pb.UserImageResponse
		wantError                bool
	}{
		{
			name: "Success",
			request: &pb.UserImageRequest{
				UserId:    testUUID.String(),
				ImageData: []byte("test_image_data"),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockUploadImage: func(req *pb.UserImageRequest, profile string) (string, error) {
				return "http://example.com/image.jpg", nil
			},
			expectedResult: &pb.UserImageResponse{
				Url: "http://example.com/image.jpg",
			},
			wantError: false,
		},
		{
			name: "Fail - GetProfileIdByUserId error",
			request: &pb.UserImageRequest{
				UserId:    testUUID.String(),
				ImageData: []byte("test_image_data"),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return "", fmt.Errorf("profile not found")
			},
			mockUploadImage: func(req *pb.UserImageRequest, profile string) (string, error) {
				return "", nil
			},
			expectedResult: nil,
			wantError:      true,
		},
		{
			name: "Fail - UploadImage error",
			request: &pb.UserImageRequest{
				UserId:    testUUID.String(),
				ImageData: []byte("test_image_data"),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockUploadImage: func(req *pb.UserImageRequest, profile string) (string, error) {
				return "", fmt.Errorf("upload failed")
			},
			expectedResult: nil,
			wantError:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAdapters.EXPECT().GetProfileIdByUserId(gomock.Any()).DoAndReturn(test.mockGetProfileIdByUserId).Times(1)
			if !test.wantError || test.name == "Fail - UploadImage error" {
				mockUsecases.EXPECT().UploadImage(test.request, profileTestUUID.String()).DoAndReturn(test.mockUploadImage).Times(1)
			}

			result, err := userService.UserUploadProfileImage(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}
}

func TestUserGetProfilePic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapters := mock_adapters.NewMockAdapterInterface(ctrl)
	userService := service.NewUserService(mockAdapters, nil)

	testUUID := uuid.New()
	profileTestUUID := uuid.New()

	tests := []struct {
		name                     string
		request                  *pb.GetUserById
		mockGetProfileIdByUserId func(string) (string, error)
		mockGetProfilePic        func(string) (string, error)
		expectedResult           *pb.UserImageResponse
		wantError                bool
	}{
		{
			name: "Success",
			request: &pb.GetUserById{
				Id: testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockGetProfilePic: func(profile string) (string, error) {
				return "http://example.com/profile.jpg", nil
			},
			expectedResult: &pb.UserImageResponse{
				Url: "http://example.com/profile.jpg",
			},
			wantError: false,
		},
		{
			name: "Fail - GetProfileIdByUserId error",
			request: &pb.GetUserById{
				Id: testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return "", fmt.Errorf("profile not found")
			},
			mockGetProfilePic: func(profile string) (string, error) {
				return "", nil
			},
			expectedResult: nil,
			wantError:      true,
		},
		{
			name: "Fail - GetProfilePic error",
			request: &pb.GetUserById{
				Id: testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockGetProfilePic: func(profile string) (string, error) {
				return "", fmt.Errorf("profile pic not found")
			},
			expectedResult: nil,
			wantError:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAdapters.EXPECT().GetProfileIdByUserId(gomock.Any()).DoAndReturn(test.mockGetProfileIdByUserId).Times(1)
			if !test.wantError || test.name == "Fail - GetProfilePic error" {
				mockAdapters.EXPECT().GetProfilePic(gomock.Any()).DoAndReturn(test.mockGetProfilePic).Times(1)
			}

			result, err := userService.UserGetProfilePic(context.Background(), test.request)
			if test.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}
}

func TestUserAddAge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapters := mock_adapters.NewMockAdapterInterface(ctrl)
	userService := service.NewUserService(mockAdapters, nil)

	testUUID := uuid.New()
	profileTestUUID := uuid.New()

	tests := []struct {
		name                     string
		request                  *pb.UserAgeRequest
		mockGetProfileIdByUserId func(string) (string, error)
		mockUpdateAge            func(int, string) error
		expectedError            bool
	}{
		{
			name: "Success",
			request: &pb.UserAgeRequest{
				UserId: testUUID.String(),
				Dob:    "1990-01-01T00:00:00.000000Z",
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockUpdateAge: func(age int, profile string) error {
				return nil
			},
			expectedError: false,
		},
		{
			name: "Fail - GetProfileIdByUserId error",
			request: &pb.UserAgeRequest{
				UserId: testUUID.String(),
				Dob:    "1990-01-01T00:00:00.000000Z",
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return "", fmt.Errorf("profile not found")
			},
			mockUpdateAge: func(age int, profile string) error {
				return nil
			},
			expectedError: true,
		},
		{
			name: "Fail - Invalid DOB format",
			request: &pb.UserAgeRequest{
				UserId: testUUID.String(),
				Dob:    "invalid-date-format",
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockUpdateAge: func(age int, profile string) error {
				return nil
			},
			expectedError: true,
		},
		{
			name: "Fail - UpdateAge error",
			request: &pb.UserAgeRequest{
				UserId: testUUID.String(),
				Dob:    "1990-01-01T00:00:00.000000Z",
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockUpdateAge: func(age int, profile string) error {
				return fmt.Errorf("update age failed")
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAdapters.EXPECT().GetProfileIdByUserId(test.request.UserId).DoAndReturn(test.mockGetProfileIdByUserId).AnyTimes().Times(1)
			if !test.expectedError || test.name == "Fail - UpdateAge error" {
				mockAdapters.EXPECT().UpdateAge(gomock.Any(), profileTestUUID.String()).DoAndReturn(test.mockUpdateAge).AnyTimes().Times(1)
			}

			_, err := userService.UserAddAge(context.Background(), test.request)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserGetAge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapters := mock_adapters.NewMockAdapterInterface(ctrl)
	userService := service.NewUserService(mockAdapters, nil)

	testUUID := uuid.New()
	profileTestUUID := uuid.New()
	expectedAge := 30

	tests := []struct {
		name                     string
		request                  *pb.GetUserById
		mockGetProfileIdByUserId func(string) (string, error)
		mockGetAge               func(string) (int, error)
		expectedAge              int32
		expectedError            bool
	}{
		{
			name: "Success",
			request: &pb.GetUserById{
				Id: testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockGetAge: func(profile string) (int, error) {
				return expectedAge, nil
			},
			expectedAge:   int32(expectedAge),
			expectedError: false,
		},
		{
			name: "Fail - GetProfileIdByUserId error",
			request: &pb.GetUserById{
				Id: testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return "", fmt.Errorf("profile not found")
			},
			mockGetAge: func(profile string) (int, error) {
				return 0, nil
			},
			expectedAge:   0,
			expectedError: true,
		},
		{
			name: "Fail - GetAge error",
			request: &pb.GetUserById{
				Id: testUUID.String(),
			},
			mockGetProfileIdByUserId: func(s string) (string, error) {
				return profileTestUUID.String(), nil
			},
			mockGetAge: func(profile string) (int, error) {
				return 0, fmt.Errorf("could not retrieve age")
			},
			expectedAge:   0,
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAdapters.EXPECT().GetProfileIdByUserId(gomock.Any()).DoAndReturn(test.mockGetProfileIdByUserId).AnyTimes().Times(1)
			if !test.expectedError || test.name == "Fail - GetAge error" {
				mockAdapters.EXPECT().GetAge(gomock.Any()).DoAndReturn(test.mockGetAge).AnyTimes().Times(1)
			}

			res, err := userService.UserGetAge(context.Background(), test.request)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedAge, res.Age)
			}
		})
	}
}

func TestGetUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapters := mock_adapters.NewMockAdapterInterface(ctrl)
	userService := service.NewUserService(mockAdapters, nil)

	testUUID := uuid.New()

	tests := []struct {
		name             string
		request          *pb.GetUserById
		mockGetUserById  func(string) (entities.User, error)
		expectedResponse *pb.UserDataResponse
		expectedError    bool
	}{
		{
			name: "Success",
			request: &pb.GetUserById{
				Id: testUUID.String(),
			},
			mockGetUserById: func(s string) (entities.User, error) {
				return entities.User{
					ID:           testUUID,
					Name:         "John Doe",
					Email:        "john.doe@example.com",
					Phone:        "1234567890",
					IsBlocked:    false,
					LikeCount:    3,
					IsSubscribed: true,
				}, nil
			},
			expectedResponse: &pb.UserDataResponse{
				Id:           testUUID.String(),
				Name:         "John Doe",
				Email:        "john.doe@example.com",
				Phone:        "1234567890",
				IsBlocked:    false,
				LikeCount:    3,
				IsSubscribed: true,
			},
			expectedError: false,
		},
		{
			name: "Fail - GetUserById error",
			request: &pb.GetUserById{
				Id: testUUID.String(),
			},
			mockGetUserById: func(s string) (entities.User, error) {
				return entities.User{}, fmt.Errorf("user not found")
			},
			expectedResponse: nil,
			expectedError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAdapters.EXPECT().GetUserById(tt.request.Id).DoAndReturn(tt.mockGetUserById).AnyTimes().Times(1)

			res, err := userService.GetUserData(context.Background(), tt.request)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, res)
			}
		})
	}
}

func TestDecrementLikeCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapters := mock_adapters.NewMockAdapterInterface(ctrl)
	userService := service.NewUserService(mockAdapters, nil)

	testUUID := uuid.New()

	tests := []struct {
		name                   string
		request                *pb.GetUserById
		mockDecrementLikeCount func(string) error
		expectedError          bool
	}{
		{
			name: "Success",
			request: &pb.GetUserById{
				Id: testUUID.String(),
			},
			mockDecrementLikeCount: func(s string) error {
				return nil
			},
			expectedError: false,
		},
		{
			name: "Fail - DecrementLikeCount error",
			request: &pb.GetUserById{
				Id: testUUID.String(),
			},
			mockDecrementLikeCount: func(s string) error {
				return fmt.Errorf("decrement like count failed")
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAdapters.EXPECT().DecrementLikeCount(gomock.Any()).DoAndReturn(test.mockDecrementLikeCount).AnyTimes().Times(1)

			_, err := userService.DecrementLikeCount(context.Background(), test.request)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateSubscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapters := mock_adapters.NewMockAdapterInterface(ctrl)
	userService := service.NewUserService(mockAdapters, nil)

	testUUID := uuid.New()

	tests := []struct {
		name                   string
		request                *pb.UpdateSubscriptionRequest
		mockUpdateSubscription func(string, bool) error
		expectedError          bool
	}{
		{
			name: "Success",
			request: &pb.UpdateSubscriptionRequest{
				UserId:       testUUID.String(),
				Subscription: true,
			},
			mockUpdateSubscription: func(userId string, subscription bool) error {
				return nil
			},
			expectedError: false,
		},
		{
			name: "Fail - UpdateSubscription error",
			request: &pb.UpdateSubscriptionRequest{
				UserId:       testUUID.String(),
				Subscription: true,
			},
			mockUpdateSubscription: func(userId string, subscription bool) error {
				return fmt.Errorf("update subscription failed")
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAdapters.EXPECT().UpdateSubscription(gomock.Any(), gomock.Any()).DoAndReturn(test.mockUpdateSubscription).AnyTimes().Times(1)

			_, err := userService.UpdateSubscription(context.Background(), test.request)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
