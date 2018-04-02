package node

import (
	"commons/errors"
	"commons/results"
	appDbmocks "db/mongo/app/mocks"
	nodeDbmocks "db/mongo/node/mocks"
	groupDbmocks "db/mongo/group/mocks"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

const (
	appId1       = "000000000000000000000000"
	appId2       = "111111111111111111111111"
	nodeId1      = "000000000000000000000001"
	nodeId2      = "000000000000000000000002"
	nodeId3      = "000000000000000000000003"
	groupId1     = "000000000000000000000011"
	host         = "192.168.0.1"
	port         = "8888"
	status       = "connected"
	groupId1Name = "testGroup"
)

var (
	appId1Images = []string{commonImage, "etc1"}
	appId2Images = []string{commonImage, "etc2"}
	commonImage  = "testimage1"
	services     = []map[string]interface{}{
		{
			"name": "test",
			"state": map[string]interface{}{
				"Status":   "teststatus",
				"ExitCode": "testexitcode",
			},
		},
	}

	group1 = map[string]interface{}{
		"id":      groupId1,
		"name":    groupId1Name,
		"members": []string{nodeId1, nodeId2},
	}

	node1 = map[string]interface{}{
		"apps":   []string{appId1},
		"id":     nodeId1,
		"ip":     host,
		"status": status,
	}

	node2 = map[string]interface{}{
		"apps":   []string{appId2},
		"id":     nodeId2,
		"ip":     host,
		"status": status,
	}

	node3 = map[string]interface{}{
		"apps":   []string{appId2},
		"id":     nodeId3,
		"ip":     host,
		"status": status,
	}

	allQuery = map[string][]string{
		GROUP_ID:   []string{groupId1},
		NODE_ID:    []string{nodeId1},
		APP_ID:     []string{appId1},
		IMAGE_NAME: []string{commonImage},
	}

	app1 = map[string]interface{}{
		"id":       appId1,
		"images":   appId1Images,
		"services": services,
	}

	app2 = map[string]interface{}{
		"id":       appId2,
		"images":   appId2Images,
		"services": services,
	}

	nodes = []map[string]interface{}{
		node1,
		node2,
	}

	groups = []map[string]interface{}{
		group1,
	}
)

var searchExecutor Command

func init() {
	searchExecutor = Executor{}
}

func TestSearchNodesWithAllQuery_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appDbExecutorMockObj := appDbmocks.NewMockCommand(ctrl)
	groupDbExecutorMockObj := groupDbmocks.NewMockCommand(ctrl)
	nodeDbExecutorMockObj := nodeDbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodeDbExecutorMockObj.EXPECT().GetNodes().Return(nodes, nil),
		groupDbExecutorMockObj.EXPECT().GetGroups().Return(groups, nil),
		appDbExecutorMockObj.EXPECT().GetApp(appId1).Return(app1, nil),
	)

	// pass mockObj to a real object.
	appDbExecutor = appDbExecutorMockObj
	groupDbExecutor = groupDbExecutorMockObj
	nodeDbExecutor = nodeDbExecutorMockObj

	code, res, err := searchExecutor.SearchNodes(allQuery)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}

	expectedResult := make(map[string]interface{})
	expectedResult["nodes"] = make([]map[string]interface{}, 1)
	expectedResult["nodes"].([]map[string]interface{})[0] = node1

	if !reflect.DeepEqual(expectedResult, res) {
		t.Errorf("Expected res: %s\n actual res: %s", expectedResult, res)
	}
}

func TestSearchNodesWithAllQueryWhenGetNodesFailed_ExpectRetrunError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeDbExecutorMockObj := nodeDbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodeDbExecutorMockObj.EXPECT().GetNodes().Return(nil, errors.DBOperationError{}),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodeDbExecutorMockObj

	code, _, err := searchExecutor.SearchNodes(allQuery)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBOperationError", "nil")

	}

	if code != results.ERROR {
		t.Errorf("Expected return code : %d, actual err: %d", 500, code)
	}
}

func TestSearchNodesWithAllQueryWhenGetGroupsFailed_ExpectRetrunError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeDbExecutorMockObj := nodeDbmocks.NewMockCommand(ctrl)
	groupDbExecutorMockObj := groupDbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodeDbExecutorMockObj.EXPECT().GetNodes().Return(nodes, nil),
		groupDbExecutorMockObj.EXPECT().GetGroups().Return(nil, errors.DBOperationError{}),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodeDbExecutorMockObj
	groupDbExecutor = groupDbExecutorMockObj

	code, _, err := searchExecutor.SearchNodes(allQuery)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBOperationError", "nil")

	}

	if code != results.ERROR {
		t.Errorf("Expected return code : %d, actual err: %d", 500, code)
	}
}

