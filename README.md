## sing-box 简易托盘应用

![](https://thumbsnap.com/i/cxhVn14K.png)

webui 订阅管理，支持添加多个订阅，切换订阅，自动更新订阅

基于singbox核心运行，可自由更换

采用go开发，体积小，占用少

不往注册表拉屎，但为了显示 windows 通知，需要注册 APPID 以显示通知图标和标题

卸载后，可自行清理：HKEY_CURRENT_USER\Software\Classes\AppUserModelId\cat-box

不往用户家目录拉屎，所有运行数据都保存在程序当前目录

无需安装，不乱拉屎，真绿色运行



### 快速上手

![](https://thumbsnap.com/i/dxj7G6J4.jpg)

### 模板修改

模板路径：resources\template\template.json

需要修改模板的地方基本上是策略组了，通过自定义不同的策略组进行分流

### 订阅

数据库路径：resources\db\app.db

自动更新：

支持以m（分钟）或h（小时）为单位自定义更新延时

例如：1m，30m，1h，24h，48h

订阅代理：

当开启代理时，订阅会自动使用代理，未开启则直连

### 代理模式

代理模式：系统代理和TUN模式

TUN模式需要管理员权限，开启前需要求程序以管理员模式运行

### 规则文件

本地规则：模板内的本地规则应放在此目录 `resources\geo`，不要瞎几吧乱放，我有强迫症

在线规则：可以自动更新，详见sing-box文档

### 托盘图标

路径：resources\icons

可自定义，但图片格式必须是.ico，且不可重命名

### 开机自启

手动修改注册表方式实现

注册表路径：HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run

进入 Run 项目，右键新建“字符串值”，右键新建的数值，点击“修改”：

数值名称：cat-box

数值数据："cat-box.exe路径" --enable-workspace

数值数据例子："E:\cat-box\cat-box.exe" --enable-workspace

参数说明：--enable-workspace 参数表示开启工作目录，否则不可以通过绝对路径运行
