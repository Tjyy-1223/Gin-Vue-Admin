package handle

import (
	"context"
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/model"
	"gin-blog-server/internal/utils"
	"github.com/gin-gonic/gin"
	"log/slog"
	"strings"
)

type BlogInfo struct{}

type AboutReq struct {
	Content string `json:"content"`
}

// Report 上报用户信息，进行相应的统计操作
// @Summary 上报用户信息
// @Description 用户登进后台时上报信息
// @Tags blog_info
// @Accept json
// @Produce json
// @Param data body object true "用户信息"
// @Success 0 {object} Response[any]
// @Router /report [post]
func (*BlogInfo) Report(c *gin.Context) {
	rdb := GetRDB(c)

	ipAddress := utils.IP.GetIpAddress(c)
	userAgent := utils.IP.GetUserAgent(c)
	browser := userAgent.Name + " " + userAgent.Version.String()
	os := userAgent.OS + " " + userAgent.OSVersion.String()
	uuid := utils.MD5(ipAddress + browser + os)

	ctx := context.Background()

	// 当前用户没有被统计成为访问人数（不在 用户set 中）
	if !rdb.SIsMember(ctx, global.KEY_UNIQUE_VISITOR_SET, uuid).Val() {
		// 统计地域信息: 中国|0|江苏省|苏州市|电信
		ipSource := utils.IP.GetIpSource(ipAddress)
		// 获取到具体的位置, 提取出其中的 省份
		if ipSource != "" {
			address := strings.Split(ipSource, "|")
			province := strings.ReplaceAll(address[2], "省", "")
			rdb.HIncrBy(ctx, global.VISITOR_AREA, province, 1)
		} else {
			rdb.HIncrBy(ctx, global.VISITOR_AREA, "未知", 1)
		}

		// 后台访问数量 + 1
		rdb.Incr(ctx, global.VIEW_COUNT)
		// 将当前用户记录到 用户 set 中
		rdb.SAdd(ctx, global.KEY_UNIQUE_VISITOR_SET, uuid)
	}

	ReturnSuccess(c, nil)
}

// GetConfigMap 获取配置
// @Summary 获取配置信息
// @Description 获取配置信息
// @Tags blog_info
// @Accept json
// @Produce json
// @Param data body object true "配置信息"
// @Success 0 {object} Response[map[string]string]
// @Router /config [get]
func (*BlogInfo) GetConfigMap(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)

	// get from redis cache
	cache, err := getConfigCache(rdb)
	if err != nil {
		ReturnError(c, global.ErrRedisOp, err)
		return
	}

	if len(cache) > 0 {
		slog.Debug("get config from redis cache")
		ReturnSuccess(c, cache)
		return
	}

	// get from db
	data, err := model.GetConfigMap(db)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	// add to redis cache
	if err := addConfigCache(rdb, data); err != nil {
		ReturnError(c, global.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, data)
}

// UpdateConfig 更新配置
// @Summary 更新配置信息
// @Description 更新配置信息
// @Tags blog_info
// @Accept json
// @Produce json
// @Param data body map[string]string true "更新配置信息"
// @Success 0 {object} Response[any]
// @Router /config [patch]
func (*BlogInfo) UpdateConfig(c *gin.Context) {
	var m map[string]string
	if err := c.ShouldBindJSON(&m); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	if err := model.CheckConfigMap(GetDB(c), m); err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	// delete cache
	if err := removeConfigCache(GetRDB(c)); err != nil {
		ReturnError(c, global.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// GetAbout 获取 About 信息
// @Summary 获取关于
// @Description 获取关于
// @Tags blog_info
// @Produce json
// @Success 0 {object} Response[string]
// @Router /about [get]
func (*BlogInfo) GetAbout(c *gin.Context) {
	ReturnSuccess(c, model.GetConfig(GetDB(c), global.CONFIG_ABOUT))
}

// UpdateAbout 更新 About 信息
// @Summary 更新关于
// @Description 更新关于
// @Tags blog_info
// @Accept json
// @Produce json
// @Param data body object true "关于"
// @Success 0 {object} Response[string]
// @Router /about [put]
func (*BlogInfo) UpdateAbout(c *gin.Context) {
	var req AboutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	err := model.CheckConfig(GetDB(c), global.CONFIG_ABOUT, req.Content)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, req.Content)
}
