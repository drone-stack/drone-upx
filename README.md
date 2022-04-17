# docker upx插件

## 参数

```yaml
  - name: 构建开发镜像
    image: ysicing/drone-plugin-upx
    pull: always
    privileged: true
    settings:
      registry: ccr.ccs.tencentyun.com
      repo: ccr.ccs.tencentyun.com/ysicing/drone-plugin-builder
      debug: true
      mode: dev
      tags: develop
      purge: false
      no_cache: false
      dockerfile: Dockerfile

```