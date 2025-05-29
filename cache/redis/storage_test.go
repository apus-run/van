package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/apus-run/van/cache/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	errs "github.com/apus-run/van/cache/internal/errs"
	cache "github.com/apus-run/van/cache/redis"
)

func TestCache_Set(t *testing.T) {
	testCases := []struct {
		name string

		mock func(*gomock.Controller) redis.Cmdable

		key        string
		value      string
		expiration time.Duration

		wantErr error
	}{
		{
			name: "set value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetVal("OK")
				cmd.EXPECT().
					Set(context.Background(), "name", "foo", time.Minute).
					Return(status)
				return cmd
			},
			key:        "name",
			value:      "foo",
			expiration: time.Minute,
		},
		{
			name: "timeout",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().
					Set(context.Background(), "name", "foo", time.Minute).
					Return(status)
				return cmd
			},
			key:        "name",
			value:      "foo",
			expiration: time.Minute,

			wantErr: context.DeadlineExceeded,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := cache.New(tc.mock(ctrl))
			err := c.Set(context.Background(), tc.key, tc.value, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestCache_Get(t *testing.T) {
	testCases := []struct {
		name string

		mock func(*gomock.Controller) redis.Cmdable

		key string

		wantErr error
		wantVal string
	}{
		{
			name: "get value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStringCmd(context.Background())
				status.SetVal("foo")
				cmd.EXPECT().
					Get(context.Background(), "name").
					Return(status)
				return cmd
			},
			key: "name",

			wantVal: "foo",
		},
		{
			name: "get error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStringCmd(context.Background())
				status.SetErr(redis.Nil)
				cmd.EXPECT().
					Get(context.Background(), "name").
					Return(status)
				return cmd
			},
			key: "name",

			wantErr: errs.ErrKeyNotExist,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := cache.New(tc.mock(ctrl))
			val, err := c.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val.(string))
		})
	}
}

func TestCache_GetAny(t *testing.T) {
	testCases := []struct {
		name string

		mock func(*gomock.Controller) redis.Cmdable

		key string

		wantErr error
		wantVal string
	}{
		{
			name: "get value",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStringCmd(context.Background())
				status.SetVal("foo")
				cmd.EXPECT().
					Get(context.Background(), "name").
					Return(status)
				return cmd
			},
			key: "name",

			wantVal: "foo",
		},
		{
			name: "get error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStringCmd(context.Background())
				status.SetErr(redis.Nil)
				cmd.EXPECT().
					Get(context.Background(), "name").
					Return(status)
				return cmd
			},
			key: "name",

			wantErr: errs.ErrKeyNotExist,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := cache.New(tc.mock(ctrl))
			val := c.GetAny(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, val.Error)
			if val.Error != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val.Value.(string))
		})
	}
}

func TestCache_Delete(t *testing.T) {
	testCases := []struct {
		name string

		mock func(*gomock.Controller) redis.Cmdable

		key string

		wantN   int64
		wantErr error
	}{
		{
			name: "delete single existed key",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewIntCmd(context.Background())
				status.SetVal(int64(1))
				status.SetErr(nil)
				cmd.EXPECT().
					Del(context.Background(), gomock.Any()).
					Return(status)
				return cmd
			},
			key:   "name",
			wantN: 1,
		},
		{
			name: "delete single does not existed key",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewIntCmd(context.Background())
				status.SetVal(int64(0))
				status.SetErr(nil)
				cmd.EXPECT().
					Del(context.Background(), gomock.Any()).
					Return(status)
				return cmd
			},
			key: "name",
		},
		{
			name: "delete single existed key",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewIntCmd(context.Background())
				status.SetVal(int64(2))
				status.SetErr(nil)
				cmd.EXPECT().
					Del(context.Background(), gomock.Any()).
					Return(status)
				return cmd
			},
			key:   "age",
			wantN: 2,
		},
		{
			name: "delete single do not existed key",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewIntCmd(context.Background())
				status.SetVal(0)
				status.SetErr(nil)
				cmd.EXPECT().
					Del(context.Background(), gomock.Any()).
					Return(status)
				return cmd
			},
			key: "age",
		},
		{
			name: "delete key, some do not existed key",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewIntCmd(context.Background())
				status.SetVal(1)
				status.SetErr(nil)
				cmd.EXPECT().
					Del(context.Background(), gomock.Any()).
					Return(status)
				return cmd
			},
			key:   "name",
			wantN: 1,
		},
		{
			name: "timeout",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewIntCmd(context.Background())
				status.SetVal(0)
				status.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().
					Del(context.Background(), gomock.Any()).
					Return(status)
				return cmd
			},
			key:     "name",
			wantErr: context.DeadlineExceeded,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := cache.New(tc.mock(ctrl))
			err := c.Delete(context.Background(), tc.key)

			assert.Equal(t, tc.wantErr, err)
		})
	}
}
