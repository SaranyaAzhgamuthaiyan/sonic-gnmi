package client

import (
    "sync"
    "errors"
	"testing"
	"os"
	"time"
	"reflect"
	"io/ioutil"
	"encoding/json"
	"fmt"

	"github.com/jipanyang/gnxi/utils/xpath"
	"github.com/sonic-net/sonic-gnmi/swsscommon"
	"github.com/sonic-net/sonic-gnmi/test_utils"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

var testFile string = "/etc/sonic/ut.cp.json"

func JsonEqual(a, b []byte) (bool, error) {
	var j1, j2 interface{}
	var err error
	if err = json.Unmarshal(a, &j1); err != nil {
		return false, err
	}
	if err = json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j1, j2), nil
}

func TestJsonClientNegative(t *testing.T) {
	os.Remove(testFile)
	_, err := NewJsonClient(testFile)
	if err == nil {
		t.Errorf("Should fail without checkpoint")
	}

	text := "{"
	err = ioutil.WriteFile(testFile, []byte(text), 0644)
	if err != nil {
		t.Errorf("Fail to create test file")
	}
	_, err = NewJsonClient(testFile)
	if err == nil {
		t.Errorf("Should fail with invalid checkpoint")
	}
}

func TestJsonAdd(t *testing.T) {
	text := "{}"
	err := ioutil.WriteFile(testFile, []byte(text), 0644)
	if err != nil {
		t.Errorf("Fail to create test file")
	}
	client, err := NewJsonClient(testFile)
	if err != nil {
		t.Errorf("Create client fail: %v", err)
	}
	path_list := [][]string {
		[]string {
			"DASH_QOS",
		},
		[]string {
			"DASH_QOS",
			"qos_02",
		},
		[]string {
			"DASH_QOS",
			"qos_03",
			"bw",
		},
		[]string {
			"DASH_VNET",
			"vnet001",
			"address_spaces",
		},
		[]string {
			"DASH_VNET",
			"vnet002",
			"address_spaces",
			"0",
		},
	}
	value_list := []string {
		`{"qos_01": {"bw": "54321", "cps": "1000", "flows": "300"}}`,
		`{"bw": "10001", "cps": "1001", "flows": "101"}`,
		`"20001"`,
		`["10.250.0.0", "192.168.3.0", "139.66.72.9"]`,
		`"6.6.6.6"`,
	}
	for i := 0; i < len(path_list); i++ {
		path := path_list[i]
		value := value_list[i]
		err = client.Add(path, value)
		if err != nil {
			t.Errorf("Add %v fail: %v", path, err)
		}
		res, err := client.Get(path)
		if err != nil {
			t.Errorf("Get %v fail: %v", path, err)
		}
		ok, err := JsonEqual([]byte(value), res)
		if err != nil {
			t.Errorf("Compare json fail: %v", err)
			return
		}
		if ok != true {
			t.Errorf("%v and %v do not match", value, string(res))
		}
	}
	path := []string{}
	res, err := client.Get(path)
	if err != nil {
		t.Errorf("Get %v fail: %v", path, err)
	}
	t.Logf("Result %s", string(res))
}

