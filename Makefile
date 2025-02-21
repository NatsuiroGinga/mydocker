# 定义编译输出目录和可执行文件名
BIN_DIR := bin
APP := mydocker
TARGET := $(BIN_DIR)/$(APP)

# 收集所有 Go 源文件（主文件和依赖包）
GO_SOURCES := $(shell find . -type f -name "*.go")

# 默认目标：编译整个项目
all: clean uninstall install

# 编译主程序并输出到 bin/app
$(TARGET): $(GO_SOURCES)
	@mkdir -p $(BIN_DIR)
	go build -o $@ .

# 清理编译生成的文件
clean:
	rm -rf $(BIN_DIR)/$(APP)

# 安装到系统目录（需要 sudo 权限）
install: $(TARGET)
	cp $(TARGET) /usr/local/bin
	cp $(TARGET) /usr/bin/

# 卸载（从所有可能的安装路径删除）
uninstall:
	rm -f /usr/local/bin/$(APP) /usr/bin/$(APP)