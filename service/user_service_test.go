package service_test

import (
	"context"
	"testing"
	"user-mgmt/domain"
	mock_repository "user-mgmt/repository/mock"
	"user-mgmt/service"
	"user-mgmt/utils"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestCreateUser(t *testing.T) {
	testcases := []struct {
		name      string
		email     string
		password  string
		expected  error
		buildStub func(mockUserRepository *mock_repository.MockUserRepository)
	}{
		{
			name:     "SUCCESS",
			email:    "valid@example.com",
			password: "validpassword",
			expected: nil,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetByEmail(gomock.Any(), "valid@example.com").
					Return(nil, mongo.ErrNoDocuments)
				mockUserRepository.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(domain.NewUser("valid", "valid@example.com", "validpassword"), nil)
			},
		},
		{
			name:     "FAILED:EMAIL_EXISTS",
			email:    "valid@example.com",
			password: "validpassword",
			expected: domain.ErrEmailExists,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetByEmail(gomock.Any(), "valid@example.com").
					Return(domain.NewUser("valid", "valid@example.com", "validpassword"), nil)
			},
		},
		{
			name:     "FAILED:GET_BY_EMAIL_ERROR",
			email:    "valid@example.com",
			password: "validpassword",
			expected: mongo.ErrWrongClient,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetByEmail(gomock.Any(), "valid@example.com").
					Return(nil, mongo.ErrWrongClient)
			},
		},
		{
			name:     "FAILED:CREATE_ERROR",
			email:    "valid@example.com",
			password: "validpassword",
			expected: mongo.ErrWrongClient,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetByEmail(gomock.Any(), "valid@example.com").
					Return(nil, mongo.ErrNoDocuments)
				mockUserRepository.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, mongo.ErrWrongClient)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init test env
			mockUserRepository := mock_repository.NewMockUserRepository(ctrl)

			tc.buildStub(mockUserRepository)

			// init test env
			userService := service.NewUserService(mockUserRepository)
			// test test env
			user, err := userService.Create(context.Background(), tc.email, tc.name, tc.password)

			if tc.expected != nil {
				assert.Error(t, tc.expected, err)
				assert.Nil(t, user)
			} else {
				assert.NotNil(t, user.ID)
				assert.Equal(t, tc.email, user.Email)

				if tc.password != "" {
					// we did not return hash password so compare with raw password
					assert.Equal(t, tc.password, user.Password)
				}
			}

		})
	}
}

func TestGetByID(t *testing.T) {

	mockUUID := uuid.New()

	testcases := []struct {
		name      string
		id        string
		expected  error
		buildStub func(mockUserRepository *mock_repository.MockUserRepository)
	}{
		{
			name:     "SUCCESS",
			id:       mockUUID.String(),
			expected: nil,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetByID(gomock.Any(), mockUUID.String()).
					Return(&domain.User{
						ID:       mockUUID,
						Name:     "valid",
						Email:    "valid@example.com",
						Password: "validpassword",
					}, nil)
			},
		},
		{
			name:     "FAILED:GET_BY_ID_ERROR",
			id:       mockUUID.String(),
			expected: mongo.ErrWrongClient,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetByID(gomock.Any(), mockUUID.String()).
					Return(nil, mongo.ErrWrongClient)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init test env
			mockUserRepository := mock_repository.NewMockUserRepository(ctrl)

			tc.buildStub(mockUserRepository)

			// init test env
			userService := service.NewUserService(mockUserRepository)
			// test test env
			user, err := userService.GetByID(context.Background(), tc.id)

			if tc.expected != nil {
				assert.Error(t, tc.expected, err)
				assert.Nil(t, user)
			} else {
				assert.NotNil(t, user.ID)
				assert.Equal(t, tc.id, user.ID.String())
			}
		})
	}
}

func TestGetUserList(t *testing.T) {

	mockUUID1 := uuid.New()
	mockUUID2 := uuid.New()

	mockUser1 := domain.User{
		ID:       mockUUID1,
		Name:     "valid1",
		Email:    "valid1@example.com",
		Password: "validpassword1",
	}

	mockUser2 := domain.User{
		ID:       mockUUID2,
		Name:     "valid2",
		Email:    "valid2@example.com",
		Password: "validpassword2",
	}

	mockUserList := []*domain.User{
		&mockUser1,
		&mockUser2,
	}

	testcases := []struct {
		name           string
		expectedResult []*domain.User
		expectedError  error
		buildStub      func(mockUserRepository *mock_repository.MockUserRepository)
	}{
		{
			name:           "SUCCESS",
			expectedResult: mockUserList,
			expectedError:  nil,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetUserList(gomock.Any(), 1, 10).
					Return(mockUserList, utils.Pagination{1, 10, int64(len(mockUserList))}, nil)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init test env
			mockUserRepository := mock_repository.NewMockUserRepository(ctrl)

			tc.buildStub(mockUserRepository)

			// init test env
			userService := service.NewUserService(mockUserRepository)
			// test test env
			userList, pagination, err := userService.GetUserList(context.Background(), 1, 10)

			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, err)
				assert.Nil(t, userList)
			} else {
				assert.NotNil(t, userList)
				assert.Equal(t, int64(len(tc.expectedResult)), pagination.Total)
				assert.Equal(t, tc.expectedResult, userList)
			}
		})
	}

}