func TestJsonAddNegative(t *testing.T) {
	text := "{}"
	err := ioutil.WriteFile(testFile, []byte(text), 0644)
	if err != nil {
		t.Errorf("Fail to create test file")
	}
	client, err := NewJsonClient(testFile)
	if err != nil {
		t.Errorf("Create client fail: %v", err)
	}
	path_list := [][]string {
		[]string {
			"DASH_QOS",
		},
		[]string {
			"DASH_QOS",
			"qos_02",
		},
		[]string {
			"DASH_QOS",
			"qos_03",
			"bw",
		},
		[]string {
			"DASH_VNET",
			"vnet001",
			"address_spaces",
		},
		[]string {
			"DASH_VNET",
			"vnet002",
			"address_spaces",
			"0",
		},
		[]string {
			"DASH_VNET",
			"vnet002",
			"address_spaces",
			"abc",
		},
		[]string {
			"DASH_VNET",
			"vnet002",
			"address_spaces",
			"100",
		},
	}
	value_list := []string {
		`{"qos_01": {"bw": "54321", "cps": "1000", "flows": "300"}`,
		`{"bw": "10001", "cps": "1001", "flows": "101"`,
		`20001`,
		`["10.250.0.0", "192.168.3.0", "139.66.72.9"`,
		`"6.6.6.6`,
		`"6.6.6.6"`,
		`"6.6.6.6"`,
	}
	for i := 0; i < len(path_list); i++ {
		path := path_list[i]
		value := value_list[i]
		err = client.Add(path, value)
		if err == nil {
			t.Errorf("Add %v should fail: %v", path, err)
		}
	}
}

func TestJsonRemove(t *testing.T) {
	text := "{}"
	err := ioutil.WriteFile(testFile, []byte(text), 0644)
	if err != nil {
		t.Errorf("Fail to create test file")
	}
	client, err := NewJsonClient(testFile)
	if err != nil {
		t.Errorf("Create client fail: %v", err)
	}
	path_list := [][]string {
		[]string {
			"DASH_QOS",
		},
		[]string {
			"DASH_QOS",
			"qos_02",
		},
		[]string {
			"DASH_QOS",
			"qos_03",
			"bw",
		},
		[]string {
			"DASH_VNET",
			"vnet001",
			"address_spaces",
		},
		[]string {
			"DASH_VNET",
			"vnet002",
			"address_spaces",
			"0",
		},
	}
	value_list := []string {
		`{"qos_01": {"bw": "54321", "cps": "1000", "flows": "300"}}`,
		`{"bw": "10001", "cps": "1001", "flows": "101"}`,
		`"20001"`,
		`["10.250.0.0", "192.168.3.0", "139.66.72.9"]`,
		`"6.6.6.6"`,
	}
	for i := 0; i < len(path_list); i++ {
		path := path_list[i]
		value := value_list[i]
		err = client.Add(path, value)
		if err != nil {
			t.Errorf("Add %v fail: %v", path, err)
		}
		err = client.Remove(path)
		if err != nil {
			t.Errorf("Remove %v fail: %v", path, err)
		}
		_, err := client.Get(path)
		if err == nil {
			t.Errorf("Get %v should fail: %v", path, err)
		}
	}
}

func TestJsonRemoveNegative(t *testing.T) {
	text := "{}"
	err := ioutil.WriteFile(testFile, []byte(text), 0644)
	if err != nil {
		t.Errorf("Fail to create test file")
	}
	client, err := NewJsonClient(testFile)
	if err != nil {
		t.Errorf("Create client fail: %v", err)
	}
	path_list := [][]string {
		[]string {
			"DASH_QOS",
		},
		[]string {
			"DASH_VNET",
			"vnet001",
			"address_spaces",
		},
	}
	value_list := []string {
		`{"qos_01": {"bw": "54321", "cps": "1000", "flows": "300"}}`,
		`["10.250.0.0", "192.168.3.0", "139.66.72.9"]`,
	}
	for i := 0; i < len(path_list); i++ {
		path := path_list[i]
		value := value_list[i]
		err = client.Add(path, value)
		if err != nil {
			t.Errorf("Add %v fail: %v", path, err)
		}
	}

	remove_list := [][]string {
		[]string {
			"DASH_QOS",
			"qos_02",
		},
		[]string {
			"DASH_QOS",
			"qos_03",
			"bw",
		},
		[]string {
			"DASH_VNET",
			"vnet001",
			"address_spaces",
			"abc",
		},
		[]string {
			"DASH_VNET",
			"vnet001",
			"address_spaces",
			"100",
		},
	}
	for i := 0; i < len(remove_list); i++ {
		path := remove_list[i]
		err = client.Remove(path)
		if err == nil {
			t.Errorf("Remove %v should fail: %v", path, err)
		}
	}
}

