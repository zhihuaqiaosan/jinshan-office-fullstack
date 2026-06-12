package initialize

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/glebarez/sqlite"
)

// GormSqlite 初始化Sqlite数据库
func GormSqlite() *gorm.DB {
	s := global.GVA_CONFIG.Sqlite
	if s.Dbname == "" {
		global.GVA_LOG.Error("请配置Sqlite数据库名称")
		return nil
	}
	return initSqliteDatabase(s)
}

func initSqliteDatabase(s config.Sqlite) *gorm.DB {
	// sqlite配置
	dsn := s.Path + "/" + s.Dbname + "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=cache_size(2000)"
	global.GVA_LOG.Info("Sqlite数据库连接信息: " + dsn)

	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: false,
	}
	if s.LogMode == "info" {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	if db, err := gorm.Open(sqlite.Open(dsn), config); err != nil {
		global.GVA_LOG.Error("Sqlite数据库启动异常!")
		global.GVA_LOG.Error(err.Error())
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(s.MaxIdleConns)
		sqlDB.SetMaxOpenConns(s.MaxOpenConns)
		return db
	}
}
