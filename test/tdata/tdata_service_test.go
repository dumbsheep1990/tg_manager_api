package tdata_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tg_manager_api/services/tdata/service"
)

// 模拟TData数据库接口
type MockTDataRepository struct {
	mock.Mock
}

func (m *MockTDataRepository) Create(tdata interface{}) error {
	args := m.Called(tdata)
	return args.Error(0)
}

func (m *MockTDataRepository) Update(tdata interface{}) error {
	args := m.Called(tdata)
	return args.Error(0)
}

func (m *MockTDataRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTDataRepository) FindByID(id uint) (interface{}, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func (m *MockTDataRepository) FindAll(query interface{}) ([]interface{}, int64, error) {
	args := m.Called(query)
	return args.Get(0).([]interface{}), int64(args.Int(1)), args.Error(2)
}

// 模拟文件处理接口
type MockFileHandler struct {
	mock.Mock
}

func (m *MockFileHandler) SaveFile(data []byte, filename string) (string, error) {
	args := m.Called(data, filename)
	return args.String(0), args.Error(1)
}

func (m *MockFileHandler) DeleteFile(filepath string) error {
	args := m.Called(filepath)
	return args.Error(0)
}

func (m *MockFileHandler) ExtractTData(filepath string) (string, error) {
	args := m.Called(filepath)
	return args.String(0), args.Error(1)
}

// 测试上传TData文件
func TestUploadTDataFile(t *testing.T) {
	mockRepo := new(MockTDataRepository)
	mockFileHandler := new(MockFileHandler)
	tdataService := service.NewTDataService()
	
	// 模拟文件数据和元数据
	fileData := []byte("file data content")
	fileInfo := map[string]interface{}{
		"filename": "test.zip",
		"size": int64(len(fileData)),
		"content_type": "application/zip",
	}
	
	// 设置模拟行为
	mockFileHandler.On("SaveFile", fileData, "test.zip").Return("/path/to/test.zip", nil)
	mockRepo.On("Create", mock.Anything).Return(nil)
	
	// 测试上传文件
	result, err := tdataService.UploadTDataFile(fileData, fileInfo)
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockFileHandler.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// 测试获取TData文件列表
func TestGetTDataFiles(t *testing.T) {
	mockRepo := new(MockTDataRepository)
	tdataService := service.NewTDataService()
	
	// 模拟查询参数
	query := map[string]interface{}{
		"page": 1,
		"page_size": 10,
		"filename": "test",
		"import_status": "notImported",
	}
	
	// 模拟返回数据
	mockTDatas := []interface{}{
		map[string]interface{}{
			"id": uint(1),
			"filename": "test1.zip",
			"file_path": "/path/to/test1.zip",
			"file_size": int64(1024),
			"import_status": "notImported",
			"uploaded_at": time.Now(),
		},
		map[string]interface{}{
			"id": uint(2),
			"filename": "test2.zip",
			"file_path": "/path/to/test2.zip",
			"file_size": int64(2048),
			"import_status": "notImported",
			"uploaded_at": time.Now(),
		},
	}
	
	// 设置模拟行为
	mockRepo.On("FindAll", mock.Anything).Return(mockTDatas, 2, nil)
	
	// 测试获取文件列表
	files, total, err := tdataService.GetTDataFiles(query)
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 2, len(files))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

// 测试导入TData文件
func TestImportTDataFile(t *testing.T) {
	mockRepo := new(MockTDataRepository)
	mockFileHandler := new(MockFileHandler)
	tdataService := service.NewTDataService()
	
	// 模拟TData数据
	mockTData := map[string]interface{}{
		"id": uint(1),
		"filename": "test.zip",
		"file_path": "/path/to/test.zip",
		"file_size": int64(1024),
		"import_status": "notImported",
	}
	
	importParams := map[string]interface{}{
		"group_id": uint(1),
		"remark": "测试导入",
	}
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(mockTData, nil)
	mockFileHandler.On("ExtractTData", "/path/to/test.zip").Return("/path/to/extracted", nil)
	mockRepo.On("Update", mock.Anything).Return(nil)
	
	// 测试导入TData
	err := tdataService.ImportTDataFile(uint(1), importParams)
	
	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockFileHandler.AssertExpectations(t)
}

// 测试删除TData文件
func TestDeleteTDataFile(t *testing.T) {
	mockRepo := new(MockTDataRepository)
	mockFileHandler := new(MockFileHandler)
	tdataService := service.NewTDataService()
	
	// 模拟TData数据
	mockTData := map[string]interface{}{
		"id": uint(1),
		"filename": "test.zip",
		"file_path": "/path/to/test.zip",
		"file_size": int64(1024),
	}
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(mockTData, nil)
	mockFileHandler.On("DeleteFile", "/path/to/test.zip").Return(nil)
	mockRepo.On("Delete", uint(1)).Return(nil)
	
	// 测试删除TData
	err := tdataService.DeleteTDataFile(uint(1))
	
	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockFileHandler.AssertExpectations(t)
}

// 测试获取TData导入状态
func TestGetTDataImportStatus(t *testing.T) {
	mockRepo := new(MockTDataRepository)
	tdataService := service.NewTDataService()
	
	// 模拟TData数据
	mockTData := map[string]interface{}{
		"id": uint(1),
		"filename": "test.zip",
		"import_status": "importing",
		"import_progress": 45,
	}
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(mockTData, nil)
	
	// 测试获取导入状态
	status, err := tdataService.GetTDataImportStatus(uint(1))
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "importing", status["import_status"])
	assert.Equal(t, 45, status["import_progress"])
	mockRepo.AssertExpectations(t)
}
