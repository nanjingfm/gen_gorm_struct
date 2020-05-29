# gen_gorm_struct
生成gorm结构体

如果你使用gorm并且你也在使用mysql，而且你也需要根据mysql表生成结构体，那么就可以使用整个工具。
可以选择增加form和json tag

## 用法
```shell script
./cmd -h
Usage of ./cmd:
  -P int
        db port (default 3306)
  -d string
        db name (default "test")
  -form
        with form tag
  -h string
        db host (default "127.0.0.1")
  -json
        with json tag
  -p string
        db password (default "123456")
  -t string
        tableName name
  -u string
        db user (default "root")
```
