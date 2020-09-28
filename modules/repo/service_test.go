package repo_test

import (
	"context"
	"testing"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/mock"
	"github.com/covergates/covergates/modules/repo"
	"github.com/golang/mock/gomock"
)

func TestSynchronize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config := &config.Config{
		Github: config.Github{Server: "url"},
	}

	// test data
	user := &core.User{Login: "user"}
	repos := []*core.Repo{
		{
			URL: "http://github/repo1",
		},
		{
			URL: "http://github/repo2",
		},
	}

	// mock
	mockSCM := mock.NewMockSCMService(ctrl)
	mockClient := mock.NewMockClient(ctrl)
	mockGitRepoService := mock.NewMockGitRepoService(ctrl)
	mockUserStore := mock.NewMockUserStore(ctrl)
	mockRepoStore := mock.NewMockRepoStore(ctrl)

	mockSCM.EXPECT().Client(gomock.Eq(core.Github)).Return(mockClient, nil)
	mockClient.EXPECT().Repositories().Return(mockGitRepoService)
	mockGitRepoService.EXPECT().List(gomock.Any(), gomock.Eq(user)).Return(repos, nil)
	mockRepoStore.EXPECT().BatchUpdateOrCreate(gomock.Eq(repos)).Return(nil)
	mockUserStore.EXPECT().UpdateRepositories(gomock.Eq(user), gomock.Eq(repos)).Return(nil)
	// testing
	service := repo.NewService(config, mockSCM, mockUserStore, mockRepoStore)
	if err := service.Synchronize(context.Background(), user); err != nil {
		t.Fatal(err)
	}
}
