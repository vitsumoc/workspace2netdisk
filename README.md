# workspace2netdisk

工作目录备份到网盘的工具

运行时指定工作目录根路径，运行后会生成两个 md 文件：delete.md 和 folder.md

delete.md 用来提示上传网盘前需要删除的文件位置，例如 node_modules 文件夹、.svn 文件夹、内容过大的文件等

folder.md 用 md 格式表示文件夹下所有内容的路径，这样就可以打包上传单文件，以后搜索 md 文件找内容就行

使用 autoDelete 标志可以自动把不需要的文件夹删除掉

用法大概就是这样：

```go
go run . C:\Users\vc\Desktop\2022msj -autoDelete=true
```