func TestUserUpdate(t *testing.T) {
	mockUUID := uuid.New()
	mockUser := domain.User{
		ID:       mockUUID,
		Name:     "valid",
		Email:    "valid@example.com",
		Password: "validpassword",
	}

	updateName := "valid2"
	updateEmail := "valid2@example.com"
	mockUpdatedUser := domain.User{
		ID:       mockUUID,
		Name:     updateName,
		Email:    updateEmail,
		Password: "validpassword",
	}

	testcases := []struct {
		name          string
		id            string
		username      string
		email         string
		expectedError error
		buildStub     func(mockUserRepository *mock_repository.MockUserRepository)
	}{
		{
			name:          "SUCCESS",
			id:            mockUUID.String(),
			username:      updateName,
			email:         updateEmail,
			expectedError: nil,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetByID(gomock.Any(), mockUUID.String()).
					Return(&mockUser, nil)
				mockUserRepository.EXPECT().Update(gomock.Any(), &mockUpdatedUser).
					Return(&mockUpdatedUser, nil)
			},
		},
		{
			name:          "FAILED:GET_BY_ID_ERROR",
			id:            mockUUID.String(),
			username:      updateName,
			email:         updateEmail,
			expectedError: mongo.ErrWrongClient,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetByID(gomock.Any(), mockUUID.String()).
					Return(nil, mongo.ErrWrongClient)
			},
		},
		{
			name:          "FAILED:UPDATE_ERROR",
			id:            mockUUID.String(),
			username:      updateName,
			email:         updateEmail,
			expectedError: mongo.ErrWrongClient,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetByID(gomock.Any(), mockUUID.String()).
					Return(&mockUser, nil)
				mockUserRepository.EXPECT().Update(gomock.Any(), &mockUpdatedUser).
					Return(nil, mongo.ErrWrongClient)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init test env
			mockUserRepository := mock_repository.NewMockUserRepository(ctrl)

			tc.buildStub(mockUserRepository)

			// init test env
			userService := service.NewUserService(mockUserRepository)
			// test test env
			updatedUser, err := userService.Update(context.Background(), tc.id, &tc.username, &tc.email)

			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, err)
			} else {
				assert.Equal(t, tc.id, updatedUser.ID.String())
				assert.Equal(t, tc.username, updatedUser.Name)
				assert.Equal(t, tc.email, updatedUser.Email)
			}
		})
	}
}

func TestUserDelete(t *testing.T) {
	mockUUID := uuid.New()
	mockUser := domain.User{
		ID:       mockUUID,
		Name:     "valid",
		Email:    "valid@example.com",
		Password: "validpassword",
	}

	testcases := []struct {
		name          string
		id            string
		expectedError error
		buildStub     func(mockUserRepository *mock_repository.MockUserRepository)
	}{
		{
			name:          "SUCCESS",
			id:            mockUUID.String(),
			expectedError: nil,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetByID(gomock.Any(), mockUUID.String()).
					Return(&mockUser, nil)
				mockUserRepository.EXPECT().Delete(gomock.Any(), mockUUID.String()).
					Return(nil)
			},
		},
		{
			name:          "FAILED:GET_BY_ID_ERROR",
			id:            mockUUID.String(),
			expectedError: mongo.ErrNoDocuments,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetByID(gomock.Any(), mockUUID.String()).
					Return(nil, mongo.ErrNoDocuments)
			},
		},
		{
			name:          "FAILED:DELETE_ERROR",
			id:            mockUUID.String(),
			expectedError: mongo.ErrWrongClient,
			buildStub: func(mockUserRepository *mock_repository.MockUserRepository) {
				mockUserRepository.EXPECT().GetByID(gomock.Any(), mockUUID.String()).
					Return(&mockUser, nil)
				mockUserRepository.EXPECT().Delete(gomock.Any(), mockUUID.String()).
					Return(mongo.ErrWrongClient)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init test env
			mockUserRepository := mock_repository.NewMockUserRepository(ctrl)

			tc.buildStub(mockUserRepository)

			// init test env
			userService := service.NewUserService(mockUserRepository)
			// test test env
			err := userService.Delete(context.Background(), tc.id)

			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
