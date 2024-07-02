## sing-box 简易托盘应用

![](https://thumbsnap.com/i/cxhVn14K.png)

webui 订阅管理，支持添加多个订阅，切换订阅，自动更新订阅

基于singbox核心运行，可自由更换

采用go开发，体积小，占用少

不往注册表拉屎，但为了显示 windows 通知，需要注册 APPID 以显示通知图标和标题

卸载后，可自行清理：HKEY_CURRENT_USER\Software\Classes\AppUserModelId\cat-box

不往用户家目录拉屎，所有运行数据都保存在程序当前目录

无需安装，不乱拉屎，真绿色运行

### 使用教程

#### 快速上手

![](https://thumbsnap.com/i/dxj7G6J4.jpg)

#### 模板修改

模板路径：resources\template\template.json

注意：

不支持自定义策略组

不可编辑 "experimental":{...} 配置，以免出现不必要的错误

mixed 入站可编辑，但不可删掉配置项 set_system_proxy

tun 入站可编辑

#### yacd webui

切换节点立即生效：点击代理页面右上角设置，开启”切换代理时自动断开旧连接“

#### sub-box 订阅管理

账户注册：

初次打开，需要注册账号，作为网页，需要保证访问安全

账号保存在本地sqlite数据库内，无法找回，切记不要忘记，如果忘记，删除数据库重新注册

数据库路径：resources\db\app.db



自动更新延时：

支持以m（分钟）或h（小时）为单位自定义更新延时

例如：1m，30m，1h，24h，48h

#### 代理模式

代理模式：系统代理，TUN模式

TUN模式需要管理员权限，开启前需要求程序以管理员模式运行

勾选指定模式后，左键单击托盘图标开启代理，已开启代理，需要左键单击停止再开启以生效所选代理模式

#### GEO文件更新

路径：resources\geo

当前不支持自动更新，需手动下载替换，且文件命名不可更改，如需更改，需保证模板内同步更改

#### 托盘图标

路径：resources\icons

可自定义，但图片格式必须是.ico，且不可重命名
