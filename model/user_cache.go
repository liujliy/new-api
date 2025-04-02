package model

import (
	"encoding/json"
	"fmt"
	"one-api/common"
	"one-api/constant"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/bytedance/gopkg/util/gopool"
)

// UserBase struct remains the same as it represents the cached data structure
type UserBase struct {
	Id               int    `json:"id"`
	Group            string `json:"group"`
	Email            string `json:"email"`
	Quota            int    `json:"quota"`
	Status           int    `json:"status"`
	Username         string `json:"username"`
	Setting          string `json:"setting"`
	StartTimeLimit   int64  `json:"start_time_limit"`
	EndTimeLimit     int64  `json:"end_time_limit"`
	InputLengthLimit int    `json:"input_length_limit"`
	OutputImageLimit int    `json:"output_image_limit"`
}

func (user *UserBase) WriteContext(c *gin.Context) {
	c.Set(constant.ContextKeyUserGroup, user.Group)
	c.Set(constant.ContextKeyUserQuota, user.Quota)
	c.Set(constant.ContextKeyUserStatus, user.Status)
	c.Set(constant.ContextKeyUserEmail, user.Email)
	c.Set("username", user.Username)
	c.Set(constant.ContextKeyUserSetting, user.GetSetting())
	c.Set("user_start_time_limit", user.StartTimeLimit)
	c.Set("user_end_time_limit", user.EndTimeLimit)
	c.Set("user_input_length_limit", user.InputLengthLimit)
	c.Set("user_output_image_limit", user.OutputImageLimit)
}

func (user *UserBase) GetSetting() map[string]interface{} {
	if user.Setting == "" {
		return nil
	}
	return common.StrToMap(user.Setting)
}

func (user *UserBase) SetSetting(setting map[string]interface{}) {
	settingBytes, err := json.Marshal(setting)
	if err != nil {
		common.SysError("failed to marshal setting: " + err.Error())
		return
	}
	user.Setting = string(settingBytes)
}

// getUserCacheKey returns the key for user cache
func getUserCacheKey(userId int) string {
	return fmt.Sprintf("user:%d", userId)
}

// invalidateUserCache clears user cache
func invalidateUserCache(userId int) error {
	if !common.RedisEnabled {
		return nil
	}
	return common.RedisHDelObj(getUserCacheKey(userId))
}

// updateUserCache updates all user cache fields using hash
func updateUserCache(user User) error {
	if !common.RedisEnabled {
		return nil
	}

	return common.RedisHSetObj(
		getUserCacheKey(user.Id),
		user.ToBaseUser(),
		time.Duration(constant.UserId2QuotaCacheSeconds)*time.Second,
	)
}

// GetUserCache gets complete user cache from hash
func GetUserCache(userId int) (userCache *UserBase, err error) {
	var user *User
	var fromDB bool
	defer func() {
		// Update Redis cache asynchronously on successful DB read
		if shouldUpdateRedis(fromDB, err) && user != nil {
			gopool.Go(func() {
				if err := updateUserCache(*user); err != nil {
					common.SysError("failed to update user status cache: " + err.Error())
				}
			})
		}
	}()

	// Try getting from Redis first
	userCache, err = cacheGetUserBase(userId)
	if err == nil {
		return userCache, nil
	}

	// If Redis fails, get from DB
	fromDB = true
	user, err = GetUserById(userId, false)
	if err != nil {
		return nil, err // Return nil and error if DB lookup fails
	}

	// Create cache object from user data
	userCache = &UserBase{
		Id:               user.Id,
		Group:            user.Group,
		Quota:            user.Quota,
		Status:           user.Status,
		Username:         user.Username,
		Setting:          user.Setting,
		Email:            user.Email,
		InputLengthLimit: user.InputLengthLimit,
		OutputImageLimit: user.OutputImageLimit,
	}
	if user.StartTimeLimit != nil {
		userCache.StartTimeLimit = user.StartTimeLimit.Unix()
	}
	if user.EndTimeLimit != nil {
		userCache.EndTimeLimit = user.EndTimeLimit.Unix()
	}

	return userCache, nil
}

func cacheGetUserBase(userId int) (*UserBase, error) {
	if !common.RedisEnabled {
		return nil, fmt.Errorf("redis is not enabled")
	}
	var userCache UserBase
	// Try getting from Redis first
	err := common.RedisHGetObj(getUserCacheKey(userId), &userCache)
	if err != nil {
		return nil, err
	}
	return &userCache, nil
}

// Add atomic quota operations using hash fields
func cacheIncrUserQuota(userId int, delta int64) error {
	if !common.RedisEnabled {
		return nil
	}
	return common.RedisHIncrBy(getUserCacheKey(userId), "Quota", delta)
}

func cacheDecrUserQuota(userId int, delta int64) error {
	return cacheIncrUserQuota(userId, -delta)
}

// Helper functions to get individual fields if needed
func getUserGroupCache(userId int) (string, error) {
	cache, err := GetUserCache(userId)
	if err != nil {
		return "", err
	}
	return cache.Group, nil
}

func getUserQuotaCache(userId int) (int, error) {
	cache, err := GetUserCache(userId)
	if err != nil {
		return 0, err
	}
	return cache.Quota, nil
}

func getUserStatusCache(userId int) (int, error) {
	cache, err := GetUserCache(userId)
	if err != nil {
		return 0, err
	}
	return cache.Status, nil
}

func getUserNameCache(userId int) (string, error) {
	cache, err := GetUserCache(userId)
	if err != nil {
		return "", err
	}
	return cache.Username, nil
}

func getUserSettingCache(userId int) (map[string]interface{}, error) {
	setting := make(map[string]interface{})
	cache, err := GetUserCache(userId)
	if err != nil {
		return setting, err
	}
	return cache.GetSetting(), nil
}

// New functions for individual field updates
func updateUserStatusCache(userId int, status bool) error {
	if !common.RedisEnabled {
		return nil
	}
	statusInt := common.UserStatusEnabled
	if !status {
		statusInt = common.UserStatusDisabled
	}
	return common.RedisHSetField(getUserCacheKey(userId), "Status", fmt.Sprintf("%d", statusInt))
}

func updateUserQuotaCache(userId int, quota int) error {
	if !common.RedisEnabled {
		return nil
	}
	return common.RedisHSetField(getUserCacheKey(userId), "Quota", fmt.Sprintf("%d", quota))
}

func updateUserGroupCache(userId int, group string) error {
	if !common.RedisEnabled {
		return nil
	}
	return common.RedisHSetField(getUserCacheKey(userId), "Group", group)
}

func updateUserNameCache(userId int, username string) error {
	if !common.RedisEnabled {
		return nil
	}
	return common.RedisHSetField(getUserCacheKey(userId), "Username", username)
}

func updateUserSettingCache(userId int, setting string) error {
	if !common.RedisEnabled {
		return nil
	}
	return common.RedisHSetField(getUserCacheKey(userId), "Setting", setting)
}
