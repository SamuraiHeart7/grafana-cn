package api

import (
	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/middleware"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/setting"
)

func setIndexViewData(c *middleware.Context) (*dtos.IndexViewData, error) {
	settings, err := getFrontendSettingsMap(c)
	if err != nil {
		return nil, err
	}

	prefsQuery := m.GetPreferencesWithDefaultsQuery{OrgId: c.OrgId, UserId: c.UserId}
	if err := bus.Dispatch(&prefsQuery); err != nil {
		return nil, err
	}
	prefs := prefsQuery.Result

	var data = dtos.IndexViewData{
		User: &dtos.CurrentUser{
			Id:             c.UserId,
			IsSignedIn:     c.IsSignedIn,
			Login:          c.Login,
			Email:          c.Email,
			Name:           c.Name,
			OrgId:          c.OrgId,
			OrgName:        c.OrgName,
			OrgRole:        c.OrgRole,
			GravatarUrl:    dtos.GetGravatarUrl(c.Email),
			IsGrafanaAdmin: c.IsGrafanaAdmin,
			LightTheme:     prefs.Theme == "light",
			Timezone:       prefs.Timezone,
		},
		Settings:           settings,
		AppUrl:             setting.AppUrl,
		AppSubUrl:          setting.AppSubUrl,
		GoogleAnalyticsId:  setting.GoogleAnalyticsId,
		GoogleTagManagerId: setting.GoogleTagManagerId,
	}

	if setting.DisableGravatar {
		data.User.GravatarUrl = setting.AppSubUrl + "/public/img/transparent.png"
	}

	if len(data.User.Name) == 0 {
		data.User.Name = data.User.Login
	}

	themeUrlParam := c.Query("theme")
	if themeUrlParam == "light" {
		data.User.LightTheme = true
	}

	dashboardChildNavs := []*dtos.NavLink{
		{Text: "主页", Url: setting.AppSubUrl + "/"},
		{Text: "播放列表", Url: setting.AppSubUrl + "/playlists"},
		{Text: "快照", Url: setting.AppSubUrl + "/dashboard/snapshots"},
	}

	if c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR {
		dashboardChildNavs = append(dashboardChildNavs, &dtos.NavLink{Divider: true})
		dashboardChildNavs = append(dashboardChildNavs, &dtos.NavLink{Text: "创建", Icon: "fa fa-plus", Url: setting.AppSubUrl + "/dashboard/new"})
		dashboardChildNavs = append(dashboardChildNavs, &dtos.NavLink{Text: "导入", Icon: "fa fa-download", Url: setting.AppSubUrl + "/import/dashboard"})
	}

	data.MainNavLinks = append(data.MainNavLinks, &dtos.NavLink{
		Text:     "仪表盘",
		Icon:     "icon-gf icon-gf-dashboard",
		Url:      setting.AppSubUrl + "/",
		Children: dashboardChildNavs,
	})

	if c.OrgRole == m.ROLE_ADMIN {
		data.MainNavLinks = append(data.MainNavLinks, &dtos.NavLink{
			Text: "数据源",
			Icon: "icon-gf icon-gf-datasources",
			Url:  setting.AppSubUrl + "/datasources",
		})

		 data.MainNavLinks = append(data.MainNavLinks, &dtos.NavLink{
		 	Text: "插件",
		 	Icon: "icon-gf icon-gf-apps",
		 	Url:  setting.AppSubUrl + "/plugins",
		 })
	}

	enabledPlugins, err := plugins.GetEnabledPlugins(c.OrgId)
	if err != nil {
		return nil, err
	}

	for _, plugin := range enabledPlugins.Apps {
		if plugin.Pinned {
			appLink := &dtos.NavLink{
				Text: plugin.Name,
				Url:  plugin.DefaultNavUrl,
				Img:  plugin.Info.Logos.Small,
			}

			for _, include := range plugin.Includes {
				if !c.HasUserRole(include.Role) {
					continue
				}

				if include.Type == "page" && include.AddToNav {
					link := &dtos.NavLink{
						Url:  setting.AppSubUrl + "/plugins/" + plugin.Id + "/page/" + include.Slug,
						Text: include.Name,
					}
					appLink.Children = append(appLink.Children, link)
				}

				if include.Type == "dashboard" && include.AddToNav {
					link := &dtos.NavLink{
						Url:  setting.AppSubUrl + "/dashboard/db/" + include.Slug,
						Text: include.Name,
					}
					appLink.Children = append(appLink.Children, link)
				}
			}

			if c.OrgRole == m.ROLE_ADMIN {
				appLink.Children = append(appLink.Children, &dtos.NavLink{Divider: true})
				appLink.Children = append(appLink.Children, &dtos.NavLink{Text: "Plugin Config", Icon: "fa fa-cog", Url: setting.AppSubUrl + "/plugins/" + plugin.Id + "/edit"})
			}

			if len(appLink.Children) > 0 {
				data.MainNavLinks = append(data.MainNavLinks, appLink)
			}
		}
	}

	if c.IsGrafanaAdmin {
		data.MainNavLinks = append(data.MainNavLinks, &dtos.NavLink{
			Text: "高级选项",
			Icon: "fa fa-fw fa-cogs",
			Url:  setting.AppSubUrl + "/admin",
			Children: []*dtos.NavLink{
				{Text: "用户管理", Url: setting.AppSubUrl + "/admin/users"},
				{Text: "部门管理", Url: setting.AppSubUrl + "/admin/orgs"},
				{Text: "系统设置", Url: setting.AppSubUrl + "/admin/settings"},
				{Text: "系统状态", Url: setting.AppSubUrl + "/admin/stats"},
			},
		})
	}

	return &data, nil
}

func Index(c *middleware.Context) {
	if data, err := setIndexViewData(c); err != nil {
		c.Handle(500, "Failed to get settings", err)
		return
	} else {
		c.HTML(200, "index", data)
	}
}

func NotFoundHandler(c *middleware.Context) {
	if c.IsApiRequest() {
		c.JsonApiErr(404, "Not found", nil)
		return
	}

	if data, err := setIndexViewData(c); err != nil {
		c.Handle(500, "Failed to get settings", err)
		return
	} else {
		c.HTML(404, "index", data)
	}
}
