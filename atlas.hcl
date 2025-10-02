data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlas-loader",
  ]
}

env "local" {
  src = data.external_schema.gorm.url
  url = "mysql://nebula:nebula@127.0.0.1:3307/nebula"
  dev = "docker://mysql/8/dev"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
  # 安全配置 - 防止意外删除
  diff {
    skip {
      drop_schema = true
      drop_table  = true
    }
  }
}

env "dev" {
  src = data.external_schema.gorm.url
  dev = "docker://mysql/8/dev"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "prod" {
  src = data.external_schema.gorm.url
  url = env("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
  # 生产环境安全配置
  diff {
    skip {
      drop_schema = true
      drop_table  = true
      drop_column = true
    }
  }
}

# 全局 lint 规则
lint {
  # 命名规范
  naming {
    match   = "^[a-z]+(_[a-z]+)*$"
    message = "table and column names should be snake_case"
  }
  
  # 破坏性操作检查
  destructive {
    error = true
  }
  
  # 数据依赖检查
  data_depend {
    error = true
  }
}