func TestParseDatabase(t *testing.T) {
	var test_paths []*gnmipb.Path
	var prefix *gnmipb.Path
	var err error

	client := MixedDbClient {
		namespace_cnt : 1,
		container_cnt : 1,
	}
	_, _, err = client.ParseDatabase(prefix, test_paths)
	if err == nil {
		t.Errorf("ParseDatabase should fail for empty path: %v", err)
	}

	test_target := "TEST_DB"
	path, err := xpath.ToGNMIPath("sonic-db:" + test_target + "/localhost" + "/VLAN")
	test_paths = append(test_paths, path)
	target, dbkey1, err := client.ParseDatabase(prefix, test_paths)
	if err != nil {
		t.Errorf("ParseDatabase failed to get target: %v", err)
	}
	defer swsscommon.DeleteSonicDBKey(dbkey1)
	if target != test_target {
		t.Errorf("ParseDatabase return wrong target: %v", target)
	}

	// Smartswitch with multiple asic NPU
	client = MixedDbClient {
		namespace_cnt : 2,
		container_cnt : 2,
	}

	test_target = "TEST_DB"
	path, err = xpath.ToGNMIPath("sonic-db:" + test_target + "/localhost" + "/VLAN")
	test_paths = append(test_paths, path)
	target, dbkey2, err := client.ParseDatabase(prefix, test_paths)
	if err != nil {
		t.Errorf("ParseDatabase failed to get target: %v", err)
	}
	defer swsscommon.DeleteSonicDBKey(dbkey2)
	if target != test_target {
		t.Errorf("ParseDatabase return wrong target: %v", target)
	}

	test_target = "TEST_DB"
	path, err = xpath.ToGNMIPath("sonic-db:" + test_target + "/xyz" + "/VLAN")
	test_paths = append(test_paths, path)
	target, _, err = client.ParseDatabase(prefix, test_paths)
	if err == nil {
		t.Errorf("ParseDatabase should fail for namespace/container")
	}
}

func mockGetFunc() ([]byte, error) {
	return nil, errors.New("mock error")
}

func TestNonDbClientGetError(t *testing.T) {
	var gnmipbPath *gnmipb.Path = &gnmipb.Path{
		Element: []string{"mockPath"},
	}

	path2Getter := map[*gnmipb.Path]dataGetFunc{
		gnmipbPath: mockGetFunc,
	}

	// Create a NonDbClient with the mocked dataGetFunc
	client := &NonDbClient{
		path2Getter: path2Getter,
	}

	var w *sync.WaitGroup
	_, err := client.Get(w)
	if errors.Is(err, errors.New("mock error")) {
		t.Errorf("Expected error from NonDbClient.Get, got nil")
	}
}

/*
	Helper method for receive data from ZmqConsumerStateTable
		consumer: Receive data from consumer
		return:
			true: data received
			false: not receive any data after retry
*/
func ReceiveFromZmq(consumer swsscommon.ZmqConsumerStateTable) (bool) {
	receivedData := swsscommon.NewKeyOpFieldsValuesQueue()
	defer swsscommon.DeleteKeyOpFieldsValuesQueue(receivedData)
	retry := 0;
	for {
		// sender's ZMQ may disconnect, wait and retry for reconnect 
		time.Sleep(time.Duration(1000) * time.Millisecond)
		consumer.Pops(receivedData)
		if receivedData.Size() == 0 {
			retry++
			if retry >= 10 {
				return false
			}
		} else {
			return true
		}
	}
}