func TestSearchNodesWithAllQueryWhenGetAppFailed_ExpectRetrunError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeDbExecutorMockObj := nodeDbmocks.NewMockCommand(ctrl)
	groupDbExecutorMockObj := groupDbmocks.NewMockCommand(ctrl)
	appDbExecutorMockObj := appDbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		nodeDbExecutorMockObj.EXPECT().GetNodes().Return(nodes, nil),
		groupDbExecutorMockObj.EXPECT().GetGroups().Return(groups, nil),
		appDbExecutorMockObj.EXPECT().GetApp(gomock.Any()).Return(nil, errors.DBOperationError{}),
	)

	// pass mockObj to a real object.
	nodeDbExecutor = nodeDbExecutorMockObj
	groupDbExecutor = groupDbExecutorMockObj
	appDbExecutor = appDbExecutorMockObj

	code, _, err := searchExecutor.SearchNodes(allQuery)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBOperationError", "nil")

	}

	if code != results.ERROR {
		t.Errorf("Expected return code : %d, actual err: %d", 500, code)
	}
}

func TestFilterByGroupId_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupDbmocks.NewMockCommand(ctrl)
	
	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroups().Return(groups, nil),
	)

	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	nodes := make([]map[string]interface{}, 3)
	nodes[0] = node1
	nodes[1] = node2
	nodes[2] = node3

	res, err := filterByGroupId(nodes, groupId1)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	expectedResult := make([]map[string]interface{}, 2)
	expectedResult[0] = node1
	expectedResult[1] = node2

	if !reflect.DeepEqual(expectedResult, res) {
		t.Errorf("Expected res: %s\n actual res: %s", expectedResult, res)
	}
}

func TestFilterByGroupIdWhenGetGroupsFailed_ExpectReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupDbExecutorMockObj := groupDbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		groupDbExecutorMockObj.EXPECT().GetGroups().Return(nil, errors.DBOperationError{}),
	)

	// pass mockObj to a real object.
	groupDbExecutor = groupDbExecutorMockObj

	_, err := filterByGroupId(nodes, groupId1)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBOperationError", "nil")
	}
}

func TestFilterByAppId_CheckReturnCorrect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodes := make([]map[string]interface{}, 3)
	nodes[0] = node1
	nodes[1] = node2
	nodes[2] = node3

	res := filterByAppId(nodes, appId1)

	expectedResult := make([]map[string]interface{}, 1)
	expectedResult[0] = node1

	if !reflect.DeepEqual(expectedResult, res) {
		t.Errorf("Expected res: %s\n actual res: %s", expectedResult, res)
	}
}

func TestFilterByImageName_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appDbExecutorMockObj := appDbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		appDbExecutorMockObj.EXPECT().GetApp(gomock.Any()).Return(app1, nil),
		appDbExecutorMockObj.EXPECT().GetApp(gomock.Any()).Return(app2, nil),
		appDbExecutorMockObj.EXPECT().GetApp(gomock.Any()).Return(app2, nil),
	)

	// pass mockObj to a real object.
	appDbExecutor = appDbExecutorMockObj

	nodes := make([]map[string]interface{}, 3)
	nodes[0] = node1
	nodes[1] = node2
	nodes[2] = node3

	res, err := filterByImageName(nodes, commonImage)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())

	}
	expectedResult := make([]map[string]interface{}, 3)
	expectedResult[0] = node1
	expectedResult[1] = node2
	expectedResult[2] = node3

	if !reflect.DeepEqual(expectedResult, res) {
		t.Errorf("Expected res: %s\n actual res: %s", expectedResult, res)
	}
}

func TestFilterByImageNameWhenGetAppFailed_ExpectReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appDbExecutorMockObj := appDbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		appDbExecutorMockObj.EXPECT().GetApp(gomock.Any()).Return(nil, errors.DBOperationError{}),
	)

	// pass mockObj to a real object.
	appDbExecutor = appDbExecutorMockObj

	nodes := make([]map[string]interface{}, 3)
	nodes[0] = node1
	nodes[1] = node2
	nodes[2] = node3

	_, err := filterByImageName(nodes, commonImage)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "DBOperationError", "nil")

	}
}

func TestDoesContainInvalidQuery_ExpectReturnTrue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arr := make(map[string][]string, 2)
	arr[NODE_ID] = []string{nodeId1}
	arr["invalidQuery"] = []string{"invalid"}

	ret := doesContainInvalidQuery(arr)

	if ret != true {
		t.Errorf("Expected err: %s, actual err: %s", "true", "false")
	}
}

func TestDoesContainInvalidQuery_ExpectReturnFalse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arr := make(map[string][]string, 2)
	arr[NODE_ID] = []string{nodeId1}
	arr[GROUP_ID] = []string{groupId1}
	arr[APP_ID] = []string{appId1}
	arr[IMAGE_NAME] = []string{commonImage}

	ret := doesContainInvalidQuery(arr)

	if ret != false {
		t.Errorf("Expected err: %s, actual err: %s", "false", "true")
	}
}

func TestDoesContain_ExpectReturnTrue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arr := make([]string, 3)
	arr[0] = "one"
	arr[1] = "two"
	arr[2] = "three"

	ret := doesContain(arr, "one")
	if ret != true {
		t.Errorf("Expected err: %s, actual err: %s", "true", "false")
	}
}

func TestDoesContain_ExpectReturnFalse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arr := make([]string, 3)
	arr[0] = "one"
	arr[1] = "two"
	arr[2] = "three"

	ret := doesContain(arr, "four")
	if ret != false {
		t.Errorf("Expected err: %s, actual err: %s", "false", "true")
	}
}
