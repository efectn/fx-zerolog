package fxzerolog

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

// Dummy error
var someError = errors.New("some error")

// Test every condition
func TestZeroLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		give        fxevent.Event
		wantMessage string
	}{
		{
			name: "OnStartExecuting",
			give: &fxevent.OnStartExecuting{
				FunctionName: "hook.onStart",
				CallerName:   "bytes.NewBuffer",
			},
			wantMessage: `{"level":"info","callee":"hook.onStart","caller":"bytes.NewBuffer","message":"OnStart hook executing"}`,
		},
		{
			name: "OnStopExecuting",
			give: &fxevent.OnStopExecuting{
				FunctionName: "hook.onStop1",
				CallerName:   "bytes.NewBuffer",
			},
			wantMessage: `{"level":"info","callee":"hook.onStop1","caller":"bytes.NewBuffer","message":"OnStop hook executing"}`,
		},
		{

			name: "OnStopExecutedError",
			give: &fxevent.OnStopExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Err:          fmt.Errorf("some error"),
			},
			wantMessage: `{"level":"error","error":"some error","callee":"hook.onStart1","caller":"bytes.NewBuffer","message":"OnStop hook failed"}`,
		},
		{
			name: "OnStopExecuted",
			give: &fxevent.OnStopExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Runtime:      time.Millisecond * 3,
			},
			wantMessage: `{"level":"info","callee":"hook.onStart1","caller":"bytes.NewBuffer","runtime":"3ms","message":"OnStop hook executed"}`,
		},
		{

			name: "OnStartExecutedError",
			give: &fxevent.OnStartExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Err:          fmt.Errorf("some error"),
			},
			wantMessage: `{"level":"error","error":"some error","callee":"hook.onStart1","caller":"bytes.NewBuffer","message":"OnStart hook failed"}`,
		},
		{
			name: "OnStartExecuted",
			give: &fxevent.OnStartExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Runtime:      time.Millisecond * 3,
			},
			wantMessage: `{"level":"info","callee":"hook.onStart1","caller":"bytes.NewBuffer","runtime":"3ms","message":"OnStart hook executed"}`,
		},
		{
			name:        "Supplied",
			give:        &fxevent.Supplied{TypeName: "*bytes.Buffer"},
			wantMessage: `{"level":"info","type":"*bytes.Buffer","module":"","message":"supplied"}`,
		},
		{
			name:        "SuppliedError",
			give:        &fxevent.Supplied{TypeName: "*bytes.Buffer", Err: someError},
			wantMessage: `{"level":"error","error":"some error","type":"*bytes.Buffer","module":"","message":"supplied"}`,
		},
		{
			name: "Provide",
			give: &fxevent.Provided{
				ConstructorName: "bytes.NewBuffer()",
				ModuleName:      "myModule",
				OutputTypeNames: []string{"*bytes.Buffer"},
			},
			wantMessage: `{"level":"info","constructor":"bytes.NewBuffer()","module":"myModule","type":"*bytes.Buffer","message":"provided"}`,
		},
		{
			name:        "Provide with Error",
			give:        &fxevent.Provided{Err: someError},
			wantMessage: `{"level":"error","error":"some error","module":"","message":"error encountered while applying options"}`,
		},
		{
			name: "Decorate",
			give: &fxevent.Decorated{
				DecoratorName:   "bytes.NewBuffer()",
				ModuleName:      "myModule",
				OutputTypeNames: []string{"*bytes.Buffer"},
			},
			wantMessage: `{"level":"info","decorator":"bytes.NewBuffer()","module":"myModule","type":"*bytes.Buffer","message":"decorated"}`,
		},
		{
			name:        "Decorate with Error",
			give:        &fxevent.Decorated{Err: someError},
			wantMessage: `{"level":"error","error":"some error","module":"","message":"error encountered while applying options"}`,
		},
		{
			name:        "Invoking/Success",
			give:        &fxevent.Invoking{ModuleName: "myModule", FunctionName: "bytes.NewBuffer()"},
			wantMessage: `{"level":"info","function":"bytes.NewBuffer()","module":"myModule","message":"invoking"}`,
		},
		{
			name:        "Invoked/Error",
			give:        &fxevent.Invoked{FunctionName: "bytes.NewBuffer()", Err: someError},
			wantMessage: `{"level":"error","error":"some error","stack":"","function":"bytes.NewBuffer()","message":"invoke failed"}`,
		},
		{
			name:        "StartError",
			give:        &fxevent.Started{Err: someError},
			wantMessage: `{"level":"error","error":"some error","message":"start failed"}`,
		},
		{
			name:        "Stopping",
			give:        &fxevent.Stopping{Signal: os.Interrupt},
			wantMessage: `{"level":"info","signal":"INTERRUPT","message":"received signal"}`,
		},
		{
			name:        "Stopped",
			give:        &fxevent.Stopped{Err: someError},
			wantMessage: `{"level":"error","error":"some error","message":"stop failed"}`,
		},
		{
			name:        "RollingBack",
			give:        &fxevent.RollingBack{StartErr: someError},
			wantMessage: `{"level":"error","error":"some error","message":"start failed, rolling back"}`,
		},
		{
			name:        "RolledBackError",
			give:        &fxevent.RolledBack{Err: someError},
			wantMessage: `{"level":"error","error":"some error","message":"rollback failed"}`,
		},
		{
			name:        "Started",
			give:        &fxevent.Started{},
			wantMessage: `{"level":"info","message":"started"}`,
		},
		{
			name:        "LoggerInitialized Error",
			give:        &fxevent.LoggerInitialized{Err: someError},
			wantMessage: `{"level":"error","error":"some error","message":"custom logger initialization failed"}`,
		},
		{
			name:        "LoggerInitialized",
			give:        &fxevent.LoggerInitialized{ConstructorName: "bytes.NewBuffer()"},
			wantMessage: `{"level":"info","function":"bytes.NewBuffer()","message":"initialized custom fxevent.Logger"}`,
		},
	}

	// Check tests one-by-one
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBufferString("")
			(&ZeroLogger{Logger: zerolog.New(buf)}).LogEvent(tt.give)

			assertEqual(t, tt.wantMessage+"\n", buf.String())
		})
	}
}
