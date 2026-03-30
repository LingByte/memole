package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

// 全局配置存储
var configMap = make(map[string]string)

// 加载配置文件
func LoadConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("加载 .env 文件失败: %v", err)
	}
	
	// 将环境变量加载到内存中
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			configMap[pair[0]] = pair[1]
		}
	}
	
	return nil
}

// 从指定文件加载配置
func LoadConfigFromFile(filename string) error {
	err := godotenv.Load(filename)
	if err != nil {
		return fmt.Errorf("加载配置文件 %s 失败: %v", filename, err)
	}
	
	// 将环境变量加载到内存中
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			configMap[pair[0]] = pair[1]
		}
	}
	
	return nil
}

// 获取配置项的字符串值
func Get(key string) string {
	// 优先从内存中获取
	if val, exists := configMap[key]; exists {
		return val
	}
	// 从环境变量获取
	return os.Getenv(key)
}

// 设置配置项
func Set(key, value string) {
	configMap[key] = value
	os.Setenv(key, value)
}

// 获取配置项的整数值
func GetInt(key string) (int, error) {
	val := Get(key)
	if val == "" {
		return 0, fmt.Errorf("配置项 %s 未设置", key)
	}
	result, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("配置项 %s 解析失败: %v", key, err)
	}
	return result, nil
}

// 设置整型配置项
func SetInt(key string, value int) {
	Set(key, strconv.Itoa(value))
}

// 获取配置项的浮点数值
func GetFloat(key string) (float64, error) {
	val := Get(key)
	if val == "" {
		return 0, fmt.Errorf("配置项 %s 未设置", key)
	}
	result, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, fmt.Errorf("配置项 %s 解析失败: %v", key, err)
	}
	return result, nil
}

// 设置浮点型配置项
func SetFloat(key string, value float64) {
	Set(key, strconv.FormatFloat(value, 'f', -1, 64))
}

// 获取配置项的布尔值
func GetBool(key string) (bool, error) {
	val := Get(key)
	if val == "" {
		return false, fmt.Errorf("配置项 %s 未设置", key)
	}
	result, err := strconv.ParseBool(val)
	if err != nil {
		return false, fmt.Errorf("配置项 %s 解析失败: %v", key, err)
	}
	return result, nil
}

// 设置布尔型配置项
func SetBool(key string, value bool) {
	Set(key, strconv.FormatBool(value))
}

// 检查配置项是否存在
func Has(key string) bool {
	val := Get(key)
	return val != ""
}

// 删除配置项
func Delete(key string) {
	delete(configMap, key)
	os.Unsetenv(key)
}

// 获取所有配置项
func GetAll() map[string]string {
	result := make(map[string]string)
	for k, v := range configMap {
		result[k] = v
	}
	return result
}

// 保存配置到文件
func SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建配置文件失败: %v", err)
	}
	defer file.Close()

	for key, value := range configMap {
		_, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return fmt.Errorf("写入配置文件失败: %v", err)
		}
	}

	return nil
}

// 清空所有配置
func Clear() {
	configMap = make(map[string]string)
}