func TestZmqReconnect(t *testing.T) {
	// create ZMQ server
	db := swsscommon.NewDBConnector(APPL_DB_NAME, SWSS_TIMEOUT, false)
	zmqServer := swsscommon.NewZmqServer("tcp://*:1234")
	var TEST_TABLE string = "DASH_ROUTE"
	consumer := swsscommon.NewZmqConsumerStateTable(db, TEST_TABLE, zmqServer)

	// create ZMQ client side
	zmqAddress := "tcp://127.0.0.1:1234"
	client := MixedDbClient {
		applDB : swsscommon.NewDBConnector(APPL_DB_NAME, SWSS_TIMEOUT, false),
		tableMap : map[string]swsscommon.ProducerStateTable{},
		zmqTableMap : map[string]swsscommon.ZmqProducerStateTable{},
		zmqClient : swsscommon.NewZmqClient(zmqAddress),
	}

	data := map[string]string{}
	var TEST_KEY string = "TestKey"
	client.DbSetTable(TEST_TABLE, TEST_KEY, data)
	if !ReceiveFromZmq(consumer) {
		t.Errorf("Receive data from ZMQ failed")
	}

	// recreate ZMQ server to trigger re-connect
	swsscommon.DeleteZmqConsumerStateTable(consumer)
	swsscommon.DeleteZmqServer(zmqServer)
	zmqServer = swsscommon.NewZmqServer("tcp://*:1234")
	consumer = swsscommon.NewZmqConsumerStateTable(db, TEST_TABLE, zmqServer)

	// send data again, client will reconnect
	client.DbSetTable(TEST_TABLE, TEST_KEY, data)
	if !ReceiveFromZmq(consumer) {
		t.Errorf("Receive data from ZMQ failed")
	}

	client.Close()
	swsscommon.DeleteZmqConsumerStateTable(consumer)
	swsscommon.DeleteZmqClient(client.zmqClient)
	swsscommon.DeleteZmqServer(zmqServer)
	swsscommon.DeleteDBConnector(db)

	for _, client := range zmqClientMap {
		swsscommon.DeleteZmqClient(client)
	}
}

func TestRetryHelper(t *testing.T) {
	// create ZMQ server
	zmqServer := swsscommon.NewZmqServer("tcp://*:2234")

	// create ZMQ client side
	zmqAddress := "tcp://127.0.0.1:2234"
	zmqClient := swsscommon.NewZmqClient(zmqAddress)
	returnError := true
	exeCount := 0
    RetryHelper(
		zmqClient,
		func () (err error) {
			exeCount++
			if returnError {
				returnError = false
				return fmt.Errorf("connection_reset")
			}
			return nil
	})

	if exeCount == 1 {
		t.Errorf("RetryHelper does not retry")
	}

	if exeCount > 2 {
		t.Errorf("RetryHelper retry too much")
	}

	swsscommon.DeleteZmqClient(zmqClient)
	swsscommon.DeleteZmqServer(zmqServer)
}

