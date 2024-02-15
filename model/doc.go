package model

/*

1.覆盖默认的 ID:
  - 给新的 primaryKey 字段增加 gorm tag "gorm:primaryKey"
  - 覆盖默认的 ID 并设置 json 和 gorm tag 为 "-"
    ID string `json:"-" gorm:"-"`
  - 覆盖默认的 SetID() 方法
    func(* FeishuUser)SetID(...string){}
2.model 结构体对象的字段如果要通过 query parameter 作为查询参数的话, 需要增加 schema tag
3.自己定义的 model 继承 Base model 时必须时匿名继承,并且不能加 json tag
  否则在 gin ShouldBindJson 和 json.Unmarshal 时会出问题
4.如果自定义了 ID 的规则, 那么记住几点:
  1.一定不要在自定义 Model 中增加 ID 字段,
  2.不要重写 GetD() 和 SetID() 方法
5.如果要 UpdatePartial, 修改的字段如果是基本类型,比如 int, string 等, 如果修改的值是默认值(zero value),
  那么必须该类型改成指针类型,否则无法修改.

5.外键只允许多集关系的表,并且这种表的 children 也是自己

rbac
	g hybfkuf admin                  // hybfkuf 属于 admin 组
	g user1 admin                    // user1 属于 admin 组
	p admin /api/asset/asset GET     // admin 组允许 GET
	p admin /api/asset/asset POST    // admin 组允许 POST
	p admin /api/asset/asset PATCH   // admin 组允许 PATCH
	p admin /api/asset/asset DELETE  // admin 组允许 DELETE

	Group -> Policy: 给角色组设置策略(group 就是 role)
	Group <- User: 向角色组中添加和删除用户

*/
