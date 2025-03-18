# 16 gin-blog-server 菜单模块-用户模块

## 1 菜单模块 menu

### 1.1 前端调用

在后台登陆成功的时候，会自动调用接口 http://localhost:3334/api/menu/user/list

**触发条件：**

1. 侧边栏 SideMenu 中调用了 import { usePermissionStore, useTagStore } from '@/store'

2. 然后在 permissionStore 中，根据后端传来数据构建出前端路由

   ```javascript
   export const usePermissionStore = defineStore('permission', {
     persist: {
       key: 'gvb_admin_permission',
     },
     state: () => ({
       accessRoutes: [], // 可访问的路由
     }),
     getters: {
       // 最终可访问路由 = 基础路由 + 可访问的路由
       routes: state => basicRoutes.concat(state.accessRoutes),
       // 过滤掉 hidden 的路由作为左侧菜单显示
       menus: state => state.routes.filter(route => route.name && !route.isHidden),
     },
     actions: {
       // ! 后端生成路由: 后端返回的就是最终路由, 处理成前端格式
       async generateRoutesBack() {
         const resp = await api.getUserMenus() // 调用接口获取后端传来的路由
         this.accessRoutes = buildRoutes(resp.data) // 处理成前端路由格式
         return this.accessRoutes
       },
       ...
     },
   })
   ```

3. api.getUserMenus 会调用 getUserMenus: () => request.get('/menu/user/list'), // 获取当前用户的菜单，来获取信息构建菜单栏目

4. **构建好菜单栏目之后，可以根据菜单栏构建对应的侧边栏，并构建相关的路由负责后台界面的跳转**

**后台管理的 api.js 文件中，菜单模块的相关接口如下：**

```javascript
 // 权限管理相关接口
  // 菜单
  getUserMenus: () => request.get('/menu/user/list'), // 获取当前用户的菜单
  getMenus: (params = {}) => request.get('/menu/list', { params }),
  saveOrUpdateMenu: data => request.post('/menu', data),
  deleteMenu: id => request.delete(`/menu/${id}`),
  getMenuOption: () => request.get('/menu/option'),
```

**后端中的菜单模块相关接口如下：**

```go
// 菜单模块
	menu := auth.Group("/menu")
	{
		menu.GET("/list", menuAPI.GetTreeList)      // 菜单列表
		menu.POST("", menuAPI.SaveOrUpdate)         // 新增/编辑菜单
		menu.DELETE("/:id", menuAPI.Delete)         // 删除菜单
		menu.GET("/user/list", menuAPI.GetUserMenu) // 获取当前用户的菜单
		menu.GET("/option", menuAPI.GetOption)      // 菜单选项列表(树形)
	}
```



### 1.2 获取当前用户的菜单 /menu/user/list

manager.go 中添加接口：

```go
// 菜单模块
	menu := auth.Group("/menu")
	{
		menu.GET("/user/list", menuAPI.GetUserMenu) // 获取当前用户的菜单
	}
```

之后我们添加一个对应的 handle_menu.go 来处理 menu 相关的接口

```go
package handle

import (
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/model"
	"github.com/gin-gonic/gin"
	"sort"
)

type Menu struct{}

// MenuTreeVO Menu的树形结构，主要添加的属性为 children 列表
type MenuTreeVO struct {
	model.Menu
	Children []MenuTreeVO `json:"children"`
}

// GetUserMenu 获取当前用户菜单: 生成后台管理界面的菜单
func (*Menu) GetUserMenu(c *gin.Context) {
	db := GetDB(c)
	auth, _ := CurrentUserAuth(c)

	var menus []model.Menu
	var err error

	if auth.IsSuper { // 如果当前用户是超级管理员
		menus, err = model.GetAllMenuList(db)
	} else {
		menus, err = model.GetMenuListByUserId(GetDB(c), auth.ID)
	}

	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, menus2MenuVos(menus))
}

// 构建菜单列表的树形数据结构, []Menu => []MenuVo
func menus2MenuVos(menus []model.Menu) []MenuTreeVO {
	result := make([]MenuTreeVO, 0)

	// 筛选出一级菜单 (parentId == 0 的菜单)
	firstLevelMenus := getFirstLevelMenus(menus)
	// key 是菜单 ID, value 是该菜单对应的子菜单列表
	childrenMap := getMenuChildrenMap(menus)

	for _, first := range firstLevelMenus {
		menu := MenuTreeVO{Menu: first}
		for _, childMenu := range childrenMap[first.ID] {
			menu.Children = append(menu.Children, MenuTreeVO{Menu: childMenu})
		}
		delete(childrenMap, first.ID)
		result = append(result, menu)
	}

	sortMenu(result)
	return result
}

// 筛选出一级菜单 (parentId == 0 的菜单)
func getFirstLevelMenus(menuList []model.Menu) []model.Menu {
	firstLevelMenus := make([]model.Menu, 0)
	for _, menu := range menuList {
		if menu.ParentId == 0 {
			firstLevelMenus = append(firstLevelMenus, menu)
		}
	}

	return firstLevelMenus
}

// key 是菜单 ID, value 是该菜单对应的子菜单列表
func getMenuChildrenMap(menuList []model.Menu) map[int][]model.Menu {
	childrenMap := make(map[int][]model.Menu)
	for _, menu := range menuList {
		if menu.ParentId != 0 {
			childrenMap[menu.ParentId] = append(childrenMap[menu.ParentId], menu)
		}
	}
	return childrenMap
}

// 以 orderNum 升序排序，包括子菜单
func sortMenu(menus []MenuTreeVO) {
	// 如果 menus[i].OrderNum 小于 menus[j].OrderNum，则返回 true，表示 menus[i] 应该排在 menus[j] 前面；
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].OrderNum < menus[j].OrderNum
	})

	for i := range menus {
		sort.Slice(menus[i].Children, func(j, k int) bool {
			return menus[i].Children[j].OrderNum < menus[i].Children[k].OrderNum
		})
	}
}

```