func TestGetDpuAddress(t *testing.T) {
	// prepare data according to design doc
	// Design doc: https://github.com/sonic-net/SONiC/blob/master/doc/smart-switch/ip-address-assigment/smart-switch-ip-address-assignment.md?plain=1

	if !swsscommon.SonicDBConfigIsInit() {
		swsscommon.SonicDBConfigInitialize()
	}

	var configDb = swsscommon.NewDBConnector("CONFIG_DB", uint(0), true)
	configDb.Flushdb()
	
	var midPlaneTable = swsscommon.NewTable(configDb, "MID_PLANE_BRIDGE")
	var dpusTable = swsscommon.NewTable(configDb, "DPUS")
	var dhcpPortTable = swsscommon.NewTable(configDb, "DHCP_SERVER_IPV4_PORT")

	// test get DPU address when database not ready
	address, err := getDpuAddress("dpu0")
	if err == nil {
		t.Errorf("get DPU address should failed: %v, but get %s", err, address)
	}

	midPlaneTable.Hset("GLOBAL", "bridge", "bridge_midplane")
	dpusTable.Hset("dpu0", "midplane_interface", "dpu0")

	// test get DPU address when DHCP_SERVER_IPV4_PORT table not ready
	address, err = getDpuAddress("dpu0")
	if err == nil {
		t.Errorf("get DPU address should failed: %v, but get %s", err, address)
	}

	dhcpPortTable.Hset("bridge_midplane|dpu0", "invalidfield", "")

	// test get DPU address when DHCP_SERVER_IPV4_PORT table broken
	address, err = getDpuAddress("dpu0")
	if err == nil {
		t.Errorf("get DPU address should failed: %v, but get %s", err, address)
	}

	dhcpPortTable.Hset("bridge_midplane|dpu0", "ips", "127.0.0.2,127.0.0.1")

	// test get valid DPU address
	address, err = getDpuAddress("dpu0")
	if err != nil {
		t.Errorf("get DPU address failed: %v", err)
	}

	if address != "127.0.0.2" {
		t.Errorf("get DPU address failed: %v", address)
	}

	// test get invalid DPU address
	address, err = getDpuAddress("dpu_x")
	if err == nil {
		t.Errorf("get invalid DPU address failed")
	}

	if address != "" {
		t.Errorf("get invalid DPU address failed: %v", address)
	}

	// test get ZMQ address
	address, err = getZmqAddress("dpu0", "1234")
	if address != "tcp://127.0.0.2:1234" {
		t.Errorf("get invalid DPU address failed")
	}

	address, err = getZmqAddress("dpu0", "")
	if err == nil {
		t.Errorf("get invalid ZMQ address failed")
	}

	address, err = getZmqAddress("", "1234")
	if err == nil {
		t.Errorf("get invalid ZMQ address failed")
	}
	
	swsscommon.DeleteTable(midPlaneTable)
	swsscommon.DeleteTable(dpusTable)
	swsscommon.DeleteTable(dhcpPortTable)
	swsscommon.DeleteDBConnector(configDb)
}

func TestGetZmqClient(t *testing.T) {
	if !swsscommon.SonicDBConfigIsInit() {
		swsscommon.SonicDBConfigInitialize()
	}

	var configDb = swsscommon.NewDBConnector("CONFIG_DB", uint(0), true)
	configDb.Flushdb()

	var midPlaneTable = swsscommon.NewTable(configDb, "MID_PLANE_BRIDGE")
	var dpusTable = swsscommon.NewTable(configDb, "DPUS")
	var dhcpPortTable = swsscommon.NewTable(configDb, "DHCP_SERVER_IPV4_PORT")

	midPlaneTable.Hset("GLOBAL", "bridge", "bridge_midplane")
	dpusTable.Hset("dpu0", "midplane_interface", "dpu0")
	dhcpPortTable.Hset("bridge_midplane|dpu0", "ips", "127.0.0.2,127.0.0.1")

	client, err := getZmqClient("dpu0", "")
	if client != nil || err != nil {
		t.Errorf("empty ZMQ port should not get ZMQ client")
	}

	client, err = getZmqClient("dpu0", "1234")
	if client == nil {
		t.Errorf("get ZMQ client failed")
	}

	client, err = getZmqClient("", "1234")
	if client == nil {
		t.Errorf("get ZMQ client failed")
	}

	err = removeZmqClient(client)
	if err != nil {
		t.Errorf("Remove ZMQ client failed")
	}

	// Remove a removed client should failed
	err = removeZmqClient(client)
	if err == nil {
		t.Errorf("Remove ZMQ client should failed")
	}
	
	swsscommon.DeleteTable(midPlaneTable)
	swsscommon.DeleteTable(dpusTable)
	swsscommon.DeleteTable(dhcpPortTable)
	swsscommon.DeleteDBConnector(configDb)

	for _, client := range zmqClientMap {
		swsscommon.DeleteZmqClient(client)
	}
}

func TestMain(m *testing.M) {
	defer test_utils.MemLeakCheck()
	m.Run()
}
