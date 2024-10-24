package sync

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"
)

func TestBlockHeightExtractNumber(t *testing.T) {
	testStr1 := "sd12345"
	testStr2 := "bh12345"
	testStr3 := "bh12345ds"
	testStr4 := "webh12345ds"
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "test case 1",
			args: args{
				str: testStr1,
			},
			want: 0,
		},
		{
			name: "test case 2",
			args: args{
				str: testStr2,
			},
			want: 12345,
		},
		{
			name: "test case 3",
			args: args{
				str: testStr3,
			},
			want: 0,
		},
		{
			name: "test case 4",
			args: args{
				str: testStr4,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BlockHeightExtractNumber(tt.args.str); got != tt.want {
				t.Errorf("BlockHeightExtractNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterContent(t *testing.T) {
	type args struct {
		originStr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test case 1: Only Chinese characters",
			args: args{originStr: "测试用例"},
			want: "测试用例",
		},
		{
			name: "Test case 2: Only English letters",
			args: args{originStr: "HelloWorld"},
			want: "HelloWorld",
		},
		{
			name: "Test case 3: Only digits",
			args: args{originStr: "1234567890"},
			want: "1234567890",
		},
		{
			name: "Test case 4: Mixed content",
			args: args{originStr: "测试Hello123!@#$%^&*()_+-=[]{}|;':\"<>?,./"},
			want: "测试Hello123_",
		},
		{
			name: "Test case 5: Empty input",
			args: args{originStr: ""},
			want: "",
		},
		{
			name: "Test case 6: Spaces and newlines",
			args: args{originStr: "Hello World\n测试 用例\n123 456"},
			want: "Hello World\n测试 用例\n123 456",
		},
		{
			name: "Test case 7: Special characters",
			args: args{originStr: "!@#$%^&*()_+-=[]{}|;':\"<>?,./"},
			want: "_",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterContent(tt.args.originStr); got != tt.want {
				t.Errorf("FilterContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilteringSensitive(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name            string
		args            args
		wantRetKeyWords []string
		wantFlag        bool
	}{
		{
			name: "test case 1",
			args: args{
				input: "aaaaaa",
			},
			wantFlag: false,
		},
		{
			name: "test case 2",
			args: args{
				input: "诈骗",
			},
			wantFlag: false,
		},
		{
			name: "test case 2",
			args: args{
				input: "诈骗",
			},
			wantFlag: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRetKeyWords, gotFlag := FilteringSensitive(tt.args.input)
			if !reflect.DeepEqual(gotRetKeyWords, tt.wantRetKeyWords) {
				t.Errorf("FilteringSensitive() gotRetKeyWords = %v, want %v", gotRetKeyWords, tt.wantRetKeyWords)
			}
			if gotFlag != tt.wantFlag {
				t.Errorf("FilteringSensitive() gotFlag = %v, want %v", gotFlag, tt.wantFlag)
			}
		})
	}
}

func TestGatewayIdExtractNumber(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Test case 1: Normal case",
			args: args{str: "g00123"},
			want: 123,
		},
		{
			name: "Test case 2: No leading zeros",
			args: args{str: "g123"},
			want: 123,
		},
		{
			name: "Test case 3: Only prefix",
			args: args{str: "g"},
			want: 0,
		},
		{
			name: "Test case 4: Empty input",
			args: args{str: ""},
			want: 0,
		},
		{
			name: "Test case 5: Invalid input",
			args: args{str: "g00abc"},
			want: 0,
		},
		{
			name: "Test case 6: Maximum int64 value",
			args: args{str: "g9223372036854775807"},
			want: 9223372036854775807,
		},
		{
			name: "Test case 7: Overflow int64 value",
			args: args{str: "g9223372036854775808"},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GatewayIdExtractNumber(tt.args.str); got != tt.want {
				t.Errorf("GatewayIdExtractNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMaxHeight(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Test case 1: Normal case",
			args: args{
				chainId: "chain1",
			},
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setMaxHeight("chain1", tt.want)
			if got := GetMaxHeight(tt.args.chainId); got != tt.want {
				t.Errorf("GetMaxHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test case 1: x < y",
			args: args{x: 3, y: 5},
			want: 3,
		},
		{
			name: "Test case 2: x > y",
			args: args{x: 7, y: 2},
			want: 2,
		},
		{
			name: "Test case 3: x == y",
			args: args{x: 4, y: 4},
			want: 4,
		},
		{
			name: "Test case 4: Negative numbers",
			args: args{x: -3, y: -7},
			want: -7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestNewPool(t *testing.T) {
//	type args struct {
//		num int
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "Test case 1: Create pool with positive worker number",
//			args: args{num: 5},
//		},
//		{
//			name: "Test case 2: Create pool with zero worker number",
//			args: args{num: 0},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := NewPool(tt.args.num)
//			if got.workerNum != tt.args.num {
//				t.Errorf("NewPool() workerNum = %v, want %v", got.workerNum, tt.args.num)
//			}
//		})
//	}
//}

//func TestNewTask(t *testing.T) {
//	type args struct {
//		f func() error
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "Test case 1: Create task with function",
//			args: args{f: func() error { return nil }},
//		},
//		{
//			name: "Test case 2: Create task with nil function",
//			args: args{f: nil},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := NewTask(tt.args.f)
//			if got.f == nil && tt.args.f != nil {
//				t.Errorf("NewTask() function should not be nil")
//			}
//			if got.f != nil && tt.args.f == nil {
//				t.Errorf("NewTask() function should be nil")
//			}
//		})
//	}
//}

func TestPool_Run(t *testing.T) {
	type fields struct {
		workerNum  int
		EntryChan  chan *Task
		workerChan chan *Task
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pool{
				workerNum:  tt.fields.workerNum,
				EntryChan:  tt.fields.EntryChan,
				workerChan: tt.fields.workerChan,
			}
			p.Run()
		})
	}
}

func TestPool_worker(t *testing.T) {
	type fields struct {
		workerNum  int
		EntryChan  chan *Task
		workerChan chan *Task
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pool{
				workerNum:  tt.fields.workerNum,
				EntryChan:  tt.fields.EntryChan,
				workerChan: tt.fields.workerChan,
			}
			p.worker()
		})
	}
}

//
//func TestPool_Workflow(t *testing.T) {
//	// Define a counter and a mutex for synchronization
//	var counter int
//	var mu sync.Mutex
//
//	// Define a simple task function that increments the counter
//	taskFunc := func() error {
//		mu.Lock()
//		counter++
//		mu.Unlock()
//		return nil
//	}
//
//	// Create a new pool with 3 workers
//	pool := NewPool(3)
//	pool.Run()
//
//	// Add 10 tasks to the pool
//	for i := 0; i < 10; i++ {
//		task := NewTask(taskFunc)
//		pool.EntryChan <- task
//	}
//
//	// Wait for the tasks to complete
//	time.Sleep(2 * time.Second)
//
//	// Check if the counter value is equal to the number of tasks
//	if counter != 10 {
//		t.Errorf("Pool workflow test failed, expected counter = 10, got %d", counter)
//	}
//}

func Test_setMaxHeight(t *testing.T) {
	type args struct {
		chainId string
		height  int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1: Normal case",
			args: args{
				chainId: "chain1",
				height:  12,
			},
		},
		{
			name: "Test case 2: Normal case",
			args: args{
				chainId: "chain2",
				height:  13,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setMaxHeight(tt.args.chainId, tt.args.height)
		})
	}
}

func Test_filterTxAndEvent(t *testing.T) {
	type args struct {
		transactions map[string]*db.Transaction
		events       []*db.ContractEvent
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1: Sample transactions and events",
			args: args{
				transactions: map[string]*db.Transaction{
					"tx1": {},
					"tx2": {},
				},
				events: []*db.ContractEvent{
					{},
					{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterTxAndEvent(tt.args.transactions, tt.args.events)
		})
	}
}

func TestGatewayIdExtractNumber1(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Test case 1: Normal case",
			args: args{str: "g00123"},
			want: 123,
		},
		{
			name: "Test case 2: No leading zeros",
			args: args{str: "g123"},
			want: 123,
		},
		{
			name: "Test case 3: Only prefix",
			args: args{str: "g"},
			want: 0,
		},
		{
			name: "Test case 4: Empty input",
			args: args{str: ""},
			want: 0,
		},
		{
			name: "Test case 5: Invalid input",
			args: args{str: "g00abc"},
			want: 0,
		},
		{
			name: "Test case 6: Maximum int64 value",
			args: args{str: "g9223372036854775807"},
			want: 9223372036854775807,
		},
		{
			name: "Test case 7: Overflow int64 value",
			args: args{str: "g9223372036854775808"},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GatewayIdExtractNumber(tt.args.str); got != tt.want {
				t.Errorf("GatewayIdExtractNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