**以下是上述代码中各个函数的功能说明**

1. **`GetUserMenu` 函数**
   - **功能**：处理获取当前用户菜单的请求
   - 步骤：
     - 根据用户是否为超级管理员，从数据库获取所有菜单或特定用户的菜单
     - 调用 `menus2MenuVos` 函数将菜单列表转换为树形结构
     - 返回树形结构的菜单数据
2. **`menus2MenuVos` 函数**
   - **功能**：将扁平的菜单列表转换为树形结构
   - 步骤：
     - 调用 `getFirstLevelMenus` 获取一级菜单（`parentId == 0`）
     - 调用 `getMenuChildrenMap` 生成子菜单映射表
     - 遍历一级菜单，组装树形结构并递归添加子菜单
     - 调用 `sortMenu` 对菜单进行排序
3. **`getFirstLevelMenus` 函数**
   - **功能**：筛选出一级菜单
   - 步骤：
     - 遍历所有菜单，收集 `parentId == 0` 的菜单
     - 返回一级菜单列表
4. **`getMenuChildrenMap` 函数**
   - **功能**：生成子菜单映射表
   - 步骤：
     - 遍历所有菜单，将子菜单（`parentId != 0`）按父菜单 ID 分组
     - 返回键为父菜单 ID、值为子菜单列表的映射表
5. **`sortMenu` 函数**
   - **功能**：对菜单进行排序
   - 步骤：
     - 按 `orderNum` 升序排列一级菜单
     - 递归对每个菜单的子菜单按 `orderNum` 升序排列

**整体流程**：

1. 获取原始菜单数据 → 2. 转换为树形结构 → 3. 筛选一级菜单 → 4. 构建子菜单映射 → 5. 组装树形结构 → 6. 排序菜单

**之后，我们尝试登陆后台首页，会看到其登陆成功之后发送如下的请求：**

![image-20250318111313272](./assets/image-20250318111313272.png)

然而，其返回的响应为空

![image-20250318111340051](./assets/image-20250318111340051.png)

这就导致后端界面不能获取到对应的菜单列表，从而无法获取到有用的菜单信息，也就不能正确的将首页路由 / 注册到 router 中

我们可以看原项目数据库中，对应的 menu 信息如下：

+ 通过获取属性 path 和 redirect，后台可以根据返回的路径自动组装 router，实现不同用户可以访问的后台视图是不一样的
+ GetMenuListByUserId 方法中，根据获取用户 userAuth 之后，获取其 role 角色 -> 然后获取 role 角色拥有的菜单 menu

![image-20250318111603989](./assets/image-20250318111603989.png)

因此，到这里，我们要补全 role 和 menu 的一些数据库对应关系

+ role 即当前项目可能存在的角色：超级管理员？普通用户？
+ role 和菜单 menu 是多对多的关系，且关系基本不会变，可以将原项目中的信息直接导入到新项目库中
+ 当前从前台注册的用户统一都是普通用户

我们将原数据库中的 role 表、menu 表、role_menu表进行补全，**对应的 sql 文件在文件夹 sql 下的：menu.sql / role.sql / role_menu.sql 中**

同时，由于中间件 Permission Check 调用 CheckRoleAuth 函数根据 role id 判断是否有运行 uri+method 的权限，因此我们需要对如下三个表进行补全，导入的方式和上面类似：

+ resource 表
+ role_resource 表

然后，刷新后台界面，查看响应如下：

<img src="./assets/image-20250318125517484.png" alt="image-20250318125517484" style="zoom:67%;" />

然后我们可以正常进入后台界面，展示如下：

<img src="./assets/image-20250318130037375.png" alt="image-20250318130037375" style="zoom:67%;" />

至此，由于我们导入了 resource 数据库中的数据，我们可以将之前 middleware/auth.go 中注释掉的代码进行还原。还原后功能也可以正常运行



### 1.3 菜单列表 /menu/list





### 1.4 新增菜单 /menu



### 1.5 删除菜单 /menu/:id



### 1.6 菜单选项列表 /option





## 2 用户模块 user