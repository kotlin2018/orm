# 核心结构简介
与GoFrame/orm兼容，主要是修改了gf/database/gdb目录

[github地址](https://github.com/kotlin2018/orm.git)
  
* `DB interface`是直接操作数据库的接口，里面定义了一系列操作数据库的方法。
````go
type DB interface {
        Insert(table string, data interface{}, batch ...int) (sql.Result, error)
        InsertIgnore(table string, data interface{}, batch ...int) (sql.Result, error)
        Replace(table string, data interface{}, batch ...int) (sql.Result, error)
        Save(table string, data interface{}, batch ...int) (sql.Result, error)
        ...
} 
````
* `Core`是数据库管理的基础结构。只实现了 DB接口的一部分方法，是所有SQL驱动的基类(父类)。
 
  各种数据库驱动要想操作数据，必须完全实现DB接口，但这必定会造成代码冗余。
  
  例如: 
  
  Mysql驱动实现了 Insert()、InsertIgnore()、Replace()、Save()...
  
  Pgsql也实现了 Insert()、InsertIgnore()、Replace()、Save()...
  
  这部分公共的方法完全可以抽取出来，只让Core结构体实现，所有的数据库驱动只要继承(内嵌)Core结构体即可；
  
  这样各种数据库驱动就间接的实现了 Insert()、InsertIgnore()、Replace()、Save()...等放法。

````go
type Core struct {
	DB     DB              // DB 接口对象。持有DB实例，使Model结构体具备操作数据库的能力。
    group  string          // 配置组的组名。
    debug  *gtype.Bool     // 为数据库启用debug模式，该模式可以在运行时更改。
    cache  *gcache.Cache   // 缓存管理器，仅缓存SQL结果。
    schema *gtype.String   // 当前对象对象的自定义架构，可以方便的在运行时切换数据库。
    logger *glog.Logger    // Logger.
    config *ConfigNode     // 当前配置节点
    ctx    context.Context // 仅链接操作的上下文。
}  
                 
````

* `Model`是`orm`的 `dao` `(Data Access Object)`数据库访问对象

````go
type Model struct {
    db    DB       // 底层数据库接口。持有这个实例就具备了操作数据的能力。
    tx    *TX      // 底层事务管理接口。   
    ...
}
````
