# docker upx插件

## 参数

```yaml
  - name: 构建开发镜像
    image: ysicing/drone-upx
    pull: always
    privileged: true
    settings:
      level: 9 # default 
      path: ./dist # 必填参数
      include: linux # 多文件文件名包含
      exclude: windows # 多文件文件名排除
```

> 如果为文件夹时，最好设置include或exclude
